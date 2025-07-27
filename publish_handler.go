package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	srt "github.com/datarhei/gosrt"
	"github.com/datarhei/joy4/av"
	"github.com/datarhei/joy4/format/flv"
	"github.com/datarhei/joy4/format/rtmp"
	"github.com/datarhei/joy4/format/ts"
)

// OutputManager управляет выходами для одного входа
// Позволяет добавлять/удалять выходы "на лету"
type OutputManager struct {
	mu      sync.RWMutex
	outputs map[string]*outputWorker // url → worker
}

type outputWorker struct {
	ch   chan av.Packet
	stop chan struct{}
}

func NewOutputManager() *OutputManager {
	return &OutputManager{
		outputs: make(map[string]*outputWorker),
	}
}

// Добавить выход (если уже есть — ничего не делает)
func (om *OutputManager) AddOutput(url string, bufSize int, startPush func(ch <-chan av.Packet, stop <-chan struct{})) {
	om.mu.Lock()
	defer om.mu.Unlock()
	if _, exists := om.outputs[url]; exists {
		return
	}
	ch := make(chan av.Packet, bufSize)
	stop := make(chan struct{})
	om.outputs[url] = &outputWorker{ch: ch, stop: stop}
	go startPush(ch, stop)
}

// Удалить выход (если нет — ничего не делает)
func (om *OutputManager) RemoveOutput(url string) {
	om.mu.Lock()
	defer om.mu.Unlock()
	if w, ok := om.outputs[url]; ok {
		close(w.stop)
		delete(om.outputs, url)
	}
}

// Получить все выходы (копия map)
func (om *OutputManager) AllOutputs() map[string]*outputWorker {
	om.mu.RLock()
	defer om.mu.RUnlock()
	res := make(map[string]*outputWorker)
	for k, v := range om.outputs {
		res[k] = v
	}
	return res
}

func handlePublish(sm *StreamManager, cfg *Config) func(conn *rtmp.Conn) {
	return func(srcConn *rtmp.Conn) {
		// Добавляем обработку паники
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] Publish handler panic: %v", r)
			}
			log.Printf("[DEBUG] Publish handler finished for: %s", srcConn.URL)
		}()

		log.Printf("Publish started: %s", srcConn.URL)

		inputCfg := sm.GetInputByPath(srcConn.URL.Path)
		if inputCfg == nil {
			log.Printf("Unknown input URL path: %s. Rejecting publish.", srcConn.URL.Path)
			srcConn.Close()
			return
		}
		log.Printf("Matched input config: %s with %d outputs", inputCfg.Name, len(inputCfg.Outputs))

		log.Printf("[DEBUG] Starting publish handling for input: %s", inputCfg.Name)
		log.Printf("[DEBUG] Output URLs: %v", inputCfg.Outputs)

		sm.SetStatusActive(inputCfg.Name, true)
		defer sm.SetStatusActive(inputCfg.Name, false)

		streams, err := srcConn.Streams()
		if err != nil {
			log.Printf("Failed to get streams from source: %v", err)
			return
		}

		stopChan := make(chan struct{})
		outputMgr := NewOutputManager()
		bufSize := 5000 // Увеличили с 3000 до 5000 для лучшей устойчивости

		// Heartbeat для отслеживания состояния трансляции
		heartbeatCtx, heartbeatCancel := context.WithCancel(context.Background())
		defer heartbeatCancel()

		go func() {
			ticker := time.NewTicker(2 * time.Minute)
			defer ticker.Stop()

			startTime := time.Now()
			log.Printf("[HEARTBEAT] Publish started for '%s' at %v", inputCfg.Name, startTime)

			for {
				select {
				case <-heartbeatCtx.Done():
					return
				case <-ticker.C:
					uptime := time.Since(startTime)
					outputs := outputMgr.AllOutputs()
					activeOutputs := 0
					for range outputs {
						activeOutputs++
					}
					log.Printf("[HEARTBEAT] Publish '%s' uptime: %v, Active outputs: %d/%d",
						inputCfg.Name, uptime, activeOutputs, len(inputCfg.Outputs))

					// Проверяем состояние соединений
					for url, w := range outputs {
						bufferSize := len(w.ch)
						// Логируем только если буфер заполнен больше чем на 50%
						if bufferSize > bufSize/2 {
							log.Printf("[WARNING] Output %s buffer filling up: %d/%d", url, bufferSize, bufSize)
						}
					}
				}
			}
		}()

		// Функция для старта push-горутины для выхода
		startPush := func(url string) func(<-chan av.Packet, <-chan struct{}) {
			return func(ch <-chan av.Packet, stop <-chan struct{}) {
				// Добавляем обработку паники для output горутин
				defer func() {
					if r := recover(); r != nil {
						log.Printf("[PANIC] Output goroutine panic for %s: %v", url, r)
					}
					log.Printf("[DEBUG] Output goroutine finished for: %s", url)
				}()

				var totalBytes int64
				var packetCount int64
				var lastLogTime time.Time

				for {
					select {
					case <-stop:
						log.Printf("Output stopped: %s", url)
						return
					default:
					}

					// Убираем спам логи подключения
					// log.Printf("[DEBUG] Attempting to connect to output URL: %s", url)

					if strings.HasPrefix(url, "file://") {
						log.Printf("[DEBUG] Detected file output for URL: %s", url)
						filename := strings.TrimPrefix(url, "file://")
						log.Printf("[DEBUG] Creating file: %s", filename)
						file, err := os.Create(filename)
						if err != nil {
							log.Printf("[ERROR] Failed to create file %s: %v", filename, err)
							time.Sleep(5 * time.Second)
							return
						}

						log.Printf("[DEBUG] Writing FLV to file: %s", filename)
						sm.SetOutputActive(inputCfg.Name, url, true)
						muxer := flv.NewMuxer(file)
						err = muxer.WriteHeader(streams)
						if err != nil {
							log.Printf("[ERROR] Failed to write FLV header to file %s: %v", filename, err)
							file.Close()
							sm.SetOutputActive(inputCfg.Name, url, false)
							time.Sleep(5 * time.Second)
							return
						}

						fileDone := false
						for !fileDone {
							select {
							case <-stop:
								if err := muxer.WriteTrailer(); err != nil {
									log.Printf("Failed to write FLV trailer to file %s: %v", filename, err)
								}
								file.Close()
								sm.SetOutputActive(inputCfg.Name, url, false)
								log.Printf("File output stopped: %s", filename)
								fileDone = true
							case pkt, ok := <-ch:
								if !ok {
									if err := muxer.WriteTrailer(); err != nil {
										log.Printf("Failed to write FLV trailer to file %s: %v", filename, err)
									}
									file.Close()
									sm.SetOutputActive(inputCfg.Name, url, false)
									fileDone = true
									break
								}
								err = muxer.WritePacket(pkt)
								if err != nil {
									log.Printf("Write error to file %s: %v", filename, err)
									if err := muxer.WriteTrailer(); err != nil {
										log.Printf("Failed to write FLV trailer to file %s: %v", filename, err)
									}
									file.Close()
									sm.SetOutputActive(inputCfg.Name, url, false)
									time.Sleep(5 * time.Second)
									fileDone = true
									break
								}
								totalBytes += int64(len(pkt.Data))
								sm.UpdateOutputBitrate(inputCfg.Name, url, totalBytes)
							}
						}
						return
					} else if strings.HasPrefix(url, "rtmp://") {
						dstConn, err := rtmp.Dial(url, rtmp.DialOptions{})
						if err != nil {
							log.Printf("Failed to connect to %s: %v", url, err)
							// Получаем актуальный интервал переподключения
							sm.mu.RLock()
							reconnectInterval := sm.config.ReconnectInterval
							sm.mu.RUnlock()
							time.Sleep(time.Duration(reconnectInterval) * time.Second)
							continue
						}
						sm.SetOutputActive(inputCfg.Name, url, true)
						err = dstConn.WriteHeader(streams)
						if err != nil {
							log.Printf("Failed to write header to %s: %v", url, err)
							dstConn.Close()
							sm.SetOutputActive(inputCfg.Name, url, false)
							// Получаем актуальный интервал переподключения
							sm.mu.RLock()
							reconnectInterval := sm.config.ReconnectInterval
							sm.mu.RUnlock()
							time.Sleep(time.Duration(reconnectInterval) * time.Second)
							continue
						}
					rtmpLoop:
						for {
							select {
							case <-stop:
								dstConn.Close()
								sm.SetOutputActive(inputCfg.Name, url, false)
								return
							case pkt, ok := <-ch:
								if !ok {
									dstConn.Close()
									sm.SetOutputActive(inputCfg.Name, url, false)
									return
								}
								err = dstConn.WritePacket(pkt)
								if err != nil {
									log.Printf("Write error to %s: %v", url, err)
									dstConn.Close()
									sm.SetOutputActive(inputCfg.Name, url, false)
									// Получаем актуальный интервал переподключения
									sm.mu.RLock()
									reconnectInterval := sm.config.ReconnectInterval
									sm.mu.RUnlock()
									time.Sleep(time.Duration(reconnectInterval) * time.Second)
									break rtmpLoop
								}
								totalBytes += int64(len(pkt.Data))
								sm.UpdateOutputBitrate(inputCfg.Name, url, totalBytes)
							}
						}
					} else if strings.HasPrefix(url, "srt://") {
						srtAddr := strings.TrimPrefix(url, "srt://")
						if idx := strings.Index(srtAddr, "?"); idx != -1 {
							srtAddr = srtAddr[:idx]
						}
						cfgSRT := srt.DefaultConfig()

						// Получаем актуальные настройки SRT из StreamManager
						sm.mu.RLock()
						srtSettings := sm.config.SRTSettings
						reconnectInterval := sm.config.ReconnectInterval
						sm.mu.RUnlock()

						if srtSettings.Latency > 0 {
							cfgSRT.Latency = time.Duration(srtSettings.Latency) * time.Millisecond
						}
						if srtSettings.Passphrase != "" {
							cfgSRT.Passphrase = srtSettings.Passphrase
						}
						if srtSettings.StreamID != "" {
							cfgSRT.StreamId = srtSettings.StreamID
						}
						if srtSettings.ConnectTimeout > 0 {
							cfgSRT.ConnectionTimeout = time.Duration(srtSettings.ConnectTimeout) * time.Millisecond
						}
						log.Printf("SRT connecting to %s with latency=%v, streamid=%s, timeout=%v", srtAddr, cfgSRT.Latency, cfgSRT.StreamId, cfgSRT.ConnectionTimeout)
						conn, err := srt.Dial("srt", srtAddr, cfgSRT)
						if err != nil {
							log.Printf("Failed to connect to SRT %s: %v", url, err)
							time.Sleep(time.Duration(reconnectInterval) * time.Second)
							continue
						}

						// Проверяем состояние соединения
						sm.SetOutputActive(inputCfg.Name, url, true)

						// Правильная обработка временных меток для SRT
						var baseTime time.Duration
						var baseTimeSet bool
						var lastVideoTime time.Duration
						var lastAudioTime time.Duration
						// Убираем ресинхронизацию - она ломает кодировщик
						// const resyncInterval = 10 * time.Minute
						// const maxNormalizedTime = 30 * time.Minute

						// Используем буферизованный подход для стабильности
						var tsBuf bytes.Buffer
						muxer := ts.NewMuxer(&tsBuf)
						err = muxer.WriteHeader(streams)
						if err != nil {
							log.Printf("TS WriteHeader error for %s: %v", url, err)
							conn.Close()
							sm.SetOutputActive(inputCfg.Name, url, false)
							time.Sleep(time.Duration(reconnectInterval) * time.Second)
							continue
						}

						// Отправляем заголовок TS сразу
						headerData := tsBuf.Bytes()
						if len(headerData) > 0 {
							_, err = conn.Write(headerData)
							if err != nil {
								log.Printf("SRT Write header error for %s: %v", url, err)
								conn.Close()
								sm.SetOutputActive(inputCfg.Name, url, false)
								time.Sleep(time.Duration(reconnectInterval) * time.Second)
								continue
							}
							totalBytes += int64(len(headerData))
							tsBuf.Reset()
						}

					srtLoop:
						for {
							select {
							case <-stop:
								conn.Close()
								sm.SetOutputActive(inputCfg.Name, url, false)
								return
							case pkt, ok := <-ch:
								if !ok {
									conn.Close()
									sm.SetOutputActive(inputCfg.Name, url, false)
									return
								}

								// Обработка временных меток для SRT
								if !baseTimeSet {
									// Ждем первый ключевой кадр для установки базового времени
									if pkt.IsKeyFrame {
										// Не устанавливаем baseTime если временная метка уже 0
										if pkt.Time > 0 {
											baseTime = pkt.Time
											log.Printf("SRT base time set to: %v", baseTime)
										} else {
											log.Printf("SRT first keyframe has zero timestamp, skipping base time setup")
										}
										baseTimeSet = true
									} else {
										// Пропускаем пакеты до первого ключевого кадра
										continue
									}
								}

								// Простая нормализация без ресинхронизации
								normalizedTime := pkt.Time - baseTime
								if normalizedTime < 0 {
									normalizedTime = 0
								}

								// Логируем детали нормализации для отладки
								//if pkt.IsKeyFrame {
								//	log.Printf("[SRT] Keyframe - Original: %v, Base: %v, Normalized: %v", pkt.Time, baseTime, normalizedTime)
								//}

								// Создаем новый пакет с нормализованной временной меткой
								normalizedPkt := av.Packet{
									Data:            pkt.Data,
									Time:            normalizedTime,
									Idx:             pkt.Idx,
									IsKeyFrame:      pkt.IsKeyFrame,
									CompositionTime: pkt.CompositionTime,
								}

								// Проверяем и корректируем монотонность временных меток
								if pkt.Idx == 0 { // Video
									if normalizedTime < lastVideoTime {
										// Если таймстемп меньше предыдущего, плавно увеличиваем его,
										// чтобы избежать дублирования и рывков.
										normalizedTime = lastVideoTime + time.Millisecond
									}
									lastVideoTime = normalizedTime
									normalizedPkt.Time = normalizedTime
								} else { // Audio
									if normalizedTime < lastAudioTime {
										// То же самое для аудио
										normalizedTime = lastAudioTime + time.Millisecond
									}
									lastAudioTime = normalizedTime
									normalizedPkt.Time = normalizedTime
								}

								// Записываем пакет в TS муксер
								err = muxer.WritePacket(normalizedPkt)
								if err != nil {
									log.Printf("TS WritePacket error for %s: %v", url, err)
									conn.Close()
									sm.SetOutputActive(inputCfg.Name, url, false)
									time.Sleep(time.Duration(reconnectInterval) * time.Second)
									break srtLoop
								}

								// Логируем обработку пакетов для диагностики блокировки
								//if pkt.IsKeyFrame {
								//	log.Printf("[DEBUG] Processing keyframe for %s", url)
								//}

								// Подсчитываем статистику
								packetCount++
								if time.Since(lastLogTime) > 30*time.Second {
									log.Printf("[DEBUG] Output %s processed %d packets in 30s", url, packetCount)
									packetCount = 0
									lastLogTime = time.Now()
								}

								// Отправляем все данные из буфера с таймаутом
								tsData := tsBuf.Bytes()
								if len(tsData) > 0 {
									// Неблокирующая запись с таймаутом
									writeDone := make(chan error, 1)

									go func() {
										conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
										_, err := conn.Write(tsData)
										writeDone <- err
									}()

									// Ждем результат с таймаутом
									select {
									case err := <-writeDone:
										if err != nil {
											log.Printf("SRT Write error for %s: %v", url, err)

											// Детальная диагностика ошибки
											if strings.Contains(err.Error(), "timeout") {
												log.Printf("[ERROR] SRT write timeout for %s - connection may be slow", url)
											} else if strings.Contains(err.Error(), "connection") {
												log.Printf("[ERROR] SRT connection lost for %s", url)
											} else if strings.Contains(err.Error(), "broken") {
												log.Printf("[ERROR] SRT pipe broken for %s", url)
											} else {
												log.Printf("[ERROR] SRT unknown error for %s: %T", url, err)
											}

											conn.Close()
											sm.SetOutputActive(inputCfg.Name, url, false)
											time.Sleep(time.Duration(reconnectInterval) * time.Second)
											break srtLoop
										}

										totalBytes += int64(len(tsData))
										sm.UpdateOutputBitrate(inputCfg.Name, url, totalBytes)
										tsBuf.Reset()

										// Небольшая задержка для стабильности SRT - убираем, т.к. вносит ненужную задержку
										// time.Sleep(1 * time.Millisecond)

									case <-time.After(3 * time.Second):
										// Таймаут записи - закрываем соединение
										log.Printf("[ERROR] SRT write timeout for %s - forcing reconnect", url)
										conn.Close()
										sm.SetOutputActive(inputCfg.Name, url, false)
										time.Sleep(time.Duration(reconnectInterval) * time.Second)
										break srtLoop
									}
								} else {
									// Если буфер пустой, ничего не делаем, просто ждем следующий пакет.
									// Задержка здесь приводила к замедлению всего пайплайна.
									// time.Sleep(10 * time.Millisecond)
								}
							}
						}
					}
				}
			}
		}

		// Инициализация выходов из inputCfg
		for _, url := range inputCfg.Outputs {
			sm.RegisterOutput(inputCfg.Name, url)
			outputMgr.AddOutput(url, bufSize, startPush(url))
		}

		// Горутина для динамического обновления выходов
		updateTicker := time.NewTicker(2 * time.Second)
		defer updateTicker.Stop()

		// Горутина мониторинга буферов
		bufferMonitorTicker := time.NewTicker(30 * time.Second)
		defer bufferMonitorTicker.Stop()

		go func() {
			for {
				select {
				case <-stopChan:
					return
				case <-updateTicker.C:
					// Синхронизируем выходы с inputCfg.Outputs
					current := make(map[string]struct{})
					for _, url := range inputCfg.Outputs {
						current[url] = struct{}{}
						sm.RegisterOutput(inputCfg.Name, url)
						outputMgr.AddOutput(url, bufSize, startPush(url))
					}
					for url := range outputMgr.AllOutputs() {
						if _, ok := current[url]; !ok {
							outputMgr.RemoveOutput(url)
						}
					}
				case <-bufferMonitorTicker.C:
					// Мониторинг состояния буферов каждые 30 секунд
					outputs := outputMgr.AllOutputs()
					for url, w := range outputs {
						bufferSize := len(w.ch)
						fillPercentage := (bufferSize * 100) / bufSize

						if bufferSize > bufSize*3/4 {
							log.Printf("[WARNING] Output buffer for %s is %d%% full: %d/%d packets", url, fillPercentage, bufferSize, bufSize)
						} else {
							log.Printf("[MONITOR] Output buffer for %s is %d%% full: %d/%d packets", url, fillPercentage, bufferSize, bufSize)
						}
					}
				}
			}
		}()

		// Горутина чтения пакетов из источника и рассылки по выходам
		for {
			pkt, err := srcConn.ReadPacket()
			if err != nil {
				log.Printf("Source disconnected or error: %v", err)
				close(stopChan)
				break
			}

			// Неблокирующая отправка пакетов в выходы
			outputs := outputMgr.AllOutputs()
			droppedCount := 0

			for _, w := range outputs {
				select {
				case w.ch <- pkt:
					// Успешно отправлен
				default:
					// буфер заполнен — дропаем пакет для этого выхода
					droppedCount++
					log.Printf("[WARNING] Output buffer full, dropping packet for output")

					// Логируем размер буфера для диагностики
					log.Printf("[DEBUG] Channel buffer size: %d", len(w.ch))
				}
			}

			// Если дропнули пакеты для всех выходов, это серьезная проблема
			if droppedCount == len(outputs) && len(outputs) > 0 {
				log.Printf("[ERROR] All outputs are blocked! Dropped packet for all %d outputs", len(outputs))

				// Логируем состояние всех выходов для диагностики
				for url, w := range outputs {
					bufferSize := len(w.ch)
					log.Printf("[DEBUG] Output %s buffer state: %d/%d packets", url, bufferSize, bufSize)
				}

				// Не очищаем буферы - это может сломать выходы
				// Вместо этого просто логируем проблему
				log.Printf("[WARNING] All outputs are slow, packets will be dropped until they catch up")
			}
		}
		log.Printf("Publish finished: %s", srcConn.URL)

		// Сбрасываем статус всех выходов при завершении трансляции
		for _, url := range inputCfg.Outputs {
			sm.SetOutputActive(inputCfg.Name, url, false)
		}
		// ЯВНО закрываем все выходные каналы, чтобы завершились все горутины записи
		for _, w := range outputMgr.AllOutputs() {
			close(w.ch)
		}
		log.Printf("[DEBUG] All output channels closed for: %s", srcConn.URL)
	}
}

package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	srt "github.com/datarhei/gosrt"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/format/flv"
	"github.com/nareix/joy4/format/rtmp"
	"github.com/nareix/joy4/format/ts"
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
		log.Printf("Publish started: %s", srcConn.URL)

		inputCfg := sm.GetInputByPath(srcConn.URL.Path)
		if inputCfg == nil {
			log.Printf("Unknown input URL path: %s. Rejecting publish.", srcConn.URL.Path)
			srcConn.Close()
			return
		}
		log.Printf("Matched input config: %s with %d outputs", inputCfg.Name, len(inputCfg.Outputs))

		sm.SetStatusActive(inputCfg.Name, true)
		defer sm.SetStatusActive(inputCfg.Name, false)

		streams, err := srcConn.Streams()
		if err != nil {
			log.Printf("Failed to get streams from source: %v", err)
			return
		}

		stopChan := make(chan struct{})
		outputMgr := NewOutputManager()
		bufSize := 500

		// Функция для старта push-горутины для выхода
		startPush := func(url string) func(<-chan av.Packet, <-chan struct{}) {
			return func(ch <-chan av.Packet, stop <-chan struct{}) {
				var totalBytes int64
				for {
					select {
					case <-stop:
						log.Printf("Output stopped: %s", url)
						return
					default:
					}

					log.Printf("Trying to connect to %s", url)

					if strings.HasPrefix(url, "file://") {
						filename := strings.TrimPrefix(url, "file://")
						file, err := os.Create(filename)
						if err != nil {
							log.Printf("Failed to create file %s: %v", filename, err)
							time.Sleep(5 * time.Second)
							continue // Повторная попытка
						}

						sm.SetOutputActive(inputCfg.Name, url, true)
						log.Printf("Writing FLV to file: %s", filename)

						muxer := flv.NewMuxer(file)
						err = muxer.WriteHeader(streams)
						if err != nil {
							log.Printf("Failed to write FLV header to file %s: %v", filename, err)
							file.Close()
							sm.SetOutputActive(inputCfg.Name, url, false)
							time.Sleep(5 * time.Second)
							continue
						}

						for {
							select {
							case <-stop:
								if err := muxer.WriteTrailer(); err != nil {
									log.Printf("Failed to write FLV trailer to file %s: %v", filename, err)
								}
								file.Close()
								sm.SetOutputActive(inputCfg.Name, url, false)
								log.Printf("File output stopped: %s", filename)
								return
							case pkt, ok := <-ch:
								if !ok {
									if err := muxer.WriteTrailer(); err != nil {
										log.Printf("Failed to write FLV trailer to file %s: %v", filename, err)
									}
									file.Close()
									sm.SetOutputActive(inputCfg.Name, url, false)
									return
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
									break // Выход из внутреннего цикла для переподключения
								}
								totalBytes += int64(len(pkt.Data))
								sm.UpdateOutputBitrate(inputCfg.Name, url, totalBytes)
							}
						}
					} else if strings.HasPrefix(url, "rtmp://") {
						dstConn, err := rtmp.Dial(url)
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
									break
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
						sm.SetOutputActive(inputCfg.Name, url, true)
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
								tsBuf.Reset()
								err = muxer.WritePacket(pkt)
								if err != nil {
									log.Printf("TS WritePacket error for %s: %v", url, err)
									conn.Close()
									sm.SetOutputActive(inputCfg.Name, url, false)
									time.Sleep(time.Duration(reconnectInterval) * time.Second)
									break
								}
								if tsBuf.Len() > 0 {
									_, err = conn.Write(tsBuf.Bytes())
									if err != nil {
										log.Printf("SRT Write error for %s: %v", url, err)
										conn.Close()
										sm.SetOutputActive(inputCfg.Name, url, false)
										time.Sleep(time.Duration(reconnectInterval) * time.Second)
										break
									}
									totalBytes += int64(tsBuf.Len())
									sm.UpdateOutputBitrate(inputCfg.Name, url, totalBytes)
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
			for _, w := range outputMgr.AllOutputs() {
				select {
				case w.ch <- pkt:
				default:
					// буфер заполнен — дропаем пакет для этого выхода
				}
			}
		}
		log.Printf("Publish finished: %s", srcConn.URL)

		// Сбрасываем статус всех выходов при завершении трансляции
		for _, url := range inputCfg.Outputs {
			sm.SetOutputActive(inputCfg.Name, url, false)
		}
	}
}

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	srt "github.com/datarhei/gosrt"
	"github.com/datarhei/joy4/av"
	"github.com/datarhei/joy4/format/flv"
	"github.com/datarhei/joy4/format/rtmp"
	"github.com/datarhei/joy4/format/ts"
	"github.com/pion/webrtc/v3"
)

type WHIPServer struct {
	port      int
	manager   *StreamManager
	server    *http.Server
	wg        sync.WaitGroup
	sessions  map[string]*WHIPSession
	sessionMu sync.RWMutex
}

type WHIPSession struct {
	inputName           string
	combinedSdpFilename string
	ffmpegCmd           *exec.Cmd
	outputMgr           *OutputManager
	stopCh              chan struct{}

	// Pipe для отправки RTP пакетов в ffmpeg
	audioWriter io.WriteCloser
	videoWriter io.WriteCloser

	// Codec data для RTMP/SRT хэдеров
	streams     []av.CodecData
	streamsOnce sync.Once
	streamsCh   chan struct{}
}

func NewWHIPServer(port int, manager *StreamManager) *WHIPServer {
	return &WHIPServer{
		port:     port,
		manager:  manager,
		sessions: make(map[string]*WHIPSession),
	}
}

func (w *WHIPServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/whip/", w.handleWHIP)

	w.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", w.port),
		Handler: mux,
	}

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		log.Printf("[WHIP] server started on :%d endpoint: /whip/{name}", w.port)
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[WHIP] server error: %v", err)
		}
	}()

	return nil
}

func (w *WHIPServer) Stop() error {
	w.sessionMu.Lock()
	for _, session := range w.sessions {
		w.stopSession(session)
	}
	w.sessions = make(map[string]*WHIPSession)
	w.sessionMu.Unlock()

	if w.server != nil {
		w.server.Close()
	}
	w.wg.Wait()
	return nil
}

func (w *WHIPServer) handleWHIP(wr http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(wr, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimPrefix(r.URL.Path, "/whip/")
	if name == "" {
		http.Error(wr, "Missing stream name in URL", http.StatusBadRequest)
		return
	}

	// Проверяем input в конфиге по url_path
	inputCfg := w.manager.GetInputByPath(r.URL.Path)
	if inputCfg == nil {
		http.Error(wr, "Input not found in config", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(wr, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	sdpOffer := string(body)
	if strings.HasPrefix(sdpOffer, "{\"type\":") {
		var offer webrtc.SessionDescription
		if err := json.Unmarshal(body, &offer); err == nil {
			sdpOffer = offer.SDP
		}
	}

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		http.Error(wr, "Failed to create PeerConnection", http.StatusInternalServerError)
		return
	}

	// Создаём сессию
	session := &WHIPSession{
		inputName: inputCfg.Name,
		outputMgr: NewOutputManager(),
		stopCh:    make(chan struct{}),
		streamsCh: make(chan struct{}),
	}

	// Запускаем ffmpeg pipeline
	go w.startFFmpegPipeline(session, inputCfg)

	// Ждём немного, чтобы ffmpeg запустился и pipe'ы создались
	time.Sleep(200 * time.Millisecond)

	// Обработчик изменения состояния соединения
	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Printf("[WHIP] Connection state changed to: %s", state)

		if state == webrtc.PeerConnectionStateClosed ||
			state == webrtc.PeerConnectionStateDisconnected ||
			state == webrtc.PeerConnectionStateFailed {
			log.Printf("[WHIP] Connection ended, stopping session for '%s'", inputCfg.Name)

			// Останавливаем сессию
			w.sessionMu.Lock()
			if session, exists := w.sessions[inputCfg.Name]; exists {
				w.stopSession(session)
				delete(w.sessions, inputCfg.Name)
			}
			w.sessionMu.Unlock()
		}
	})

	// OnTrack - записываем RTP пакеты в pipe
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("[WHIP] Received track: %s, kind: %s", track.ID(), track.Kind())

		if track.Kind() == webrtc.RTPCodecTypeVideo {
			// Обработка видео
			writer := session.videoWriter
			if writer == nil {
				log.Printf("[WHIP] Video writer not ready for track %s", track.ID())
				return
			}

			go w.handleVideoTrack(track, writer, session)

		} else if track.Kind() == webrtc.RTPCodecTypeAudio {
			// Обработка аудио
			writer := session.audioWriter
			if writer == nil {
				log.Printf("[WHIP] Audio writer not ready for track %s", track.ID())
				return
			}

			go w.handleAudioTrack(track, writer, session)
		}
	})

	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  sdpOffer,
	}
	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		http.Error(wr, "Failed to set remote description", http.StatusBadRequest)
		return
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		http.Error(wr, "Failed to create answer", http.StatusInternalServerError)
		return
	}
	if err := peerConnection.SetLocalDescription(answer); err != nil {
		http.Error(wr, "Failed to set local description", http.StatusInternalServerError)
		return
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	<-gatherComplete

	// Сохраняем сессию
	w.sessionMu.Lock()
	w.sessions[inputCfg.Name] = session
	w.sessionMu.Unlock()

	resp := peerConnection.LocalDescription().SDP
	wr.Header().Set("Location", r.URL.Path)
	wr.Header().Set("Content-Type", "application/sdp")
	wr.WriteHeader(http.StatusCreated)
	_, _ = wr.Write([]byte(resp))

	log.Printf("[WHIP] SDP answer sent for stream '%s'", inputCfg.Name)
}

func (w *WHIPServer) startFFmpegPipeline(session *WHIPSession, inputCfg *InputCfg) {
	rand.Seed(time.Now().UnixNano())
	videoUDPPort := rand.Intn(10000) + 40000 // 40000-49999
	audioUDPPort := rand.Intn(10000) + 50000 // 50000-59999

	session.combinedSdpFilename = fmt.Sprintf("whip_combined_%s.sdp", session.inputName)

	// Создаём единый SDP файл для видео и аудио, чтобы ffmpeg понимал, что они связаны
	combinedSdpContent := []byte(fmt.Sprintf(
		"v=0\r\n"+
			"o=- 0 0 IN IP4 127.0.0.1\r\n"+
			"s=WHIP\r\n"+
			"c=IN IP4 127.0.0.1\r\n"+
			"t=0 0\r\n"+
			"m=video %d RTP/AVP 96\r\n"+
			"a=rtpmap:96 H264/90000\r\n"+
			"m=audio %d RTP/AVP 111\r\n"+
			"a=rtpmap:111 opus/48000/2\r\n",
		videoUDPPort,
		audioUDPPort,
	))
	if err := ioutil.WriteFile(session.combinedSdpFilename, combinedSdpContent, 0644); err != nil {
		log.Printf("[WHIP] Failed to write Combined SDP file: %v", err)
		return
	}

	args := []string{
		"-y", // Перезаписывать существующие файлы
		"-thread_queue_size", "1024",
		"-rtbufsize", "2M",
		"-protocol_whitelist", "file,udp,rtp",
		"-i", session.combinedSdpFilename,
		"-map", "0:v",
		"-map", "0:a",
		// Копируем видео, конвертируем аудио
		"-c:v", "copy",
		"-c:a", "aac",
		"-bsf:v", "h264_mp4toannexb",
	}

	// Если есть хотя бы один выход, направляем ffmpeg в pipe
	hasOutputs := len(inputCfg.Outputs) > 0
	if hasOutputs {
		args = append(args, "-f", "flv", "pipe:1")
	}

	session.ffmpegCmd = exec.Command("./bin/ffmpeg.exe", args...)

	// Создаём UDP соединение для видео
	videoAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", videoUDPPort))
	if err != nil {
		log.Printf("[WHIP] Failed to resolve video UDP addr: %v", err)
		return
	}
	videoConn, err := net.DialUDP("udp", nil, videoAddr)
	if err != nil {
		log.Printf("[WHIP] Failed to create video UDP connection: %v", err)
		return
	}
	session.videoWriter = videoConn

	// Создаём UDP соединение для аудио
	audioAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", audioUDPPort))
	if err != nil {
		log.Printf("[WHIP] Failed to resolve audio UDP addr: %v", err)
		return
	}
	audioConn, err := net.DialUDP("udp", nil, audioAddr)
	if err != nil {
		log.Printf("[WHIP] Failed to create audio UDP connection: %v", err)
		return
	}
	session.audioWriter = audioConn

	// Получаем stdout для FLV только если есть выходы
	var stdout io.ReadCloser
	if hasOutputs {
		var err error
		stdout, err = session.ffmpegCmd.StdoutPipe()
		if err != nil {
			log.Printf("[WHIP] Failed to create stdout pipe: %v", err)
			return
		}
	}

	stderr, err := session.ffmpegCmd.StderrPipe()
	if err != nil {
		log.Printf("[WHIP] Failed to create stderr pipe: %v", err)
		return
	}

	// Логируем stderr
	go func() {
		// Просто читаем из stderr, чтобы не блокировать ffmpeg, но ничего не выводим
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
		}
	}()

	if err := session.ffmpegCmd.Start(); err != nil {
		log.Printf("[WHIP] Failed to start ffmpeg: %v", err)
		return
	}

	log.Printf("[WHIP] FFmpeg started for stream '%s' (video + audio)", session.inputName)

	// Инициализируем выходы
	w.manager.SetStatusActive(session.inputName, true)

	// Инициализируем ВСЕ выходы через наш менеджер
	for _, url := range inputCfg.Outputs {
		w.manager.RegisterOutput(session.inputName, url)
		session.outputMgr.AddOutput(url, 500, w.createOutputPusher(session, url))
	}

	// Горутина для динамического обновления выходов
	updateTicker := time.NewTicker(2 * time.Second)
	defer updateTicker.Stop()

	go func() {
		for {
			select {
			case <-session.stopCh:
				return
			case <-updateTicker.C:
				// Синхронизируем выходы с inputCfg.Outputs
				current := make(map[string]struct{})
				for _, url := range inputCfg.Outputs {
					current[url] = struct{}{}
					w.manager.RegisterOutput(session.inputName, url)
					session.outputMgr.AddOutput(url, 500, w.createOutputPusher(session, url))
				}
				for url := range session.outputMgr.AllOutputs() {
					if _, ok := current[url]; !ok {
						session.outputMgr.RemoveOutput(url)
					}
				}
			}
		}
	}()

	// Читаем FLV из stdout и отправляем в выходы только если они есть
	if hasOutputs && stdout != nil {
		w.processFLVStream(stdout, session, inputCfg)
	}

	// Обновляем статус в менеджере
	w.manager.SetStatusActive(session.inputName, false)

	// Удаляем временный SDP файл
	if session.combinedSdpFilename != "" {
		if err := os.Remove(session.combinedSdpFilename); err != nil {
			// Это не критичная ошибка, просто логируем, т.к. файл мог быть удален вручную или не создаться
			log.Printf("[WHIP] Failed to remove temporary SDP file %s: %v", session.combinedSdpFilename, err)
		}
	}

	log.Printf("[WHIP] Session stopped for '%s'", session.inputName)
}

func (w *WHIPServer) createOutputPusher(session *WHIPSession, url string) func(<-chan av.Packet, <-chan struct{}) {
	return func(ch <-chan av.Packet, stop <-chan struct{}) {
		var totalBytes int64
		inputName := session.inputName

		// Ждём, пока появятся stream-данные от ffmpeg, прежде чем подключаться
		select {
		case <-session.streamsCh:
			// Данные готовы
		case <-stop:
			return
		case <-session.stopCh:
			return
		}

		for {
			select {
			case <-stop:
				log.Printf("[WHIP] Output stopped: %s", url)
				return
			case <-session.stopCh:
				log.Printf("[WHIP] Session stopped, stopping output: %s", url)
				return
			default:
			}

			log.Printf("[WHIP] Trying to connect to %s", url)

			if strings.HasPrefix(url, "file://") {
				filename := strings.TrimPrefix(url, "file://")
				file, err := os.Create(filename)
				if err != nil {
					log.Printf("[WHIP] Failed to create file %s: %v", filename, err)
					time.Sleep(5 * time.Second)
					continue // Повторная попытка
				}

				w.manager.SetOutputActive(inputName, url, true)
				log.Printf("[WHIP] Writing FLV to file: %s", filename)

				muxer := flv.NewMuxer(file)
				err = muxer.WriteHeader(session.streams)
				if err != nil {
					log.Printf("[WHIP] Failed to write FLV header to file %s: %v", filename, err)
					file.Close()
					w.manager.SetOutputActive(inputName, url, false)
					time.Sleep(5 * time.Second)
					continue
				}

				for {
					select {
					case <-stop:
						if err := muxer.WriteTrailer(); err != nil {
							log.Printf("[WHIP] Failed to write FLV trailer to file %s: %v", filename, err)
						}
						file.Close()
						w.manager.SetOutputActive(inputName, url, false)
						log.Printf("[WHIP] File output stopped: %s", filename)
						return
					case <-session.stopCh:
						if err := muxer.WriteTrailer(); err != nil {
							log.Printf("[WHIP] Failed to write FLV trailer to file %s: %v", filename, err)
						}
						file.Close()
						w.manager.SetOutputActive(inputName, url, false)
						return
					case pkt, ok := <-ch:
						if !ok {
							if err := muxer.WriteTrailer(); err != nil {
								log.Printf("[WHIP] Failed to write FLV trailer to file %s: %v", filename, err)
							}
							file.Close()
							w.manager.SetOutputActive(inputName, url, false)
							return
						}
						err = muxer.WritePacket(pkt)
						if err != nil {
							log.Printf("[WHIP] Write error to file %s: %v", filename, err)
							if err := muxer.WriteTrailer(); err != nil {
								log.Printf("[WHIP] Failed to write FLV trailer to file %s: %v", filename, err)
							}
							file.Close()
							w.manager.SetOutputActive(inputName, url, false)
							time.Sleep(5 * time.Second)
							break // Выход из внутреннего цикла для переподключения
						}
						totalBytes += int64(len(pkt.Data))
						w.manager.UpdateOutputBitrate(inputName, url, totalBytes)
					}
				}
			} else if strings.HasPrefix(url, "rtmp://") {
				dstConn, err := rtmp.Dial(url, rtmp.DialOptions{})
				if err != nil {
					log.Printf("[WHIP] Failed to connect to %s: %v", url, err)
					w.manager.mu.RLock()
					reconnectInterval := w.manager.config.ReconnectInterval
					w.manager.mu.RUnlock()
					time.Sleep(time.Duration(reconnectInterval) * time.Second)
					continue
				}
				w.manager.SetOutputActive(inputName, url, true)

				err = dstConn.WriteHeader(session.streams)
				if err != nil {
					log.Printf("[WHIP] Failed to write header to %s: %v", url, err)
					dstConn.Close()
					w.manager.SetOutputActive(inputName, url, false)
					w.manager.mu.RLock()
					reconnectInterval := w.manager.config.ReconnectInterval
					w.manager.mu.RUnlock()
					time.Sleep(time.Duration(reconnectInterval) * time.Second)
					continue
				}
				for {
					select {
					case <-stop:
						dstConn.Close()
						w.manager.SetOutputActive(inputName, url, false)
						return
					case <-session.stopCh:
						dstConn.Close()
						w.manager.SetOutputActive(inputName, url, false)
						return
					case pkt, ok := <-ch:
						if !ok {
							dstConn.Close()
							w.manager.SetOutputActive(inputName, url, false)
							return
						}
						err = dstConn.WritePacket(pkt)
						if err != nil {
							log.Printf("[WHIP] Write error to %s: %v", url, err)
							dstConn.Close()
							w.manager.SetOutputActive(inputName, url, false)
							w.manager.mu.RLock()
							reconnectInterval := w.manager.config.ReconnectInterval
							w.manager.mu.RUnlock()
							time.Sleep(time.Duration(reconnectInterval) * time.Second)
							break
						}
						totalBytes += int64(len(pkt.Data))
						w.manager.UpdateOutputBitrate(inputName, url, totalBytes)
					}
				}
			} else if strings.HasPrefix(url, "srt://") {
				srtAddr := strings.TrimPrefix(url, "srt://")
				if idx := strings.Index(srtAddr, "?"); idx != -1 {
					srtAddr = srtAddr[:idx]
				}
				cfgSRT := srt.DefaultConfig()

				w.manager.mu.RLock()
				srtSettings := w.manager.config.SRTSettings
				reconnectInterval := w.manager.config.ReconnectInterval
				w.manager.mu.RUnlock()

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
				log.Printf("[WHIP] SRT connecting to %s with latency=%v, streamid=%s, timeout=%v", srtAddr, cfgSRT.Latency, cfgSRT.StreamId, cfgSRT.ConnectionTimeout)
				conn, err := srt.Dial("srt", srtAddr, cfgSRT)
				if err != nil {
					log.Printf("[WHIP] Failed to connect to SRT %s: %v", url, err)
					time.Sleep(time.Duration(reconnectInterval) * time.Second)
					continue
				}
				w.manager.SetOutputActive(inputName, url, true)
				var tsBuf bytes.Buffer
				muxer := ts.NewMuxer(&tsBuf)

				err = muxer.WriteHeader(session.streams)
				if err != nil {
					log.Printf("[WHIP] TS WriteHeader error for %s: %v", url, err)
					conn.Close()
					w.manager.SetOutputActive(inputName, url, false)
					time.Sleep(time.Duration(reconnectInterval) * time.Second)
					continue
				}
				for {
					select {
					case <-stop:
						tsBuf.Reset()
						if err := muxer.WriteTrailer(); err != nil {
							log.Printf("[WHIP] TS WriteTrailer error for %s: %v", url, err)
						}
						if tsBuf.Len() > 0 {
							_, _ = conn.Write(tsBuf.Bytes())
						}
						conn.Close()
						w.manager.SetOutputActive(inputName, url, false)
						return
					case <-session.stopCh:
						tsBuf.Reset()
						if err := muxer.WriteTrailer(); err != nil {
							log.Printf("[WHIP] TS WriteTrailer error for %s: %v", url, err)
						}
						if tsBuf.Len() > 0 {
							_, _ = conn.Write(tsBuf.Bytes())
						}
						conn.Close()
						w.manager.SetOutputActive(inputName, url, false)
						return
					case pkt, ok := <-ch:
						if !ok {
							tsBuf.Reset()
							if err := muxer.WriteTrailer(); err != nil {
								log.Printf("[WHIP] TS WriteTrailer error for %s: %v", url, err)
							}
							if tsBuf.Len() > 0 {
								_, _ = conn.Write(tsBuf.Bytes())
							}
							conn.Close()
							w.manager.SetOutputActive(inputName, url, false)
							return
						}
						// Сохраняем текущую позицию в буфере перед записью
						bufferPosBefore := tsBuf.Len()
						err = muxer.WritePacket(pkt)
						if err != nil {
							log.Printf("[WHIP] TS WritePacket error for %s: %v", url, err)
							conn.Close()
							w.manager.SetOutputActive(inputName, url, false)
							time.Sleep(time.Duration(reconnectInterval) * time.Second)
							break
						}
						// Отправляем только новые данные, добавленные этим пакетом
						if tsBuf.Len() > bufferPosBefore {
							newData := tsBuf.Bytes()[bufferPosBefore:]
							_, err = conn.Write(newData)
							if err != nil {
								log.Printf("[WHIP] SRT Write error for %s: %v", url, err)
								conn.Close()
								w.manager.SetOutputActive(inputName, url, false)
								time.Sleep(time.Duration(reconnectInterval) * time.Second)
								break
							}
							totalBytes += int64(len(newData))
							w.manager.UpdateOutputBitrate(inputName, url, totalBytes)

							// Сбрасываем буфер если он стал слишком большим (более 1MB)
							if tsBuf.Len() > 1024*1024 {
								tsBuf.Reset()
							}
						}
					}
				}
			}
		}
	}
}

func (w *WHIPServer) processFLVStream(reader io.Reader, session *WHIPSession, inputCfg *InputCfg) {
	defer w.manager.SetStatusActive(session.inputName, false)

	log.Printf("[WHIP] Starting FLV processing for '%s'", session.inputName)

	flvReader := flv.NewDemuxer(reader)
	streams, err := flvReader.Streams()
	if err != nil {
		log.Printf("[WHIP] Failed to read FLV streams: %v", err)
		return
	}

	// Сохраняем streams и оповещаем все выходы, что хэдеры готовы
	session.streamsOnce.Do(func() {
		session.streams = streams
		close(session.streamsCh)
		log.Printf("[WHIP] Codec data ready for stream '%s'", session.inputName)
	})

	log.Printf("[WHIP] FLV streams detected: %d", len(streams))
	for i, stream := range streams {
		log.Printf("[WHIP] Stream %d: %T", i, stream)
	}

	packetCount := 0

	for {
		select {
		case <-session.stopCh:
			return
		default:
		}

		pkt, err := flvReader.ReadPacket()
		if err != nil {
			if err == io.EOF {
				log.Printf("[WHIP] FLV stream ended for '%s'", session.inputName)
			} else {
				log.Printf("[WHIP] FLV read error: %v", err)
			}
			return
		}

		packetCount++
		if packetCount%100 == 0 {
			// log.Printf("[WHIP] FLV packets processed: %d (stream idx: %d)", packetCount, pkt.Idx)
		}

		// Отправляем пакет во все выходы
		for _, w := range session.outputMgr.AllOutputs() {
			select {
			case w.ch <- pkt:
			default:
				// буфер заполнен — дропаем пакет для этого выхода
			}
		}
	}
}

func (w *WHIPServer) stopSession(session *WHIPSession) {
	log.Printf("[WHIP] Stopping session for '%s'", session.inputName)

	// Сначала закрываем канал остановки
	select {
	case <-session.stopCh:
		// Уже закрыт
	default:
		close(session.stopCh)
	}

	// Корректно закрываем writers чтобы FFmpeg получил EOF
	if session.videoWriter != nil {
		session.videoWriter.Close()
		log.Printf("[WHIP] Video writer closed for '%s'", session.inputName)
	}
	if session.audioWriter != nil {
		session.audioWriter.Close()
		log.Printf("[WHIP] Audio writer closed for '%s'", session.inputName)
	}

	// Ждём немного чтобы FFmpeg корректно завершился
	if session.ffmpegCmd != nil && session.ffmpegCmd.Process != nil {
		// Даём FFmpeg 2 секунды на корректное завершение
		done := make(chan error, 1)
		go func() {
			done <- session.ffmpegCmd.Wait()
		}()

		select {
		case <-time.After(2 * time.Second):
			log.Printf("[WHIP] FFmpeg timeout, killing process for '%s'", session.inputName)
			session.ffmpegCmd.Process.Kill()
		case err := <-done:
			if err != nil {
				log.Printf("[WHIP] FFmpeg finished with error for '%s': %v", session.inputName, err)
			} else {
				log.Printf("[WHIP] FFmpeg finished normally for '%s'", session.inputName)
			}
		}
	}

	// Останавливаем все выходы
	for url := range session.outputMgr.AllOutputs() {
		session.outputMgr.RemoveOutput(url)
	}

	// Обновляем статус в менеджере
	w.manager.SetStatusActive(session.inputName, false)

	// Удаляем временный SDP файл
	if session.combinedSdpFilename != "" {
		if err := os.Remove(session.combinedSdpFilename); err != nil {
			// Это не критичная ошибка, просто логируем, т.к. файл мог быть удален вручную или не создаться
			log.Printf("[WHIP] Failed to remove temporary SDP file %s: %v", session.combinedSdpFilename, err)
		}
	}

	log.Printf("[WHIP] Session stopped for '%s'", session.inputName)
}

func (w *WHIPServer) handleVideoTrack(track *webrtc.TrackRemote, writer io.WriteCloser, session *WHIPSession) {
	packetCount := 0
	log.Printf("[WHIP] Started handling video track: %s", track.ID())

	for {
		select {
		case <-session.stopCh:
			return
		default:
		}

		packet, _, err := track.ReadRTP()
		if err != nil {
			log.Printf("[WHIP] Video track %s ended: %v", track.ID(), err)
			return
		}

		// Отправляем полный RTP пакет в ffmpeg по UDP, как ожидает SDP
		payload, err := packet.Marshal()
		if err != nil {
			log.Printf("[WHIP] Failed to marshal video RTP packet: %v", err)
			continue
		}

		if len(payload) > 0 {
			_, err := writer.Write(payload)
			if err != nil {
				// Игнорируем ошибки, т.к. ffmpeg может быть еще не готов
			}

			packetCount++
			if packetCount%100 == 0 {
				// log.Printf("[WHIP] Video packets written: %d (last payload size: %d)", packetCount, len(payload))
			}
		}
	}
}

func (w *WHIPServer) handleAudioTrack(track *webrtc.TrackRemote, writer io.WriteCloser, session *WHIPSession) {
	packetCount := 0
	log.Printf("[WHIP] Started handling audio track: %s", track.ID())

	for {
		select {
		case <-session.stopCh:
			return
		default:
		}

		packet, _, err := track.ReadRTP()
		if err != nil {
			log.Printf("[WHIP] Audio track %s ended: %v", track.ID(), err)
			return
		}

		// Отправляем полный RTP пакет в ffmpeg по UDP, как ожидает SDP
		payload, err := packet.Marshal()
		if err != nil {
			log.Printf("[WHIP] Failed to marshal audio RTP packet: %v", err)
			continue
		}

		if len(payload) > 0 {
			_, err := writer.Write(payload)
			if err != nil {
				// Игнорируем ошибки, т.к. ffmpeg может быть еще не готов
			}

			packetCount++
			if packetCount%200 == 0 {
				// log.Printf("[WHIP] Audio packets written: %d (last payload size: %d)", packetCount, len(payload))
			}
		}
	}
}

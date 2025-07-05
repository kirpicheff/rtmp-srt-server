package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/asticode/go-astits"
	srt "github.com/datarhei/gosrt"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/codec/aacparser"
	"github.com/nareix/joy4/codec/h264parser"
	"github.com/nareix/joy4/format/rtmp"
)

type SRTServer struct {
	port        int
	config      *Config
	manager     *StreamManager
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	server      *srt.Server
	connections map[string]srt.Conn
	mu          sync.RWMutex
}

func NewSRTServer(port int, config *Config, manager *StreamManager) *SRTServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &SRTServer{
		port:        port,
		config:      config,
		manager:     manager,
		ctx:         ctx,
		cancel:      cancel,
		connections: make(map[string]srt.Conn),
	}
}

func (s *SRTServer) Start() error {
	// Создаем SRT сервер
	serverConfig := srt.DefaultConfig()
	if s.config.SRTSettings.Latency > 0 {
		serverConfig.Latency = time.Duration(s.config.SRTSettings.Latency) * time.Millisecond
	}
	if s.config.SRTSettings.Passphrase != "" {
		serverConfig.Passphrase = s.config.SRTSettings.Passphrase
	}

	s.server = &srt.Server{
		Addr:   fmt.Sprintf(":%d", s.port),
		Config: &serverConfig,
		HandleConnect: func(req srt.ConnRequest) srt.ConnType {
			streamID := req.StreamId()
			log.Printf("[SRT] Incoming connection with streamID: %s", streamID)
			return srt.PUBLISH
		},
		HandlePublish: func(conn srt.Conn) {
			s.wg.Add(1)
			go s.handleConnection(conn)
		},
	}

	// Запускаем сервер
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			if s.ctx.Err() != nil {
				return
			}
			log.Printf("[SRT] Server error: %v", err)
		}
	}()

	return nil
}

func (s *SRTServer) Stop() error {
	s.cancel()
	if s.server != nil {
		s.server.Shutdown()
	}
	s.wg.Wait()
	log.Printf("[SRT] SRT server stopped")
	return nil
}

func (s *SRTServer) handleConnection(conn srt.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	streamID := conn.StreamId()
	inputName := "obs"
	if streamID != "" {
		inputName = streamID
	}
	log.Printf("[SRT] New SRT connection from %s with streamID: %s", conn.RemoteAddr(), streamID)

	inputCfg := s.manager.GetInputByName(inputName)
	if inputCfg == nil {
		log.Printf("[SRT] No input config found for %s, creating default", inputName)
		inputCfg = &InputCfg{
			Name:    inputName,
			URLPath: "/live/stream",
			Outputs: []string{},
		}
		s.manager.AddInput(*inputCfg)
	}

	s.manager.SetStatusActive(inputName, true)
	defer s.manager.SetStatusActive(inputName, false)

	// Находим SRT и RTMP выходы
	var srtOutputs, rtmpOutputs, fileOutputs []string
	for _, output := range inputCfg.Outputs {
		if strings.HasPrefix(output, "srt://") {
			srtOutputs = append(srtOutputs, output)
		} else if strings.HasPrefix(output, "rtmp://") {
			rtmpOutputs = append(rtmpOutputs, output)
		} else if strings.HasPrefix(output, "file://") {
			fileOutputs = append(fileOutputs, output)
		}
	}

	if len(srtOutputs) == 0 && len(rtmpOutputs) == 0 && len(fileOutputs) == 0 {
		log.Printf("[SRT] No outputs configured for %s", inputName)
		return
	}

	if len(srtOutputs) > 0 {
		log.Printf("[SRT] Found %d SRT outputs for %s", len(srtOutputs), inputName)
		for _, outputURL := range srtOutputs {
			s.manager.RegisterOutput(inputName, outputURL)
		}
	}
	if len(rtmpOutputs) > 0 {
		log.Printf("[SRT] Found %d RTMP outputs for %s", len(rtmpOutputs), inputName)
		for _, outputURL := range rtmpOutputs {
			s.manager.RegisterOutput(inputName, outputURL)
		}
	}
	if len(fileOutputs) > 0 {
		log.Printf("[SRT] Found %d file outputs for %s", len(fileOutputs), inputName)
		for _, outputURL := range fileOutputs {
			s.manager.RegisterOutput(inputName, outputURL)
		}
	}

	// Создаем каналы для каждого выхода
	outputChannels := make(map[string]chan []byte)
	stopChannels := make(map[string]chan struct{})

	// Функция для создания выхода
	createOutput := func(outputURL string) {
		if _, exists := outputChannels[outputURL]; exists {
			return // Выход уже создан
		}

		ch := make(chan []byte, 1000)
		if strings.HasPrefix(outputURL, "rtmp://") {
			ch = make(chan []byte, 5000)
		}
		stop := make(chan struct{})
		outputChannels[outputURL] = ch
		stopChannels[outputURL] = stop

		s.wg.Add(1)
		if strings.HasPrefix(outputURL, "srt://") {
			go s.handleSRTOutput(inputName, outputURL, ch, stop)
		} else if strings.HasPrefix(outputURL, "rtmp://") {
			go s.handleRTMPOutput(inputName, outputURL, ch, stop)
		} else if strings.HasPrefix(outputURL, "file://") {
			go s.handleFileOutput(inputName, outputURL, ch, stop)
		}
	}

	// Создаем начальные выходы
	for _, outputURL := range srtOutputs {
		createOutput(outputURL)
	}
	for _, outputURL := range rtmpOutputs {
		createOutput(outputURL)
	}
	for _, outputURL := range fileOutputs {
		createOutput(outputURL)
	}

	// Горутина для динамического обновления выходов
	updateTicker := time.NewTicker(2 * time.Second)
	defer updateTicker.Stop()

	stopUpdateChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-stopUpdateChan:
				return
			case <-updateTicker.C:
				// Получаем актуальный список выходов
				currentOutputs := make(map[string]struct{})
				for _, url := range inputCfg.Outputs {
					if strings.HasPrefix(url, "srt://") || strings.HasPrefix(url, "rtmp://") || strings.HasPrefix(url, "file://") {
						currentOutputs[url] = struct{}{}
						s.manager.RegisterOutput(inputName, url)
						createOutput(url)
					}
				}

				// Удаляем неактуальные выходы
				for url, stopCh := range stopChannels {
					if _, exists := currentOutputs[url]; !exists {
						close(stopCh)
						delete(outputChannels, url)
						delete(stopChannels, url)
					}
				}
			}
		}
	}()

	// Читаем данные из входящего SRT соединения и отправляем во все выходы
	buffer := make([]byte, 1316)
	var totalBytes int64
	var packetCount int64
	startTime := time.Now()
	lastLogTime := time.Now()

	for {
		select {
		case <-s.ctx.Done():
			log.Printf("[SRT] Context cancelled, stopping connection")
			close(stopUpdateChan)
			return
		default:
		}

		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Printf("[SRT] Connection closed by client")
			} else {
				log.Printf("[SRT] Read error: %v", err)
			}
			break
		}

		if n == 0 {
			continue
		}

		data := make([]byte, n)
		copy(data, buffer[:n])

		totalBytes += int64(n)
		packetCount++

		for outputURL, ch := range outputChannels {
			select {
			case ch <- data:
				// ok
			default:
				if strings.HasPrefix(outputURL, "rtmp://") {
					log.Printf("[SRT] RTMP Output buffer full for %s, dropping packet", outputURL)
				} else {
					log.Printf("[SRT] Output buffer full for %s, dropping packet", outputURL)
				}
			}
		}

		if time.Since(lastLogTime) >= 5*time.Second {
			elapsed := time.Since(startTime).Seconds()
			if elapsed > 0 {
				// log.Printf("[SRT] Input stats - Packets: %d, Bytes: %d, Bitrate: %.1f kbps, AvgPacketSize: %.1f bytes",
				// 	packetCount, totalBytes, bitrate, avgPacketSize)
			}
			lastLogTime = time.Now()
		}
	}

	close(stopUpdateChan)
	for _, stop := range stopChannels {
		close(stop)
	}

	log.Printf("[SRT] Connection closed: %s", conn.RemoteAddr())
}

func (s *SRTServer) handleSRTOutput(inputName, outputURL string, dataCh <-chan []byte, stopCh <-chan struct{}) {
	defer s.wg.Done()

	srtAddr := strings.TrimPrefix(outputURL, "srt://")
	if idx := strings.Index(srtAddr, "?"); idx != -1 {
		srtAddr = srtAddr[:idx]
	}

	log.Printf("[SRT] Starting SRT output to %s", outputURL)

	var totalBytes int64
	var packetCount int64

	for {
		select {
		case <-stopCh:
			log.Printf("[SRT] SRT output stopped: %s", outputURL)
			s.manager.SetOutputActive(inputName, outputURL, false)
			return
		default:
		}

		// Подключаемся к SRT выходу
		config := srt.DefaultConfig()
		if s.config.SRTSettings.Latency > 0 {
			config.Latency = time.Duration(s.config.SRTSettings.Latency) * time.Millisecond
		}
		if s.config.SRTSettings.Passphrase != "" {
			config.Passphrase = s.config.SRTSettings.Passphrase
		}
		if s.config.SRTSettings.StreamID != "" {
			config.StreamId = s.config.SRTSettings.StreamID
		}

		conn, err := srt.Dial("srt", srtAddr, config)
		if err != nil {
			log.Printf("[SRT] Failed to connect to SRT output %s: %v", outputURL, err)
			s.manager.SetOutputActive(inputName, outputURL, false)
			s.manager.IncrementOutputError(inputName, outputURL)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("[SRT] Connected to SRT output: %s", outputURL)
		s.manager.SetOutputActive(inputName, outputURL, true)

		// Отправляем данные
		for {
			select {
			case <-stopCh:
				conn.Close()
				s.manager.SetOutputActive(inputName, outputURL, false)
				return
			case data, ok := <-dataCh:
				if !ok {
					conn.Close()
					s.manager.SetOutputActive(inputName, outputURL, false)
					return
				}

				// Отправляем данные с таймаутом
				conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
				_, err := conn.Write(data)
				if err != nil {
					log.Printf("[SRT] Write error to %s: %v", outputURL, err)
					s.manager.IncrementOutputError(inputName, outputURL)
					conn.Close()
					s.manager.SetOutputActive(inputName, outputURL, false)
					break
				}

				totalBytes += int64(len(data))
				packetCount++

				// Обновляем битрейт в StreamManager
				s.manager.UpdateOutputBitrate(inputName, outputURL, totalBytes)
			}
		}
	}
}

// RTMP output: repack TS to FLV and send to RTMP server
func (s *SRTServer) handleRTMPOutput(inputName, outputURL string, dataCh <-chan []byte, stopCh <-chan struct{}) {
	defer s.wg.Done()

	reconnectInterval := 5 // seconds
	log.Printf("[SRT] Starting RTMP output to %s", outputURL)

	for {
		select {
		case <-stopCh:
			log.Printf("[SRT] RTMP output stopped: %s", outputURL)
			s.manager.SetOutputActive(inputName, outputURL, false)
			return
		default:
		}

		dstConn, err := rtmp.Dial(outputURL)
		if err != nil {
			log.Printf("[SRT] RTMP Dial error for %s: %v", outputURL, err)
			s.manager.SetOutputActive(inputName, outputURL, false)
			s.manager.IncrementOutputError(inputName, outputURL)
			time.Sleep(time.Duration(reconnectInterval) * time.Second)
			continue
		}
		s.manager.SetOutputActive(inputName, outputURL, true)
		log.Printf("[SRT] Connected to RTMP output: %s", outputURL)

		pipeReader, pipeWriter := io.Pipe()

		writerDone := make(chan struct{})
		go func() {
			defer close(writerDone)
			defer pipeWriter.Close()
			for {
				select {
				case <-stopCh:
					return
				case data, ok := <-dataCh:
					if !ok {
						return
					}
					_, err := pipeWriter.Write(data)
					if err != nil {
						// Pipe closed, probably by the reader.
						return
					}
				}
			}
		}()

		err = s.processRTMPStream(pipeReader, dstConn, inputName, outputURL)
		if err != nil {
			log.Printf("[SRT] RTMP stream processing error for %s: %v. Reconnecting...", outputURL, err)
		}

		dstConn.Close()
		pipeReader.Close()
		<-writerDone // Wait for writer to finish before looping to reconnect

		s.manager.SetOutputActive(inputName, outputURL, false)
		s.manager.IncrementOutputError(inputName, outputURL)

		select {
		case <-stopCh:
			log.Printf("[SRT] RTMP output stopped during reconnect wait: %s", outputURL)
			return
		case <-time.After(time.Duration(reconnectInterval) * time.Second):
		}
	}
}

func (s *SRTServer) processRTMPStream(reader io.Reader, dstConn *rtmp.Conn, inputName, outputURL string) error {
	demuxer := astits.NewDemuxer(context.Background(), reader)
	var videoPID, audioPID uint16
	var videoCodecData av.VideoCodecData
	var audioCodecData av.AudioCodecData
	var streams []av.CodecData
	var flvHeaderWritten bool

	var baseTime time.Duration
	var baseTimeSet bool

	var lastVideoTime time.Duration
	var lastAudioTime time.Duration

	var totalBytes int64

	for {
		data, err := demuxer.NextData()
		if err != nil {
			if err == io.EOF || err == io.ErrClosedPipe {
				return fmt.Errorf("demuxer stopped: %w", err)
			}
			return fmt.Errorf("demuxer error: %w", err)
		}

		if data == nil {
			continue
		}

		// First, we need to detect PIDs and create codec data before processing any packets.
		// This part runs until the RTMP header is written.
		if !flvHeaderWritten {
			// 1. Detect PIDs from PMT
			if data.PMT != nil {
				for _, es := range data.PMT.ElementaryStreams {
					if es.StreamType == astits.StreamTypeH264Video {
						videoPID = es.ElementaryPID
					} else if es.StreamType == astits.StreamTypeAACAudio {
						audioPID = es.ElementaryPID
					}
				}
			}

			if data.PES != nil {
				// 2. Fallback PID detection
				if videoPID == 0 || audioPID == 0 {
					pes := data.PES.Data
					if videoPID == 0 && len(pes) > 4 && pes[0] == 0x00 && pes[1] == 0x00 && pes[2] == 0x00 && pes[3] == 0x01 {
						nalType := pes[4] & 0x1F
						if nalType >= 1 && nalType <= 12 { // Valid H264 NAL types
							videoPID = data.PID
						}
					} else if audioPID == 0 && len(pes) > 2 && (pes[0] == 0xFF) && ((pes[1] & 0xF0) == 0xF0) {
						audioPID = data.PID
					}
				}

				// 3. Create CodecData
				if data.PID == videoPID && videoCodecData == nil {
					nalus, _ := h264parser.SplitNALUs(data.PES.Data)
					var sps, pps []byte
					for _, nalu := range nalus {
						if len(nalu) == 0 {
							continue
						}
						naluType := nalu[0] & 0x1f
						switch naluType {
						case 7: // SPS
							sps = nalu
						case 8: // PPS
							pps = nalu
						}
					}
					if sps != nil && pps != nil {
						vcd, err := h264parser.NewCodecDataFromSPSAndPPS(sps, pps)
						if err == nil {
							videoCodecData = vcd
						}
					}
				} else if data.PID == audioPID && audioCodecData == nil {
					if cfg, _, _, _, err := aacparser.ParseADTSHeader(data.PES.Data); err == nil {
						acd, err := aacparser.NewCodecDataFromMPEG4AudioConfig(cfg)
						if err == nil {
							audioCodecData = acd
						}
					}
				}
			}

			// 4. Write Header when both video and audio codecs are ready
			if videoCodecData != nil && audioCodecData != nil {
				streams = []av.CodecData{videoCodecData, audioCodecData}
				if err := dstConn.WriteHeader(streams); err != nil {
					return fmt.Errorf("failed to write RTMP header: %w", err)
				}
				flvHeaderWritten = true
			}
		}

		// If the header hasn't been written, or we don't have a PES packet, skip to the next data.
		if !flvHeaderWritten || data.PES == nil {
			continue
		}

		// Now, process the packets.
		// We must wait for the first keyframe to set the base timestamp.
		// All packets before the first keyframe will be dropped.

		var isKeyFrame bool
		if data.PID == videoPID {
			nalus, _ := h264parser.SplitNALUs(data.PES.Data)
			for _, nalu := range nalus {
				if len(nalu) > 0 && (nalu[0]&0x1F) == 5 {
					isKeyFrame = true
					break
				}
			}
		}

		if !baseTimeSet {
			if isKeyFrame {
				if data.PES.Header.OptionalHeader != nil && data.PES.Header.OptionalHeader.PTS != nil {
					pts := data.PES.Header.OptionalHeader.PTS.Base
					baseTime = time.Duration(pts/90) * time.Millisecond
					baseTimeSet = true
				} else {
					continue // Keyframe without a timestamp, drop it.
				}
			} else {
				continue // Not a keyframe, and we are not started yet. Drop it.
			}
		}

		// At this point, baseTime is set, and we can process all packets.
		var pkt av.Packet
		var streamIndex int8

		if data.PID == videoPID {
			streamIndex = 0
		} else if data.PID == audioPID {
			streamIndex = 1
		} else {
			continue // Not a stream we are interested in
		}

		if data.PES.Header.OptionalHeader == nil || data.PES.Header.OptionalHeader.PTS == nil {
			continue // No timestamp
		}
		pts := data.PES.Header.OptionalHeader.PTS.Base
		tsMilli := time.Duration(pts/90) * time.Millisecond

		// Normalize timestamp against baseTime
		pkt.Time = tsMilli - baseTime
		if pkt.Time < 0 {
			// This can happen if audio packets arrive with an earlier timestamp than the first keyframe.
			// Set their timestamp to 0 to ensure they are sent immediately with the first video frame.
			pkt.Time = 0
		}

		pkt.Idx = streamIndex
		pkt.Data = data.PES.Data

		if pkt.Idx == 0 { // Video
			if pkt.Time < lastVideoTime {
				pkt.Time = lastVideoTime
			}
			lastVideoTime = pkt.Time
			pkt.IsKeyFrame = isKeyFrame
		} else { // Audio
			if pkt.Time < lastAudioTime {
				pkt.Time = lastAudioTime
			}
			lastAudioTime = pkt.Time
			// Strip ADTS header from audio packet
			if len(pkt.Data) > 7 {
				pkt.Data = pkt.Data[7:]
			}
		}

		if err := dstConn.WritePacket(pkt); err != nil {
			return fmt.Errorf("failed to write RTMP packet: %w", err)
		}

		totalBytes += int64(len(pkt.Data))
		s.manager.UpdateOutputBitrate(inputName, outputURL, totalBytes)
	}
}

// Функция для обработки записи в файл
func (s *SRTServer) handleFileOutput(inputName, outputURL string, dataCh <-chan []byte, stopCh <-chan struct{}) {
	defer s.wg.Done()

	filePath := strings.TrimPrefix(outputURL, "file://")
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("[SRT] Failed to create file %s: %v", filePath, err)
		s.manager.SetOutputActive(inputName, outputURL, false)
		return
	}
	defer file.Close()

	log.Printf("[SRT] Writing to file: %s", filePath)
	s.manager.SetOutputActive(inputName, outputURL, true)

	var totalBytes int64

	for {
		select {
		case <-stopCh:
			log.Printf("[SRT] File output stopped: %s", filePath)
			s.manager.SetOutputActive(inputName, outputURL, false)
			return
		case data, ok := <-dataCh:
			if !ok {
				s.manager.SetOutputActive(inputName, outputURL, false)
				return
			}
			_, err := file.Write(data)
			if err != nil {
				log.Printf("[SRT] Write error to file %s: %v", filePath, err)
				s.manager.SetOutputActive(inputName, outputURL, false)
				return
			}
			totalBytes += int64(len(data))
			s.manager.UpdateOutputBitrate(inputName, outputURL, totalBytes)
		}
	}
}

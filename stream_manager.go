package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Stream представляет поток данных
type Stream struct {
	Name    string
	URLPath string
	mu      sync.Mutex
	data    []byte
}

// WriteData добавляет данные в поток
func (s *Stream) WriteData(data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, data...)
}

// GetData возвращает данные потока
func (s *Stream) GetData() []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data
}

type StreamStatus struct {
	Name        string          `json:"name"`
	URLPath     string          `json:"url_path"`
	Active      bool            `json:"active"`
	Connections int             `json:"connections"`
	ErrorCount  int             `json:"error_count"`
	Outputs     []*OutputStatus `json:"outputs,omitempty"`
}

type OutputStatus struct {
	URL         string  `json:"url"`
	Active      bool    `json:"active"`
	BitrateKbps float64 `json:"bitrate_kbps"`
	ErrorCount  int     `json:"error_count"`
	Uptime      string  `json:"uptime"`

	// Внутренние поля для подсчёта битрейта
	prevBytes int64
	prevTime  time.Time
	startTime time.Time // время старта активности
}

type StreamManager struct {
	mu      sync.RWMutex
	inputs  map[string]*InputCfg
	status  map[string]*StreamStatus
	outputs map[string]map[string]*OutputStatus // inputName -> url -> OutputStatus
	config  *Config                             // ссылка на глобальную конфигурацию
}

func validateRTMPURL(rawURL string) error {
	if !strings.HasPrefix(rawURL, "rtmp://") {
		return errors.New("URL must start with rtmp://")
	}
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return errors.New("URL must be a valid RTMP address")
	}
	return nil
}

// Конструктор
func NewStreamManager(cfgs []InputCfg, cfg *Config) *StreamManager {
	sm := &StreamManager{
		inputs:  make(map[string]*InputCfg),
		status:  make(map[string]*StreamStatus),
		outputs: make(map[string]map[string]*OutputStatus),
		config:  cfg,
	}
	for _, c := range cfgs {
		c := c
		sm.inputs[c.Name] = &c
		sm.status[c.Name] = &StreamStatus{
			Name:    c.Name,
			URLPath: c.URLPath,
			Active:  false,
		}
	}
	return sm
}

func (sm *StreamManager) AddInput(input InputCfg) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.inputs[input.Name] = &input
	sm.status[input.Name] = &StreamStatus{Name: input.Name, URLPath: input.URLPath}
	log.Printf("[API] Added input %s", input.Name)
}

func (sm *StreamManager) RemoveInput(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.inputs, name)
	delete(sm.status, name)
	delete(sm.outputs, name)
	log.Printf("[API] Removed input %s", name)
}

func (sm *StreamManager) ListInputs() []InputCfg {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	var inputs []InputCfg
	for _, v := range sm.inputs {
		inputs = append(inputs, *v)
	}
	return inputs
}

func (sm *StreamManager) GetInputByPath(path string) *InputCfg {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	for _, input := range sm.inputs {
		if strings.EqualFold(input.URLPath, path) {
			return input
		}
	}
	return nil
}

func (sm *StreamManager) GetInputByName(name string) *InputCfg {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if input, ok := sm.inputs[name]; ok {
		return input
	}
	return nil
}

func (sm *StreamManager) SetStatusActive(name string, active bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if s, ok := sm.status[name]; ok {
		s.Active = active
		if active {
			s.Connections++
		} else if s.Connections > 0 {
			s.Connections--
		}
	}
}

func (sm *StreamManager) IncrementError(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if s, ok := sm.status[name]; ok {
		s.ErrorCount++
	}
}

func (sm *StreamManager) GetStatus(name string) *StreamStatus {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if s, ok := sm.status[name]; ok {
		copy := *s
		copy.Outputs = sm.GetOutputsStatus(name)
		return &copy
	}
	return nil
}

func (sm *StreamManager) GetAllStatuses() []*StreamStatus {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	var list []*StreamStatus
	for _, s := range sm.status {
		copy := *s
		copy.Outputs = sm.GetOutputsStatus(s.Name)
		list = append(list, &copy)
	}
	return list
}

// Управление статусами выходов

func (sm *StreamManager) RegisterOutput(inputName, url string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.outputs[inputName] == nil {
		sm.outputs[inputName] = make(map[string]*OutputStatus)
	}
	if _, exists := sm.outputs[inputName][url]; !exists {
		sm.outputs[inputName][url] = &OutputStatus{URL: url}
	}
}

func (sm *StreamManager) SetOutputActive(inputName, url string, active bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if outMap, ok := sm.outputs[inputName]; ok {
		if out, ok2 := outMap[url]; ok2 {
			out.Active = active
			if active {
				out.startTime = time.Now()
			} else {
				out.startTime = time.Time{}
				out.Uptime = "00:00:00"
				out.BitrateKbps = 0
			}
		}
	}
}

func (sm *StreamManager) IncrementOutputError(inputName, url string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if outMap, ok := sm.outputs[inputName]; ok {
		if out, ok2 := outMap[url]; ok2 {
			out.ErrorCount++
			out.Active = false
		}
	}
}

// Корректный подсчёт битрейта
func (sm *StreamManager) UpdateOutputBitrate(inputName, url string, bytes int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	outMap, ok := sm.outputs[inputName]
	if !ok {
		log.Printf("[Bitrate] No outputs for input %s", inputName)
		return
	}
	out, ok2 := outMap[url]
	if !ok2 {
		// Не логируем, если выход уже удалён — это штатная ситуация при динамическом удалении
		return
	}

	now := time.Now()

	if out.prevTime.IsZero() {
		out.prevBytes = bytes
		out.prevTime = now
		out.BitrateKbps = 0
		return
	}

	deltaBytes := bytes - out.prevBytes
	deltaSeconds := now.Sub(out.prevTime).Seconds()

	if deltaSeconds < 0.5 {
		return
	}

	if deltaBytes < 0 {
		log.Printf("[Bitrate] Bytes counter reset detected for input=%s url=%s", inputName, url)
		out.prevBytes = bytes
		out.prevTime = now
		out.BitrateKbps = 0
		return
	}

	bitrate := float64(deltaBytes*8) / deltaSeconds / 1000.0
	if bitrate < 0 {
		bitrate = 0
	}

	out.BitrateKbps = bitrate
	out.prevBytes = bytes
	out.prevTime = now

	// log.Printf("[Bitrate] Calculated bitrate for input=%s url=%s: %.3f kbps", inputName, url, bitrate)
}

func (sm *StreamManager) GetOutputsStatus(inputName string) []*OutputStatus {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	var list []*OutputStatus
	if outMap, ok := sm.outputs[inputName]; ok {
		for _, out := range outMap {
			copy := *out
			// Сбрасываем внутренние поля при отдаче статуса
			copy.prevBytes = 0
			copy.prevTime = time.Time{}
			if out.Active && !out.startTime.IsZero() {
				d := time.Since(out.startTime)
				h := int(d.Hours())
				m := int(d.Minutes()) % 60
				s := int(d.Seconds()) % 60
				copy.Uptime = fmt.Sprintf("%02d:%02d:%02d", h, m, s)
			} else {
				copy.Uptime = "00:00:00"
			}
			list = append(list, &copy)
		}
	}
	return list
}

func (sm *StreamManager) ForceReconnectOutput(inputName, url string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if outMap, ok := sm.outputs[inputName]; ok {
		if out, ok2 := outMap[url]; ok2 {
			out.Active = false
			out.startTime = time.Time{}
			out.Uptime = "00:00:00"
			out.ErrorCount++
			// Здесь можно добавить дополнительную логику для сигнала перезапуска, если потребуется
		}
	}
}

// Удаляет выходы, которых нет в inputCfg.Outputs
func (sm *StreamManager) CleanupRemovedOutputs(inputName string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	outMap, ok := sm.outputs[inputName]
	if !ok {
		return
	}
	input, ok := sm.inputs[inputName]
	if !ok || input == nil {
		return
	}
	// Собираем set актуальных выходов
	actual := make(map[string]struct{})
	for _, url := range input.Outputs {
		actual[url] = struct{}{}
	}
	for url := range outMap {
		if _, found := actual[url]; !found {
			delete(outMap, url)
		}
	}
}

// Методы для работы с глобальными настройками

func (sm *StreamManager) GetGlobalSettings() *Config {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if sm.config == nil {
		return &Config{}
	}
	return sm.config
}

func (sm *StreamManager) UpdateGlobalSettings(srtSettings SRTSettings) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.config == nil {
		sm.config = &Config{}
	}
	sm.config.SRTSettings = srtSettings
	log.Printf("[API] Global SRT settings updated: latency=%d, passphrase=%s, streamid=%s",
		srtSettings.Latency, srtSettings.Passphrase, srtSettings.StreamID)
}

func (sm *StreamManager) ReloadSettingsFromConfig() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Перезагружаем конфигурацию из файла
	newConfig, err := LoadConfig("config.yaml")
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	sm.config = newConfig
	log.Printf("[API] Global settings reloaded from config.yaml")
	return nil
}

// Методы для SRT сервера
func (sm *StreamManager) GetOrCreateStream(name, urlPath string) *Stream {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Проверяем, существует ли уже поток
	if status, exists := sm.status[name]; exists {
		status.Active = true
		status.Connections++
		// Возвращаем существующий поток или создаем новый
		return &Stream{
			Name:    name,
			URLPath: urlPath,
		}
	}

	// Создаем новый статус
	sm.status[name] = &StreamStatus{
		Name:        name,
		URLPath:     urlPath,
		Active:      true,
		Connections: 1,
	}

	return &Stream{
		Name:    name,
		URLPath: urlPath,
	}
}

func (sm *StreamManager) RemoveStream(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if status, exists := sm.status[name]; exists {
		status.Active = false
		if status.Connections > 0 {
			status.Connections--
		}
	}
}

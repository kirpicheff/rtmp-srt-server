package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type APIServer struct {
	SM       *StreamManager
	User     string
	Password string
}

func NewAPIServer(sm *StreamManager, user, password string) *APIServer {
	return &APIServer{SM: sm, User: user, Password: password}
}

func (api *APIServer) basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if api.User == "" && api.Password == "" {
			h(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Basic ") {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(string(payload), ":", 2)
		if len(parts) != 2 || parts[0] != api.User || parts[1] != api.Password {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h(w, r)
	}
}

func (api *APIServer) routes() http.Handler {
	mux := http.NewServeMux()

	// Статические файлы веб-интерфейса
	mux.HandleFunc("/", api.handleWebInterface)

	// API маршруты
	mux.HandleFunc("/api/inputs", api.basicAuth(api.handleListInputs))                      // GET
	mux.HandleFunc("/api/inputs/add", api.basicAuth(api.handleAddInput))                    // POST
	mux.HandleFunc("/api/inputs/remove", api.basicAuth(api.handleRemoveInput))              // GET ?name=
	mux.HandleFunc("/api/inputs/update_outputs", api.basicAuth(api.handleUpdateOutputs))    // POST
	mux.HandleFunc("/api/status", api.basicAuth(api.handleGetStatus))                       // GET ?name=
	mux.HandleFunc("/api/status/all", api.basicAuth(api.handleGetAllStatuses))              // GET
	mux.HandleFunc("/api/outputs/reconnect", api.basicAuth(api.handleForceReconnectOutput)) // POST
	mux.HandleFunc("/api/outputs/add", api.basicAuth(api.handleAddOutput))                  // POST
	mux.HandleFunc("/api/outputs/remove", api.basicAuth(api.handleRemoveOutput))            // POST
	mux.HandleFunc("/api/settings", api.basicAuth(api.handleSettings))                      // GET/PUT
	mux.HandleFunc("/api/settings/reload", api.basicAuth(api.handleReloadSettings))         // POST

	return mux
}

// --- Handlers ---

func (api *APIServer) handleListInputs(w http.ResponseWriter, r *http.Request) {
	inputs := api.SM.ListInputs()
	writeJSON(w, inputs)
}

func (api *APIServer) handleAddInput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024) // Ограничиваем тело запроса 1 МБ

	var input InputCfg
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Name == "" || input.URLPath == "" {
		http.Error(w, "Input name and URL path are required", http.StatusBadRequest)
		return
	}

	// Проверка на дублирование имени
	if existing := api.SM.GetInputByName(input.Name); existing != nil {
		http.Error(w, "Input with this name already exists", http.StatusConflict)
		return
	}

	// Валидация всех выходов
	for _, outURL := range input.Outputs {
		if err := validateRTMPURL(outURL); err != nil {
			http.Error(w, "Invalid output URL: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	api.SM.AddInput(input)
	w.WriteHeader(http.StatusCreated)
	log.Printf("[API] New input added: %s", input.Name)
}

func (api *APIServer) handleRemoveInput(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing 'name' parameter", http.StatusBadRequest)
		return
	}
	api.SM.RemoveInput(name)
	w.WriteHeader(http.StatusOK)
}

// Обновление списка выходов у входа
func (api *APIServer) handleUpdateOutputs(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string   `json:"name"`
		Outputs []string `json:"outputs"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	api.SM.mu.RLock()
	input, ok := api.SM.inputs[req.Name]
	api.SM.mu.RUnlock()
	if !ok {
		http.Error(w, "Input not found", http.StatusNotFound)
		return
	}

	// Обновляем Outputs (указатель, потокобезопасно)
	input.Outputs = req.Outputs

	// Регистрируем новые выходы
	for _, url := range req.Outputs {
		api.SM.RegisterOutput(req.Name, url)
	}
	// Очищаем неактуальные выходы
	api.SM.CleanupRemovedOutputs(req.Name)

	log.Printf("[API] Updated outputs for input %s: %v", req.Name, req.Outputs)
	w.WriteHeader(http.StatusOK)
}

func (api *APIServer) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing 'name' parameter", http.StatusBadRequest)
		return
	}
	status := api.SM.GetStatus(name)
	if status == nil {
		http.Error(w, "Stream not found", http.StatusNotFound)
		return
	}
	writeJSON(w, status)
}

func (api *APIServer) handleGetAllStatuses(w http.ResponseWriter, r *http.Request) {
	statuses := api.SM.GetAllStatuses()
	writeJSON(w, statuses)
}

func (api *APIServer) handleForceReconnectOutput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.URL == "" {
		http.Error(w, "Missing name or url", http.StatusBadRequest)
		return
	}
	go api.SM.ForceReconnectOutput(req.Name, req.URL)
	w.WriteHeader(http.StatusOK)
}

func (api *APIServer) handleAddOutput(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	api.SM.mu.RLock()
	input, ok := api.SM.inputs[req.Name]
	api.SM.mu.RUnlock()
	if !ok {
		http.Error(w, "Input not found", http.StatusNotFound)
		return
	}
	// Проверяем, есть ли уже такой выход
	for _, out := range input.Outputs {
		if out == req.URL {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	// Добавляем выход в память
	input.Outputs = append(input.Outputs, req.URL)
	api.SM.RegisterOutput(req.Name, req.URL)
	// Обновляем config.yaml
	go updateConfigYAML(api.SM.inputs)
	w.WriteHeader(http.StatusOK)
}

func (api *APIServer) handleRemoveOutput(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	api.SM.mu.RLock()
	input, ok := api.SM.inputs[req.Name]
	api.SM.mu.RUnlock()
	if !ok {
		http.Error(w, "Input not found", http.StatusNotFound)
		return
	}
	// Удаляем выход из памяти
	newOuts := make([]string, 0, len(input.Outputs))
	for _, out := range input.Outputs {
		if out != req.URL {
			newOuts = append(newOuts, out)
		}
	}
	input.Outputs = newOuts
	// Очищаем неактуальные выходы
	api.SM.CleanupRemovedOutputs(req.Name)
	// Обновляем config.yaml
	go updateConfigYAML(api.SM.inputs)
	w.WriteHeader(http.StatusOK)
}

func (api *APIServer) handleSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Получить текущие настройки
		api.SM.mu.RLock()
		settings := api.SM.GetGlobalSettings()
		api.SM.mu.RUnlock()
		writeJSON(w, settings)

	case http.MethodPut:
		// Обновить настройки
		var req struct {
			SRTSettings       SRTSettings `json:"srt_settings,omitempty"`
			LogToFile         bool        `json:"log_to_file,omitempty"`
			LogFile           string      `json:"log_file,omitempty"`
			ReconnectInterval int         `json:"reconnect_interval,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Валидация
		if req.ReconnectInterval < 1 {
			http.Error(w, "Reconnect interval must be >= 1", http.StatusBadRequest)
			return
		}
		if req.SRTSettings.Latency < 0 {
			http.Error(w, "Latency must be >= 0", http.StatusBadRequest)
			return
		}
		if req.SRTSettings.ConnectTimeout < 0 {
			http.Error(w, "Connect timeout must be >= 0", http.StatusBadRequest)
			return
		}

		// Обновляем настройки
		api.SM.mu.Lock()
		if req.SRTSettings != (SRTSettings{}) {
			api.SM.UpdateGlobalSettings(req.SRTSettings)
		}
		if req.LogToFile != api.SM.config.LogToFile {
			api.SM.config.LogToFile = req.LogToFile
			log.Printf("[API] LogToFile updated: %v", req.LogToFile)
		}
		if req.LogFile != "" && req.LogFile != api.SM.config.LogFile {
			api.SM.config.LogFile = req.LogFile
			log.Printf("[API] LogFile updated: %s", req.LogFile)
		}
		if req.ReconnectInterval > 0 && req.ReconnectInterval != api.SM.config.ReconnectInterval {
			api.SM.config.ReconnectInterval = req.ReconnectInterval
			log.Printf("[API] ReconnectInterval updated: %d", req.ReconnectInterval)
		}
		api.SM.mu.Unlock()

		// Получаем копию актуальных настроек
		settings := api.SM.GetGlobalSettings()
		// Сохраняем в config.yaml
		go api.saveSettingsToConfig(settings)

		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (api *APIServer) handleReloadSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Перезагружаем настройки из config.yaml
	if err := api.SM.ReloadSettingsFromConfig(); err != nil {
		http.Error(w, "Failed to reload settings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *APIServer) saveSettingsToConfig(settings *Config) {
	// Читаем текущий config.yaml
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("[API] Failed to read config.yaml: %v", err)
		return
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Printf("[API] Failed to parse config.yaml: %v", err)
		return
	}

	// Обновляем только безопасные настройки (не трогаем ListenPort, APIAuthUser, APIAuthPassword)
	cfg.SRTSettings = settings.SRTSettings
	cfg.LogToFile = settings.LogToFile
	cfg.LogFile = settings.LogFile
	cfg.ReconnectInterval = settings.ReconnectInterval

	// Сохраняем обратно
	out, err := yaml.Marshal(&cfg)
	if err != nil {
		log.Printf("[API] Failed to marshal config.yaml: %v", err)
		return
	}
	if err := os.WriteFile("config.yaml", out, 0644); err != nil {
		log.Printf("[API] Failed to write config.yaml: %v", err)
		return
	}
	log.Printf("[API] Global settings saved to config.yaml")
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// Обновляет config.yaml на лету (через структуру Config)
func updateConfigYAML(inputs map[string]*InputCfg) {
	// Читаем текущий config.yaml в структуру
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("[API] Failed to read config.yaml: %v", err)
		return
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Printf("[API] Failed to parse config.yaml: %v", err)
		return
	}
	// Обновляем outputs для каждого входа
	for i := range cfg.Inputs {
		name := cfg.Inputs[i].Name
		if in, ok := inputs[name]; ok {
			cfg.Inputs[i].Outputs = in.Outputs
		}
	}
	// Сохраняем обратно
	out, err := yaml.Marshal(&cfg)
	if err != nil {
		log.Printf("[API] Failed to marshal config.yaml: %v", err)
		return
	}
	if err := os.WriteFile("config.yaml", out, 0644); err != nil {
		log.Printf("[API] Failed to write config.yaml: %v", err)
		return
	}
	log.Printf("[API] config.yaml updated (pretty)")
}

// Обработчик веб-интерфейса
func (api *APIServer) handleWebInterface(w http.ResponseWriter, r *http.Request) {
	// Если запрос к API - пропускаем
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}

	// Basic auth для web-интерфейса
	if api.User != "" || api.Password != "" {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Basic ") {
			w.Header().Set("WWW-Authenticate", `Basic realm="WebUI"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(string(payload), ":", 2)
		if len(parts) != 2 || parts[0] != api.User || parts[1] != api.Password {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Определяем путь к файлу
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	// Читаем файл из папки web
	filePath := "web" + path
	data, err := os.ReadFile(filePath)
	if err != nil {
		// Если файл не найден, возвращаем index.html
		if path != "/index.html" {
			data, err = os.ReadFile("web/index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
		} else {
			http.NotFound(w, r)
			return
		}
	}

	// Определяем Content-Type
	contentType := "text/plain"
	switch {
	case strings.HasSuffix(path, ".html"):
		contentType = "text/html"
	case strings.HasSuffix(path, ".js"):
		contentType = "application/javascript"
	case strings.HasSuffix(path, ".css"):
		contentType = "text/css"
	case strings.HasSuffix(path, ".json"):
		contentType = "application/json"
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

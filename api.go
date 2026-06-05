package main

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

//go:embed web/*
var webFS embed.FS

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

	// Health check endpoint (без аутентификации)
	mux.HandleFunc("/health", api.handleHealthCheck)

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

// Health check endpoint для мониторинга
func (api *APIServer) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем статус всех входов
	allStatuses := api.SM.GetAllStatuses()

	// Подсчитываем статистику
	var totalInputs, activeInputs, totalOutputs, activeOutputs int
	for _, status := range allStatuses {
		totalInputs++
		if status.Active {
			activeInputs++
		}
		for _, output := range status.Outputs {
			totalOutputs++
			if output.Active {
				activeOutputs++
			}
		}
	}

	// Формируем ответ
	healthStatus := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"inputs": map[string]interface{}{
			"total":  totalInputs,
			"active": activeInputs,
		},
		"outputs": map[string]interface{}{
			"total":  totalOutputs,
			"active": activeOutputs,
		},
	}

	// Если нет активных входов, считаем сервер "нездоровым"
	if activeInputs == 0 {
		healthStatus["status"] = "warning"
		healthStatus["message"] = "No active inputs"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthStatus)
}

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

	// Обновляем Outputs потокобезопасно
	if ok := api.SM.UpdateInputOutputs(req.Name, req.Outputs); !ok {
		http.Error(w, "Input not found", http.StatusNotFound)
		return
	}

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

	ok, alreadyExists := api.SM.AddOutputToInput(req.Name, req.URL)
	if !ok {
		http.Error(w, "Input not found", http.StatusNotFound)
		return
	}
	if alreadyExists {
		w.WriteHeader(http.StatusOK)
		return
	}

	api.SM.RegisterOutput(req.Name, req.URL)
	// Обновляем config.yaml, используя безопасную копию
	go updateInputsInConfig(api.SM.GetInputsCopy())
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

	if ok := api.SM.RemoveOutputFromInput(req.Name, req.URL); !ok {
		http.Error(w, "Input not found", http.StatusNotFound)
		return
	}

	// Очищаем неактуальные выходы
	api.SM.CleanupRemovedOutputs(req.Name)
	// Обновляем config.yaml, используя безопасную копию
	go updateInputsInConfig(api.SM.GetInputsCopy())
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

// updateInputsInConfig обновляет только секцию inputs в config.yaml, сохраняя комментарии.
func updateInputsInConfig(inputs map[string]*InputCfg) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("[API] Failed to read config.yaml for update: %v", err)
		return
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		log.Printf("[API] Failed to parse config.yaml for update: %v", err)
		return
	}

	// Находим узел 'inputs'
	var inputsNode *yaml.Node
	if root.Content[0].Kind == yaml.MappingNode {
		for i := 0; i < len(root.Content[0].Content); i += 2 {
			if root.Content[0].Content[i].Value == "inputs" {
				inputsNode = root.Content[0].Content[i+1]
				break
			}
		}
	}

	if inputsNode == nil || inputsNode.Kind != yaml.SequenceNode {
		log.Printf("[API] 'inputs' section not found or not a sequence in config.yaml")
		return
	}

	// Обновляем дочерние узлы
	for _, inputNode := range inputsNode.Content {
		if inputNode.Kind != yaml.MappingNode {
			continue
		}
		var name string
		var outputsNode *yaml.Node
		for i := 0; i < len(inputNode.Content); i += 2 {
			keyNode := inputNode.Content[i]
			valNode := inputNode.Content[i+1]
			if keyNode.Value == "name" {
				name = valNode.Value
			}
			if keyNode.Value == "outputs" {
				outputsNode = valNode
			}
		}

		if name != "" && outputsNode != nil {
			if memInput, ok := inputs[name]; ok {
				// Создаем новый узел для outputs
				newOutputsNode := yaml.Node{Kind: yaml.SequenceNode}
				for _, outURL := range memInput.Outputs {
					newOutputsNode.Content = append(newOutputsNode.Content, &yaml.Node{
						Kind:  yaml.ScalarNode,
						Value: outURL,
					})
				}
				*outputsNode = newOutputsNode
			}
		}
	}

	// Сохраняем обновленный YAML
	file, err := os.Create("config.yaml")
	if err != nil {
		log.Printf("[API] Failed to open config.yaml for writing: %v", err)
		return
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2) // Сохраняем форматирование
	if err := encoder.Encode(&root); err != nil {
		log.Printf("[API] Failed to write updated config.yaml: %v", err)
		return
	}

	log.Printf("[API] config.yaml updated preserving comments.")
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

	// Читаем файл из встроенной папки web
	filePath := "web" + path
	data, err := webFS.ReadFile(filePath)
	if err != nil {
		// Если файл не найден, возвращаем index.html
		if path != "/index.html" {
			data, err = webFS.ReadFile("web/index.html")
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

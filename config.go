package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"

	"gopkg.in/yaml.v3"
)

type SRTSettings struct {
	ConnectTimeout int    `yaml:"connect_timeout"`
	Latency        int    `yaml:"latency"`
	Encryption     string `yaml:"encryption"`
	Passphrase     string `yaml:"passphrase"`
	StreamID       string `yaml:"streamid"`
}

type ServerConfig struct {
	Port            int    `yaml:"port"`
	RTMPPort        int    `yaml:"rtmp_port"`
	SRTPort         int    `yaml:"srt_port"`
	WHIPPort        int    `yaml:"whip_port"`
	APIAuthUser     string `yaml:"api_username"`
	APIAuthPassword string `yaml:"api_password"`
}

type WHIPSettings struct {
	ICEServers []string `yaml:"ice_servers"`
}

type Config struct {
	Server            ServerConfig `yaml:"server"`
	SRTPort           int          `yaml:"srt_port"`
	ReconnectInterval int          `yaml:"reconnect_interval"`
	LogToFile         bool         `yaml:"log_to_file"`
	LogFile           string       `yaml:"log_file"`
	MinimizeToTray    bool         `yaml:"minimize_to_tray"`
	Inputs            []InputCfg   `yaml:"inputs"`
	APIAuthUser       string       `yaml:"api_auth_user" json:"-"`
	APIAuthPassword   string       `yaml:"api_auth_password" json:"-"`
	SRTSettings       SRTSettings  `yaml:"srt_settings"`
	WHIPSettings      WHIPSettings `yaml:"whip_settings"`
}

type InputCfg struct {
	Name    string   `yaml:"name"`
	URLPath string   `yaml:"url_path"`
	Outputs []string `yaml:"outputs"`
}

func LoadConfig(path string) (*Config, error) {
	// Создаем дефолтный конфиг, если файла нет
	if _, err := ioutil.ReadFile(path); err != nil {
		defaultConfig := `server:
  port: 8080
  rtmp_port: 1935
  srt_port: 9000
  whip_port: 8084
  api_username: admin
  api_password: secret
srt_port: 0
reconnect_interval: 5
log_to_file: true
log_file: server.log
minimize_to_tray: false
inputs:
  - name: obs
    url_path: /live/stream
    outputs: []
api_auth_user: admin
api_auth_password: secret
srt_settings:
  connect_timeout: 5000
  latency: 200
  encryption: none
  passphrase: ""
  streamid: ""`
		if errWrite := ioutil.WriteFile(path, []byte(defaultConfig), 0644); errWrite != nil {
			return nil, fmt.Errorf("failed to create default config: %w", errWrite)
		}
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (cfg *Config) Validate() error {
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return errors.New("server.port must be between 1 and 65535")
	}
	if cfg.Server.RTMPPort <= 0 || cfg.Server.RTMPPort > 65535 {
		return errors.New("server.rtmp_port must be between 1 and 65535")
	}
	if cfg.Server.SRTPort != 0 && (cfg.Server.SRTPort < 1 || cfg.Server.SRTPort > 65535) {
		return errors.New("server.srt_port must be between 1 and 65535")
	}
	if cfg.ReconnectInterval < 1 {
		return errors.New("reconnect_interval must be > 0")
	}
	seenPaths := make(map[string]struct{})
	for _, input := range cfg.Inputs {
		if input.Name == "" {
			return errors.New("input name cannot be empty")
		}
		if input.URLPath == "" || input.URLPath[0] != '/' {
			return fmt.Errorf("invalid url_path in input %s", input.Name)
		}
		if _, ok := seenPaths[input.URLPath]; ok {
			return fmt.Errorf("duplicate url_path: %s", input.URLPath)
		}
		seenPaths[input.URLPath] = struct{}{}

		for _, out := range input.Outputs {
			if _, err := url.ParseRequestURI(out); err != nil {
				return fmt.Errorf("invalid output URL '%s' in input %s", out, input.Name)
			}
		}
	}

	return nil
}

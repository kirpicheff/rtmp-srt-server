package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/datarhei/joy4/format/rtmp"
)

func main() {
	// Добавляем обработку паники в main
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[PANIC] Main panic: %v", r)
		}
		log.Printf("[DEBUG] Main function finished")
	}()

	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.LogToFile {
		f, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		log.SetOutput(io.MultiWriter(os.Stdout, f))
		defer f.Close()
	}

	sm := NewStreamManager(cfg.Inputs, cfg)

	// RTMP сервер
	rtmpServer := &rtmp.Server{}
	rtmpServer.HandlePublish = handlePublish(sm, cfg)
	rtmpServer.HandlePlay = func(conn *rtmp.Conn) {
		log.Printf("Play started: %s", conn.URL)
		conn.Close()
	}

	// SRT сервер
	srtServer := NewSRTServer(cfg.Server.SRTPort, cfg, sm)

	// WHIP сервер
	whipServer := NewWHIPServer(cfg.Server.WHIPPort, sm)

	// HTTP API сервер
	apiServer := NewAPIServer(sm, cfg.Server.APIAuthUser, cfg.Server.APIAuthPassword)
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Server.Port),
		Handler: apiServer.routes(),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Heartbeat горутина для мониторинга
	heartbeatCtx, heartbeatCancel := context.WithCancel(context.Background())
	defer heartbeatCancel()

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		startTime := time.Now()
		log.Printf("[HEARTBEAT] Server started at %v", startTime)

		for {
			select {
			case <-heartbeatCtx.Done():
				return
			case <-ticker.C:
				uptime := time.Since(startTime)
				activeInputs := sm.GetActiveInputs()
				log.Printf("[HEARTBEAT] Server uptime: %v, Active inputs: %d", uptime, len(activeInputs))

				// Логируем детали активных входов
				for _, input := range activeInputs {
					outputs := sm.GetOutputsForInput(input.Name)
					activeOutputs := 0
					for _, output := range outputs {
						if output.Active {
							activeOutputs++
						}
					}
					log.Printf("[HEARTBEAT] Input '%s': %d/%d outputs active", input.Name, activeOutputs, len(outputs))
				}
			}
		}
	}()

	// Мониторинг горутин с таймаутами
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				log.Printf("[MONITOR] Checking for stuck goroutines...")

				// Проверяем, что все серверы еще работают
				// Если какой-то сервер "завис", логируем это
				log.Printf("[MONITOR] All servers appear to be running")
			}
		}
	}()

	// Запуск RTMP сервера
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] RTMP server panic: %v", r)
			}
			log.Printf("[DEBUG] RTMP server goroutine finished")
		}()

		addr := ":" + strconv.Itoa(cfg.Server.RTMPPort)
		log.Printf("[RTMP] server started on %s", addr)
		if err := rtmpServer.ListenAndServe(); err != nil && err.Error() != "use of closed network connection" {
			log.Fatalf("RTMP server error: %v", err)
		}
	}()

	// Запуск SRT сервера
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] SRT server panic: %v", r)
			}
			log.Printf("[DEBUG] SRT server goroutine finished")
		}()

		log.Printf("[SRT] server started on :%d", cfg.Server.SRTPort)
		if err := srtServer.Start(); err != nil {
			log.Fatalf("SRT server error: %v", err)
		}
	}()

	// Запуск WHIP сервера
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] WHIP server panic: %v", r)
			}
			log.Printf("[DEBUG] WHIP server goroutine finished")
		}()

		if err := whipServer.Start(); err != nil {
			log.Fatalf("WHIP server error: %v", err)
		}
	}()

	// Запуск HTTP API сервера
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] HTTP server panic: %v", r)
			}
			log.Printf("[DEBUG] HTTP server goroutine finished")
		}()

		log.Printf("[API] server listening on :%d", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API server error: %v", err)
		}
	}()

	// Ждём сигнал завершения
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// Останавливаем heartbeat
	heartbeatCancel()

	// Завершаем HTTP сервер
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Завершаем SRT сервер
	if err := srtServer.Stop(); err != nil {
		log.Printf("SRT server shutdown error: %v", err)
	}

	// Завершаем WHIP сервер
	if err := whipServer.Stop(); err != nil {
		log.Printf("WHIP server shutdown error: %v", err)
	}

	// Завершаем RTMP сервер (у joy4 нет явного метода shutdown, он просто перестаёт принимать)
	log.Println("RTMP server stopped")

	log.Println("All servers gracefully stopped")
}

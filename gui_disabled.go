//go:build !gui

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// startGUI в консольном режиме просто ждет сигналов SIGINT/SIGTERM для корректного завершения
func startGUI(url string, cfg *Config, stop func()) {
	log.Printf("[CONSOLE] Running in headless mode. GUI is disabled.")
	
	// Ждем системный сигнал завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	select {
	case sig := <-sigChan:
		log.Printf("[CONSOLE] Received signal: %v. Shutting down...", sig)
	}
	
	stop()
}

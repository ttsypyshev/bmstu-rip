package main

import (
	"log"
	"os"
	"os/signal"
	"rip/src/backend"
	"syscall"
)

// @title BITOP
// @version 1.0
// @description Bmstu Open IT Platform

// @contact.name ttsypyshev
// @contact.url https://vk.com/ttsypyshev
// @contact.email ttsypyshev01@gmail.com

// @license.name AS IS (NO WARRANTY)

// @host localhost:8080
// @schemes http
// @BasePath /

func main() {
	log.Println("App start")

	err := backend.Migrate()
	if err != nil {
		log.Fatalf("Failed to migrate the database: %v", err)
		return
	}

	go func() {
		if err := backend.Run(); err != nil {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("App down")
}

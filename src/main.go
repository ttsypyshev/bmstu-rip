package main

import (
	"log"
	"os"
	"os/signal"
	"rip/src/backend"
	"syscall"
)

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

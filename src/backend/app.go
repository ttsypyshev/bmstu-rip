package backend

import (
	"log"

	"github.com/gin-gonic/gin"
)

type App struct {
	db *Db
	userID uint
}

func Run() error {
	log.Println("Server starting up")

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	app, err := NewDB(FromEnv())
	if err != nil {
		log.Fatalf("Error initializing the database: %v", err)
		return err
	}

	app.userID = 1

	app.SetupRoutes(r)
	if err := r.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
		return err
	}

	log.Println("Server stopped")
	return nil
}

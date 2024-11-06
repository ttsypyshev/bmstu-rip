package backend

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type App struct {
	db          *Db
	ModeratorID uint
	userID      uint
	minioClient *minio.Client
}

func Run() error {
	log.Println("Server starting up")

	r := gin.Default()

	app, err := NewDB(FromEnvDB())
	if err != nil {
		log.Fatalf("Error initializing the database: %v", err)
		return err
	}

	app.minioClient, err = InitializeMinIO()
	if err != nil {
		log.Fatalf("Error initializing MinIO: %v", err)
		return err
	}

	app.ModeratorID = 1
	app.userID = 2

	app.SetupRoutes(r)
	if err := r.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
		return err
	}

	log.Println("Server stopped")
	return nil
}

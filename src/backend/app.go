package backend

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"

	env "rip/pkg/settings"
)

type App struct {
	db          *Db
	minioClient *minio.Client
	redisClient *redis.Client
	secret      string
}

func Run() error {
	log.Println("Server starting up")

	r := gin.Default()

	app, err := NewDB(env.FromEnvDB())
	if err != nil {
		log.Fatalf("Error initializing the database: %v", err)
		return err
	}

	app.minioClient, err = InitializeMinIO()
	if err != nil {
		log.Fatalf("Error initializing MinIO: %v", err)
		return err
	}

	app.redisClient, err = InitializeRedis()
	if err != nil {
		log.Fatalf("Error initializing Redis: %v", err)
		return err
	}

	app.SetupRoutes(r)
	if err := r.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
		return err
	}

	log.Println("Server stopped")
	return nil
}

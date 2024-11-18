package backend

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"

	"rip/pkg/database"
	env "rip/pkg/settings"
)

type App struct {
	db          *Db
	minioClient *minio.Client
	redisClient *redis.Client
	userID      *uuid.UUID
	role        *database.Role
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
	r.GET("/ping/:name", app.Ping)
	if err := r.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
		return err
	}

	log.Println("Server stopped")
	return nil
}

type pingReq struct{}
type pingResp struct {
	Status string `json:"status"`
}

// Ping godoc
// @Summary      Show hello text
// @Description  very very friendly response
// @Tags         Tests
// @Produce      json
// @Success      200  {object}  pingResp
// @Router       /ping/{name} [get]
func (a *App) Ping(gCtx *gin.Context) {
	name := gCtx.Param("name")
	gCtx.String(http.StatusOK, "Hello %s", name)
}

package backend

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Парсинг списка в карту
func ParseList(listStr string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(listStr, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, ":") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0]) + ":"
		value := strings.TrimSpace(parts[1])
		result[key] = value
	}
	return result
}

func extractObjectNameFromURL(url string) string {
	// Извлекаем имя объекта из URL (например, "service_images/filename.jpg")
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

type ErrorResponse struct {
	Message string `json:"message" example:"[err] invalid request format"`
	Status  bool   `json:"status" example:"false"`
}

// handleError обрабатывает и логирует ошибки, отправляет ответ
func handleError(c *gin.Context, statusCode int, err error, additionalErrs ...error) {
	if err == nil {
		log.Panic("Logging error: empty main error argument")
		return
	}

	c.JSON(statusCode, ErrorResponse{Status: false, Message: err.Error()})

	var errorMessages strings.Builder
	errorMessages.WriteString(err.Error())
	for _, additionalErr := range additionalErrs {
		if additionalErr != nil {
			errorMessages.WriteString(": " + additionalErr.Error())
		}
	}

	log.Printf("Error: %s", errorMessages.String())
}

// func ParseQueryParam(c *gin.Context, key string) (int, error) {
// 	param := c.Query(key)
// 	if param == "" {
// 		return 0, nil
// 	}
// 	return strconv.Atoi(param)
// }

func (app *App) GetFilteredLangs(query string) ([]DbLang, error) {
	if query != "" {
		return app.filterLangsByQuery(query)
	}
	return app.getLangs(func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", true)
	})
}

func ExtractUserID(c *gin.Context) (uuid.UUID, error) {
	idAny, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, errors.New("userID is missing from context")
	}

	idStr, ok := idAny.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("expected string for userID, but got %T", idAny)
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format for userID: %v", err)
	}

	return userID, nil
}

package backend

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
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

// handleError обрабатывает и логирует ошибки, отправляет ответ
func handleError(c *gin.Context, statusCode int, err error, additionalErrs ...error) {
	if err == nil {
		log.Panic("Logging error: empty main error argument")
		return
	}

	c.JSON(statusCode, gin.H{"status": false, "message": err.Error()})

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

// // Генерация JWT-токена (пример с использованием "github.com/golang-jwt/jwt/v4")
// func generateJWTToken(userID uint) (string, error) {
// 	// Определяем стандартные параметры токена
// 	claims := jwt.MapClaims{
// 		"user_id": userID,
// 		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Время жизни токена — 72 часа
// 	}

// 	// Создаем токен с использованием секретного ключа
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	secretKey := []byte("your-secret-key") // Замените на более безопасный способ хранения ключа
// 	return token.SignedString(secretKey)
// }

// func (app *App) blockToken(tokenString string) error {
// 	// В простом варианте — добавить токен в базу заблокированных токенов с временем истечения
// 	// Либо использовать Redis или другую систему для временного хранения токенов
// 	expirationTime := time.Now().Add(72 * time.Hour) // Токен будет храниться до истечения срока его действия

// 	blockedToken := &database.BlockedToken{
// 		Token:          tokenString,
// 		ExpirationTime: expirationTime,
// 	}

// 	return app.db.Create(blockedToken).Error
// }

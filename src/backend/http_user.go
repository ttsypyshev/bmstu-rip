package backend

// import (
// 	"errors"
// 	"log"
// 	"net/http"
// 	"rip/database"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"golang.org/x/crypto/bcrypt"
// )

// // Структура запроса для регистрации пользователя
// type RegisterUserRequest struct {
// 	Name     string `json:"name" binding:"required"`     // Имя пользователя
// 	Login    string `json:"login" binding:"required"`    // Логин пользователя
// 	Password string `json:"password" binding:"required"` // Пароль
// }

// // Функция для регистрации нового пользователя
// func (app *App) RegisterUser(c *gin.Context) {
// 	var req RegisterUserRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
// 		return
// 	}
// 	log.Printf("[info] RegisterUser called: Login=%s", req.Login)

// 	// Проверяем, существует ли пользователь с таким же логином
// 	if exists, err := app.isUserLoginTaken(req.Login); err != nil {
// 		handleError(c, http.StatusInternalServerError, errors.New("[err] database error"), err)
// 		return
// 	} else if exists {
// 		handleError(c, http.StatusConflict, errors.New("[err] login already taken"), nil)
// 		return
// 	}

// 	// Хешируем пароль
// 	hashedPassword, err := hashPassword(req.Password)
// 	if err != nil {
// 		handleError(c, http.StatusInternalServerError, errors.New("[err] password hashing error"), err)
// 		return
// 	}

// 	// Создаем пользователя
// 	user := &database.User{
// 		Name:     req.Name,
// 		Login:    req.Login,
// 		Password: hashedPassword,
// 	}

// 	// Сохраняем пользователя в базе данных
// 	if err := app.createUser(user); err != nil {
// 		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to create user"), err)
// 		return
// 	}

// 	log.Printf("[info] User registered successfully: Login=%s", req.Login)
// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "User registered successfully",
// 	})
// }

// type UpdateUserProfileRequest struct {
// 	Name     *string `json:"name"`     // Имя пользователя
// 	Password *string `json:"password"` // Пароль пользователя
// }

// func (app *App) UpdateUserProfile(c *gin.Context) {
// 	var req UpdateUserProfileRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
// 		return
// 	}

// 	// Получаем ID текущего пользователя (например, из сессии или токена)
// 	userID, err := getCurrentUserID(c)
// 	if err != nil {
// 		handleError(c, http.StatusUnauthorized, errors.New("[err] unauthorized user"), err)
// 		return
// 	}

// 	log.Printf("[info] UpdateUserProfile called: UserID=%d", userID)

// 	// Находим пользователя в базе данных
// 	user, err := app.getUserByID(userID)
// 	if err != nil {
// 		handleError(c, http.StatusNotFound, errors.New("[err] user not found"), err)
// 		return
// 	}

// 	// Обновляем только указанные поля
// 	updatedFields := make(map[string]interface{})
// 	if req.Name != nil {
// 		updatedFields["name"] = *req.Name
// 	}
// 	if req.Password != nil {
// 		hashedPassword, err := hashPassword(*req.Password)
// 		if err != nil {
// 			handleError(c, http.StatusInternalServerError, errors.New("[err] password hashing error"), err)
// 			return
// 		}
// 		updatedFields["password"] = hashedPassword
// 	}

// 	// Сохраняем обновления в базе данных
// 	if err := app.updateUserProfile(userID, updatedFields); err != nil {
// 		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update user profile"), err)
// 		return
// 	}

// 	log.Printf("[info] User profile updated successfully: UserID=%d", userID)
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "User profile updated successfully",
// 	})
// }

// type LoginRequest struct {
// 	Login    string `json:"login" binding:"required"`    // Логин пользователя
// 	Password string `json:"password" binding:"required"` // Пароль пользователя
// }

// type LoginResponse struct {
// 	Token string `json:"token"` // JWT-токен или другой аутентификационный токен
// }

// func (app *App) UserLogin(c *gin.Context) {
// 	var req LoginRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
// 		return
// 	}

// 	// Проверяем, существует ли пользователь с таким логином
// 	user, err := app.getUserByLogin(req.Login)
// 	if err != nil {
// 		handleError(c, http.StatusUnauthorized, errors.New("[err] login or password incorrect"), err)
// 		return
// 	}

// 	// Сравниваем пароль пользователя с хешем в базе данных
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
// 		handleError(c, http.StatusUnauthorized, errors.New("[err] login or password incorrect"), err)
// 		return
// 	}

// 	// Генерируем JWT-токен
// 	token, err := generateJWTToken(user.ID)
// 	if err != nil {
// 		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to generate token"), err)
// 		return
// 	}

// 	log.Printf("[info] User login successful: Login=%s", req.Login)
// 	c.JSON(http.StatusOK, LoginResponse{
// 		Token: token,
// 	})
// }

// func (app *App) UserLogout(c *gin.Context) {
// 	// Извлечение токена аутентификации (например, из заголовка Authorization)
// 	tokenString := c.GetHeader("Authorization")
// 	if tokenString == "" {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] no token provided"), nil)
// 		return
// 	}

// 	// Проверяем, что токен имеет префикс "Bearer" и извлекаем его
// 	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
// 	if tokenString == "" {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid token format"), nil)
// 		return
// 	}

// 	// Блокируем токен (например, добавляем его в список заблокированных токенов)
// 	if err := app.blockToken(tokenString); err != nil {
// 		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to logout"), err)
// 		return
// 	}

// 	log.Printf("[info] User logged out successfully")
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "User logged out successfully",
// 	})
// }

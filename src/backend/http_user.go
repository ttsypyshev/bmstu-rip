package backend

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// type RegisterUserRequest struct {
// 	Name     string `json:"name"`
// 	Email    string `json:"email"`
// 	Login    string `json:"login"`
// 	Password string `json:"password"`
// }

// func (app *App) RegisterUser(c *gin.Context) {
// 	var req RegisterUserRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
// 		return
// 	}

// 	userID, err := app.createUser(req.Login, req.Password, req.Name, req.Email)
// 	if err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] failed to save user in the database"), err)
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "User registered successfully",
// 		"userID":  userID,
// 	})
// }

// type UpdateUserProfileRequest struct {
// 	Password string `json:"password"`
// 	Email    string `json:"email"`
// 	Name     string `json:"name"`
// }

// func (app *App) UpdateUserProfile(c *gin.Context) {
// 	var req UpdateUserProfileRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
// 		return
// 	}

// 	user, err := app.getUserByID(*app.userID)
// 	if err != nil {
// 		handleError(c, http.StatusNotFound, errors.New("[err] user not found"), err)
// 		return
// 	}

// 	if req.Name != "" {
// 		user.Name = req.Name
// 	}
// 	if req.Email != "" {
// 		user.Email = &req.Email
// 	}
// 	if req.Password != "" {
// 		user.Password = []byte(req.Password)
// 	}

// 	if err := app.updateUser(&user); err != nil {
// 		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update user profile"), err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "User profile updated successfully",
// 		"user":    user,
// 	})
// }

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResp struct {
	ExpiresIn   time.Duration `json:"expires_in"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
}

func (app *App) UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	isMatch, user, err := app.matchPassword(req.Login, req.Password)
	if err != nil || !isMatch {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
		return
	}

	token, err := GenerateJWT(user.ID, string(user.Role), app.secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx := context.Background()
	expiration := time.Hour * 24
	err = SaveSession(ctx, app.redisClient, user.ID, string(user.Role), expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	response := loginResp{
		ExpiresIn:   expiration,
		AccessToken: token,
		TokenType:   "Bearer",
	}

	c.JSON(http.StatusOK, response)
}

func (app *App) UserLogout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized1"})
		return
	}

	ctx := context.Background()
	err := DeleteSession(ctx, app.redisClient, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
		return
	}

	// Удаление токена из cookie
	c.SetCookie("auth_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

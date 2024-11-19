package backend

import (
	"context"
	"errors"
	"net/http"
	"rip/pkg/auth"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (app *App) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	userID, err := app.createUser(req.Login, req.Password, req.Name, req.Email)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] failed to save user in the database"), err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"userID":  userID,
	})
}

type UpdateUserProfileRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

func (app *App) UpdateUserProfile(c *gin.Context) {
	idAny, exists := c.Get("userID")
	if !exists {
		handleError(c, http.StatusUnauthorized, errors.New("[err] Unauthorized"))
		return
	}
	requestUserID, ok := idAny.(uuid.UUID)
	if !ok {
		handleError(c, http.StatusBadRequest, errors.New("[err] Unauthorized"), errors.New("userID is not of type *uuid.UUID"))
		return
	}

	var req UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	user, err := app.getUserByID(requestUserID)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] user not found"), err)
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.Password != "" {
		user.Password = []byte(req.Password)
	}

	if err := app.updateUser(&user); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update user profile"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile updated successfully",
		"user":    user,
	})
}

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (app *App) UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request"), err)
		return
	}

	isMatch, user, err := app.matchPassword(req.Login, req.Password)
	if err != nil || !isMatch {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid login or password"), err)
		return
	}

	token, err := auth.GenerateJWT(user.ID, string(user.Role), app.secret)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] failed to generate token"), err)
		return
	}

	ctx := context.Background()
	expiration := time.Hour * 24
	err = SaveSession(ctx, app.redisClient, user.ID, string(user.Role), expiration)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] failed to save session"), err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ExpiresIn":   expiration,
		"AccessToken": token,
		"TokenType":   "Bearer",
	})
}

func (app *App) UserLogout(c *gin.Context) {
	idAny, exists := c.Get("userID")
	if !exists {
		handleError(c, http.StatusUnauthorized, errors.New("[err] Unauthorized"))
		return
	}
	requestUserID, ok := idAny.(string)
	if !ok {
		handleError(c, http.StatusBadRequest, errors.New("[err] Unauthorized"), errors.New("userID is not of type *uuid.UUID"))
		return
	}

	ctx := context.Background()
	err := DeleteSession(ctx, app.redisClient, requestUserID)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] failed to clear session"), err)
		return
	}

	// Удаление токена из cookie
	c.SetCookie("auth_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

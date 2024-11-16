package backend

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
	var req UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	user, err := app.getUserByID(app.userID)
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
		user.Password = req.Password
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
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	isMatch, user, err := app.matchPassword(req.Login, req.Password)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] login or password entered incorrectly"), err)
		return
	}
	if !isMatch {
		handleError(c, http.StatusNotFound, fmt.Errorf("[info] login or password entered incorrectly"), fmt.Errorf("log: %s, pass: %s", req.Login, req.Password))
		return
	}

	if app.userID == user.ID {
		handleError(c, http.StatusNotFound, errors.New("[info] user is already logged in"))
		return
	}

	app.userID = user.ID
	app.isAdmin = user.IsAdmin

	log.Printf("[INFO] user %d logged in. Admin: %t", user.ID, user.IsAdmin)

	c.JSON(http.StatusOK, gin.H{
		"message": "User has successfully logged in",
		"status":  true,
	})
}

type UserLogoutRequest struct {
	Login string `json:"login"`
}

func (app *App) UserLogout(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	user, err := app.findUserByLogin(req.Login)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] user not found"), err)
		return
	}

	if app.userID != user.ID {
		handleError(c, http.StatusNotFound, errors.New("[info] user are not logged in"))
		return
	}

	app.userID = 0
	app.isAdmin = false

	log.Printf("[INFO] user %d logged out. Admin: %t", user.ID, user.IsAdmin)

	c.JSON(http.StatusOK, gin.H{
		"message": "User has successfully logged out",
		"status":  true,
	})
}

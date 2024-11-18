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

	user, err := app.getUserByID(*app.userID)
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

	app.userID = &user.ID
	app.role = &user.Role

	log.Printf("[INFO] user %d logged in. Role: %s", user.ID, user.Role)

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

	if *app.userID != user.ID {
		handleError(c, http.StatusNotFound, errors.New("[info] user are not logged in"))
		return
	}

	app.userID = nil
	app.role = nil

	log.Printf("[INFO] user %d logged out. Admin: %t", user.ID, user.Role)

	c.JSON(http.StatusOK, gin.H{
		"message": "User has successfully logged out",
		"status":  true,
	})
}

// type JWTClaims struct {
// 	jwt.StandardClaims           // все что точно необходимо по RFC
// 	UserUUID           uuid.UUID `json:"user_uuid"`            // наши данные - uuid этого пользователя в базе данных
// 	Scopes             []string  `json:"scopes" json:"scopes"` // список доступов в нашей системе
// }

// type loginReq struct {
// 	Login    string `json:"login"`
// 	Password string `json:"password"`
// }

// type loginResp struct {
// 	ExpiresIn   time.Duration `json:"expires_in"`
// 	AccessToken string        `json:"access_token"`
// 	TokenType   string        `json:"token_type"`
// }

// func (a *App) Login(gCtx *gin.Context) {
// 	cfg := a.config
// 	req := &loginReq{}

// 	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
// 	if err != nil {
// 		gCtx.AbortWithError(http.StatusBadRequest, err)
// 		return
// 	}

// 	if req.Login == login && req.Password == password {
// 		// значит проверка пройдена
// 		// генерируем ему jwt
// 		token := jwt.NewWithClaims(cfg.JWT.SigningMethod, &ds.JWTClaims{
// 			StandardClaims: jwt.StandardClaims{
// 				ExpiresAt: time.Now().Add(cfg.JWT.ExpiresIn).Unix(),
// 				IssuedAt:  time.Now().Unix(),
// 				Issuer:    "bitop-admin",
// 			},
// 			UserUUID: uuid.New(), // test uuid
// 			Scopes:   []string{}, // test data
// 		})

// 		if token == nil {
// 			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))
// 			return
// 		}

// 		strToken, err := token.SignedString([]byte(cfg.JWT.Token))
// 		if err != nil {
// 			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant create str token"))
// 			return
// 		}

// 		gCtx.JSON(http.StatusOK, loginResp{
// 			ExpiresIn:   cfg.JWT.ExpiresIn,
// 			AccessToken: strToken,
// 			TokenType:   "Bearer",
// 		})
// 	}

// 	gCtx.AbortWithStatus(http.StatusForbidden) // отдаем 403 ответ в знак того что доступ запрещен
// }

// type registerReq struct {
// 	Name string `json:"name"` // лучше назвать то же самое что login
// 	Pass string `json:"pass"`
// }

// type registerResp struct {
// 	Ok bool `json:"ok"`
// }

// func (a *App) Register(gCtx *gin.Context) {
// 	req := &registerReq{}

// 	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
// 	if err != nil {
// 		gCtx.AbortWithError(http.StatusBadRequest, err)
// 		return
// 	}

// 	if req.Pass == "" {
// 		gCtx.AbortWithError(http.StatusBadRequest, fmt.Errorf("pass is empty"))
// 		return
// 	}

// 	if req.Name == "" {
// 		gCtx.AbortWithError(http.StatusBadRequest, fmt.Errorf("name is empty"))
// 		return
// 	}

// 	err = a.repo.Register(&ds.User{
// 		UUID: uuid.New(),
// 		Role: role.Buyer,
// 		Name: req.Name,
// 		Pass: generateHashString(req.Pass), // пароли делаем в хешированном виде и далее будем сравнивать хеши, чтобы их не угнали с базой вместе
// 	})
// 	if err != nil {
// 		gCtx.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}

// 	gCtx.JSON(http.StatusOK, &registerResp{
// 		Ok: true,
// 	})
// }

// func generateHashString(s string) string {
// 	h := sha1.New()
// 	h.Write([]byte(s))
// 	return hex.EncodeToString(h.Sum(nil))
// }

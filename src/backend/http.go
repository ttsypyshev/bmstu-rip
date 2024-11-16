package backend

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) SetupRoutes(r *gin.Engine) {
	// Услуги (домен: `/info`)
	r.GET("/info", app.CheckAuth(), app.GetServiceList)                                  // Получить список услуг с фильтрацией, включая ID заявки-черновика пользователя и количество услуг в этой заявке.
	r.GET("/info/:id", app.CheckAuth(), app.GetServiceByID)                              // Получить данные конкретной услуги по ее ID.
	r.POST("/info", app.CheckAuth(), app.CheckAdmin(), app.CreateService)                // Добавить новую услугу (без изображения).
	r.PUT("/info/:id", app.CheckAuth(), app.CheckAdmin(), app.UpdateService)             // Изменить данные услуги по ее ID.
	r.POST("/info/:id", app.CheckAuth(), app.CheckAdmin(), app.UpdateServiceImage)       // Добавить или заменить изображение для услуги с указанным ID. Если изображение уже существует, оно заменяется.
	r.DELETE("/info/:id", app.CheckAuth(), app.CheckAdmin(), app.DeleteService)          // Удалить услугу вместе с изображением.
	r.POST("/info/add-service", app.CheckAuth(), app.CheckUser(), app.AddServiceToDraft) // Добавить услугу в заявку-черновик. Если черновик отсутствует, создается новый с указанными значениями.

	// Заявки (домен: `/project`)
	r.GET("/project", app.CheckAuth(), app.GetProjectList)                                 // Получить список заявок с фильтрацией по диапазону даты формирования и статусу (исключая удаленные и черновики, поля модератора и создателя отображаются через логины).
	r.GET("/project/:id", app.CheckAuth(), app.GetProjectByID)                             // Получить данные конкретной заявки по ее ID, включая список услуг и изображения.
	r.PUT("/project/:id", app.CheckAuth(), app.CheckUser(), app.UpdateProject)             // Изменить поля заявки по теме.
	r.PUT("/project/:id/submit", app.CheckAuth(), app.CheckUser(), app.SubmitProject)      // Сформировать заявку создателем с установкой даты формирования. Происходит проверка обязательных полей.
	r.PUT("/project/:id/complete", app.CheckAuth(), app.CheckAdmin(), app.CompleteProject) // Завершить или отклонить заявку модератором, указываются модератор и дата завершения. При завершении производится расчет дополнительных полей.
	r.DELETE("/project/:id", app.CheckAuth(), app.CheckUser(), app.DeleteProject)          // Удалить заявку (удаляется только при отсутствии даты формирования).

	// ММ (домен: `/file`)
	r.DELETE("/file/delete", app.CheckAuth(), app.CheckUser(), app.DeleteFileFromProject) // Удалить файл из заявки (без первичного ключа файла).
	r.PUT("/file/update", app.CheckAuth(), app.CheckUser(), app.UpdateFileInProject)      // Изменить количество, порядок или значение файла в заявке (без первичного ключа файла).

	// Пользователи (домен: `/user`)
	r.POST("/user/register", app.RegisterUser)                    // Регистрация нового пользователя.
	r.PUT("/user/update", app.CheckAuth(), app.UpdateUserProfile) // Изменить данные пользователя (личный кабинет).
	r.POST("/user/login", app.UserLogin)                          // Аутентификация пользователя.
	r.POST("/user/logout", app.CheckAuth(), app.UserLogout)       // Деавторизация пользователя.
}

func (app *App) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if app.userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (app *App) CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !app.isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (app *App) CheckUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if app.isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "User access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

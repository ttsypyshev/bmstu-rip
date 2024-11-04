package backend

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteFileRequest struct {
	ProjectID uint `json:"project_id" binding:"required"` // ID проекта
	LangID    uint `json:"lang_id" binding:"required"`    // ID услуги
}

func (app *App) DeleteFileFromProject(c *gin.Context) {
	var req DeleteFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}
	log.Printf("[info] DeleteFileFromProject called: ProjectID=%d, LangID=%d", req.ProjectID, req.LangID)

	// Проверяем, существует ли файл с указанным ProjectID и LangID
	file, err := app.findFile(req.ProjectID, req.LangID)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] file not found"), err)
		return
	}

	// Удаляем файл
	if err := app.deleteFile(file.ID); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to delete file"), err)
		return
	}

	log.Printf("[info] File successfully deleted from project: ProjectID=%d, LangID=%d", req.ProjectID, req.LangID)
	c.JSON(http.StatusOK, gin.H{
		"message": "File successfully deleted from project",
	})
}

type UpdateFileRequest struct {
	ProjectID uint   `json:"project_id" binding:"required"` // ID проекта
	LangID    uint   `json:"lang_id" binding:"required"`    // ID услуги
	Code      string `json:"code"`                          // Новое значение кода файла
	AutoCheck *int   `json:"auto_check"`                    // Новый статус автопроверки
	Order     *int   `json:"order"`                         // Новый порядок файла, если применимо
	Count     *int   `json:"count"`                         // Новое количество, если применимо
}

func (app *App) UpdateFileInProject(c *gin.Context) {
	var req UpdateFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}
	log.Printf("[info] UpdateFileInProject called: ProjectID=%d, LangID=%d", req.ProjectID, req.LangID)

	// Находим файл по ProjectID и LangID
	file, err := app.findFile(req.ProjectID, req.LangID)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] file not found"), err)
		return
	}

	// Обновляем указанные поля файла
	file = DbFile{
		ProjectID: req.ProjectID,
		LangID:    req.LangID,
		Code:      req.Code,
		AutoCheck: *req.AutoCheck,
	}

	// Обновляем запись в базе данных
	if err := app.updateFile(&file); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update file"), err)
		return
	}

	log.Printf("[info] File successfully updated in project: ProjectID=%d, LangID=%d", req.ProjectID, req.LangID)
	c.JSON(http.StatusOK, gin.H{
		"message": "File successfully updated in project",
	})
}

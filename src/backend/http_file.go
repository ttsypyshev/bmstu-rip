package backend

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteFileRequest struct {
	ProjectID uint `json:"project_id" binding:"required"`
	LangID    uint `json:"lang_id" binding:"required"`
}

func (app *App) DeleteFileFromProject(c *gin.Context) {
	var req DeleteFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	project, err := app.getProjectByID(req.ProjectID)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	if app.userID != project.UserID {
		handleError(c, http.StatusNotFound, errors.New("[err] project does not belong to the user"), err)
		return
	}

	file, err := app.findFile(req.ProjectID, req.LangID)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] file not found"), err)
		return
	}

	if err := app.deleteFile(req.ProjectID, file.ID); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to delete file"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File successfully deleted from project",
		"project": project,
	})
}

type UpdateFileRequest struct {
	ProjectID uint   `json:"project_id" binding:"required"`
	LangID    uint   `json:"lang_id" binding:"required"`
	Code      string `json:"code"`
	AutoCheck *int   `json:"auto_check"`
}

func (app *App) UpdateFileInProject(c *gin.Context) {
	var req UpdateFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	project, err := app.getProjectByID(req.ProjectID)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	if app.userID != project.UserID {
		handleError(c, http.StatusNotFound, errors.New("[err] project does not belong to the user"), err)
		return
	}

	file, err := app.findFile(req.ProjectID, req.LangID)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] file not found"), err)
		return
	}

	if req.ProjectID == 0 {
		handleError(c, http.StatusBadRequest, errors.New("[err] project ID is required"))
		return
	} else {
		file.ProjectID = req.ProjectID
	}

	if req.LangID == 0 {
		handleError(c, http.StatusBadRequest, errors.New("[err] language ID is required"))
		return
	} else {
		file.LangID = req.LangID
	}

	file.Code = req.Code

	if req.AutoCheck != nil {
		file.AutoCheck = *req.AutoCheck
	}

	if err := app.updateFile(&file); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update file"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File successfully updated in project",
		"file":    file,
	})
}

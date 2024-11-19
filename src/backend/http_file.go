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

// DeleteFileFromProject godoc
// @Summary Delete a file from a project
// @Description Deletes a file from a project based on the provided project ID and language ID. Only the project owner can delete files from their project.
// @Tags Files
// @Accept json
// @Produce json
// @Param request body DeleteFileRequest true "File deletion request data"
// @Success 200 {object} gin.H "File successfully deleted from project"
// @Failure 400 {object} ErrorResponse "Invalid request format or missing fields"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "User is not the project owner"
// @Failure 404 {object} ErrorResponse "Project or file not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /file/delete [delete]
func (app *App) DeleteFileFromProject(c *gin.Context) {
	requestUserID, err := ExtractUserID(c)
	if err != nil {
		handleError(c, http.StatusUnauthorized, errors.New("[err] Unauthorized"), err)
		return
	}

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

	if requestUserID != project.UserID {
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
		"status":  true,
	})
}

type UpdateFileRequest struct {
	ProjectID uint   `json:"project_id" binding:"required"`
	LangID    uint   `json:"lang_id" binding:"required"`
	Code      string `json:"code"`
	FileName  string `json:"filename"`
	Comment   string `json:"comment"`
}

// UpdateFileInProject godoc
// @Summary Update a file in a project
// @Description Updates the details of a file within a project. The user must be the owner of the project to update the file.
// @Tags Files
// @Accept json
// @Produce json
// @Param request body UpdateFileRequest true "File update request data"
// @Success 200 {object} gin.H "File successfully updated in project"
// @Failure 400 {object} ErrorResponse "Invalid request format or missing fields"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "User is not the project owner"
// @Failure 404 {object} ErrorResponse "Project or file not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /file/update [put]
func (app *App) UpdateFileInProject(c *gin.Context) {
	requestUserID, err := ExtractUserID(c)
	if err != nil {
		handleError(c, http.StatusUnauthorized, errors.New("[err] Unauthorized"), err)
		return
	}

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

	if requestUserID != project.UserID {
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

	if req.Code != "" {
		file.Code = req.Code
	}

	if req.FileName != "" {
		file.FileName = req.FileName
	}

	if req.Comment != "" {
		file.Comment = req.Comment
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

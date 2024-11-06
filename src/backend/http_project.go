package backend

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *App) GetProjectList(c *gin.Context) {
	startDateStr := c.Query("start_date") // формат: "YYYY-MM-DD"
	endDateStr := c.Query("end_date")     // формат: "YYYY-MM-DD"
	statusStr := c.Query("status")

	projects, err := app.filterProjects(startDateStr, endDateStr, statusStr)
	if err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to retrieve projects"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}

func (app *App) GetProjectByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid project ID format"), err)
		return
	}

	project, err := app.getProjectByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	files, err := app.getFilesForProject(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] files not found for project"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project": project,
		"files":   files,
	})
}

type UpdateProjectRequest struct {
	Status  *int   `json:"status"`
	Comment string `json:"comment"`
}

func (app *App) UpdateProject(c *gin.Context) {
	if app.isAdmin {
		handleError(c, http.StatusBadRequest, errors.New("[err] this is not the task of this user"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid project ID"), err)
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	project, err := app.getProjectByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	if app.userID != project.UserID {
		handleError(c, http.StatusNotFound, errors.New("[err] project does not belong to the user"), err)
		return
	}

	if req.Status != nil {
		project.Status = *req.Status
	} else {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid status"))
		return
	}

	if req.Comment != "" {
		project.Comment = req.Comment
	}

	if err := app.updateProject(&project); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update project"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project updated successfully",
		"project": project,
	})
}

type AddProjectRequest struct {
	IDProject uint            `form:"id_project" json:"id_project"`
	FileCodes map[uint]string `form:"file_codes" json:"file_codes"`
}

func (app *App) SubmitProject(c *gin.Context) {
	if app.isAdmin {
		handleError(c, http.StatusBadRequest, errors.New("[err] this is not the task of this user"))
		return
	}

	var req AddProjectRequest
	if err := c.ShouldBind(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid data format"), err)
		return
	}

	project, err := app.getProjectByID(req.IDProject)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	if app.userID != project.UserID {
		handleError(c, http.StatusNotFound, errors.New("[err] project does not belong to the user"), err)
		return
	}

	if err := app.updateFilesCode(req.IDProject, req.FileCodes); err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to update file"), err)
		return
	}

	// Обновляем статус проекта на "сформирован" (или статус 2)
	project.Status = 2
	if err := app.updateProject(&project); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update project"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project submitted successfully",
		"project": project,
	})
}

type DeleteProjectRequest struct {
	IDProject uint            `form:"id_project" json:"id_project"`
	FileCodes map[uint]string `form:"file_codes" json:"file_codes"`
}

func (app *App) DeleteProject(c *gin.Context) {
	if app.isAdmin {
		handleError(c, http.StatusBadRequest, errors.New("[err] this is not the task of this user"))
		return
	}

	var req DeleteProjectRequest
	if err := c.ShouldBind(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid data format"), err)
		return
	}

	project, err := app.getProjectByID(req.IDProject)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	if app.userID != project.UserID {
		handleError(c, http.StatusNotFound, errors.New("[err] project does not belong to the user"), err)
		return
	}

	if project.FormationTime != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] project cannot be deleted, formation date found"))
		return
	}

	if err := app.updateFilesCode(req.IDProject, req.FileCodes); err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to update file"), err)
		return
	}

	// Обновляем статус проекта на "удалён" (или статус 1)
	project.Status = 1
	if err := app.updateProject(&project); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update project"), err)
		return
	}

	log.Printf("[info] Project %d deleted successfully", req.IDProject)
	c.JSON(http.StatusOK, gin.H{
		"message":   "Project deleted successfully",
		"projectID": req.IDProject,
	})
}

type CompleteProjectRequest struct {
	IDProject   uint   `form:"id_project" json:"id_project"`
	ModeratorID uint   `json:"moderator_id"`
	Status      int    `json:"status"`
	Comment     string `json:"comment"`
}

func (app *App) CompleteProject(c *gin.Context) {
	if !app.isAdmin {
		handleError(c, http.StatusBadRequest, errors.New("[err] user does not have sufficient rights"))
		return
	}

	var req CompleteProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	project, err := app.getProjectByID(req.IDProject)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	if project.FormationTime == nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] project cannot be complete, formation date exists"))
		return
	}

	if req.ModeratorID != 0 {
		project.ModeratorID = &req.ModeratorID
	} else {
		handleError(c, http.StatusBadRequest, errors.New("[err] moderator ID is required"))
		return
	}

	if req.Status == 3 || req.Status == 4 {
		project.Status = req.Status
	} else {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid status"))
		return
	}

	if req.Comment != "" {
		project.Comment = req.Comment
	}

	if err := app.updateProject(&project); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to complete project"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project completed successfully",
		"project": project,
	})
}

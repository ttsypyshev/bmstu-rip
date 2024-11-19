package backend

// import (
// 	"errors"
// 	"net/http"
// 	database "rip/pkg/database"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// func (app *App) GetProjectList(c *gin.Context) {
// 	startDateStr := c.Query("start_date") // формат: "YYYY-MM-DD"
// 	endDateStr := c.Query("end_date")     // формат: "YYYY-MM-DD"
// 	status := c.Query("status")

// 	projects, err := app.filterProjects(startDateStr, endDateStr, status)
// 	if err != nil {
// 		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to retrieve projects"), err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"projects": projects,
// 	})
// }

// func (app *App) GetProjectByID(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid project ID format"), err)
// 		return
// 	}

// 	project, err := app.getProjectByID(uint(id))
// 	if err != nil {
// 		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
// 		return
// 	}

// 	files, err := app.getFilesForProject(uint(id))
// 	if err != nil {
// 		handleError(c, http.StatusNotFound, errors.New("[err] files not found for project"), err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"project": project,
// 		"files":   files,
// 	})
// }

// type UpdateProjectRequest struct {
// 	Status  database.Status `json:"status"`
// 	Comment string          `json:"comment"`
// }

// func (app *App) UpdateProject(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid project ID"), err)
// 		return
// 	}

// 	var req UpdateProjectRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
// 		return
// 	}

// 	project, err := app.getProjectByID(uint(id))
// 	if err != nil {
// 		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
// 		return
// 	}

// 	if *app.userID != project.UserID {
// 		handleError(c, http.StatusNotFound, errors.New("[err] project does not belong to the user"), err)
// 		return
// 	}

// 	if req.Status != "" {
// 		project.Status = req.Status
// 	} else {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid status"))
// 		return
// 	}

// 	if req.Comment != "" {
// 		project.ModeratorComment = req.Comment
// 	}

// 	if err := app.updateProject(&project); err != nil {
// 		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update project"), err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Project updated successfully",
// 		"project": project,
// 	})
// }

// type AddProjectRequest struct {
// 	FileCodes map[uint]string `json:"file_codes"`
// }

// func (app *App) SubmitProject(c *gin.Context) {
// 	idStr := c.Param("id")
// 	projectID, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid id project"), err)
// 	}

// 	var req AddProjectRequest
// 	if err = c.ShouldBind(&req); err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid data format"), err)
// 		return
// 	}

// 	project, err := app.getProjectByID(uint(projectID))
// 	if err != nil {
// 		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
// 		return
// 	}

// 	if *app.userID != project.UserID {
// 		handleError(c, http.StatusNotFound, errors.New("[err] project does not belong to the user"), err)
// 		return
// 	}

// 	if err := app.updateFilesCode(project.ID, req.FileCodes); err != nil {
// 		handleError(c, http.StatusNotFound, errors.New("[err] failed to update file"), err)
// 		return
// 	}

// 	project.Status = database.Created
// 	if err := app.updateProject(&project); err != nil {
// 		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update project"), err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Project submitted successfully",
// 		"project": project,
// 	})
// }

// type CompleteProjectRequest struct {
// 	Status  database.Status `json:"status"`
// 	Comment string          `json:"comment"`
// }

// func (app *App) CompleteProject(c *gin.Context) {
// 	idStr := c.Param("id")
// 	projectID, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid id project"), err)
// 	}

// 	var req CompleteProjectRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
// 		return
// 	}

// 	project, err := app.getProjectByID(uint(projectID))
// 	if err != nil {
// 		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
// 		return
// 	}

// 	if project.FormationTime == nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] project cannot be complete, formation date exists"))
// 		return
// 	}

// 	project.ModeratorID = app.userID

// 	if req.Status == database.Completed || req.Status == database.Rejected {
// 		project.Status = req.Status
// 	} else {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid status"))
// 		return
// 	}

// 	if req.Comment != "" {
// 		project.ModeratorComment = req.Comment
// 	}

// 	if err := app.updateProject(&project); err != nil {
// 		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to complete project"), err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Project completed successfully",
// 		"project": project,
// 	})
// }

// type DeleteProjectRequest struct {
// 	FileCodes map[uint]string `json:"file_codes"`
// }

// func (app *App) DeleteProject(c *gin.Context) {
// 	idStr := c.Param("id")
// 	projectID, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid id project"), err)
// 	}

// 	var req DeleteProjectRequest
// 	if err := c.ShouldBind(&req); err != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] invalid data format"), err)
// 		return
// 	}

// 	project, err := app.getProjectByID(uint(projectID))
// 	if err != nil {
// 		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
// 		return
// 	}

// 	if *app.userID != project.UserID {
// 		handleError(c, http.StatusNotFound, errors.New("[err] project does not belong to the user"), err)
// 		return
// 	}

// 	if project.FormationTime != nil {
// 		handleError(c, http.StatusBadRequest, errors.New("[err] project cannot be deleted, formation date found"))
// 		return
// 	}

// 	if err := app.updateFilesCode(project.ID, req.FileCodes); err != nil {
// 		handleError(c, http.StatusNotFound, errors.New("[err] failed to update file"), err)
// 		return
// 	}

// 	// Обновляем статус проекта на "удалён" (или статус 1)
// 	project.Status = database.Deleted
// 	if err := app.updateProject(&project); err != nil {
// 		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update project"), err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Project deleted successfully",
// 		"status":  true,
// 	})
// }

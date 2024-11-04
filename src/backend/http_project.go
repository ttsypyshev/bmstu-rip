package backend

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (app *App) GetProjectList(c *gin.Context) {
	// Извлекаем параметры запроса для фильтрации
	startDateStr := c.Query("start_date") // формат: "YYYY-MM-DD"
	endDateStr := c.Query("end_date")     // формат: "YYYY-MM-DD"
	statusStr := c.Query("status")        // статус заявки

	// Преобразуем статус в целое число
	status, err := strconv.Atoi(statusStr)
	if err != nil && statusStr != "" {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid status value"), err)
		return
	}

	// Получаем список проектов с фильтрацией
	projects, err := app.filterProjects(startDateStr, endDateStr, status)
	if err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to retrieve projects"), err)
		return
	}

	// Формируем ответ в формате JSON
	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}

func (app *App) GetProjectByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid project ID"), err)
		return
	}

	// Получаем проект по ID
	project, err := app.getProjectByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to retrieve project information"), err)
		return
	}

	// Получаем файлы, связанные с проектом
	files, err := app.getFilesForProject(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to retrieve project files"), err)
		return
	}

	// Получаем языки, которые активны
	langs, err := app.getLangs(func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", true)
	})
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to retrieve language information"), err)
		return
	}

	// Формируем ответ в формате JSON
	c.JSON(http.StatusOK, gin.H{
		"project": project,
		"files":   files,
		"langs":   langs,
	})
}

type UpdateProjectRequest struct {
	Status         *int       `json:"status"`
	CompletionTime *time.Time `json:"completion_time"`
	DeletionTime   *time.Time `json:"deletion_time"`
}

func (app *App) UpdateProject(c *gin.Context) {
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

	// Получаем проект из базы данных
	project, err := app.getProjectByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	project = DbProject{
		Status:         *req.Status,
		CompletionTime: req.CompletionTime,
		DeletionTime:   req.DeletionTime,
	}

	// Сохраняем изменения в базе данных
	if err := app.updateProject(&project); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update project"), err)
		return
	}

	// Возвращаем обновлённый проект
	c.JSON(http.StatusOK, gin.H{
		"message": "Project updated successfully",
		"project": project,
	})
}

func (app *App) SubmitProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid project ID"), err)
		return
	}

	// Получаем проект из базы данных
	project, err := app.getProjectByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	// Проверяем обязательные поля
	if project.UserID == 0 {
		handleError(c, http.StatusBadRequest, errors.New("[err] UserID is required"), nil)
		return
	}
	// Здесь вы можете добавить дополнительные проверки для других обязательных полей
	// например, проверка наличия связанных файлов или других необходимых данных.

	now := time.Now()
	project.CreationTime = now
	project.Status = 2 // 2 - сформирован

	// Сохраняем изменения в базе данных
	if err := app.updateProject(&project); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to submit project"), err)
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"message": "Project submitted successfully",
		"project": project,
	})
}

type CompleteProjectRequest struct {
	ModeratorID    uint       `json:"moderator_id"`
	Status         int        `json:"status"`
	CompletionTime *time.Time `json:"completion_time"`
}

func (app *App) CompleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid project ID"), err)
		return
	}

	var req CompleteProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	// Получаем проект из базы данных
	project, err := app.getProjectByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	// Проверяем, может ли модератор завершить проект
	if project.Status != 2 { // 2 - статус "сформирован"
		handleError(c, http.StatusBadRequest, errors.New("[err] project cannot be completed"), nil)
		return
	}

	// Устанавливаем информацию о модераторе и дату завершения
	project.ModeratorID = &req.ModeratorID
	if req.CompletionTime != nil {
		project.CompletionTime = req.CompletionTime
	} else {
		now := time.Now()
		project.CompletionTime = &now // Устанавливаем текущее время, если не указано
	}

	// Обновляем статус проекта
	project.Status = req.Status

	// Выполняем расчеты для дополнительных полей (если необходимо)
	// Например, можем подсчитать общее количество файлов или другие метрики
	// app.calculateProjectMetrics(project)

	// Сохраняем изменения в базе данных
	if err := app.updateProject(&project); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to complete project"), err)
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"message": "Project completed successfully",
		"project": project,
	})
}

type RequestDelete struct {
	IDProject uint            `form:"id_project" json:"id_project"`
	FileCodes map[uint]string `form:"file_codes" json:"file_codes"`
}

func (app *App) DeleteProject(c *gin.Context) {
	var req RequestDelete
	if err := c.ShouldBind(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid data format"), err)
		return
	}
	log.Printf("[info] DeleteProject called: IDProject=%d, FileCodes=%v", req.IDProject, req.FileCodes)

	// Проверяем, существует ли проект
	project, err := app.getProjectByID(req.IDProject)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	// Проверяем наличие даты формирования
	if project.CreationTime.IsZero() {
		handleError(c, http.StatusBadRequest, errors.New("[err] project cannot be deleted, creation date exists"), nil)
		return
	}

	// Обновляем статус проекта на "удалён" (или статус 1)
	if err := app.updateProjectStatus(req.IDProject, 1); err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to update project status"), err)
		return
	}

	// Обновляем информацию о файлах
	if err := app.updateFilesCode(req.FileCodes); err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to update file"), err)
		return
	}

	log.Printf("[info] Project %d deleted successfully", req.IDProject)
	c.JSON(http.StatusOK, gin.H{
		"message":   "Project deleted successfully",
		"projectID": req.IDProject,
	})
}

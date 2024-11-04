package backend

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *App) GetServiceList(c *gin.Context) {
	query := c.Query("langname")
	filteredLangs, err := app.GetFilteredLangs(query)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to retrieve language information"), err)
		return
	}

	projectID, err := findLastDraft(app, app.userID)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project was not created 1"), err)
		return
	}

	count, err := app.getProjectCount(projectID)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project was not created 2"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"langs":   filteredLangs,
		"draftID": projectID,
		"count":   count,
	})
}

func (app *App) GetServiceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid language ID"), err)
		return
	}

	lang, err := app.getLangByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] language information not available"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"info": lang,
		// "list":  ParseList(lang.List),
	})
}

type ServiceInput struct {
	Name             string `json:"name"`
	ShortDescription string
	Description      string `json:"description"`
	Author           string
	Year             string
	Version          string
	List             string
}

func (app *App) CreateService(c *gin.Context) {
	var input ServiceInput

	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid input data"), err)
		return
	}

	service := DbLang{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		Author:           input.Author,
		Year:             input.Year,
		Version:          input.Version,
		List:             input.List,
	}

	if err := app.saveService(&service); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to save service"), err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Service created successfully",
		"service": service,
	})
}

type ServiceUpdateInput struct {
	Name             string `json:"name"`
	ShortDescription string
	Description      string `json:"description"`
	Author           string
	Year             string
	Version          string
	List             string
}

func (app *App) UpdateService(c *gin.Context) {
	// Получаем ID услуги из параметров запроса
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid service ID"), err)
		return
	}

	// Парсим данные для обновления из JSON-запроса
	var input ServiceUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid input data"), err)
		return
	}

	// Получаем текущую запись услуги по ID
	service, err := app.getLangByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] service not found"), err)
		return
	}

	service = DbLang{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		Author:           input.Author,
		Year:             input.Year,
		Version:          input.Version,
		List:             input.List,
	}

	// Сохраняем обновленную услугу в базе данных
	if err := app.updateService(&service); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update service"), err)
		return
	}

	// Возвращаем JSON-ответ с обновленными данными услуги
	c.JSON(http.StatusOK, gin.H{
		"message": "Service updated successfully",
		"service": service,
	})
}

func (app *App) UpdateServiceImage(c *gin.Context) {
	// Извлекаем ID услуги из параметров запроса
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid service ID"), err)
		return
	}

	// Проверяем наличие услуги с указанным ID
	service, err := app.getLangByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] service not found"), err)
		return
	}

	// Получаем файл изображения из запроса
	file, err := c.FormFile("image")
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] image file is required"), err)
		return
	}

	// Сохраняем изображение
	//! minio add image

	// Получаем ссылку
	// !minio get imagePath

	// Обновляем путь к изображению в данных услуги (если оно хранится в базе данных)
	service.ImagePath = imagePath
	if err := app.updateService(&service); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update service image path"), err)
		return
	}

	// Возвращаем JSON-ответ с подтверждением успешного обновления изображения
	c.JSON(http.StatusOK, gin.H{
		"message": "Service image updated successfully",
		"service": service,
	})
}

func (app *App) DeleteService(c *gin.Context) {
	// Извлекаем ID услуги из параметров запроса
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid service ID"), err)
		return
	}

	// Проверяем наличие услуги с указанным ID
	service, err := app.getLangByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] service not found"), err)
		return
	}

	// Удаляем изображение услуги
	//! minio delete image

	// Удаляем запись услуги из базы данных
	if err := app.deleteServiceByID(uint(id)); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to delete service"), err)
		return
	}

	// Возвращаем JSON-ответ с подтверждением успешного удаления
	c.JSON(http.StatusOK, gin.H{
		"message":   "Service deleted successfully",
		"serviceID": id,
	})
}

type RequestAdd struct {
	IDLang uint `form:"id_lang" json:"id_lang" binding:"required"`
}

func (app *App) AddServiceToDraft(c *gin.Context) {
	var req RequestAdd

	// Парсим данные запроса в структуру RequestAdd
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid data format"), err)
		return
	}
	log.Printf("[info] AddService called: IDUser=%d, IDLang=%d", app.userID, req.IDLang)

	// Создаем черновик, если он еще не существует, и получаем его ID
	projectID, err := createDraft(app, app.userID)
	if err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] error creating project"), err)
		return
	}

	// Добавляем услугу в черновик
	if err := app.addFile(projectID, req.IDLang, app.userID); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to add service to draft"), err)
		return
	}

	// Успешное добавление услуги в черновик
	c.JSON(http.StatusOK, gin.H{
		"message":   "Service added to draft successfully",
		"projectID": projectID,
	})
}

package backend

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *App) GetServiceList(c *gin.Context) {
	query := c.Query("langname")
	filteredLangs, err := app.GetFilteredLangs(query)
	if err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to retrieve filtered languages"), err)
		return
	}

	projectID, err := findLastDraft(app, *app.userID)
	if err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] unable to find the last draft project for the user"), err)
		return
	}

	count, err := app.getProjectCount(projectID)
	if err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] unable to retrieve project count for the draft"), err)
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
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid language ID format"), err)
		return
	}

	lang, err := app.getLangByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] language not found for the given ID"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"info": lang,
	})
}

type CreateServiceRequest struct {
	Name             string            `json:"name"`
	ShortDescription string            `json:"short_description"`
	Description      string            `json:"description"`
	Author           string            `json:"author"`
	Year             string            `json:"year"`
	Version          string            `json:"version"`
	List             map[string]string `json:"list"`
}

func (app *App) CreateService(c *gin.Context) {
	var input CreateServiceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid input data: unable to parse JSON"), err)
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

	langID, err := app.createLang(&service)
	if err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to save service in the database"), err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Service successfully created.",
		"serviceID": langID,
	})
}

type UpdateServiceRequest struct {
	Name             string            `json:"name"`
	ShortDescription string            `json:"short_description"`
	Description      string            `json:"description"`
	Author           string            `json:"author"`
	Year             string            `json:"year"`
	Version          string            `json:"version"`
	List             map[string]string `json:"list"`
}

func (app *App) UpdateService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid service ID format"), err)
		return
	}

	var input UpdateServiceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid input data"), err)
		return
	}

	service, err := app.getLangByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] service not found"), err)
		return
	}

	if input.Name != "" {
		service.Name = input.Name
	}
	if input.ShortDescription != "" {
		service.ShortDescription = input.ShortDescription
	}
	if input.Description != "" {
		service.Description = input.Description
	}
	if input.Author != "" {
		service.Author = input.Author
	}
	if input.Year != "" {
		service.Year = input.Year
	}
	if input.Version != "" {
		service.Version = input.Version
	}
	if len(input.List) != 0 {
		service.List = input.List
	}

	if err := app.updateLang(&service); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update service"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service updated successfully",
		"service": service,
	})
}

func (app *App) UpdateServiceImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid service ID"), err)
		return
	}

	service, err := app.getLangByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] service not found"), err)
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] image file is required"), err)
		return
	}

	imageURL, err := app.uploadImageToMinIO(file)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, nil)
		return
	}

	service.ImgLink = imageURL
	service.Status = false

	if err := app.updateLang(&service); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update service image path"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service image updated successfully",
		"service": service,
	})
}

func (app *App) DeleteService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid service ID"), err)
		return
	}

	service, err := app.getLangByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] service not found"), err)
		return
	}

	if err := app.deleteImageFromMinIO(service.ImgLink); err != nil {
		handleError(c, http.StatusInternalServerError, err, nil)
		return
	}

	if err := app.deleteLang(uint(id)); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to delete service"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service deleted successfully",
		"status":  true,
	})
}

type AddServiceRequest struct {
	IDLang uint `json:"id_lang"`
}

func (app *App) AddServiceToDraft(c *gin.Context) {
	var req AddServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid JSON format or missing fields"), err)
		return
	}

	projectID, err := app.createProject(*app.userID)
	if err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] unable to create project"), err)
		return
	}

	if err := app.createFile(projectID, req.IDLang); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to create file for service draft"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Service successfully added to draft",
		"projectID": projectID,
	})
}

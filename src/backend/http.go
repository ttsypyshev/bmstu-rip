package backend

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (app *App) SetupRoutes(r *gin.Engine) {
	r.GET("/home", app.handleHome)
	r.GET("/info/:id", app.handleInfo)
	r.GET("/project/:id", app.handleApp)
	r.POST("/add-service", app.handleAddService)
	r.POST("/del-project", app.handleDeleteProject)
}

func (app *App) handleHome(c *gin.Context) {
	langID, err := ParseQueryParam(c, "langID")
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid langID"), err)
		return
	}

	status, err := ParseQueryParam(c, "status")
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid status"), err)
		return
	}

	query := c.Query("langname")
	filteredLangs, err := app.GetFilteredLangs(query)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to retrieve language information"), err)
		return
	}

	projectID, err := FindLastDraft(app, app.userID)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project was not created 1"), err)
		return
	}

	count, err := app.GetProjectCount(projectID)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project was not created 2"), err)
		return
	}

	c.HTML(http.StatusOK, "services.tmpl", gin.H{
		"Title":     "Langs",
		"Langs":     filteredLangs,
		"Count":     count,
		"UserID":    app.userID,
		"ProjectID": projectID,
		"LangID":    langID,
		"Status":    status,
	})
}

func (app *App) handleInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] invalid language ID"), err)
		return
	}

	lang, err := app.GetLangByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] language information not available"), err)
		return
	}

	c.HTML(http.StatusOK, "information.tmpl", gin.H{
		"Title": lang.Name,
		"Info":  lang,
		"List":  ParseList(lang.List),
	})
}

func (app *App) handleApp(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] invalid project ID"), err)
		return
	}

	project, err := app.GetProjectByID(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to retrieve project information"), err)
		return
	}

	files, err := app.GetFilesForProject(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to retrieve project files"), err)
		return
	}

	langs, err := app.GetLangs(func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", true)
	})
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to retrieve language information"), err)
		return
	}

	c.HTML(http.StatusOK, "applications.tmpl", gin.H{
		"Title":   "Project",
		"Project": project,
		"Files":   files,
		"Langs":   langs,
	})
}

type RequestAdd struct {
	IDUser uint `form:"id_user" json:"id_user"`
	IDLang uint `form:"id_lang" json:"id_lang"`
}

func (app *App) handleAddService(c *gin.Context) {
	var req RequestAdd
	if err := c.ShouldBind(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid data format"), err)
		return
	}
	log.Printf("[info] AddService called : IDUser=%d, IDLang=%d", req.IDUser, req.IDLang)

	projectID, err := CreateDraft(app, req.IDUser)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] error creating project"), err)
		return
	}

	if err := app.AddFile(projectID, req.IDLang, req.IDUser); err != nil {
		log.Printf("[info] Redirecting to home page at URL: /home?langID=%d&status=%d", req.IDLang, 2)
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/home?langID=%d&status=%d", req.IDLang, 2))
		return
	}

	log.Printf("[info] Redirecting to home page at URL: /home?langID=%d&status=%d", req.IDLang, 1)
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/home?langID=%d&status=%d", req.IDLang, 1))
}

type RequestDelete struct {
	IDProject uint            `form:"id_project" json:"id_project"`
	FileCodes map[uint]string `form:"file_codes" json:"file_codes"`
}

func (app *App) handleDeleteProject(c *gin.Context) {
	var req RequestDelete
	if err := c.ShouldBind(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid data format"), err)
		return
	}
	log.Printf("[info] DeleteProject called: IDProject=%d, FileCodes=%v", req.IDProject, req.FileCodes)

	if err := app.UpdateProjectStatus(req.IDProject, 1); err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to update project status"), err)
		return
	}

	if err := app.UpdateFilesCode(req.FileCodes); err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] failed to update file"), err)
		return
	}

	log.Printf("[info] Redirecting to home page at URL: /home")
	c.Redirect(http.StatusSeeOther, "/home")
}

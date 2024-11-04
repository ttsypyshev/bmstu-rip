package backend

import (
	"fmt"
	"log"
	"rip/database"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	DbLang    = database.Lang
	DbProject = database.Project
	DbFile    = database.File
	DbUser    = database.User
)

type Db struct {
	db *gorm.DB
}

func Migrate() error {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(FromEnv()), &gorm.Config{})
	if err != nil {
		return err
	}

	// Migrate the schema
	err = db.AutoMigrate(&DbLang{}, &DbProject{}, &DbFile{}, &DbUser{})
	if err != nil {
		return err
	}

	return nil
}

func NewDB(dsn string) (*App, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &App{
		db: &Db{db},
	}, nil
}

// Получение всех сущностей
func getAll[T any](app *App, filter func(*gorm.DB) *gorm.DB) ([]T, error) {
	var items []T

	query := app.db.db.Model(&items)
	if filter != nil {
		query = filter(query)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return items, nil
}

func getEntities[T any](app *App, filterFunc ...func(db *gorm.DB) *gorm.DB) ([]T, error) {
	var filter func(db *gorm.DB) *gorm.DB
	if len(filterFunc) > 0 {
		filter = filterFunc[0]
	}
	return getAll[T](app, filter)
}

func (app *App) GetLangs(filterFunc ...func(db *gorm.DB) *gorm.DB) ([]DbLang, error) {
	return getEntities[DbLang](app, filterFunc...)
}

func (app *App) GetProjects(filterFunc ...func(db *gorm.DB) *gorm.DB) ([]DbProject, error) {
	return getEntities[DbProject](app, filterFunc...)
}

func (app *App) GetFiles(filterFunc ...func(db *gorm.DB) *gorm.DB) ([]DbFile, error) {
	return getEntities[DbFile](app, filterFunc...)
}

// Получение сущностей по ID
func getByID[T any](app *App, id uint) (T, error) {
	var item T

	err := app.db.db.First(&item, "id = ?", id).Error
	if err != nil {
		return item, err
	}

	return item, nil
}

func (app *App) GetLangByID(langID uint) (DbLang, error) {
	return getByID[DbLang](app, langID)
}

func (app *App) GetProjectByID(projectID uint) (DbProject, error) {
	return getByID[DbProject](app, projectID)
}

func (app *App) GetFileByID(fileID uint) (DbFile, error) {
	return getByID[DbFile](app, fileID)
}

// Получение файлов для проекта
func (app *App) GetFilesForProject(projectID uint) ([]DbFile, error) {
	var matchedFiles []DbFile
	
	if err := app.db.db.Where("project_id = ?", projectID).Find(&matchedFiles).Error; err != nil {
		return nil, err
	}

	return matchedFiles, nil
}

// Подсчет количества файлов в черновике пользователя
func (app *App) GetProjectCount(projectID uint) (int64, error) {
	var count int64
	if err := app.db.db.Model(&DbProject{}).Select("count").Where("id = ?", projectID).Scan(&count).Error; err != nil {
		return -1, err
	}

	return count, nil
}

// Фильтрация языков по запросу
func (app *App) FilterLangsByQuery(query string) ([]DbLang, error) {
	var filteredLangs []DbLang
	lowerQuery := "%" + strings.ToLower(query) + "%"

	result := app.db.db.Where("LOWER(name) LIKE ? AND status = ?", lowerQuery, true).Find(&filteredLangs)
	if result.Error != nil {
		return nil, result.Error
	}

	return filteredLangs, nil
}

// Создание нового черновика или возврат существующего
func CreateDraft(app *App, userID uint) (uint, error) {
	projectID, err := FindLastDraft(app, userID)
	if err != nil {
		return 0, err
	} else if projectID == 0 {
		newProject := DbProject{
			UserID:       userID,
			CreationTime: time.Now(),
			Status:       0,
		}

		if err := app.db.db.Create(&newProject).Error; err != nil {
			return 0, err
		}

		return newProject.ID, nil
	}
	return projectID, nil
}

// Поиск последнего черновика для пользователя
func FindLastDraft(app *App, userID uint) (uint, error) {
	var lastProject DbProject

	if err := app.db.db.Where("status = ? AND user_id = ?", 0, userID).First(&lastProject).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return lastProject.ID, nil
}

// Добавление файла к проекту для пользователя
func (app *App) AddFile(projectID, langID, userID uint) error {
	newFile := DbFile{
		LangID:    langID,
		ProjectID: projectID,
	}

	var count int64
	if err := app.db.db.Model(&DbFile{}).Where("lang_id = ? AND project_id = ?", newFile.LangID, newFile.ProjectID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("file with IDLang %d and IDProject %d already exists", newFile.LangID, newFile.ProjectID)
	}

	if err := app.db.db.Create(&newFile).Error; err != nil {
		return err
	}

	if err := app.db.db.Model(&DbProject{}).Where("id = ?", projectID).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
		return err
	}

	return nil
}

// Обновление статуса проекта
func (app *App) UpdateProjectStatus(projectID uint, newStatus uint) error {
	query := "UPDATE projects SET status = ? WHERE id = ?"

	result := app.db.db.Exec(query, newStatus, projectID)
	if result.Error != nil {
		return fmt.Errorf("failed to update project status: %w", result.Error)
	}

	return nil
}

// Обновление кода файлов по предоставленным мапам идентификаторов и кода
func (app *App) UpdateFilesCode(idToCodeMap map[uint]string) error {
	for id, newCode := range idToCodeMap {
		var file DbFile
		if err := app.db.db.Where("id = ?", id).First(&file).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Printf("file with id %d not found, skipping...", id)
				continue
			}
			return err
		}

		file.Code = newCode

		if err := app.db.db.Save(&file).Error; err != nil {
			return fmt.Errorf("failed to update file with id %d: %v", id, err)
		}
	}
	return nil
}

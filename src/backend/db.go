package backend

import (
	"errors"
	"fmt"
	"log"
	"rip/database"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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

func (app *App) getLangs(filterFunc ...func(db *gorm.DB) *gorm.DB) ([]DbLang, error) {
	return getEntities[DbLang](app, filterFunc...)
}

func (app *App) getProjects(filterFunc ...func(db *gorm.DB) *gorm.DB) ([]DbProject, error) {
	return getEntities[DbProject](app, filterFunc...)
}

func (app *App) getFiles(filterFunc ...func(db *gorm.DB) *gorm.DB) ([]DbFile, error) {
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

func (app *App) getLangByID(langID uint) (DbLang, error) {
	return getByID[DbLang](app, langID)
}

func (app *App) getProjectByID(projectID uint) (DbProject, error) {
	return getByID[DbProject](app, projectID)
}

func (app *App) getFileByID(fileID uint) (DbFile, error) {
	return getByID[DbFile](app, fileID)
}

// Получение пользователя по ID
func (app *App) getUserByID(userID uint) (*database.User, error) {
	var user database.User
	err := app.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Получение файлов для проекта
func (app *App) getFilesForProject(projectID uint) ([]DbFile, error) {
	var matchedFiles []DbFile

	if err := app.db.db.Where("project_id = ?", projectID).Find(&matchedFiles).Error; err != nil {
		return nil, err
	}

	return matchedFiles, nil
}

// Подсчет количества файлов в черновике пользователя
func (app *App) getProjectCount(projectID uint) (int64, error) {
	var count int64
	if err := app.db.db.Model(&DbProject{}).Select("count").Where("id = ?", projectID).Scan(&count).Error; err != nil {
		return -1, err
	}

	return count, nil
}

// Фильтрация языков по запросу
func (app *App) filterLangsByQuery(query string) ([]DbLang, error) {
	var filteredLangs []DbLang
	lowerQuery := "%" + strings.ToLower(query) + "%"

	result := app.db.db.Where("LOWER(name) LIKE ? AND status = ?", lowerQuery, true).Find(&filteredLangs)
	if result.Error != nil {
		return nil, result.Error
	}

	return filteredLangs, nil
}

// Создание нового черновика или возврат существующего
func createDraft(app *App, userID uint) (uint, error) {
	projectID, err := findLastDraft(app, userID)
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
func findLastDraft(app *App, userID uint) (uint, error) {
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
func (app *App) addFile(projectID, langID, userID uint) error {
	newFile := DbFile{
		LangID:    langID,
		ProjectID: projectID,
	}

	var count int64
	// todo заменить на findProjects
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
func (app *App) updateProjectStatus(projectID uint, newStatus uint) error {
	query := "UPDATE projects SET status = ? WHERE id = ?"

	result := app.db.db.Exec(query, newStatus, projectID)
	if result.Error != nil {
		return fmt.Errorf("failed to update project status: %w", result.Error)
	}

	return nil
}

// Обновление кода файлов по предоставленным мапам идентификаторов и кода
func (app *App) updateFilesCode(idToCodeMap map[uint]string) error {
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

// Функция для фильтрации проектов по дате и статусу
func (app *App) filterProjects(startDateStr, endDateStr string, status int) ([]Project, error) {
	var projects []Project
	query := app.db.Model(&Project{}).Where("status NOT IN (?)", []int{1, 0}) // исключаем удаленные (1) и черновики (0)

	// Фильтрация по диапазону дат
	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			query = query.Where("creation_time >= ?", startDate)
		}
	}

	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			query = query.Where("creation_time <= ?", endDate)
		}
	}

	// Фильтрация по статусу, если указан
	if status != 0 {
		query = query.Where("status = ?", status)
	}

	// Выполняем запрос к базе данных
	if err := query.Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (app *App) findProjects(projectID, langID uint) (DbFile, error) {
	var file DbFile
	if err := app.db.Where("project_id = ? AND lang_id = ?", req.ProjectID, req.LangID).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return DbFile{}, gorm.ErrRecordNotFound
		}
		return DbFile{}, err
	}
	return file, nil
}

// Проверка, занят ли логин
func (app *App) isUserLoginTaken(login string) (bool, error) {
	var count int64
	err := app.db.Model(&database.User{}).Where("login = ?", login).Count(&count).Error
	return count > 0, err
}

// Хеширование пароля
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// Создание пользователя в базе данных
func (app *App) createUser(user *database.User) error {
	return app.db.Create(user).Error
}

// Получение ID текущего пользователя (например, из токена)
func getCurrentUserID(c *gin.Context) (uint, error) {
	// Пример: получаем ID пользователя из токена или сессии
	// Это зависит от того, как реализована аутентификация
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("user ID not found in context")
	}
	return userID.(uint), nil
}

// Обновление профиля пользователя
func (app *App) updateUserProfile(userID uint, updatedFields map[string]interface{}) error {
	return app.db.Model(&database.User{}).Where("id = ?", userID).Updates(updatedFields).Error
}

// Получение пользователя по логину
func (app *App) getUserByLogin(login string) (*database.User, error) {
	var user database.User
	err := app.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

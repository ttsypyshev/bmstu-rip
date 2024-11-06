package backend

import (
	"errors"
	"fmt"
	"rip/database"
	"strconv"
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
	db, err := gorm.Open(postgres.Open(FromEnvDB()), &gorm.Config{})
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

func (app *App) getUserByID(userID uint) (DbUser, error) {
	return getByID[DbUser](app, userID)
}

func (app *App) createLang(lang *DbLang) (uint, error) {
	if err := app.db.db.Create(lang).Error; err != nil {
		return 0, err
	}

	return lang.ID, nil
}

func (app *App) updateLang(lang *DbLang) error {
	if err := app.db.db.Save(lang).Error; err != nil {
		return err
	}

	return nil
}

func (app *App) deleteLang(langID uint) error {
	var lang DbLang
	if err := app.db.db.First(&lang, langID).Error; err != nil {
		return err
	}

	lang.Status = false
	if err := app.db.db.Save(&lang).Error; err != nil {
		return err
	}

	return nil
}

func (app *App) createProject(userID uint) (uint, error) {
	projectID, err := findLastDraft(app, userID)
	if err != nil {
		return 0, err
	}

	if projectID == 0 {
		newProject := DbProject{
			UserID:       userID,
			CreationTime: time.Now(),
			Status:       0, // 0 - черновик
		}

		if err := app.db.db.Create(&newProject).Error; err != nil {
			return 0, err
		}

		return newProject.ID, nil
	}

	return projectID, nil
}

func (app *App) updateProject(project *DbProject) error {
	if project.ID == 0 {
		return fmt.Errorf("project ID is required for update")
	}

	if project.Status == 1 || project.Status == 2 {
		now := time.Now()
		project.FormationTime = &now
	} else {
		project.FormationTime = nil
	}

	if project.Status == 3 || project.Status == 4 {
		now := time.Now()
		project.CompletionTime = &now
	} else {
		project.CompletionTime = nil
	}

	if err := app.db.db.Save(project).Error; err != nil {
		return err
	}

	return nil
}

func (app *App) createFile(projectID, langID uint) error {
	newFile := DbFile{
		LangID:    langID,
		ProjectID: projectID,
	}

	if _, err := app.findFile(projectID, langID); err == nil {
		return fmt.Errorf("file with LangID %d and ProjectID %d already exists", newFile.LangID, newFile.ProjectID)
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	if err := app.db.db.Create(&newFile).Error; err != nil {
		return err
	}

	if err := app.db.db.Model(&DbProject{}).Where("id = ?", projectID).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
		return err
	}

	return nil
}

func (app *App) updateFile(file *DbFile) error {
	if file.ID == 0 {
		return fmt.Errorf("file ID is required for update")
	}

	if err := app.db.db.Save(file).Error; err != nil {
		return err
	}

	return nil
}

func (app *App) deleteFile(projectID, fileID uint) error {
	if fileID == 0 {
		return fmt.Errorf("file ID must be greater than 0")
	}

	if err := app.db.db.Delete(&DbFile{}, fileID).Error; err != nil {
		return err
	}

	if err := app.db.db.Model(&DbProject{}).Where("id = ?", projectID).Update("count", gorm.Expr("count - ?", 1)).Error; err != nil {
		return err
	}

	return nil
}

func (app *App) createUser(login, password, name string, is_admin bool) (uint, error) {
	var user DbUser

	if login == "" {
		return 0, fmt.Errorf("login is required")
	}
	user.Login = login

	if password == "" {
		return 0, fmt.Errorf("password is required")
	}
	user.Password = password

	if name != "" {
		user.Name = name
	}

	user.IsAdmin = is_admin

	var existingUser DbUser
	if err := app.db.db.Where("login = ?", login).First(&existingUser).Error; err == nil {
		return 0, fmt.Errorf("user with login %s already exists", login)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	if err := app.db.db.Create(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (app *App) updateUser(user *DbUser) error {
	if user.ID == 0 {
		return fmt.Errorf("user ID is required for update")
	}

	if err := app.db.db.Model(&DbUser{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"name":     user.Name,
		"login":    user.Login,
		"is_admin": user.IsAdmin,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (app *App) deleteUser(userID uint) error {
	if userID == 0 {
		return fmt.Errorf("user ID must be greater than 0")
	}

	if err := app.db.db.Delete(&DbUser{}, userID).Error; err != nil {
		return err
	}

	return nil
}

// Фильтрация услег по запросу
func (app *App) filterLangsByQuery(query string) ([]DbLang, error) {
	var filteredLangs []DbLang
	lowerQuery := "%" + strings.ToLower(query) + "%"

	result := app.db.db.Where("LOWER(name) LIKE ? AND status = ?", lowerQuery, true).Find(&filteredLangs)
	if result.Error != nil {
		return nil, result.Error
	}

	return filteredLangs, nil
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

// Функция для фильтрации проектов по дате и статусу
func (app *App) filterProjects(startDate, endDate, statusStr string) ([]DbProject, error) {
	var projects []DbProject
	query := app.db.db.Model(&DbProject{}).Preload("User").Preload("Moderator").Where("status NOT IN (?)", []int{1, 0})

	if startDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, err
		}
		query = query.Where("creation_time >= ?", start)
	}

	if endDate != "" {
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, err
		}
		query = query.Where("creation_time <= ?", end)
	}

	if statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err != nil {
			return nil, err
		}
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

// Подсчет количества файлов в проекте
func (app *App) getProjectCount(projectID uint) (int64, error) {
	var count int64
	if err := app.db.db.Model(&DbProject{}).Select("count").Where("id = ?", projectID).Scan(&count).Error; err != nil {
		return -1, err
	}

	return count, nil
}

// Получение файлов для проекта
func (app *App) getFilesForProject(projectID uint) ([]DbFile, error) {
	var matchedFiles []DbFile

	if err := app.db.db.Where("project_id = ?", projectID).Preload("Lang").Find(&matchedFiles).Error; err != nil {
		return nil, err
	}

	return matchedFiles, nil
}

// Обновление кода файлов по предоставленным мапам идентификаторов и кода
func (app *App) updateFilesCode(projectID uint, idToCodeMap map[uint]string) error {
	if len(idToCodeMap) == 0 {
		return fmt.Errorf("the map of IDs to codes is empty")
	}

	for id, newCode := range idToCodeMap {
		var file DbFile
		if err := app.db.db.Where("id = ? AND project_id = ?", id, projectID).First(&file).Error; err != nil {
			return err
		}

		file.Code = newCode

		if err := app.db.db.Save(&file).Error; err != nil {
			return fmt.Errorf("failed to update file with id %d: %v", id, err)
		}
	}
	return nil
}

// Поиск конкретного файла
func (app *App) findFile(projectID, langID uint) (DbFile, error) {
	var file DbFile
	if err := app.db.db.Where("project_id = ? AND lang_id = ?", projectID, langID).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return DbFile{}, gorm.ErrRecordNotFound
		}
		return DbFile{}, err
	}
	return file, nil
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

// Поиск юзера по логину
func (app *App) findUserByLogin(login string) (DbUser, error) {
	var user DbUser
	if err := app.db.db.Where("login = ?", login).First(&user).Error; err != nil {
		return DbUser{}, err
	}
	return user, nil
}

// Сравнение переданного пароля с паролем пользователя
func (app *App) matchPassword(login, password string) (bool, DbUser, error) {
	var user DbUser
	if err := app.db.db.Where("login = ?", login).First(&user).Error; err != nil {
		return false, DbUser{}, err
	}

	if user.Password == password {
		return true, user, nil
	}

	return false, DbUser{}, nil
}

// // Хеширование пароля
// func hashPassword(password string) (string, error) {
// 	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	return string(hashed), err
// }

// // Получение ID текущего пользователя (например, из токена)
// func getCurrentUserID(c *gin.Context) (uint, error) {
// 	// Пример: получаем ID пользователя из токена или сессии
// 	// Это зависит от того, как реализована аутентификация
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		return 0, errors.New("user ID not found in context")
// 	}
// 	return userID.(uint), nil
// }

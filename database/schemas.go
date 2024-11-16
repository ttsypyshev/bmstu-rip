package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Status string

const (
	Draft     Status = "draft"
	Deleted   Status = "deleted"
	Created   Status = "created"
	Completed Status = "completed"
	Rejected  Status = "rejected"
)

type JSONB map[string]string

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	var objMap map[string]string
	if err := json.Unmarshal(bytes, &objMap); err == nil {
		*j = objMap
		return nil
	}

	var arr []interface{}
	if err := json.Unmarshal(bytes, &arr); err == nil {
		objMap := make(map[string]string)
		for _, val := range arr {
			if strVal, ok := val.(string); ok {
				parts := strings.SplitN(strVal, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					objMap[key] = value
				}
			} else {
				return errors.New("array elements must be strings in the form 'key:value'")
			}
		}
		*j = objMap
		return nil
	}

	return errors.New("failed to unmarshal JSONB data as either an object or an array")
}

// User представляет пользователя в системе
type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"size:50"`
	Email    string `gorm:"type:text;unique"`
	Login    string `gorm:"size:50;unique;not null"`
	Password string `gorm:"type:text;not null" json:"-"`
	IsAdmin  bool   `gorm:"default:false"`
}

// Lang представляет услугу
type Lang struct {
	ID               uint   `gorm:"primaryKey;autoIncrement"`
	Name             string `gorm:"size:50;not null"`
	ShortDescription string `gorm:"size:255;not null"`
	Description      string `gorm:"type:text;not null"`
	ImgLink          string `gorm:"size:255"`
	Author           string `gorm:"size:50"`
	Year             string `gorm:"size:4"`
	Version          string `gorm:"size:50"`
	List             JSONB  `gorm:"type:jsonb"`
	Status           bool   `gorm:"default:true;not null"`
}

// Project представляет заявку
type Project struct {
	ID               uint      `gorm:"primaryKey;autoIncrement"`
	UserID           uint      `gorm:"not null"`
	User             *User     `gorm:"foreignKey:UserID" json:"User,omitempty"`
	CreationTime     time.Time `gorm:"default:current_timestamp"`
	FormationTime    *time.Time
	CompletionTime   *time.Time
	Status           Status `gorm:"type:project_status;not null"`
	ModeratorID      *uint
	Moderator        *User  `gorm:"foreignKey:ModeratorID" json:"Moderator,omitempty"`
	ModeratorComment string `gorm:"type:text"`
	Count            int    `gorm:"default:0"`
}

// File представляет файл
type File struct {
	ID        uint   `gorm:"primaryKey"`
	LangID    uint   `gorm:"not null"`
	Lang      *Lang  `gorm:"foreignKey:LangID" json:"Lang,omitempty"`
	ProjectID uint   `gorm:"not null"`
	Code      string `gorm:"type:text"`
	FileName  string `gorm:"size:255"`
	FileSize  int64  `gorm:"default:0"`
	Comment   string `gorm:"type:text"`
	AutoCheck int    `gorm:"default:0"`
}

package database

import "time"

// User представляет пользователя в системе
type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"size:255"`
	Login    string `gorm:"size:255;unique;not null"`
	Password string `gorm:"size:255;not null" json:"-"`
	IsAdmin  bool   `gorm:"default:false"`
}

// Lang представляет услугу
type Lang struct {
	ID               uint   `gorm:"primaryKey;autoIncrement"`
	Name             string `gorm:"size:255;not null"`
	ShortDescription string `gorm:"size:255;not null"`
	Description      string `gorm:"type:text;not null"`
	ImgLink          string `gorm:"size:255"`
	Author           string `gorm:"size:255"`
	Year             string `gorm:"type:char(4)"`
	Version          string `gorm:"size:50"`
	List             string `gorm:"type:text"`
	Status           bool   `gorm:"default:true;not null"`
}

// Project представляет заявку
type Project struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	UserID         uint      `gorm:"not null"`
	User           *User     `gorm:"foreignKey:UserID" json:"User,omitempty"`
	CreationTime   time.Time `gorm:"default:current_timestamp"`
	FormationTime  *time.Time
	CompletionTime *time.Time
	Status         int `gorm:"not null"` // 0 - черновик, 1 - удалён, 2 - сформирован, 3 - завершён, 4 - отклонён
	ModeratorID    *uint
	Moderator      *User  `gorm:"foreignKey:ModeratorID" json:"Moderator,omitempty"`
	Comment        string `gorm:"type:text"`
	Count          int    `gorm:"default:0"`
}

// File представляет файл
type File struct {
	ID        uint   `gorm:"primaryKey"`
	LangID    uint   `gorm:"not null"`
	Lang      *Lang  `gorm:"foreignKey:LangID" json:"Lang,omitempty"`
	ProjectID uint   `gorm:"not null"`
	Code      string `gorm:"type:text"`
	AutoCheck int    `gorm:"default:0"`
}

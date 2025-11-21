package models

import "time"

type User struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	Login        string `json:"login" gorm:"unique;size:16"`
	PasswordHash string `json:"password_hash"`
}

type Task struct {
	ID        int        `json:"id" gorm:"primaryKey"`
	UserID    int        `json:"user_id"`
	Task      string     `json:"task"`
	Completed bool       `json:"completed" gorm:"default:false"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt *time.Time `json:"updated_at"`
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

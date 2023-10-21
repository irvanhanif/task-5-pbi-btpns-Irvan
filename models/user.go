package models

import "gorm.io/gorm"

type User struct {
	Id          int64  `gorm:"primaryKey" json:"id"`
	Username    string `gorm:"type:varchar(300);not null" json:"username"`
	Email    string `gorm:"type:varchar(100); not null; unique" json:"email"`
	Password    string `gorm:"type:varchar(300); not null; " json:"-"`
	gorm.Model
}

type UserAuth struct {
	Username    string `json:"username"`
	Email    string `json:"email"`
	Password    string `json:"password"`
}
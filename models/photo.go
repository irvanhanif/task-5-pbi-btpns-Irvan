package models

import "gorm.io/gorm"

type Photo struct {
	Id       int64  `gorm:"primaryKey" json:"id"`
	Title    string `gorm:"type:varchar(300)" json:"title"`
	Caption  string `gorm:"type:varchar(300)" json:"caption"`
	PhotoUrl string `gorm:"type:text" json:"photo_url"`
	UserID   int64 `gorm:"not null" json:"user_id"`
	User User `gorm:"foreignKey:Id;references:UserID; constraint:onUpdate:CASCASE,onDelete:CASCADE"`
	gorm.Model
}

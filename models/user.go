package models

import "gorm.io/gorm"

//数据库user模型
type User struct {
	gorm.Model
	Username      string `gorm:"not null"`
	Password      string `gorm:"not null"`
	FollowCount   int    `gorm:"not null;default:0"`
	FollowerCount int    `gorm:"not null;default:0"`
}

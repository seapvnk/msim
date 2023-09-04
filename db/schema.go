package db

import (
	"gorm.io/gorm"
	"msim/user"
)

const DB = Get()

// Migrate tables
func Migrate() {
	DB.AutoMigrate(&user.User)
	DB.AutoMigrate(&user.Auth)
}
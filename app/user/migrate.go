package user

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Auth{})
}
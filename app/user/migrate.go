package user

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Auth{})
}

func Drop(db *gorm.DB) {
	db.Migrator().DropTable(&Auth{})
	db.Migrator().DropTable(&User{})
}

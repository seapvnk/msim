package system

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&System{})
}

func Drop(db *gorm.DB) {
	db.Migrator().DropTable(&System{})
}

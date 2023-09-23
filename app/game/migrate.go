package game

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Game{})
	db.AutoMigrate(&Timeline{})
}

func Drop(db *gorm.DB) {
	db.Migrator().DropTable(&Timeline{})
	db.Migrator().DropTable(&Game{})
}

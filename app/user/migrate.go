package user

import "msim/db"


func Migrate() {
	DB := db.ORM()

	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Auth{})
}
package db

import "gorm.io/gorm"

// Switch between database used in build
// TODO: create another databases to switch in depending on environment
func ORM() *gorm.DB {
	localDB, _ := LocalDB()
	return localDB
}
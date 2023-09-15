package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Get ORM instance for in-memory database
func InMemoryDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	return db, err
}

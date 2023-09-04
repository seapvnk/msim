package db

import (
	"gorm.io/gorm"
	"msim/config"
)

// Get database, base case
func Get() (*gorm.DB, error) {
	return LocalDB()
}

// Get database for another environments
// TODO: database for each environment
//	* server:	Postgresql or MySQL.
//	* tests:	In-memory setup for sqlite 
//				or temporary file automatically generated.
func Get(environment config.Environment) (*gorm.DB, error) {
	switch environment {
	case config.Local:
		return LocalDB()
	// TODO
	case config.Server:
		return LocalDB()
	// TODO
	case config.Test:
		return LocalDB()
	default:
		return Get()
	}
}
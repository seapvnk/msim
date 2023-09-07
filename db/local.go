package db

import (
	"os"
	"fmt"
	"errors"

	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

const storageFilename = "storage"
const databaseFilename = "db.sqlite"

// Get ORM instance
func LocalDB() (*gorm.DB, error) {
	filePath := fmt.Sprintf("%s/%s", storageFilename, databaseFilename)
	pathWithForeignKeyArgs := fmt.Sprintf("%s?_foreign_keys=on", filePath)
	db, err := gorm.Open(sqlite.Open(pathWithForeignKeyArgs), &gorm.Config{})

	if err != nil {
		fmt.Println("Error openning local database")
	}

	return db, err
}

// Setup environment for sqlite database
func LocalDBSetup() error {
	if createStorageFolderIfDoesntExists(storageFilename) {
		filePath := fmt.Sprintf("%s/%s", storageFilename, databaseFilename)
		createSqliteDatabase(filePath)

		return nil
	}

	return errors.New("can't setup database environment")
}

// Create storage folder if doesnt exists, return true if folder exists
func createStorageFolderIfDoesntExists(path string) bool {
	if _, err := os.Stat(storageFilename); os.IsNotExist(err) {
		if err := os.Mkdir(storageFilename, os.ModePerm); err != nil {
			fmt.Println("Error creating storage folder", err)
			return false
		}
		fmt.Println("Storage folder created")
	} else if err != nil {
		fmt.Println("Error verifying storage folder", err)
		return false
	} else {
		fmt.Println("Storage folder already exists")
	}

	return true
}

// Create sqlite database, return true if file exists
func createSqliteDatabase(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		newFile, err := os.Create(path)
		if err != nil {
			errorMessage := fmt.Sprintf("Error creating database file in \"%s\"", path)
			fmt.Println(errorMessage)
			return false
		}
		defer newFile.Close()

		fmt.Println("Storage file created")
	} else if err != nil {
		fmt.Println("Error verifying storage file")
		return false
	} else {
		fmt.Println("Storage file already exists")
	}

	return true
}
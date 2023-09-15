package db

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

// Test LocalDB
func TestLocalDB(t *testing.T) {
	t.Run("Should return *gorm.DB instance when database exists", func(t *testing.T) {
		setupLocal()

		result, _ := LocalDB()
		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*gorm.DB)(nil))

		if resultType != expectedType {
			errFormated := `LocalDB() returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		teardownLocal()
	})

	t.Run("Should return error when database doesnt exists", func(t *testing.T) {
		_, err := LocalDB()

		if err == nil {
			t.Fatal("LocalDB() expects error when theres no database")
		}

		teardownLocal()
	})
}

// Test LocalDBSetup
func TestLocalDBSetup(t *testing.T) {
	t.Run("Should setup environment for local database", func(t *testing.T) {
		result := LocalDBSetup()

		if result != nil {
			t.Fatal("LocalDBSetup() must create database file without errors")
		}

		teardownLocal()
	})
}

// Runs after each tests
func teardownLocal() {
	folder := "storage"

	err := os.RemoveAll(folder)
	if err != nil {
		panic("Error excluding database folder")
	} else {
		fmt.Println("Folder excluded")
	}
}

// Create database file for tests
func setupLocal() {
	if err := os.Mkdir("storage", os.ModePerm); err != nil {
		panic("Error creating storage folder")
	}
	fmt.Println("Storage folder created")

	newFile, err := os.Create("storage/db.sqlite")
	if err != nil {
		panic("Error creating database file")
	}
	defer newFile.Close()
}

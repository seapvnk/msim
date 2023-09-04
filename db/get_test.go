package db

import (
	"fmt"
	"testing"
	"os"
	"reflect"

	"gorm.io/gorm"
)

// Test Get
// TODO: Test for each environment
func TestLocalDB(t *testing.T) {
	t.Run("Should return *gorm.DB instance when database exists", func(t *testing.T) {
		setupCreateDatabaseFile()
	
		result, _ := Get()
		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*gorm.DB)(nil))
	
		if resultType != expectedType {
			errFormated := `LocalDB() returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}
	
		afterEach()
	})

	t.Run("Should return error when database doesnt exists", func(t *testing.T) {
		_, err := Get()

		if err == nil {
			t.Fatal("LocalDB() expects error when theres no database")
		}
	
		afterEach()
	})
}

// Runs after each tests
func afterEach() {
	folder := "storage"

	err := os.RemoveAll(folder)
	if err != nil {
		panic("Error excluding database folder")
	} else {
		fmt.Println("Folder excluded")
	}
}

// Create database file for tests
func setupCreateDatabaseFile() {
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
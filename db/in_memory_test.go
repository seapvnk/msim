package db

import (
	"reflect"
	"testing"

	"gorm.io/gorm"
)

// Test InMemoryDB
func TestInMemoryDB(t *testing.T) {
	t.Run("Should return *gorm.DB instance", func(t *testing.T) {

		result, _ := InMemoryDB()
		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*gorm.DB)(nil))

		if resultType != expectedType {
			errFormated := `InMemoryDB() returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}
	})
}

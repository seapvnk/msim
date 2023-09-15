package system

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"msim/db"
)

// Test create.
func TestCreateRepository(t *testing.T) {
	t.Run("Should create a system variable", func(t *testing.T) {
		repository, DB := CreateSystemRepository()
		result, err := repository.Create(&SystemEntity{ID: uuid.New(), Key: "abcd", Value: "abcd", Type: "string"})

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*SystemEntity)(nil))

		if err != nil {
			t.Fatal(err)
		}

		if resultType != expectedType {
			errFormated := `create returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		var countEnvs int64
		DB.Model(&System{}).Count(&countEnvs)

		if countEnvs != 1 {
			errFormated := `created count returns %d, expects 1`
			t.Fatalf(errFormated, countEnvs)
		}
	})

	t.Run("Should not create a system variable when key already exists", func(t *testing.T) {
		repository, DB := CreateSystemRepository()

		DB.Create(&System{ID: uuid.New(), Key: "abcd", Value: "abcd", Type: "string"})
		_, err := repository.Create(&SystemEntity{ID: uuid.New(), Key: "abcd", Value: "abcd", Type: "string"})

		if err == nil {
			t.Fatal("Should not create another env with same key")
		}
	})
}

// Test getByKey.
func TestGetByKey(t *testing.T) {
	t.Run("Should get an env by key when exists", func(t *testing.T) {
		repository, DB := CreateSystemRepository()

		created := System{ID: uuid.New(), Key: "key", Value: "val", Type: "type"}
		DB.Create(&created)

		result, err := repository.GetByKey(created.Key)

		if err != nil {
			t.Fatal(err)
		}

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*SystemEntity)(nil))

		if resultType != expectedType {
			errFormated := `find returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		if result.Key != created.Key || result.Value != created.Value || result.Type != created.Type {
			t.Fatal("System should should match")
		}
	})

	t.Run("Should not get an env by key when doesnt exists", func(t *testing.T) {
		repository, _ := CreateSystemRepository()
		result, err := repository.GetByKey("key2")

		if err == nil {
			t.Fatal("Should throw error")
		}

		if result != nil {
			t.Fatal("System should be nil")
		}
	})
}

// Test GetAll.
func TestGetAll(t *testing.T) {
	t.Run("Should get all system variables", func(t *testing.T) {
		repository, DB := CreateSystemRepository()

		variables := []System{
			{ID: uuid.New(), Key: "key1", Value: "value1", Type: "string"},
			{ID: uuid.New(), Key: "key2", Value: "value2", Type: "string"},
			{ID: uuid.New(), Key: "key3", Value: "value3", Type: "string"},
		}
		DB.Create(&variables)

		result, err := repository.GetAll()

		if err != nil {
			t.Fatal(err)
		}

		expectedLength := len(variables)
		if len(result) != expectedLength {
			t.Fatalf("Expected %d variables, got %d", expectedLength, len(result))
		}

		for i := 0; i < expectedLength; i++ {
			if result[i].Key != variables[i].Key || result[i].Value != variables[i].Value || result[i].Type != variables[i].Type {
				t.Fatalf("Mismatch in variable %d", i+1)
			}
		}
	})

	t.Run("Should retrieve an empty list of system variables", func(t *testing.T) {
		repository, _ := CreateSystemRepository()

		result, err := repository.GetAll()

		if err != nil {
			t.Fatal(err)
		}

		expectedLength := 0
		if len(result) != expectedLength {
			t.Fatalf("Expected %d variables, got %d", expectedLength, len(result))
		}
	})
}

// Test UpdateValueByKey.
func TestUpdateValueByKey(t *testing.T) {
	t.Run("Should update the value of an existing system variable", func(t *testing.T) {
		repository, DB := CreateSystemRepository()

		initialVariable := System{ID: uuid.New(), Key: "key1", Value: "value1", Type: "string"}
		DB.Create(&initialVariable)

		newValue := "updatedValue"
		result, err := repository.UpdateValueByKey(initialVariable.Key, newValue)

		if err != nil {
			t.Fatal(err)
		}

		if result.Value != newValue {
			t.Fatalf("Expected value to be %s, got %s", newValue, result.Value)
		}

		var updatedVariable System
		DB.Where("key = ?", initialVariable.Key).First(&updatedVariable)
		if updatedVariable.Value != newValue {
			t.Fatalf("Expected value in database to be %s, got %s", newValue, updatedVariable.Value)
		}
	})

	t.Run("Should return an error if key does not exist", func(t *testing.T) {
		repository, _ := CreateSystemRepository()
		newValue := "updatedValue"

		_, err := repository.UpdateValueByKey("NonExistentKey", newValue)

		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
	})
}

// Create repository and test database.
func CreateSystemRepository() (*SystemRepository, *gorm.DB) {
	DB, _ := db.InMemoryDB()

	Drop(DB)
	Migrate(DB)

	return &SystemRepository{db: DB}, DB
}

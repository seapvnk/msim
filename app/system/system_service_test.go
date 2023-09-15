package system

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"msim/db"
)

// Test create.
func TestCreateService(t *testing.T) {
	t.Run("Should create a system env", func(t *testing.T) {
		service, DB := CreateSystemService()
		result, err := service.Create(&SystemEnvDTO{"test", "teste2", "string"})

		if err != nil {
			t.Fatal(err)
		}

		var find System
		DB.First(&find)

		if find.Key != result.Key || find.Value != result.Value || find.Type != result.Type {
			t.Fatal("Should create this exact model in database")
		}
	})

	t.Run("Should not create a system env when already exists", func(t *testing.T) {
		service, DB := CreateSystemService()

		DB.Create(&System{Key: "test", Value: "12345", Type: "int"})
		result, err := service.Create(&SystemEnvDTO{"test", "teste2", "string"})

		if err == nil {
			t.Fatal("Should not create a system env")
		}

		if result != nil {
			t.Fatal("Result should be nil")
		}
	})
}

// Test GetByKey.
func TestGetSystemByKeyService(t *testing.T) {
	t.Run("Should get system variable by key when it exists", func(t *testing.T) {
		service, DB := CreateSystemService()

		created := System{ID: uuid.New(), Key: "test", Value: "123", Type: "string"}
		DB.Create(&created)

		result, _ := service.GetByKey(&SystemKeyDTO{Key: created.Key})

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*SystemEntity)(nil))

		if resultType != expectedType {
			errFormatted := `find returns type of %s, expects type of %s`
			t.Fatalf(errFormatted, resultType, expectedType)
		}

		if result.Key != created.Key || result.Value != created.Value || result.Type != created.Type {
			t.Fatal("System variable should match")
		}
	})

	t.Run("Should not get system variable by key when it doesn't exist", func(t *testing.T) {
		service, _ := CreateSystemService()
		result, err := service.GetByKey(&SystemKeyDTO{Key: "NonExistentKey"})

		if err == nil {
			t.Fatal("Should throw error")
		}

		if result != nil {
			t.Fatal("System variable should be nil")
		}
	})
}

// Test GetAll
func TestGetAllService(t *testing.T) {
	t.Run("Should retrieve all system variables", func(t *testing.T) {
		service, DB := CreateSystemService()

		// Create some system variables
		variables := []System{
			{ID: uuid.New(), Key: "key1", Value: "value1", Type: "string"},
			{ID: uuid.New(), Key: "key2", Value: "value2", Type: "string"},
			{ID: uuid.New(), Key: "key3", Value: "value3", Type: "string"},
		}
		DB.Create(&variables)

		result, _ := service.GetAll()

		expectedLength := len(variables)
		if len(result) != expectedLength {
			t.Fatalf("Expected %d variables, got %d", expectedLength, len(result))
		}

		// Verify if keys and values match
		for i := 0; i < expectedLength; i++ {
			if result[i].Key != variables[i].Key || result[i].Value != variables[i].Value || result[i].Type != variables[i].Type {
				t.Fatalf("Mismatch in variable %d", i+1)
			}
		}
	})

	t.Run("Should retrieve an empty list of system variables", func(t *testing.T) {
		service, _ := CreateSystemService()

		result, _ := service.GetAll()

		expectedLength := 0
		if len(result) != expectedLength {
			t.Fatalf("Expected %d variables, got %d", expectedLength, len(result))
		}
	})
}

// Test UpdateValueByKey.
func TestUpdateValueByKeyService(t *testing.T) {
	t.Run("Should update the value of an existing system variable", func(t *testing.T) {
		service, DB := CreateSystemService()

		initialVariable := System{ID: uuid.New(), Key: "key1", Value: "value1", Type: "string"}
		DB.Create(&initialVariable)

		newValue := "updatedValue"
		dto := &SystemKeyUpdateDTO{
			Key:   initialVariable.Key,
			Value: newValue,
		}
		result, err := service.UpdateValueByKey(dto)

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
		service, _ := CreateSystemService()
		dto := &SystemKeyUpdateDTO{
			Key:   "NonExistentKey",
			Value: "updatedValue",
		}

		_, err := service.UpdateValueByKey(dto)

		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
	})
}

// Create service and test database.
func CreateSystemService() (*SystemService, *gorm.DB) {
	DB, _ := db.InMemoryDB()

	Drop(DB)
	Migrate(DB)

	systemRepo := &SystemRepository{db: DB}
	return &SystemService{systemRepository: systemRepo}, DB
}

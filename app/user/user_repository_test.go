package user

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"msim/db"
)

// Test create.
func TestCreate(t *testing.T) {
	t.Run("Should create an user and return an user entity", func(t *testing.T) {
		repository, DB := CreateUserRepository()
		result, err := repository.Create(&UserEntity{Name: "test", password: "12345"})

		if err != nil {
			t.Fatal(err)
		}

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*UserEntity)(nil))

		if resultType != expectedType {
			errFormated := `create returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		var countUsers int64
		DB.Model(&User{}).Count(&countUsers)

		if countUsers != 1 {
			errFormated := `created count returns %d, expects 1`
			t.Fatalf(errFormated, countUsers)
		}
	})

	t.Run("Should not create an user when user name already exists", func(t *testing.T) {
		repository, DB := CreateUserRepository()

		DB.Create(&User{Name: "test", Password: "12345"})
		_, err := repository.Create(&UserEntity{Name: "test", password: "12345"})

		if err == nil {
			t.Fatal("Should not create another user with same name")
		}
	})
}

// Test getAll.
func TestGetAll(t *testing.T) {
	t.Run("Should get all users when theres users", func(t *testing.T) {
		repository, DB := CreateUserRepository()

		DB.Create(&User{ID: uuid.New(), Name: "Test", Password: "12345"})
		DB.Create(&User{Name: "test2", Password: "12345"})

		result, err := repository.GetAll()

		if err != nil {
			t.Fatal(err)
		}

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf(([]*UserEntity)(nil))

		if resultType != expectedType {
			errFormated := `find returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		if len(result) != 2 {
			t.Fatal("Should return the only 2 created users")
		}

		if result[0].Name != "test1" && result[1].Name != "test2" {
			t.Fatal("Users should match")
		}
	})

	t.Run("Should get an empty array when theres no user", func(t *testing.T) {
		repository, _ := CreateUserRepository()

		result, err := repository.GetAll()

		if err != nil {
			t.Fatal(err)
		}

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf(([]*UserEntity)(nil))

		if resultType != expectedType {
			errFormated := `find returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		if len(result) != 0 {
			t.Fatal("Should return zero users")
		}
	})
}

// Test getById.
func TestGetById(t *testing.T) {
	t.Run("Should get an user by id when exists", func(t *testing.T) {
		repository, DB := CreateUserRepository()

		created := User{ID: uuid.New(), Name: "Testll", Password: "12345"}
		DB.Create(&created)

		result, err := repository.GetById(created.ID)

		if err != nil {
			t.Fatal(err)
		}

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*UserEntity)(nil))

		if resultType != expectedType {
			errFormated := `find returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		if result.Name != created.Name {
			t.Fatal("User name should match should match")
		}
	})

	t.Run("Should not get an user by id when doesnt exists", func(t *testing.T) {
		repository, _ := CreateUserRepository()

		result, err := repository.GetById(uuid.New())

		if err == nil {
			t.Fatal("Should throw error")
		}

		if result != nil {
			t.Fatal("User should be nil")
		}
	})
}

// Test getByName.
func TestGetByName(t *testing.T) {
	t.Run("Should get an user by name when exists", func(t *testing.T) {
		repository, DB := CreateUserRepository()

		created := User{ID: uuid.New(), Name: "Test867", Password: "12345"}
		DB.Create(&created)

		result, err := repository.GetByName(created.Name)

		if err != nil {
			t.Fatal(err)
		}

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*UserEntity)(nil))

		if resultType != expectedType {
			errFormated := `find returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		if result.Name != created.Name {
			t.Fatal("User name should match")
		}
	})

	t.Run("Should not get an user by name when doesnt exists", func(t *testing.T) {
		repository, _ := CreateUserRepository()
		result, err := repository.GetByName("Test2")

		if err == nil {
			t.Fatal("Should throw error")
		}

		if result != nil {
			t.Fatal("User should be nil")
		}
	})
}

// Create repository and test database.
func CreateUserRepository() (*UserRepository, *gorm.DB) {
	DB, _ := db.InMemoryDB()

	Drop(DB)
	Migrate(DB)

	return &UserRepository{db: DB}, DB
}

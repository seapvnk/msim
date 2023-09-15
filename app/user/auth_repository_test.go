package user

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"msim/db"
)

// Test create.
func TestCreateAuth(t *testing.T) {
	t.Run("Should create an auth and return a code when theres an user", func(t *testing.T) {
		repository, DB := CreateAuthRepository()

		created := &User{Name: "test1", Password: "12345"}
		DB.Create(&created)

		result, err := repository.Create(created.ID)

		var find Auth
		DB.First(&find)

		if err != nil {
			t.Fatal(err)
		}

		if find.Code != result {
			t.Fatal("Created code should match returned code")
		}
	})

	t.Run("Should not create an auth when theres no user", func(t *testing.T) {
		repository, _ := CreateAuthRepository()

		_, err := repository.Create(uuid.Nil)

		if err != nil {
			t.Fatal("Should return an error")
		}
	})
}

// Test getAuthUser.
func TestGetAuthUser(t *testing.T) {
	t.Run("Should get auth user when theres an user", func(t *testing.T) {
		repository, DB := CreateAuthRepository()

		createdUser := User{Name: "test1a", Password: "12345"}
		DB.Create(&createdUser)

		createdAuth := Auth{ID: uuid.New(), Code: uuid.New(), User: createdUser}
		DB.Create(&createdAuth)

		result, err := repository.GetAuthUser(createdAuth.Code)

		if err != nil {
			t.Fatal(err)
		}

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*UserEntity)(nil))

		if resultType != expectedType {
			errFormated := `find returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		if result.Name != createdUser.Name {
			t.Fatal("User name should match should match")
		}
	})

	t.Run("Should not get auth user when theres no auth code or user", func(t *testing.T) {
		repository, _ := CreateAuthRepository()

		result, err := repository.GetAuthUser(uuid.New())

		if err == nil {
			t.Fatal(err)
		}

		if result != nil {
			t.Fatal("Should not find am user")
		}
	})

	t.Run("Should not get auth user when auth code is expired", func(t *testing.T) {
		repository, DB := CreateAuthRepository()

		createdUser := User{Name: "test1", Password: "12345"}
		DB.Create(&createdUser)

		mockTimeCreated := time.Date(2002, time.January, 9, 0, 0, 0, 0, time.UTC)
		createdAuth := Auth{ID: uuid.New(), Code: uuid.New(), User: createdUser}
		createdAuth.CreatedAt = mockTimeCreated

		DB.Create(&createdAuth)

		result, err := repository.GetAuthUser(createdAuth.Code)

		if err == nil {
			t.Fatal(err)
		}

		if result != nil {
			t.Fatal("Should not find am user")
		}
	})
}

// Create repository and test database.
func CreateAuthRepository() (*AuthRepository, *gorm.DB) {
	DB, _ := db.InMemoryDB()

	Drop(DB)
	Migrate(DB)

	return &AuthRepository{db: DB}, DB
}

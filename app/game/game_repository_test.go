package game

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"msim/db"
)

// Test create.
func TestCreateRepository(t *testing.T) {
	t.Run("Should create a game slot", func(t *testing.T) {
		repository, DB := CreateGameRepository()
		result, err := repository.Create(&GameEntity{ID: uuid.New(), Name: "abcd"})

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*GameEntity)(nil))

		if err != nil {
			t.Fatal(err)
		}

		if resultType != expectedType {
			errFormated := `create returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		var countGames int64
		DB.Model(&Game{}).Count(&countGames)

		if countGames != 1 {
			errFormated := `created count returns %d, expects 1`
			t.Fatalf(errFormated, countGames)
		}
	})
}

// Test GetById.
func TestGetById(t *testing.T) {
	t.Run("Should get a game by ID when it exists", func(t *testing.T) {
		repository, DB := CreateGameRepository()

		created := Game{ID: uuid.New(), Name: "test"}
		DB.Create(&created)

		result, err := repository.GetById(created.ID)

		if err != nil {
			t.Fatal(err)
		}

		if result.ID != created.ID || result.Name != created.Name {
			t.Fatal("Game should match")
		}
	})

	t.Run("Should return an error if game with given ID doesn't exist", func(t *testing.T) {
		repository, _ := CreateGameRepository()
		nonExistentID := uuid.New()

		result, err := repository.GetById(nonExistentID)

		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		if result != nil {
			t.Fatal("Game should be nil")
		}
	})
}

// Create repository and test database.
func CreateGameRepository() (*GameRepository, *gorm.DB) {
	DB, _ := db.InMemoryDB()

	Drop(DB)
	Migrate(DB)

	return &GameRepository{db: DB}, DB
}

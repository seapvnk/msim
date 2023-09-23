package game

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"msim/db"
)

// Test Create.
func TestCreateService(t *testing.T) {
	t.Run("Should create a game", func(t *testing.T) {
		service, DB := CreateGameService()
		result, err := service.Create(&GameCreateDTO{Name: "test"})

		if err != nil {
			t.Fatal(err)
		}

		var find Game
		DB.First(&find)

		if find.Name != result.Name {
			t.Fatal("Should create this exact model in database")
		}
	})
}

// Test GetById.
func TestGetGameByIdService(t *testing.T) {
	t.Run("Should get a game by ID when it exists", func(t *testing.T) {
		service, DB := CreateGameService()

		created := Game{ID: uuid.New(), Name: "test"}
		DB.Create(&created)

		result, _ := service.GetById(&GameGetDTO{ID: created.ID})

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*GameEntity)(nil))

		if resultType != expectedType {
			errFormatted := `find returns type of %s, expects type of %s`
			t.Fatalf(errFormatted, resultType, expectedType)
		}

		if result.ID != created.ID || result.Name != created.Name {
			t.Fatal("Game should match")
		}
	})

	t.Run("Should not get a game by ID when it doesn't exist", func(t *testing.T) {
		service, _ := CreateGameService()
		result, err := service.GetById(&GameGetDTO{ID: uuid.New()})

		if err == nil {
			t.Fatal("Should throw error")
		}

		if result != nil {
			t.Fatal("Game should be nil")
		}
	})
}

// Create service and test database.
func CreateGameService() (*GameService, *gorm.DB) {
	DB, _ := db.InMemoryDB()

	Drop(DB)
	Migrate(DB)

	gameRepo := &GameRepository{db: DB}
	return &GameService{gameRepository: gameRepo}, DB
}

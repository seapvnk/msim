package game

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"msim/db"
)

// Test Create.
func TestCreateTimelineRepository(t *testing.T) {
	t.Run("Should create a timeline", func(t *testing.T) {
		repository, DB := CreateTimelineRepositoryMock()

		game := Game{ID: uuid.New(), Name: "test"}
		DB.Create(&game)

		timelineEntity := &TimelineEntity{
			ID:          uuid.New(),
			Signature:   uuid.New(),
			StartDate:   time.Now(),
			CurrentDate: time.Now(),
			GameID:      game.ID,
		}

		result, err := repository.Create(timelineEntity)

		if err != nil {
			t.Fatal(err)
		}

		if result.ID != timelineEntity.ID || result.Signature != timelineEntity.Signature {
			t.Fatal("Created timeline does not match")
		}
	})
}

// Test FindBySignature.
func TestFindBySignature(t *testing.T) {
	t.Run("Should find a timeline by signature when it exists", func(t *testing.T) {
		repository, DB := CreateTimelineRepositoryMock()

		game := Game{ID: uuid.New(), Name: "test"}
		DB.Create(&game)

		timeline := &Timeline{
			ID:          uuid.New(),
			Signature:   uuid.New(),
			StartDate:   time.Now(),
			CurrentDate: time.Now(),
			GameID:      uuid.New(),
		}
		DB.Create(timeline)

		result, err := repository.FindBySignature(timeline.Signature)

		if err != nil {
			t.Fatal(err)
		}

		if result.ID != timeline.ID || result.Signature != timeline.Signature {
			t.Fatal("Found timeline does not match")
		}
	})

	t.Run("Should return an error if timeline with signature does not exist", func(t *testing.T) {
		repository, _ := CreateTimelineRepositoryMock()
		nonExistentSignature := uuid.New()
		_, err := repository.FindBySignature(nonExistentSignature)

		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
	})
}

// Test FindBySignatureInIndex.
func TestFindBySignatureInIndex(t *testing.T) {
	t.Run("Should find a timeline by signature when it exists", func(t *testing.T) {
		repository, DB := CreateTimelineRepositoryMock()

		game := Game{ID: uuid.New(), Name: "test"}
		DB.Create(&game)

		id1 := uuid.New()
		id2 := uuid.New()
		signature := uuid.New()

		DB.Create(&Timeline{
			ID:          id2,
			Signature:   signature,
			StartDate:   time.Now(),
			CurrentDate: time.Now(),
			GameID:      uuid.New(),
		})

		DB.Create(&Timeline{
			ID:          id1,
			Signature:   signature,
			StartDate:   time.Now(),
			CurrentDate: time.Now(),
			GameID:      uuid.New(),
		})

		result, err := repository.FindBySignatureInIndex(signature, uint(1))

		if err != nil {
			t.Fatal(err)
		}

		if result.ID != id2 || result.Signature != signature {
			t.Fatal("Found timeline does not match")
		}
	})

	t.Run("Should return an error if timeline with signature does not exist", func(t *testing.T) {
		repository, _ := CreateTimelineRepositoryMock()
		nonExistentSignature := uuid.New()
		_, err := repository.FindBySignatureInIndex(nonExistentSignature, uint(2))

		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
	})
}

// Test FindAllByGameID.
func TestFindAllByGameID(t *testing.T) {
	t.Run("Should find all timelines by game ID", func(t *testing.T) {
		repository, DB := CreateTimelineRepositoryMock()

		game := Game{ID: uuid.New(), Name: "test"}
		DB.Create(&game)

		gameID := game.ID
		timelines := []*Timeline{
			{ID: uuid.New(), Signature: uuid.New(), StartDate: time.Now().Add(2 * time.Minute), CurrentDate: time.Now(), GameID: gameID},
			{ID: uuid.New(), Signature: uuid.New(), StartDate: time.Now(), CurrentDate: time.Now(), GameID: gameID},
		}
		DB.Create(&timelines)

		result, err := repository.FindAllByGameID(gameID)

		if err != nil {
			t.Fatal(err)
		}

		if len(result) != len(timelines) {
			t.Fatalf("Expected %d timelines, got %d", len(timelines), len(result))
		}

		for i, res := range result {
			if res.ID != timelines[i].ID || res.Signature != timelines[i].Signature {
				t.Fatalf("Timeline %d does not match", i+1)
			}
		}
	})

	t.Run("Should return an empty list if no timelines exist for game ID", func(t *testing.T) {
		repository, DB := CreateTimelineRepositoryMock()

		game := Game{ID: uuid.New(), Name: "test"}
		DB.Create(&game)

		gameID := game.ID
		result, err := repository.FindAllByGameID(gameID)

		if err != nil {
			t.Fatal(err)
		}

		if len(result) != 0 {
			t.Fatalf("Expected 0 timelines, got %d", len(result))
		}
	})
}

// Create repository and test database.
func CreateTimelineRepositoryMock() (*TimelineRepository, *gorm.DB) {
	DB, _ := db.InMemoryDB()

	Drop(DB)
	Migrate(DB)

	return &TimelineRepository{db: DB}, DB
}

package game

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Game struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string
	Timelines []Timeline `gorm:"foreignKey:GameID"`
}

type GameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (repository *GameRepository) Create(g *GameEntity) (*GameEntity, error) {
	gameModel := &Game{ID: g.ID, Name: g.Name}
	result := repository.db.Create(gameModel)

	if result.Error != nil {
		return nil, result.Error
	}

	return g, nil
}

func (repository *GameRepository) GetById(id uuid.UUID) (*GameEntity, error) {
	var model Game

	result := repository.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &GameEntity{ID: model.ID, Name: model.Name}, nil
}

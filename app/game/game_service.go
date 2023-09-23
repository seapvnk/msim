package game

import (
	"github.com/google/uuid"
	"msim/app/shared"
)

type GameEntity struct {
	ID   uuid.UUID
	Name string
}

type GameService struct {
	gameRepository *GameRepository
}

type GameCreateDTO struct {
	Name string
}

// Create and return an save slot.
func (service *GameService) Create(g *GameCreateDTO) (*GameEntity, *shared.Exception) {
	entity := &GameEntity{ID: uuid.New(), Name: g.Name}
	result, err := service.gameRepository.Create(entity)

	if err != nil {
		return nil, shared.InternalErrorException()
	}

	return result, nil
}

type GameGetDTO struct {
	ID uuid.UUID
}

// Get a save slot by id.
func (service *GameService) GetById(g *GameGetDTO) (*GameEntity, *shared.Exception) {
	result, err := service.gameRepository.GetById(g.ID)

	if err != nil {
		return nil, shared.DefaultException(shared.NOT_FOUND_EX, "save slot")
	}

	return result, nil
}

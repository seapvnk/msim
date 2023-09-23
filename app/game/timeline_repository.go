package game

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Timeline struct {
	gorm.Model
	ID          uuid.UUID `gorm:"primaryKey"`
	StartDate   time.Time
	CurrentDate time.Time
	Signature   uuid.UUID
	GameID      uuid.UUID `gorm:"foreignKey"`
	Game        Game      `gorm:"foreignKey:GameID"`
	ParentID    uuid.UUID `gorm:"type:uuid;index"`
}

type TimelineRepository struct {
	db *gorm.DB
}

func CreateTimelineRepository(db *gorm.DB) *TimelineRepository {
	return &TimelineRepository{db: db}
}

// Create a timeline.
func (repository *TimelineRepository) Create(t *TimelineEntity) (*TimelineEntity, error) {
	timeline := &Timeline{
		ID:          t.ID,
		Signature:   t.Signature,
		StartDate:   t.StartDate,
		CurrentDate: t.CurrentDate,
		GameID:      t.GameID,
		ParentID:    t.ParentID,
	}

	result := repository.db.Create(timeline)
	if result.Error != nil {
		return nil, result.Error
	}

	return &TimelineEntity{
		ID:          timeline.ID,
		StartDate:   timeline.StartDate,
		CurrentDate: timeline.CurrentDate,
		Signature:   timeline.Signature,
		ParentID:    timeline.ParentID,
	}, nil
}

// Find a timeline by signature.
func (repository *TimelineRepository) FindBySignature(signature uuid.UUID) (*TimelineEntity, error) {
	var model Timeline

	result := repository.db.
		Where("signature = ?", signature).
		Order("created_at desc").
		First(&model)

	if result.Error != nil {
		return nil, result.Error
	}

	return &TimelineEntity{
		ID:          model.ID,
		StartDate:   model.StartDate,
		CurrentDate: model.CurrentDate,
		Signature:   model.Signature,
		ParentID:    model.ParentID,
	}, nil
}

// Find a timeline by signature in a index.
func (repository *TimelineRepository) FindBySignatureInIndex(signature uuid.UUID, index uint) (*TimelineEntity, error) {
	var model Timeline

	result := repository.db.
		Where("signature = ?", signature).
		Order("created_at desc").
		Limit(1).
    	Offset(int(index)).
		First(&model)

	if result.Error != nil {
		return nil, result.Error
	}

	return &TimelineEntity{
		ID:          model.ID,
		StartDate:   model.StartDate,
		CurrentDate: model.CurrentDate,
		Signature:   model.Signature,
		ParentID:    model.ParentID,
	}, nil
}

// Find all timelines by game id.
func (repository *TimelineRepository) FindAllByGameID(gameID uuid.UUID) ([]*TimelineEntity, error) {
	var (
		models    []*Timeline
		timelines []*TimelineEntity
	)

	result := repository.db.
		Where("game_id = ?", gameID).
		Order("created_at desc").
		Group("signature").
		Find(&models)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, model := range models {
		timelines = append(timelines, &TimelineEntity{
			ID:          model.ID,
			StartDate:   model.StartDate,
			CurrentDate: model.CurrentDate,
			Signature:   model.Signature,
			ParentID:    model.ParentID,
		})
	}

	return timelines, nil
}

package game

import (
	"time"

	"github.com/google/uuid"
	"msim/app/shared"
)

type TimelineEntity struct {
	ID          uuid.UUID
	StartDate   time.Time
	CurrentDate time.Time
	Signature   uuid.UUID
	ParentID    uuid.UUID
	GameID      uuid.UUID
}

type TimelineService struct {
	timelineRepository *TimelineRepository
}

func NewTimelineService(repository *TimelineRepository) *TimelineService {
	return &TimelineService{timelineRepository: repository}
}

type CreateTimelineDTO struct {
	GameID    uuid.UUID
	StartDate time.Time
}

// Create a new timeline.
func (service *TimelineService) Create(t *CreateTimelineDTO) (*TimelineEntity, *shared.Exception) {
	result, err := service.timelineRepository.Create(&TimelineEntity{
		ID:          uuid.New(),
		Signature:   uuid.New(),
		GameID:      t.GameID,
		StartDate:   t.StartDate,
		CurrentDate: t.StartDate,
	})

	if err != nil {
		return nil, shared.InternalErrorException()
	}

	return result, nil
}

type BranchTimelineDTO struct {
	Signature uuid.UUID
}

// Branch timeline.
func (service *TimelineService) Branch(t *BranchTimelineDTO) (*TimelineEntity, *shared.Exception) {
	timeline, err := service.timelineRepository.FindBySignature(t.Signature)

	if err != nil {
		return nil, shared.DefaultException(shared.NOT_FOUND_EX, "timeline")
	}

	result, err := service.timelineRepository.Create(&TimelineEntity{
		ID:          uuid.New(),
		Signature:   uuid.New(),
		GameID:      timeline.GameID,
		StartDate:   timeline.StartDate,
		CurrentDate: timeline.StartDate,
		ParentID:    timeline.ID,
	})

	if err != nil {
		return nil, shared.InternalErrorException()
	}

	return result, nil
}

type UpdateTimelineDTO struct {
	Signature     uuid.UUID
	TimeInSeconds uint
}

// Branch timeline.
func (service *TimelineService) Update(t *UpdateTimelineDTO) (*TimelineEntity, *shared.Exception) {
	timeline, err := service.timelineRepository.FindBySignature(t.Signature)

	if err != nil {
		return nil, shared.DefaultException(shared.NOT_FOUND_EX, "timeline")
	}

	result, err := service.timelineRepository.Create(&TimelineEntity{
		ID:          uuid.New(),
		Signature:   timeline.Signature,
		GameID:      timeline.GameID,
		StartDate:   timeline.StartDate,
		CurrentDate: timeline.StartDate.Add(time.Second * time.Duration(t.TimeInSeconds)),
		ParentID:    timeline.ID,
	})

	if err != nil {
		return nil, shared.InternalErrorException()
	}

	return result, nil
}

type TimelineDTO struct {
	Signature uuid.UUID
}

// Find timeline by signature.
func (service *TimelineService) FindBySignature(t *TimelineDTO) (*TimelineEntity, *shared.Exception) {
	result, err := service.timelineRepository.FindBySignature(t.Signature)

	if err != nil {
		return nil, shared.DefaultException(shared.NOT_FOUND_EX, "timeline")
	}

	return result, nil
}

type TimelinesDTO struct {
	GameID uuid.UUID
}

// Find all timelines by gameID.
func (service *TimelineService) FindAllByGameID(t *TimelinesDTO) ([]*TimelineEntity, *shared.Exception) {
	result, err := service.timelineRepository.FindAllByGameID(t.GameID)

	if err != nil {
		return nil, shared.DefaultException(shared.NOT_FOUND_EX, "timeline")
	}

	return result, nil
}

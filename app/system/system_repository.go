package system

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type System struct {
	gorm.Model
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Key   string    `gorm:"unique"`
	Value string
	Type  string
}

type SystemRepository struct {
	db *gorm.DB
}

// Create a SystemRepository instance.
func NewSystemRepository(database *gorm.DB) *SystemRepository {
	return &SystemRepository{db: database}
}

// Create system variable.
func (repository *SystemRepository) Create(s *SystemEntity) (*SystemEntity, error) {
	envModel := &System{ID: s.ID, Key: s.Key, Value: s.Value, Type: s.Type}
	result := repository.db.Create(&envModel)

	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}

	return s, nil
}

// Get system variable by key.
func (repository *SystemRepository) GetByKey(key string) (*SystemEntity, error) {
	var model System

	result := repository.db.Where("key = ?", key).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &SystemEntity{Key: model.Key, Value: model.Value, Type: model.Type}, nil
}

// Get all system variables.
func (repository *SystemRepository) GetAll() ([]*SystemEntity, error) {
	var models []*System

	result := repository.db.Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	var entities []*SystemEntity
	for _, model := range models {
		entities = append(entities, &SystemEntity{
			Key:   model.Key,
			Value: model.Value,
			Type:  model.Type,
		})
	}

	return entities, nil
}

// Update system variable value by key
func (repository *SystemRepository) UpdateValueByKey(key string, value string) (*SystemEntity, error) {
	tx := repository.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var system System
	result := tx.Where("key = ?", key).First(&system)

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	system.Value = value
	err := tx.Save(&system).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &SystemEntity{
		Key:   system.Key,
		Value: system.Value,
		Type:  system.Type,
	}, nil
}

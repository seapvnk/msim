package user

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name     string    `gorm:"unique"`
	Password string
}

type UserRepository struct {
	db *gorm.DB
}

// Create an UserRepository instance
func (repository *UserRepository) New(database *gorm.DB) *UserRepository {
	return &UserRepository{db: database}
}

// Create an user in database
func (repository *UserRepository) Create(u *UserEntity) (*UserEntity, error) {
	userModel := &User{ID: u.ID, Name: u.Name, Password: u.password}
	result := repository.db.Create(&userModel)

	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}

	return u, nil
}

// Get all users in database
func (repository *UserRepository) GetAll() ([]*UserEntity, error) {
	var (
		userModels []User
		users      []*UserEntity
	)

	result := repository.db.Find(&userModels)

	if result.Error != nil {
		empty := []*UserEntity{}
		return empty, result.Error
	}

	for _, model := range userModels {
		users = append(users, &UserEntity{ID: model.ID, Name: model.Name})
	}

	return users, nil
}

// Get an user by id
func (repository *UserRepository) GetById(id uuid.UUID) (*UserEntity, error) {
	var model User

	result := repository.db.Where("id = ?", id).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &UserEntity{ID: model.ID, Name: model.Name}, nil
}

// Get an user by name
func (repository *UserRepository) GetByName(name string) (*UserEntity, error) {
	var model User

	result := repository.db.Where("name = ?", name).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &UserEntity{ID: model.ID, Name: model.Name, password: model.Password}, nil
}

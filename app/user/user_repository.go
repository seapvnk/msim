package user

import (
	"fmt"
	"errors"

	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"msim/db"
)

type User struct {
	gorm.Model
	Name	 string `gorm:"unique"`
	Password string
}

// Create an user in database
func create(u *UserEntity) (*UserEntity, error) {
	DB := db.ORM()

	cost := bcrypt.DefaultCost
	bytePassword := []byte(u.password)
	hashPasswordByte, err := bcrypt.GenerateFromPassword(bytePassword, cost)

	if err != nil {
		errorMessage := "Error generating hash for password"
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	hashPassword := string(hashPasswordByte)
	userModel := &User{Name: u.Name, Password: hashPassword}
	result := DB.Create(&userModel)

	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}

	return &UserEntity{ID: userModel.ID, Name: userModel.Name}, nil
}

// Get all users in database
func getAll() ([]*UserEntity, error) {
	DB := db.ORM()

	var (
		userModels	[]User
		users		[]*UserEntity
	)

	result := DB.Find(&userModels)

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
func getById(id uint) (*UserEntity, error) {
	var model User
	
	result := db.ORM().Where("id = ?", id).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &UserEntity{ID: model.ID, Name: model.Name}, nil
}

// Get an user by name
func getByName(name string) (*UserEntity, error) {
	var model User

	result := db.ORM().Where("name = ?", name).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &UserEntity{ID: model.ID, Name: model.Name}, nil
}
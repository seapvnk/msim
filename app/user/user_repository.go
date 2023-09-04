package user

import (
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

const DB = db.Get()

func create(u *UserEntity) (*UserEntity, error) {
	cost := bcrypt.DefaultCost
	bytePassword := []byte(password)
	hashPasswordByte, err := bcrypt.GenerateFromPassword(bytePassword, cost)

	if err != nil {
		errorMessage := "Error generating hash for password"
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	hashPassword := string(hashPasswordByte)
	userModel := DB.Create(&User{Name: name, password: hashPassword})

	if userModel.Error != nil {
		fmt.Println(userModel.Error)
		return nil, userModel.Error
	}

	return &User{ID: userModel.ID, Name: name}, nil
}

func getAll(db *gorm.DB) (*UserEntity[], error) {
	var (
		userModels	[]User
		users		[]UserEntity
	)

	result := DB.Find(&userModels)

	if result.Error != nil {
		return [], result.Error
	}

	for _, model := range result {
		users = append(users, &UserEntity{ID: mode.ID, Name: model.name})
	}

	return users, nil
}

func getById(id uint) (*UserEntity, error) {
	result := DB.First(&UserModel{ID: id})
	if result.Error != nil {
		return nil, result.Error
	}

	user := &UserEntity{ID: result.ID, Name: result.Name}
	return user, nil
}

func getByName(name string) (*UserEntity, error) {
	result := DB.Where("name = ?", name).First(&User{})

	if result.Error != nil {
		return nil, result.Error
	}

	user := &UserEntity{ID: result.ID, Name: result.Name}
	return user, nil
}
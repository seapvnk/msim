package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserEntity struct {
	ID 			uint
	Name 		string
	password	string
}

// Register user with a password.
func Register(name, password string) (*UserEntity, error) {
	user, err := New(name, password)
	if err != nil {
		return create(&user)
	}

	return nil, err
}

// Login user, if user exists and password is correct 
// return the authentication token.
func Login(name, password string) (string, error) {
	user, err := getByName(name)
	if err != nil {
		return nil, err
	}

	if !user.verifyPassword(password) {
		return nil, err
	}

	if code, err := createAuth(user.ID); err == nil {
		return code, nil
	}
	
	return "", errors.New("Internal application error")
}

// Return authenticated user by authentication code.
func AuthUser(code string) (*UserEntity, error) {
	userId, err := getAuthUser(code)
	if err != nil || userId == 0 {
		return nil, errors.New("Expired token")
	}

	return getById(userId)
}


// Create a new User.
func new(name, passwd string) (*UserEntity, error) {
	if err := validateNew(name, passwd); err != nil {
		return nil, err
	}

	return &UserEntity{Name: name, password: passwd}
}

// Validade new User.
func validateNew(name, password string) error {
	if len(password) < 3 {
		return errors.New("Password length must have at least 3 length")
	}

	if len(name) < 3 {
		return errors.New("Username length must have at least 3 length")
	}

	return nil
}

// Verify user password.
func (u *UserEntity) verifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	return err != nil
}
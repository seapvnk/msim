package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserEntity struct {
	ID 			uint
	Name 		string
	password	string
}

// Register user with a password.
func Register(name, password string) (*UserEntity, error) {
	user, err := new(name, password)
	if err != nil {
		return nil, err
	}

	return create(user)
}

// Login user, if user exists and password is correct 
// return the authentication token.
func Login(name, password string) (string, error) {
	user, err := getByName(name)

	if err != nil {
		return "", err
	}

	if !user.verifyPassword(password) {
		return "", errors.New("Incorrect password")
	}

	code, err := createAuth(user.ID)
	
	if err != nil {
		return "", errors.New("Internal application error")
	}
	
	return code, nil
}

// Return authenticated user by authentication code.
func AuthUser(code string) (*UserEntity, error) {
	user, err := getAuthUser(code)
	if err != nil || user.ID == 0 {
		return nil, errors.New("Expired token or doesnt exists")
	}

	return user, nil
}

// Create a new User.
func new(name, passwd string) (*UserEntity, error) {
	if err := validateNew(name, passwd); err != nil {
		return nil, err
	}

	cost := bcrypt.DefaultCost
	bytePassword := []byte(passwd)
	hashPasswordByte, err := bcrypt.GenerateFromPassword(bytePassword, cost)

	if err != nil {
		errorMessage := "Error generating hash for password"
		return nil, errors.New(errorMessage)
	}

	return &UserEntity{Name: name, password: string(hashPasswordByte)}, nil
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
	return err == nil
}
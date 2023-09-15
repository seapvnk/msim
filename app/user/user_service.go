package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"msim/app/shared"
)

type UserEntity struct {
	ID       uuid.UUID
	Name     string
	password string
}

// Validade new User.
func (u *UserEntity) validate() *shared.Exception {
	if len(u.password) < 3 {
		return shared.FormException(shared.MIN_LENGTH_EX, "password")
	}

	if len(u.Name) < 3 {
		return shared.FormException(shared.MIN_LENGTH_EX, "name")
	}

	return nil
}

// Verify user password.
func (u *UserEntity) verifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	return err == nil
}

type UserService struct {
	userRepository *UserRepository
	authRepository *AuthRepository
}

type UserAuthDTO struct {
	Name     string
	Password string
}

// Register user with a password.
func (service *UserService) Register(u *UserAuthDTO) (*UserEntity, *shared.Exception) {
	user, ex := new(u.Name, u.Password)
	if ex != nil {
		return nil, ex
	}

	result, err := service.userRepository.Create(user)
	if err != nil {
		return nil, shared.DefaultException(shared.ALREADY_CREATED_EX, "user")
	}

	return result, nil
}

// Login user, if user exists and password is correct
// return the authentication token.
func (service *UserService) Login(u *UserAuthDTO) (uuid.UUID, *shared.Exception) {
	user, err := service.userRepository.GetByName(u.Name)

	if err != nil {
		return uuid.Nil, shared.FormException(shared.NOT_FOUND_EX, "user")
	}

	if !user.verifyPassword(u.Password) {
		return uuid.Nil, shared.FormException(shared.UNAUTHORIZED_EX, "password")
	}

	code, err := service.authRepository.Create(user.ID)
	if err != nil {
		return uuid.Nil, shared.InternalErrorException()
	}

	return code, nil
}

type AuthDTO struct {
	Code uuid.UUID
}

// Return authenticated user by authentication code.
func (service *UserService) GetAuthUser(auth *AuthDTO) (*UserEntity, *shared.Exception) {
	user, err := service.authRepository.GetAuthUser(auth.Code)
	if err != nil || user.ID == uuid.Nil {
		msg := "Expired token or user doesnt exists"
		return nil, shared.DefaultException(shared.UNAUTHORIZED_EX, msg)
	}

	return user, nil
}

// PRIVATE:

// Create a new User.
func new(name, passwd string) (*UserEntity, *shared.Exception) {
	user := UserEntity{ID: uuid.New(), Name: name, password: passwd}
	if ex := user.validate(); ex != nil {
		return nil, ex
	}

	hashPassword, ex := getPasswordHash(passwd)
	if ex != nil {
		return nil, ex
	}

	user.password = hashPassword

	return &user, nil
}

// Get password hash
func getPasswordHash(password string) (string, *shared.Exception) {
	cost := bcrypt.DefaultCost

	if hashByte, err := bcrypt.GenerateFromPassword([]byte(password), cost); err == nil {
		return string(hashByte), nil
	}

	return "", shared.DefaultException(shared.UNKNOWN_EX, "Can't hash password")
}

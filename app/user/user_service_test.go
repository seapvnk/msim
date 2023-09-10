package user

import (
	"testing"
	"reflect"

	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"msim/db"
)

// Test register
func TestRegister(t *testing.T) {
	t.Run("Should register an user", func(t *testing.T) {
		service, DB := CreateUserService()
		result, err := service.Register(&UserAuthDTO{"Test", "passwd"})

		if err != nil {
			t.Fatal(err)
		}

		var find User
		DB.First(&find)

		if find.Name != result.Name {
			t.Fatal("Should create this exact model in database")
		}
	})
}

// Test login
func TestLogin(t *testing.T) {
	t.Run("Should login an user when exists and password match", func(t *testing.T) {
		service, DB := CreateUserService()

		cost := bcrypt.DefaultCost
		bytePassword := []byte("passwd")
		hashPasswordByte, _ := bcrypt.GenerateFromPassword(bytePassword, cost)
		
		userModel := &User{ID: uuid.New(), Name: "Test99", Password: string(hashPasswordByte)}
		DB.Create(userModel)

		result, ex := service.Login(&UserAuthDTO{"Test99", "passwd"})

		if ex != nil {
			t.Fatal(ex)
		}

		if result == uuid.Nil {
			t.Fatal("Should return a code")
		}
	})

	t.Run("Should not login an user when password doesnt match", func(t *testing.T) {
		service, DB := CreateUserService()

		cost := bcrypt.DefaultCost
		bytePassword := []byte("passwd")
		hashPasswordByte, _ := bcrypt.GenerateFromPassword(bytePassword, cost)
		
		userModel := &User{ID: uuid.New(), Name: "Test2", Password: string(hashPasswordByte)}
		DB.Create(userModel)

		result, ex := service.Login(&UserAuthDTO{"Test", "passwd1"})

		if ex == nil {
			t.Fatal(ex)
		}

		if result != uuid.Nil {
			t.Fatal("Should return a code")
		}
	})

	t.Run("Should not login an user when doesnt exists", func(t *testing.T) {
		service, _ := CreateUserService()
		result, err := service.Login(&UserAuthDTO{"Test", "passwd"})

		if err == nil {
			t.Fatal(err)
		}

		if result != uuid.Nil {
			t.Fatal("Should return a code")
		}
	})
}

// Test GetAuthUser
func TestGetAuthUserService(t *testing.T) {
	t.Run("Should get auth user when theres one", func(t *testing.T) {
		service, DB := CreateUserService()

		createdUser := User{ID: uuid.New(), Name: "test15", Password: "12345"}
		DB.Create(&createdUser)
		
		createdAuth := Auth{Code: uuid.New(), User: createdUser}
		DB.Create(&createdAuth)

		result, _ := service.GetAuthUser(&AuthDTO{createdAuth.Code})

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*UserEntity)(nil))
	
		if resultType != expectedType {
			errFormated := `create returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}
	})

	t.Run("Should not get auth user when doesnt have one", func(t *testing.T) {
		service, _ := CreateUserService()
		result, err := service.GetAuthUser(&AuthDTO{uuid.New()})

		if err == nil {
			t.Fatal("Should throw error")
		}

		if result != nil {
			t.Fatal("User should be nil")
		}
	})
}

func CreateUserService() (*UserService, *gorm.DB) {
	DB, _ := db.InMemoryDB()
	Migrate(DB)

	userRepo, authRepo := &UserRepository{db: DB}, &AuthRepository{db: DB}
	return &UserService{userRepository: userRepo, authRepository: authRepo}, DB
}
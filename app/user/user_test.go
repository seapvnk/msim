package user

import (
	"os"
	"testing"
	"fmt"
	"reflect"

	"golang.org/x/crypto/bcrypt"
	"msim/db"
)

// Test register
func TestRegister(t *testing.T) {
	t.Run("Should register an user", func(t *testing.T) {
		setup()
	
		DB := db.ORM()

		result, err := Register("Test", "passwd")

		if err != nil {
			t.Fatal(err)
		}

		var find User
		DB.First(&find)

		if find.Name != result.Name {
			t.Fatal("Should create this exact model in database")
		}
	
		teardown()
	})
}

// Test login
func TestLogin(t *testing.T) {
	t.Run("Should login an user when exists and password match", func(t *testing.T) {
		setup()
	
		DB := db.ORM()

		cost := bcrypt.DefaultCost
		bytePassword := []byte("passwd")
		hashPasswordByte, err := bcrypt.GenerateFromPassword(bytePassword, cost)
		
		userModel := &User{Name: "Test", Password: string(hashPasswordByte)}
		DB.Create(userModel)

		result, err := Login("Test", "passwd")

		if err != nil {
			t.Fatal(err)
		}

		if result == "" {
			t.Fatal("Should return a code")
		}
	
		teardown()
	})

	t.Run("Should not login an user when password doesnt match", func(t *testing.T) {
		setup()
	
		DB := db.ORM()

		cost := bcrypt.DefaultCost
		bytePassword := []byte("passwd")
		hashPasswordByte, err := bcrypt.GenerateFromPassword(bytePassword, cost)
		
		userModel := &User{Name: "Test", Password: string(hashPasswordByte)}
		DB.Create(userModel)

		result, err := Login("Test", "passwd1")

		if err == nil {
			t.Fatal(err)
		}

		if result != "" {
			t.Fatal("Should return a code")
		}
	
		teardown()
	})

	t.Run("Should not login an user when doesnt exists", func(t *testing.T) {
		setup()
	
		result, err := Login("Test", "passwd")

		if err == nil {
			t.Fatal(err)
		}

		if result != "" {
			t.Fatal("Should return a code")
		}
	
		teardown()
	})
}

// Test AuthUser
func TestAuthUser(t *testing.T) {
	t.Run("Should get auth user when theres one", func(t *testing.T) {
		DB := db.ORM()

		createdUser := User{Name: "test1", Password: "12345"}
		DB.Create(&createdUser)
		
		createdAuth := Auth{Code: "ABCD", User: createdUser}
		DB.Create(&createdAuth)

		result, _ := AuthUser(createdAuth.Code)

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*UserEntity)(nil))
	
		if resultType != expectedType {
			errFormated := `create returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}
	})

	t.Run("Should not get auth user when doesnt have one", func(t *testing.T) {
		result, err := AuthUser("ABCDE")

		if err == nil {
			t.Fatal("Should throw error")
		}

		if result != nil {
			t.Fatal("User should be nil")
		}
	})
}

// Runs after each tests
func teardown() {
	folder := "storage"

	err := os.RemoveAll(folder)
	if err != nil {
		panic("Error excluding database folder")
	} else {
		fmt.Println("Folder excluded")
	}
}

// Setup environment for test
func setup() {
	if err := os.Mkdir("storage", os.ModePerm); err != nil {
		panic("Error creating storage folder")
	}
	fmt.Println("Storage folder created")

	newFile, err := os.Create("storage/db.sqlite")
	if err != nil {
		panic("Error creating database file")
	}
	defer newFile.Close()

	Migrate()
}
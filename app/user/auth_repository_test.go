package user

import (
	"fmt"
	"testing"
	"os"
	"reflect"
	"time"

	"msim/db"
)

// Test create
func TestCreateAuth(t *testing.T) {
	t.Run("Should create an auth and return a code when theres an user", func(t *testing.T) {
		setupAuthRepo()
	
		DB := db.ORM()

		created := &User{Name: "test1", Password: "12345"}
		DB.Create(&created)

		result, err := createAuth(created.ID)

		var find Auth
		DB.First(&find)

		if err != nil {
			t.Fatal(err)
		}

		if find.Code != result {
			t.Fatal("Created code should match returned code")
		}
	
		teardownAuthRepo()
	})

	t.Run("Should not create an auth when theres no user", func(t *testing.T) {
		setupAuthRepo()
	
		_, err := createAuth(1)

		if err == nil {
			t.Fatal("Should return an error")
		}

		teardownAuthRepo()
	})
}

// Test getAuthUser
func TestGetAuthUser(t *testing.T) {
	t.Run("Should get auth user when theres an user", func(t *testing.T) {
		setupAuthRepo()
	
		DB := db.ORM()

		createdUser := User{Name: "test1", Password: "12345"}
		DB.Create(&createdUser)
		
		createdAuth := Auth{Code: "ABCD", User: createdUser}
		DB.Create(&createdAuth)

		result, err := getAuthUser(createdAuth.Code)

		if err != nil {
			t.Fatal(err)
		}

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*UserEntity)(nil))
	
		if resultType != expectedType {
			errFormated := `find returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		if result.Name != createdUser.Name {
			t.Fatal("User name should match should match")
		}

		teardownAuthRepo()
	})

	t.Run("Should not get auth user when theres no auth code or user", func(t *testing.T) {
		setupAuthRepo()

		result, err := getAuthUser("ABCD")

		if err == nil {
			t.Fatal(err)
		}

		if result != nil {
			t.Fatal("Should not find am user")
		}

		teardownAuthRepo()
	})

	t.Run("Should not get auth user when auth code is expired", func(t *testing.T) {
		setupAuthRepo()
		
		DB := db.ORM()

		createdUser := User{Name: "test1", Password: "12345"}
		DB.Create(&createdUser)

		mockTimeCreated := time.Date(2002, time.January, 9, 0, 0, 0, 0, time.UTC)		
		createdAuth := Auth{Code: "ABCD", User: createdUser}
		createdAuth.CreatedAt = mockTimeCreated

		DB.Create(&createdAuth)

		result, err := getAuthUser(createdAuth.Code)

		if err == nil {
			t.Fatal(err)
		}

		if result != nil {
			t.Fatal("Should not find am user")
		}

		teardownAuthRepo()
	})
}

// Test generateRandomString
func TestGenerateRandomString(t *testing.T) {
	t.Run("Should generate a random string with length 10", func(t *testing.T) {
		length := uint(10)
		result := generateRandomString(length)
		lenResult := uint(len(result))

		if lenResult != length {
			t.Fatalf(`Expected length %d, received length %d`, length, lenResult)
		}
	})

	t.Run("Should generate a random string with length 0", func(t *testing.T) {
		length := uint(10)
		result := generateRandomString(length)
		lenResult := uint(len(result))

		if lenResult != length {
			t.Fatalf(`Expected length %d, received length %d`, length, lenResult)
		}
	})

	t.Run("Each random string should be different", func(t *testing.T) {
		result1 := generateRandomString(10)
		result2 := generateRandomString(10)

		if result1 == result2 {
			t.Fatal("Random strings cant be the same")
		}
	})
}

// Runs after each tests
func teardownAuthRepo() {
	folder := "storage"

	err := os.RemoveAll(folder)
	if err != nil {
		panic("Error excluding database folder")
	} else {
		fmt.Println("Folder excluded")
	}
}

// Setup environment for test
func setupAuthRepo() {
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
package user

import (
	"fmt"
	"testing"
	"os"
	"reflect"

	"msim/db"
)

// Test create
func TestCreate(t *testing.T) {
	t.Run("Should create an user and return an user entity", func(t *testing.T) {
		setupUserRepo()
		DB := db.ORM()
	
		result, _ := create(&UserEntity{Name: "test", password: "12345"})

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*UserEntity)(nil))
	
		if resultType != expectedType {
			errFormated := `create returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		var countUsers int64
		DB.Model(&User{}).Count(&countUsers)

		if countUsers != 1 {
			errFormated := `created count returns %d, expects 1`
			t.Fatalf(errFormated, countUsers)
		}
	
		teardownUserRepo()
	})
}

// Test getAll
func TestGetAll(t *testing.T) {
	t.Run("Should get all users when theres users", func(t *testing.T) {
		setupUserRepo()
		DB := db.ORM()
	
		DB.Create(&User{Name: "test1", Password: "12345"})
		DB.Create(&User{Name: "test2", Password: "12345"})

		result, err := getAll()

		if err != nil {
			t.Fatal(err)
		}

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf(([]*UserEntity)(nil))
	
		if resultType != expectedType {
			errFormated := `find returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		if len(result) != 2 {
			t.Fatal("Should return the only 2 created users")
		}

		if result[0].Name != "test1" && result[1].Name != "test2" {
			t.Fatal("Users should match")
		}
		
		teardownUserRepo()
	})

	t.Run("Should get an empty array when theres no user", func(t *testing.T) {
		setupUserRepo()

		result, err := getAll()

		if err != nil {
			t.Fatal(err)
		}

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf(([]*UserEntity)(nil))
	
		if resultType != expectedType {
			errFormated := `find returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		if len(result) != 0 {
			t.Fatal("Should return zero users")
		}
		
		teardownUserRepo()
	})
}

// Test getById
func TestGetById(t *testing.T) {
	t.Run("Should get an user by id when exists", func(t *testing.T) {
		setupUserRepo()
		DB := db.ORM()
	
		created := User{Name: "test1", Password: "12345"} 
		DB.Create(&created)

		result, err := getById(created.ID)

		if err != nil {
			t.Fatal(err)
		}

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*UserEntity)(nil))
	
		if resultType != expectedType {
			errFormated := `find returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		if result.Name != created.Name {
			t.Fatal("User name should match should match")
		}
		
		teardownUserRepo()
	})

	t.Run("Should not get an user by id when doesnt exists", func(t *testing.T) {
		setupUserRepo()
		result, err := getById(1)

		if err == nil {
			t.Fatal("Should throw error")
		}

		if result != nil {
			t.Fatal("User should be nil")
		}
		
		teardownUserRepo()
	})
}

// Test getByName
func TestGetByName(t *testing.T) {
	t.Run("Should get an user by name when exists", func(t *testing.T) {
		setupUserRepo()
		DB := db.ORM()
	
		created := User{Name: "test1", Password: "12345"} 
		DB.Create(&created)

		result, err := getByName(created.Name)
		fmt.Println(result)

		if err != nil {
			t.Fatal(err)
		}

		resultType := reflect.TypeOf(result)
		expectedType := reflect.TypeOf((*UserEntity)(nil))
	
		if resultType != expectedType {
			errFormated := `find returns type of %s, expects type of %s`
			t.Fatalf(errFormated, resultType, expectedType)
		}

		if result.Name != created.Name {
			t.Fatal("User name should match should match")
		}
		
		teardownUserRepo()
	})

	t.Run("Should not get an user by name when doesnt exists", func(t *testing.T) {
		setupUserRepo()
		result, err := getByName("Test2")

		if err == nil {
			t.Fatal("Should throw error")
		}

		if result != nil {
			t.Fatal("User should be nil")
		}
		
		teardownUserRepo()
	})
}

// Runs after each tests
func teardownUserRepo() {
	folder := "storage"

	err := os.RemoveAll(folder)
	if err != nil {
		panic("Error excluding database folder")
	} else {
		fmt.Println("Folder excluded")
	}
}

// Setup environment for test
func setupUserRepo() {
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
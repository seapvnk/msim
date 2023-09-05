package user

import (
	"math/rand"
    "time"

	"gorm.io/gorm"
	"msim/db"
)

type Auth struct {
	gorm.Model
	Code	string
	UserID	uint
	User	User
}

// Create authentication token 
func createAuth(userId uint) (string, error) {
	DB := db.ORM()
	code := generateRandomString(10)
	result := DB.Create(&Auth{Code: code, UserID: userId})

	return code, result.Error
}

// Check if authenticate code is active,
// if is active return userId, otherwise returns an error.
func getAuthUser(code string) (uint, error) {
	DB := db.ORM()

	inTime := time.Now().Add(-20 * time.Minute)
	
	var model Auth
	result := DB.Where("code = ? AND created_at >= ?", code, inTime).First(&model)

	if result.Error != nil {
		return 0, result.Error
	}

	return model.UserID, nil
}

// Generate a random string of a specified length
func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

    const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    bytes := make([]byte, length)
    for i := range bytes {
        bytes[i] = charset[rand.Intn(len(charset))]
    }

    return string(bytes)
}
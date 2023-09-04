package user

import (
	"errors"
	"math/rand"
    "time"

	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"msim/db"
)

type Auth struct {
	gorm.Model
	Code	string
	UserID	uint
	User	User
}

const DB = db.Get()

// Create authentication token 
func createAuth(userId uint) (string, error) {
	code := generateRandomString(10)
	_, err := DB.Create(&AuthModel{Code: code, UserID: userId})

	return code, err
}

// Check if authenticate code is active,
// if is active return userId, otherwise returns an error.
func getAuthUser(code string) (uint, error) {
	inTime := time.Now().Add(-20 * time.Minute)
	result := DB.Where("code = ? AND created_at >= ?", code, inTime).First(&Auth{})

	if result.Error != nil {
		return 0, result.Error
	}

	return result.UserID, nil
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
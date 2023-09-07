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
func getAuthUser(code string) (*UserEntity, error) {
	DB := db.ORM()
	
	query := `
		SELECT user.* FROM users as user 
		LEFT JOIN auths as auth
			ON user.id = auth.user_id
		WHERE auth.code = ? AND auth.created_at >= ?
		ORDER BY auth.created_at DESC
		LIMIT 1
	`
		
	var model User
	inTime := time.Now().Add(-20 * time.Minute)
	result := DB.Raw(query, code, inTime).First(&model)

	if result.Error != nil {
		return nil, result.Error
	}

	return &UserEntity{ID: model.ID, Name: model.Name}, nil
}

// Generate a random string of a specified length
func generateRandomString(length uint) string {
	rand.Seed(time.Now().UnixNano())

    const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    bytes := make([]byte, length)
    for i := range bytes {
        bytes[i] = charset[rand.Intn(len(charset))]
    }

    return string(bytes)
}
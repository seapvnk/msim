package user

import (
	"time"
	"errors"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Auth struct {
	gorm.Model
	ID      uuid.UUID	`gorm:"type:uuid;primaryKey"`
	Code	uuid.UUID
	UserID	uuid.UUID
	User	User
}

type AuthRepository struct {
	db	*gorm.DB
}

// Create an AuthRepository instance.
func (repository *AuthRepository) New(database *gorm.DB) *AuthRepository {
	return &AuthRepository{db: database}
}

// Create authentication token.
func (repository *AuthRepository) Create(userId uuid.UUID) (uuid.UUID, error) {
	code := uuid.New()
	result := repository.db.Create(&Auth{ID: uuid.New(), Code: code, UserID: userId})

	return code, result.Error
}

// Check if authenticate code is active,
// if is active return userId, otherwise returns an error.
func (repository *AuthRepository) GetAuthUser(code uuid.UUID) (*UserEntity, error) {
	var models []User
	
	inTime := time.Now().Add(-20 * time.Minute)
	result := repository.db.Raw(`
		SELECT user.* FROM users as user 
		LEFT JOIN auths as auth
			ON user.id = auth.user_id
		WHERE auth.code = ? AND auth.created_at >= ?
		ORDER BY auth.created_at DESC
		LIMIT 1
	`, code, inTime).Find(&models)

	if len(models) == 0 || result.Error != nil {
		return nil, errors.New("Record not found")
	}

	model := models[0]
	return &UserEntity{ID: model.ID, Name: model.Name}, nil
}
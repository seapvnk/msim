package setup

import (
	"msim/app/user"
)

// Migrate tables
func Migrate() {
	user.Migrate()
}
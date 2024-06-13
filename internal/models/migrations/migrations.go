package migrations

import "github.com/SineChat/auth-ms/internal/models"

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		&models.User{},
		&models.BannedAccount{},
		&models.LoginToken{},
		&models.AccessToken{},
		&models.Subscription{},
		&models.OTP{},
		&models.Payment{},
		&models.WebhookLog{},

		// Admin models
		&models.FlaggedAccount{},
	}
}

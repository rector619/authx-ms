package migrations

import (
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
)

func RunAllMigrations(db *mongodb.Database) {

	// migration
	MigrateModels(db, AuthMigrationModels())

}

func MigrateModels(db *mongodb.Database, models []interface{}) {
	_ = db.AutoMigrate(models)
}

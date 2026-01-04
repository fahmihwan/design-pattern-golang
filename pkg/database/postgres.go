package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"best-pattern/internal/config"
	"best-pattern/migrations"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	gormConfig := &gorm.Config{
		// TranslateError: true,
		Logger: newLogger, // Set the desired log level
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// // Perform auto-migration
	err = autoMigrate(db)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	return db, nil

}

func autoMigrate(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not apply migrations: %v", err)
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}

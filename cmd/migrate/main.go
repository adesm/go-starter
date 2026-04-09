package main

import (
	"log/slog"
	"os"

	"boilerplate/internal/config"
	"boilerplate/internal/module/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	dsn := "host=" + cfg.Database.Host + " user=" + cfg.Database.User + " password=" + cfg.Database.Password + " dbname=" + cfg.Database.Name + " port=" + cfg.Database.Port + " sslmode=" + cfg.Database.SSLMode + " TimeZone=UTC"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database for migration", "error", err)
		os.Exit(1)
	}

	slog.Info("Running migrations...")

	// List all models that need to be migrated here
	err = db.AutoMigrate(&user.User{})
	if err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	slog.Info("Migrations completed successfully!")
}

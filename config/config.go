package config

import (
	"log"
	"os"
	"time"
)

type Configuration struct {
	DatabaseName        string
	DatabaseHost        string
	DatabaseUser        string
	DatabasePassword    string
	MigrateToVersion    string
	MigrationLocation   string
	FileStorageLocation string
	JwtSecret           string
	JwtTTL              time.Duration
}

func GetConfiguration() Configuration {
	migrationLocation, set := os.LookupEnv("MIGRATION_LOCATION")
	if !set {
		migrationLocation = "internal/infra/database/migrations"
	}
	migrateToVersion, set := os.LookupEnv("MIGRATE")
	if !set {
		migrateToVersion = "latest"
	}
	staticFilesLocation, set := os.LookupEnv("FILES_LOCATION")
	if !set {
		staticFilesLocation = "file_storage"
	}
	jwtSecret, set := os.LookupEnv("JWT_SECRET")
	if !set || jwtSecret == "" {
		log.Fatal("JWT_SECRET env vat is missing")
	}
	return Configuration{
		DatabaseName: os.Getenv("DB_CHAT"),
		DatabaseHost: os.Getenv("DB_HOST"),
		//DatabaseHost:        "localhost:54322",
		DatabaseUser:        os.Getenv("DB_USER"),
		DatabasePassword:    os.Getenv("DB_PASSWORD"),
		MigrateToVersion:    migrateToVersion,
		MigrationLocation:   migrationLocation,
		FileStorageLocation: staticFilesLocation,
		JwtSecret:           jwtSecret,
		JwtTTL:              72 * time.Hour,
	}
}

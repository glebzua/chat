package config

import (
	"github.com/pusher/pusher-http-go/v5"
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
	Pusher              pusher.Client
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
	//pusher conf
	appID, set := os.LookupEnv("chatprjkt_pusher_appID")
	if !set || appID == "" {
		log.Fatal("appID env is missing")
	}
	key, set := os.LookupEnv("chatprjkt_pusher_key")
	if !set || key == "" {
		log.Fatal("key env is missing")
	}
	secret, set := os.LookupEnv("chatprjkt_pusher_secret")
	if !set || secret == "" {
		log.Fatal("secret env is missing")
	}
	cluster, set := os.LookupEnv("chatprjkt_pusher_cluster")
	if !set || cluster == "" {
		log.Fatal("cluster env is missing")
	}
	secure := true
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
		Pusher: pusher.Client{
			AppID:   appID,
			Key:     key,
			Secret:  secret,
			Cluster: cluster,
			Secure:  secure,
		},
	}
}

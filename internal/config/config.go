package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBHost           string `mapstructure:"PG_DB_HOST"`
	DBUser           string `mapstructure:"PG_DB_USER"`
	DBPassword       string `mapstructure:"PG_DB_PASSWORD"`
	DBName           string `mapstructure:"DB_PG_RAIJAI_DB_NAME"`
	DBPort           string `mapstructure:"PG_DB_PORT"`
	DBSSLMode        string `mapstructure:"PG_SSL_MODE"`
	KeycloakIssuer   string `mapstructure:"KEYCLOAK_ISSUER"`
	KeycloakClientID string `mapstructure:"KEYCLOAK_RAIJAI_CLIENT_ID"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
	}

	return &Config{
		DBHost:           viper.GetString("PG_DB_HOST"),
		DBUser:           viper.GetString("PG_DB_USER"),
		DBPassword:       viper.GetString("PG_DB_PASSWORD"),
		DBName:           viper.GetString("DB_PG_RAIJAI_DB_NAME"),
		DBPort:           viper.GetString("PG_DB_PORT"),
		DBSSLMode:        viper.GetString("PG_SSL_MODE"),
		KeycloakIssuer:   viper.GetString("KEYCLOAK_ISSUER"),
		KeycloakClientID: viper.GetString("KEYCLOAK_RAIJAI_CLIENT_ID"),
	}
}

func ConnectDB(cfg *Config) *gorm.DB {
	sslMode := cfg.DBSSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, sslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
	return db
}

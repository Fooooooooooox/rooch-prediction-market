package testutils

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("..")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file: %s", err)
	}
}

func AssertNil(err error) {
	if err != nil {
		panic(err)
	}
}

func GetTestDSN() string {
	// Read PostgreSQL configuration from environment variables
	host := viper.GetString("CONFIG_ENV_POSTGRES_HOST")
	port := viper.GetInt("CONFIG_ENV_POSTGRES_PORT")
	user := viper.GetString("CONFIG_ENV_POSTGRES_USER")
	password := viper.GetString("CONFIG_ENV_POSTGRES_PASSWORD")
	database := viper.GetString("CONFIG_ENV_POSTGRES_DBNAME")
	sslmode := viper.GetString("CONFIG_ENV_POSTGRES_SSLMODE")
	url := viper.GetString("CONFIG_ENV_PG_URL")

	if url != "" {
		return url
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, database, sslmode)
}

func GenerateTestJWT(userId uint, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"userid": userId,
		"exp":    time.Now().Add(time.Hour * 1).Unix(), // Token valid for 1 hour
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func UintPtr(i uint) *uint {
	return &i
}

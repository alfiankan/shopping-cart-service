package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// ApplicationConfig contain all application config loaded from env or (dot)env file
type ApplicationConfig struct {
	PostgreeHost string
	PostgreeUser string
	PostgreePass string
	PostgreeDb   string
	PostgreePort int
	PostgreeSsl  string
	RedisHost    []string
	RedisPass    string
	HTTPApiPort  string
}

// Load load config from env
func Load(configFile ...string) ApplicationConfig {

	if err := godotenv.Load(configFile...); err != nil {
		//load from os env
		log.Println(err.Error(), "trying to load config from os env instead")
	}

	postgreeDbPort, err := strconv.Atoi(os.Getenv("PG_DATABASE_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	return ApplicationConfig{
		PostgreeHost: os.Getenv("PG_DATABASE_HOST"),
		PostgreeUser: os.Getenv("PG_DATABASE_USERNAME"),
		PostgreePass: os.Getenv("PG_DATABASE_PASSWORD"),
		PostgreeDb:   os.Getenv("PG_DATABASE_NAME"),
		PostgreePort: postgreeDbPort,
		PostgreeSsl:  os.Getenv("PG_DATABASE_SSL_MODE"),
		RedisHost:    strings.Split(os.Getenv("REDIS_HOST"), ";"),
		RedisPass:    os.Getenv("REDIS_PASSWORD"),
		HTTPApiPort:  os.Getenv("HTTP_API_PORT"),
	}

}

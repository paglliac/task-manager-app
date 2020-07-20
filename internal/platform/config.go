package platform

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type DbConfig struct {
	Host string
	Port string

	Db string

	Username string
	Password string
}

type Config struct {
	DbConfig DbConfig
}

func InitConfig() Config {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	err := godotenv.Load(basepath + "/../../.env")

	if err != nil {
		log.Fatal("Config file can not be loaded", err)
	}

	c := Config{
		DbConfig: DbConfig{
			Host:     getEnv("DB_HOST"),
			Port:     getEnv("DB_PORT"),
			Db:       getEnv("DB_NAME"),
			Username: getEnv("DB_USER"),
			Password: getEnv("DB_PASSWORD"),
		},
	}

	return c
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Printf("Env variable %s does not represented", key)
	panic("Env variable does not represented")
}

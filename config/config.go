package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const (
	TimeExpiredAt = time.Hour * 24
)

type Config struct {
	Environment string

	ServerHost string
	ServerPort string

	Redis Redis

	Postgres Postgres

	Minio Minio

	SekretKey string

	TelegramToken string
}

type Redis struct {
	Host     string
	Port     int
	Password string
	DataBase int
}

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	DataBase string
}

type Minio struct {
	Host            string
	Port            int
	AccessKeyID     string
	SecretAccessKey string
	Secure          bool
}

func Load() *Config {
	if err := godotenv.Load("config/.env"); err != nil {
		fmt.Println("NO .env file  not found")
	}

	cfg := Config{}
	cfg.ServerHost = cast.ToString(getOrDefaultValue("SERVER_HOST", "3.123.128.20"))
	cfg.Postgres = Postgres{
		Host:     cast.ToString(getOrDefaultValue("POSTGRES_HOST", "3.123.128.20")),
		Port:     cast.ToInt(getOrDefaultValue("POSTGRES_PORT", "5432")),
		User:     cast.ToString(getOrDefaultValue("POSTGRES_USER", "maxmurjon")),
		Password: cast.ToString(getOrDefaultValue("POSTGRES_PASSWORD", "max22012004")),
		DataBase: cast.ToString(getOrDefaultValue("POSTGRES_DATABASE", "comics"))}
	cfg.Redis = Redis{
		Host:     cast.ToString(getOrDefaultValue("REDIS_HOST", "localhost")),
		Port:     cast.ToInt(getOrDefaultValue("REDIS_PORT", "6379")),
		Password: cast.ToString(getOrDefaultValue("REDIS_PASSWORD", "")),
		DataBase: cast.ToInt(getOrDefaultValue("REDIS_DATABASE", "0")),}

	cfg.Minio = Minio{
		Host:            cast.ToString(getOrDefaultValue("MINIO_HOST", "3.123.128.20")),
		Port:            cast.ToInt(getOrDefaultValue("MINIO_PORT", "9000")),
		AccessKeyID:     cast.ToString(getOrDefaultValue("ACCESSKEY", "maxmurjon")),
		SecretAccessKey: cast.ToString(getOrDefaultValue("SECRETKEY", "max22012004")),
		Secure:          cast.ToBool(getOrDefaultValue("SECURE", "false")),
	}

	cfg.SekretKey = cast.ToString(getOrDefaultValue("SEKRET_KEY", "sekret"))
	cfg.TelegramToken = cast.ToString(getOrDefaultValue("TELEGRAM_BOT_TOKEN", "token"))

	fmt.Println(cfg)
	return &cfg
}

func getOrDefaultValue(key string, defaultValue string) interface{} {
	val, exists := os.LookupEnv(key)
	if exists {
		return val
	}

	return defaultValue
}

package main

import (
	"fmt"
	"log"
	"morc/bot"
	"morc/config"
	postgres "morc/storage/postges"
)

func main() {
	cfg := config.Load()
	psqlConnString := fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.User,
		cfg.Postgres.DataBase,
		cfg.Postgres.Password,
		cfg.Postgres.Port,
	)
	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)

	strg := postgres.NewPostgres(psqlConnString, redisAddr, cfg.Redis.Password, cfg.Redis.DataBase)

	b, err := bot.NewBot(cfg, strg)
	if err != nil {
		log.Fatal("❌ Botni yaratishda xatolik:", err)
	}

	b.Start()
}

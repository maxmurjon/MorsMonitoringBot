package bot

import (
	"log"
	"morc/bot/handlers"
	"morc/config"
	"morc/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	handlers *handlers.Handlers
	strg     storage.StorageRepoI
}

// NewBot – Botni yaratish funksiyasi
func NewBot(cfg *config.Config, strg storage.StorageRepoI) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}

	// Handlersni yaratamiz, va `storage`ni ham ulaymiz
	handler := handlers.NewHandlers(botAPI, strg)

	return &Bot{
		api:      botAPI,
		handlers: handler,
		strg:     strg,
	}, nil
}

// Start – Botni ishga tushiradigan funksiya
func (b *Bot) Start() {
	log.Printf("🚀 Bot @%s ishga tushdi", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	// Xabarlarni qabul qilish va handlerga yuborish
	for update := range updates {
		b.handlers.HandleUpdate(b.api, update)
	}
}

package handlers

import (
	"context"
	"log"
	"morc/models"
	"morc/storage"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ROLE_COURIER  = "🚴 Kuryer"
	ROLE_SELLER   = "🏣 Sotuvchi"
	ADMIN_CHAT_ID = 1023205318
)

type Handlers struct {
	bot           *tgbotapi.BotAPI
	strg          storage.StorageRepoI
	userHandler   *UserHandler
	barrelHandler *BarrelHandler
}

func NewHandlers(bot *tgbotapi.BotAPI, strg storage.StorageRepoI) *Handlers {
	return &Handlers{
		bot:           bot,
		userHandler:   NewUserHandler(bot, strg),
		barrelHandler: NewBarrelHandler(bot, strg),
		strg:          strg,
	}
}

func (h *Handlers) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Error sending message to chatID %d: %v", chatID, err)
	}
}

func (h *Handlers) getUserRole(userID int64) string {
	telegramID := strconv.Itoa(int(userID))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.strg.User().GetUserByTelegramID(ctx, &models.UserPrimaryKey{TelegramId: telegramID})
	if err != nil {
		log.Printf("❌ User not found or not in the database: %v", err)
		return "user"
	}

	return user.Role
}

func (h *Handlers) HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Println("Update Received:", update.UpdateID)
	if update.Message != nil {
		log.Println("Message Received:", update.Message.Text)
		h.HandleMessage(update)
	} else if update.CallbackQuery != nil {
		log.Println("Callback Query Received:", update.CallbackQuery.Data)
		h.HandleCallback(bot,update.CallbackQuery)
	}
}

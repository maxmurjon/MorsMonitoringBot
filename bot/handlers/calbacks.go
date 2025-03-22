package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handlers) HandleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	fmt.Println("Callback Query Received:", callback.Data)
	if callback.Message == nil {
		h.sendMessage(int64(callback.From.ID), "⚠️ Xatolik: Xabar topilmadi.")
		return
	}

	chatID := callback.Message.Chat.ID
	data := callback.Data

	log.Println("📩 Callback Query:", data)

	switch {
	case strings.HasPrefix(data, "edit_barrel_"):
		barrelIDStr := strings.TrimPrefix(data, "edit_barrel_")
		if barrelIDStr == "" {
			h.sendMessage(chatID, "⚠️ Xatolik: ID topilmadi.")
			return
		}
		h.barrelHandler.startEditingBarrel(chatID, callback.Message.MessageID, barrelIDStr)

	case strings.HasPrefix(data, "delete_barrel_"):
		barrelIDStr := strings.TrimPrefix(data, "delete_barrel_")
		if barrelIDStr == "" {
			h.sendMessage(chatID, "⚠️ Xatolik: ID topilmadi.")
			return
		}
		barrelID, err := strconv.Atoi(barrelIDStr)
		if err != nil {
			h.sendMessage(chatID, "⚠️ Xatolik: ID noto‘g‘ri formatda.")
			return
		}
		h.barrelHandler.confirmDeleteBarrel(chatID, callback.Message.MessageID, barrelID)

	case strings.HasPrefix(data, "confirm_delete_"):
		barrelIDStr := strings.TrimPrefix(data, "confirm_delete_")
		if barrelIDStr == "" {
			h.sendMessage(chatID, "⚠️ Xatolik: ID topilmadi.")
			return
		}
		barrelID, err := strconv.Atoi(barrelIDStr)
		if err != nil {
			h.sendMessage(chatID, "⚠️ Xatolik: ID noto‘g‘ri formatda.")
			return
		}
		h.barrelHandler.deleteBarrel(chatID, int64(barrelID))

	case strings.HasPrefix(data, "select_barrel_"):
		barrelIDStr := strings.TrimPrefix(data, "select_barrel_")
		barrelID, err := strconv.Atoi(barrelIDStr)
		if err != nil {
			h.sendMessage(chatID, "⚠️ Xatolik: ID noto‘g‘ri formatda.")
			return
		}
		h.barrelHandler.SelectBarrelForSeller(chatID, barrelID)

	case strings.HasPrefix(data, "assign_barrel_"):
		parts := strings.Split(data, "_")

		barrelID, err := strconv.Atoi(parts[2])
		if err != nil {
			h.sendMessage(chatID, "⚠️ Xatolik: barrel ID noto‘g‘ri formatda.")
			return
		}
		sellerID, err := strconv.Atoi(parts[4])
		if err != nil {
			h.sendMessage(chatID, "⚠️ Xatolik: sotuvchi ID noto‘g‘ri formatda.")
			return
		}
		h.barrelHandler.AssignBarrelToSeller(chatID, barrelID, sellerID)

	case strings.HasPrefix(data, "confirm_user_"):
		parts := strings.Split(data, "_")

		h.userHandler.ConfirmeUser(parts[2],callback)
	
	case strings.HasPrefix(data, "delete_user_"):
		parts := strings.Split(data, "_")

		h.userHandler.NonActivateUser(parts[2],callback)
	
	case strings.HasPrefix(data, "select_empty_barrel_"):
		parts := strings.Split(data, "_")

		h.barrelHandler.sendLocation(chatID,parts[3])

	default:
		h.sendMessage(chatID, "⚠️ Noma'lum callback!")
	}
}

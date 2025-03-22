package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"morc/models"
	"morc/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BarrelHandler struct {
	bot  *tgbotapi.BotAPI
	strg storage.StorageRepoI
}

var barrelCreationState = make(map[int64]*models.CreateBarrel)

func NewBarrelHandler(bot *tgbotapi.BotAPI, strg storage.StorageRepoI) *BarrelHandler {
	return &BarrelHandler{bot: bot, strg: strg}
}

func (h *BarrelHandler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	h.bot.Send(msg)
}

func (h *BarrelHandler) handleNewBarrel(chatID int64) {
	barrelCreationState[chatID] = &models.CreateBarrel{}

	locationKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation("📍 Joylashuvni yuborish"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, "📍 Iltimos, bochkani joylashuvini yuboring yoki qo‘lda kiriting:")
	msg.ReplyMarkup = locationKeyboard
	h.bot.Send(msg)
}

func (h *BarrelHandler) handleLocation(chatID int64, location *tgbotapi.Location) {
	barrel, exists := barrelCreationState[chatID]
	if !exists {
		h.sendMessage(chatID, "⚠️ Avval barrel qo‘shish jarayonini boshlang.")
		return
	}

	barrel.Latitude = location.Latitude
	barrel.Longitude = location.Longitude
	h.sendMessage(chatID, "📏 Iltimos, barrel hajmini litrda kiriting:")
}

func (h *BarrelHandler) handleBarrelVolume(chatID int64, text string) {
	barrel, exists := barrelCreationState[chatID]
	if !exists {
		h.sendMessage(chatID, "⚠️ Avval barrel qo‘shish jarayonini boshlang.")
		return
	}

	volume, err := strconv.ParseFloat(text, 64)
	if err != nil || volume <= 0 {
		h.sendMessage(chatID, "⚠️ Iltimos, to‘g‘ri hajm kiriting (masalan, 200.5):")
		return
	}

	barrel.VolumeLiters = volume
	barrel.CurrentVolume = volume
	barrel.CurrentVolume = 0.00 // Barrel bo‘sh holatda yaratiladi

	h.sendMessage(chatID, "✏️ Iltimos, barrel nomini kiriting:")
}

func (h *BarrelHandler) handleBarrelName(chatID int64, text string) {
	barrel, exists := barrelCreationState[chatID]
	if !exists {
		h.sendMessage(chatID, "⚠️ Avval barrel qo‘shish jarayonini boshlang.")
		return
	}

	barrel.Name = text
	h.sendMessage(chatID, "🏠 Iltimos, barrelning joylashuv nomini kiriting:")
}

func (h *BarrelHandler) handleLocationName(chatID int64, text string) {
	barrel, exists := barrelCreationState[chatID]
	if !exists {
		h.sendMessage(chatID, "⚠️ Avval barrel qo‘shish jarayonini boshlang.")
		return
	}

	barrel.LocationName = text

	// ❗ Agar sotuvchi ID kiritish kerak bo‘lsa
	h.sendMessage(chatID, "👤 Iltimos, barrelga tayinlanadigan sotuvchi ID ni kiriting yoki 'yo'q' deb yozing:")
}

func (h *BarrelHandler) handleSellerID(chatID int64, text string) {
	barrel, exists := barrelCreationState[chatID]
	if !exists {
		h.sendMessage(chatID, "⚠️ Avval barrel qo‘shish jarayonini boshlang.")
		return
	}

	if strings.EqualFold(text, "yo'q") {
		sellerID, err := strconv.Atoi(text)
		if err != nil || sellerID <= 0 {
			h.sendMessage(chatID, "⚠️ Iltimos, to‘g‘ri sotuvchi ID kiriting yoki 'yo'q' deb yozing:")
			return
		}
		barrel.AssignedSellerId = &sellerID
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createdBarrel, err := h.strg.Barrel().Create(ctx, barrel)
	if err != nil {
		h.sendMessage(chatID, "❌ Barrel yaratishda xatolik yuz berdi.")
		log.Println("Error creating barrel:", err)
		return
	}

	h.sendMessage(chatID, fmt.Sprintf(
		"✅ Yangi barrel yaratildi!\n📍 Joylashuv: %s (%.6f, %.6f)\n📏 Hajm: %.2f litr\n👤 Sotuvchi ID: %v",
		createdBarrel.LocationName, createdBarrel.Latitude, createdBarrel.Longitude,
		createdBarrel.VolumeLiters, createdBarrel.AssignedSellerId,
	))

	delete(barrelCreationState, chatID)
}

func (h *BarrelHandler) handleGetBarrels(chatID int64) {
	barrels, err := h.strg.Barrel().GetList(context.Background(), &models.GetListBarrelRequest{})
	if err != nil {
		h.sendMessage(chatID, "❌ Bochka ma'lumotlari olinmadi. Keyinroq urinib ko'ring.")
		log.Println("Error getting barrels:", err)
		return
	}

	if len(barrels.Barrels) == 0 {
		h.sendMessage(chatID, "📭 Hozircha hech qanday bochka mavjud emas.")
		return
	}

	var messages []tgbotapi.MessageConfig

	for _, barrel := range barrels.Barrels {
		text := fmt.Sprintf(
			"🛢 *%s*\n📍 *Manzil:* %s (%.6f, %.6f)\n📏 *Hajm:* %.2f L\n👤 *Sotuvchi ID:* %v",
			barrel.Name, barrel.LocationName, barrel.Latitude, barrel.Longitude,
			barrel.VolumeLiters, barrel.AssignedSellerId,
		)

		// Inline tugmalar
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✏️ Tahrirlash", fmt.Sprintf("edit_barrel_%d", barrel.Id)),
				tgbotapi.NewInlineKeyboardButtonData("🗑 O‘chirish", fmt.Sprintf("delete_barrel_%d", barrel.Id)),
			),
		)

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = inlineKeyboard
		messages = append(messages, msg)
	}

	// Foydalanuvchiga barcha xabarlarni yuborish
	for _, msg := range messages {
		h.bot.Send(msg)
	}
}


func (h *BarrelHandler) startEditingBarrel(chatID int64, messageID int, barrelID string) {
	h.sendMessage(chatID, "✏️ Barrelni tahrirlash uchun yangi ma'lumotlarni yuboring. Hozircha faqat nomini o‘zgartirish mumkin.")
	// barrelCreationState[chatID] = &models.CreateBarrel{Name: barrelID}
}

func (h *BarrelHandler) confirmDeleteBarrel(chatID int64, messageID int, barrelID int) {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Ha, o‘chirish", fmt.Sprintf("confirm_delete_%d", barrelID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Yo‘q, bekor qilish", "cancel_delete"),
		),
	)

	msg := tgbotapi.NewEditMessageText(chatID, messageID, "⚠️ Siz ushbu barrelni o‘chirmoqchimisiz?")
	msg.ReplyMarkup = &inlineKeyboard
	h.bot.Send(msg)
}

func (h *BarrelHandler) deleteBarrel(chatID int64, barrelID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.strg.Barrel().Delete(ctx, barrelID)
	if err != nil {
		h.sendMessage(chatID, "❌ Barrelni o‘chirishda xatolik yuz berdi.")
		log.Println("Error deleting barrel:", err)
		return
	}

	h.sendMessage(chatID, "✅ Barrel muvaffaqiyatli o‘chirildi.")
}

func (h *BarrelHandler) GaveBarrelToSaller(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	// 1. Biriktirilmagan bochkalarni olish
	barrels, err := h.strg.Barrel().GetListSellerId(context.Background(), &models.GetListBarrelRequest{})
	if err != nil {
		h.sendMessage(chatID, "⚠️ Xatolik: bochkalar olinmadi.")
		return
	}

	if len(barrels.Barrels) == 0 {
		h.sendMessage(chatID, "✅ Hozircha barcha bochkalar biriktirilgan.")
		return
	}

	// 2. Foydalanuvchiga tanlash uchun inline tugmalar yaratish
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	for _, barrel := range barrels.Barrels {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("🛢 %s (%fL)", barrel.Name, barrel.CurrentVolume),
			fmt.Sprintf("select_barrel_%d", barrel.Id),
		)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{btn})
	}

	msg := tgbotapi.NewMessage(chatID, "🛢 Biriktirilmagan bochkani tanlang:")
	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}

// 3. Barrel tanlanganidan keyin sotuvchilar ro‘yxatini ko‘rsatish
func (h *BarrelHandler) SelectBarrelForSeller(chatID int64, barrelID int) {
	// 4. Sotuvchilarni olish
	sellers, err := h.strg.User().GetList(context.Background(), &models.GetListUserRequest{})
	if err != nil {
		h.sendMessage(chatID, "⚠️ Xatolik: sotuvchilar olinmadi.")
		return
	}

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	for _, seller := range sellers.Users {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("👤 %s\n %s", seller.FirstName, seller.PhoneNumber),
			fmt.Sprintf("assign_barrel_%d_to_%s", barrelID, seller.Id),
		)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{btn})
	}

	msg := tgbotapi.NewMessage(chatID, "👤 Sotuvchini tanlang:")
	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}

// 5. Barrelni sotuvchiga biriktirish
func (h *BarrelHandler) AssignBarrelToSeller(chatID int64, barrelID int, sellerID int) {
	_, err := h.strg.Barrel().Update(context.Background(), &models.UpdateBarrel{
		Id:               barrelID,
		AssignedSellerId: &sellerID,
	})
	if err != nil {
		h.sendMessage(chatID, "⚠️ Xatolik: bochka sotuvchiga biriktirilmadi.")
		return
	}

	h.sendMessage(chatID, "✅ Bochka sotuvchiga muvaffaqiyatli biriktirildi!")
}

func (h *BarrelHandler) handleGetEmptyBarrels(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	barrels, err := h.strg.Barrel().GetListEmpty(context.Background(), &models.GetListBarrelRequest{})
	if err != nil {
		h.sendMessage(chatID, "⚠️ Xatolik: bochkalar olinmadi.")
		return
	}

	if len(barrels.Barrels) == 0 {
		h.sendMessage(chatID, "✅ Hozircha bo‘sh bochka mavjud emas.")
		return
	}

	// 2. Foydalanuvchiga tanlash uchun inline tugmalar yaratish
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	for _, barrel := range barrels.Barrels {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("🛢 %s (%fL)", barrel.Name, barrel.CurrentVolume),
			fmt.Sprintf("select_empty_barrel_%d", barrel.Id),
		)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{btn})
	}

	msg := tgbotapi.NewMessage(chatID, "🛢 Bo'sh bo'chkalar ro'yhati:")
	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}

func (h *BarrelHandler) sendLocation(chatID int64, barrelID string) {
	// 1. Bochka joylashuvi
	id, err := strconv.ParseInt(barrelID, 10, 64)
	if err != nil {
		h.sendMessage(chatID, "❌ Xatolik: barrel ID noto‘g‘ri.")
		log.Println("Error converting barrelID to int64:", err)
		return
	}

	barrel, err := h.strg.Barrel().GetByID(context.Background(), id)
	if err != nil {
		h.sendMessage(chatID, "❌ Xatolik: barrel olinmadi.")
		log.Println("Error getting barrel:", err)
		return
	}

	msg := tgbotapi.NewLocation(chatID, barrel.Latitude, barrel.Longitude)
	h.bot.Send(msg)
}
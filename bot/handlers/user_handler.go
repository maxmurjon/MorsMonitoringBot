package handlers

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"morc/bot/keyboards"
	"morc/models"
	"morc/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	AdminChatID = 1023205318
	Timeout     = 5 * time.Minute
)

type UserHandler struct {
	bot  *tgbotapi.BotAPI
	strg storage.StorageRepoI
	mux  sync.Mutex
}

var userCreationState = make(map[int64]*models.CreateUser)
var userLogin = make(map[int64]*models.User)

func NewUserHandler(bot *tgbotapi.BotAPI, strg storage.StorageRepoI) *UserHandler {
	return &UserHandler{bot: bot, strg: strg}
}

func (h *UserHandler) HandleRegistration(msg *tgbotapi.Message) {
	h.mux.Lock()
	defer h.mux.Unlock()

	tgID := fmt.Sprint(msg.From.ID)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingUser, err := h.strg.User().GetUserByTelegramID(ctx, &models.UserPrimaryKey{TelegramId: tgID})
	if err == nil && existingUser != nil {
		if existingUser.IsVerified {
			h.sendMessage(msg.Chat.ID, "✅ Siz allaqachon ro'yxatdan o'tgansiz!")
		} else {
			h.sendMessage(msg.Chat.ID, "⏳ Admin tasdiqlashini kuting...")
		}
		return
	}

	userCreationState[msg.Chat.ID] = &models.CreateUser{
		TelegramId: tgID,
		FirstName:  msg.From.FirstName,
	}

	h.requestLastName(msg.Chat.ID)
	h.startTimeout(msg.Chat.ID)
}

func (h *UserHandler) HandleFirstName(msg *tgbotapi.Message) {
	if user, exists := userCreationState[msg.Chat.ID]; exists {
		user.FirstName = msg.Text
		h.requestLastName(msg.Chat.ID)
	} else {
		h.sendMessage(msg.Chat.ID, "❌ Ro'yxatdan o'tish jarayoni boshlanmagan")
	}
}

func (h *UserHandler) HandleLastName(msg *tgbotapi.Message) {
	if user, exists := userCreationState[msg.Chat.ID]; exists {
		user.LastName = msg.Text
		h.requestPhoneNumber(msg.Chat.ID)
	} else {
		h.sendMessage(msg.Chat.ID, "❌ Ro'yxatdan o'tish jarayoni boshlanmagan")
	}
}

func (h *UserHandler) HandlePhoneNumber(msg *tgbotapi.Message) {
	if user, exists := userCreationState[msg.Chat.ID]; exists {
		if msg.Contact == nil {
			h.sendMessage(msg.Chat.ID, "❌ Iltimos, telefon raqamingizni yuboring!")
			h.requestPhoneNumber(msg.Chat.ID)
			return
		}
		user.PhoneNumber = msg.Contact.PhoneNumber
		h.sendMessage(msg.Chat.ID, "✅ Telefon raqamingiz qabul qilindi.")
		h.requestRoleSelection(msg.Chat.ID)
	} else {
		h.sendMessage(msg.Chat.ID, "❌ Ro'yxatdan o'tish jarayoni boshlanmagan")
	}
}

func (h *UserHandler) HandleRoleSelection(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	text := msg.Text

	// Faqat ruxsat etilgan rollarni qabul qilish
	allowedRoles := map[string]bool{
		"🏣 Sotuvchi": true,
		"🚴 Kuryer":   true,
	}

	if _, exists := allowedRoles[text]; !exists {
		h.bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Noto‘g‘ri rol! Iltimos, tugmalardan birini tanlang."))
		return
	}

	// Foydalanuvchi rolini saqlash
	userCreationState[chatID].Role = text

	// Ro‘yxatdan o‘tish yakunlanganini bildirish
	h.CompleteRegistration(chatID)
}

func (h *UserHandler) requestLastName(chatID int64) {
	h.sendMessage(chatID, "✍️ Iltimos, familiyangizni kiriting:")
}

func (h *UserHandler) requestPhoneNumber(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "📞 Telefon raqamingizni yuboring:")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButtonContact("Raqamni ulashish")),
	)
	h.bot.Send(msg)
}

func (h *UserHandler) requestRoleSelection(chatID int64) {
	h.sendMessageWithKeyboard(chatID, "👔 Iltimos, rolingizni tanlang:", keyboards.GetRoleSelectionMenu())
}

func (h *UserHandler) sendMessage(chatID int64, text string) {
	h.bot.Send(tgbotapi.NewMessage(chatID, text))
}

func (h *UserHandler) sendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}

func (h *UserHandler) startTimeout(chatID int64) {
	go func() {
		time.Sleep(Timeout)
		h.mux.Lock()
		defer h.mux.Unlock()
		if _, exists := userCreationState[chatID]; exists {
			delete(userCreationState, chatID)
			h.sendMessage(chatID, "⏳ Ro'yxatdan o'tish vaqti tugadi. Iltimos, qayta urinib ko'ring!")
		}
	}()
}

func (h *UserHandler) CompleteRegistration(chatID int64) {
	user := userCreationState[chatID]
	_, err := h.strg.User().Create(context.Background(), user)
	if err != nil {
		h.sendMessage(chatID, "❌ Ro‘yxatdan o‘tishda xatolik yuz berdi. Iltimos, qayta urinib ko‘ring!")
		return
	}

	msg := tgbotapi.NewMessage(chatID, "✅ Ro‘yxatdan o‘tish muvaffaqiyatli yakunlandi!\n\n"+
		"👤 Ism: "+user.FirstName+"\n"+
		"👤 Familiya: "+user.LastName+"\n"+
		"📞 Telefon: "+user.PhoneNumber+"\n"+
		"🛠 Rol: "+user.Role)

	// ReplyKeyboardni o'chirish
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	// Xabarni jo‘natish
	h.bot.Send(msg)

	h.sendMessage(chatID, "Admin tasdiqlashini kuting tasdiqlaganidan so'ng, sizga habar keladi va menu ochiladi")
	// Foydalanuvchini bazaga saqlash yoki keyingi jarayonlarga yo‘naltirish
	delete(userCreationState, chatID)
}

func (h *UserHandler) LoginSeller(msg *tgbotapi.Message) {
	chatID := strconv.Itoa(int(msg.Chat.ID))
	fmt.Println(chatID)

	userLogin[msg.Chat.ID] = &models.User{
		TelegramId: chatID,
	}

	h.strg.User().GetByID(context.Background(), &models.UserPrimaryKey{TelegramId: chatID})

	// Agar joylashuv hali kelmagan bo'lsa, so'rash
	if msg.Location == nil {
		locationKeyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonLocation("📍 Joylashuvni yuborish"),
			),
		)
		message := tgbotapi.NewMessage(msg.From.ID, "📍 Iltimos, bochkani joylashuvini yuboring yoki qo‘lda kiriting:")
		message.ReplyMarkup = locationKeyboard
		h.bot.Send(message)
		return
	}

	telegramId := fmt.Sprint(msg.From.ID)

	// Joylashuvni tekshirish
	if h.isValidLocation(telegramId, msg.Location) {
		fmt.Println("Joylashuv")
		h.sendMessage(msg.From.ID, "✅ Joylashuvingiz qabul qilindi.")
		message := tgbotapi.NewMessage(msg.From.ID, "📍 Joylashuvingiz qabul qilindi.")
		message.ReplyMarkup = keyboards.SellerMenu()
		h.bot.Send(message)
	} else {
		h.sendMessage(msg.From.ID, "❌ Noto'g'ri joylashuv. Iltimos, belgilangan joyda ekanligingizga ishonch hosil qiling.")
	}
}

func (h *UserHandler) isValidLocation(userId string, location *tgbotapi.Location) bool {
	user, err := h.strg.User().GetUserByTelegramID(context.Background(), &models.UserPrimaryKey{TelegramId: userId})
	if err != nil {
		return false
	}
	barrel, err := h.strg.Barrel().GetBarrelBySellerId(context.Background(), user.Id)
	if err != nil {
		return false
	}
	fmt.Println(barrel.Latitude, barrel.Longitude, location.Latitude, location.Longitude)
	// Koordinatalar farqi 0.0001 ichida bo‘lsa, joylashuv to‘g‘ri deb hisoblanadi
	const epsilon = 0.0001
	if math.Abs(barrel.Latitude-location.Latitude) < epsilon && math.Abs(barrel.Longitude-location.Longitude) < epsilon {
		return true
	}
	return false
}

func (h *UserHandler) handleGetCouriers(chatID int64) {
	users, err := h.strg.User().GetByRole(context.Background(), "🚴 Kuryer")
	if err != nil {
		h.sendMessage(chatID, "❌ Foydalanuvchi ma'lumotlari olinmadi. Keyinroq urinib ko'ring.")
		log.Println("Error getting users:", err)
		return
	}

	if len(users) == 0 {
		h.sendMessage(chatID, "📭 Hozircha hech qanday foydalanuvchi mavjud emas.")
		return
	}

	h.sendUserData(users, chatID)

}

func (h *UserHandler) handleGetSellers(chatID int64) {
	users, err := h.strg.User().GetByRole(context.Background(), "🏣 Sotuvchi")
	if err != nil {
		h.sendMessage(chatID, "❌ Foydalanuvchi ma'lumotlari olinmadi. Keyinroq urinib ko'ring.")
		log.Println("Error getting users:", err)
		return
	}

	if len(users) == 0 {
		h.sendMessage(chatID, "📭 Hozircha hech qanday foydalanuvchi mavjud emas.")
		return
	}

	h.sendUserData(users, chatID)

}


func (h *UserHandler) sendUserData(users []models.User, chatID int64) {
	var messages []tgbotapi.MessageConfig

	for _, user := range users {
		text := fmt.Sprintf(
			"👤 *Ism:* %s\n👤 *Familiya:* %s\nId *Id:* %s\n📞 *Telefon:* %s\n🛠 *Rol:* %s",
			user.FirstName, user.LastName, user.TelegramId, user.PhoneNumber, user.Role,
		)

		status := ""
		callback := ""
		if !user.IsVerified {
			status = "\n✅ Tasdiqlash"
			callback = fmt.Sprintf("confirm_user_%s", user.TelegramId)
		} else {
			status = "\n🗑 O‘chirish"
			callback = fmt.Sprintf("delete_user_%s", user.TelegramId)
		}
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(status, callback),
			),
		)

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = inlineKeyboard
		messages = append(messages, msg)
	}

	for _, msg := range messages {
		h.bot.Send(msg)
	}
}


func (h *UserHandler) ConfirmeUser(userID string, callback *tgbotapi.CallbackQuery) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user,err:=h.strg.User().GetUserByTelegramID(context.Background(),&models.UserPrimaryKey{TelegramId: userID})
	if err!=nil{
		h.sendMessage(int64(callback.From.ID), "❌ Foydalanuvchi tasdiqlanmadi. Keyin qayta urinib ko'ring.")
		log.Println("Error confirming user:", err) 		
		return
	}

	err = h.strg.User().Approve(ctx, user.Id)
	if err != nil {
		h.sendMessage(int64(callback.From.ID), "❌ Foydalanuvchi tasdiqlanmadi. Keyin qayta urinib ko'ring.")
		log.Println("Error confirming user:", err) 		
		return
	}

	h.sendMessage(int64(callback.From.ID), "✅ Foydalanuvchi muvaffaqiyatli tasdiqlandi.")
}

func (h *UserHandler) NonActivateUser(userID string, callback *tgbotapi.CallbackQuery) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user,err:=h.strg.User().GetUserByTelegramID(context.Background(),&models.UserPrimaryKey{TelegramId: userID})
	if err!=nil{
		h.sendMessage(int64(callback.From.ID), "❌ Foydalanuvchi o'chirilmadi. Keyin qayta urinib ko'ring.")
		log.Println("Error confirming user:", err) 		
		return
	}

	err = h.strg.User().Reject(ctx, user.Id)
	if err != nil {
		h.sendMessage(int64(callback.From.ID), "❌ Foydalanuvchi o'chirilmadi. Keyin qayta urinib ko'ring.")
		log.Println("Error confirming user:", err) 		
		return
	}

	h.sendMessage(int64(callback.From.ID), "✅ Foydalanuvchi muvaffaqiyatli o'chirildi.")
}


func (h *UserHandler) checkUserValidation(telegramId string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.strg.User().GetUserByTelegramID(ctx, &models.UserPrimaryKey{TelegramId: telegramId})
	if err != nil {
		log.Println("Error checking user validation:", err)
		return false
	}

	return user.IsVerified
}

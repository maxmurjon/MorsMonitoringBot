package handlers

import (
	"fmt"
	"log"
	"morc/bot/keyboards"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handlers) HandleMessage(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	role := h.getUserRole(userID)
	text := update.Message.Text

	// Agar foydalanuvchi joylashuv jo‘natsa
	if update.Message.Location != nil {
		if _, exists := barrelCreationState[chatID]; exists {
			h.barrelHandler.handleLocation(chatID, update.Message.Location)
			return
		} else if _, exists := userLogin[chatID]; exists {
			h.userHandler.LoginSeller(update.Message)
			return
		}
	}

	// Agar foydalanuvchi barrel qo'shish jarayonida bo'lsa
	if barrel, exists := barrelCreationState[chatID]; exists {
		fmt.Println("Barrel")
		switch {
		case barrel.VolumeLiters == 0:
			h.barrelHandler.handleBarrelVolume(chatID, text)
		case barrel.Name == "":
			h.barrelHandler.handleBarrelName(chatID, text)
		case barrel.LocationName == "":
			h.barrelHandler.handleLocationName(chatID, text)
		default:
			h.barrelHandler.handleSellerID(chatID, text)
		}
		return
	}

	// Agar foydalanuvchi joylashuv jo‘natsa
	if update.Message.Contact != nil {
		if _, exists := userCreationState[chatID]; exists {
			h.userHandler.HandlePhoneNumber(update.Message)
		}
		return
	}

	// Agar foydalanuvchi ro‘yxatdan o‘tish jarayonida bo‘lsa
	if user, exists := userCreationState[update.Message.Chat.ID]; exists {
		chatID := update.Message.Chat.ID

		switch {
		case user.FirstName == "":
			h.userHandler.HandleFirstName(update.Message)
		case user.LastName == "":
			h.userHandler.HandleLastName(update.Message)
		case user.PhoneNumber == "":
			h.userHandler.HandlePhoneNumber(update.Message)
		case user.Role == "":
			h.userHandler.HandleRoleSelection(update.Message)
		default:
			h.userHandler.CompleteRegistration(chatID)
		}

		return
	}

	// Asosiy menyu
	var msg tgbotapi.MessageConfig
	switch text {
	case "/start":
		msg = tgbotapi.NewMessage(chatID, "👋 Salom! Kerakli bo‘limni tanlang:")
		routeMenu := h.RouteMenu(strconv.Itoa(int(chatID)),role)
		msg.ReplyMarkup = routeMenu

	case "🧑‍💼 Hodimlar":
		msg = tgbotapi.NewMessage(chatID, "🧑‍💼 Hodimlar bo‘limi:")
		msg.ReplyMarkup = keyboards.AdminUsersMenu()

	case "🛢 Bochkalar":
		msg = tgbotapi.NewMessage(chatID, "🛢 Bochkalar bo‘limi:")
		msg.ReplyMarkup = keyboards.AdminBarrelsMenu()

	case "🔙 Ortga qaytish":
		msg = tgbotapi.NewMessage(chatID, "🏠 Asosiy menyu:")
		roleMenu := h.RouteMenu(strconv.Itoa(int(chatID)),role)
		msg.ReplyMarkup = roleMenu

	case "➕ Bochka qo'shish":
		h.barrelHandler.handleNewBarrel(chatID)
		return

	case "📝 Bochkalar ro'yxati":
		h.barrelHandler.handleGetBarrels(chatID)
		return

	case "🧑‍💼 Bochkani biriktirish":
		h.barrelHandler.GaveBarrelToSaller(update)
		return

	case "Ro'yhatdan o'tish":
		h.userHandler.HandleRegistration(update.Message)
		return

	case "🔑 Kirish":
		h.userHandler.LoginSeller(update.Message)
		return
	
	case "🚴 Kuryerlar ro'yxati":
		h.userHandler.handleGetCouriers(chatID)
		return
	
	case "🏣 Sotuvchi ro'yxati":
		h.userHandler.handleGetSellers(chatID)

	case "📥 Bosh bo'chkalar":
		h.barrelHandler.handleGetEmptyBarrels(update)
		return
	
	case "🛍 Sotish":
		msg = tgbotapi.NewMessage(chatID, "🥤 Sotish")
		msg.ReplyMarkup = keyboards.SellerSellMenu()

	default:
		msg = tgbotapi.NewMessage(chatID, "⚠️ Noto‘g‘ri buyruq!")
	}

	// Xabarni jo'natish
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Error sending message to chatID %d: %v", chatID, err)
	}
}

func (h *Handlers) RouteMenu(telegramId,role string) tgbotapi.ReplyKeyboardMarkup {
	switch role {
	case "admin":
		return keyboards.AdminMenu()
	case ROLE_SELLER:
		if h.userHandler.checkUserValidation(telegramId){
			return keyboards.LoginSellerMenu()
		}else{
			id, _ := strconv.ParseInt(telegramId, 10, 64)
			h.sendMessage(id, "⏳ Admin tasdiqlashini kuting")
			return tgbotapi.ReplyKeyboardMarkup{}
		}
	case ROLE_COURIER:
		if h.userHandler.checkUserValidation(telegramId){
			return keyboards.CourierMenu()
		}else{
			id, _ := strconv.ParseInt(telegramId, 10, 64)
			h.sendMessage(id, "⏳ Admin tasdiqlashini kuting")
			return tgbotapi.ReplyKeyboardMarkup{}
		}
	default:
		return keyboards.GetRegistrationMenu()
	}
}


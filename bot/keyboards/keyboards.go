package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	ROLE_COURIER  = "🚴 Kuryer"
	ROLE_SELLER   = "🏣 Sotuvchi"
	ADMIN_CHAT_ID = 1023205318
)

// **Admin asosiy menyusi**
func AdminMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🛢 Bochkalar"),
			tgbotapi.NewKeyboardButton("🧑‍💼 Hodimlar"),
			tgbotapi.NewKeyboardButton("📊 Statistikalar"),
		),
		tgbotapi.NewKeyboardButtonRow(

			tgbotapi.NewKeyboardButton("🔙 Ortga qaytish"),
		),
	)
}

// **Admin -> Buyurtmalar sub-menu**
func AdminBarrelsMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("➕ Bochka qo'shish"),
			tgbotapi.NewKeyboardButton("📝 Bochkalar ro'yxati"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🧑‍💼 Bochkani biriktirish"),
			tgbotapi.NewKeyboardButton("🔙 Ortga qaytish"),
		),
	)
}

func AdminUsersMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🚴 Kuryerlar ro'yxati"),
			tgbotapi.NewKeyboardButton("🏣 Sotuvchi ro'yxati"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔙 Ortga qaytish"),
		),
	)
}

func LoginSellerMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔑 Kirish"),
		),
	)
}

// **User asosiy menyusi**
func SellerMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🛍 Sotish"),
			tgbotapi.NewKeyboardButton("🚛 Bochkani kuryerga berish"),
		),
	)
}

func SellerSellMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🥤\n(200ml)"),
			tgbotapi.NewKeyboardButton("🥤\n(400ml)"),
			tgbotapi.NewKeyboardButton("🥤\n(1L)"),
			tgbotapi.NewKeyboardButton("🥤\n(5L)"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔙 Ortga qaytish"),
		),
	)
}

// **User -> Sozlamalar sub-menu**
func CourierMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📥 Bosh bo'chkalar"),
			// tgbotapi.NewKeyboardButton("🔑 Parolni almashtirish"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔙 Ortga qaytish"), // Asosiy menu
		),
	)
}

func UserMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🛍 Buyurtma berish"),
			tgbotapi.NewKeyboardButton("🛒 Savatcha"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔙 Ortga qaytish"),
		),
	)
}

func GetRegistrationMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Ro'yhatdan o'tish"),
		),
	)
}

// GetMainMenu - Asosiy menyuni hosil qiladi
func GetMainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("👤 Profil"),
			tgbotapi.NewKeyboardButton("📋 Buyurtmalar"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📞 Bog‘lanish"),
			tgbotapi.NewKeyboardButton("⚙️ Sozlamalar"),
		),
	)
}

func GetRoleSelectionMenu() tgbotapi.ReplyKeyboardMarkup {
	buttons := [][]tgbotapi.KeyboardButton{
		{tgbotapi.NewKeyboardButton(ROLE_COURIER), tgbotapi.NewKeyboardButton(ROLE_SELLER)},
	}
	return tgbotapi.NewReplyKeyboard(buttons...)

}

// func CheckSituation() tgbotapi.ReplyKeyboardMarkup {
//     buttons := [][]tgbotapi.KeyboardButton{
//         {tgbotapi.NewKeyboardButton("🔄 Yangilash")},
//     }
//     return tgbotapi.ReplyKeyboardMarkup{
//         Keyboard:        buttons,
//         ResizeKeyboard:  true,  // Katta keyboardni kichraytiradi
//         OneTimeKeyboard: false, // Foydalanuvchi tanlagandan keyin ham ko‘rinib turadi
//     }
// }


package Telegram

import (
	"github.com/SadikSunbul/TelegramUrlBot/Database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ProcessUserInput(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *Database.DataBase) {
	chatID := update.Message.Chat.ID

	// Kullanıcının adımını kontrol et
	if data, ok := userData[chatID]; ok {
		switch data["step"] {
		case "name":
			// Kullanıcı ismini al
			userData[chatID]["name"] = update.Message.Text // Kullanıcının girdiği ismi al
			msg := tgbotapi.NewMessage(chatID, "Mail adresinizi giriniz.")
			bot.Send(msg)
			userData[chatID]["step"] = "email" // Adımı "email" olarak güncelle
		case "email":
			// Kullanıcı e-posta adresini al
			userData[chatID]["email"] = update.Message.Text
			// Kullanıcıdan bilgileri onaylamasını iste
			msg := tgbotapi.NewMessage(chatID,
				"Bu bilgiler sizin mi?\nİsim: "+userData[chatID]["name"]+"\nE-posta: "+userData[chatID]["email"]+"\n\nEvet için 'evet', hayır için 'hayır' yazın.")
			bot.Send(msg)
			userData[chatID]["step"] = "confirm"
		case "confirm":
			if update.Message.Text == "evet" {
				msg := tgbotapi.NewMessage(chatID, "Kaydınız başarıyla yapıldı!")
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(chatID, "Kaydınız iptal edildi.")
				bot.Send(msg)
			}
			delete(userData, chatID) // Kullanıcı verisini temizle
		}
	}
}

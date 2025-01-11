package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func HandleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userData map[int64]map[string]string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "İsminizi giriniz!")
	bot.Send(msg)

	// Kullanıcı verilerini başlat
	userData[message.Chat.ID] = make(map[string]string)
	userData[message.Chat.ID]["step"] = "name" // Kullanıcının adımını "name" olarak ayarla
}

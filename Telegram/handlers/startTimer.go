package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

func StartTimer(chatID int64, bot *tgbotapi.BotAPI, userData map[int64]map[string]string) {
	go func() {
		time.Sleep(10 * time.Second) // 1 dakika bekle
		if _, ok := userData[chatID]; ok {
			delete(userData, chatID) // Kullanıcı verisini sil
			msg := tgbotapi.NewMessage(chatID, "Zaman doldu, lütfen en baştan kaydınızı yapın.")
			bot.Send(msg) // Mesajı gönder
		}
	}()
}

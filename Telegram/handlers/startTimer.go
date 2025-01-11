package handlers

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func StartTimer(chatID int64, bot *tgbotapi.BotAPI, userData map[int64]map[string]string) {
	go func() {
		time.Sleep(2 * time.Minute)
		if _, ok := userData[chatID]; ok {
			msg := tgbotapi.NewMessage(chatID, "İşlem zaman aşımına uğradı, lütfen en baştan kaydınızı yapın.")
			bot.Send(msg)
			delete(userData, chatID) // chatID'ye ait veriyi sil
		}
	}()
}

package handlers

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func StartTimer(chatID int64, bot *tgbotapi.BotAPI, userData map[int64]map[string]string) {
	go func() {
		time.Sleep(1 * time.Minute)
		if _, ok := userData[chatID]; ok {
			delete(userData, chatID)
			msg := tgbotapi.NewMessage(chatID, "Zaman doldu, lütfen en baştan kaydınızı yapın.")
			bot.Send(msg)
		}
	}()
}

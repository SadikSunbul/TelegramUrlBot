package Telegram

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/SadikSunbul/TelegramUrlBot/Database"
	"github.com/SadikSunbul/TelegramUrlBot/Telegram/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func handleCallbackQuery(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *Database.DataBase) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	parts := strings.Split(data, ":")

	if len(parts) != 2 {
		return
	}

	action := parts[0]
	value := parts[1]

	switch action {
	case "auto_name":
		if value == "yes" {
			var shortUrl string
			for {
				shortUrl = generateShortUrl()
				available, err := isShortUrlAvailable(db, shortUrl)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, handlers.ErorrTelegram("Veri tabanı hatası oluştu.")))
					return
				}
				if available {
					break
				}
			}
			handlers.UserData[chatID]["shortUrl"] = shortUrl
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Kısa URL'niz: %s\nZamanlı mı yoksa zamansız mı olsun?", shortUrl))
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("⏱️ Zamanlı", "time_limit:timed"),
					tgbotapi.NewInlineKeyboardButtonData("♾️ Zamansız", "time_limit:unlimited"),
				),
			)
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, "Kısa URL'nizi giriniz (sadece şu kısmı girin xyz.com/'burayı'):"))
			handlers.UserData[chatID]["step"] = "get_custom_short_url"
		}

	case "time_limit":
		if value == "timed" {
			msg := tgbotapi.NewMessage(chatID, "Ne kadar süre geçerli olsun?")
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("1 Saat", "duration:1"),
					tgbotapi.NewInlineKeyboardButtonData("6 Saat", "duration:6"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("12 Saat", "duration:12"),
					tgbotapi.NewInlineKeyboardButtonData("24 Saat", "duration:24"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("📝 Manuel Tarih Gir", "duration:manual"),
				),
			)
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		} else {
			err := saveUrlToDatabase(handlers.UserData[chatID], db, chatID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, handlers.ErorrTelegram("Veri tabanı hatası oluştu.")))
				return
			}
			bot.Send(tgbotapi.NewMessage(chatID, "Kısa URL başarıyla oluşturuldu ve sınırsız olarak kullanılabilir."))
			delete(handlers.UserData, chatID)
		}

	case "duration":
		if value == "manual" {
			bot.Send(tgbotapi.NewMessage(chatID, "Lütfen son kullanma tarihini giriniz (YYYY-MM-DD HH:MM):"))
			handlers.UserData[chatID]["step"] = "get_expiration_date"
		} else {
			hours, _ := strconv.Atoi(value)
			expirationTime := time.Now().Add(time.Duration(hours) * time.Hour)
			handlers.UserData[chatID]["expirationDate"] = expirationTime.Format("2006-01-02 15:04")

			err := saveUrlToDatabase(handlers.UserData[chatID], db, chatID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, handlers.ErorrTelegram("Veri tabanı hatası oluştu.")))
				return
			}
			bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Kısa URL başarıyla oluşturuldu ve %d saat geçerli olacak.", hours)))
			delete(handlers.UserData, chatID)
		}
	}

	// Callback query'yi yanıtla
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	bot.AnswerCallbackQuery(callback)
}

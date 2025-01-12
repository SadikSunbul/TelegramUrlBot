package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/SadikSunbul/TelegramUrlBot/Database"
	"github.com/SadikSunbul/TelegramUrlBot/Models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleMyLinks(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *Database.DataBase, aktifmi bool) {
	mesaj := "🟢 Aktif Url leriniz:"
	if !aktifmi {
		mesaj = "🔴 Pasif Url leriniz:"
	}

	var urls []Models.Url
	currentTime := time.Now()

	if aktifmi {
		filter := bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "endDate", Value: bson.D{{Key: "$exists", Value: false}}}},
				bson.D{{Key: "endDate", Value: bson.D{{Key: "$gt", Value: currentTime}}}},
			}},
		}

		data, err := db.GetList(Database.Url, filter)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, ErorrTelegram("Veri tabanı hatası oluştu. Lütfen daha sonra tekrar deneyin.")))
			return
		}

		// URL'leri almak için bir döngü kullan
		for data.Next(context.TODO()) {
			var url Models.Url
			err := data.Decode(&url)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(message.Chat.ID, ErorrTelegram(fmt.Sprintf("Decode hatası oluştu: %s. Lütfen daha sonra tekrar deneyin.", err.Error()))))
				return
			}
			urls = append(urls, url)
		}

		if err := data.Err(); err != nil {
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, ErorrTelegram(fmt.Sprintf("Veri okuma hatası: %s", err.Error()))))
			return
		}
	} else {
		filter := bson.D{
			{Key: "endDate", Value: bson.D{{Key: "$lt", Value: currentTime}}},
		}

		data, err := db.GetList(Database.Url, filter)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, ErorrTelegram("Veri tabanı hatası oluştu. Lütfen daha sonra tekrar deneyin.")))
			return
		}

		// URL'leri almak için bir döngü kullan
		for data.Next(context.TODO()) {
			var url Models.Url
			err := data.Decode(&url)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(message.Chat.ID, ErorrTelegram(fmt.Sprintf("Decode hatası oluştu: %s. Lütfen daha sonra tekrar deneyin.", err.Error()))))
				return
			}
			urls = append(urls, url)
		}

		if err := data.Err(); err != nil {
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, ErorrTelegram(fmt.Sprintf("Veri okuma hatası: %s", err.Error()))))
			return
		}
	}

	// Butonları oluştur

	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, u := range urls {
		// Buton verisi olarak URL'nin ID'sini kullan
		button := tgbotapi.NewInlineKeyboardButtonData(u.ShortUrl, fmt.Sprintf("%s:%s", u.Id.Hex(), u.ShortUrl))
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{button})
	}

	if len(urls) > 0 {
		replyMarkup := tgbotapi.NewInlineKeyboardMarkup(keyboard...)
		msg := tgbotapi.NewMessage(message.Chat.ID, mesaj)
		msg.ReplyMarkup = replyMarkup
		bot.Send(msg)
	} else {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Kısaltılmış URL'niz bulunmamaktadır."))
	}
}

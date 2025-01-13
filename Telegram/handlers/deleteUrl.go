package handlers

import (
	"fmt"

	"github.com/SadikSunbul/TelegramUrlBot/Database"
	"github.com/SadikSunbul/TelegramUrlBot/Models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleDeleteUrl(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *Database.DataBase) {
	chatID := message.Chat.ID

	// Kullanıcının tüm URL'lerini getir
	urls, err := db.GetList(Database.Url, bson.D{{"userTelegramId", chatID}})
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, ErorrTelegram("URL'ler getirilirken bir hata oluştu.")))
		return
	}

	var keyboard [][]tgbotapi.InlineKeyboardButton
	for urls.Next(nil) {
		var url Models.Url
		if err := urls.Decode(&url); err != nil {
			continue
		}

		// Her URL için bir buton oluştur
		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("🗑️ %s", url.ShortUrl),
			fmt.Sprintf("delete_confirm:%s", url.Id.Hex()),
		)
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{button})
	}

	if len(keyboard) == 0 {
		bot.Send(tgbotapi.NewMessage(chatID, "Silinebilecek URL'niz bulunmamaktadır."))
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Silmek istediğiniz URL'yi seçin:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	bot.Send(msg)
}

package Telegram

import (
	"fmt"
	"strings"

	"github.com/SadikSunbul/TelegramUrlBot/Database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func handleCallbackQuery(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *Database.DataBase) {
	if update.CallbackQuery != nil {
		data := update.CallbackQuery.Data
		parts := strings.Split(data, ":")
		if len(parts) == 2 {
			urlId := parts[0]    // URL'nin ID'si
			shortUrl := parts[1] // Kısa URL

			// Kullanıcıya hangi verileri görmek istediğini sor
			keyboard := [][]tgbotapi.InlineKeyboardButton{
				{
					tgbotapi.NewInlineKeyboardButtonData("Toplam Kaç kişi tıkladı?", fmt.Sprintf("clicks:%s:%s", urlId, shortUrl)),
				},
				{
					tgbotapi.NewInlineKeyboardButtonData("Hangi ülkelerden tıklandı?", fmt.Sprintf("countries:%s:%s", urlId, shortUrl)),
				},
				{
					tgbotapi.NewInlineKeyboardButtonData("Ortalama tıklanma zamanları?", fmt.Sprintf("average_times:%s:%s", urlId, shortUrl)),
				},
				{
					tgbotapi.NewInlineKeyboardButtonData("Bu linkin uzun hali?", fmt.Sprintf("long_link:%s:%s", urlId, shortUrl)),
				},
				{
					tgbotapi.NewInlineKeyboardButtonData("Bu linkin bitiş zamanı?", fmt.Sprintf("end_date:%s:%s", urlId, shortUrl)),
				},
			}
			replyMarkup := tgbotapi.NewInlineKeyboardMarkup(keyboard...)
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("👀 (%s) Lütfen görmek istediğiniz veriyi seçin:", shortUrl))
			msg.ReplyMarkup = replyMarkup
			bot.Send(msg)
		} else if len(parts) == 3 {
			action := parts[0]
			urlId := parts[1]
			handleAction(action, urlId, update, bot, db)
		}
	}
}

func handleAction(action, urlId string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *Database.DataBase) {
	switch action {
	case "clicks":
		clicksCount := getClicksCount(urlId, db)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Bu URL'ye toplam %d kişi tıkladı.", clicksCount))
		bot.Send(msg)
	case "countries":
		countries := getCountries(urlId, db)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Bu URL'ye tıklayan ülkeler: %s", strings.Join(countries, ", ")))
		bot.Send(msg)
	case "average_times":
		averageTimes := getAverageClickTimes(urlId, db)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Bu URL için ortalama tıklanma zamanları: %s", averageTimes))
		bot.Send(msg)
	case "long_link":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Bu URL'nin uzun hali: %s", getLongLink(urlId, db)))
		bot.Send(msg)
	case "end_date":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Bu URL bitiş zamanı: %s", getEndDate(urlId, db)))
		bot.Send(msg)
	}
}

func getClicksCount(urlId string, db *Database.DataBase) int {
	// TODO: Veritabanından tıklama sayısını al
	return 0 // Şimdilik 0 döndür
}

func getCountries(urlId string, db *Database.DataBase) []string {
	// TODO: Veritabanından ülke bilgilerini al
	return []string{"Türkiye"} // Şimdilik sabit değer döndür
}

func getAverageClickTimes(urlId string, db *Database.DataBase) string {
	// TODO: Veritabanından ortalama tıklanma zamanlarını al
	return "Henüz veri yok" // Şimdilik sabit değer döndür
}

func getLongLink(urlId string, db *Database.DataBase) string {
	return "...x.com"
}

func getEndDate(urlId string, db *Database.DataBase) string {
	return "sonsuz.."
}

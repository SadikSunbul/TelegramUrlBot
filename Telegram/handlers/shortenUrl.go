package handlers

import (
	"github.com/SadikSunbul/TelegramUrlBot/Database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
)

var UserData = make(map[int64]map[string]string)

func HandleShortenUrl(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *Database.DataBase) {
	chatID := message.Chat.ID
	_, err := db.GetBy(Database.User, bson.D{{"telegramId", chatID}})
	if err != nil {
		// Eğer hata "no documents in result" ise, kullanıcı yok demektir
		if err.Error() == "mongo: no documents in result" {
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, ErorrTelegram("Kaydınız bulunmamaktadır lütfen /start diyerek kaydınızı yaptıktan sonra tekrar deneyiniz.")))
			return
		}
	}
	// Kullanıcıdan URL'yi iste
	bot.Send(tgbotapi.NewMessage(chatID, "🔗 Url kısaltma işlemine hoşgeldiniz. Kısaltmak istediğiniz URL'yi giriniz:"))

	// Kullanıcının girdiği URL'yi al
	UserData[chatID] = make(map[string]string)
	UserData[chatID]["step"] = "get_url"
	StartTimer(chatID, bot, UserData)
}

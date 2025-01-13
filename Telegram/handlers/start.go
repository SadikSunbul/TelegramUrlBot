package handlers

import (
	"fmt"

	"github.com/SadikSunbul/TelegramUrlBot/Database"
	"github.com/SadikSunbul/TelegramUrlBot/Models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *Database.DataBase) {
	telegramID := message.Chat.ID
	data, err := db.GetBy(Database.User, bson.D{{"telegramId", telegramID}})
	if err != nil {
		// Eğer hata "no documents in result" ise, kullanıcı yok demektir
		if err.Error() == "mongo: no documents in result" {
			// Kullanıcı kaydı yapılmamış demek
			user := Models.User{
				TelegramId: telegramID,
				Name:       message.Chat.UserName,
			}
			_, err := db.Add(Database.User, user)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(message.Chat.ID, ErorrTelegram("Veri tabanı hatası oluştu. Lütfen daha sonra tekrar deneyin.")))
				return
			}
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, SuccessfulTelegram("Kaydınız başarılı bir şekilde yapıldı.")))
			return
		}

		// Diğer hatalar için hata mesajı gönder
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, ErorrTelegram(fmt.Sprintf("Veri tabanı hatası oluştu. Lütfen daha sonra tekrar deneyin.: %s", err.Error()))))
		return
	}

	var user Models.User
	if err := data.Decode(&user); err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, ErorrTelegram("Veri tabanı hatası oluştu. Lütfen daha sonra tekrar deneyin.")))
		return
	}

	// Kullanıcı zaten kayıtlı
	bot.Send(tgbotapi.NewMessage(message.Chat.ID, NotificationTelegram(fmt.Sprintf("Hoşgeldin %s , daha fazla bilgi için ' /help ' yaz !", user.Name))))
}

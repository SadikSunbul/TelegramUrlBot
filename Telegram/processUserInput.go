package Telegram

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"time"

	"github.com/SadikSunbul/TelegramUrlBot/Database"
	"github.com/SadikSunbul/TelegramUrlBot/Models"
	"github.com/SadikSunbul/TelegramUrlBot/Telegram/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/rand"
)

func ProcessUserInput(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *Database.DataBase) {
	chatID := update.Message.Chat.ID

	if data, ok := handlers.UserData[chatID]; ok {
		switch data["step"] {
		case "get_url":
			handlers.UserData[chatID]["originalUrl"] = update.Message.Text
			msg := tgbotapi.NewMessage(chatID, "Otomatik bir ad atansın mı?")
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("✅ Evet", "auto_name:yes"),
					tgbotapi.NewInlineKeyboardButtonData("❌ Hayır", "auto_name:no"),
				),
			)
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
			handlers.UserData[chatID]["step"] = "ask_auto_name"

		case "ask_auto_name":
			if update.Message.Text == "evet" {
				var shortUrl string

				for {
					shortUrl = generateShortUrl() // Kısa URL oluşturma fonksiyonu
					available, err := isShortUrlAvailable(db, shortUrl)
					if err != nil {
						bot.Send(tgbotapi.NewMessage(chatID, handlers.ErorrTelegram(fmt.Sprintf("Veri tabanı hatası oluştu. Lütfen daha sonra tekrar deneyin.:", err.Error()))))
						delete(handlers.UserData, chatID)
						return
					}
					if available {
						break
					}
				}

				handlers.UserData[chatID]["shortUrl"] = shortUrl
				bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Kısa URL'niz: %s", shortUrl)))
				msg := tgbotapi.NewMessage(chatID, "Zamanlı mı yoksa zamansız mı olsun?")
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("⏱️ Zamanlı", "time_limit:timed"),
						tgbotapi.NewInlineKeyboardButtonData("♾️ Zamansız", "time_limit:unlimited"),
					),
				)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
				handlers.UserData[chatID]["step"] = "ask_time_limit"
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Kısa URL'nizi giriniz (sadeceşu kısmı girin ....com/'burayı'):"))
				handlers.UserData[chatID]["step"] = "get_custom_short_url"
			}

		case "get_custom_short_url":
			shortUrl := update.Message.Text
			available, err := isShortUrlAvailable(db, shortUrl)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, handlers.ErorrTelegram(fmt.Sprintf("Veri tabanı hatası oluştu. Lütfen daha sonra tekrar deneyin.:", err.Error()))))
				return
			}
			if !available {
				bot.Send(tgbotapi.NewMessage(chatID, "Bu URL zaten kullanılıyor. Lütfen başka bir kısa URL girin:"))
				handlers.UserData[chatID]["step"] = "get_custom_short_url" // Kullanıcıdan yeni URL girmesini iste
				return
			}
			handlers.UserData[chatID]["shortUrl"] = shortUrl
			msg := tgbotapi.NewMessage(chatID, "Zamanlı mı yoksa zamansız mı olsun?")
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("⏱️ Zamanlı", "time_limit:timed"),
					tgbotapi.NewInlineKeyboardButtonData("♾️ Zamansız", "time_limit:unlimited"),
				),
			)
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
			handlers.UserData[chatID]["step"] = "ask_time_limit"

		case "ask_time_limit":
			if update.Message.Text == "zamanlı" {
				bot.Send(tgbotapi.NewMessage(chatID, "Lütfen son kullanma tarihini giriniz (YYYY-MM-DD HH:MM):"))
				handlers.UserData[chatID]["step"] = "get_expiration_date"
			} else {
				// Sınırsız kullanım
				err := saveUrlToDatabase(handlers.UserData[chatID], db, chatID) // Veritabanına kaydet
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, handlers.ErorrTelegram("Veri tabanı hatası oluştu. Lütfen daha sonra tekrar deneyin.")))
					delete(handlers.UserData, chatID) // Kullanıcı verisini temizle
					return
				}
				bot.Send(tgbotapi.NewMessage(chatID, "Kısa URL başarıyla oluşturuldu ve sınırsız olarak kullanılabilir."))
				delete(handlers.UserData, chatID) // Kullanıcı verisini temizle
			}

		case "get_expiration_date":
			// Kullanıcının girdiği son kullanma tarihini al
			expirationDate := update.Message.Text
			handlers.UserData[chatID]["expirationDate"] = expirationDate
			err := saveUrlToDatabase(handlers.UserData[chatID], db, chatID) // Veritabanına kaydet
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, handlers.ErorrTelegram("Veri tabanı hatası oluştu. Lütfen daha sonra tekrar deneyin.")))
				return
			}
			bot.Send(tgbotapi.NewMessage(chatID, "Kısa URL başarıyla oluşturuldu."))
			delete(handlers.UserData, chatID) // Kullanıcı verisini temizle
		}
	}
}

func generateShortUrl() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	shortUrl := ""

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		shortUrl += string(charset[randomIndex])
	}

	return shortUrl
}
func isShortUrlAvailable(db *Database.DataBase, shortUrl string) (bool, error) {
	data, err := db.GetBy(Database.Url, bson.D{
		{Key: "shortUrl", Value: shortUrl},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "endDate", Value: bson.D{{Key: "$gt", Value: time.Now()}}}},
			bson.D{{Key: "endDate", Value: bson.D{{Key: "$exists", Value: false}}}},
		}},
	})
	if err != nil {
		// Eğer hata "no documents in result" ise, bunu göz ardı et
		if err.Error() == "mongo: no documents in result" {
			return true, nil // URL kullanılabilir
		}
		return false, err // Diğer hataları döndür
	}
	var url Models.Url
	data.Decode(&url)
	fmt.Sprintf("data:", data)
	return data == nil, nil // Eğer data nil ise, shortUrl kullanılabilir
}
func saveUrlToDatabase(data map[string]string, db *Database.DataBase, chanId int64) error {
	var url Models.Url
	url.OriginalUrl = data["originalUrl"]
	url.ShortUrl = data["shortUrl"]
	url.UserTelegramId = chanId

	expTimeStr := data["expirationDate"]
	if expTimeStr == "" {
		// sınırsız zaman
	} else {
		expirationTime, err := time.Parse("2006-01-02 15:04", expTimeStr)
		if err != nil {
			return err
		}
		url.EndDate = primitive.NewDateTimeFromTime(expirationTime)
	}

	_, err := db.Add(Database.Url, url)
	return err
}

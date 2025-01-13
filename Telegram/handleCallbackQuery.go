package Telegram

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	chart "github.com/SadikSunbul/TelegramUrlBot/Chart"
	"github.com/SadikSunbul/TelegramUrlBot/Database"
	"github.com/SadikSunbul/TelegramUrlBot/Models"
	"github.com/SadikSunbul/TelegramUrlBot/Telegram/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func handleCallbackQuery(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *Database.DataBase) {
	if update.CallbackQuery != nil {
		data := update.CallbackQuery.Data
		chatID := update.CallbackQuery.Message.Chat.ID
		parts := strings.Split(data, ":")

		if len(parts) == 2 {
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
					bot.Send(tgbotapi.NewMessage(chatID, "Kısa URL'nizi giriniz (Örnek: 'sadik' yazarsanız -> kısaurl.com/sadik şeklinde olacaktır):"))
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
			default:
				// URL detayları için
				urlId := parts[0]    // URL'nin ID'si
				shortUrl := parts[1] // Kısa URL

				// URL'nin kullanıcıya ait olup olmadığını kontrol et
				objID, err := primitive.ObjectIDFromHex(urlId)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, handlers.ErorrTelegram("Geçersiz URL ID'si.")))
					return
				}

				var url Models.Url
				result, err := db.GetBy(Database.Url, bson.D{
					{"_id", objID},
					{"userTelegramId", chatID},
				})
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, handlers.ErorrTelegram("Bu URL'ye erişim yetkiniz yok veya URL bulunamadı.")))
					return
				}
				if err := result.Decode(&url); err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, handlers.ErorrTelegram("Bu URL'ye erişim yetkiniz yok veya URL bulunamadı.")))
					return
				}

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
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("👀 (%s) Lütfen görmek istediğiniz veriyi seçin:", shortUrl))
				msg.ReplyMarkup = replyMarkup
				bot.Send(msg)
			}
		} else if len(parts) == 3 {
			action := parts[0]
			urlId := parts[1]

			// URL'nin kullanıcıya ait olup olmadığını kontrol et
			objID, err := primitive.ObjectIDFromHex(urlId)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, handlers.ErorrTelegram("Geçersiz URL ID'si.")))
				return
			}

			var url Models.Url
			result, err := db.GetBy(Database.Url, bson.D{
				{"_id", objID},
				{"userTelegramId", chatID},
			})
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, handlers.ErorrTelegram("Bu URL'ye erişim yetkiniz yok veya URL bulunamadı.")))
				return
			}
			if err := result.Decode(&url); err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, handlers.ErorrTelegram("Bu URL'ye erişim yetkiniz yok veya URL bulunamadı.")))
				return
			}

			handleAction(action, urlId, update, bot, db)
		}

		// Callback query'yi yanıtla
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		bot.AnswerCallbackQuery(callback)
	}
}

func handleAction(action, urlId string, update tgbotapi.Update, bot *tgbotapi.BotAPI, db *Database.DataBase) {
	switch action {
	case "clicks":
		err := CreateChart(update, bot, []string{"2025"}, []int{120}, "Kişi") // TODO : veriler değişicek
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, handlers.ErorrTelegram("Grafik oluşturulurken hata oluştu.")))
			return
		}
	case "countries":
		err := CreateChart(update, bot, []string{"Türkiye", "Almanya", "ABD", "Suriye"}, []int{120, 12, 78, 99}, "Ülkeler") // TODO : veriler değişicek
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, handlers.ErorrTelegram("Grafik oluşturulurken hata oluştu.")))
			return
		}
	case "average_times":
		err := CreateChart(update, bot, []string{"09:00", "20:11", "00:01"}, []int{120, 50, 20}, "Top 3 tıklanma zamanı") // TODO : veriler değişicek
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, handlers.ErorrTelegram("Grafik oluşturulurken hata oluştu.")))
			return
		}
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

func CreateChart(update tgbotapi.Update, bot *tgbotapi.BotAPI, xExsenData []string, yEksenData []int, title string) error {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(title, "grafiğini hazırlıyorum..."))
	bot.Send(msg)

	// Grafiği oluştur
	chartBuffer := chart.CreateChart(xExsenData, yEksenData, title)

	file := tgbotapi.FileBytes{
		Name:  "chart.html",
		Bytes: chartBuffer.Bytes(),
	}

	doc := tgbotapi.NewDocumentUpload(update.CallbackQuery.Message.Chat.ID, file)
	doc.Caption = "İşte grafiğiniz!"
	_, err := bot.Send(doc)
	return err
}

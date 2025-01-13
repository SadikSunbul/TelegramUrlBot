package Telegram

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/SadikSunbul/TelegramUrlBot/Database"
	"github.com/SadikSunbul/TelegramUrlBot/Models"
	"github.com/SadikSunbul/TelegramUrlBot/Telegram/analysis"
	"github.com/SadikSunbul/TelegramUrlBot/Telegram/handlers"
	"github.com/SadikSunbul/TelegramUrlBot/config"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
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
					cfg := config.GetConfig()
					bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("%s%s URL başarıyla oluşturuldu ve sınırsız olarak kullanılabilir.", cfg.ApiDomain, handlers.UserData[chatID]["shortUrl"])))
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
					cfg := config.GetConfig()
					bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("%s%s URL başarıyla oluşturuldu ve %d saat geçerli olacak.", cfg.ApiDomain, handlers.UserData[chatID]["shortUrl"], hours)))
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
						tgbotapi.NewInlineKeyboardButtonData("Toplam Kaç kere tıkladı?", fmt.Sprintf("clicks:%s:%s", urlId, shortUrl)),
					},
					{
						tgbotapi.NewInlineKeyboardButtonData("Hangi ülkelerden tıklandı?", fmt.Sprintf("countries:%s:%s", urlId, shortUrl)),
					},
					{
						tgbotapi.NewInlineKeyboardButtonData("Tıklama Analizi?", fmt.Sprintf("average_times:%s:%s", urlId, shortUrl)),
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
	case "clicks": // ✅

		datas, err := analysis.GetUrlInfo(db, urlId)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, handlers.ErorrTelegram("Veri tabanı hatası oluştu.")))
			return
		}
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("👆🏻 %v defa bu url ye tıklandı.", len(datas))))
		return

	case "countries": // ✅

		ulke, sayisi, err := analysis.CountriesAnalyse(db, urlId)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, handlers.ErorrTelegram("Veri tabanı hatası oluştu.")))
			return
		}

		err = CreateChart(update, bot, ulke, sayisi, "🇹🇷 Ülkeler")
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, handlers.ErorrTelegram("Grafik oluşturulurken hata oluştu.")))
			return
		}
	case "average_times": // ✅
		// Tıklanma zamanalrı
		saat, sayi, err := analysis.TimeAnalysis(db, urlId)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, handlers.ErorrTelegram("Veri tabanı hatası oluştu.")))
			return
		}
		err = CreateChart(update, bot, saat, sayi, "📈 Tıklama Analizi")
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, handlers.ErorrTelegram("Grafik oluşturulurken hata oluştu.")))
			return
		}
	case "long_link": // ✅

		urldate, err := db.Get(Database.Url, urlId)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, handlers.ErorrTelegram("Veri tabanı hatası oluştu.")))
			return
		}
		var urldecode Models.Url
		err = urldate.Decode(&urldecode)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, handlers.ErorrTelegram("Decode hatası oluştu.")))
			return
		}
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Bu URL'nin uzun hali: %s", urldecode.OriginalUrl))
		bot.Send(msg)
	case "end_date": // ✅
		urldate, err := db.Get(Database.Url, urlId)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, handlers.ErorrTelegram("Veri tabanı hatası oluştu.")))
			return
		}
		var urldecode Models.Url
		err = urldate.Decode(&urldecode)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, handlers.ErorrTelegram("Decode hatası oluştu.")))
			return
		}
		if urldecode.EndDate == 0 {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Bu URL bitiş zamanı: ∞ (sonsuz)")
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Bu URL bitiş zamanı: %s", urldecode.EndDate.Time().Format("2006-01-02 15:04")))
		bot.Send(msg)
	}
}

func generateBarItems(values []int) []opts.BarData {
	items := make([]opts.BarData, len(values))
	for i := 0; i < len(values); i++ {
		items[i] = opts.BarData{Value: values[i]}
	}
	return items
}

func CreateChartDik(xExsenData []string, yEksenData []int, title string) *bytes.Buffer {
	if len(xExsenData) != len(yEksenData) {
		// hatalı
		return nil
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: title}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{Rotate: 90},
		}),
	)
	bar.SetXAxis(xExsenData).
		AddSeries("Değerler", generateBarItems(yEksenData))

	buf := new(bytes.Buffer)
	bar.Render(buf)
	return buf
}

func CreateChart(update tgbotapi.Update, bot *tgbotapi.BotAPI, xExsenData []string, yEksenData []int, title string) error {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("%s grafiğini hazırlıyorum...", title))
	bot.Send(msg)

	// Grafiği oluştur
	chartBuffer := CreateChartDik(xExsenData, yEksenData, title)
	if chartBuffer == nil {
		return fmt.Errorf("grafik oluşturulamadı")
	}

	file := tgbotapi.FileBytes{
		Name:  "chart.html",
		Bytes: chartBuffer.Bytes(),
	}

	doc := tgbotapi.NewDocumentUpload(update.CallbackQuery.Message.Chat.ID, file)
	doc.Caption = fmt.Sprintf("%s grafiğiniz", title)
	_, err := bot.Send(doc)
	return err
}

package handlers

import (
	"time"

	"github.com/SadikSunbul/TelegramUrlBot/Database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleClear(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *Database.DataBase) {
	// Önce bir bilgi mesajı gönder
	msg := tgbotapi.NewMessage(message.Chat.ID, "🗑️ Mesajlar temizleniyor...")
	statusMsg, _ := bot.Send(msg)

	// Mesajları sil
	for i := 0; i < 100; i++ { // Daha fazla mesaj için döngüyü artır
		if i%50 == 0 { // Her 50 silme işleminde bir bekleme yap
			time.Sleep(500 * time.Millisecond) // Rate limit'e takılmamak için
		}

		// Sadece aşağı doğru mesajları sil, en üstteki mesajı koru
		if message.MessageID-i > 1 { // MessageID 1'den büyükse sil
			deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID-i)
			bot.Send(deleteMsg)
		}
	}

	// Silme işlemi bitti mesajı
	completeMsg := tgbotapi.NewMessage(message.Chat.ID, "🗑️ Mesajlar temizlendi!")
	bot.Send(completeMsg)

	// 3 saniye bekle ve bilgi mesajlarını sil
	time.Sleep(3 * time.Second)
	bot.Send(tgbotapi.NewDeleteMessage(message.Chat.ID, statusMsg.MessageID))
}

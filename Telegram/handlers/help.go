package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func HandleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msgTxt := "💡 *Yardım Menüsüne Hoş Geldiniz* \n\n" +
		"*Temel Komutlar:*\n" +
		"➡️ /start - Botu başlatır ve kullanıcı kaydınızı yapar\n" +
		"➡️ /help - Bu yardım menüsünü gösterir\n\n" +
		"*URL İşlemleri:*\n" +
		"➡️ /shortenurl - Yeni bir URL kısaltmak için kullanılır\n" +
		"➡️ /mylinksactive - Aktif olan kısa URL'lerinizi listeler\n" +
		"➡️ /mylinkspassive - Süresi dolmuş URL'lerinizi listeler\n\n" +
		"*Diğer Komutlar:*\n" +
		"➡️ /clear - Sohbetteki mesajları temizler\n\n" +
		"*URL Özellikleri:*\n" +
		"• Zamanlı veya zamansız URL oluşturabilirsiniz\n" +
		"• Özel kısa URL belirleyebilirsiniz\n" +
		"• URL'lerinizin tıklanma istatistiklerini görebilirsiniz\n" +
		"• Hangi ülkelerden tıklandığını takip edebilirsiniz\n" +
		"• Tıklanma zamanlarını analiz edebilirsiniz"

	msg := tgbotapi.NewMessage(message.Chat.ID, msgTxt)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

/*
Komut Listesi ve Açıklamaları:

/start - Telegram botunu başlatır ve kullanıcı kaydını oluşturur
/help - Kullanılabilir komutları ve özelliklerini listeler
/shortenurl - Uzun URL'leri kısaltmak için kullanılır
/mylinksactive - Süresi devam eden aktif URL'leri gösterir
/mylinkspassive - Süresi dolmuş pasif URL'leri gösterir
/clear - Sohbet geçmişini temizler (en üstteki mesaj hariç)

URL Detay Komutları (Buton olarak gelir):
clicks - URL'nin toplam tıklanma sayısını gösterir
countries - URL'ye hangi ülkelerden erişildiğini gösterir
average_times - URL'nin tıklanma zamanlarının analizini gösterir
long_link - Kısa URL'nin uzun/orijinal halini gösterir
end_date - URL'nin bitiş tarihini gösterir (varsa)
*/

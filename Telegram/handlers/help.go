package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func HandleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msgTxt := "💡 *Commands:* \n\n" + "help sayfasına hoş geldin..."

	msg := tgbotapi.NewMessage(message.Chat.ID, msgTxt)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

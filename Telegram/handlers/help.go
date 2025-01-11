package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func HandleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msgTxt := "💡 *Help sayfasına hoş geldiniz:* \n\n"

	msg := tgbotapi.NewMessage(message.Chat.ID, msgTxt)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

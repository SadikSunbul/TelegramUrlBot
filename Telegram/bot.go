package Telegram

import (
	"github.com/SadikSunbul/TelegramUrlBot/Database"
	"github.com/SadikSunbul/TelegramUrlBot/Telegram/handlers"
	"log"

	"github.com/SadikSunbul/TelegramUrlBot/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ConnectTelegram(db *Database.DataBase) {
	config := *config.GetConfig()

	bot, err := tgbotapi.NewBotAPI(config.BootIdTelegram)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Error getting updates: %v", err)
	}
	for update := range updates {
		go handleUpdate(update, bot, db)
	}
}

func handleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *Database.DataBase) {

	if update.Message == nil {
		return
	}

	switch update.Message.Text {
	case "/help":
		handlers.HandleHelp(bot, update.Message)
	case "/start":
		handlers.HandleStart(bot, update.Message, db)
	case "/shortenurl":
		handlers.HandleShortenUrl(bot, update.Message, db)

	default:
		ProcessUserInput(update, bot, db) // Kullanıcıdan gelen mesajı işle
	}
}

package main

import (
	"github.com/SadikSunbul/TelegramUrlBot/Database"
	"github.com/SadikSunbul/TelegramUrlBot/Telegram"
	"github.com/SadikSunbul/TelegramUrlBot/config"
)

func main() {
	config.LoadConfig("config.yaml")
	db := Database.ConnectionDatabase()
	Telegram.ConnectTelegram(db)
}

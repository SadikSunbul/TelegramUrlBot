package main

import (
	"fmt"

	"github.com/SadikSunbul/TelegramUrlBot/config"
)

func main() {
	config := config.LoadConfig() // Yapılandırmayı yükle
	fmt.Printf(config.BootIdTelegram)
}

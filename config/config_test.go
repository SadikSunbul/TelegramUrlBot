package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test için geçici bir config.yaml dosyası oluştur
	file, err := os.Create("configtest.yaml")
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	// defer os.Remove("configtest.yaml") // Test sonunda dosyayı sil

	// Örnek içerik yaz
	file.WriteString("mongoDbConnect: \"mongodb://localhost:27017\"\n")
	file.WriteString("bootIdTelegram: \"your-telegram-bot-token\"\n")
	file.WriteString("dbName: \"telegram\"")
	file.Close()

	// Yapılandırmayı yükle
	cfg := LoadConfig("configtest.yaml")

	// Beklenen değerlerle karşılaştır
	if cfg.MongoDbConnect != "mongodb://localhost:27017" {
		t.Errorf("Expected mongoDbConnect to be 'mongodb://localhost:27017', got '%s'", cfg.MongoDbConnect)
	}
	if cfg.BootIdTelegram != "your-telegram-bot-token" {
		t.Errorf("Expected bootIdTelegram to be 'your-telegram-bot-token', got '%s'", cfg.BootIdTelegram)
	}
	if cfg.DbName != "telegram" {
		t.Errorf("Expected dbName to be 'telegram', got '%s'", cfg.DbName)
	}
}

func TestGetConfig(t *testing.T) {
	// İlk önce LoadConfig fonksiyonunu çağır
	LoadConfig("configtest.yaml")

	// GetConfig ile yapılandırmayı al
	cfg := GetConfig()

	// Beklenen değerlerle karşılaştır
	if cfg.MongoDbConnect == "" {
		t.Error("Expected MongoDbConnect to be set, but it is empty")
	}
	if cfg.BootIdTelegram == "" {
		t.Error("Expected BootIdTelegram to be set, but it is empty")
	}
}

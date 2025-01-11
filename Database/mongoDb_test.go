package Database

import (
	"testing"

	"github.com/SadikSunbul/TelegramUrlBot/config"
)

func TestConnectionDatabase(t *testing.T) {
	// Test için config yükle
	cfg := config.LoadConfig("configtest.yaml")
	cfg.MongoDbConnect = "mongodb://localhost:27017"
	cfg.DbName = "test_db"

	// Veritabanı bağlantısını test et
	db := ConnectionDatabase()
	if db.Client == nil {
		t.Fatal("Veritabanı bağlantısı başarısız")
	}

	// Ping ile bağlantıyı test et
	if err := db.Client.RunCommand(nil, map[string]interface{}{"ping": 1}).Err(); err != nil {
		t.Errorf("Veritabanı ping hatası: %v", err)
	}
}

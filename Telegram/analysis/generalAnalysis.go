package analysis

import (
	"context"
	"sort"
	"time"

	"github.com/SadikSunbul/TelegramUrlBot/Database"
	"github.com/SadikSunbul/TelegramUrlBot/Models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUrlInfo(db *Database.DataBase, urlId string) ([]Models.UserDeviceInfo, error) {
	objID, err := primitive.ObjectIDFromHex(urlId)
	if err != nil {
		return nil, err
	}
	data, err := db.GetList(Database.UrlIfo, bson.D{{"urlId", objID}})
	if err != nil {
		return nil, err
	}

	var urlInfos []Models.UserDeviceInfo
	for data.Next(context.TODO()) {
		var urlInfo Models.UserDeviceInfo
		if err := data.Decode(&urlInfo); err != nil {
			return nil, err
		}
		urlInfos = append(urlInfos, urlInfo)
	}

	return urlInfos, nil
}

func TimeAnalysis(db *Database.DataBase, urlId string) ([]string, []int, error) {
	urlInfos, err := GetUrlInfo(db, urlId)
	if err != nil {
		return nil, nil, err
	}

	// Zamanları 15 dakikalık dilimlere yuvarlama ve sayma
	timeCount := make(map[string]int)
	for _, info := range urlInfos {
		// MongoDB timestamp'ini time.Time'a çevir
		t := info.ClickTime.Time()

		// 15 dakikalık dilime yuvarla
		minute := t.Minute()
		roundedMinute := (minute / 15) * 15
		roundedTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), roundedMinute, 0, 0, t.Location())

		// Formatla ve say
		timeKey := roundedTime.Format("15:04")
		timeCount[timeKey]++
	}

	// Map'i slice'a çevir ve zamanları sırala
	type timeEntry struct {
		time  string
		count int
	}
	var entries []timeEntry
	for t, count := range timeCount {
		entries = append(entries, timeEntry{t, count})
	}
	// Zamanları kronolojik olarak sırala
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].time < entries[j].time
	})

	// Tüm zamanları al
	var times []string
	var counts []int
	for _, entry := range entries {
		times = append(times, entry.time)
		counts = append(counts, entry.count)
	}

	return times, counts, nil
}

func CountriesAnalyse(db *Database.DataBase, urlId string) ([]string, []int, error) {
	data, err := GetUrlInfo(db, urlId)
	if err != nil {
		return nil, nil, err
	}

	// Ülke sayımını tutacak harita
	countryCount := make(map[string]int)

	// Her kullanıcı bilgisi için ülke sayısını artır
	for _, userInfo := range data {
		countryCount[userInfo.Country]++
	}

	// Haritayı dilimlere çevir
	var countries []string
	var counts []int
	for country, count := range countryCount {
		countries = append(countries, country)
		counts = append(counts, count)
	}

	return countries, counts, nil
}

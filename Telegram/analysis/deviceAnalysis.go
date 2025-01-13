package analysis

import (
	"github.com/SadikSunbul/TelegramUrlBot/Database"
)

func DeviceAnalysis(db *Database.DataBase, urlId string) ([]string, []int, error) {
	data, err := GetUrlInfo(db, urlId)
	if err != nil {
		return nil, nil, err
	}

	deviceCount := make(map[string]int)
	for _, info := range data {
		deviceCount[info.Device]++
	}

	var devices []string
	var counts []int
	for device, count := range deviceCount {
		devices = append(devices, device)
		counts = append(counts, count)
	}
	return devices, counts, nil
}

func BrowserAnalysis(db *Database.DataBase, urlId string) ([]string, []int, error) {
	data, err := GetUrlInfo(db, urlId)
	if err != nil {
		return nil, nil, err
	}

	browserCount := make(map[string]int)
	for _, info := range data {
		browserCount[info.Browser]++
	}

	var browsers []string
	var counts []int
	for browser, count := range browserCount {
		browsers = append(browsers, browser)
		counts = append(counts, count)
	}
	return browsers, counts, nil
}

func OsAnalysis(db *Database.DataBase, urlId string) ([]string, []int, error) {
	data, err := GetUrlInfo(db, urlId)
	if err != nil {
		return nil, nil, err
	}

	osCount := make(map[string]int)
	for _, info := range data {
		osCount[info.OS]++
	}

	var systems []string
	var counts []int
	for os, count := range osCount {
		systems = append(systems, os)
		counts = append(counts, count)
	}
	return systems, counts, nil
}

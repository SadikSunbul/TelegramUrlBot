package handlers

func ErorrTelegram(mesaj string) string {
	str := "⚠️ "
	str += mesaj
	return str
}

func NotificationTelegram(mesaj string) string {
	str := "🔔 "
	str += mesaj
	return str
}

func SuccessfulTelegram(mesaj string) string {
	str := "✅ "
	str += mesaj
	return str
}

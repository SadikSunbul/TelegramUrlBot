[![Go CI/CD](https://github.com/SadikSunbul/TelegramUrlBot/actions/workflows/go.yml/badge.svg)](https://github.com/SadikSunbul/TelegramUrlBot/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/SadikSunbul/TelegramUrlBot)](https://goreportcard.com/report/github.com/SadikSunbul/TelegramUrlBot)
[![Version](https://img.shields.io/badge/Version-1.0-blue)]()
# Telegram Bot Üzerinden URL Kısaltma ve Analiz Hizmeti
[TelegramUrlBotServer](https://github.com/SadikSunbul/TelegramUrlBotServer) dan kısaltılmış url lere ulaşılabilir.
## Proje Özeti
Bu projenin amacı, URL kısaltma hizmetleri sunan ve kısaltılmış linkler için detaylı analizler sağlayan bir Telegram botu geliştirmektir. Geleneksel web tabanlı URL kısaltma platformlarından farklı olarak, bu proje Telegram'ı birincil kullanıcı arayüzü olarak kullanmayı hedeflemektedir. Bu sayede işlemler mesajlaşma uygulaması içerisinde kolayca gerçekleştirilebilecektir.

## Proje Amaçları
- Telegram botu aracılığıyla uzun URL'lerin kısaltılması.
- Kullanıcıların kısaltılan URL'ler için son kullanma tarihi belirleyebilmesi. Süresi dolan URL'ler orijinal linke yönlendirme yapmayacaktır.
- Her bir kısaltılan URL için detaylı analiz sağlanması:
  - Tıklama sayısı.
  - Kullanıcıların coğrafi konumu (ülkeler).
  - Zaman bazlı analiz (en yoğun kullanım saatleri).
- Kullanıcıların:
  - Aktif ve süresi dolmuş linklerini görebilmesi.
  - URL'lerine ait analiz verilerine erişebilmesi.
- Kullanıcıların aktif bir URL kullanmaması durumunda, özel bir "custom URL" belirleyebilmesi.

## Teknolojiler ve Araçlar
- **Programlama Dili:** Go (Golang)
- **Veritabanı:** MongoDB veya PostgreSQL (proje gereksinimlerine göre karar verilecek).
- **Telegram Bot API:** Kullanıcı etkileşimi ve komut yönetimi için.

## Kurulum
1. Bu projeyi klonlayın:
    ```bash
    git clone https://github.com/SadikSunbul/TelegramUrlBot.git
    cd TelegramUrlBot
    ```
2. Gereksinimleri yükleyin:
    ```bash
    go mod tidy
    ```
3. Telegram botunuzu oluşturun ve token'ınızı alın. BotFather ile Telegram üzerinden botunuzu oluşturabilirsiniz.

4. Proje ayarlarını yapılandırın:
    ```bash
    cp config.example.yaml config.yaml
    ```
    `config.yaml` dosyasını düzenleyerek Telegram bot token'ınızı ve veritabanı ayarlarınızı girin.

5. Projeyi başlatın:
    ```bash
    go run main.go
    ```

## Kullanım
Telegram botunuza `/start` komutunu göndererek botu başlatabilirsiniz. Botun sunduğu komutlar ve işlevler hakkında bilgi almak için `/help` komutunu kullanabilirsiniz.

## Katkıda Bulunma
Katkılarınızı memnuniyetle karşılıyoruz! Lütfen katkıda bulunmadan önce bir issue açarak ne üzerinde çalışmak istediğinizi belirtin.

1. Fork yapın.
2. Kendi branşınızı oluşturun (`git checkout -b feature/AmazingFeature`).
3. Değişikliklerinizi commitleyin (`git commit -m 'Add some AmazingFeature'`).
4. Branşınıza push yapın (`git push origin feature/AmazingFeature`).
5. Bir Pull Request açın.

## Lisans
Bu proje MIT Lisansı ile lisanslanmıştır. Daha fazla bilgi için `LICENSE` dosyasına bakın.

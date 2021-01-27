package utils

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"time"
)

func Handlers(){
		var videoID,videoTitle string
		bot, err := tb.NewBot(tb.Settings{
			Token: os.Getenv("BOT_TOKEN"),
			Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		})
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println("Bot is running")
		//**********************************************
		//Mesaj butonları
		selector := &tb.ReplyMarkup{}
		btnmp3 := selector.Data("MP3", "btnmp3")
		btnmp4 := selector.Data("MP4", "btnmp4")

		selector.Inline(
			selector.Row(btnmp3, btnmp4),
		)
		//**********************************************

		bot.Handle("/start", func(m *tb.Message) {
			bot.Send(m.Sender, "İndirmek istediğiniz YouTube linkini gönderin")
		})

		bot.Handle("/help", func(m *tb.Message) {
			bot.Send(m.Sender, "Youtube2MP3 botuna hoşgeldiniz :)\n"+
				"İndirmek istediğiniz YouTube videosunun linkini mesaj olarak göndermeniz yeterlidir.\n"+
				"Katkıda bulunmak için : github.com/msrexe/yt2mp3-telegram-bot\n"+
				"@msrexe")
		})


		bot.Handle(tb.OnText, func(m *tb.Message) {
			bot.Send(m.Sender, "İsteğiniz işleniyor...")
			videoID,videoTitle, err = UrlFormatter(m.Text)
			if err == nil {
				bot.Send(m.Sender, "Video Adı : "+ videoTitle+"\nİndirmek istediğiniz formatı seçin", selector)
			} else {
				log.Println(err)
				bot.Send(m.Sender, "Video bulunamadı. Kaldırılmış veya hatalı URL olabilir!!")
				return
			}
		})


		bot.Handle(&btnmp4, func(c *tb.Callback) {
			bot.Send(c.Sender, "Medya hazırlanıyor. Bu işlem biraz zaman alabilir. Hazır olduğunda size mesaj ile bildirilecektir....")
			mp4, err := DownloadMP4fromURL(videoTitle,videoID)
			if err != nil {
				bot.Send(c.Sender, "Bir hata oluştu")
				log.Println(err)
				return
			}
			log.Println("Gönderim başladı")
			_, err = bot.Send(c.Sender, mp4)
			os.Remove(videoTitle+".mp4")
			if err != nil {
				bot.Send(c.Sender, "Bir hata oluştu")
				log.Println(err)
				return
			}
		})


		bot.Handle(&btnmp3, func(c *tb.Callback) {
			bot.Send(c.Sender, "Medya hazırlanıyor. Bu işlem biraz zaman alabilir. Hazır olduğunda size mesaj ile bildirilecektir....")

			mp3, err := DownloadMP3fromURL(videoTitle,videoID)
			if err != nil {
				bot.Send(c.Sender, "Bir hata oluştu")
				log.Println(err)
				return
			}
			log.Println("Gönderim başladı")
			_, err = bot.Send(c.Sender, mp3)
			os.Remove(videoTitle+".mp3")
			if err != nil {
				bot.Send(c.Sender, "Bir hata oluştu")
				log.Println(err)
				return
			}

		})

		bot.Start()
}
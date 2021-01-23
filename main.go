package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	SendMP3()
}

func SendMP3(){
	bot, err := tb.NewBot(tb.Settings{
		Token: "1514859884:AAEQpQkXr1iMAcH6BvQ70Q_-RMHREP9y79A",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Bot is running")

	bot.Handle("/start", func(m *tb.Message) {
		bot.Send(m.Sender,"Youtube2MP3 botuna hoşgeldiniz :)\n" +
			"İndirmek istediğiniz YouTube videosunun linkini mesaj olarak göndermeniz yeterlidir.\n" +
			"Katkıda bulunmak için : github.com/msrexe/yt2mp3-telegram-bot\n" +
			"@msrexe")
	})
	bot.Handle(tb.OnText, func(m *tb.Message) {
		if urlCheck(m.Text) {
			bot.Send(m.Sender,"İsteğiniz işleniyor...")
			dir,_ := ioutil.TempDir("/","prefix")
			err := downloadMP3(m.Text)
			if err == nil {
				bot.Send(m.Sender,"Dosya hazırlanıyor (Bu işlem biraz zaman alabilir)...")
				mp3 := &tb.Audio{
					File: tb.FromDisk("file.mp3"),
				}
				bot.Send(m.Sender,mp3)
				os.RemoveAll(dir)
			}else{
				bot.Send(m.Sender,"Video bulunamadı. Kaldırılmış veya hatalı URL olabilir!!")
			}
		}else{
			bot.Send(m.Sender,"Geçersiz URL!! \n Örnek istek => 'https://www.youtube.com/watch?v=jHjFxJVeCQs'")
		}
	})
	bot.Start()
}

func downloadMP3(url string) error{
	command1 := exec.Command("cd","/utils")
	command2 := exec.Command("youtube-dl","-x","--audio-format","mp3","-o","/file.mp4",url)
	command1.Run()
	_,err := command2.Output()
	return err
}

func urlCheck(urlString string) bool {
	_ ,err := http.Get(urlString)
	if err != nil {
		return false
	}
	return true
}

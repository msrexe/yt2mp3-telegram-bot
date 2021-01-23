package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
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
		bot.Send(m.Sender,"Youtube2MP3 botuna hoşgeldiniz :)\nİndirmek istediğiniz YouTube videosunun linkini mesaj olarak göndermeniz yeterlidir.\n ")
	})
	bot.Handle(tb.OnText, func(m *tb.Message) {
		if urlCheck(m.Text) {
			bot.Send(m.Sender,"İsteğiniz işleniyor...")
			bot.Send(m.Sender,"Dosya hazırlanıyor (Bu işlem biraz zaman alabilir)...")
			dir,_ := ioutil.TempDir("/","prefix")
			downloadMP3(m.Text)
			bot.Send(m.Sender,"Az kaldı ...")
			mp3 := &tb.Audio{
				File: tb.FromDisk("file.mp3"),
			}
			bot.Send(m.Sender,mp3)
			os.RemoveAll(dir)
		}else{
			bot.Send(m.Sender,"Geçersiz URL!! \n Örnek istek => 'https://www.youtube.com/watch?v=jHjFxJVeCQs'")
		}
	})

	bot.Start()
}

func downloadMP3(url string){
	command1 := exec.Command("cd","/utils")
	command2 := exec.Command("youtube-dl","-x","--audio-format","mp3","-o","/file.mp4",url)
	log.Println(command2)
	command1.Run()
	command2.Run()
}

func urlCheck(url string) bool {
	return true
}

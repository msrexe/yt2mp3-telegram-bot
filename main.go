package main

import (
	"errors"
	"github.com/kkdai/youtube/v2"
	tb "gopkg.in/tucnak/telebot.v2"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	handlers()
}

func handlers(){
	var videoID string
	bot, err := tb.NewBot(tb.Settings{
		//Token: os.Getenv("BOT_TOKEN"),
		Token: "1514859884:AAEQpQkXr1iMAcH6BvQ70Q_-RMHREP9y79A",
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
		selector.Row(btnmp3,btnmp4),
	)
	//**********************************************

	bot.Handle("/start", func(m *tb.Message) {
		bot.Send(m.Sender,"İndirmek istediğiniz YouTube linkini gönderin")
	})

	bot.Handle("/help", func(m *tb.Message) {
		bot.Send(m.Sender,"Youtube2MP3 botuna hoşgeldiniz :)\n" +
			"İndirmek istediğiniz YouTube videosunun linkini mesaj olarak göndermeniz yeterlidir.\n" +
			"Katkıda bulunmak için : github.com/msrexe/yt2mp3-telegram-bot\n" +
			"@msrexe")
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		bot.Send(m.Sender,"İsteğiniz işleniyor...")
		videoID,err = urlFormatter(m.Text)
		if err == nil {
			bot.Send(m.Sender,"İndirmek istediğiniz formatı seçin",selector)
		}else{
			log.Println(err)
			bot.Send(m.Sender,"Video bulunamadı. Kaldırılmış veya hatalı URL olabilir!!")
			return
		}
	})

	bot.Handle(&btnmp4, func(c *tb.Callback) {
		bot.Send(c.Sender,"Medya hazırlanıyor (Bu işlem biraz zaman alabilir)...")
		videoTitle,err := downloadMP4(videoID)
		if err != nil{
			bot.Send(c.Sender,"Bir hata oluştu")
			log.Println(err)
			return
		}
		mp4 := &tb.Video{
			File: tb.FromDisk(videoTitle + ".mp4"),
		}
		_,err = bot.Send(c.Sender,mp4)
		bot.Send(c.Sender,videoTitle)
		if err != nil {
			bot.Send(c.Sender,"Bir hata oluştu")
			log.Println(err)
			return
		}
	})

	bot.Handle(&btnmp3, func(c *tb.Callback) {
		bot.Send(c.Sender,"Medya hazırlanıyor (Bu işlem biraz zaman alabilir)...")

		audioTitle,err := downloadMP3(videoID)
		if err != nil{
			bot.Send(c.Sender,"Bir hata oluştu")
			log.Println(err)
			return
		}

		mp3 := &tb.Audio{
			File: tb.FromDisk(audioTitle + ".mp3"),
		}

		_,err = bot.Send(c.Sender,mp3)
		bot.Send(c.Sender,audioTitle)
		if err != nil {
			bot.Send(c.Sender,"Bir hata oluştu")
			log.Println(err)
			return
		}

	})


	bot.Start()
}


func downloadMP4(videoID string) (string,error){
	client := youtube.Client{}

	video,err := client.GetVideo(videoID)
	if err!= nil {
		log.Println(err)
	}
	resp,err := client.GetStream(video,&video.Formats[0])
	if err!= nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	file, err := os.Create(video.Title+".mp4")
	if err!= nil {
		log.Println(err)
	}
	defer file.Close()

	_, err = io.Copy(file,resp.Body)
	if err!= nil {
		log.Println(err)
	}

	return video.Title,err
}

func downloadMP3(videoID string) (string,error){
	client := youtube.Client{}

	video,err := client.GetVideo(videoID)
	if err!= nil {
		log.Println(err)
	}
	resp,err := client.GetStream(video,&video.Formats[0])
	if err!= nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	/*
	file, err := os.Create(video.Title+".mp4")
		if err!= nil {
			log.Println(err)
		}
	defer file.Close()
	 */
	/*
	pr,pw := io.Pipe()

	go func() {
		_, err = io.Copy(pw,resp.Body)
		if err!= nil {
			log.Println(err)
		}
	}()
	*/


	return video.Title,err
}


func urlFormatter(urlString string) (string,error) {
	if strings.Contains(urlString,"youtu.be/") {
		output, _:= url.Parse(urlString)
		videoID := output.Path
		if videoID == ""{
			return "",errors.New("hatalı url")
		}
		return videoID,nil
	}else{
		output, _:= url.Parse(urlString)
		videoID := output.Query().Get("v")
		if videoID == ""{
			return "",errors.New("hatalı url")
		}
		return videoID,nil
	}
}


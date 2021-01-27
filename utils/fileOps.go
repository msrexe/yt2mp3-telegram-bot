package utils

import (
	"encoding/json"
	"errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func DownloadMP4fromURL(title string ,url string) (*tb.Video,error)  {
	pr1, pw1 := io.Pipe()
	file,err := os.Create(title+".mp4")
	if err != nil {
		log.Println(err)
	}

	ytdl := exec.Command("youtube-dl", url, "-o-")
	ytdl.Stdout = pw1
	ytdl.Stderr = os.Stderr

	log.Println("indirme basladi")

	go func() {
		if err = ytdl.Run(); err != nil {
			log.Println(err)
		}
		pw1.Close()
	}()



	_,err = io.Copy(io.Writer(file),io.Reader(pr1))
	pr1.Close()
	log.Println("Pipedan kopyalama tamamlandı")
	mp4 := &tb.Video{
		File: tb.FromDisk(title+".mp4"),
	}
	return mp4,err
}

func DownloadMP3fromURL(title string,url string)  (*tb.Audio,error) {
	r1, w1 := io.Pipe()

	ytdl := exec.Command("youtube-dl","--audio-format","mp3", url, "-o-")
	ytdl.Stdout = w1
	ytdl.Stderr = os.Stderr

	ffmpeg := exec.Command("ffmpeg", "-i", "/dev/stdin", "-f", "mp3", "-ab","96000", "-vn", title+".mp3")
	ffmpeg.Stdin = r1  // PIPE OUTPUT
	ffmpeg.Stderr = os.Stderr

	go func() {
		if err := ytdl.Run(); err != nil {
			log.Println(err)
		}
		w1.Close()
	}()
	err := ffmpeg.Run()
	mp3 := &tb.Audio{
		Title: title,
		File: tb.FromDisk(title + ".mp3"),
	}

	return mp3, err
}

func UrlFormatter(urlString string) (string,string, error) {
	if strings.Contains(urlString, "youtu.be/") {
		output, _ := url.Parse(urlString)
		videoID := strings.Replace(output.Path,"/","",1)
		if videoID == "" {
			return "","", errors.New("hatalı url")
		}
		url := "https://youtube.com/watch?v=" + videoID
		title,_:= getVideoTitle(url)
		return url,title, nil
	} else {
		output, _ := url.Parse(urlString)
		videoID := output.Query().Get("v")
		if videoID == "" {
			return "","", errors.New("hatalı url")
		}
		url := "https://youtube.com/watch?v=" + videoID
		title,_:= getVideoTitle(url)
		return url,title, nil
	}
}

func getVideoTitle(url string) (string,error){
	type VideoInfo struct {
		Title   string `json:"fulltitle"`
	}

	ytdlTitle := exec.Command("youtube-dl",url,"--print-json","--no-warnings","--skip-download")
	outputJSON,err := ytdlTitle.Output()
	if err != nil {
		log.Println("getVideoTitle "+err.Error())
		return "", err
	}
	var video VideoInfo
	json.Unmarshal(outputJSON,&video)
	return video.Title,nil
}

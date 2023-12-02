package main

import (
	"fmt"
	"github.com/XiaoMiku01/douyin-live-go/player"
	dyproto "github.com/XiaoMiku01/douyin-live-go/protobuf"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/wav"
	"log"
	"os"
	"time"
)

type SingleData struct {
	path       string
	createTime time.Time
}

var singleCh = make(chan SingleData, 20)

var UserEnterRoomChan = make(chan dyproto.MemberMessage, 100)

func init() {
	go doUserEnterRoom()
	go doUserChatRoom()
	go doUserGiftRoom()
	go doUserLikeRoom()
	go single()
}

func single() {
	var streamer beep.StreamSeekCloser
	var format beep.Format
	var err error
	sr := beep.SampleRate(44100)
	//speaker.Init(sr, sr.N(time.Second/10))
	//defer streamer.Close()
	for {
		select {
		case s := <-singleCh:
			if time.Now().Add(-time.Second * 10).Before(s.createTime) {
				//if streamer != nil {
				//	streamer.Seek(0)
				//	//time.Sleep(time.Second * 5)
				//	continue
				//}
				//beep.Dup()
				f, _ := os.Open(s.path)
				streamer, format, err = wav.Decode(f)
				if err != nil {
					log.Println(err)
					continue
				}
				//f.Close()
				//streamer, format, err = mp3.Decode(f)
				resampled := beep.Resample(4, format.SampleRate, sr, streamer)
				//if err != nil {
				//	log.Println(err)
				//	continue
				//}
				player.Play(resampled, func() {
					if err = f.Close(); err != nil {
						log.Println(err)
					} else {
						os.Remove(s.path)
					}
				})
				//err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
				//if err != nil {
				//	fmt.Println(err.Error())
				//	continue
				//}
				//speaker.Play(streamer)
				//volume.Silent = true
				//time.Sleep(5 * time.Millisecond)
				//speaker.Clear()
				//volume.Silent = false
				//volume.Volume = 0
				//
				//// skipping some resampling but I update the Volume streamer here.
				//
				//speaker.Play(beep.Seq(&volume, beep.Callback(func() {
				//	done <- true
				//})))
				//time.Sleep(time.Second * 1000)
			} else {
				os.Remove(s.path)
			}
			//default:

		}
	}
}

func doUserEnterRoom() {
	var list []dyproto.MemberMessage
	skip := time.NewTicker(time.Second * 3)
	for {
		select {
		case m := <-UserEnterRoomChan:
			list = append(list, m)
		case <-skip.C:
			if len(list) != 0 {
				if len(list) == 1 {
					gender := ""
					if list[0].User.Gender == 1 {
						gender = "小哥哥"
					} else if list[0].User.Gender == 2 {
						gender = "小姐姐"
					}
					path, err := tts("demo", 1, TtsData{
						Text:       fmt.Sprintf("欢迎%s%s进入直播间", gender, list[0].User.NickName),
						NoiseScale: 1,
					})
					if err != nil {
						fmt.Println(err.Error())
					} else {
						singleCh <- SingleData{
							path:       path,
							createTime: time.Now(),
						}
					}
				} else {
					path, err := tts("demo", 1, TtsData{
						Text: fmt.Sprintf("欢迎新进来的%d位朋友", len(list)),
					})
					if err != nil {
						fmt.Println(err.Error())
					} else {
						singleCh <- SingleData{
							path:       path,
							createTime: time.Now(),
						}
					}
				}
				list = nil
			}
		}
	}
}

func UserEnterRoom(enterMsg dyproto.MemberMessage) {
	select {
	case UserEnterRoomChan <- enterMsg:
	default:
	}
}

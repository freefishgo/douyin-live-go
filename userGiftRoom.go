package main

import (
	"fmt"
	dyproto "github.com/XiaoMiku01/douyin-live-go/protobuf"
	"time"
)

var UserGiftRoomChan = make(chan dyproto.GiftMessage, 100)

func doUserGiftRoom() {
	for {
		select {
		case m := <-UserGiftRoomChan:
			if m.SendTime < uint64(time.Now().Add(-time.Second*20).Unix()) {
				break
			}
			gender := ""
			if m.User.Gender == 1 {
				gender = "小哥哥"
			} else if m.User.Gender == 2 {
				gender = "小姐姐"
			}
			path, err := tts("demo", 1, TtsData{
				Text:       fmt.Sprintf("谢谢%s%s的%d个%s", gender, m.User.NickName, m.ComboCount, m.Gift.Name),
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
		}
	}
}

func UserGiftRoom(msg dyproto.GiftMessage) {
	select {
	case UserGiftRoomChan <- msg:
	default:
	}
}

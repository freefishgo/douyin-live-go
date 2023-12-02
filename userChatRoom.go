package main

import (
	"fmt"
	dyproto "github.com/XiaoMiku01/douyin-live-go/protobuf"
	"time"
)

var UserChatRoomChan = make(chan dyproto.ChatMessage, 100)

func doUserChatRoom() {
	for {
		select {
		case m := <-UserChatRoomChan:
			if m.EventTime < uint64(time.Now().Add(-time.Second*20).Unix()) {
				break
			}
			gender := ""
			if m.User.Gender == 1 {
				gender = "小哥哥"
			} else if m.User.Gender == 2 {
				gender = "小姐姐"
			}
			path, err := tts("demo", 1, TtsData{
				Text:       fmt.Sprintf("%s%s %s", gender, m.User.NickName, m.Content),
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

func UserChatRoom(msg dyproto.ChatMessage) {
	select {
	case UserChatRoomChan <- msg:
	default:
	}
}

package main

import (
	"fmt"
	dyproto "github.com/XiaoMiku01/douyin-live-go/protobuf"
	"time"
)

var UserLikeRoomChan = make(chan dyproto.LikeMessage, 100)

func doUserLikeRoom() {
	var list []dyproto.LikeMessage
	skip := time.NewTicker(time.Second * 3)
	for {
		select {
		case m := <-UserLikeRoomChan:
			list = append(list, m)
		case <-skip.C:
			if len(list) != 0 {
				count := uint64(0)
				for _, v := range list {
					count += v.Count
				}
				gender := ""
				if len(list) == 1 {
					if list[0].User.Gender == 1 {
						gender = "小哥哥"
					} else if list[0].User.Gender == 2 {
						gender = "小姐姐"
					}
					gender += list[0].User.NickName
				} else {
					gender = "你们"
				}
				path, err := tts("demo", 1, TtsData{
					Text:       fmt.Sprintf("谢谢%s的%d个赞", gender, count),
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
				list = nil
			}
		}
	}
}

func UserLikeRoom(msg dyproto.LikeMessage) {
	select {
	case UserLikeRoomChan <- msg:
	default:
	}
}

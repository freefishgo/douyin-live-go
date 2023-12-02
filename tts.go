package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// TtsData
// text 需要转换的文字(必须), 如果是 language=Mix 的情况下就需要带上语言标签
// speed 语速(可选), 越小越慢, 默认为0.95
// noise_scale 感情变化程度(可选) 默认为 0.667
// noise_scale_w 音素发音长度(可选) 默认为0.8
// language 语言(可选) 默认为:简体中文, 选项有: 日本語, 简体中文,English
type TtsData struct {
	Text        string  `json:"text"`
	Language    string  `json:"language,omitempty"`
	Speed       float32 `json:"speed,omitempty"`
	NoiseScale  float32 `json:"noise_scale,omitempty"`
	NoiseScaleW float32 `json:"noise_scale_w,omitempty"`
}

func tts(modelName string, speakerId int, data TtsData) (path string, err error) {
	if data.Language == "" {
		data.Language = "简体中文"
	}
	url := fmt.Sprintf("http://127.0.0.1:3232/models/%s/speakers/%d", modelName, speakerId)
	b, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", io.NopCloser(bytes.NewBuffer(b)))
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	path = filepath.Join("audio", strconv.Itoa(int(time.Now().UnixMilli()))+".wav")
	f, _ := os.Create(path)
	defer f.Close()
	defer resp.Body.Close()
	//beep.NewBuffer()
	b, _ = io.ReadAll(resp.Body)
	//wav.Encode(resp.Body)
	f.Write(b)
	return
}

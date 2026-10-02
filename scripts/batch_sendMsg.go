package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Message struct {
	TargetId uint   `json:"targetId"`
	Type     string `json:"type"` //单聊，群聊 ,好友请求 chat ,groupMessage,friendRequest
	Content  string `json:"content"`
	MsgType  int    `json:"msgType"` // 0 文本 1 图片  (后续可以扩展为2 音频)
}

func main() {
	client := &http.Client{}
	url := "http://127.0.0.1:8080/chat/send"
	f, _ := os.Open("GroupTokens.txt")
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		msgBody, _ := json.Marshal(Message{
			TargetId: 5,
			Type:     "groupMessage",
			Content:  "hello world",
			MsgType:  0,
		})
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(msgBody))
		if err != nil {
			fmt.Println("创建请求失败:", err)
			continue
		}
		req.Header.Set("Authorization", sc.Text())
		_, err = client.Do(req)
		if err != nil {
			fmt.Println("发送请求失败:", err)
			continue
		}

	}
	fmt.Println("发送完成")

}

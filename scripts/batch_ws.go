package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	f, _ := os.Open("tokens.txt")
	defer f.Close()
	sc := bufio.NewScanner(f)
	var tokens []string
	for sc.Scan() {
		if t := sc.Text(); t != "" {
			tokens = append(tokens, t)
		}
	}
	// 限制最多50个同时连，解决端口爆满
	ch := make(chan struct{}, 50)
	for _, t := range tokens {
		ch <- struct{}{}
		go func(token string) {
			defer func() { <-ch }()
			u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/ws/" + token}
			c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				fmt.Println("失败", token, err)
				return
			}
			fmt.Println("用户上线：", token)
			// 保活读消息
			go func() {
				defer c.Close()
				for {
					if _, _, e := c.ReadMessage(); e != nil {
						return
					}
				}
			}()
		}(t)
	}
	for range ch {
	}
	time.Sleep(999 * time.Second)
}

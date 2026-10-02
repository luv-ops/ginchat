package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

const workerCount = 10 // 可自行调整20~50
const baseUrl = "http://127.0.0.1:8080"

type LoginReq struct {
	Name     string
	Password string
}
type Res struct {
	Code int       `json:"code"`
	Data LoginResp `json:"data"`
	Msg  string    `json:"msg"`
}
type LoginResp struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Token  string `json:"token"`
	Avatar string `json:"avatar"`
}

func main() {
	f, err := os.OpenFile("tokens.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Sprintf("打开文件失败: %v", err))
	}
	defer func() {
		_ = f.Sync()
		_ = f.Close()
	}()

	var wg sync.WaitGroup
	var fileMu sync.Mutex // 多协程写文件加锁，防止换行错乱
	taskChan := make(chan int, 100)

	// 启动工作协程
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{}
			for uid := range taskChan {
				uname := fmt.Sprintf("stress_user_%d", uid)
				// 登录拿token
				loginBody, _ := json.Marshal(LoginReq{Name: uname, Password: "123456"})
				resp, err := client.Post(baseUrl+"/user/login", "application/json", bytes.NewBuffer(loginBody))
				if err != nil {
					fmt.Printf("用户 %d 登录失败:%v\n", i, err)
					continue
				}

				var res Res
				err = json.NewDecoder(resp.Body).Decode(&res)
				_ = resp.Body.Close()

				// 写文件加互斥锁
				line := res.Data.Token + "\n"
				fileMu.Lock()
				_, err = f.WriteString(line)
				fileMu.Unlock()
				if err != nil {
					fmt.Printf("用户%d 写入token失败: %v\n", uid, err)
					continue
				}
				fmt.Printf("用户 %d token生成成功\n", uid)
			}
		}()
	}

	// 投递任务
	start, end := 10000, 10001
	for uid := start; uid <= end; uid++ {
		taskChan <- uid
	}
	close(taskChan)
	wg.Wait()

	fmt.Printf("全部 %d~%d token 追加完成\n", start, end)
}

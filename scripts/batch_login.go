package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const baseUrl = "http://127.0.0.1:8080"
const testCount = 5000 // 需要多少并发就生成多少用户

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
	f, err := os.Create("tokens.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	client := &http.Client{}
	for i := 1; i <= testCount; i++ {
		uname := fmt.Sprintf("stress_user_%d", i)
		pwd := "123456"

		// 登录拿token
		loginBody, _ := json.Marshal(LoginReq{Name: uname, Password: pwd})
		resp, err := client.Post(baseUrl+"/user/login", "application/json", bytes.NewBuffer(loginBody))
		if err != nil {
			fmt.Printf("用户 %d 登录失败:%v\n", i, err)
			continue
		}

		var res Res
		_ = json.NewDecoder(resp.Body).Decode(&res)
		io.Copy(f, bytes.NewBufferString(res.Data.Token+"\n"))
		resp.Body.Close()
		fmt.Printf("生成用户 %d token成功\n", i)
	}
	fmt.Println("全部token已写入 tokens.txt")
}

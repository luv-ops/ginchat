package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RegisterReq struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

const burl = "http://127.0.0.1:8080"
const count = 5000

func main() {
	client := &http.Client{}
	for i := 1; i <= 10001; i++ {
		uname := fmt.Sprintf("stress_user_%d", i)
		pwd := "123456"

		req := RegisterReq{
			Name:            uname,
			Password:        pwd,
			ConfirmPassword: pwd,
		}

		body, _ := json.Marshal(req)
		_, err := client.Post(burl+"/user/register", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("register failed:", err)
			continue
		}
		fmt.Printf("register success %d\n", i)
	}
	fmt.Printf("register success %s\n", "用户注册完成")
}

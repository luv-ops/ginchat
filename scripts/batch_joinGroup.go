package main

import (
	"fmt"
	"net/http"
)

func main() {
	url := "http://127.0.0.1:8080/group/test/join"

	client := &http.Client{}

	_, err := client.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println("register failed:", err)
	}

}

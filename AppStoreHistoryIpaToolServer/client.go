package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

// 1. 声明保存命令行参数的变量
var (
	op      string
	appleId string
	taskId  string
)

func test() error {
	//go run main.go --name ls --value /etc -r
	flag.StringVar(&op, "Op", "", "Op")
	flag.StringVar(&appleId, "AppleId", "", "AppleId")
	flag.StringVar(&taskId, "TaskId", "", "TaskId")

	if op == "requestTask" {
		requestTask()
	} else if op == "requestTask" {

	}
	return nil
}

func requestTask() {
	/*	urlValues := url.Values{}
		urlValues.Add("name", "zhaofan")
		urlValues.Add("age", "22")

		reqBody := urlValues.Encode()*/

	response, err := http.Post("http://127.0.0.1/requestTask", "text/html", nil)
	if err != nil {
		log.Printf("requestTask failed, err:%v", err)
		return
	}

	body, _ := io.ReadAll(response.Body)
	fmt.Println(string(body))
}

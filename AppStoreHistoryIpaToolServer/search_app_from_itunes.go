package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func searchApp(appName string, country string) (*WebSearchAppResult, error) {
	log.Printf("searchApp in")
	defer log.Printf("searchApp out")

	if 0 == len(appName) || 0 == len(country) {
		return nil, fmt.Errorf("param invalid")
	}

	// 进行 URL 编码
	encodedAppName := url.QueryEscape(appName)
	// 发送 GET 请求
	urlString := fmt.Sprintf("https://itunes.apple.com/search?term=%s&entity=software&country=%s&limit=10", encodedAppName, country)

	resp, err := http.Get(urlString)
	if err != nil {
		log.Printf("searchApp get request err:%v", err)
		return nil, fmt.Errorf("searchApp get request err:%v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("searchApp read body err:%v", err)
		return nil, fmt.Errorf("searchApp read body err:%v", err)
	}

	var searchAppResult WebSearchAppResult
	err = json.Unmarshal(body, &searchAppResult)
	if err != nil {
		log.Printf("searchApp Unmarshal err:%v", err)
		return nil, fmt.Errorf("searchApp Unmarshal err:%v", err)
	}

	//fmt.Println("Response:", string(body))
	return &searchAppResult, nil
}

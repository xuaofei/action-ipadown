package main

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
)

func runServer(c *cli.Context) error {
	go startWebServer()
	startScriptServer()

	return nil
}

func startWebServer() {
	log.Printf("startWebServer in")
	defer log.Printf("startWebServer out")

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.templates", "templates/template2.templates")

	router.GET("/", webHomeHandler)                                          // 网站主页
	router.POST("/login", webLoginHandler)                                   // appleid applepwd
	router.GET("/loginResult", webLoginResultHandler)                        // 查询登录结果
	router.POST("/verifyCode", webVerifyCodeHandler)                         // 输入二次验证码
	router.POST("/verifyCodeResult", webVerifyCodeResultHandler)             // 查询二次验证码登录结果
	router.POST("/searchApp", webSearchAppHandler)                           // 搜索应用
	router.POST("/searchAppVersion", webSearchAppVersionHandler)             // 搜索应用版本号
	router.POST("/searchAppVersionResult", webSearchAppVersionResultHandler) // 搜索应用版本号结果

	_ = router.Run(":8080")
}

func startScriptServer() {
	log.Printf("startScriptServer in")
	defer log.Printf("startScriptServer out")

	http.HandleFunc("/scriptTaskIdRequest", scriptTaskIdHandler)
	http.HandleFunc("/scriptLoginInfoRequest", scriptLoginInfoHandler)
	http.HandleFunc("/script2FARequest", script2FAHandler)
	http.HandleFunc("/scriptUploadVersionsInfoRequest", scriptUploadVersionsInfoHandler)
	http.HandleFunc("/scriptDownloadListRequest", scriptDownloadListHandler)
	http.HandleFunc("/scriptReportResult", scriptReportResultHandler)

	_ = http.ListenAndServe(":80", nil)
}

func main() {
	log.Printf("args:%v", os.Args)

	app := &cli.App{
		Name:   "AppStoreHistoryIpaTool",
		Usage:  "run as client",
		Action: runServer,
	}

	_ = app.Run(os.Args)
}

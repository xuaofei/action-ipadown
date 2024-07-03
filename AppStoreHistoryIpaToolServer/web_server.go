package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

//func makeResponse(result int, msg string) map[string]any {
//	return gin.H{
//		"code":    result,
//		"message": msg,
//	}
//}

func webHomeHandler(c *gin.Context) {
	log.Printf("webHomeHandler in")
	defer log.Printf("webHomeHandler out")

	c.HTML(http.StatusOK, "home.html", nil)
}

func webLoginHandler(c *gin.Context) {
	log.Printf("webLoginHandler in")
	defer log.Printf("webLoginHandler out")

	// 获取表单数据
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 这里可以添加验证逻辑
	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("用户名和密码不能为空")})
		return
	}

	taskId := makeTaskId()
	err := GetDBInstance().InsertTask(taskId)
	if err != nil {
		// 返回JSON响应
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("taksid make failed:%v", err),
		})
		return
	}

	err = GetDBInstance().UpdateTaskLoginInfo(taskId, username, password)
	if err != nil {
		// 返回JSON响应
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("Update Login failed:%v", err),
		})
		return
	}
	// 开启登录
	GetSMInstance().StartTask()

	// 设置 cookie
	c.SetCookie("task_id", taskId, 3600, "/", "", false, false)
	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"code":    SUCCESS,
		"message": "",
	})
}

func webLoginResultHandler(c *gin.Context) {
	log.Printf("webLoginResultHandler in")
	defer log.Printf("webLoginResultHandler out")

	task_id, err := c.Cookie("task_id")
	if err != nil || len(task_id) == 0 {
		log.Printf("webLoginResultHandler not find cookie task_id, err:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("query login status failed:%v", err),
		})
		return
	}
	log.Printf("webLoginResultHandler task_id:%v", task_id)

	loginStatus, err := GetDBInstance().QueryTaskLoginStatus(task_id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("query login status failed:%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        SUCCESS,
		"message":     "",
		"loginStatus": loginStatus,
	})
}

func webVerifyCodeHandler(c *gin.Context) {
	log.Printf("webVerifyCodeHandler in")
	defer log.Printf("webVerifyCodeHandler out")

	task_id, err := c.Cookie("task_id")
	if err != nil || len(task_id) == 0 {
		log.Printf("webVerifyCodeHandler not find cookie task_id, err:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("webVerifyCodeHandler failed:%v", err),
		})
		return
	}
	log.Printf("webVerifyCodeHandler task_id:%v", task_id)

	// 获取表单数据
	verifyCode := c.PostForm("verifyCode")
	// 这里可以添加验证逻辑
	if verifyCode == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("verifyCode is empty"),
		})
		return
	}

	err = GetDBInstance().UpdateTask2FA(task_id, verifyCode)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("webVerifyCodeHandler failed:%v", err),
		})
		return
	}

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"code":    SUCCESS,
		"message": "",
	})
}

func webVerifyCodeResultHandler(c *gin.Context) {
	log.Printf("webVerifyCodeResultHandler in")
	defer log.Printf("webVerifyCodeResultHandler out")

	task_id, err := c.Cookie("task_id")
	if err != nil || len(task_id) == 0 {
		log.Printf("webVerifyCodeResultHandler not find cookie task_id, err:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("query VerifyCode status failed:%v", err),
		})
		return
	}
	log.Printf("webVerifyCodeResultHandler task_id:%v", task_id)

	verifyCodeStatus, err := GetDBInstance().QueryTaskVerifyCodeStatus(task_id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("query verifyCodeStatus failed:%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":             SUCCESS,
		"message":          "",
		"verifyCodeStatus": verifyCodeStatus,
	})
}

func webSearchAppHandler(c *gin.Context) {
	log.Printf("webSearchAppHandler in")
	defer log.Printf("webSearchAppHandler out")

	// 获取表单数据
	appName := c.PostForm("appName")
	country := c.PostForm("country")

	// 这里可以添加验证逻辑
	if appName == "" || country == "" {
		c.JSON(http.StatusOK, gin.H{"error": "应用名称或国家不能为空"})
		return
	}

	log.Printf("webSearchAppHandler appName:%v country:%v", appName, country)

	webSearchAppResult, err := searchApp(appName, country)
	if err != nil {
		log.Printf("webSearchAppHandler in err:%v", err)
		return
	}

	// 返回JSON响应
	c.JSON(http.StatusOK, webSearchAppResult)
}

func webSearchAppVersionHandler(c *gin.Context) {
	log.Printf("webSearchAppVersionHandler in")
	defer log.Printf("webSearchAppVersionHandler out")

	var appData AppData

	// 解析传入的 JSON 数据
	if err := c.ShouldBindJSON(&appData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task_id, err := c.Cookie("task_id")
	if err != nil || len(task_id) == 0 {
		log.Printf("webSearchAppVersionHandler not find cookie task_id, err:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("query SearchAppVersion status failed:%v", err),
		})
		return
	}
	log.Printf("webVerifyCodeResultHandler task_id:%v", task_id)

	err = GetDBInstance().UpdateTaskDownloadIpaInfo(task_id, strconv.Itoa(appData.TrackId), appData.Price)
	if err != nil {
		log.Printf("webSearchAppVersionHandler UpdateTaskDownloadIpaInfo, err:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("UpdateTaskDownloadIpaInfo failed:%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    SUCCESS,
		"message": "",
	})
}

func webSearchAppVersionResultHandler(c *gin.Context) {
	log.Printf("webSearchAppVersionResultHandler in")
	defer log.Printf("webSearchAppVersionResultHandler out")

	task_id, err := c.Cookie("task_id")
	if err != nil || len(task_id) == 0 {
		log.Printf("webSearchAppVersionResultHandler not find cookie task_id, err:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("query AppVersion failed:%v", err),
		})
		return
	}
	log.Printf("webSearchAppVersionResultHandler task_id:%v", task_id)

	verifyCodeStatus, err := GetDBInstance().QueryTaskVerifyCodeStatus(task_id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    FAILED,
			"message": fmt.Sprintf("query verifyCodeStatus failed:%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":             SUCCESS,
		"message":          "",
		"verifyCodeStatus": verifyCodeStatus,
	})

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "App data received successfully",
	})
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
)

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
		webResponse(c, FAILED, fmt.Sprintf("username or password is empty"), "")
		return
	}

	taskId := makeTaskId()
	err := GetDBInstance().InsertTask(taskId)
	if err != nil {
		webResponse(c, FAILED, fmt.Sprintf("taksid make failed:%v", err), "")
		return
	}

	err = GetDBInstance().UpdateAppleIDAndPassword(taskId, username, password)
	if err != nil {
		webResponse(c, FAILED, fmt.Sprintf("Update Login failed:%v", err), taskId)
		return
	}

	err = GetDBInstance().UpdateTaskStatus(taskId, TASK_INPUT_LOGIN_INFO)
	if err != nil {
		webResponse(c, FAILED, fmt.Sprintf("update task status failed:%v", err), taskId)
		return
	}

	// 开启登录
	GetSMInstance().StartTask()
	// 设置 cookie
	c.SetCookie("task_id", taskId, 3600, "/", "", false, false)
	// 返回JSON响应
	webResponse(c, SUCCESS, "", taskId)
}

func webLoginResultHandler(c *gin.Context) {
	log.Printf("webLoginResultHandler in")
	defer log.Printf("webLoginResultHandler out")

	taskId, err := c.Cookie("task_id")
	if err != nil || len(taskId) == 0 {
		webResponse(c, FAILED, fmt.Sprintf("query login status failed:%v", err), taskId)
		return
	}
	log.Printf("webLoginResultHandler task_id:%v", taskId)

	// 返回JSON响应
	webResponse(c, SUCCESS, "", taskId)
}

func webVerifyCodeHandler(c *gin.Context) {
	log.Printf("webVerifyCodeHandler in")
	defer log.Printf("webVerifyCodeHandler out")

	taskId, err := c.Cookie("task_id")
	if err != nil || len(taskId) == 0 {
		webResponse(c, FAILED, fmt.Sprintf("webVerifyCodeHandler not find taskId failed:%v", err), taskId)
		return
	}
	log.Printf("webVerifyCodeHandler task_id:%v", taskId)

	// 获取表单数据
	verifyCode := c.PostForm("verifyCode")
	// 这里可以添加验证逻辑
	if verifyCode == "" {
		webResponse(c, FAILED, fmt.Sprintf("verifyCode is empty"), taskId)
		return
	}

	err = GetDBInstance().Update2FA(taskId, verifyCode)
	if err != nil {
		webResponse(c, FAILED, fmt.Sprintf("webVerifyCodeHandler failed:%v", err), taskId)
		return
	}

	webResponse(c, SUCCESS, "", taskId)
}

func webVerifyCodeResultHandler(c *gin.Context) {
	log.Printf("webVerifyCodeResultHandler in")
	defer log.Printf("webVerifyCodeResultHandler out")

	taskId, err := c.Cookie("task_id")
	if err != nil || len(taskId) == 0 {
		webResponse(c, FAILED, fmt.Sprintf("webVerifyCodeResultHandler not find taskId failed:%v", err), taskId)
		return
	}
	log.Printf("webVerifyCodeResultHandler task_id:%v", taskId)

	webResponse(c, SUCCESS, "", taskId)
}

func webSearchAppHandler(c *gin.Context) {
	log.Printf("webSearchAppHandler in")
	defer log.Printf("webSearchAppHandler out")

	// 获取表单数据
	appName := c.PostForm("appName")
	country := c.PostForm("country")

	// 这里可以添加验证逻辑
	if appName == "" || country == "" {
		webResponse(c, FAILED, fmt.Sprintf("appName or country is empty"), "")
		return
	}

	log.Printf("webSearchAppHandler appName:%v country:%v", appName, country)

	webSearchAppResult, err := searchApp(appName, country)
	if err != nil {
		log.Printf("webSearchAppHandler in err:%v", err)
		return
	}

	c.JSON(http.StatusOK, webSearchAppResult)
}

func webSelectAppHandler(c *gin.Context) {
	log.Printf("webSelectAppHandler in")
	defer log.Printf("webSelectAppHandler out")

	taskId, err := c.Cookie("task_id")
	if err != nil || len(taskId) == 0 {
		webResponse(c, FAILED, fmt.Sprintf("webSelectAppHandler not find taskId failed:%v", err), taskId)
		return
	}
	log.Printf("webSelectAppHandler task_id:%v", taskId)

	var appData AppData

	// 解析传入的 JSON 数据
	if err := c.ShouldBindJSON(&appData); err != nil {
		webResponse(c, FAILED, fmt.Sprintf("request data invalid:%v", err), taskId)
		return
	}

	log.Printf("webSelectAppHandler task_id:%v", taskId)

	err = GetDBInstance().UpdateDownloadIpaInfo(taskId, strconv.Itoa(appData.TrackId), appData.Price)
	if err != nil {
		webResponse(c, FAILED, fmt.Sprintf("UpdateDownloadIpaInfo failed:%v", err), taskId)
		return
	}

	err = GetDBInstance().UpdateTaskStatus(taskId, TASK_SELECTED_APP)
	if err != nil {
		webResponse(c, FAILED, fmt.Sprintf("UpdateTaskStatus failed:%v", err), taskId)
		return
	}

	webResponse(c, SUCCESS, "", taskId)
}

func webSearchAppVersionResultHandler(c *gin.Context) {
	log.Printf("webSearchAppVersionResultHandler in")
	defer log.Printf("webSearchAppVersionResultHandler out")

	taskId, err := c.Cookie("task_id")
	if err != nil || len(taskId) == 0 {
		webResponse(c, FAILED, fmt.Sprintf("webSearchAppVersionResultHandler not find taskId failed:%v", err), taskId)
		return
	}
	log.Printf("webSearchAppVersionResultHandler task_id:%v", taskId)

	webResponse(c, SUCCESS, "", taskId)
}

func webSelectAppAppVersionHandler(c *gin.Context) {
	log.Printf("webSelectAppAppVersionHandler in")
	defer log.Printf("webSelectAppAppVersionHandler out")

	taskId, err := c.Cookie("task_id")
	if err != nil || len(taskId) == 0 {
		webResponse(c, FAILED, fmt.Sprintf("webSelectAppAppVersionHandler not find taskId failed:%v", err), taskId)
		return
	}
	log.Printf("webSelectAppAppVersionHandler task_id:%v", taskId)

	selectAppVersion, err := io.ReadAll(c.Request.Body)
	if err != nil {
		webResponse(c, FAILED, fmt.Sprintf("webSelectAppAppVersionHandler Read body failed:%v", err), taskId)
		return
	}

	log.Printf("webSelectAppAppVersionHandler task_id:%v", taskId)

	err = GetDBInstance().UpdateDownloadVersion(taskId, string(selectAppVersion))
	if err != nil {
		webResponse(c, FAILED, fmt.Sprintf("webSelectAppAppVersionHandler failed:%v", err), taskId)
		return
	}

	err = GetDBInstance().UpdateTaskStatus(taskId, TASK_SELECT_VERSION_LIST)
	if err != nil {
		webResponse(c, FAILED, fmt.Sprintf("UpdateTaskStatus failed:%v", err), taskId)
		return
	}

	webResponse(c, SUCCESS, "", taskId)
}

func webTaskInfoHandler(c *gin.Context) {
	log.Printf("webTaskInfoHandler in")
	defer log.Printf("webTaskInfoHandler out")

	taskId, err := c.Cookie("task_id")
	if err != nil || len(taskId) == 0 {
		webResponse(c, FAILED, fmt.Sprintf("webTaskInfoHandler not find taskId failed:%v", err), taskId)
		return
	}
	log.Printf("webTaskInfoHandler task_id:%v", taskId)

	webResponse(c, SUCCESS, "", taskId)
}

func webResponse(c *gin.Context, code int, msg string, taskId string) {
	var webResponse WebServerCommonResponse
	webResponse.Code = code
	webResponse.Message = msg

	if webResponse.Code == FAILED {
		log.Printf("web request failed, url:%v code:%v msg:%v, taskid:%v", c.Request.URL.Path, code, msg, taskId)
	} else {
		if len(taskId) == 0 {
			webResponse.Code = FAILED
			webResponse.Message = fmt.Sprintf("taskId is empty")
			log.Printf("web request query taskid is empty, url:%v code:%v msg:%v", c.Request.URL.Path, code, msg)
		} else {
			taskInfo, err := GetDBInstance().QueryTaskData(taskId)
			if err != nil {
				webResponse.Code = FAILED
				webResponse.Message = fmt.Sprintf("query task data failed:%v", err)
				log.Printf("web request query task data failed:%v, url:%v code:%v msg:%v, taskid:%v", err, c.Request.URL.Path, code, msg, taskId)
			} else {
				webResponse.TaskInfo = *taskInfo
				log.Printf("web request success, url:%v taskid:%v", c.Request.URL.Path, taskId)
			}
		}
	}

	c.JSON(http.StatusOK, webResponse)
}

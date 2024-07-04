package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

var task IPATaskDetail = IPATaskDetail{TaskId: "20240620_123456789_10003", AppleID: "fsi2clsfiz8@163.com", ApplePwd: "Ls112211", AppID: "583376064", Country: "JP", DownloadVersionList: nil}
var allVersion ipaVersionsRequest

var lockDB sync.Mutex

func scriptTaskIdHandler(c *gin.Context) {
	// 脚本获取taskid,需要加锁，避免多个脚本获取到同一个taskid.
	lockDB.Lock()
	defer lockDB.Unlock()

	log.Printf("scriptTaskIdHandler in")
	defer log.Printf("scriptTaskIdHandler out")

	taskId, err := GetDBInstance().QueryNotStartedTask()
	if err != nil {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(""))
		return
	}

	// lock task
	err = GetDBInstance().UpdateLoginStatus(taskId, EMAIL_LOGINING)
	if err != nil {
		c.HTML(http.StatusOK, "", nil)
		return
	}

	// 直接返回taskId
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(taskId))
	//c.HTML(http.StatusOK, taskId, nil)

	log.Printf("scriptTaskIdHandler success taskid:%v", taskId)
}

func scriptLoginInfoHandler(c *gin.Context) {
	log.Printf("scriptLoginInfoHandler in")
	defer log.Printf("scriptLoginInfoHandler out")

	// TODO:loginInfoRequestScript和其他的 合成一个就可以了，只需要taskid
	var request loginInfoRequestScript

	// 解析传入的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	appleId, pwd, err := GetDBInstance().QueryAppleIDAndPassword(request.TaskId)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	var task IPATaskDetail
	task.TaskId = request.TaskId
	task.AppleID = appleId
	task.ApplePwd = pwd

	c.JSON(http.StatusOK, task)
	log.Printf("scriptLoginInfoHandler:%v", task.TaskId)
}

func script2FAHandler(c *gin.Context) {
	log.Printf("script2FAHandler in")
	defer log.Printf("script2FAHandler out")

	var request scriptCommomRequest

	// 解析传入的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

		log.Printf("script2FAHandler failed,TaskID:%v err:%v", request.TaskId, err)
		return
	}

	tfa, err := GetDBInstance().Query2FA(request.TaskId)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

		log.Printf("script2FAHandler failed,TaskID:%v err:%v", request.TaskId, err)
		return
	}

	err = GetDBInstance().UpdateLoginStatus(request.TaskId, EMAIL_LOGINED)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

		log.Printf("script2FAHandler failed,TaskID:%v err:%v", request.TaskId, err)
		return
	}

	err = GetDBInstance().Update2FAStatus(request.TaskId, TFA_LOGINING)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

		log.Printf("script2FAHandler failed,TaskID:%v err:%v", request.TaskId, err)
		return
	}

	log.Printf("script2FAHandler,TaskID:%v tfa:%v", request.TaskId, tfa)
	var task IPATaskDetail
	task.TaskId = request.TaskId
	task.TwoFaCode = tfa

	c.JSON(http.StatusOK, task)
	log.Printf("scriptLoginInfoHandler:%v tfa:%v", task.TaskId, tfa)
}

func scriptItunesLoginResultHandler(c *gin.Context) {
	log.Printf("scriptItunesLoginResultHandler in")
	defer log.Printf("scriptItunesLoginResultHandler out")

	var request scriptCommomRequest

	// 解析传入的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

		log.Printf("scriptItunesLoginResultHandler failed,TaskID:%v err:%v", request.TaskId, err)
		return
	}

	err := GetDBInstance().Update2FAStatus(request.TaskId, TFA_LOGINED)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

		log.Printf("scriptItunesLoginResultHandler failed,TaskID:%v err:%v", request.TaskId, err)
		return
	}

	log.Printf("scriptItunesLoginResultHandler,TaskID:%v", request.TaskId)

	c.JSON(http.StatusOK, nil)
}

func scriptCommandHandler(c *gin.Context) {
	log.Printf("scriptCommandHandler in")
	defer log.Printf("scriptCommandHandler out")

	var request scriptCommomRequest

	// 解析传入的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

		log.Printf("scriptItunesLoginResultHandler failed,TaskID:%v err:%v", request.TaskId, err)
		return
	}

	err := GetDBInstance().Update2FAStatus(request.TaskId, TFA_LOGINED)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

		log.Printf("scriptItunesLoginResultHandler failed,TaskID:%v err:%v", request.TaskId, err)
		return
	}

	log.Printf("scriptItunesLoginResultHandler,TaskID:%v", request.TaskId)

	c.JSON(http.StatusOK, nil)
}

// 上传ipa版本信息
func scriptUploadVersionsInfoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("scriptUploadVersionsInfoHandler in")
	defer log.Printf("scriptUploadVersionsInfoHandler out")

	//w.Header().Set("Content-Type", "application/json")
	request := ipaVersionsRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("scriptUploadVersionsInfoHandler failed, param decode err:%v", err)
		return
	}

	allVersion = request

	log.Printf("scriptUploadVersionsInfoHandler TaskID:%v", request.TaskID)

	//fmt.Println(string(content))

	_, err = w.Write(nil)
	if err != nil {
		log.Printf("script2FAHandler failed, Write err:%v", err)
		return
	}
}

func scriptDownloadListHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("scriptDownloadListHandler in")
	defer log.Printf("scriptDownloadListHandler out")

	request := downloadIpaVersionsRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("scriptDownloadListHandler failed, param err:%v", err)
		return
	}

	if request.TaskId != task.TaskId {
		log.Printf("scriptDownloadListHandler failed, task id not find:%v", request.TaskId)
		return
	}

	log.Printf("scriptDownloadListHandler:%v", task.TaskId)
	task.DownloadVersionList = allVersion.AllVersionList

	data, err := json.Marshal(task)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("scriptDownloadListHandler failed, Write err:%v", err)
		return
	}

	log.Printf("scriptDownloadListHandler success")
}

func scriptReportResultHandler(c *gin.Context) {
	log.Printf("scriptReportResultHandler in")
	defer log.Printf("scriptReportResultHandler out")

	var request reportResultRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	log.Printf("scriptReportResultHandler taskid:%v code:%d msg:%v", request.TaskId, request.Code, request.Msg)
	c.JSON(http.StatusOK, nil)
}

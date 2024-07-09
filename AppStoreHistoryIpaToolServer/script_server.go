package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

//var task IPATaskDetail = IPATaskDetail{TaskId: "20240620_123456789_10003", AppleID: "fsi2clsfiz8@163.com", ApplePwd: "Ls112211", AppID: "583376064", Country: "JP", DownloadVersionList: nil}

var lockDB sync.Mutex

func scriptLockTaskIdHandler(c *gin.Context) {
	// 脚本获取taskid,需要加锁，避免多个脚本获取到同一个taskid.
	lockDB.Lock()
	defer lockDB.Unlock()

	log.Printf("scriptLockTaskIdHandler in")
	defer log.Printf("scriptLockTaskIdHandler out")

	taskId, err := GetDBInstance().QueryNotStartedTask()
	if err != nil {
		scriptResponse(c, FAILED, fmt.Sprintf("not find valid task"), taskId)
		return
	}

	// lock task
	err = GetDBInstance().UpdateTaskStatus(taskId, TASK_LOGINING)
	if err != nil {
		scriptResponse(c, FAILED, fmt.Sprintf("modify task status failed"), taskId)
		return
	}

	scriptResponse(c, SUCCESS, "", taskId)
}

func scriptTaskInfoHandler(c *gin.Context) {
	// 脚本获取taskid,需要加锁，避免多个脚本获取到同一个taskid.
	lockDB.Lock()
	defer lockDB.Unlock()

	log.Printf("scriptTaskInfoHandler in")
	defer log.Printf("scriptTaskInfoHandler out")

	var loginResult scriptCommomRequest
	err := c.ShouldBindJSON(&loginResult)
	if err != nil {
		scriptResponse(c, FAILED, fmt.Sprintf("scriptLoginResultHandler failed, invalid data:%v", err), "")
		return
	}
	taskId := loginResult.TaskId

	//taskId, err := GetDBInstance().QueryNotStartedTask()
	//if err != nil {
	//	scriptResponse(c, FAILED, fmt.Sprintf("not find valid task"), taskId)
	//	return
	//}

	//// modify task status
	//err = GetDBInstance().UpdateTaskStatus(taskId, TASK_LOGINING)
	//if err != nil {
	//	scriptResponse(c, FAILED, fmt.Sprintf("modify task status failed"), taskId)
	//	return
	//}

	scriptResponse(c, SUCCESS, "", taskId)
}

func scriptLoginResultHandler(c *gin.Context) {
	log.Printf("scriptLoginResultHandler in")
	defer log.Printf("scriptLoginResultHandler out")

	var loginResult scriptLoginResultRequest
	err := c.ShouldBindJSON(&loginResult)
	if err != nil {
		scriptResponse(c, FAILED, fmt.Sprintf("scriptLoginResultHandler failed, invalid data:%v", err), "")
		return
	}
	taskId := loginResult.TaskId

	result := TASK_LOGINING
	if loginResult.Result == 0 {
		// 需要二次认证或者账号密码错误
		result = TASK_LOGINING_2FA
	} else {
		// 登录成功
		result = TASK_LOGINED
	}

	log.Printf("login status:%v", result)

	err = GetDBInstance().UpdateTaskStatus(taskId, result)
	if err != nil {
		scriptResponse(c, FAILED, fmt.Sprintf("Update login status failed:%v", err), taskId)
		return
	}

	scriptResponse(c, SUCCESS, "", taskId)
}

func scriptUploadAllVersionHandler(c *gin.Context) {
	log.Printf("scriptUploadAllVersionHandler in")
	defer log.Printf("scriptUploadAllVersionHandler out")

	var request scriptUploadAllVersionsRequest

	// 解析传入的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		scriptResponse(c, FAILED, fmt.Sprintf("invalid data format:%v", err), "")
		return
	}
	taskId := request.TaskId

	data, err := json.Marshal(request.AllVersion)
	if err != nil {
		scriptResponse(c, FAILED, fmt.Sprintf("Marshal AllVersion failed:%v", err), taskId)
		return
	}

	// 将 data 编码为 Base64 字符串
	allVersionStr := base64.StdEncoding.EncodeToString(data)
	err = GetDBInstance().UpdateAllVersion(taskId, allVersionStr)

	if err != nil {
		scriptResponse(c, FAILED, fmt.Sprintf("UpdateAllVersion err:%v", err), taskId)
		return
	}

	err = GetDBInstance().UpdateTaskStatus(taskId, TASK_GET_VERSION_LIST)
	if err != nil {
		scriptResponse(c, FAILED, fmt.Sprintf("UpdateTaskStatus err:%v", err), taskId)
		return
	}

	scriptResponse(c, SUCCESS, "", taskId)
}

func scriptDownloadIPAResultHandler(c *gin.Context) {
	log.Printf("scriptDownloadIPAResultHandler in")
	defer log.Printf("scriptDownloadIPAResultHandler out")

	var request scriptDownloadIPAResultRequest

	// 解析传入的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		scriptResponse(c, FAILED, fmt.Sprintf("invalid data format:%v", err), "")
		return
	}
	taskId := request.TaskId

	if request.Result == 0 {
		// 失败
	} else {
		// 成功
		err := GetDBInstance().UpdateTaskStatus(taskId, TASK_FINISH_DOWNLOAD_IPA)
		if err != nil {
			scriptResponse(c, FAILED, fmt.Sprintf("UpdateTaskStatus err:%v", err), taskId)
			return
		}
	}

	err := GetDBInstance().UpdateTaskCompleted(taskId, TASK_COMPLETED)
	if err != nil {
		scriptResponse(c, FAILED, fmt.Sprintf("UpdateTaskCompleted err:%v", err), taskId)
		return
	}

	scriptResponse(c, SUCCESS, "", taskId)
}

func scriptLoginInfoHandler(c *gin.Context) {
	log.Printf("scriptLoginInfoHandler in")
	defer log.Printf("scriptLoginInfoHandler out")

	var request scriptCommomRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		scriptResponse(c, FAILED, fmt.Sprintf("scriptLoginInfoHandler request data invalid:%v", err), "")
		return
	}

	scriptResponse(c, SUCCESS, "", request.TaskId)
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

	//err = GetDBInstance().UpdateLoginStatus(request.TaskId, TASK_LOGINING)
	//if err != nil {
	//	c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
	//
	//	log.Printf("script2FAHandler failed,TaskID:%v err:%v", request.TaskId, err)
	//	return
	//}

	//err = GetDBInstance().UpdateTaskStatus(request.TaskId, TFA_LOGINING)
	//if err != nil {
	//	c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
	//
	//	log.Printf("script2FAHandler failed,TaskID:%v err:%v", request.TaskId, err)
	//	return
	//}

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

	err := GetDBInstance().Update2FAStatus(request.TaskId, TASK_LOGINING)
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

		log.Printf("scriptCommandHandler failed,TaskID:%v err:%v", request.TaskId, err)
		return
	}

	allVersion, err := GetDBInstance().QueryAllVersion(request.TaskId)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

		log.Printf("scriptCommandHandler failed,TaskID:%v err:%v", request.TaskId, err)
		return
	}

	command := ""
	appId := ""
	if len(allVersion) == 0 {
		// free app
		command = "uploadAllVersionInfo"

		appId, err = GetDBInstance().QueryAppId(request.TaskId)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})

			log.Printf("scriptCommandHandler failed,TaskID:%v err:%v", request.TaskId, err)
			return
		}
	} else {
		command = "downloadIpa"
	}

	var response scriptCommonResponse
	response.Command = command
	response.AppId = appId

	log.Printf("scriptCommandHandler,TaskID:%v", request.TaskId)

	c.JSON(http.StatusOK, response)
}

// 上传ipa版本信息
func scriptUploadVersionsInfoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("scriptUploadVersionsInfoHandler in")
	defer log.Printf("scriptUploadVersionsInfoHandler out")

	//w.Header().Set("Content-Type", "application/json")
	request := scriptUploadAllVersionsRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("scriptUploadVersionsInfoHandler failed, param decode err:%v", err)
		return
	}

	log.Printf("scriptUploadVersionsInfoHandler TaskID:%v", request.TaskId)

	//fmt.Println(string(content))

	_, err = w.Write(nil)
	if err != nil {
		log.Printf("script2FAHandler failed, Write err:%v", err)
		return
	}
}

//func scriptDownloadListHandler(w http.ResponseWriter, r *http.Request) {
//	log.Printf("scriptDownloadListHandler in")
//	defer log.Printf("scriptDownloadListHandler out")
//
//	request := downloadIpaVersionsRequest{}
//	err := json.NewDecoder(r.Body).Decode(&request)
//	if err != nil {
//		log.Printf("scriptDownloadListHandler failed, param err:%v", err)
//		return
//	}
//
//	if request.TaskId != task.TaskId {
//		log.Printf("scriptDownloadListHandler failed, task id not find:%v", request.TaskId)
//		return
//	}
//
//	log.Printf("scriptDownloadListHandler:%v", task.TaskId)
//	//task.DownloadVersionList = allVersion.AllVersion
//
//	data, err := json.Marshal(task)
//	_, err = w.Write(data)
//	if err != nil {
//		log.Printf("scriptDownloadListHandler failed, Write err:%v", err)
//		return
//	}
//
//	log.Printf("scriptDownloadListHandler success")
//}

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

func scriptResponse(c *gin.Context, code int, msg string, taskId string) {
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

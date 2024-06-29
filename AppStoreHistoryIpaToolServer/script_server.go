package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var task IPATaskDetail = IPATaskDetail{TaskId: "20240620_123456789_10003", AppleID: "fsi2clsfiz8@163.com", ApplePwd: "Ls112211", AppID: "583376064", Country: "JP", DownloadVersionList: nil}
var allVersion ipaVersionsRequest

func scriptTaskIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("scriptTaskIdHandler in")
	defer log.Printf("scriptTaskIdHandler out")

	task.TaskId = makeTaskId()

	data := []byte(task.TaskId)

	log.Printf("scriptTaskIdHandler:%v", task.TaskId)

	_, err := w.Write(data)
	if err != nil {
		log.Printf("scriptTaskIdHandler failed, Write err:%v", err)
		return
	}

	log.Printf("scriptTaskIdHandler success")
}

func scriptLoginInfoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("scriptLoginInfoHandler in")
	defer log.Printf("scriptLoginInfoHandler out")

	request := loginInfoRequestScript{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("scriptLoginInfoHandler failed, param err:%v", err)
		return
	}

	if request.TaskId != task.TaskId {
		log.Printf("scriptLoginInfoHandler failed, task id not find:%v", request.TaskId)
		return
	}

	log.Printf("scriptLoginInfoHandler:%v", task.TaskId)

	data, err := json.Marshal(task)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("scriptLoginInfoHandler failed, Write err:%v", err)
		return
	}

	log.Printf("scriptLoginInfoHandler success")
}

func script2FAHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("script2FAHandler in")
	defer log.Printf("script2FAHandler out")

	request := twoFAInfoRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)

	log.Printf("script2FAHandler AppleID:%v ,TaskID:%v", request.AppleID, request.TaskID)

	filePath := "C:\\Users\\xuaofei\\Desktop\\AppStoreHistoryIpa\\db.txt"
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("script2FAHandler failed ,读取文件失败：%v", err)
		return
	}

	log.Printf("script2FAHandler read file content:%v", string(content))

	response := twoFAInfoResponse{}
	response.TwoFACode = string(content)

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("script2FAHandler failed, Marshal taskInfo err:%v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Printf("script2FAHandler failed, Write err:%v", err)
		return
	}
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

func scriptReportResultHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("scriptReportResultHandler in")
	defer log.Printf("scriptReportResultHandler out")

	request := reportResultRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("scriptReportResultHandler failed, param err:%v", err)
		return
	}

	log.Printf("scriptReportResultHandler taskid:%v code:%d msg:%v", request.TaskId, request.Code, request.Msg)
	w.WriteHeader(200)
}

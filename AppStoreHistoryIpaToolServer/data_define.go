package main

import (
	"time"
)

// web status
const (
	SUCCESS = 0
	FAILED  = 1
)

//// login status
//const (
//	EMAIL_UNLOGIN    = 0
//	EMAIL_LOGINING   = 1
//	EMAIL_LOGINED    = 2
//	EMAIL_LOGINERROR = 3
//)
//
//// 2fa status
//const (
//	TFA_UNLOGIN     = 0
//	TFA_LOGINING    = 1
//	TFA_LOGINED     = 2
//	TFA_LOGIN_ERROR = 3
//)

// task status
const (
	TASK_INIT = iota
	TASK_INPUT_LOGIN_INFO
	TASK_LOGINING
	TASK_LOGINED
	TASK_LOGINING_2FA
	TASK_LOGINED_2FA
	TASK_SELECTED_APP
	TASK_GET_VERSION_LIST
	TASK_SELECT_VERSION_LIST
	TASK_START_DOWNLOAD_IPA
	TASK_FINISHED_DOWNLOAD_IPA
)

// 请求任务
type taskInfoRequest struct {
	AppleID     string `json:"apple_id"`
	ApplePwd    string `json:"apple_pwd"`
	AppBundleID string `json:"app_bundle_id"`
	Country     string `json:"country"`
	TaskID      string `json:"task_id"`
}

type scriptCommomRequest struct {
	TaskId string `json:"task_id"`
}

type scriptCommonResponse struct {
	AppleID     string `json:"apple_id"`
	ApplePwd    string `json:"apple_pwd"`
	AppBundleID string `json:"app_bundle_id"`
	AppId       string `json:"app_id"`
	Country     string `json:"country"`
	TaskID      string `json:"task_id"`
	Command     string `json:"command"`
	AllVersion  string `json:"all_version"`
}

type scriptLoginResultRequest struct {
	TaskId string `json:"task_id"`
	Result int    `json:"result"`
}

// 上报ipa版本信息
type scriptUploadAllVersionsRequest struct {
	TaskId       string `json:"task_id"`
	VersionCount int    `json:"version_count"`
	AllVersion   []struct {
		AppVer                   string `json:"app_ver"`
		AppVerID                 string `json:"app_ver_id"`
		BundleShortVersionString string `json:"bundle_short_version_string"`
	} `json:"all_version"`
}

// 要下载的ipa版本
type downloadIpaVersionsRequest struct {
	TaskId string `json:"task_id"`
}

type downloadIpaVersionsResponse struct {
	DownloadVersionList struct {
		AppVer                   string `json:"app_ver"`
		AppVerID                 string `json:"app_ver_id"`
		BundleShortVersionString string `json:"bundle_short_version_string"`
	} `json:"download_version_list"`
}

type reportResultRequest struct {
	TaskId string `json:"task_id"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
}

type taskIdRequestScript struct {
	//TaskId string `json:"task_id"`
}

type taskIdResponseScript struct {
	TaskId string `json:"task_id"`
}

type IPATaskDetail struct {
	TaskId string `json:"task_id"`

	AppleID     string `json:"apple_id"`
	ApplePwd    string `json:"apple_pwd"`
	TwoFaCode   string `json:"two_fa_code"`
	AppBundleID string `json:"app_bundle_id"`
	AppID       string `json:"app_id"`
	Country     string `json:"country"`

	DownloadVersionList []struct {
		AppVer                   string `json:"app_ver"`
		AppVerID                 string `json:"app_ver_id"`
		BundleShortVersionString string `json:"bundle_short_version_string"`
	} `json:"download_version_list"`
}

type ITunesSearchAppResult struct {
	ResultCount int `json:"resultCount"`
	Results     []struct {
		Features                           []string      `json:"features"`
		SupportedDevices                   []string      `json:"supportedDevices"`
		Advisories                         []interface{} `json:"advisories"`
		IsGameCenterEnabled                bool          `json:"isGameCenterEnabled"`
		ScreenshotUrls                     []string      `json:"screenshotUrls"`
		IpadScreenshotUrls                 []string      `json:"ipadScreenshotUrls"`
		AppletvScreenshotUrls              []interface{} `json:"appletvScreenshotUrls"`
		ArtworkUrl60                       string        `json:"artworkUrl60"`
		ArtworkUrl512                      string        `json:"artworkUrl512"`
		ArtworkUrl100                      string        `json:"artworkUrl100"`
		ArtistViewUrl                      string        `json:"artistViewUrl"`
		Kind                               string        `json:"kind"`
		ReleaseNotes                       string        `json:"releaseNotes"`
		ArtistId                           int           `json:"artistId"`
		ArtistName                         string        `json:"artistName"`
		Genres                             []string      `json:"genres"`
		Price                              float32       `json:"price"`
		Description                        string        `json:"description"`
		GenreIds                           []string      `json:"genreIds"`
		ReleaseDate                        time.Time     `json:"releaseDate"`
		BundleId                           string        `json:"bundleId"`
		PrimaryGenreName                   string        `json:"primaryGenreName"`
		PrimaryGenreId                     int           `json:"primaryGenreId"`
		CurrentVersionReleaseDate          time.Time     `json:"currentVersionReleaseDate"`
		IsVppDeviceBasedLicensingEnabled   bool          `json:"isVppDeviceBasedLicensingEnabled"`
		SellerName                         string        `json:"sellerName"`
		TrackId                            int           `json:"trackId"`
		TrackName                          string        `json:"trackName"`
		Currency                           string        `json:"currency"`
		MinimumOsVersion                   string        `json:"minimumOsVersion"`
		AverageUserRatingForCurrentVersion float64       `json:"averageUserRatingForCurrentVersion"`
		AverageUserRating                  float64       `json:"averageUserRating"`
		TrackCensoredName                  string        `json:"trackCensoredName"`
		LanguageCodesISO2A                 []string      `json:"languageCodesISO2A"`
		FileSizeBytes                      string        `json:"fileSizeBytes"`
		SellerUrl                          string        `json:"sellerUrl"`
		FormattedPrice                     string        `json:"formattedPrice"`
		ContentAdvisoryRating              string        `json:"contentAdvisoryRating"`
		UserRatingCountForCurrentVersion   int           `json:"userRatingCountForCurrentVersion"`
		TrackViewUrl                       string        `json:"trackViewUrl"`
		TrackContentRating                 string        `json:"trackContentRating"`
		Version                            string        `json:"version"`
		WrapperType                        string        `json:"wrapperType"`
		UserRatingCount                    int           `json:"userRatingCount"`
	} `json:"results"`
}

type AppData struct {
	ArtworkUrl60   string  `json:"artworkUrl60"`
	ReleaseNotes   string  `json:"releaseNotes"`
	Price          float32 `json:"price"`
	Description    string  `json:"description"`
	BundleId       string  `json:"bundleId"`
	SellerName     string  `json:"sellerName"`
	TrackId        int     `json:"trackId"`
	TrackName      string  `json:"trackName"`
	FileSizeBytes  string  `json:"fileSizeBytes"`
	FormattedPrice string  `json:"formattedPrice"`
	TrackViewUrl   string  `json:"trackViewUrl"`
	Version        string  `json:"version"`
}

type WebSearchAppResult struct {
	ResultCount int       `json:"resultCount"`
	Results     []AppData `json:"results"`
}

type Task struct {
	ID              int     `json:"id"`
	TaskID          string  `json:"task_id"`
	AppleID         string  `json:"apple_id"`
	Password        string  `json:"password"`
	LoginStatus     int     `json:"login_status"`
	TFA             string  `json:"tfa"`
	TFAStatus       int     `json:"tfa_status"`
	AppBoundID      string  `json:"app_boundId"`
	AppID           string  `json:"app_id"`
	Price           float64 `json:"price"`
	AllVersion      string  `json:"all_version"`
	DownloadVersion string  `json:"download_version"`
	Duration        int     `json:"duration"`
	SaveDirectory   string  `json:"save_directory"`
	StartTime       string  `json:"start_time"`
	TaskStatus      int     `json:"task_status"`
	Completed       int     `json:"completed"`
}

type WebServerCommonResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	TaskInfo Task   `json:"taskInfo"`
}

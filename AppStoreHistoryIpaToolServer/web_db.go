package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

// Database 结构体用于封装数据库操作
type Database struct {
	db *sql.DB
}

var instanceDB *Database
var onceDB sync.Once

func GetDBInstance() *Database {
	onceDB.Do(func() {
		instanceDB = &Database{}
	})
	return instanceDB
}

// CreateTable 创建表
func (d *Database) CreateTable() error {
	db, err := sql.Open("sqlite3", "./appstore.db")
	if err != nil {
		log.Printf("open db failed:%v", err)
		return err
	}

	d.db = db

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id TEXT NOT NULL,
		apple_id TEXT DEFAULT "",
		password TEXT DEFAULT "",
		login_status INTEGER DEFAULT 0,
		tfa TEXT DEFAULT "",
		tfa_status INTEGER DEFAULT 0,
		app_boundid TEXT DEFAULT "",
		app_id TEXT DEFAULT "",
		price REAL DEFAULT 0,
		all_version TEXT DEFAULT "",
		download_version TEXT DEFAULT "",
		duration INTEGER DEFAULT 0,
		save_directory TEXT DEFAULT "",
		start_time TEXT DEFAULT "",
		task_status INTEGER DEFAULT 0,
		completed INTEGER DEFAULT 0
	);`

	_, err = d.db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Exec db failed:%v", err)
		return err
	}

	log.Println("Table created successfully")
	return nil
}

// Close 关闭数据库连接
func (d *Database) Close() {
	if d.db != nil {
		_ = d.db.Close()
	}
}

// InsertTask
func (d *Database) InsertTask(taskID string) error {
	log.Printf("InsertTask in Taskid:%v", taskID)
	defer log.Printf("InsertTask out Taskid:%v", taskID)

	insertSQL := `
	INSERT INTO tasks (task_id)
	VALUES (?);`

	_, err := d.db.Exec(insertSQL, taskID)
	if err != nil {
		return err
	}

	log.Printf("Task inserted successfully")
	return nil
}

// UpdateTask
func (d *Database) UpdateAppleIDAndPassword(taskID, appleID, password string) error {
	log.Printf("UpdateAppleIDAndPassword Taskid:%v, appleID:%v, password len:%d", taskID, appleID, len(password))

	updateSQL := `UPDATE tasks SET apple_id = ?, password = ? WHERE task_id = ?;`

	_, err := d.db.Exec(updateSQL, appleID, password, taskID)
	if err != nil {
		log.Printf("UpdateAppleIDAndPassword db failed:%v", err)
		return err
	}

	log.Printf("UpdateAppleIDAndPassword inserted successfully")
	return nil
}

//func (d *Database) UpdateLoginStatus(taskID string, loginStatus int) error {
//	log.Printf("UpdateLoginStatus Taskid:%v, loginStatus:%v", taskID, loginStatus)
//
//	updateSQL := `UPDATE tasks SET login_status = ? WHERE task_id = ?;`
//
//	_, err := d.db.Exec(updateSQL, loginStatus, taskID)
//	if err != nil {
//		log.Printf("UpdateLoginStatus db failed:%v", err)
//		return err
//	}
//
//	log.Printf("UpdateLoginStatus inserted successfully")
//	return nil
//}

func (d *Database) Update2FA(taskID, tfa string) error {
	log.Printf("Update2FA Taskid:%v, tfa:%v", taskID, tfa)

	updateSQL := `UPDATE tasks SET tfa = ? WHERE task_id = ?;`

	_, err := d.db.Exec(updateSQL, tfa, taskID)
	if err != nil {
		log.Printf("Update2FA db failed:%v", err)
		return err
	}

	log.Printf("Update2FA inserted successfully")
	return nil
}

func (d *Database) Update2FAStatus(taskID string, tfaStatus int) error {
	log.Printf("UpdateLoginStatus Taskid:%v, tfaStatus:%v", taskID, tfaStatus)

	updateSQL := `UPDATE tasks SET tfa_status = ? WHERE task_id = ?;`

	_, err := d.db.Exec(updateSQL, tfaStatus, taskID)
	if err != nil {
		log.Printf("Update2FAStatus db failed:%v", err)
		return err
	}

	log.Printf("Update2FAStatus inserted successfully")
	return nil
}

func (d *Database) UpdateDownloadIpaInfo(taskID string, appid string, price float32) error {
	log.Printf("UpdateDownloadIpaInfo Taskid:%v, appid:%v, price:%v", taskID, appid, price)

	updateSQL := `UPDATE tasks SET app_id = ?, price = ? WHERE task_id = ?;`

	_, err := d.db.Exec(updateSQL, appid, price, taskID)
	if err != nil {
		log.Printf("UpdateDownloadIpaInfo db failed:%v", err)
		return err
	}

	log.Printf("UpdateDownloadIpaInfo inserted successfully")
	return nil
}

func (d *Database) UpdateAllVersion(taskID string, allVersion string) error {
	log.Printf("UpdateAllVersion Taskid:%v, allVersion:%v", taskID, allVersion)

	updateSQL := `UPDATE tasks SET all_version = ? WHERE task_id = ?;`

	_, err := d.db.Exec(updateSQL, allVersion, taskID)
	if err != nil {
		log.Printf("UpdateAllVersion db failed:%v", err)
		return err
	}

	log.Printf("UpdateAllVersion inserted successfully")
	return nil
}

func (d *Database) UpdateDownloadVersion(taskID string, downloadVersion string) error {
	log.Printf("UpdateAllVersion Taskid:%v, downloadVersion:%v", taskID, downloadVersion)

	updateSQL := `UPDATE tasks SET download_version = ? WHERE task_id = ?;`

	_, err := d.db.Exec(updateSQL, downloadVersion, taskID)
	if err != nil {
		log.Printf("UpdateDownloadVersion db failed:%v", err)
		return err
	}

	log.Printf("UpdateDownloadVersion inserted successfully")
	return nil
}

func (d *Database) UpdateTaskStatus(taskID string, status int) error {
	log.Printf("UpdateAllVersion Taskid:%v, status:%v", taskID, status)

	updateSQL := `UPDATE tasks SET task_status = ? WHERE task_id = ?;`

	_, err := d.db.Exec(updateSQL, status, taskID)
	if err != nil {
		log.Printf("UpdateTaskStatus db failed:%v", err)
		return err
	}

	log.Printf("UpdateTaskStatus inserted successfully")
	return nil
}

// QueryTask
func (d *Database) QueryNotStartedTask() (string, error) {
	querySQL := `SELECT task_id FROM tasks WHERE task_status = 1 LIMIT 1;`

	var taskID string
	err := d.db.QueryRow(querySQL).Scan(&taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No task found with login_status 0")
			return "", fmt.Errorf("no task found with login_status 0")
		}
		log.Printf("GetTaskWithLoginStatusZero db failed:%v", err)
		return "", err
	}

	log.Printf("Task found with login_status 0: %v", taskID)
	return taskID, nil
}

func (d *Database) QueryLoginStatus(taskID string) (int, error) {
	log.Printf("QueryLoginStatus Taskid:%v", taskID)

	querySQL := `SELECT login_status FROM tasks WHERE task_id = ?;`

	var loginStatus int
	err := d.db.QueryRow(querySQL, taskID).Scan(&loginStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No task found with taskID: %v", taskID)
			return 0, fmt.Errorf("no task found with taskID: %v", taskID)
		}
		log.Printf("QueryLoginStatus db failed:%v", err)
		return 0, err
	}

	log.Printf("Task login status queried successfully:Taskid:%v loginStatus:%v", taskID, loginStatus)
	return loginStatus, nil
}

func (d *Database) QueryVerifyCodeStatus(taskID string) (int, error) {
	log.Printf("QueryVerifyCodeStatus Taskid:%v", taskID)

	querySQL := `SELECT tfa_status FROM tasks WHERE task_id = ?;`

	var verifyCodeStatus int
	err := d.db.QueryRow(querySQL, taskID).Scan(&verifyCodeStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No task found with taskID: %v", taskID)
			return 0, fmt.Errorf("no task found with taskID: %v", taskID)
		}
		log.Printf("QueryVerifyCodeStatus db failed:%v", err)
		return 0, err
	}

	log.Printf("Task VerifyCode queried successfully:Taskid:%v verifyCodeStatus:%v", taskID, verifyCodeStatus)
	return verifyCodeStatus, nil
}

func (d *Database) Query2FA(taskID string) (string, error) {
	querySQL := `SELECT tfa FROM tasks WHERE task_id = ?;`

	var tfa string
	err := d.db.QueryRow(querySQL, taskID).Scan(&tfa)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No task found with taskID: %v", taskID)
			return "", nil
		}
		log.Printf("Query2FA db failed:%v", err)
		return "", err
	}

	log.Printf("Task found with taskID %v: tfa=%v", taskID, tfa)
	return tfa, nil
}

func (d *Database) QueryAppleIDAndPassword(taskID string) (string, string, error) {
	querySQL := `SELECT apple_id, password FROM tasks WHERE task_id = ?;`

	var appleID, password string
	err := d.db.QueryRow(querySQL, taskID).Scan(&appleID, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No task found with taskID: %v", taskID)
			return "", "", fmt.Errorf("no task found with taskID: %v", taskID)
		}
		log.Printf("QueryAppleIDAndPassword db failed:%v", err)
		return "", "", err
	}

	log.Printf("Task found with taskID %v: apple_id=%v, password=%v", taskID, appleID, password)
	return appleID, password, nil
}

func (d *Database) QueryPrice(taskID string) (float32, error) {
	querySQL := `SELECT price FROM tasks WHERE task_id = ?;`

	var price float32
	err := d.db.QueryRow(querySQL, taskID).Scan(&price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No task found with taskID: %v", taskID)
			return price, fmt.Errorf("no task found with taskID: %v", taskID)
		}
		log.Printf("QueryAppleIDAndPassword db failed:%v", err)
		return price, err
	}

	log.Printf("Task found with taskID %v: price=%v", taskID, price)
	return price, nil
}

func (d *Database) QueryAppId(taskID string) (string, error) {
	querySQL := `SELECT app_id FROM tasks WHERE task_id = ?;`

	var app_id string
	err := d.db.QueryRow(querySQL, taskID).Scan(&app_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No task found with taskID: %v", taskID)
			return "", fmt.Errorf("no task found with taskID: %v", taskID)
		}
		log.Printf("QueryAllVersion db failed:%v", err)
		return "", err
	}

	log.Printf("QueryAllVersion taskID %v: app_id=%v", taskID, app_id)
	return app_id, nil
}

func (d *Database) QueryAllVersion(taskID string) (string, error) {
	querySQL := `SELECT all_version FROM tasks WHERE task_id = ?;`

	var all_version string
	err := d.db.QueryRow(querySQL, taskID).Scan(&all_version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No task found with taskID: %v", taskID)
			return "", fmt.Errorf("no task found with taskID: %v", taskID)
		}
		log.Printf("QueryAllVersion db failed:%v", err)
		return "", err
	}

	log.Printf("QueryAllVersion taskID %v: QueryAllVersion=%v", taskID, all_version)
	return all_version, nil
}

func (d *Database) QueryTaskData(taskID string) (*Task, error) {
	var task Task
	query := `SELECT id, task_id, apple_id, password, login_status, tfa, tfa_status, app_boundid, app_id, price, all_version, download_version, duration, save_directory, start_time, task_status, completed FROM tasks WHERE task_id = ?`

	err := d.db.QueryRow(query, taskID).Scan(
		&task.ID,
		&task.TaskID,
		&task.AppleID,
		&task.Password,
		&task.LoginStatus,
		&task.TFA,
		&task.TFAStatus,
		&task.AppBoundID,
		&task.AppID,
		&task.Price,
		&task.AllVersion,
		&task.DownloadVersion,
		&task.Duration,
		&task.SaveDirectory,
		&task.StartTime,
		&task.TaskStatus,
		&task.Completed,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No task found with taskID: %v", taskID)
			return nil, fmt.Errorf("no task found with taskID: %v", taskID)
		}
		log.Printf("QueryTaskData db failed: %v", err)
		return nil, err
	}

	log.Printf("QueryTaskData taskID %v: Task=%+v", taskID, task)
	return &task, nil
}

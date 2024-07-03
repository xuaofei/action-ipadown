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
		apple_id TEXT,
		password TEXT,
		login_status INTEGER DEFAULT 0,
		tfa TEXT DEFAULT "",
		tfa_status INTEGER DEFAULT 0,
		app_boundid TEXT,
		app_id TEXT,
		price REAL DEFAULT 0,
		all_version TEXT,
		download_version TEXT,
		duration INTEGER,
		save_directory TEXT,
		start_time TEXT,
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
func (d *Database) UpdateTaskLoginInfo(taskID, appleID, password string) error {
	log.Printf("UpdateTaskLoginInfo Taskid:%v, appleID:%v, password len:%d", taskID, appleID, len(password))

	updateSQL := `UPDATE tasks SET apple_id = ?, password = ? WHERE task_id = ?;`

	_, err := d.db.Exec(updateSQL, appleID, password, taskID)
	if err != nil {
		log.Printf("UpdateTaskLoginInfo db failed:%v", err)
		return err
	}

	log.Printf("UpdateTaskLoginInfo inserted successfully")
	return nil
}

func (d *Database) UpdateTaskLoginStatus(taskID string, loginStatus int) error {
	log.Printf("UpdateTaskLoginStatus Taskid:%v, loginStatus:%v", taskID, loginStatus)

	updateSQL := `UPDATE tasks SET login_status = ? WHERE task_id = ?;`

	_, err := d.db.Exec(updateSQL, loginStatus, taskID)
	if err != nil {
		log.Printf("UpdateTaskLoginStatus db failed:%v", err)
		return err
	}

	log.Printf("UpdateTaskLoginStatus inserted successfully")
	return nil
}

func (d *Database) UpdateTask2FA(taskID, tfa string) error {
	log.Printf("UpdateTask2FA Taskid:%v, tfa:%v", taskID, tfa)

	updateSQL := `UPDATE tasks SET tfa = ? WHERE task_id = ?;`

	_, err := d.db.Exec(updateSQL, tfa, taskID)
	if err != nil {
		log.Printf("UpdateTask2FA db failed:%v", err)
		return err
	}

	log.Printf("UpdateTask2FA inserted successfully")
	return nil
}

func (d *Database) UpdateTask2FAStatus(taskID string, tfaStatus int) error {
	log.Printf("UpdateTaskLoginStatus Taskid:%v, tfaStatus:%v", taskID, tfaStatus)

	updateSQL := `UPDATE tasks SET tfa_status = ? WHERE task_id = ?;`

	_, err := d.db.Exec(updateSQL, tfaStatus, taskID)
	if err != nil {
		log.Printf("UpdateTask2FAStatus db failed:%v", err)
		return err
	}

	log.Printf("UpdateTask2FAStatus inserted successfully")
	return nil
}

func (d *Database) UpdateTaskDownloadIpaInfo(taskID string, appid string, price float32) error {
	log.Printf("UpdateTaskDownloadIpaInfo Taskid:%v, appid:%v, price:%v", taskID, appid, price)

	updateSQL := `UPDATE tasks SET app_id = ?, price = ? WHERE task_id = ?;`

	_, err := d.db.Exec(updateSQL, appid, price, taskID)
	if err != nil {
		log.Printf("UpdateTaskDownloadIpaInfo db failed:%v", err)
		return err
	}

	log.Printf("UpdateTaskDownloadIpaInfo inserted successfully")
	return nil
}

// QueryTask
func (d *Database) QueryNotStartedTask() (string, error) {
	querySQL := `SELECT task_id FROM tasks WHERE login_status = 0 LIMIT 1;`

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

func (d *Database) QueryTaskLoginStatus(taskID string) (int, error) {
	log.Printf("QueryTaskLoginStatus Taskid:%v", taskID)

	querySQL := `SELECT login_status FROM tasks WHERE task_id = ?;`

	var loginStatus int
	err := d.db.QueryRow(querySQL, taskID).Scan(&loginStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No task found with taskID: %v", taskID)
			return 0, fmt.Errorf("no task found with taskID: %v", taskID)
		}
		log.Printf("QueryTaskLoginStatus db failed:%v", err)
		return 0, err
	}

	log.Printf("Task login status queried successfully:Taskid:%v loginStatus:%v", taskID, loginStatus)
	return loginStatus, nil
}

func (d *Database) QueryTaskVerifyCodeStatus(taskID string) (int, error) {
	log.Printf("QueryTaskVerifyCodeStatus Taskid:%v", taskID)

	querySQL := `SELECT tfa_status FROM tasks WHERE task_id = ?;`

	var verifyCodeStatus int
	err := d.db.QueryRow(querySQL, taskID).Scan(&verifyCodeStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No task found with taskID: %v", taskID)
			return 0, fmt.Errorf("no task found with taskID: %v", taskID)
		}
		log.Printf("QueryTaskVerifyCodeStatus db failed:%v", err)
		return 0, err
	}

	log.Printf("Task VerifyCode queried successfully:Taskid:%v verifyCodeStatus:%v", taskID, verifyCodeStatus)
	return verifyCodeStatus, nil
}

// Get
func (d *Database) GetTFAByTaskID(taskID string) (string, error) {
	querySQL := `SELECT tfa FROM tasks WHERE task_id = ?;`

	var tfa string
	err := d.db.QueryRow(querySQL, taskID).Scan(&tfa)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No task found with taskID: %v", taskID)
			return "", nil
		}
		log.Printf("GetTFAByTaskID db failed:%v", err)
		return "", err
	}

	log.Printf("Task found with taskID %v: tfa=%v", taskID, tfa)
	return tfa, nil
}

func (d *Database) GetAppleIDAndPasswordByTaskID(taskID string) (string, string, error) {
	querySQL := `SELECT apple_id, password FROM tasks WHERE task_id = ?;`

	var appleID, password string
	err := d.db.QueryRow(querySQL, taskID).Scan(&appleID, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No task found with taskID: %v", taskID)
			return "", "", fmt.Errorf("no task found with taskID: %v", taskID)
		}
		log.Printf("GetAppleIDAndPasswordByTaskID db failed:%v", err)
		return "", "", err
	}

	log.Printf("Task found with taskID %v: apple_id=%v, password=%v", taskID, appleID, password)
	return appleID, password, nil
}

package main

import (
	"fmt"
	"time"
)

func makeTaskId() string {
	unix := time.Now()
	return fmt.Sprintf(unix.Format("2006.01.02.15.04.05.999999999"))
}

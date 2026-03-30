package debug

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	enabled bool
	file    *os.File
	mu      sync.Mutex
)

func Init() error {
	if os.Getenv("GH_DEBUG") == "" {
		return nil
	}

	enabled = true
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	logPath := filepath.Join(cwd, "ghr-debug.log")
	file, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Close() {
	if file != nil {
		file.Close()
	}
}

func Print(format string, args ...interface{}) {
	if !enabled || file == nil {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	fmt.Fprintf(file, "[%s] "+format+"\n", append([]interface{}{timestamp}, args...)...)
}

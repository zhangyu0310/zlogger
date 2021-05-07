// Package zlogger is a go log package
/*
   Auto print file name & line at call point.
   Can split log file automatically. (Prevent a single log file from being too large)
*/
package zlogger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

// DefaultLogger is a default logger for sample function.
var DefaultLogger *Logger

// Logger contain log of go & file handler.
// Use logger.xxx() to log & set right prefix.
type Logger struct {
	logger   *log.Logger // Use Logger of golang
	file     *os.File    // File handler of Logger
	Path     string      // The path of Logger
	Name     string      // The name of Logger without day
	FileName string      // The name of Logger with day info
	mutex    sync.Mutex  // Mutex for logger (Prefix order & update file safe)
}

// New create a new logger handler.
// @path: dir of logs.
// @name: prefix of logs.
// Log file name just have year-month-day
// Time of logs record is microseconds.
func New(path, name string, autoUpdate bool) (*Logger, error) {
	l := Logger{
		Path: path,
		Name: name,
	}
	l.FileName = name + "_" + time.Now().Format("2006-01-02")
	filePath := l.Path + "/" + l.FileName

	var err error
	l.file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	l.logger = log.New(l.file, "", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)
	if l.logger != nil && autoUpdate {
		go func() {
			// Check time and update logger file.
			for true {
				if l.FileName != l.Name+"_"+time.Now().Format("2006-01-02") {
					if err := UpdateLoggerFile(&l); err != nil {
						l.Error("Update logger file failed.", err)
						break
					}
				}
				time.Sleep(time.Minute * 10)
			}
		}()
	}
	return &l, nil
}

// UpdateLoggerFile update the log file name. (Date suffix)
func UpdateLoggerFile(logger *Logger) error {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	logger.FileName = logger.Name + "_" + time.Now().Format("2006-01-02")
	filePath := logger.Path + "/" + logger.FileName
	// Create new file handler & new logger
	nFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	nLogger := log.New(nFile, "", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)
	// Set new file handler/logger to old logger & close old file handler.
	logger.logger = nLogger
	oldFileHandler := logger.file
	logger.file = nFile
	if err := oldFileHandler.Close(); err != nil {
		logger.Error("Old logger file handler close failed.", err)
	}
	return nil
}

// getFileAndLinePrefix get file name & line of call function.
func getFileAndLinePrefix(depth int) string {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}
	// I don't like log path. Use short.
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return fmt.Sprintf("%s:%d: ", short, line)
}

func (logger *Logger) Info(msg ...interface{}) {
	logger.InfoN(3, msg...)
}

func (logger *Logger) InfoN(n int, msg ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	prefix := getFileAndLinePrefix(n)
	logger.logger.SetPrefix(prefix + "[INFO] ")
	logger.logger.Println(msg...)
}

func (logger *Logger) Debug(msg ...interface{}) {
	logger.DebugN(3, msg...)
}

func (logger *Logger) DebugN(n int, msg ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	prefix := getFileAndLinePrefix(n)
	logger.logger.SetPrefix(prefix + "[DEBUG] ")
	logger.logger.Println(msg...)
}

func (logger *Logger) Warn(msg ...interface{}) {
	logger.WarnN(3, msg...)
}

func (logger *Logger) WarnN(n int, msg ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	prefix := getFileAndLinePrefix(n)
	logger.logger.SetPrefix(prefix + "[WARN] ")
	logger.logger.Println(msg...)
}

func (logger *Logger) Error(msg ...interface{}) {
	logger.ErrorN(3, msg...)
}

func (logger *Logger) ErrorN(n int, msg ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	prefix := getFileAndLinePrefix(n)
	logger.logger.SetPrefix(prefix + "[ERROR] ")
	logger.logger.Println(msg...)
}

func (logger *Logger) Fatal(msg ...interface{}) {
	logger.FatalN(3, msg...)
}

func (logger *Logger) FatalN(n int, msg ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	prefix := getFileAndLinePrefix(n)
	logger.logger.SetPrefix(prefix + "[FATAL] ")
	logger.logger.Fatalln(msg...)
}

func (logger *Logger) Panic(msg ...interface{}) {
	logger.PanicN(3, msg...)
}

func (logger *Logger) PanicN(n int, msg ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	prefix := getFileAndLinePrefix(n)
	logger.logger.SetPrefix(prefix + "[PANIC] ")
	logger.logger.Panicln(msg...)
}

func Info(msg ...interface{}) {
	DefaultLogger.InfoN(3, msg...)
}

func Debug(msg ...interface{}) {
	DefaultLogger.DebugN(3, msg...)
}

func Warn(msg ...interface{}) {
	DefaultLogger.WarnN(3, msg...)
}

func Error(msg ...interface{}) {
	DefaultLogger.ErrorN(3, msg...)
}

func Fatal(msg ...interface{}) {
	DefaultLogger.FatalN(3, msg...)
}

func Panic(msg ...interface{}) {
	DefaultLogger.PanicN(3, msg...)
}

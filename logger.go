package zlogger

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

// Logger contain log of go & file handler.
// Use logger.xxx() to log & set right prefix.
type Logger struct {
	logger   *log.Logger // Use Logger of golang
	file     *os.File    // File handler of Logger
	Path     string      // The path of Logger
	Name     string      // The name of Logger without day
	FileName string      // The name of Logger with day info
}

// @path: dir of logs.
// @name: prefix of logs.
// Log file name just have year-month-day
// Time of logs record is microseconds.
func New(path string, name string) *Logger {
	l := Logger{
		Path: path,
		Name: name,
	}
	l.FileName = name + "_" + time.Now().Format("2006-01-02")
	filePath := l.Path + "/" + l.FileName

	var err error
	l.file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}
	l.logger = log.New(l.file, "", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)
	if l.logger != nil {
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
	return &l
}

func UpdateLoggerFile(logger *Logger) error {
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

// GetFileAndLinePrefix get file name & line of call function.
func GetFileAndLinePrefix() string {
	_, file, line, ok := runtime.Caller(2)
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
	file = short
	return file + ":" + strconv.Itoa(line) + ": "
}

func (logger *Logger) Info(msg ...interface{}) {
	prefix := GetFileAndLinePrefix()
	logger.logger.SetPrefix(prefix + "[INFO] ")
	logger.logger.Println(msg...)
}

func (logger *Logger) Debug(msg ...interface{}) {
	prefix := GetFileAndLinePrefix()
	logger.logger.SetPrefix(prefix + "[DEBUG] ")
	logger.logger.Println(msg...)
}

func (logger *Logger) Warn(msg ...interface{}) {
	prefix := GetFileAndLinePrefix()
	logger.logger.SetPrefix(prefix + "[WARN] ")
	logger.logger.Println(msg...)
}

func (logger *Logger) Error(msg ...interface{}) {
	prefix := GetFileAndLinePrefix()
	logger.logger.SetPrefix(prefix + "[ERROR] ")
	logger.logger.Println(msg...)
}

func (logger *Logger) Fatal(msg ...interface{}) {
	prefix := GetFileAndLinePrefix()
	logger.logger.SetPrefix(prefix + "[FATAL] ")
	logger.logger.Fatalln(msg...)
}

func (logger *Logger) Panic(msg ...interface{}) {
	prefix := GetFileAndLinePrefix()
	logger.logger.SetPrefix(prefix + "[PANIC] ")
	logger.logger.Panicln(msg...)
}

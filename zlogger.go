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
	"sync/atomic"
	"time"
)

// Only print log which bigger than set.
// We don't need trace!!!
const (
	LogLevelAll   = 0
	LogLevelDebug = 1
	LogLevelInfo  = 2
	LogLevelWarn  = 3
	LogLevelError = 4
	LogLevelFatal = 5
	LogLevelPanic = 6
	LogLevelOff   = 7
)

// defaultLogger is a default logger for sample function.
var defaultLogger *Logger

// Logger contain log of go & file handler.
// Use logger.xxx() to log & set right prefix.
type Logger struct {
	logger     *log.Logger  // Use Logger of golang
	file       *os.File     // File handler of Logger
	Path       string       // The path of Logger
	Name       string       // The name of Logger without day
	FileName   string       // The name of Logger with day info
	close      chan bool    // The logger is closed
	autoUpdate bool         // logger can auto update log file
	logLevel   atomic.Value // The level of log to print
}

// New create a new logger handler.
// @path: dir of logs.
// @name: prefix of logs.
// Log file name just have year-month-day
// Time of logs record is microseconds.
func New(path, name string, autoUpdate bool, logLevel uint8) (err error) {
	if defaultLogger != nil {
		defaultLogger.Close()
	}
	defaultLogger, err = NewInternal(path, name, autoUpdate, logLevel)
	return
}

// defaultNew create a default logger handler.
// If defaultLogger is nil and log function been called, this function will be called.
func defaultNew() {
	err := New("./", "zlogger", false, LogLevelAll)
	if err != nil {
		panic(err)
	}
}

func getLogFileName(name string) string {
	return name + "." + time.Now().Format("2006-01-02_15")
}

// NewInternal the implement of New
func NewInternal(path, name string, autoUpdate bool, logLevel uint8) (*Logger, error) {
	l := Logger{
		Path:       path,
		Name:       name,
		close:      make(chan bool, 0),
		autoUpdate: autoUpdate,
	}
	l.SetLogLevel(logLevel)
	l.FileName = getLogFileName(name)
	filePath := l.Path + "/" + l.FileName

	var err error
	l.file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	l.logger = log.New(l.file, "", log.LstdFlags|log.Lmicroseconds)
	if l.logger != nil && autoUpdate {
		go func() {
			// Check time and update logger file.
			t := time.NewTicker(time.Minute * 10)
			defer t.Stop()
			for {
				select {
				case <-l.close:
					return
				case <-t.C:
					if l.FileName != getLogFileName(l.Name) {
						if err := l.updateLoggerFile(); err != nil {
							l.Error("Update logger file failed.", err)
							break
						}
					}
				}
			}
		}()
	}
	return &l, nil
}

func ForceUpdateLoggerFile() error {
	return defaultLogger.updateLoggerFile()
}

// updateLoggerFile update the log file name. (Date suffix)
func (logger *Logger) updateLoggerFile() error {
	logger.FileName = getLogFileName(logger.Name)
	filePath := logger.Path + "/" + logger.FileName
	// Create new file handler & new logger
	nFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	// Set new file handler/logger to old logger & close old file handler.
	logger.logger.SetOutput(nFile)
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
	return fmt.Sprintf("%s:%d:", short, line)
}

func (logger *Logger) SetLogLevel(logLevel uint8) {
	logger.logLevel.Store(logLevel)
}

func (logger *Logger) GetLogLevel() uint8 {
	return logger.logLevel.Load().(uint8)
}

func getMsgSlice(prefix, logType string, msg []interface{}) []interface{} {
	msgLocal := make([]interface{}, 0, 4)
	msgLocal = append(msgLocal, prefix)
	msgLocal = append(msgLocal, logType)
	msgLocal = append(msgLocal, msg...)
	return msgLocal
}

func (logger *Logger) Debug(msg ...interface{}) {
	logger.DebugN(3, msg...)
}

func (logger *Logger) DebugF(format string, v ...interface{}) {
	logger.DebugNF(3, format, v...)
}

func (logger *Logger) DebugN(n int, msg ...interface{}) {
	if logger.GetLogLevel() > LogLevelDebug {
		return
	}
	prefix := getFileAndLinePrefix(n)
	msgLocal := getMsgSlice(prefix, "[DEBUG]", msg)
	logger.logger.Println(msgLocal...)
}

func (logger *Logger) DebugNF(n int, format string, v ...interface{}) {
	if logger.GetLogLevel() > LogLevelDebug {
		return
	}
	prefix := getFileAndLinePrefix(n)
	f := fmt.Sprintf("%s %s %s\n", prefix, "[DEBUG]", format)
	logger.logger.Printf(f, v...)
}

func (logger *Logger) Info(msg ...interface{}) {
	logger.InfoN(3, msg...)
}

func (logger *Logger) InfoF(format string, v ...interface{}) {
	logger.InfoNF(3, format, v...)
}

func (logger *Logger) InfoN(n int, msg ...interface{}) {
	if logger.GetLogLevel() > LogLevelInfo {
		return
	}
	prefix := getFileAndLinePrefix(n)
	msgLocal := getMsgSlice(prefix, "[INFO]", msg)
	logger.logger.Println(msgLocal...)
}

func (logger *Logger) InfoNF(n int, format string, v ...interface{}) {
	if logger.GetLogLevel() > LogLevelInfo {
		return
	}
	prefix := getFileAndLinePrefix(n)
	f := fmt.Sprintf("%s %s %s\n", prefix, "[INFO]", format)
	logger.logger.Printf(f, v...)
}

func (logger *Logger) Warn(msg ...interface{}) {
	logger.WarnN(3, msg...)
}

func (logger *Logger) WarnF(format string, v ...interface{}) {
	logger.WarnNF(3, format, v...)
}

func (logger *Logger) WarnN(n int, msg ...interface{}) {
	if logger.GetLogLevel() > LogLevelWarn {
		return
	}
	prefix := getFileAndLinePrefix(n)
	msgLocal := getMsgSlice(prefix, "[WARN]", msg)
	logger.logger.Println(msgLocal...)
}

func (logger *Logger) WarnNF(n int, format string, v ...interface{}) {
	if logger.GetLogLevel() > LogLevelWarn {
		return
	}
	prefix := getFileAndLinePrefix(n)
	f := fmt.Sprintf("%s %s %s\n", prefix, "[WARN]", format)
	logger.logger.Printf(f, v...)
}

func (logger *Logger) Error(msg ...interface{}) {
	logger.ErrorN(3, msg...)
}

func (logger *Logger) ErrorF(format string, v ...interface{}) {
	logger.ErrorNF(3, format, v...)
}

func (logger *Logger) ErrorN(n int, msg ...interface{}) {
	if logger.GetLogLevel() > LogLevelError {
		return
	}
	prefix := getFileAndLinePrefix(n)
	msgLocal := getMsgSlice(prefix, "[ERROR]", msg)
	logger.logger.Println(msgLocal...)
}

func (logger *Logger) ErrorNF(n int, format string, v ...interface{}) {
	if logger.GetLogLevel() > LogLevelError {
		return
	}
	prefix := getFileAndLinePrefix(n)
	f := fmt.Sprintf("%s %s %s\n", prefix, "[ERROR]", format)
	logger.logger.Printf(f, v...)
}

func (logger *Logger) Fatal(msg ...interface{}) {
	logger.FatalN(3, msg...)
}

func (logger *Logger) FatalF(format string, v ...interface{}) {
	logger.FatalNF(3, format, v...)
}

func (logger *Logger) FatalN(n int, msg ...interface{}) {
	if logger.GetLogLevel() > LogLevelFatal {
		return
	}
	prefix := getFileAndLinePrefix(n)
	msgLocal := getMsgSlice(prefix, "[FATAL]", msg)
	logger.logger.Fatalln(msgLocal...)
}

func (logger *Logger) FatalNF(n int, format string, v ...interface{}) {
	if logger.GetLogLevel() > LogLevelFatal {
		return
	}
	prefix := getFileAndLinePrefix(n)
	f := fmt.Sprintf("%s %s %s\n", prefix, "[FATAL]", format)
	logger.logger.Printf(f, v...)
}

func (logger *Logger) Panic(msg ...interface{}) {
	logger.PanicN(3, msg...)
}

func (logger *Logger) PanicF(format string, v ...interface{}) {
	logger.PanicNF(3, format, v...)
}

func (logger *Logger) PanicN(n int, msg ...interface{}) {
	if logger.GetLogLevel() > LogLevelPanic {
		return
	}
	prefix := getFileAndLinePrefix(n)
	msgLocal := getMsgSlice(prefix, "[PANIC]", msg)
	logger.logger.Panicln(msgLocal...)
}

func (logger *Logger) PanicNF(n int, format string, v ...interface{}) {
	if logger.GetLogLevel() > LogLevelFatal {
		return
	}
	prefix := getFileAndLinePrefix(n)
	f := fmt.Sprintf("%s %s %s\n", prefix, "[PANIC]", format)
	logger.logger.Printf(f, v...)
}

// Close stop update log file coroutine & close log file handler.
// You don't need to call this function on exit.
func (logger *Logger) Close() {
	_ = logger.file.Close()
	if logger.autoUpdate {
		logger.close <- true
		close(logger.close)
	}
}

func SetLogLevel(logLevel uint8) {
	if defaultLogger == nil {
		defaultNew()
	}
	defaultLogger.SetLogLevel(logLevel)
}

func GetLogLevel() uint8 {
	if defaultLogger == nil {
		defaultNew()
	}
	return defaultLogger.GetLogLevel()
}

func Debug(msg ...interface{}) {
	if defaultLogger == nil {
		defaultNew()
	}
	defaultLogger.DebugN(3, msg...)
}

func Info(msg ...interface{}) {
	if defaultLogger == nil {
		defaultNew()
	}
	defaultLogger.InfoN(3, msg...)
}

func Warn(msg ...interface{}) {
	if defaultLogger == nil {
		defaultNew()
	}
	defaultLogger.WarnN(3, msg...)
}

func Error(msg ...interface{}) {
	if defaultLogger == nil {
		defaultNew()
	}
	defaultLogger.ErrorN(3, msg...)
}

func Fatal(msg ...interface{}) {
	if defaultLogger == nil {
		defaultNew()
	}
	defaultLogger.FatalN(3, msg...)
}

func Panic(msg ...interface{}) {
	if defaultLogger == nil {
		defaultNew()
	}
	defaultLogger.PanicN(3, msg...)
}

func DebugF(format string, v ...interface{}) {
	if defaultLogger == nil {
		defaultNew()
	}
	defaultLogger.DebugNF(3, format, v...)
}

func InfoF(format string, v ...interface{}) {
	if defaultLogger == nil {
		defaultNew()
	}
	defaultLogger.InfoNF(3, format, v...)
}

func WarnF(format string, v ...interface{}) {
	if defaultLogger == nil {
		defaultNew()
	}
	defaultLogger.WarnNF(3, format, v...)
}

func ErrorF(format string, v ...interface{}) {
	if defaultLogger == nil {
		defaultNew()
	}
	defaultLogger.ErrorNF(3, format, v...)
}

func FatalF(format string, v ...interface{}) {
	if defaultLogger == nil {
		defaultNew()
	}
	defaultLogger.FatalNF(3, format, v...)
}

func PanicF(format string, v ...interface{}) {
	if defaultLogger == nil {
		defaultNew()
	}
	defaultLogger.PanicNF(3, format, v...)
}

func LogLevel2Str(level uint8) string {
	switch level {
	case LogLevelAll:
		return "all"
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	case LogLevelFatal:
		return "fatal"
	case LogLevelPanic:
		return "panic"
	case LogLevelOff:
		return "off"
	default:
		return "unknown"
	}
}

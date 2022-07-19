package zlogger

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMultiOpenAndWrite(t *testing.T) {
	l1, _ := NewInternal("./", "zlogger", false, LogLevelAll)
	l2, _ := NewInternal("./", "zlogger", false, LogLevelAll)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		for i := 0; i < 100; i++ {
			l1.Info("Info l3", i)
			l1.Debug("Debug l3", i)
			l1.Warn("Warn l3", i)
			l1.Error("Error l3", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 100; i++ {
			l2.Info("Info l4", i)
			l2.Debug("Debug l4", i)
			l2.Warn("Warn l4", i)
			l2.Error("Error l4", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 100; i++ {
			Info("Info default", i)
			Debug("Debug default", i)
			Warn("Warn default", i)
			Error("Error default", i)
		}
		wg.Done()
	}()
	wg.Wait()

	_ = os.Remove(defaultLogger.Path + defaultLogger.FileName)
}

func newForTest(t *testing.T) {
	err := New("./", "zlogger", true, LogLevelAll)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateLoggerFile(t *testing.T) {
	newForTest(t)
	var wg sync.WaitGroup
	wg.Add(9)
	begin := time.Now().UnixNano()
	go func() {
		for i := 0; i < 10000; i++ {
			Debug("Debug l1", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			Info("Info l1", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			Warn("Warn l1", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			Error("Error l1", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			DebugF("This is Debug %d", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			InfoF("This is Info %d", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			WarnF("This is Warn %d", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			ErrorF("This is Error %d", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			err := ForceUpdateLoggerFile()
			if err != nil {
				t.Error("updateLoggerFile failed.", err)
			}
		}
		wg.Done()
	}()
	wg.Wait()
	end := time.Now().UnixNano()
	t.Log("Time cost:", end-begin)

	f, _ := os.Open(defaultLogger.Path + defaultLogger.FileName)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Log("File close failed.", err)
		}
	}(f)
	reader := bufio.NewReader(f)
	debugSize := 0
	infoSize := 0
	warnSize := 0
	errorSize := 0
	debugFSize := 0
	infoFSize := 0
	warnFSize := 0
	errorFSize := 0
	for {
		line, _, err := reader.ReadLine()
		lineStr := string(line)
		if strings.Contains(lineStr, "DEBUG") {
			if strings.Contains(lineStr, "Debug l1") {
				if !strings.Contains(lineStr,
					fmt.Sprintf("Debug l1 %d", debugSize)) {
					t.Error("Log debug order is wrong!")
				}
				debugSize++
			}
			if strings.Contains(lineStr, "This is Debug") {
				if !strings.Contains(lineStr,
					fmt.Sprintf("This is Debug %d", debugFSize)) {
					t.Error("Log debugF order is wrong!")
				}
				debugFSize++
			}
		} else if strings.Contains(lineStr, "INFO") {
			if strings.Contains(lineStr, "Info l1") {
				if !strings.Contains(lineStr,
					fmt.Sprintf("Info l1 %d", infoSize)) {
					t.Error("Log info order is wrong!")
				}
				infoSize++
			}
			if strings.Contains(lineStr, "This is Info") {
				if !strings.Contains(lineStr,
					fmt.Sprintf("This is Info %d", infoFSize)) {
					t.Error("Log infoF order is wrong!")
				}
				infoFSize++
			}
		} else if strings.Contains(lineStr, "WARN") {
			if strings.Contains(lineStr, "Warn l1") {
				if !strings.Contains(lineStr,
					fmt.Sprintf("Warn l1 %d", warnSize)) {
					t.Error("Log warn order is wrong!")
				}
				warnSize++
			}
			if strings.Contains(lineStr, "This is Warn") {
				if !strings.Contains(lineStr,
					fmt.Sprintf("This is Warn %d", warnFSize)) {
					t.Error("Log warnF order is wrong!")
				}
				warnFSize++
			}
		} else if strings.Contains(lineStr, "ERROR") {
			if strings.Contains(lineStr, "Error l1") {
				if !strings.Contains(lineStr,
					fmt.Sprintf("Error l1 %d", errorSize)) {
					t.Error("Log error order is wrong!")
				}
				errorSize++
			}
			if strings.Contains(lineStr, "This is Error") {
				if !strings.Contains(lineStr,
					fmt.Sprintf("This is Error %d", errorFSize)) {
					t.Error("Log errorF order is wrong!")
				}
				errorFSize++
			}
		} else if lineStr == "" {
			t.Log("Empty log line")
		} else {
			t.Error("String Contains error")
		}
		if err != nil {
			break
		}
	}
	if debugSize != 10000 {
		t.Error("Debug size is", debugSize, "not 10000")
	}
	if infoSize != 10000 {
		t.Error("Info size is", infoSize, "not 10000")
	}
	if warnSize != 10000 {
		t.Error("Warn size is", warnSize, "not 10000")
	}
	if errorSize != 10000 {
		t.Error("Error size is", errorSize, "not 10000")
	}
	if debugFSize != 10000 {
		t.Error("DebugF size is", debugFSize, "not 10000")
	}
	if infoFSize != 10000 {
		t.Error("InfoF size is", infoFSize, "not 10000")
	}
	if warnFSize != 10000 {
		t.Error("WarnF size is", warnFSize, "not 10000")
	}
	if errorFSize != 10000 {
		t.Error("ErrorF size is", errorFSize, "not 10000")
	}

	t.Log("Check success.")

	_ = os.Remove(defaultLogger.Path + defaultLogger.FileName)
}

func readLine(reader *bufio.Reader) string {
	line, _, _ := reader.ReadLine()
	return string(line)
}

func checkResult(t *testing.T, level uint8) {
	f, _ := os.Open(defaultLogger.Path + defaultLogger.FileName)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Log("File close failed.", err)
		}
	}(f)

	debugMsg := fmt.Sprintf("level %s %s", LogLevel2Str(level), "debug")
	infoMsg := fmt.Sprintf("level %s %s", LogLevel2Str(level), "info")
	warnMsg := fmt.Sprintf("level %s %s", LogLevel2Str(level), "warn")
	errorMsg := fmt.Sprintf("level %s %s", LogLevel2Str(level), "error")

	debugFMsg := fmt.Sprintf("level %s %s", LogLevel2Str(level), "F debug F")
	infoFMsg := fmt.Sprintf("level %s %s", LogLevel2Str(level), "F info F")
	warnFMsg := fmt.Sprintf("level %s %s", LogLevel2Str(level), "F warn F")
	errorFMsg := fmt.Sprintf("level %s %s", LogLevel2Str(level), "F error F")

	reader := bufio.NewReader(f)
	switch level {
	case LogLevelAll:
		line := readLine(reader)
		if !strings.Contains(line, debugMsg) {
			t.Error("log level", LogLevel2Str(level), "write debug error")
		}
		line = readLine(reader)
		if !strings.Contains(line, debugFMsg) {
			t.Error("log level", LogLevel2Str(level), "write debugF error")
		}
		line = readLine(reader)
		if !strings.Contains(line, infoMsg) {
			t.Error("log level", LogLevel2Str(level), "write info error")
		}
		line = readLine(reader)
		if !strings.Contains(line, infoFMsg) {
			t.Error("log level", LogLevel2Str(level), "write infoF error")
		}
		line = readLine(reader)
		if !strings.Contains(line, warnMsg) {
			t.Error("log level", LogLevel2Str(level), "write warn error")
		}
		line = readLine(reader)
		if !strings.Contains(line, warnFMsg) {
			t.Error("log level", LogLevel2Str(level), "write warnF error")
		}
		line = readLine(reader)
		if !strings.Contains(line, errorMsg) {
			t.Error("log level", LogLevel2Str(level), "write error error")
		}
		line = readLine(reader)
		if !strings.Contains(line, errorFMsg) {
			t.Error("log level", LogLevel2Str(level), "write errorF error")
		}
	case LogLevelDebug:
		line := readLine(reader)
		if !strings.Contains(line, debugMsg) {
			t.Error("log level", LogLevel2Str(level), "write debug error")
		}
		line = readLine(reader)
		if !strings.Contains(line, debugFMsg) {
			t.Error("log level", LogLevel2Str(level), "write debugF error")
		}
		line = readLine(reader)
		if !strings.Contains(line, infoMsg) {
			t.Error("log level", LogLevel2Str(level), "write info error")
		}
		line = readLine(reader)
		if !strings.Contains(line, infoFMsg) {
			t.Error("log level", LogLevel2Str(level), "write infoF error")
		}
		line = readLine(reader)
		if !strings.Contains(line, warnMsg) {
			t.Error("log level", LogLevel2Str(level), "write warn error")
		}
		line = readLine(reader)
		if !strings.Contains(line, warnFMsg) {
			t.Error("log level", LogLevel2Str(level), "write warnF error")
		}
		line = readLine(reader)
		if !strings.Contains(line, errorMsg) {
			t.Error("log level", LogLevel2Str(level), "write error error")
		}
		line = readLine(reader)
		if !strings.Contains(line, errorFMsg) {
			t.Error("log level", LogLevel2Str(level), "write errorF error")
		}
	case LogLevelInfo:
		line := readLine(reader)
		if !strings.Contains(line, infoMsg) {
			t.Error("log level", LogLevel2Str(level), "write info error")
		}
		line = readLine(reader)
		if !strings.Contains(line, infoFMsg) {
			t.Error("log level", LogLevel2Str(level), "write infoF error")
		}
		line = readLine(reader)
		if !strings.Contains(line, warnMsg) {
			t.Error("log level", LogLevel2Str(level), "write warn error")
		}
		line = readLine(reader)
		if !strings.Contains(line, warnFMsg) {
			t.Error("log level", LogLevel2Str(level), "write warnF error")
		}
		line = readLine(reader)
		if !strings.Contains(line, errorMsg) {
			t.Error("log level", LogLevel2Str(level), "write error error")
		}
		line = readLine(reader)
		if !strings.Contains(line, errorFMsg) {
			t.Error("log level", LogLevel2Str(level), "write errorF error")
		}
	case LogLevelWarn:
		line := readLine(reader)
		if !strings.Contains(line, warnMsg) {
			t.Error("log level", LogLevel2Str(level), "write warn error")
		}
		line = readLine(reader)
		if !strings.Contains(line, warnFMsg) {
			t.Error("log level", LogLevel2Str(level), "write warnF error")
		}
		line = readLine(reader)
		if !strings.Contains(line, errorMsg) {
			t.Error("log level", LogLevel2Str(level), "write error error")
		}
		line = readLine(reader)
		if !strings.Contains(line, errorFMsg) {
			t.Error("log level", LogLevel2Str(level), "write errorF error")
		}
	case LogLevelError:
		line := readLine(reader)
		if !strings.Contains(line, errorMsg) {
			t.Error("log level", LogLevel2Str(level), "write error error")
		}
		line = readLine(reader)
		if !strings.Contains(line, errorFMsg) {
			t.Error("log level", LogLevel2Str(level), "write errorF error")
		}
	case LogLevelOff:
		line := readLine(reader)
		if line != "" {
			t.Error("log level", LogLevel2Str(level), "error")
		}
	}
}

func writeTestLog(level uint8) {
	SetLogLevel(level)
	prefix := fmt.Sprintf("level %s", LogLevel2Str(level))
	Debug(prefix, "debug")
	DebugF(prefix+" %s", "F debug F")
	Info(prefix, "info")
	InfoF(prefix+" %s", "F info F")
	Warn(prefix, "warn")
	WarnF(prefix+" %s", "F warn F")
	Error(prefix, "error")
	ErrorF(prefix+" %s", "F error F")
}

func TestSetLogLevel(t *testing.T) {
	newForTest(t)
	_ = os.Remove(defaultLogger.Path + defaultLogger.FileName)
	err := ForceUpdateLoggerFile()
	if err != nil {
		t.Fatal("Update log file failed.", err)
	}
	writeTestLog(LogLevelAll)
	checkResult(t, LogLevelAll)

	_ = os.Remove(defaultLogger.Path + defaultLogger.FileName)
	err = ForceUpdateLoggerFile()
	if err != nil {
		t.Fatal("Update log file failed.", err)
	}
	writeTestLog(LogLevelDebug)
	checkResult(t, LogLevelDebug)

	_ = os.Remove(defaultLogger.Path + defaultLogger.FileName)
	err = ForceUpdateLoggerFile()
	if err != nil {
		t.Fatal("Update log file failed.", err)
	}
	writeTestLog(LogLevelInfo)
	checkResult(t, LogLevelInfo)

	_ = os.Remove(defaultLogger.Path + defaultLogger.FileName)
	err = ForceUpdateLoggerFile()
	if err != nil {
		t.Fatal("Update log file failed.", err)
	}
	writeTestLog(LogLevelWarn)
	checkResult(t, LogLevelWarn)

	_ = os.Remove(defaultLogger.Path + defaultLogger.FileName)
	err = ForceUpdateLoggerFile()
	if err != nil {
		t.Fatal("Update log file failed.", err)
	}
	writeTestLog(LogLevelError)
	checkResult(t, LogLevelError)

	_ = os.Remove(defaultLogger.Path + defaultLogger.FileName)
	err = ForceUpdateLoggerFile()
	if err != nil {
		t.Fatal("Update log file failed.", err)
	}
	writeTestLog(LogLevelOff)
	checkResult(t, LogLevelOff)

	_ = os.Remove(defaultLogger.Path + defaultLogger.FileName)
}

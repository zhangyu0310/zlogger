package zlogger

import (
	"bufio"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMultiOpenAndWrite(t *testing.T) {
	l1, _ := realNew("./", "zlogger", false)
	l2, _ := realNew("./", "zlogger", false)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		for i := 0; i < 100; i++ {
			l1.Info("Info l1", i)
			l1.Debug("Debug l1", i)
			l1.Warn("Warn l1", i)
			l1.Error("Error l1", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 100; i++ {
			l2.Info("Info l2", i)
			l2.Debug("Debug l2", i)
			l2.Warn("Warn l2", i)
			l2.Error("Error l2", i)
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
}

func newForTest(t *testing.T) {
	err := New("./", "zlogger", true)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateLoggerFile(t *testing.T) {
	newForTest(t)
	var wg sync.WaitGroup
	wg.Add(5)
	begin := time.Now().UnixNano()
	go func() {
		for i := 0; i < 10000; i++ {
			Info("Info l1", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			Debug("Debug l1", i)
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
	for {
		line, _, err := reader.ReadLine()
		lineStr := string(line)
		if strings.Contains(lineStr, "DEBUG") {
			strings.Contains(lineStr, "Debug")
		} else if strings.Contains(lineStr, "INFO") {
			strings.Contains(lineStr, "Info")
		} else if strings.Contains(lineStr, "WARN") {
			strings.Contains(lineStr, "Warn")
		} else if strings.Contains(lineStr, "ERROR") {
			strings.Contains(lineStr, "Error")
		} else if lineStr == "" {

		} else {
			t.Error("String Contains error")
		}
		if err != nil {
			break
		}
	}
	t.Log("Check success.")
}

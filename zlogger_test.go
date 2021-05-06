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
	l1, _ := New("./", "zlogger", false)
	l2, _ := New("./", "zlogger", false)
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

func TestFatal(t *testing.T) {
	l1, _ := New("./", "zlogger", false)
	l1.Fatal("Fatal l1")
}

func TestDefaultFatal(t *testing.T) {
	Fatal("Fatal default")
}

func TestPanic(t *testing.T) {
	l1, _ := New("./", "zlogger", false)
	l1.Panic("Panic l1")
}

func TestDefaultPanic(t *testing.T) {
	Panic("Panic default")
}

func TestUpdateLoggerFile(t *testing.T) {
	l1, _ := New("./", "TestUpdateLoggerFile", false)
	var wg sync.WaitGroup
	wg.Add(5)
	begin := time.Now().UnixNano()
	go func() {
		for i := 0; i < 10000; i++ {
			l1.Info("Info l1", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			l1.Debug("Debug l1", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			l1.Warn("Warn l1", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			l1.Error("Error l1", i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			err := UpdateLoggerFile(l1)
			if err != nil {
				t.Error("UpdateLoggerFile failed.", err)
			}
		}
		wg.Done()
	}()
	wg.Wait()
	end := time.Now().UnixNano()
	t.Log("Time cost:", end-begin)
	//file, err := ioutil.ReadFile(l1.FileName)
	//if err != nil {
	//	t.Log("Check function ReadFile failed.", err)
	//}
	f, _ := os.Open(l1.Path + l1.FileName)
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

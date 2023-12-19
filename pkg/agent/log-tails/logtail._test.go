package logtail

import (
	"log"
	"os"
	"testing"
	"time"

	mock "github.com/samasno/monit/pkg/agent/mocks"
)

func TestOpenCloseLogtail(t *testing.T) {
	testFile := "./test.log"
	closer := writeToTestLog(testFile)
	time.Sleep(1 * time.Second)
	l := &LogTail{
		FilePath: testFile,
		Pipe:     &mock.MockEmitter{},
		Logger:   &mock.MockEmitter{},
	}
	err := l.Open()
	if err != nil {
		t.Fatal(err.Error())
	}
	err = l.Close()
	if err != nil {
		t.Fatal(err.Error())
	}
	closer()
}

func TestRunUpdate(t *testing.T) {
	testFile := "./test.log"
	closer := writeToTestLog(testFile)
	time.Sleep(1 * time.Second)
	l := &LogTail{
		FilePath: testFile,
		Pipe:     &mock.MockEmitter{},
		Logger:   &mock.MockEmitter{},
	}
	go func() {
		time.Sleep(60 * time.Second)
		l.Close()
	}()
	l.Update()
	closer()
}

func writeToTestLog(filePath string) func() {
	closer := make(chan bool)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0655)
	if err != nil {
		log.Fatal(err.Error())
	}
	testMsg := " test message test message"
	go func(closer chan bool) {
		breakLoop := false
		for {
			select {
			case _, ok := <-closer:
				if !ok {
					println("breaking test writer")
					breakLoop = true
				}
			default:
				time.Sleep(500 * time.Millisecond)
				time := time.Now().String()
				msg := time + testMsg + "\n"
				file.Write([]byte(msg))
			}
			if breakLoop {
				break
			}
		}
	}(closer)
	return func() {
		time.Sleep(2 * time.Second)
		close(closer)
		os.RemoveAll(filePath)
	}
}

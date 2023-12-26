package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/samasno/monit/pkg/agent/types"
	"github.com/samasno/monit/pkg/agent/vars"
)

type Logger struct {
	Listener      types.Listener
	LogFile       string
	logFileHandle *os.File
	shutdown      *sync.WaitGroup
	closer        chan bool
}

func (l *Logger) ListenAndLog() error {
	works := make(chan bool)
	go func(works chan bool) {
		defer l.close()
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT)
		err := l.open()
		if err != nil {
			l.log(vars.ERROR, err.Error())
			works <- false
			return
		}
		err = l.connect()
		if err != nil {
			l.log(vars.ERROR, err.Error())
			works <- false
		}
		fromListener := make(chan []byte, 1000)
		closeListener := make(chan bool)
		l.Listener.Listen(fromListener, closeListener, l.shutdown)
		go l.logWorker(fromListener)
		works <- true
		select {
		case <-l.closer:
			break
		case <-sig:
			break
		}
		close(sig)
		err = l.close()
		if err != nil {
			l.log(vars.ERROR, err.Error())
		}
	}(works)
	res := <-works
	if !res {
		return fmt.Errorf("Failed to start logger")
	}
	return nil
}

func (l *Logger) logWorker(job chan []byte) {
	closeWorker := false
	for {
		msg, ok := <-job
		if !ok {
			closeWorker = true
		}
		println(string(msg))
		formatted := formatEventToLog(msg)
		log.Println(formatted)
		if closeWorker {
			break
		}
	}
}

func (l *Logger) open() error {
	if l.closer == nil {
		l.closer = make(chan bool)
	}
	if l.logFileHandle != nil {
		l.log(vars.NOTICE, "Log file is already open")
		return nil
	}
	if l.LogFile == "" {
		return fmt.Errorf("No log file provided\n")
	}
	file, err := os.OpenFile(l.LogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open file at %s: %s\n", l.LogFile, err.Error())
	}
	l.logFileHandle = file
	log.SetOutput(file)
	return nil
}

func (l *Logger) connect() error {
	if l.Listener == nil {
		msg := "No listener attached to logger"
		l.log(vars.CRITICAL, msg)
		return fmt.Errorf(msg + "\n")
	}
	l.shutdown = &sync.WaitGroup{}
	err := l.Listener.Open(l.shutdown)
	if err != nil {
		msg := fmt.Sprintf("Failed to open listener: %s", err.Error())
		l.log(vars.CRITICAL, msg)
		return fmt.Errorf(msg + "\n")
	}
	return nil
}

func (l *Logger) close() error {
	if l.closer != nil {
		l.closer <- true
	}
	errs := []string{}
	log.SetOutput(os.Stdout)
	if l.logFileHandle != nil {
		err := l.logFileHandle.Close()
		if err != nil {
			println("close handle")
			errs = append(errs, err.Error())
		}
		l.logFileHandle = nil
	}
	err := l.Listener.Close()
	if err != nil {
		println("close listener")
		errs = append(errs, err.Error())
	}
	l.shutdown.Wait()
	if len(errs) > 0 {
		msg := strings.Join(errs, "\n")
		return fmt.Errorf(msg)
	}
	return nil
}

func formatEventToLog(eventB []byte) string {
	event := &types.Event{}
	err := json.Unmarshal(eventB, event)
	if err != nil {
		return "failed to parse log event"
	}
	bldr := strings.Builder{}
	bldr.WriteString(strconv.Itoa(event.Payload.Level) + ":")
	bldr.WriteString(" " + event.Type)
	bldr.WriteString(" " + event.Payload.Source)
	bldr.WriteString(" " + event.Payload.Message)
	return bldr.String()
}

func (l *Logger) log(level int, msg string) error {
	if l.logFileHandle == nil {
		err := l.open()
		if err != nil {
			return err
		}
	}
	msg = fmt.Sprintf("%d: %s %s", level, loggerName, msg)
	log.Println(msg)
	return nil
}

var loggerName = "monit-logger"

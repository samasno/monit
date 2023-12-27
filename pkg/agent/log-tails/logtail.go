package logtail

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/samasno/monit/pkg/agent/types"
	"github.com/samasno/monit/pkg/agent/vars"
)

type LogTail struct {
	FilePath   string
	Pipe       types.Emitter
	Logger     types.Emitter
	fileHandle *os.File
	position   int64
	mtx        *sync.Mutex
	wg         *sync.WaitGroup
}

func (l *LogTail) Open() error {
	if l.mtx == nil {
		l.mtx = &sync.Mutex{}
	}
	if l.fileHandle != nil {
		l.log(vars.NOTICE, "File already open "+l.FilePath)
		return nil
	}
	handle, err := os.OpenFile(l.FilePath, os.O_RDONLY|os.O_CREATE, 0444)
	if err != nil {
		msg := fmt.Sprintf("Failed to open file %s: %s", l.FilePath, err.Error())
		l.log(vars.CRITICAL, msg)
		return fmt.Errorf(msg + "\n")
	}
	l.log(vars.INFO, "Succesfully opened target file")
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.fileHandle = handle
	info, err := l.fileHandle.Stat()
	if err != nil {
		msg := fmt.Sprintf("Failed to stat file: %s", err.Error())
		l.log(vars.ERROR, msg)
		return fmt.Errorf(msg + "\n")
	}
	l.position = info.Size()
	msg := fmt.Sprintf("Set position to %d", l.position)
	l.log(vars.INFO, msg)
	return nil
}

func (l *LogTail) Close() error {
	l.mtx.Lock()
	if l.fileHandle == nil {
		l.log(vars.NOTICE, "File is already closed")
		return nil
	}
	defer l.mtx.Unlock()
	l.wg.Wait()
	err := l.fileHandle.Close()
	l.fileHandle = nil
	if err != nil {
		msg := fmt.Sprintf("Failed to close file: %s", err.Error())
		l.log(vars.ERROR, msg)
		return fmt.Errorf(msg)
	}
	return nil
}

func (l *LogTail) Update() error {
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGHUP)
	if l.fileHandle == nil {
		err := l.Open()
		if err != nil {
			return err
		}
	}
	sendChan := make(chan string, 100)
	for {
		select {
		case <-shutdown:
			l.log(vars.INFO, "Received signal to terminate")
			err := l.Close()
			if err != nil {
				msg := "Failed to close properly after receiving signal to terminate."
				l.log(vars.ERROR, msg)
				return err
			}
			close(sendChan)
			l.wg.Wait()
			return nil
		default:
			if l.fileHandle == nil {
				msg := "Found closed file while updating"
				l.log(vars.NOTICE, msg)
				return l.Close()
			}
			l.mtx.Lock()
			l.wg = &sync.WaitGroup{}
			l.mtx.Unlock()
			info, err := l.fileHandle.Stat()
			if err != nil {
				msg := fmt.Sprintf("Error while running update: %s", err.Error())
				l.log(vars.ERROR, msg)
				l.log(vars.NOTICE, "Skipping update due to error")
				time.Sleep(15 * time.Second)
				continue
			}
			size := info.Size()
			if l.position > size {
				l.position = 0
				l.log(vars.INFO, "File size is smaller than position, assuming log rotation and setting position back to 0")
			}
			var readN int64
			var end int64
			var limit int64 = 5000000
			diff := size - l.position
			if diff > limit {
				readN = limit
				end = l.position + limit
			} else {
				readN = size - l.position
				end = size
			}
			l.log(vars.INFO, fmt.Sprintf("Current position set at %d", l.position))
			l.log(vars.INFO, fmt.Sprintf("Current file size %d", size))
			l.log(vars.INFO, fmt.Sprintf("Reading from %d to %d", l.position, end))
			section := io.NewSectionReader(l.fileHandle, l.position, readN)
			scanner := bufio.NewScanner(section)
			go l.updateWorker(sendChan)
			var totalBytes int64
			for scanner.Scan() {
				line := scanner.Text()
				sendChan <- line
				totalBytes += int64(len(line))
			}
			l.wg.Wait()
			err = scanner.Err()
			if err != nil && err != io.EOF {
				l.log(vars.ERROR, fmt.Sprintf("Encountered non-EOF error while scanning file"))
				l.Close()
			}
			l.position = end
			l.log(vars.INFO, fmt.Sprintf("Read %d bytes", totalBytes))
			time.Sleep(15 * time.Second)
		}
	}
}

func (l *LogTail) updateWorker(in chan string) {
	defer func() {
		r := recover()
		if r != nil {
			l.log(vars.ALERT, "Recovered from panic in update worker")
			l.log(vars.ALERT, fmt.Sprintf("%v", r))
		}
	}()
	for {
		str, ok := <-in
		if !ok {
			l.log(vars.NOTICE, "Update worker's channel closed, shutting down worker")
			break
		}
		e := types.Event{
			Type: vars.LOGTAIL_UPDATE,
			Payload: types.Payload{
				Source:  l.FilePath,
				Message: str,
				Level:   vars.INFO,
			},
		}
		err := l.Pipe.Emit(e)
		if err != nil {
			l.log(vars.ERROR, "Failed to forward log tail update: "+err.Error())
		}
	}
	l.wg.Done()
}

func (l *LogTail) log(level int, message string) error {
	if l.Logger == nil {
		return fmt.Errorf(("Logtail has no logger to send messages"))
	}
	payload := types.Payload{
		Source:  NAME,
		Message: message,
		Level:   level,
	}
	event := types.Event{
		Type:    vars.LOGTAIL_LOG,
		Payload: payload,
	}
	err := l.Logger.Emit(event)
	if err != nil {
		msg := fmt.Sprintf("Failed to emit message")
		return fmt.Errorf(msg + "\n")
	}
	return nil
}

var (
	NAME = "logtail"
)

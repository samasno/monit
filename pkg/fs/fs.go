package fs

import (
	"log"
	"os"
	"path"
)

var (
	APP_DIR          = "monit-go"
	FORWARDER_SOCKET = "fwd.sock"
	LOGGER_SOCKET    = "log.sock"
	LOGGER_FILE      = "monit-go.log"
)

// starts forwarder and logger at the same time

func SetupWorkDir() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err.Error())
	}
	workingDirectory := path.Join(homedir, APP_DIR)
	_, err = os.Stat(workingDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(workingDirectory, 0700)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	forwarderSocket := path.Join(workingDirectory, FORWARDER_SOCKET)
	loggerSocket := path.Join(workingDirectory, LOGGER_SOCKET)
	loggerFile := path.Join(workingDirectory, LOGGER_FILE)
	stat, _ := os.Stat(forwarderSocket)
	if stat != nil {
		os.RemoveAll(forwarderSocket)
	}
	stat, _ = os.Stat(loggerSocket)
	if stat != nil {
		os.RemoveAll(loggerSocket)
	}
	_, err = os.Stat(loggerFile)
	if os.IsNotExist(err) {
		file, err := os.OpenFile(loggerFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0700)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func ForwarderSocket() string {
	homedir, _ := os.UserHomeDir()
	return path.Join(homedir, APP_DIR, FORWARDER_SOCKET)
}

func LoggerSocket() string {
	homedir, _ := os.UserHomeDir()
	return path.Join(homedir, APP_DIR, LOGGER_SOCKET)
}

func LogFile() string {
	homedir, _ := os.UserHomeDir()
	return path.Join(homedir, APP_DIR, LOGGER_FILE)
}

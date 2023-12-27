package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/samasno/monit/pkg/agent/emitters"
	logtail "github.com/samasno/monit/pkg/agent/log-tails"
	"github.com/samasno/monit/pkg/fs"
)

var targetFile string

func main() {
	flag.StringVar(&targetFile, "file", "", "path to target file to tail and forward")
	flag.Parse()
	println(targetFile)
	fwdsock := fs.ForwarderSocket()
	fwdPipe := &emitters.SocketEmitter{
		Raddr: fwdsock,
	}
	loggersock := fs.LoggerSocket()
	logPipe := &emitters.SocketEmitter{
		Raddr: loggersock,
	}
	tail := logtail.LogTail{
		FilePath: targetFile,
		Pipe:     fwdPipe,
		Logger:   logPipe,
	}
	go tail.Update()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP)
	<-sig
	err := tail.Close()
	if err != nil {
		println(err.Error())
		log.Fatal(err.Error())
	}
}

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/samasno/monit/pkg/agent/emitters"
	"github.com/samasno/monit/pkg/agent/forwarder"
	"github.com/samasno/monit/pkg/agent/listeners"
	"github.com/samasno/monit/pkg/agent/logger"
	"github.com/samasno/monit/pkg/agent/types"
	fs "github.com/samasno/monit/pkg/fs"
)

var upstreamUrl string
var upstreamPort int

func main() {
	flag.StringVar(&upstreamUrl, "url", "", "URL to upstream host")
	flag.IntVar(&upstreamPort, "port", 0, "Upstream host port number to connect")
	flag.Parse()
	if upstreamUrl == "" {
		log.Fatal("--url argument required")
	}
	if upstreamPort == 0 {
		log.Fatal("--port argument required")
	}
	err := fs.SetupWorkDir()
	if err != nil {
		log.Fatal(err.Error())
	}
	fwdsock := fs.ForwarderSocket()
	logsock := fs.LoggerSocket()
	logEmitter := &emitters.SocketEmitter{
		Raddr: logsock,
	}
	upstream := &types.Upstream{
		Url:  upstreamUrl,
		Port: upstreamPort,
	}
	tcpClient := &forwarder.ForwarderTcpClient{
		Upstream: upstream,
		Logger:   logEmitter,
	}
	lndn := &types.Downstream{
		Url: fwdsock,
	}
	ln := &listeners.UnixDatagramSocketListener{
		Name:       "forwarder-downstream-client",
		Downstream: lndn,
		Logger:     logEmitter,
	}
	fwd := forwarder.Forwarder{
		UpstreamClient:     tcpClient,
		DownstreamListener: ln,
	}
	logFile := fs.LogFile()
	logSocket := fs.LoggerSocket()
	loggerLn := &listeners.UnixDatagramSocketListener{
		Name: "logger-datagram-listener",
		Downstream: &types.Downstream{
			Url: logSocket,
		},
	}
	logger := logger.Logger{
		Listener: loggerLn,
		LogFile:  logFile,
	}
	err = logger.ListenAndLog()
	if err != nil {
		fwd.Close()
		log.Fatal(err.Error())
	}
	run := make(chan bool)
	go func(r chan bool) {
		err := fwd.Run()
		if err != nil {
			println(err.Error())
			run <- false
		} else {
			run <- true
		}
	}(run)
	fwdOk := <-run
	if !fwdOk {
		log.Fatal("Failed to start forwarder")
	}
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT)
	<-sig
	fwd.Close()
	println("Got signal to terminate")
	if err != nil {
		log.Fatal(err.Error())
	}
}

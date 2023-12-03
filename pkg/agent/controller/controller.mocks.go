package controller

import (
	"fmt"
	"io"

	"github.com/samasno/monit/pkg/agent/types"
)

type MockForwarder struct {
	Name   string
	URL    string
	Logger types.Logger
}

func (m *MockForwarder) Connect() error {
	m.Logger.StdOut(fmt.Sprintf("Forwarder %s connected to %s", m.Name, m.URL))
	return nil
}

func (m *MockForwarder) Close() error {
	m.Logger.StdOut(fmt.Sprintf("Forwarder %s closed connection to %s", m.Name, m.URL))
	return nil
}

func (m *MockForwarder) Push([]byte) error {
	m.Logger.StdOut(fmt.Sprintf("Forwarder %s pushed msg to connection %s", m.Name, m.URL))
	return nil
}

func (m *MockForwarder) Status() (types.Status, error) {
	return MockForwarderStatus{}, nil
}

type MockForwarderStatus struct{}

func (m MockForwarderStatus) Message() string {
	return "Mock forwarder status OK"
}

func (m MockForwarderStatus) Ok() bool {
	return true
}

type MockLogTail struct {
	Name        string
	LogLocation string
	Logger      types.Logger
}

func (m *MockLogTail) Open() error {
	m.Logger.StdOut(fmt.Sprintf("LogTail %s opened log file at %s", m.Name, m.LogLocation))
	return nil
}

func (m *MockLogTail) Close() error {
	m.Logger.StdOut(fmt.Sprintf("LogTail %s closed log file at %s", m.Name, m.LogLocation))
	return nil
}

func (m *MockLogTail) Update() error {
	m.Logger.StdOut(fmt.Sprintf("LogTail %s update sent for log file at %s", m.Name, m.LogLocation))
	return nil
}

func (m *MockLogTail) Status() (types.Status, error) {
	return MockLogTailStatus{}, nil
}

type MockLogTailStatus struct{}

func (m MockLogTailStatus) Message() string {
	return "Mock log tail status OK"
}

func (m MockLogTailStatus) Ok() bool {
	return true
}

type MockLogger struct {
	Name        string
	LogLocation string
	Output      io.ReadWriter
	Error       io.ReadWriter
}

func (m *MockLogger) StdOut(msg string) error {
	m.Output.Write([]byte(msg))
	return nil
}

func (m *MockLogger) StdErr(msg string) error {
	m.Error.Write([]byte(msg))
	return nil
}

func (m *MockLogger) Close() error {
	msg := fmt.Sprintf("Logger %s closed log file and shutting down", m.Name)
	m.Output.Write([]byte(msg))
	return nil
}

func (m *MockLogger) Status() (types.Status, error) {
	return MockLoggerStatus{}, nil
}

type MockLoggerStatus struct{}

func (m MockLoggerStatus) Message() string {
	return "Mock logger status OK"
}

func (m MockLoggerStatus) Ok() bool {
	return true
}

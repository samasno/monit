package mock

import "github.com/samasno/monit/pkg/agent/types"

type MockEmitter struct {
}

func (m *MockEmitter) Emit(e types.Event) error {
	println(e.Payload.Message)
	return nil
}

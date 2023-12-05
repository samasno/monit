package forwarder

type ForwarderStatus struct {
	IsOk        bool
	MessageText string
}

func (fs ForwarderStatus) Message() string {
	return fs.MessageText
}

func (fs ForwarderStatus) Ok() bool {
	return fs.IsOk
}

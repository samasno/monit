package i

type Controller interface { // manages the forwarder to upstream and log runners
	Init()
	Run()
	Shutdown()
	Status()
}

type Forwarder interface {
	Connect()
	Push()
	Close()
	Status()
}

type LogTail interface {
	Update()
	Open()
	Close()
	Status()
}

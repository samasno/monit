package controller

import (
	i "github.com/samasno/monit/pkg/agent/interfaces"
	t "github.com/samasno/monit/pkg/agent/types"
)

type Controller struct {
	Forwarder i.Forwarder
	LogTails  []i.LogTail
	Logger    i.Logger
}

func NewController() i.Controller {
	c := &Controller{}
	return c
}

func (c *Controller) Init(t.ControllerInitInput) error {
	return nil
}
func (c *Controller) Run() error {
	return nil
}

func (c *Controller) Shutdown() error {
	return nil
}

func (c Controller) Status() t.ControllerStatus {
	s := t.ControllerStatus{}
	return s
}

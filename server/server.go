package server

import (
	"os"
	"syscall"
	"time"
)

type IServer interface {
	Run() error
}

type Option struct {
	Post    int
	Timeout time.Duration
	Signals []os.Signal
	OnStart func()
	OnStop  func()
}

func NewDefault() (*Option, error) {
	return &Option{
		Post:    2345,
		Signals: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		Timeout: time.Second * 60, // 60 s
	}, nil
}

func NewHttpDefault() (*Option, error) {
	o, err := NewDefault()
	o.Post = 3939
	return o, err
}

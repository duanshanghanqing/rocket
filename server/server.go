package server

import (
	"github.com/duanshanghanqing/rocket/registry"
	"github.com/google/uuid"
	"os"
	"syscall"
	"time"
)

type IServer interface {
	Run() error
}

type Option struct {
	ID                    string
	Post                  int
	Timeout               time.Duration
	Signals               []os.Signal
	ServiceRegisterCenter registry.IRegistrar
	ServiceRegisterInfo   *registry.ServiceRegisterInfo
}

func NewDefault() (*Option, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &Option{
		ID:      uid.String(),
		Post:    2345,
		Signals: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		Timeout: time.Second * 30, // 30 s
	}, nil
}

func NewHttpDefault() (*Option, error) {
	o, err := NewDefault()
	o.Post = 3939
	return o, err
}

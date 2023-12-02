package server

import (
	"github.com/google/uuid"
	"os"
	"syscall"
	"time"
	"github.com/duanshanghanqing/rocket/registry"
)

type IServer interface {
	Run() error
}

type Option struct {
	ID                    string
	Name                  string
	Post                  int
	Timeout               time.Duration
	Signals               []os.Signal
	ServiceRegisterCenter registry.IRegistrar
	ServiceRegisterInfo   *registry.ServiceRegisterInfo
	ServiceRegisterHost   string
}

func NewDefault() (*Option, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &Option{
		ID:                  uid.String(),
		Post:                2345,
		Signals:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		Timeout:             time.Second * 30, // 30 s
		ServiceRegisterHost: "127.0.0.1",
	}, nil
}

func NewHttpDefault() (*Option, error) {
	o, err := NewDefault()
	o.Post = 3939
	return o, err
}

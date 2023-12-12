package httpserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/duanshanghanqing/rocket/pkg/utils"
	"github.com/duanshanghanqing/rocket/registry"
	"github.com/duanshanghanqing/rocket/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	option     *server.Option
	httpServer *http.Server
}

type ServerOption func(s *Server)

func WithServerOptionID(id string) ServerOption {
	return func(s *Server) {
		s.option.ID = id
	}
}

func WithServerOptionName(name string) ServerOption {
	return func(s *Server) {
		s.option.Name = name
	}
}

func WithServerOptionPost(post int) ServerOption {
	return func(s *Server) {
		s.option.Post = post
	}
}

func WithServerOptionTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.option.Timeout = timeout
	}
}

func WithServerOptionSignal(signals []os.Signal) ServerOption {
	return func(s *Server) {
		s.option.Signals = signals
	}
}

func WithServerOptionServiceRegisterCenter(serviceRegisterCenter registry.IRegistrar) ServerOption {
	return func(s *Server) {
		s.option.ServiceRegisterCenter = serviceRegisterCenter
	}
}

func WithServerOptionServiceRegisterInfo(serviceRegisterInfo *registry.ServiceRegisterInfo) ServerOption {
	return func(s *Server) {
		s.option.ServiceRegisterInfo = serviceRegisterInfo
	}
}

func WithServerHttpServer(httpServer *http.Server) ServerOption {
	return func(s *Server) {
		s.httpServer = httpServer
	}
}

func (s *Server) startHttpServer() error {
	// 服务注册
	if s.option.ServiceRegisterCenter != nil {
		err := s.option.ServiceRegisterCenter.Register(context.Background(), s.option.ServiceRegisterInfo)
		if err != nil {
			return err
		}
	}

	// When the user does not implement a handler, implement a default
	if s.httpServer.Handler == nil {
		mux := http.NewServeMux()
		// health examination
		mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(string("ok")))
		})
		s.httpServer.Handler = mux
	}

	log.Printf("http server start: %s", utils.HostPostToAddress("", s.option.Post))

	err := s.httpServer.ListenAndServe()
	time.Sleep(6 * time.Second)
	return err
}

func (s *Server) stopHttpServer() error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, s.option.Signals...)
	err := fmt.Errorf("%s", <-signalChan)
	log.Println("http server stopping")

	if s.option.ServiceRegisterCenter != nil {
		_ = s.option.ServiceRegisterCenter.Deregister(context.Background(), s.option.ServiceRegisterInfo)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.httpServer.Shutdown(ctx)
	time.Sleep(5 * time.Second)
	log.Println("http server stop")
	return err
}

func (s *Server) Run() error {
	errChan := make(chan error)

	go func() {
		errChan <- s.startHttpServer()
	}()

	go func() {
		errChan <- s.stopHttpServer()
	}()

	return <-errChan
}

func New(opts ...ServerOption) (server.IServer, error) {

	defaultOption, err := server.NewHttpDefault()
	if err != nil {
		return nil, err
	}

	s := &Server{
		option: defaultOption,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.option.Name == "" {
		return nil, errors.New("service name cannot be empty")
	}

	if s.httpServer == nil {
		s.httpServer = &http.Server{
			Addr: utils.HostPostToAddress("", s.option.Post),
		}
	}

	_, post, _ := utils.AddressToHostPost(s.httpServer.Addr)
	s.option.Post = post
	if s.option.ServiceRegisterInfo != nil {
		if s.option.ServiceRegisterInfo.Host == "" {
			return nil, errors.New("service register host cannot be empty")
		}
		s.option.ServiceRegisterInfo.ID = s.option.ID
		s.option.ServiceRegisterInfo.Name = s.option.Name
		s.option.ServiceRegisterInfo.Port = post
	}

	return s, nil
}

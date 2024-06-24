package httpserver

import (
	"context"
	"fmt"
	"github.com/duanshanghanqing/rocket/pkg/utils"
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
		s.option.Signals = append(s.option.Signals, signals...)
	}
}

func WithServerHttpServer(httpServer *http.Server) ServerOption {
	return func(s *Server) {
		s.httpServer = httpServer
	}
}

func WithServerOptionOnStart(onStart func()) ServerOption {
	return func(s *Server) {
		s.option.OnStart = onStart
	}
}

func WithServerOptionOnStop(onStop func()) ServerOption {
	return func(s *Server) {
		s.option.OnStop = onStop
	}
}

func (s *Server) startHttpServer() error {
	if s.option.OnStart != nil {
		s.option.OnStart()
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

	return s.httpServer.ListenAndServe()
}

func (s *Server) stopHttpServer() error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, s.option.Signals...)
	err := fmt.Errorf("%s", <-signalChan)
	log.Println("http server stopping")

	if s.option.OnStop != nil {
		s.option.OnStop()
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

	if s.httpServer == nil {
		s.httpServer = &http.Server{
			Addr: utils.HostPostToAddress("", s.option.Post),
		}
	}

	_, post, _ := utils.AddressToHostPost(s.httpServer.Addr)
	s.option.Post = post

	return s, nil
}

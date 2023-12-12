package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/duanshanghanqing/rocket/pkg/utils"
	"github.com/duanshanghanqing/rocket/registry"
	"github.com/duanshanghanqing/rocket/server"
	"github.com/duanshanghanqing/rocket/server/grpcserver/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	option                   *server.Option
	unaryServerInterceptors  []grpc.UnaryServerInterceptor
	streamServerInterceptors []grpc.StreamServerInterceptor
	serverOptions            []grpc.ServerOption
	healthServer             grpc_health_v1.HealthServer
	registerServer           func(server *grpc.Server)
	grpcServer               *grpc.Server
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

func WithServerOptionServiceRegisterHost(host string) ServerOption {
	return func(s *Server) {
		s.option.ServiceRegisterHost = host
	}
}

func WithServerOptionPost(post int) ServerOption {
	return func(s *Server) {
		s.option.Post = post
	}
}

func WithServerUnaryServerInterceptors(unaryServerInterceptors []grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryServerInterceptors = unaryServerInterceptors
	}
}

func WithServerStreamServerInterceptors(streamServerInterceptors []grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamServerInterceptors = streamServerInterceptors
	}
}

func WithServerServerOptions(serverOptions []grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.serverOptions = serverOptions
	}
}

func WithServerHealthServer(healthServer grpc_health_v1.HealthServer) ServerOption {
	return func(s *Server) {
		s.healthServer = healthServer
	}
}

func WithServerRegisterServer(registerServer func(server *grpc.Server)) ServerOption {
	return func(s *Server) {
		s.registerServer = registerServer
	}
}

func New(opts ...ServerOption) (server.IServer, error) {
	// Default options
	defaultOption, err := server.NewDefault()
	if err != nil {
		return nil, err
	}

	s := &Server{
		option:                   defaultOption,
		registerServer:           func(server *grpc.Server) {},
		healthServer:             newHealthServer(),
		unaryServerInterceptors:  make([]grpc.UnaryServerInterceptor, 0),
		streamServerInterceptors: make([]grpc.StreamServerInterceptor, 0),
		serverOptions:            make([]grpc.ServerOption, 0),
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.option.Name == "" {
		return nil, errors.New("service name cannot be empty")
	}

	// Set service registration information
	s.option.ServiceRegisterInfo = &registry.ServiceRegisterInfo{
		ID:   s.option.ID,
		Name: s.option.Name,
		Host: s.option.ServiceRegisterHost,
		Port: s.option.Post,
	}

	return s, nil
}

func (s *Server) startGrpcServer() error {
	// We now hope that if users do not set interceptors, we will automatically add some necessary interceptors by default, such as crash tracing
	unaryIns := []grpc.UnaryServerInterceptor{
		interceptors.UnaryTimeoutInterceptor(s.option.Timeout),
	}

	// Convert the interceptor we passed in to the serverOption of grpc
	var grpcOpts = []grpc.ServerOption{grpc.ChainUnaryInterceptor(unaryIns...)}

	// Interceptors passed in by the user themselves
	if len(s.unaryServerInterceptors) > 0 {
		grpcOpts = append(grpcOpts, grpc.ChainUnaryInterceptor(s.unaryServerInterceptors...))
	}

	// stream server interceptor
	if len(s.streamServerInterceptors) > 0 {
		grpcOpts = append(grpcOpts, grpc.ChainStreamInterceptor(s.streamServerInterceptors...))
	}

	// Transfer the GRPC passed in by the user themselves Putting ServerOptions Together
	if len(s.serverOptions) > 0 {
		grpcOpts = append(grpcOpts, s.serverOptions...)
	}

	// Set up listening
	lis, err := net.Listen("tcp", utils.HostPostToAddress("", s.option.Post))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	// Create Service
	s.grpcServer = grpc.NewServer(grpcOpts...)

	// Register Health Server
	grpc_health_v1.RegisterHealthServer(s.grpcServer, s.healthServer)

	// Registering services developed by users themselves,
	s.registerServer(s.grpcServer)

	// Service Registration Center, Registering GRPC Services
	if s.option.ServiceRegisterCenter != nil { // Explain that the user has implemented the service registration center themselves
		err = s.option.ServiceRegisterCenter.Register(context.Background(), s.option.ServiceRegisterInfo)
		if err != nil {
			return err
		}
	}

	log.Printf("grpc server start: %s", lis.Addr())

	return s.grpcServer.Serve(lis) // Will block
}

func (s *Server) Run() error {
	// Create save error synchronization pipeline
	errChan := make(chan error)

	// Start the GRPC service and register with the registry
	go func() {
		errChan <- s.startGrpcServer()
	}()

	// Listen and exit
	go func() {
		// Listening for exit signals
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, s.option.Signals...) // This place will not block
		err := fmt.Errorf("%s", <-signalChan)          // Block here, only execute downwards when receiving the above two signals, and retrieve the pipeline value
		// 服务优雅退出
		if s.grpcServer != nil {
			log.Println("grpc server stopping")
			// Unregister service
			if s.option.ServiceRegisterCenter != nil {
				_ = s.option.ServiceRegisterCenter.Deregister(context.Background(), s.option.ServiceRegisterInfo)
			}
			// Elegant Exit
			s.grpcServer.GracefulStop()
			log.Println("grpc server stop")
		}
		errChan <- err
	}()

	return <-errChan
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type healthServer struct{}

func (s *healthServer) Watch(request *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return nil
}

func (s *healthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	// Implement your health check logic
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func newHealthServer() grpc_health_v1.HealthServer {
	return &healthServer{}
}

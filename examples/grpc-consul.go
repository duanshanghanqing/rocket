package main

import (
	"fmt"
	"github.com/duanshanghanqing/rocket/server/grpcserver"
	consulApi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"log"
)

type GrpcConsulRegisterCenter struct {
	ID     string
	client *consulApi.Client
}

func (r *GrpcConsulRegisterCenter) Register() {
	// 1.Registration Service Information
	reg := consulApi.AgentServiceRegistration{}
	reg.ID = r.ID
	reg.Name = "grpc-consul"
	reg.Address = "127.0.0.1"
	reg.Port = 8090
	reg.Check = &consulApi.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", reg.Address, reg.Port),
		Interval:                       "10s",
		Timeout:                        "5s",
		DeregisterCriticalServiceAfter: "30s",
	}
	reg.Tags = []string{}

	// 2.Register Service
	serviceRegisterOpts := consulApi.ServiceRegisterOpts{
		//Token: "",
	}
	if err := r.client.Agent().ServiceRegisterOpts(&reg, serviceRegisterOpts); err != nil {
		log.Printf("err: %s", err.Error())
		return
	}

	//if err := r.client.Agent().ServiceRegister(&reg); err != nil {
	//	log.Printf("err: %s", err.Error())
	//	return
	//}
}

func (r *GrpcConsulRegisterCenter) Deregister() {
	if err := r.client.Agent().ServiceDeregister(r.ID); err != nil {
		log.Printf("err: %s", err.Error())
		return
	}
}

func NewGrpcConsulRegisterCenter(config *consulApi.Config) (*GrpcConsulRegisterCenter, error) {
	client, err := consulApi.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &GrpcConsulRegisterCenter{
		ID:     "grpc-consul",
		client: client,
	}, nil
}

func main() {
	// Initialize Service Registry
	consulRegisterCenter, err := NewGrpcConsulRegisterCenter(
		&consulApi.Config{
			Address: "127.0.0.1:8500",
		},
	)
	if err != nil {
		log.Printf("err: %s", err.Error())
		return
	}

	server, err := grpcserver.New(
		grpcserver.WithServerOptionPost(8090),
		grpcserver.WithServerRegisterServer(func(server *grpc.Server) {
			// register your service
			// userpb.RegisterUserServer(server, user.NewUserServer())
		}),
		grpcserver.WithServerOptionServiceRegisterCenter(consulRegisterCenter),
	)

	if err != nil {
		log.Printf("err: %s", err.Error())
		return
	}

	if err = server.Run(); err != nil {
		log.Printf("err: %s", err.Error())
	}
}

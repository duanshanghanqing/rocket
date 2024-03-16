package main

import (
	//"github.com/duanshanghanqing/rocket/pkg/utils"
	"github.com/duanshanghanqing/rocket/registry"
	"github.com/duanshanghanqing/rocket/server/grpcserver"
	"github.com/duanshanghanqing/rocket/server/grpcserver/registrycenter"
	consulApi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"log"
)

func main() {
	// Initialize Service Registry
	consulRegisterCenter, err := registrycenter.NewConsulRegisterCenter(
		registrycenter.WithConsulRegisterCenterClientConfig(&consulApi.Config{
			Address: "127.0.0.1:8500",
		}),
	)
	if err != nil {
		log.Printf("err: %s", err.Error())
		return
	}

	//externalIp, err := utils.GetExternalIp()
	//if err != nil {
	//	log.Printf("err: %s", err.Error())
	//	return
	//}

	server, err := grpcserver.New(
		grpcserver.WithServerOptionPost(8090),
		grpcserver.WithServerRegisterServer(func(server *grpc.Server) {
			// register your service
			// userpb.RegisterUserServer(server, user.NewUserServer())
		}),
		// Set up registration center
		grpcserver.WithServerOptionServiceRegisterInfo(&registry.ServiceRegisterInfo{
			Name: "grpc-server",
			Host: "127.0.0.1", // local
			//Host: externalIp, // external ip of cloud server
			Tags: []string{"dev"},
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

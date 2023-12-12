package main

import (
	"github.com/duanshanghanqing/rocket/server/grpcserver"
	"google.golang.org/grpc"
	"log"
)

func main() {
	server, err := grpcserver.New(
		grpcserver.WithServerOptionName("grpc-server"),
		grpcserver.WithServerOptionPost(8090),
		grpcserver.WithServerRegisterServer(func(server *grpc.Server) {
			// register your service
			// userpb.RegisterUserServer(server, user.NewUserServer())
		}),
	)

	if err != nil {
		log.Printf("err: %s", err.Error())
		return
	}

	if err = server.Run(); err != nil {
		log.Printf("err: %s", err.Error())
	}
}

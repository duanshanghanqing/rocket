package main

import (
	"fmt"
	consulApi "github.com/hashicorp/consul/api"

	"github.com/duanshanghanqing/rocket/server/httpserver"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ConsulRegisterCenter struct {
	ID     string
	client *consulApi.Client
}

func (r *ConsulRegisterCenter) Register() {
	// 1.Registration Service Information
	reg := consulApi.AgentServiceRegistration{}
	reg.ID = r.ID
	reg.Name = "http-consul"
	reg.Address = "127.0.0.1"
	reg.Port = 8091
	reg.Check = &consulApi.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", reg.Address, reg.Port),
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

func (r *ConsulRegisterCenter) Deregister() {
	if err := r.client.Agent().ServiceDeregister(r.ID); err != nil {
		log.Printf("err: %s", err.Error())
		return
	}
}

func NewConsulRegisterCenter(config *consulApi.Config) (*ConsulRegisterCenter, error) {
	//clientConfig := consulApi.DefaultConfig()
	client, err := consulApi.NewClient(config)
	if err != nil {
		return nil, err
	}
	consulRegisterCenter := &ConsulRegisterCenter{
		ID:     "1",
		client: client,
	}
	return consulRegisterCenter, nil
}

func main() {
	// consulRegisterCenter
	consulRegisterCenter, err := NewConsulRegisterCenter(&consulApi.Config{
		Address: "127.0.0.1:8500",
	})
	if err != nil {
		log.Printf("err: %s", err.Error())
		return
	}

	// http
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	// server
	server, err := httpserver.New(
		httpserver.WithServerHttpServer(
			&http.Server{
				Addr:    fmt.Sprintf(":%d", 8091),
				Handler: r,
			},
		),
		httpserver.WithServerOptionServiceRegisterCenter(consulRegisterCenter),
	)
	if err != nil {
		log.Printf("err: %s", err.Error())
		return
	}

	if err = server.Run(); err != nil {
		log.Printf("err: %s", err.Error())
	}
}

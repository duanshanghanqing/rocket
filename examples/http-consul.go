package main

import (
	"fmt"
	//"github.com/duanshanghanqing/rocket/pkg/utils"
	"github.com/duanshanghanqing/rocket/registry"
	"github.com/duanshanghanqing/rocket/server/httpserver"
	"github.com/duanshanghanqing/rocket/server/httpserver/registrycenter"
	"github.com/gin-gonic/gin"
	consulApi "github.com/hashicorp/consul/api"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

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

	server, err := httpserver.New(
		httpserver.WithServerHttpServer(
			&http.Server{
				Addr:    fmt.Sprintf(":%d", 8091),
				Handler: r,
			},
		),
		// Set up registration center
		httpserver.WithServerOptionServiceRegisterInfo(&registry.ServiceRegisterInfo{
			Name: "http-server",
			Host: "127.0.0.1", // local
			//Host: externalIp, // external ip of cloud server
			Tags: []string{"http-server", "8091"},
		}),
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

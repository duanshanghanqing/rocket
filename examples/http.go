package main

import (
	"fmt"
	"github.com/duanshanghanqing/rocket/server/httpserver"
	"github.com/gin-gonic/gin"
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

	server, err := httpserver.New(
		httpserver.WithServerOptionName("http-server"),
		httpserver.WithServerHttpServer(
			&http.Server{
				Addr:    fmt.Sprintf(":%d", 8091),
				Handler: r,
			},
		),
	)
	if err != nil {
		log.Printf("err: %s", err.Error())
		return
	}

	if err = server.Run(); err != nil {
		log.Printf("err: %s", err.Error())
	}
}

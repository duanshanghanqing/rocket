# a quickly build GRPC or HTTP server framework

**main features:**
- grpc
- http
- service register

**registration center:**
- consul
- etcd (todo)
- nacos (todo)

[examples](https://github.com/duanshanghanqing/rocket/examples)

## Getting rocket

### Prerequisites

- **[Go](https://go.dev/)**: any one of the **three latest major** [releases](https://go.dev/doc/devel/release) .
- Go version >= 1.18

### Getting rocket

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```
import "github.com/duanshanghanqing/rocket"
```

Otherwise, run the following Go command to install the `rocket` package:

```sh
$ go get -u github.com/duanshanghanqing/rocket
```

### grpc server

```go
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

```

### http server

```go
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
			"message": "pong",
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

```
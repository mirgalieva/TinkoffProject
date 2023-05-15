package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
)

func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panicMiddleware", c.Request.URL.Path)
			fmt.Println("recovered", err)
			c.JSON(http.StatusInternalServerError, err)
		}
	}()
	c.Next()
}
func PanicUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("gRPC panic Middleware", info.FullMethod, '\n', "recovered", err)
		}
	}()
	resp, err := handler(ctx, req)
	return resp, err
}

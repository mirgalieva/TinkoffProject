package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"time"
)

func Logger(c *gin.Context) {
	t := time.Now()
	c.Next()
	latency := time.Since(t)
	status := c.Writer.Status()
	log.Println("latency:", latency, "method:", c.Request.Method, "path:", c.Request.URL.Path, "status:", status)
}

func LoggerUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	t := time.Now()
	h, err := handler(ctx, req)
	latency := time.Since(t)
	log.Printf("Request - Method:%s\tDuration:%s\tError:%v\n\n", // no lint
		info.FullMethod, //no lint
		latency,         //no lint
		err)             //no lint
	return h, err
}

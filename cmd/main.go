package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/app"
	grpcPort "homework9/internal/ports/grpc"
	"homework9/internal/ports/httpgin"
	"homework9/middleware"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	grpcPortAdr = ":50054"
	httpPort    = ":9000"
)

func main() {
	lis, err := net.Listen("tcp", grpcPortAdr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	repoAds := adrepo.New()
	repoUsers := userrepo.New()
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(middleware.LoggerUnaryServerInterceptor, middleware.PanicUnaryInterceptor))
	svc := grpcPort.NewService(app.NewApp(repoAds, repoUsers))
	grpcPort.RegisterAdServiceServer(grpcServer, svc)

	httpServer := httpgin.NewHTTPServer(httpPort, app.NewApp(repoAds, repoUsers))

	eg, ctx := errgroup.WithContext(context.Background())

	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	// run grpc server
	eg.Go(func() error {
		log.Printf("starting grpc server, listening on %s\n", grpcPortAdr)
		defer log.Printf("close grpc server listening on %s\n", grpcPortAdr)

		errCh := make(chan error)

		defer func() {
			grpcServer.GracefulStop()
			_ = lis.Close()

			close(errCh)
		}()

		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("grpc server can't listen and serve requests: %w", err)
		}
	})

	// run http server
	eg.Go(func() error {
		log.Printf("starting http server, listening on %s\n", httpPort)
		defer log.Printf("close http server listening on %s\n", httpPort)

		errCh := make(chan error)

		defer func() {
			shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(shCtx); err != nil {
				log.Printf("can't close http server listening on %s: %s", httpServer, err.Error()) // nolint
			}

			close(errCh)
		}()

		go func() {
			if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("http server can't listen and serve requests: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
		return
	}

	log.Println("servers were successfully shutdown")
}

package tests

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework9/internal/adapters/userrepo"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/app"
	grpcPort "homework9/internal/ports/grpc"
)

func TestGRRPCCreateUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure()) //nolint
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Oleg", Email: "alncalknd"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Nickname)
}

func TestGRRPCGetUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure()) //nolint:all
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Oleg", Email: "alncalknd"})
	assert.NoError(t, err, "client.CreateUser")
	res, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: res.Id})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Nickname)
}

func TestGRRPCDeleteUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure()) //nolint:all
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Oleg", Email: "alncalknd"})
	assert.NoError(t, err, "client.CreateUser")
	res, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: res.Id})
	assert.NoError(t, err, "client.GetUser")
	assert.Equal(t, "Oleg", res.Nickname)
	_, err = client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: res.Id})
	assert.NoError(t, err, "client.DeleteUser")
	_, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: res.Id})
	assert.Error(t, err, "client.GetUser")
}

func TestGRRPCCreateAd(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure()) //nolint:all
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Oleg", Email: "alncalknd"})
	assert.NoError(t, err, "client.CreateUser")
	resAd, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: res.Id, Title: "hello", Text: "world"})
	assert.NoError(t, err, "client.CreateAd")
	assert.Equal(t, "hello", resAd.Title)
	assert.Equal(t, "world", resAd.Text)
}

func TestGRRPCChangeAdStatus(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure()) //nolint:all
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Oleg", Email: "alncalknd"})
	assert.NoError(t, err, "client.CreateUser")
	resAd, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: res.Id, Title: "hello", Text: "world"})
	assert.NoError(t, err, "client.CreateAd")
	assert.Equal(t, "hello", resAd.Title)
	assert.Equal(t, false, resAd.Published)

	resAd, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: res.Id, Published: true, AdId: resAd.Id})
	assert.NoError(t, err, "client.CreateAd")
	assert.Equal(t, "hello", resAd.Title)
	assert.Equal(t, true, resAd.Published)
}

func TestGRRPCUpdateAd(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure()) //nolint:all
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Oleg", Email: "alncalknd"})
	assert.NoError(t, err, "client.CreateUser")
	resAd, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: res.Id, Title: "hello", Text: "world"})
	assert.NoError(t, err, "client.CreateAd")
	assert.Equal(t, "hello", resAd.Title)
	assert.Equal(t, false, resAd.Published)

	resAd, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{UserId: res.Id, AdId: resAd.Id, Title: "привет", Text: "мир"})
	assert.NoError(t, err, "client.CreateAd")
	assert.Equal(t, "привет", resAd.Title)
	assert.Equal(t, "мир", resAd.Text)
}

func TestGRRPCListAds(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure()) //nolint:all
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Oleg", Email: "alncalknd"})
	assert.NoError(t, err, "client.CreateUser")

	resList, err := client.ListAds(ctx, &emptypb.Empty{})
	assert.NoError(t, err, "client.ListAd")
	assert.Len(t, resList.List, 0)

	resAd, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: res.Id, Title: "hello", Text: "world"})
	assert.NoError(t, err, "client.CreateAd")

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: res.Id, Published: true, AdId: resAd.Id})
	assert.NoError(t, err, "client.ChangeAdStatus")
	resAd, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: res.Id, Title: "hello2", Text: "world2"})
	assert.NoError(t, err, "client.CreateAd")

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: res.Id, Published: true, AdId: resAd.Id})
	assert.NoError(t, err, "client.ChangeAdStatus")

	resList, err = client.ListAds(ctx, &emptypb.Empty{})
	assert.NoError(t, err, "client.ListAd")
	assert.Len(t, resList.List, 2)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: res.Id, Title: "hello3", Text: "world3"})
	assert.NoError(t, err, "client.CreateAd")
	resList, err = client.ListAds(ctx, &emptypb.Empty{})
	assert.NoError(t, err, "client.ListAd")
	assert.Len(t, resList.List, 2)
}

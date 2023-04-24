package grpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/app"
	grpcPort "homework9/internal/ports/grpc"
	"net"
	"testing"
	"time"
)

func getTestClient(t *testing.T) (grpcPort.AdServiceClient, context.Context) {
	adApp := app.NewApp(adrepo.New())
	srv, lis := grpcPort.TestNewGRPCServer(1024*1024, adApp)

	t.Cleanup(func() {
		_ = lis.Close()
	})

	t.Cleanup(func() {
		srv.Stop()
	})

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

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		_ = conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	return client, ctx
}

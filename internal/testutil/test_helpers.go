package testutil

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc"

	"github.com/rickyreddygari/walletsdk/internal/app"
)

func NewTestServer(t *testing.T) (*httptest.Server, *grpc.ClientConn, func()) {
	t.Helper()

	container, err := app.NewContainer()
	if err != nil {
		t.Fatalf("bootstrap container: %v", err)
	}

	server := httptest.NewServer(container.HTTPServer)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen gRPC: %v", err)
	}

	go func() {
		if err := container.GRPCServer.Serve(lis); err != nil {
			t.Logf("grpc server stopped: %v", err)
		}
	}()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dial grpc: %v", err)
	}

	cleanup := func() {
		conn.Close()
		container.GRPCServer.Stop()
		server.Close()
	}

	return server, conn, cleanup
}

func MustDo(t *testing.T, client *http.Client, req *http.Request) *http.Response {
	t.Helper()

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("http request failed: %v", err)
	}
	return resp
}

func Background() context.Context {
	return context.Background()
}

package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/rickyreddygari/walletsdk/internal/app"
)

func main() {
	container, err := app.NewContainer()
	if err != nil {
		log.Fatalf("bootstrap container: %v", err)
	}

	httpSrv := &http.Server{
		Addr:    ":" + container.Config.HTTPPort,
		Handler: container.HTTPServer,
	}

	lis, err := net.Listen("tcp", ":"+container.Config.GRPCPort)
	if err != nil {
		log.Fatalf("listen grpc: %v", err)
	}

	var wg sync.WaitGroup
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("http listening on %s", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("http server error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("grpc listening on %s", lis.Addr())
		if err := container.GRPCServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Printf("grpc server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown error: %v", err)
	}
	container.GRPCServer.GracefulStop()
	lis.Close()

	wg.Wait()
}

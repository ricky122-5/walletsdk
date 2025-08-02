package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/rickyreddygari/walletsdk/internal/app"
)

func main() {
	container, err := app.NewContainer()
	if err != nil {
		log.Fatalf("bootstrap container: %v", err)
	}

	server := container.HTTPServer
	port := container.Config.HTTPPort

	go func() {
		log.Printf("listening on :%s", port)
		if err := http.ListenAndServe(":"+port, server); err != nil && err != http.ErrServerClosed {
			log.Fatalf("serve http: %v", err)
		}
	}()

	waitForShutdown()
}

func waitForShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Printf("received signal %s, shutting down\n", sig)
}

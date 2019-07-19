package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func handle(conn net.Conn, ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	if _, err := io.Copy(conn, conn); err != nil {
		log.Println(err)
	}
}

func serve(l net.Listener, ctx context.Context) error {
	var wg sync.WaitGroup
	var conn net.Conn
	var err error
	for {
		conn, err = l.Accept()
		if err != nil {
			break
		}
		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			handle(c, ctx)
		}(conn)
	}
	wg.Wait()
	return err
}

func listenAndServe(network, address string, ctx context.Context, ready chan struct{}) error {
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	close(ready)
	go func() {
		<-ctx.Done()
		l.Close()
	}()
	defer l.Close()

	return serve(l, ctx)
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigCh
		cancel()
	}()
	ready := make(chan struct{})
	go func() {
		<-ready
		log.Println("ready to accept connections")
	}()

	log.Fatal(listenAndServe("tcp", ":9090", ctx, ready))
}

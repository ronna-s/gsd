package main

import (
	"bufio"
	"context"
	"net"
	"testing"
)

func TestListenAndServe(t *testing.T) {
	const (
		network = "tcp"
		address = ":8080"
		message = "hello"
	)
	ready := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		listenAndServe(network, address, ctx, ready)
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("listenAndServe exited before ready, problem with IO?")
	case <-ready:
	}

	conn, err := net.Dial(network, address)
	if err != nil {
		t.Fatalf("unexpected error '%s'", err.Error())
	}
	conn.Write([]byte(message + "\n"))

	s := bufio.NewScanner(conn)
	s.Scan()
	if s.Text() != message {
		t.Fatalf("unexpected message received from server: '%s'", s.Text())
	}
	cancel()
	<-done
}

//wrapper for the Pipe connection to count closed
type c struct {
	net.Conn
	closed bool
}

func (c *c) Close() error {
	c.closed = true
	return c.Conn.Close()
}

func TestHandleTerminates(t *testing.T) {
	const message = "hello"
	server, client := net.Pipe()
	serverConn, clientConn := &c{Conn: server}, &c{Conn: client}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		handle(serverConn, ctx)
		close(done)
	}()
	s := bufio.NewScanner(clientConn)
	clientConn.Write([]byte(message + "\n"))
	s.Scan()

	if s.Text() != message {
		t.Fatalf("unexpected message received from server: '%s'", s.Text())
	}
	cancel()
	<-done
	if !serverConn.closed {
		t.Fatal("expected connection to be closed once")
	}
}

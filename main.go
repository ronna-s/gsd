package main

import (
	"io"
	"log"
	"net"
	"sync"
)

func handle(conn net.Conn) {
	if _, err := io.Copy(conn, conn); err != nil {
		log.Println(err)
	}
	conn.Close()
}

func serve(l net.Listener) error {
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
			handle(c)
		}(conn)
	}
	wg.Wait()
	return err
}

func listenAndServe(network, address string) error {
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	defer l.Close()

	return serve(l)
}

func main() {
	log.Fatal(listenAndServe("tcp", ":9090"))
}

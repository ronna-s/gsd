package main

import (
	"io"
	"log"
	"net"
	"sync"
)

func handle(c net.Conn) {
	if _, err := io.Copy(c, c); err != nil {
		log.Println(err)
	}
	c.Close()
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

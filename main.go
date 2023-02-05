package main

import (
	distributedcache "DistributedCache/cache"
	"log"
	"net"
	"time"
)

func main() {
	opts := ServerOpts{
		ListenAddr: ":8080",
		IsLeader:   true,
	}

	go func() {
		time.Sleep(time.Second * 2)
		conn, err := net.Dial("tcp", ":8080")
		if err != nil {
			log.Fatal(err)
		}

		conn.Write([]byte("SET Foo Bar 1000"))

	}()
	server := NewServer(opts, distributedcache.NewCache())
	server.Start()
}

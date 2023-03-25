package main

import (
	distributedcache "DistributedCache/cache"
	"DistributedCache/client"
	"context"
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	listenAddr := flag.String("listenaddr", ":8080", "listen address of the server")
	leaderAddr := flag.String("leaderaddr", "", "listen address of the leader")
	flag.Parse()

	opts := ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	go func() {
		client, err := client.NewClient(":8080", client.Options{})
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second * 2)
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond * 500)
			SendCommand(client)
		}

		client.Close()
		time.Sleep(time.Second * 2)
	}()

	server := NewServer(opts, distributedcache.NewCache())
	server.Start()
}

func SendCommand(c *client.Client) {

	resp, err := c.Set(context.Background(), []byte("Foo"), []byte("Bar"), 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("resp:", resp)

}

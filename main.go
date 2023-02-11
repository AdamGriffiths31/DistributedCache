package main

import (
	distributedcache "DistributedCache/cache"
	"flag"
	"log"
	"net"
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

	// go func() {
	// 	time.Sleep(time.Second * 2)
	// 	conn, err := net.Dial("tcp", ":8080")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	conn.Write([]byte("SET Foo Bar 999999999999"))
	// 	time.Sleep(time.Second * 2)
	// 	conn.Write([]byte("GET Foo"))
	// 	buf := make([]byte, 1000)
	// 	n, _ := conn.Read(buf)
	// 	fmt.Println(string(buf[:n]))
	// }()

	server := NewServer(opts, distributedcache.NewCache())
	server.Start()
}

func SendCommand() {
	cmd := &CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	conn.Write(cmd.Bytes())
}

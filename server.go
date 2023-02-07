package main

import (
	distributedcache "DistributedCache/cache"
	"context"
	"fmt"
	"log"
	"net"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	cache     distributedcache.Cacher
	followers map[net.Conn]struct{}
}

func NewServer(opts ServerOpts, c distributedcache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
		followers:  make(map[net.Conn]struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server staring on port [%s]\n", s.ListenAddr)

	if !s.IsLeader {
		go func() {
			conn, err := net.Dial("tcp", s.LeaderAddr)
			fmt.Printf("connected with leader: [%v]", s.LeaderAddr)
			if err != nil {
				log.Printf("Start dial [%v] error %v", s.LeaderAddr, err)
			}
			s.handleConn(conn)
		}()
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	if s.IsLeader {
		s.followers[conn] = struct{}{}
	}

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn read error: %s\n", err)
			break
		}

		msg := buf[:n]
		fmt.Println(string(msg))

		go s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {

	msg, err := ParseMessage(rawCmd)
	if err != nil {
		log.Printf("handleCommandd error: %s\n", err)
		return
	}

	switch msg.Cmd {
	case CMDSet:
		err = s.handleSetCommand(conn, msg)
	case CMDGet:
		err = s.handleGETCommand(conn, msg)

	}

	if err != nil {
		fmt.Println("handleCommand error:", err)
		conn.Write([]byte(err.Error()))
	}
}

func (s *Server) handleSetCommand(conn net.Conn, msg *Message) error {
	fmt.Printf("handle SET command %v\n", msg)

	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}

	go s.sendToFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) handleGETCommand(conn net.Conn, msg *Message) error {
	fmt.Printf("handle GET command %v\n", msg)

	value, err := s.cache.Get(msg.Key)
	if err != nil {
		return err
	}

	conn.Write(value)
	return nil
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	for conn := range s.followers {
		fmt.Printf("forwarding to followers: %v", msg)
		_, err := conn.Write(msg.ToBytes())
		if err != nil {
			log.Printf("sendToFollowers error: %v", err)
			continue
		}
	}
	return nil
}

package main

import (
	distributedcache "DistributedCache/cache"
	"DistributedCache/protocol"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	cache distributedcache.Cacher
}

func NewServer(opts ServerOpts, c distributedcache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server staring on port [%s]\n", s.ListenAddr)

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

	log.Println("Connection made:", conn.RemoteAddr())

	for {
		cmd, err := protocol.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("handleConn - parse command error: %s\n", err)
			break
		}
		fmt.Println("cmd", cmd)

		go s.handleCommand(conn, cmd)
	}
	log.Println("Connection closed:", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *protocol.CommandSet:
		s.handleSetCommand(conn, v)
	case *protocol.CommandGet:
	}
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *protocol.CommandSet) error {
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		return err
	}
	return nil
}

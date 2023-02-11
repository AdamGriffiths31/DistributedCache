package main

// import (
// 	"errors"
// 	"fmt"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// type Command string

// const (
// 	CMDSet Command = "SET"
// 	CMDGet Command = "GET"
// )

// type Message struct {
// 	Cmd   Command
// 	Key   []byte
// 	Value []byte
// 	TTL   time.Duration
// }

// func (m *Message) ToBytes() []byte {
// 	switch m.Cmd {
// 	case CMDSet:
// 		cmd := fmt.Sprintf("%s %s %s %d", m.Cmd, m.Key, m.Value, m.TTL)
// 		return []byte(cmd)
// 	case CMDGet:
// 		cmd := fmt.Sprintf("%s %s", m.Cmd, m.Key)
// 		return []byte(cmd)
// 	default:
// 		panic("unkown command")
// 	}
// }

// func ParseMessage(raw []byte) (*Message, error) {
// 	rawString := string(raw)
// 	parts := strings.Split(rawString, " ")
// 	if len(parts) == 0 {
// 		return nil, errors.New("invalid handleCommand")
// 	}

// 	msg := &Message{
// 		Cmd: Command(parts[0]),
// 		Key: []byte(parts[1]),
// 	}

// 	if msg.Cmd == CMDSet {
// 		if len(parts) != 4 {
// 			return nil, errors.New("invalid handleCommand SET invalid parts")
// 		}

// 		msg.Value = []byte(parts[2])
// 		ttl, err := strconv.Atoi(parts[3])
// 		if err != nil {
// 			return nil, errors.New("invalid handleCommand SET invalid TTL")
// 		}
// 		msg.TTL = time.Duration(ttl)
// 	}

// 	return msg, nil
// }

package main

import (
	"bufio"
	"fmt"
	"github.com/Maddoxx88/garnet/protocol"
	"github.com/Maddoxx88/garnet/store"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const GarnetVersion = "v0.1.0"

func handleConnection(conn net.Conn, db *store.GarnetStore) {
	defer conn.Close()
	conn.Write([]byte("Welcome to Garnet!\n"))

	// Replace scanner loop:
	reader := bufio.NewReader(conn)
	for {
		cmd, err := protocol.ParseRESP(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			conn.Write([]byte("-ERR invalid command\r\n"))
			continue
		}

		if len(cmd) == 0 {
			conn.Write([]byte("-ERR empty command\r\n"))
			continue
		}

		switch strings.ToUpper(cmd[0]) {
		case "PING":
			conn.Write([]byte("+PONG\r\n"))

		case "SET":
			if len(cmd) < 3 {
				conn.Write([]byte("-ERR wrong number of arguments for SET\r\n"))
				continue
			}
			key, val := cmd[1], cmd[2]
			ttl := 0
			if len(cmd) == 5 && strings.ToUpper(cmd[3]) == "EX" {
				ttl, _ = strconv.Atoi(cmd[4])
			}
			db.Set(key, val, ttl)
			conn.Write([]byte("+OK\r\n"))

		case "GET":
			if len(cmd) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for GET\r\n"))
				continue
			}
			val, ok := db.Get(cmd[1])
			if ok {
				conn.Write([]byte("$" + strconv.Itoa(len(val)) + "\r\n" + val + "\r\n"))
			} else {
				conn.Write([]byte("$-1\r\n"))
			}

		default:
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}

func main() {
	db := store.New()
	fmt.Printf("🔴 Garnet %s started at %s\n", GarnetVersion, time.Now().Format(time.RFC1123))
	db.StartTTLLoop(1)

	// Start TCP server
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	fmt.Println("Garnet TCP server is listening on port 6379...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn, db)
	}
}

package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/Maddoxx88/garnet/store"
)

func handleConnection(conn net.Conn, db *store.GarnetStore) {
	defer conn.Close()
	conn.Write([]byte("Welcome to Garnet!\n"))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		input := scanner.Text()
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		switch strings.ToUpper(parts[0]) {
		case "PING":
			conn.Write([]byte("PONG\n"))
		case "SET":
			if len(parts) < 3 {
				conn.Write([]byte("Usage: SET key value [EX seconds]\n"))
				continue
			}
			key, val := parts[1], parts[2]
			ttl := 0
			if len(parts) == 5 && strings.ToUpper(parts[3]) == "EX" {
				fmt.Sscanf(parts[4], "%d", &ttl)
			}
			db.Set(key, val, ttl)
			conn.Write([]byte("OK\n"))
		case "GET":
			if len(parts) != 2 {
				conn.Write([]byte("Usage: GET key\n"))
				continue
			}
			if val, ok := db.Get(parts[1]); ok {
				conn.Write([]byte(val + "\n"))
			} else {
				conn.Write([]byte("(nil)\n"))
			}
		case "DEL":
			if len(parts) != 2 {
				conn.Write([]byte("Usage: DEL key\n"))
				continue
			}
			if db.Del(parts[1]) {
				conn.Write([]byte("1\n"))
			} else {
				conn.Write([]byte("0\n"))
			}
		case "EXISTS":
			if len(parts) != 2 {
				conn.Write([]byte("Usage: EXISTS key\n"))
				continue
			}
			if db.Exists(parts[1]) {
				conn.Write([]byte("1\n"))
			} else {
				conn.Write([]byte("0\n"))
			}
		case "KEYS":
			keys := db.Keys()
			for _, k := range keys {
				conn.Write([]byte(k + "\n"))
			}
		case "FLUSHALL":
			db.FlushAll()
			conn.Write([]byte("OK\n"))
		case "HELP", "/H":
			conn.Write([]byte("Available Commands:\n"))
			conn.Write([]byte("SET key value [EX seconds]\n"))
			conn.Write([]byte("GET key\n"))
			conn.Write([]byte("DEL key\n"))
			conn.Write([]byte("EXISTS key\n"))
			conn.Write([]byte("KEYS\n"))
			conn.Write([]byte("FLUSHALL\n"))
			conn.Write([]byte("PING\n"))
			conn.Write([]byte("EXIT or QUIT or /exit\n"))
		case "EXIT", "QUIT", "/EXIT":
			conn.Write([]byte("Goodbye 👋\n"))
			return
		default:
			conn.Write([]byte("Unknown command. Type HELP for a list.\n"))
		}
	}
}

func main() {
	db := store.New()
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

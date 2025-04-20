package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Maddoxx88/garnet/protocol"
	"github.com/Maddoxx88/garnet/store"
)

const GarnetVersion = "v0.1.0"

func handleConnection(conn net.Conn, db *store.GarnetStore) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		cmd, err := protocol.ParseRESP(reader)
		if err != nil {
			writer.WriteString("-ERR " + err.Error() + "\r\n")
			writer.Flush()
			return
		}

		if len(cmd) == 0 {
			writer.WriteString("-ERR empty command\r\n")
			writer.Flush()
			continue
		}

		switch strings.ToUpper(cmd[0]) {
		case "PING":
			writer.WriteString("+PONG\r\n")

		case "SET":
			if len(cmd) < 3 {
				writer.WriteString("-ERR wrong number of arguments for SET\r\n")
			} else {
				key := cmd[1]
				val := cmd[2]
				ttl := 0
				if len(cmd) == 5 && strings.ToUpper(cmd[3]) == "EX" {
					ttl, _ = strconv.Atoi(cmd[4])
				}
				db.Set(key, val, ttl)
				writer.WriteString("+OK\r\n")
			}

		case "GET":
			if len(cmd) != 2 {
				writer.WriteString("-ERR wrong number of arguments for GET\r\n")
			} else {
				val, ok := db.Get(cmd[1])
				if ok {
					writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val))
				} else {
					writer.WriteString("$-1\r\n")
				}
			}

		case "GHELP":
			lines := []string{
				"SET key value [EX seconds] : Set a key with optional TTL",
				"GET key                    : Get the value of a key",
				"DEL key                    : Delete a key",
				"EXISTS key                 : Check if a key exists",
				"KEYS                       : List all keys",
				"FLUSHALL                   : Delete all keys",
				"PING                       : Returns PONG",
				"QUIT                       : Exit the connection",
			}
			writer.WriteString(fmt.Sprintf("*%d\r\n", len(lines)))
			for _, line := range lines {
				writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(line), line))
			}

		case "EXIT", "QUIT":
			writer.WriteString("+Goodbye 👋\r\n")
			writer.Flush()
			return

		default:
			writer.WriteString("-ERR unknown command\r\n")
		}

		// ✅ Always flush the response after every command
		writer.Flush()
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

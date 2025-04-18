package main

import (
	"bufio"
	"fmt"
	"github.com/Maddoxx88/garnet/store"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	db := store.New()
	db.StartTTLLoop(1 * time.Second)

	fmt.Println("Garnet v0.1 ready to accept commands. Try: SET key value EX 5, GET key")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "HELP", "/h":
			fmt.Println("Available Commands:")
			fmt.Println("  SET key value [EX seconds]   → Set a key with optional expiration (TTL)")
			fmt.Println("  GET key                      → Get the value of a key")
			fmt.Println("  DEL key                      → Delete a key")
			fmt.Println("  EXISTS key                   → Check if a key exists")
			fmt.Println("  KEYS                         → List all keys")
			fmt.Println("  FLUSHALL                     → Delete all keys")
			fmt.Println("  PING                         → Health check (returns PONG)")
			fmt.Println("  HELP or /h                   → Show this help message")
			fmt.Println("  EXIT or QUIT or /exit        → Exit Garnet")
		case "EXIT", "QUIT", "/exit":
			fmt.Println("Goodbye 👋")
			os.Exit(0)
		case "SET":
			if len(parts) < 3 {
				fmt.Println("Usage: SET key value [EX seconds]")
				continue
			}

			key := parts[1]
			val := parts[2]
			ttl := 0

			if len(parts) == 5 && strings.ToUpper(parts[3]) == "EX" {
				if parsed, err := strconv.Atoi(parts[4]); err == nil {
					ttl = parsed
				} else {
					fmt.Println("Invalid TTL value")
					continue
				}
			}

			db.Set(key, val, ttl)
			fmt.Println("OK")
		case "GET":
			if len(parts) != 2 {
				fmt.Println("Usage: GET key")
				continue
			}
			if val, ok := db.Get(parts[1]); ok {
				fmt.Println(val)
			} else {
				fmt.Println("(nil)")
			}
		case "DEL":
			if len(parts) != 2 {
				fmt.Println("Usage: DEL key")
				continue
			}
			deleted := db.Del(parts[1])
			if deleted {
				fmt.Println("key deleted ✅")
			} else {
				fmt.Println("error deleting key ❌")
			}

		case "EXISTS":
			if len(parts) != 2 {
				fmt.Println("Usage: EXISTS key")
				continue
			}
			if db.Exists(parts[1]) {
				fmt.Println("yes")
			} else {
				fmt.Println("no")
			}

		case "KEYS":
			keys := db.Keys()
			for _, k := range keys {
				fmt.Println(k)
			}

		case "FLUSHALL":
			db.FlushAll()
			fmt.Println("OK")

		case "PING":
			fmt.Println("PONG")

		default:
			fmt.Println("Unknown command. Type HELP or /h to see available commands.")

		}
	}
}

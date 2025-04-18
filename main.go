package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Maddoxx88/garnet/store"
)

func main() {
	db := store.New()
	fmt.Println("Garnet ready. Type commands (SET key val EX, GET key):")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
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
			fmt.Println("Unknown command:", parts[0])
		}
	}
}

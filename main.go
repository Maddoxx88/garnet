package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Maddoxx88/garnet/store"
)

func main() {
	db := store.New()
	fmt.Println("Garnet ready. Type commands (SET key val, GET key):")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "SET":
			if len(parts) != 3 {
				fmt.Println("Usage: SET key value")
				continue
			}
			db.Set(parts[1], parts[2])
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
		default:
			fmt.Println("Unknown command:", parts[0])
		}
	}
}

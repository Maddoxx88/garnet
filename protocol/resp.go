package protocol

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// ParseRESP reads and parses a RESP-formatted command
func ParseRESP(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if len(line) == 0 || line[0] != '*' {
		return nil, fmt.Errorf("invalid RESP array")
	}

	count, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil {
		return nil, fmt.Errorf("invalid array count: %v", err)
	}

	parts := make([]string, 0, count)
	for i := 0; i < count; i++ {
		_, err := reader.ReadString('\n') // skip $N line
		if err != nil {
			return nil, err
		}

		part, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		parts = append(parts, strings.TrimSpace(part))
	}

	return parts, nil
}

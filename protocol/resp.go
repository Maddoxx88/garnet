package protocol

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func ParseRESP(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if len(line) == 0 || line[0] != '*' {
		return nil, fmt.Errorf("invalid RESP format")
	}

	count, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil {
		return nil, err
	}

	parts := make([]string, 0, count)
	for i := 0; i < count; i++ {
		prefix, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if len(prefix) == 0 || prefix[0] != '$' {
			return nil, fmt.Errorf("invalid bulk string prefix")
		}

		length, err := strconv.Atoi(strings.TrimSpace(prefix[1:]))
		if err != nil {
			return nil, err
		}

		data := make([]byte, length+2) // +2 for \r\n
		_, err = reader.Read(data)
		if err != nil {
			return nil, err
		}

		parts = append(parts, string(data[:length]))
	}

	return parts, nil
}

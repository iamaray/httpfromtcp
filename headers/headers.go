package header

import (
	"bytes"
	"fmt"
	"strings"
)

func isToken(str []byte) bool {
	for _, ch := range str {
		found := false
		if ch >= '0' && ch <= '9' ||
			ch >= 'a' && ch <= 'z' ||
			ch >= 'A' && ch <= 'Z' {
			found = true
		}
		switch ch {
		case '#', '$', '%', '&', '\'', '*', '+', '.', '^', '_', '`', '|', '~':
			found = true
		}
		if !found {
			return false
		}
	}
	return true
}

type Headers map[string]string

func (h Headers) Get(name string) string {
	return h[strings.ToLower(name)]
}

func (h Headers) Set(name, value string) {
	name = strings.ToLower(name)

	if v, ok := h[name]; ok {
		h[name] = fmt.Sprintf("%s,%s", v, value)
	} else {
		h[name] = value
	}
}

func NewHeaders() Headers {
	return map[string]string{}
}

const RN string = "\r\n"

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed header")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}
	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	done := false
	read := 0

	for {
		idx := bytes.Index(data[read:], []byte(RN))

		if idx == -1 {
			break
		}
		// EMPTY HEADER
		if idx == 0 {
			done = true
			read += len(RN)
			break
		}
		// fmt.Printf("header: \"%s\"\n", string(data[read:idx]))
		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}
		if !isToken([]byte(name)) {
			return 0, false, fmt.Errorf("malformed header name: %s", name)
		}

		read += idx + len(RN)
		h.Set(name, value)
	}

	return read, done, nil
}

package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var ErrMalformedStartLine = fmt.Errorf("malformed starting line")
var ErrIncompleteStartLine = fmt.Errorf("malformed start line")
var ErrMalformedRequestLine = fmt.Errorf("malformed request line")
var ErrWrongHttpVersion = fmt.Errorf("HTTP version not supported")

const SEPARATOR = "\r\n"

func parseRequestLine(s string) (*RequestLine, string, error) {
	idx := strings.Index(s, SEPARATOR)
	if idx == -1 {
		return nil, s, nil
	}

	startLine := s[:idx]
	restOfMsg := s[idx+len(SEPARATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, restOfMsg, ErrMalformedStartLine
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, restOfMsg, ErrMalformedRequestLine
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	return rl, restOfMsg, nil
}

func RequestFromHeader(reader io.Reader) (*Request, error) {
	var err error

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("unable to io.ReadAll"),
			err,
		)
	}

	str := string(data)
	rl, _, err := parseRequestLine(str)

	return &Request{
		RequestLine: *rl,
	}, err
}

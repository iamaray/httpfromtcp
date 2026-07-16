package request

import (
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}
type parserState string

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

var ErrMalformedStartLine = fmt.Errorf("malformed starting line")
var ErrIncompleteStartLine = fmt.Errorf("malformed start line")
var ErrMalformedRequestLine = fmt.Errorf("malformed request line")
var ErrWrongHttpVersion = fmt.Errorf("HTTP version not supported")

const SEPARATOR = "\r\n"

func parseRequestLine(s string) (*RequestLine, int, error) {
	idx := strings.Index(s, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := s[:idx]
	read := idx + len(SEPARATOR)

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, 0, ErrMalformedStartLine
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, 0, ErrMalformedRequestLine
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	return rl, read, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		switch r.state {
		case StateInit:
			rl, n, err := parseRequestLine(string(data[read:]))
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.state = StateDone
		case StateDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func RequestFromHeader(reader io.Reader) (*Request, error) {
	// var err error
	req := newRequest()
	// NOTE: buffer could get overrun
	buf := make([]byte, 1024)
	bufLen := 0

	for !req.done() {
		n, err := reader.Read(buf[bufLen:])
		//TODO what to do here
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := req.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return req, nil
}

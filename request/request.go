package request

import (
	"bytes"
	"fmt"
	"httpfromtcp/headers"
	"io"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateDone    parserState = "done"
	StateError   parserState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	State       parserState
}

func newRequest() *Request {
	return &Request{
		State:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

var ErrMalformedStartLine = fmt.Errorf("malformed starting line")
var ErrIncompleteStartLine = fmt.Errorf("malformed start line")
var ErrMalformedRequestLine = fmt.Errorf("malformed request line")
var ErrWrongHttpVersion = fmt.Errorf("HTTP version not supported")

const SEPARATOR = "\r\n"

func parseRequestLine(s []byte) (*RequestLine, int, error) {
	idx := bytes.Index(s, []byte(SEPARATOR))
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := s[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrMalformedStartLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 ||
		!bytes.Equal(httpParts[0], []byte("HTTP")) ||
		!bytes.Equal(httpParts[1], []byte("1.1")) {
		return nil, 0, ErrMalformedRequestLine
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	return rl, read, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		currentData := data[read:]
		switch r.State {
		case StateError:
			return 0, fmt.Errorf("request in error state")
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				r.State = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.State = StateHeaders

		case StateHeaders:
			headersMap := headers.NewHeaders()
			n, done, err := headersMap.Parse(currentData)
			if err != nil {
				r.State = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.Headers = headersMap
			read += n
			if done {
				r.State = StateDone
			}

		case StateDone:
			break outer
		default:
			panic("something wrong :(")
		}

	}
	return read, nil
}

func (r *Request) done() bool {
	return r.State == StateDone || r.State == StateError
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

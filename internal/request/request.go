package request

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/lucoand/httpfromtcp/internal/headers"
)

const requestStateInitialized int = 0
const requestStateParsingHeaders = 1
const requestStateDone int = 2
const bufferSize = 8

type Request struct {
	Headers     headers.Headers
	RequestLine RequestLine
	state       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func newRequest() *Request {
	return &Request{
		state:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}
}

func isUpper(s string) bool {
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}

func (r *Request) parseSingle(dataString string) (int, error) {
	requestLine, n, err := parseRequestLine(dataString)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, nil
	}
	r.RequestLine = requestLine
	r.state = requestStateParsingHeaders
	return n, nil
}

func (r *Request) parseHeaders(data []byte) (int, error) {
	n, done, err := r.Headers.Parse(data)
	if err != nil {
		return 0, err
	}
	if done {
		r.state = requestStateDone
	}
	return n, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		fmt.Println("Parsing Request Line")
		dataString := string(data)
		return r.parseSingle(dataString)
	case requestStateParsingHeaders:
		fmt.Println("Parsing Headers")
		return r.parseHeaders(data)
	case requestStateDone:
		return 0, fmt.Errorf("Error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("Error: unknown state")
	}
}

func parseRequestLine(dataString string) (RequestLine, int, error) {
	if !strings.Contains(dataString, "\r\n") {
		return RequestLine{}, 0, nil
	}
	lines := strings.Split(dataString, "\r\n")
	parts := strings.Split(lines[0], " ")
	// fmt.Printf("Parts: ")
	// for _, part := range parts {
	// 	fmt.Printf("%s ", part)
	// }
	// fmt.Printf("\nparts length: %d\n", len(parts))
	if len(parts) != 3 {
		return RequestLine{}, 0, fmt.Errorf("incomplete HTTP request-line")
	}
	method := parts[0]
	if !isUpper(method) {
		return RequestLine{}, 0, fmt.Errorf("HTTP method not properly capitalized")
	}

	version := parts[2]
	if version != "HTTP/1.1" {
		return RequestLine{}, 0, fmt.Errorf("Unsupported HTTP version. Currently only HTTP/1.1 is supported.")
	}
	versionParts := strings.Split(version, "/")
	target := parts[1]
	return RequestLine{
		HttpVersion:   versionParts[1],
		Method:        method,
		RequestTarget: target,
	}, len(lines[0]) + 2, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	r := newRequest()
	readToIndex := 0
	numParsed := 0
	for r.state != requestStateDone {
		if len(buf) == cap(buf) {
			temp := make([]byte, len(buf)*2, cap(buf)*2)
			copy(temp, buf)
			buf = temp
		}
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		readToIndex += numBytesRead
		numParsed, err = r.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if numParsed == 0 {
			continue
		}
		temp := make([]byte, len(buf), cap(buf))
		copy(temp, buf[numParsed:])
		buf = temp
		readToIndex -= numParsed
	}
	if r.state != requestStateDone {
		return nil, fmt.Errorf("Parsing finished unexpectedly - incomplete request")
	}
	return r, nil
}

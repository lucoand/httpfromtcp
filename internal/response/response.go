package response

import (
	"fmt"
	"io"

	"github.com/lucoand/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK         StatusCode = iota // 200
	StatusBADREQUEST                   // 400
	StatusINTERNAL                     // 500
)

func writeErrorHelper(err error, n int, line []byte) error {
	if err != nil {
		return err
	}
	if n != len(line) {
		return fmt.Errorf("Status Line Write length mismatch")
	}
	return nil
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		line := []byte("HTTP/1.1 200 OK\r\n")
		n, err := w.Write(line)
		return writeErrorHelper(err, n, line)
	case StatusBADREQUEST:
		line := []byte("HTTP/1.1 400 Bad Request\r\n")
		n, err := w.Write(line)
		return writeErrorHelper(err, n, line)
	case StatusINTERNAL:
		line := []byte("HTTP/1.1 500 Internal Server Error\r\n")
		n, err := w.Write(line)
		return writeErrorHelper(err, n, line)
	default:
		lineString := "HTTP/1.1 " + fmt.Sprintf("%d \r\n", statusCode)
		line := []byte(lineString)
		n, err := w.Write(line)
		return writeErrorHelper(err, n, line)
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := make(map[string]string)
	h["content-length"] = fmt.Sprintf("%d", contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		headerString := k + ": " + v + "\r\n"
		headerBytes := []byte(headerString)
		n, err := w.Write(headerBytes)
		err = writeErrorHelper(err, n, headerBytes)
		if err != nil {
			return err
		}
	}
	headersEndString := "\r\n"
	headersEndBytes := []byte(headersEndString)
	n, err := w.Write(headersEndBytes)
	return writeErrorHelper(err, n, headersEndBytes)
}

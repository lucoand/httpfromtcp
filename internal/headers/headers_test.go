package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header with uppercase in field name
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Two valid headers with identical field names
	data = []byte("Host: localhost:1337\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069, localhost:1337", headers.Get("host"))
	assert.Equal(t, 22, n)
	assert.False(t, done)
	data = []byte("Host: remotehost:31415\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069, localhost:1337, remotehost:31415", headers.Get("hOST"))
	assert.Equal(t, 24, n)
	assert.False(t, done)

	// Test: Valid single header lowercase field name
	headers = NewHeaders()
	data = []byte("host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("host"))
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Host:       localhost:42069          \r\n\r\n")
	length := len("Host:       localhost:42069          \r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, length, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	data = []byte("Content-Type: application/json\r\nField-Name: field-data\r\n\r\n")
	n, done, err = headers.Parse(data)
	length = len("Content-Type: application/json\r\n")
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, length, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[n:])
	length = len("Field-Name: field-data\r\n")
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, length, n)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, "application/json", headers.Get("Content-Type"))
	assert.Equal(t, "field-data", headers.Get("Field-Name"))
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: No CRLF
	headers = NewHeaders()
	data = []byte("Host: localhost:42069")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Missing field name
	headers = NewHeaders()
	data = []byte(": localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Missing colon
	headers = NewHeaders()
	data = []byte("Host localhost-42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Missing field data
	headers = NewHeaders()
	data = []byte("Host:\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid character in field name
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("    Host:   localhost:42069   \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 32, n)
	assert.False(t, done)
	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	data = []byte("Content-Type: application/json\r\nAccept: */*\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 32, n)
	assert.False(t, done)
	num, done2, err2 := headers.Parse(data[n:])
	require.NoError(t, err2)
	assert.Equal(t, "application/json", headers["content-type"])
	assert.Equal(t, "*/*", headers["accept"])
	assert.Equal(t, 13, num)
	assert.False(t, done2)
	close := n + num
	num3, done3, err3 := headers.Parse(data[close:])
	assert.Equal(t, 2, num3)
	assert.NoError(t, err3)
	assert.True(t, done3)
	// Test: Invalid Character
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.False(t, done)
	// Test: Valid Character
	headers = NewHeaders()
	data = []byte("H$st: localhost:42069\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, "localhost:42069", headers["h$st"])
}

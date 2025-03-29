package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	t.Run("good headers parse with optional white space between : and val", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "localhost:42069", headers["host"])
		assert.Equal(t, 23, n)
		assert.False(t, done)
	})

	t.Run("good headers parse with no optional white space between : and val", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host:localhost:42069\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "localhost:42069", headers["host"])
		assert.Equal(t, 22, n)
		assert.False(t, done)
	})

	t.Run("valid headers with space in the middle of the key", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("       Content Length: 42069       \r\n\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		require.False(t, done)
		assert.Equal(t, 37, n)
		assert.Equal(t, headers["content length"], "42069")
	})

	t.Run("always set key header to lower case", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("       VAMO GREMIO-PORRA!: 42069       \r\n\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		require.False(t, done)
		assert.Equal(t, 41, n)
		assert.Equal(t, headers["vamo gremio-porra!"], "42069")
	})

	t.Run("return true when data is a only a crlf - the next line will be the body so parse headers is done", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		require.True(t, done)
		require.Zero(t, n)
	})

	t.Run("invalid headers with space between key name and :", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("       Host : localhost:42069       \r\n\r\n")

		n, done, err := headers.Parse(data)

		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
		assert.ErrorContains(t, err, "malformed headers key - got key ending with space - header key")
	})

	t.Run("invalid headers with no : separator", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("       Content-Length 42069       \r\n\r\n")

		n, done, err := headers.Parse(data)

		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
		assert.ErrorContains(t, err, "malformed headers - not find header separator : - headers:")
	})

	t.Run("invalid headers key char", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("H©st: localhost:42069\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
		assert.ErrorContains(t, err, "malformed headers - key header with not allowed char - header key: H©st")
	})

	t.Run("invalid headers key char wiht emoji text", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("┐( ͡◉ ͜ʖ ͡◉)┌: ligma?\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
		assert.ErrorContains(t, err, "malformed headers - key header with not allowed char - header key: ┐( ͡◉ ͜ʖ ͡◉)┌")
	})

	t.Run("invalid headers key char wiht emoji", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("🫦: 06??\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
		assert.ErrorContains(t, err, "malformed headers - key header with not allowed char - header key: 🫦")
	})

}

package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	t.Run("Good request line without request target method Get", func(t *testing.T) {
		reader := &chunkReader{
			data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 1,
		}

		r, err := ParseFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	})

	t.Run("Good request line with request target", func(t *testing.T) {
		data := "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"

		reader := &chunkReader{
			data:            data,
			numBytesPerRead: len(data),
		}

		r, err := ParseFromReader(reader)

		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	})

	t.Run("without http method", func(t *testing.T) {
		reader := &chunkReader{
			data:            "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 10,
		}

		r, err := ParseFromReader(reader)

		require.Nil(t, r)
		require.Error(t, err)
		require.ErrorContains(t, err, "request line has not 3 parts format")
	})

	t.Run("without request target", func(t *testing.T) {
		reader := &chunkReader{
			data:            "GET HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 8,
		}

		r, err := ParseFromReader(reader)

		require.Nil(t, r)
		require.Error(t, err)
		require.ErrorContains(t, err, "request line has not 3 parts format")
	})

	t.Run("without http version", func(t *testing.T) {
		reader := &chunkReader{
			data:            "GET /coffe\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 8,
		}

		r, err := ParseFromReader(reader)

		require.Nil(t, r)
		require.Error(t, err)
		require.ErrorContains(t, err, "request line has not 3 parts format")
	})

	t.Run("method malformed", func(t *testing.T) {
		reader := &chunkReader{
			data:            "get / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 8,
		}

		r, err := ParseFromReader(reader)

		require.Nil(t, r)
		require.Error(t, err)
		require.ErrorContains(t, err, "request method malformed method is not in all captal letter")
	})

	t.Run("method not mapped", func(t *testing.T) {
		reader := &chunkReader{
			data:            "PIZZA / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 8,
		}

		r, err := ParseFromReader(reader)

		require.Nil(t, r)
		require.Error(t, err)
		require.ErrorContains(t, err, "request method unsported - method got PIZZA - see AllMethods variable to suported methods")
	})

	t.Run("http version not suported", func(t *testing.T) {
		reader := &chunkReader{
			data:            "GET / HTTP/2.0\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 8,
		}

		r, err := ParseFromReader(reader)

		require.Nil(t, r)
		require.Error(t, err)
		require.ErrorContains(t, err, "unsoported http version - the httpVersion is HTTP/2.0 and only httpVersion suported is HTTP/1.1")
	})

	t.Run("http version malformed", func(t *testing.T) {
		reader := &chunkReader{
			data:            "GET / HTTP2.0\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 8,
		}

		r, err := ParseFromReader(reader)

		require.Nil(t, r)
		require.Error(t, err)
		require.ErrorContains(t, err, "malformed http version expected <HTTP-NAME>/<digit>.<digit>")
	})

	t.Run("Standard Headers", func(t *testing.T) {
		reader := &chunkReader{
			data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 3,
		}

		r, err := ParseFromReader(reader)

		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "localhost:42069", r.Headers["host"])
		assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
		assert.Equal(t, "*/*", r.Headers["accept"])
	})

	t.Run("With Duplicated Headers", func(t *testing.T) {
		// this unit test is very important
		// he is testing the whole macanism of parsing many data inside of the same chunk
		reader := &chunkReader{
			data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\n gremio: vamo0\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\ngremio: vamo1\r\ngremio: vamo2\r\ngremio: vamo3\r\n",
			numBytesPerRead: 1080,
		}

		r, err := ParseFromReader(reader)

		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "localhost:42069", r.Headers["host"])
		assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
		assert.Equal(t, "*/*", r.Headers["accept"])
		assert.Equal(t, "vamo0, vamo1, vamo2, vamo3", r.Headers["gremio"])
	})
}

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}

	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}

	n = copy(p, cr.data[cr.pos:endIndex])

	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}

	return n, nil
}

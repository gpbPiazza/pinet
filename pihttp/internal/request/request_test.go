package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFromReader(t *testing.T) {
	t.Run("Good request line without request target method Get", func(t *testing.T) {
		reader := &chunkReader{
			data: "GET / HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"User-Agent: curl/7.81.0\r\n" +
				"Accept: */*\r\n" +
				"\r\n",
			numBytesPerRead: 1,
		}

		r, err := ParseFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
		assert.Empty(t, r.Body)
	})

	t.Run("Good request with \n at the end of the body", func(t *testing.T) {
		reader := &chunkReader{
			data:            "GET /httpbin/stream/100 HTTP/1.1\r\nHost: localhost:42069\r\nConnection: close\r\n\r\n\n",
			numBytesPerRead: 18,
		}

		r, err := ParseFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/httpbin/stream/100", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
		assert.Empty(t, r.Body)
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
		assert.Empty(t, r.Body)
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
			data: "GET / HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				" gremio: vamo0\r\n" +
				"User-Agent: curl/7.81.0\r\n" +
				"Accept: */*\r\n" +
				"gremio: vamo1\r\n" +
				"gremio: vamo2\r\n" +
				"gremio: vamo3\r\n" +
				"\r\n",
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

	t.Run("Standard Body", func(t *testing.T) {
		reader := &chunkReader{
			data: "POST /submit HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"Content-Length: 13\r\n" +
				"\r\n" +
				"hello world!\n",
			numBytesPerRead: 3,
		}

		r, err := ParseFromReader(reader)

		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "hello world!\n", string(r.Body))
	})

	t.Run("request with json body", func(t *testing.T) {
		reader := &chunkReader{
			data: "POST /coffee HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"User-Agent: curl/8.6.0\r\n" +
				"Accept: */*\r\n" +
				"Content-Type: application/json\r\n" +
				"Content-Length: 22\r\n" +
				"\r\n" +
				"{\"flavor\":\"dark mode\"}",
			numBytesPerRead: 25,
		}

		r, err := ParseFromReader(reader)

		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "{\"flavor\":\"dark mode\"}", string(r.Body))
	})

	t.Run("Empty body reported zero length", func(t *testing.T) {
		reader := &chunkReader{
			data: "POST /submit HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"Content-Length: 0\r\n" +
				"\r\n",
			numBytesPerRead: 3,
		}

		r, err := ParseFromReader(reader)

		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Empty(t, r.Body)
	})

	t.Run("Empty body and NOT reported zero length", func(t *testing.T) {
		reader := &chunkReader{
			data: "POST /submit HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"\r\n",
			numBytesPerRead: 3,
		}

		r, err := ParseFromReader(reader)

		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Empty(t, r.Body)
	})

	t.Run("Body is bigger than reported content length", func(t *testing.T) {
		reader := &chunkReader{
			data: "POST /submit HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"Content-Length: 2\r\n" +
				"\r\n" +
				"partial content",
			numBytesPerRead: 3,
		}

		r, err := ParseFromReader(reader)

		require.Error(t, err)
		require.Nil(t, r)
		assert.ErrorContains(t, err, "error: content length informed is less than body length")
	})

	t.Run("content lenght is not an int", func(t *testing.T) {
		reader := &chunkReader{
			data: "POST /submit HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"Content-Length: gremio\r\n" +
				"\r\n" +
				"partial content",
			numBytesPerRead: 3,
		}

		r, err := ParseFromReader(reader)

		require.Error(t, err)
		require.Nil(t, r)
		assert.ErrorContains(t, err, "error: content length value is not an int")
	})

	t.Run("No content length reported and has body", func(t *testing.T) {
		reader := &chunkReader{
			data: "POST /submit HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"\r\n" +
				"partial content",
			numBytesPerRead: 3,
		}

		r, err := ParseFromReader(reader)

		require.Error(t, err)
		require.ErrorContains(t, err, "error: content length header is required")
		require.Nil(t, r)
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

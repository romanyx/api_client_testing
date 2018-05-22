package example

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	okResponse = `{
		"users": [
			{"id": 1, "name": "Roman"},
			{"id": 2, "name": "Dmitry"}
		]	
	}`
)

func TestClientGetUsers(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "key", r.Header.Get("Key"))
		assert.Equal(t, "secret", r.Header.Get("Secret"))
		w.Write([]byte(okResponse))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	cli := NewClient("key", "secret")
	cli.httpClient = httpClient

	users, err := cli.GetUsers()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(users))
}

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewTLSServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("tcp", s.URL[strings.LastIndex(s.URL, "/")+1:])
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return cli, func() {
		s.Close()
	}
}

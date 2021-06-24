// +build integration

package api

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)
	s, err := NewHTTPServer(
		url.URL{Scheme: "", Host: "localhost"},
		"githubsecret",
		"githubkey",
		"somerandomtext",
		"development",
	)
	assert.Nil(err)
	ts := httptest.NewServer(s.GetHandler())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal("<p><a href='/auth/github?provider=github'>Click to log in with github</a></p>", string(body))
}

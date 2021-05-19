// +build integration

package api

import (
	"net/url"
	"testing"
)

func TestInitHTTPRouter(t *testing.T) {
	myURL, e := url.Parse("http://localhost:3000")
	check(e)
	c := ConfigHTTPRouter{
		URLLoc: *myURL,
	}
	InitHTTPRouter(c)
}

package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)

	if err != nil {
		t.Fatal(err)
	}

	// response cookies are stored in this jar
	// and subsequent requests get those cookies
	ts.Client().Jar = jar

	// preventing the test server to automatically following redirects
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}

}

func (ts *testServer) get(t *testing.T, urlpath string) (int, http.Header, string) {

	rs, err := ts.Client().Get(ts.URL + urlpath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)

}

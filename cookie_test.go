package gocookie

import (
	"net/http"
	"testing"
	"time"
)

func TestGetCookies(t *testing.T) {
	cookies, err := GetCookies(Edge, WithDomainSuffix("apple.com"))
	if err != nil {
		t.Fatal(err)
	}

	for _, cookie := range cookies {
		t.Logf("%+v", cookie)
	}
}

type loggingTransport struct {
	transport http.RoundTripper
	cookies   map[string]string
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, cookie := range req.Cookies() {
		t.cookies[cookie.Name] = cookie.Value
	}
	return t.transport.RoundTrip(req)
}

func TestClientCookie(t *testing.T) {
	cookies := []*http.Cookie{
		{
			Name:   "test1",
			Value:  "nok",
			Domain: "example.com",
			Secure: true,
		},
		{
			Name:   "test2",
			Value:  "ok",
			Domain: "example.com",
			Secure: false,
		},
		{
			Name:   "test3",
			Value:  "ok",
			Domain: ".example.com",
			Secure: false,
		},
		{
			Name:   "test4",
			Value:  "nok",
			Domain: "www.example.com",
			Secure: false,
		},
		{
			Name:    "test5",
			Value:   "ok",
			Domain:  "example.com",
			Expires: time.Now().Add(time.Hour * 24),
		},
		{
			Name:    "test6",
			Value:   "nok",
			Domain:  "example.com",
			Expires: time.Now().Add(time.Hour * -24),
		},
	}

	jar, err := NewCookieJar(cookies)
	if err != nil {
		t.Fatal(err)
	}

	loggingTransport := &loggingTransport{
		transport: http.DefaultTransport,
		cookies:   make(map[string]string),
	}

	client := &http.Client{
		Jar:       jar,
		Transport: loggingTransport,
	}

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	for _, cookie := range cookies {
		if cookie.Value == "ok" {
			if _, ok := loggingTransport.cookies[cookie.Name]; !ok {
				t.Errorf("cookie %s is not set", cookie.Name)
			}
		} else {
			if _, ok := loggingTransport.cookies[cookie.Name]; ok {
				t.Errorf("cookie %s is set", cookie.Name)
			}
		}
	}
}

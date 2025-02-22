package gocookie

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func GetCookies(browserType BrowserType, opts ...Option) ([]*http.Cookie, error) {
	o := new(Options)
	for _, opt := range opts {
		opt(o)
	}
	browser, err := NewBrowser(browserType, o.CookiePath, o.SecretKey)
	if err != nil {
		return nil, err
	}
	return browser.GetCookies(o.DomainFilter())
}

func GetCookiesJar(browserType BrowserType, opts ...Option) (http.CookieJar, error) {
	cookies, err := GetCookies(browserType, opts...)
	if err != nil {
		return nil, err
	}
	return NewCookieJar(cookies)
}

func NewCookieJar(cookies []*http.Cookie) (http.CookieJar, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}

	for _, cookie := range cookies {
		domain := strings.TrimPrefix(cookie.Domain, ".")
		scheme := "http://"
		if cookie.Secure {
			scheme = "https://"
		}
		u, err := url.Parse(scheme + domain)
		if err != nil {
			continue
		}
		jar.SetCookies(u, []*http.Cookie{cookie})
	}

	return jar, nil
}

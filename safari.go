package gocookie

import (
	"net/http"
	"os"

	"github.com/cixtor/binarycookies"
)

type safari struct {
	cookiePath string
}

func newSafari(cookiePath string) Browser {
	return &safari{
		cookiePath: cookiePath,
	}
}

func (s *safari) GetName() string {
	return string(Safari)
}

func (s *safari) GetCookies(domainFilter domainFilter) ([]*http.Cookie, error) {
	f, err := os.Open(s.cookiePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cook := binarycookies.New(f)
	pages, err := cook.Decode()
	if err != nil {
		return nil, err
	}

	var cookies []*http.Cookie
	for _, page := range pages {
		for _, cookie := range page.Cookies {
			domain := string(cookie.Domain)
			if domainFilter != nil && !domainFilter(domain) {
				continue
			}

			// fmt.Println(string(cookie.Name), string(cookie.Value), domain, string(cookie.Path), cookie.Expires, cookie.Secure, cookie.HttpOnly)
			cookies = append(cookies, &http.Cookie{
				Name:     string(cookie.Name),
				Value:    string(cookie.Value),
				Domain:   domain,
				Path:     string(cookie.Path),
				Expires:  cookie.Expires,
				Secure:   cookie.Secure,
				HttpOnly: cookie.HttpOnly,
			})
		}
	}

	return cookies, nil
}

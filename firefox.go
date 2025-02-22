package gocookie

import (
	"database/sql"
	"net/http"
	"time"
)

const queryFirefoxCookie = `SELECT name, value, host, path, expiry, isSecure, isHttpOnly FROM moz_cookies`

type firefox struct {
	cookiePath string
}

func newFirefox(cookiePath string) Browser {
	return &firefox{
		cookiePath: cookiePath,
	}
}

func (f *firefox) GetName() string {
	return string(Firefox)
}

func (f *firefox) GetCookies(domainFilter domainFilter) ([]*http.Cookie, error) {
	cookiesDB, err := sql.Open("sqlite3", "file:"+f.cookiePath+"?mode=ro")
	if err != nil {
		return nil, err
	}
	defer cookiesDB.Close()

	rows, err := cookiesDB.Query(queryFirefoxCookie)
	if err != nil {
		return nil, err
	}

	var cookies []*http.Cookie
	for rows.Next() {
		var (
			name       string
			value      string
			host       string
			path       string
			expiry     int64
			isSecure   int
			isHTTPOnly int
		)
		if err = rows.Scan(&name, &value, &host, &path, &expiry, &isSecure, &isHTTPOnly); err != nil {
			return nil, err
		}

		if domainFilter != nil && !domainFilter(host) {
			continue
		}

		// fmt.Println(name, value, host, path, expiry, isSecure, isHTTPOnly)
		cookies = append(cookies, &http.Cookie{
			Name:     name,
			Value:    value,
			Domain:   host,
			Path:     path,
			Expires:  time.Unix(expiry, 0),
			Secure:   isSecure > 0,
			HttpOnly: isHTTPOnly > 0,
		})
	}
	return cookies, nil
}

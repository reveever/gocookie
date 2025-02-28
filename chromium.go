package gocookie

import (
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/pbkdf2"
	_ "modernc.org/sqlite"
)

const queryChromiumCookie = `SELECT name, encrypted_value, host_key, path, expires_utc, is_secure, is_httponly FROM cookies`

type chromium struct {
	name       string
	cookiePath string
	secretKey  []byte
}

func newChromium(name, cookiePath string, secretKey []byte) (Browser, error) {
	if cookiePath == "" || len(secretKey) == 0 {
		return nil, errors.New("empty cookie path or secret key")
	}
	return &chromium{
		name:       name,
		cookiePath: cookiePath,
		secretKey:  secretKey,
	}, nil
}

func (c *chromium) GetName() string {
	return c.name
}

func (c *chromium) GetCookies(domainFilter domainFilter) ([]*http.Cookie, error) {
	cookiesDB, err := sql.Open("sqlite", "file:"+c.cookiePath+"?mode=ro")
	if err != nil {
		return nil, err
	}
	defer cookiesDB.Close()

	rows, err := cookiesDB.Query(queryChromiumCookie)
	if err != nil {
		return nil, err
	}

	var (
		cookies []*http.Cookie
		aesKey  = pbkdf2.Key(c.secretKey, []byte("saltysalt"), 1003, 16, sha1.New)
		iv      = []byte{32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32}
	)

	for rows.Next() {
		var (
			name         string
			encryptValue []byte
			hostKey      string
			path         string
			expiresUTC   int64
			isSecure     int
			isHTTPOnly   int
		)
		if err = rows.Scan(&name, &encryptValue, &hostKey, &path, &expiresUTC, &isSecure, &isHTTPOnly); err != nil {
			return nil, err
		}

		if domainFilter != nil && !domainFilter(hostKey) {
			continue
		}

		value, err := c.decrypt(aesKey, iv, encryptValue)
		if err != nil {
			return nil, err
		}

		// fmt.Println(name, string(value), hostKey, path, expiresUTC, isSecure, isHTTPOnly)
		cookies = append(cookies, &http.Cookie{
			Name:     name,
			Value:    string(value),
			Domain:   hostKey,
			Path:     path,
			Expires:  ChromiumTimeToUnix(expiresUTC),
			Secure:   isSecure > 0,
			HttpOnly: isHTTPOnly > 0,
		})
	}
	return cookies, nil
}

func (c *chromium) decrypt(key, iv, encryptValue []byte) ([]byte, error) {
	if len(encryptValue) <= 3 {
		return nil, errors.New("encryptValue is too short")
	}
	value, err := DecryptAESCBC(key, iv, encryptValue[3:])
	if err != nil {
		return nil, err
	}
	if len(value) < 32 {
		return nil, fmt.Errorf("value is too short: %d", len(value))
	}
	return value[32:], nil
}

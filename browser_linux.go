//go:build linux

package gocookie

import (
	"errors"
	"os"
)

var homeDir, _ = os.UserHomeDir()

func GetCookieFilePath(browserType BrowserType) string {
	switch browserType {
	case Chrome:
		return homeDir + "/.config/google-chrome/Default/Cookies"
	case Edge:
		return homeDir + "/.config/microsoft-edge/Default/Cookies"
	case Firefox:
		return homeDir + "/.mozilla/firefox/cookies.sqlite"
	default:
		return ""
	}
}

func GetEdgeSecretKey() ([]byte, error) {
	return nil, errors.New("get edge secret key not implemented")
}

func GetChromeSecretKey() ([]byte, error) {
	return nil, errors.New("get chrome secret key not implemented")
}

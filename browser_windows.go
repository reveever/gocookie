//go:build windows

package gocookie

import (
	"errors"
	"os"
)

var homeDir, _ = os.UserHomeDir()

func GetCookieFilePath(browserType BrowserType) string {
	switch browserType {
	case Chrome:
		return homeDir + "/AppData/Local/Google/Chrome/User Data/Default/Cookies"
	case Edge:
		return homeDir + "/AppData/Local/Microsoft/Edge/User Data/Default/Network/Cookies"
	case Firefox:
		return homeDir + "/AppData/Roaming/Mozilla/Firefox/Profiles/cookies.sqlite"
	default:
		return ""
	}
}

func GetEdgeSecretKey(cookiePath string) ([]byte, error) {
	return nil, errors.New("get edge secret key not implemented")
}

func GetChromeSecretKey(cookiePath string) ([]byte, error) {
	return nil, errors.New("get chrome secret key not implemented")
}

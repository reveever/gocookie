//go:build darwin

package gocookie

import (
	"os"
	"strings"

	"github.com/keybase/go-keychain"
)

var homeDir, _ = os.UserHomeDir()

func GetCookieFilePath(browserType BrowserType) string {
	switch browserType {
	case Chrome:
		return homeDir + "/Library/Application Support/Google/Chrome/Default/Cookies"
	case Edge:
		return homeDir + "/Library/Application Support/Microsoft Edge/Default/Cookies"
	case Firefox:
		profiles, err := os.ReadDir(homeDir + "/Library/Application Support/Firefox/Profiles")
		if err != nil {
			return ""
		}
		for _, profile := range profiles {
			if strings.HasSuffix(profile.Name(), "default-release") {
				return homeDir + "/Library/Application Support/Firefox/Profiles/" + profile.Name() + "/cookies.sqlite"
			}
		}
		return ""
	case Safari:
		return homeDir + "/Library/Cookies/Cookies.binarycookies"
	default:
		return ""
	}
}

func GetEdgeSecretKey() ([]byte, error) {
	return keychain.GetGenericPassword(`Microsoft Edge Safe Storage`, `Microsoft Edge`, "", "")
}

func GetChromeSecretKey() ([]byte, error) {
	return keychain.GetGenericPassword(`Chrome Safe Storage`, `Chrome`, "", "")
}

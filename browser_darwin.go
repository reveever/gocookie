//go:build darwin

package gocookie

import (
	"os"
	"strings"
	"sync"

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

var (
	edgeSecretKey   []byte
	edgeSecretKeyMu sync.Mutex
)

func GetEdgeSecretKey(cookiePath string) ([]byte, error) {
	edgeSecretKeyMu.Lock()
	defer edgeSecretKeyMu.Unlock()

	if edgeSecretKey != nil {
		return edgeSecretKey, nil
	}

	key, err := keychain.GetGenericPassword(`Microsoft Edge Safe Storage`, `Microsoft Edge`, "", "")
	if err != nil {
		return nil, err
	}
	edgeSecretKey = key
	return key, nil
}

var (
	chromeSecretKey   []byte
	chromeSecretKeyMu sync.Mutex
)

func GetChromeSecretKey(cookiePath string) ([]byte, error) {
	chromeSecretKeyMu.Lock()
	defer chromeSecretKeyMu.Unlock()

	if chromeSecretKey != nil {
		return chromeSecretKey, nil
	}

	key, err := keychain.GetGenericPassword(`Chrome Safe Storage`, `Chrome`, "", "")
	if err != nil {
		return nil, err
	}
	chromeSecretKey = key
	return key, nil
}

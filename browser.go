package gocookie

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Browser interface {
	GetName() string
	GetCookies(domainFilter) ([]*http.Cookie, error)
}

type BrowserType string

const (
	Safari   BrowserType = "Safari"
	Chrome   BrowserType = "Chrome"
	Edge     BrowserType = "Edge"
	Chromium BrowserType = "Chromium"
	Firefox  BrowserType = "Firefox"
)

func ParseBrowserType(browserName string) (BrowserType, error) {
	name := strings.TrimSpace(strings.ToLower(browserName))

	if strings.Contains(name, "safari") || strings.Contains(name, "sf") || strings.Contains(name, "appl") {
		return Safari, nil
	}
	if strings.Contains(name, "chrome") || strings.Contains(name, "goo") {
		return Chrome, nil
	}
	if strings.Contains(name, "edge") || strings.Contains(name, "ms") {
		return Edge, nil
	}
	if strings.Contains(name, "chromium") {
		return Chromium, nil
	}
	if strings.Contains(name, "firefox") || strings.Contains(name, "ff") || strings.Contains(name, "moz") {
		return Firefox, nil
	}

	return "", fmt.Errorf("unsupported browser type: %s (supported: safari, chrome, edge, chromium, firefox)", browserName)
}

func NewBrowser(browserType BrowserType, cookiePath string, secretKey []byte) (Browser, error) {
	if cookiePath == "" {
		cookiePath = GetCookieFilePath(browserType)
	}
	if cookiePath == "" {
		return nil, errors.New("empty cookie path")
	}

	if _, err := os.Stat(cookiePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("cookie file not found: %s", cookiePath)
	}

	switch browserType {
	case Safari:
		return newSafari(cookiePath), nil

	case Chrome:
		if len(secretKey) == 0 {
			var err error
			if secretKey, err = GetChromeSecretKey(cookiePath); err != nil {
				return nil, err
			}
		}
		return newChromium(string(Chrome), cookiePath, secretKey)

	case Edge:
		if len(secretKey) == 0 {
			var err error
			if secretKey, err = GetEdgeSecretKey(cookiePath); err != nil {
				return nil, err
			}
		}
		return newChromium(string(Edge), cookiePath, secretKey)

	case Chromium:
		return newChromium(string(Chromium), cookiePath, secretKey)

	case Firefox:
		return newFirefox(cookiePath), nil

	default:
		return nil, errors.New("unsupported browser type")
	}
}

type Option func(*Options)

type Options struct {
	CookiePath     string
	SecretKey      []byte
	Domains        []string
	DomainSuffix   []string
	DomainContains []string
}

type domainFilter func(domain string) bool

func (o *Options) DomainFilter() domainFilter {
	if len(o.Domains) == 0 && len(o.DomainSuffix) == 0 && len(o.DomainContains) == 0 {
		return nil
	}

	return func(domain string) bool {
		if len(o.Domains) > 0 {
			for _, d := range o.Domains {
				if d == domain {
					return true
				}
			}
		}

		if len(o.DomainSuffix) > 0 {
			for _, suffix := range o.DomainSuffix {
				if strings.HasSuffix(domain, suffix) {
					return true
				}
			}
		}

		if len(o.DomainContains) > 0 {
			for _, contains := range o.DomainContains {
				if strings.Contains(domain, contains) {
					return true
				}
			}
		}

		return false
	}
}

func WithCookiePath(path string) Option {
	return func(o *Options) {
		o.CookiePath = path
	}
}

func WithSecretKey(secretKey []byte) Option {
	return func(o *Options) {
		o.SecretKey = secretKey
	}
}

func WithDomain(domains ...string) Option {
	return func(o *Options) {
		o.Domains = append(o.Domains, domains...)
	}
}

func WithDomainSuffix(suffixes ...string) Option {
	return func(o *Options) {
		o.DomainSuffix = append(o.DomainSuffix, suffixes...)
	}
}

func WithDomainContains(contains ...string) Option {
	return func(o *Options) {
		o.DomainContains = append(o.DomainContains, contains...)
	}
}

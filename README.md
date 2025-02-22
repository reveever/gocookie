# Go Cookie

[![Go Reference](https://pkg.go.dev/badge/github.com/reveever/gocookie.svg)](https://pkg.go.dev/github.com/reveever/gocookie) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go library for extracting cookies from various web browsers on macOS (Windows and Linux support coming soon).

## Features

* Extract cookies from multiple browsers on macOS:
  - Safari
  - Google Chrome
  - Microsoft Edge
  - Firefox
  - Other Chromium-based browsers
* Generate `http.CookieJar` directly for use with `http.Client`
* Filter cookies by domain suffix, name, or other criteria

## Example

```go
cookies, err := gocookie.GetCookies(gocookie.Edge, gocookie.WithDomainSuffix("example.com"))
if err != nil {
    log.Fatal(err)
}

for _, cookie := range cookies {
    fmt.Printf("%+v\n", cookie)
}
```

```go
jar, err := gocookie.GetCookiesJar(gocookie.Edge, gocookie.WithDomainSuffix("example.com"))
if err != nil {
    log.Fatal(err)
}

client := &http.Client{
    Jar: jar,
}
client.Get("https://example.com")
```

## Notes

### macOS Keychain Access

When reading cookies from Chromium-based browsers (Chrome, Edge, etc.) on macOS, if no `secretKey` is provided, the library will attempt to read the encryption key from the system keychain. This will trigger a system prompt asking for keychain access permission. To avoid this prompt, you can provide the secret key directly using the `WithSecretKey` option.

### Windows and Linux Support

Windows and Linux support is coming soon.

## Acknowledgments

This project was inspired by and builds upon the work of:

- [kooky](https://github.com/browserutils/kooky)
- [macCookies](https://github.com/kawakatz/macCookies)

Special thanks to their authors and contributors.

## License

[MIT](LICENSE)

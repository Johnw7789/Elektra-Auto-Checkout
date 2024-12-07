package login

import (
	"fmt"
	"log"
	"net/url"

	http "github.com/bogdanfinn/fhttp"
)

func (l *LoginClient) GetCookieStr() string {
	var cookies string

	url, _ := url.Parse("https://www.bestbuy.com")
	for _, c := range l.HttpClient.GetCookies(url) {
		cookies += c.Name + "=" + c.Value + "; "
	}

	return cookies
}

func (l *LoginClient) logMessage(module, msg string) {
	if !l.opts.Logging {
		log.Println(fmt.Sprintf("[%s] %s", module, msg))
	}
}

func (l *LoginClient) defaultHeaders(header http.Header) http.Header {
	header.Set("authority", "www.bestbuy.com")
	header.Set("content-type", "application/json")
	header.Set("accept", "*/*")
	header.Set("accept-language", "en-US,en;q=0.9")

	if l.opts.UserAgent != "" {
		header.Set("User-Agent", l.opts.UserAgent)
	} else {
		header.Set("User-Agent", DefaultUserAgent)
	}

	return header
}

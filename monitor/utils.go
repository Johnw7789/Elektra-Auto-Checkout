package monitor

import (
	"fmt"
	"log"

	http "github.com/bogdanfinn/fhttp"
)

func (m *MonitorClient) Cancel() {
	m.cancelChannel <- true
}

func (m *MonitorClient) logMessage(module, msg string) {
	if m.opts.Logging {
		log.Println(fmt.Sprintf("[SKU %s] [%s] %s", m.opts.Sku, module, msg))
	}
}

func (m *MonitorClient) defaultHeaders(header http.Header) http.Header {
	header.Set("accept-language", "en-US,en;q=0.9")

	if m.opts.UserAgent != "" {
		header.Set("User-Agent", m.opts.UserAgent)
	} else {
		header.Set("User-Agent", DefaultUserAgent)
	}

	return header
}

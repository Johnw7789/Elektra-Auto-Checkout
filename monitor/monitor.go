package monitor

import (
	"errors"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type MonitorClient struct {
	HttpClient   tls_client.HttpClient
	AlertChannel chan bool // * Alert channel for when the monitor finds stock

	cancelChannel chan bool

	opts MonitorOpts
}

type MonitorOpts struct {
	Sku       string
	Delay     int
	Proxy     string
	UserAgent string
	Logging   bool
}

// * NewMonitorClient intializes a new http client and returns a new monitor instance
func NewMonitorClient(opts MonitorOpts) (*MonitorClient, error) {
	if opts.Sku == "" {
		return nil, errors.New("sku is required")
	}

	jar := tls_client.NewCookieJar()

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithCookieJar(jar),
		tls_client.WithClientProfile(profiles.Chrome_124),
	}

	if opts.Proxy != "" {
		options = append(options, tls_client.WithProxyUrl(opts.Proxy))
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	return &MonitorClient{
		HttpClient:   client,
		AlertChannel: make(chan bool),
		opts:         opts,
	}, nil
}

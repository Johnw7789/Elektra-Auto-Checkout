package login

import (
	"errors"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type LoginClient struct {
	HttpClient tls_client.HttpClient

	opts LoginOpts
}

type LoginOpts struct {
	Email         string
	Password      string
	EmailPassword string
	Proxy         string
	UserAgent     string
	Logging       bool
}

// * NewLoginClient intializes a new http client and returns a new login instance
func NewLoginClient(opts LoginOpts) (*LoginClient, error) {
	if opts.Email == "" || opts.Password == "" {
		return nil, errors.New("email and password are required")
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

	return &LoginClient{
		HttpClient: client,
		opts:       opts,
	}, nil
}

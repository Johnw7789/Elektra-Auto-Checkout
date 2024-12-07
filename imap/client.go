package imap

import (
	"errors"
	"strings"

	"github.com/emersion/go-imap/client"
)

type EmailClient struct {
	Username string
	Password string
	Client   *client.Client
}

func NewEmailClient(username, password string) (*EmailClient, error) {
	var err error
	ec := &EmailClient{
		Username: username,
		Password: password,
	}

	err = ec.imapLogin()
	if err != nil {
		return nil, errors.New("failed to login with IMAP")
	}

	return ec, nil
}

func (e *EmailClient) imapLogin() error {
	// Connect to server
	var err error

	var domain string = "imap.mail.me.com"

	if strings.Contains(e.Username, "@gmail.com") {
		domain = "imap.gmail.com"
	} else if strings.Contains(e.Username, "@outlook.com") {
		domain = "imap-mail.outlook.com"
	} else if !strings.Contains(e.Username, "@icloud.com") {
		return errors.New("Invalid email domain")
	}

	e.Client, err = client.DialTLS(domain+":993", nil)
	if err != nil {
		return err
	}

	err = e.Client.Login(e.Username, e.Password)
	if err != nil {
		return err
	}

	return nil
}

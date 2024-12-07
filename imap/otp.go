package imap

import (
	"errors"
	"io"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/emersion/go-imap"
)

// * Scans the last 10 emails for a OTP, for up to a minute
func (e *EmailClient) FetchOtp() (string, error) {
	var err error
	var msg string

	for i := 0; i < 20; i++ {
		msg, err = e.searchLast10Messages()
		if err != nil {
			continue
		}

		re := regexp.MustCompile(`<span style="font-size:18px; font-weight:bold;">(.*?)</span>`)
		otp := re.FindString(msg)

		otp = strings.TrimLeft(otp, "<span style=\"font-size:18px; font-weight:bold;\">")
		otp = strings.TrimRight(otp, "</span>")

		if len(otp) == 6 {
			return otp, nil
		}

		time.Sleep(3 * time.Second)
	}

	return "", errors.New("failed to fetch otp")
}

func (e *EmailClient) searchLast10Messages() (string, error) {
	mbox, err := e.Client.Select("Inbox", false)
	if err != nil {
		return "", err
	}

	if mbox.Messages == 0 {
		return "", errors.New("no messages in inbox")
	}

	// Get the latest message
	seqset := new(imap.SeqSet)

	size := 1

	if mbox.Messages > 10 {
		seqset.AddRange(mbox.Messages-10, mbox.Messages)
		size = 10
	} else {
		seqset.AddRange(1, mbox.Messages)
		size = int(mbox.Messages)
	}

	messages := make(chan *imap.Message, size)

	// Get the whole message body
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	// messages := make(chan *imap.Message, mbox.Messages)

	done := make(chan error, 1)
	go func() {
		done <- e.Client.Fetch(seqset, items, messages)
	}()

	for msg := range messages {
		if msg == nil {
			continue
		}

		r := msg.GetBody(section)
		if r == nil {
			continue
		}

		m, err := mail.ReadMessage(r)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(m.Body)
		if err != nil {
			continue
		}

		// timeReceived, err := m.Header.Date()
		// if err != nil {
		// 	continue
		// }

		if strings.Contains(strings.ToLower(string(body)), "bestbuy") && strings.Contains(strings.ToLower(string(body)), "verification code") {
			return string(body), nil
		}
	}

	return "", nil
}

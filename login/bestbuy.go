package login

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/emersion/go-imap/client"
	"github.com/obito/cclient"
	utls "github.com/refraction-networking/utls"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type BestBuyEncryptionData struct {
  	EncryptedEmail    string
  	EncryptedAgent    string
  	EncryptedActivity string
}

type BestBuyLoginData struct {
  	VerificationCodeFieldName string
  	EncryptedPasswordField    string
  	EncryptedAlpha            string
  	EmailField                string
  	Salmon                    string
  	Token                     string
}

type BestBuyLogin struct {
  	Username      string
	Password      string
	Email         string
	Phone         string
  	UserAgent     string
  	Proxy         string       
	LoginData     BestBuyLoginData
  	EncryptedData BestBuyEncryptionData
}


func ImapLogin(email string, password string) *client.Client {
	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Login
	log.Println("Logging in")
	err = c.Login(email, password)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Logged in")

	return c
}


func GetNewestMessage(c * client.Client) string {
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	// Get the last message
	if mbox.Messages == 0 {
		return ""
	} else {
		seqset := new(imap.SeqSet)
		seqset.AddRange(mbox.Messages, mbox.Messages)

		// Get the whole message body
		section := &imap.BodySectionName{}
		items := []imap.FetchItem{section.FetchItem()}

		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)
		go func() {
			done <- c.Fetch(seqset, items, messages)
		}()

		msg := <-messages

		r := msg.GetBody(section)
		if r == nil {
			log.Fatal("Server didn't return a message body")
		}

		if err := <-done; err != nil {
			log.Fatal(err)
		}

		m, err := mail.ReadMessage(r)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(m.Body)
		if err != nil {
			log.Fatal(err)
		}

		return string(body)
	}
}




func (login *BestBuyLogin) parse(str string, start string, end string) []byte {
	var match []byte
	index := strings.Index(str, start)

	if index == -1 {
		return match
	}

	index += len(start)

	for {
		char := str[index]

		if strings.HasPrefix(str[index:index+len(match)], end) {
			break
		}

		match = append(match, char)
		index++
	}

	return match
}

func (login *BestBuyLogin) reverse(s string) string {
	rs := []rune(s)
	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		rs[i], rs[j] = rs[j], rs[i]
	}
	return string(rs)
}



func (login *BestBuyLogin) encrypt(s string, publicKey string, keyId string) string {
	block, _ := pem.Decode([]byte(publicKey))
	if block.Type != "PUBLIC KEY" {
		log.Fatal("error decoding public key from pem")
	}
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal("error parsing key")
	}
	var ok bool
	var pubkey *rsa.PublicKey
	if pubkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		log.Fatal("unable to parse public key")
	}
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha1.New(), rng, pubkey, []byte(s), nil)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Join([]string{"1", keyId, base64.StdEncoding.EncodeToString(ciphertext)}, ":")
}

func (login *BestBuyLogin) BestbuyLoginSession() {
  
  	login.EncryptedData.EncryptedEmail = encrypt(email, emailPublicKey, emailKeyId)
	login.EncryptedData.EncryptedAgent = encrypt(fmt.Sprintf("{\"user-agent\": \"%s\"}", userAgent), activityPublicKey, activityKeyId)
	login.EncryptedData.EncryptedActivity = encrypt(fmt.Sprintf("{mouseMoved\":true,\"keyboardUsed\":true,\"fieldReceivedInput\":true,\"fieldReceivedFocus\":true,\"timestamp\":\"%s\",\"email\":\"%s\"}", time.Now().UTC().Format("2006-01-02T15:04:05-0700"), email), activityPublicKey, activityKeyId)
  
}

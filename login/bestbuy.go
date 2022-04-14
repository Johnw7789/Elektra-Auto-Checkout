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
	ua "github.com/wux1an/fake-useragent"
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

func (login *BestBuyLogin) getPublicKey(client * http.Client, url string) (string, string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "www.bestbuy.com")
	req.Header.Set("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"97\", \"Chromium\";v=\"97\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", userAgent)
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "*/*")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	//req.Header.Set("referer", "https://www.bestbuy.com/identity/signin?token=tid%3A505572da-7422-11ec-80b8-12491479a2d1")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)

	publicKey := gjson.Get(string(bodyText), "publicKey").String()
	keyId := gjson.Get(string(bodyText), "keyId").String()

	return publicKey, keyId
}

func (login *BestBuyLogin) bestbuyEncrypt(s string, publicKey string, keyId string) string {
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

func (login *BestBuyLogin) BestbuyLoginSession() (bool, bool, error) {
	client, err := elektra.CreateClient(login.Proxy)
	if err != nil {
		log.Println("Error creating client")
	}
	
	if login.UserAgent == "" {
		login.UserAgent = ua.RandomType(ua.Desktop)
	}

	login.scrapeLoginData(client)
	
  	emailPublicKey, emailKeyId := getPublicKey(client, "https://www.bestbuy.com/api/csiservice/v2/key/cia-email")
	activityPublicKey, activityKeyId := getPublicKey(client, "https://www.bestbuy.com/api/csiservice/v2/key/cia-user-activity")
	
  	login.EncryptedData.EncryptedEmail = bestbuyEncrypt(email, emailPublicKey, emailKeyId)
	login.EncryptedData.EncryptedAgent = bestbuyEncrypt(fmt.Sprintf("{\"user-agent\": \"%s\"}", userAgent), activityPublicKey, activityKeyId)
	login.EncryptedData.EncryptedActivity = bestbuyEncrypt(fmt.Sprintf("{mouseMoved\":true,\"keyboardUsed\":true,\"fieldReceivedInput\":true,\"fieldReceivedFocus\":true,\"timestamp\":\"%s\",\"email\":\"%s\"}", time.Now().UTC().Format("2006-01-02T15:04:05-0700"), login.Email), activityPublicKey, activityKeyId)
}

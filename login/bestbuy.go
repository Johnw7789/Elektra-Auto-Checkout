
package login

import (
	"encoding/base64"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/ffeathers/Elektra-Auto-Checkout/elektra"
	"github.com/tidwall/gjson"
	ua "github.com/wux1an/fake-useragent"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/url"
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
	Id            		  string
	Username      		  string
	Email         		  string
	Password      		  string
	Phone         		  string
	GmailPassword 		  string
	UserAgent     		  string
	Proxy         		  string
	Cookies 		  string
	Active         		  bool
	LoggingDisabled       	  bool
	BestBuyLoginData      	  BestBuyLoginData
	BestBuyEncryptionData 	  BestBuyEncryptionData
}

func (login *BestBuyLogin) logMessage(msg string) {
	if !login.LoggingDisabled {
		log.Println(fmt.Sprintf("[Login %s] [BestBuy] %s", monitor.Id, msg))
	}
}

func (login *BestBuyLogin) Cancel() {
	login.Active = false
	login.logMessage("Login canceled")
	//add exit code
}

func ImapLogin(email string, password string) *client.Client {
	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Login to gmail
	err = c.Login(email, password)
	if err != nil {
		log.Fatal(err)
	}

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

func (login *BestBuyLogin) getPublicKey(client * http.Client, url string) (string, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("authority", "www.bestbuy.com")
	req.Header.Set("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"97\", \"Chromium\";v=\"97\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", login.UserAgent)
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "*/*")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	publicKey := gjson.Get(string(bodyText), "publicKey").String()
	keyId := gjson.Get(string(bodyText), "keyId").String()

	return publicKey, keyId, nil
}

func (login *BestBuyLogin) scrapeLoginData(client * http.Client) error {
	req, err := http.NewRequest("GET", "https://www.bestbuy.com/identity/global/signin", nil)
	if err != nil {
		return err
	}
	req.Header.Set("authority", "www.bestbuy.com")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", login.UserAgent)
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"97\", \"Chromium\";v=\"97\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}


	responseString := string(response)

	initDataBytes := login.parse(responseString, "var initData = ", "; </script>")
	initData := string(initDataBytes)

	login.BestBuyLoginData.VerificationCodeFieldName = gjson.Get(initData, "verificationCodeFieldName").String()

	passwordArray := gjson.Get(initData, "codeList")
	for _, passwordField := range passwordArray.Array() {
		decodedString, _ := base64.URLEncoding.DecodeString(passwordField.String())
		if strings.Contains(string(decodedString) , "_X_") {
			login.BestBuyLoginData.EncryptedPasswordField = passwordField.String()
			break
		}
	}

	alphaArray := gjson.Get(initData, "alpha")
	for _, alpha := range alphaArray.Array() {
		decodedString, _ := base64.URLEncoding.DecodeString(login.reverse(alpha.String()))
		if strings.Contains(string(decodedString) , "_A_") {
			login.BestBuyLoginData.EncryptedAlpha = alpha.String()
			break
		}
	}

	login.BestBuyLoginData.EmailField = gjson.Get(initData, "emailFieldName").String()
	login.BestBuyLoginData.Salmon = gjson.Get(initData, "Salmon").String()
	login.BestBuyLoginData.Token = gjson.Get(initData, "token").String()

	return nil
}


func (login *BestBuyLogin) bestbuyLogin(client * http.Client, loginJson string) (string, error) {
	var data = strings.NewReader(loginJson)
	req, err := http.NewRequest("POST", "https://www.bestbuy.com/identity/authenticate", data)
	if err != nil {
		return "", err
	}
	req.Header.Set("authority", "www.bestbuy.com")
	req.Header.Set("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"97\", \"Chromium\";v=\"97\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", login.UserAgent)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("origin", "https://www.bestbuy.com")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(response), nil
}


func (login *BestBuyLogin) submitAuthCode(client *http.Client, authCode string, flowOptions string) (string, error) {
	var data = strings.NewReader(fmt.Sprintf(`{"token":"%s","isResetFlow":false,"challengeType":"2","smsDigits":"","flowOptions":"%s","%s":"%s","%s":"%s"}`, login.BestBuyLoginData.Token, flowOptions, login.BestBuyLoginData.EmailField, login.Email, login.BestBuyLoginData.VerificationCodeFieldName, authCode))
	req, err := http.NewRequest("POST", "https://www.bestbuy.com/identity/unlock", data)
	if err != nil {
		return "", err
	}
	req.Header.Set("authority", "www.bestbuy.com")
	req.Header.Set("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"97\", \"Chromium\";v=\"97\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", login.UserAgent)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("origin", "https://www.bestbuy.com")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	authResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(authResp), nil
}


func (login *BestBuyLogin) verifyWithEmail(client *http.Client, flowOptions string, challengeType string) error {
	var data = strings.NewReader(fmt.Sprintf(`{"token":"%s","recoveryOptionType":"email","email":"%s","smsDigits":"","isResetFlow":false,"challengeType":"%s","flowOptions":"%s"}`, login.BestBuyLoginData.Token, login.Email, challengeType, flowOptions))
	req, err := http.NewRequest("POST", "https://www.bestbuy.com/identity/account/recovery/code", data)
	if err != nil {
		return err
	}
	req.Header.Set("authority", "www.bestbuy.com")
	req.Header.Set("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"97\", \"Chromium\";v=\"97\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", login.UserAgent)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("origin", "https://www.bestbuy.com")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}


func (login *BestBuyLogin) getAuthCode(c * client.Client) string {
	for {
		message := GetNewestMessage(c)
		if strings.Contains(message, "BestBuy") && strings.Contains(message, "Verification Code") {
			re := regexp.MustCompile(`<span style="font-size:18px; font-weight:bold;">(.*?)</span>`)
			authCode := re.FindString(message)

			authCode = strings.TrimLeft(authCode, "<span style=\"font-size:18px; font-weight:bold;\">")
			authCode = strings.TrimRight(authCode, "</span>")

			return authCode
		}
		time.Sleep(time.Second * 1)
	}
}

func (login *BestBuyLogin) getCookieStr(client * http.Client) string {
	var cookies string

	url, _ := url.Parse("https://www.bestbuy.com")
	for _, c := range client.Jar.Cookies(url) {
		cookies += c.Name + "=" + c.Value + "; "
	}

	return cookies
}

// BestbuyLoginSession login successful (bool), is banned (bool), error
func (login *BestBuyLogin) BestbuyLoginSession() (bool, bool, error) {
	client, err := elektra.CreateClient(login.Proxy)
	if err != nil {
		login.logMessage("Error creating client")
		return false, false, err
	}

	if login.UserAgent == "" {
		login.UserAgent = ua.RandomType(ua.Desktop)
	}

	// very repetitive, checks to see if login.Cancel() was called
	if !login.Active {return false, false, nil}
	err = login.scrapeLoginData(client)
	if err != nil {
		login.logMessage("Error scraping login data")
		return false, false, err
	}

	if !login.Active {return false, false, nil}
	emailPublicKey, emailKeyId, err := login.getPublicKey(client, "https://www.bestbuy.com/api/csiservice/v2/key/cia-email")
	if err != nil {
		login.logMessage("Error fetching email public key")
		return false, false, err
	}

	if !login.Active {return false, false, nil}
	activityPublicKey, activityKeyId, err := login.getPublicKey(client, "https://www.bestbuy.com/api/csiservice/v2/key/cia-user-activity")
	if err != nil {
		login.logMessage("Error fetching activity public key")
		return false, false, err
	}

	if !login.Active {return false, false, nil}
	login.BestBuyEncryptionData.EncryptedEmail, err = elektra.BestbuyEncrypt(login.Email, emailPublicKey, emailKeyId)
	if err != nil {
		login.logMessage("Error encrypting email")
		return false, false, err
	}

	if !login.Active {return false, false, nil}
	login.BestBuyEncryptionData.EncryptedAgent, err = elektra.BestbuyEncrypt(fmt.Sprintf("{\"user-agent\": \"%s\"}", login.UserAgent), activityPublicKey, activityKeyId)
	if err != nil {
		login.logMessage("Error encrypting useragent")
		return false, false, err
	}

	if !login.Active {return false, false, nil}
	login.BestBuyEncryptionData.EncryptedActivity, err = elektra.BestbuyEncrypt(fmt.Sprintf("{mouseMoved\":true,\"keyboardUsed\":true,\"fieldReceivedInput\":true,\"fieldReceivedFocus\":true,\"timestamp\":\"%s\",\"email\":\"%s\"}", time.Now().UTC().Format("2006-01-02T15:04:05-0700"), login.Email), activityPublicKey, activityKeyId)
	if err != nil {
		login.logMessage("Error encrypting activity")
		return false, false, err
	}

	loginJson := fmt.Sprintf("{\"token\":\"%s\",\"activity\":\"%s\",\"loginMethod\":\"UID_PASSWORD\",\"flowOptions\":\"0000000000000000\",\"alpha\":\"%s\",\"Salmon\":\"%s\",\"encryptedEmail\":\"%s\",\"%s\":\"%s\",\"info\":\"%s\",\"%s\":\"%s\"}", login.BestBuyLoginData.Token, login.BestBuyEncryptionData.EncryptedActivity, login.BestBuyLoginData.EncryptedAlpha, login.BestBuyLoginData.Salmon, login.BestBuyEncryptionData.EncryptedEmail, login.BestBuyLoginData.EncryptedPasswordField, login.Password, login.BestBuyEncryptionData.EncryptedAgent, login.BestBuyLoginData.EmailField, login.Email)
	if !login.Active {return false, false, nil}
	loginResp, err := login.bestbuyLogin(client, loginJson)
	if err != nil {
		login.logMessage("Error submitting login")
		return false, false, err
	}


	status := gjson.Get(loginResp, "status").String()
	if status == "success" {
		login.logMessage("Successfully logged in")
		login.Cookies = login.getCookieStr(client)
		return true, false, nil
	} else if status == "stepUpRequired" {
		login.logMessage("Code verification required")

		if login.GmailPassword != "" {
			if !login.Active {return false, false, nil}
			c := ImapLogin(login.Email, login.GmailPassword)


			flowOptions := gjson.Get(loginResp, "flowOptions").String()
			challengeType := gjson.Get(loginResp, "challengeType").String()

			if !login.Active {return false, false, nil}
			login.verifyWithEmail(client, flowOptions, challengeType)

			login.logMessage("Fetching auth code")
			if !login.Active {return false, false, nil}
			authCode := login.getAuthCode(c)

			if authCode != "" {
				if !login.Active {return false, false, nil}
				login.logMessage(fmt.Sprintf("Fetched authentication code: %s", authCode))
				authResp,err := login.submitAuthCode(client, authCode, flowOptions)
				if err != nil {
					return false, false, err
				}

				if !login.Active {return false, false, nil}
				if strings.Contains(authResp, "temporaryAccessToken") {
					login.logMessage("Account flagged, new password required")
				} else {
					login.Cookies = login.getCookieStr(client)
					return true, false, nil
				}
			} else {
				login.logMessage("Failed to fetch authentication code")
				return false, false, nil
			}

			//incomplete flow
		} else {
			return false, false, nil
		}


	} else if status == "failure" {
		login.logMessage("Account flagged")
		return false, true, nil
	} else {
		login.logMessage("Unknown error")
		return false, false, nil
	}

	return false, false, nil
}

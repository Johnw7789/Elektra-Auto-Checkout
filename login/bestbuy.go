package login

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"

	"github.com/Johnw7789/Elektra-Auto-Checkout/imap"
	"github.com/Johnw7789/Elektra-Auto-Checkout/shr"
	"github.com/tidwall/gjson"
)

// BestbuyLoginSession login successful (bool), is banned (bool), error
func (l *LoginClient) BestbuyTask() (bool, error) {
	l.logMessage("Bestbuy", "Getting login data")
	ld, err := l.getLoginData()
	if err != nil {
		return false, err
	}

	l.logMessage("Bestbuy", "Getting encryption data")
	loginJson, err := l.handleEncryption(ld)
	if err != nil {
		return false, err
	}

	l.logMessage("Bestbuy", "Logging in")
	success, flowOptions, challengeType, err := l.getBestbuyLogin(loginJson)
	if err != nil {
		return false, err
	}

	if success {
		l.logMessage("Bestbuy", "Successfully logged in")
		return true, nil
	}

	if errors.Is(err, errors.New("step up required")) {
		l.logMessage("Bestbuy", "Code verification required")

		if l.opts.EmailPassword == "" {
			return false, errors.New("email password required")
		}

		err = l.verifyWithEmail(ld, flowOptions, challengeType)
		if err != nil {
			return false, err
		}

		ec, err := imap.NewEmailClient(l.opts.Email, l.opts.EmailPassword)
		if err != nil {
			return false, err
		}
		defer ec.Client.Logout()

		l.logMessage("Bestbuy", "Fetching auth code")

		time.Sleep(time.Second * 5)

		authCode, err := ec.FetchOtp()
		if err != nil {
			return false, err
		}

		l.logMessage("Bestbuy", "Submitting otp code: "+authCode)
		return l.getSubmitAuthCode(ld, authCode, flowOptions)
	}

	return false, err
}

func (l *LoginClient) handleEncryption(ld BestbuyLoginData) (string, error) {
	// * Get public keys
	emailPublicKey, emailKeyId, err := l.getPublicKey(BestbuyEmailPKUrl)
	if err != nil {
		return "", err
	}

	activityPublicKey, activityKeyId, err := l.getPublicKey(BestbuyActPKUrl)
	if err != nil {
		return "", err
	}

	// * Do encryption
	var encData BestbuyEncryptionData
	encData.EncryptedEmail, err = BestbuyEncrypt(l.opts.Email, emailPublicKey, emailKeyId)
	if err != nil {
		return "", err
	}

	encData.EncryptedAgent, err = BestbuyEncrypt(fmt.Sprintf("{\"user-agent\": \"%s\"}", l.opts.UserAgent), activityPublicKey, activityKeyId)
	if err != nil {
		return "", err
	}

	encData.EncryptedActivity, err = BestbuyEncrypt(fmt.Sprintf("{mouseMoved\":true,\"keyboardUsed\":true,\"fieldReceivedInput\":true,\"fieldReceivedFocus\":true,\"timestamp\":\"%s\",\"email\":\"%s\"}", time.Now().UTC().Format("2006-01-02T15:04:05-0700"), l.opts.Email), activityPublicKey, activityKeyId)
	if err != nil {
		l.logMessage("Bestbuy", "Error encrypting activity")
		return "", err
	}

	// * Create login json with encrypted data
	loginJson := fmt.Sprintf("{\"token\":\"%s\",\"activity\":\"%s\",\"loginMethod\":\"UID_PASSWORD\",\"flowOptions\":\"0000000000000000\",\"alpha\":\"%s\",\"Salmon\":\"%s\",\"encryptedEmail\":\"%s\",\"%s\":\"%s\",\"info\":\"%s\",\"%s\":\"%s\",\"recaptchaData\": \"Error: recaptcha is not enabled.\"}", ld.Token, encData.EncryptedActivity, ld.EncryptedAlpha, ld.Salmon, encData.EncryptedEmail, ld.EncryptedPasswordField, l.opts.Password, encData.EncryptedAgent, ld.EmailField, l.opts.Email)

	return loginJson, nil
}

func (l *LoginClient) getPublicKey(url string) (string, string, error) {
	body, err := l.reqPublicKey(url)
	if err != nil {
		return "", "", err
	}

	publicKey := gjson.Get(body, "publicKey").String()
	keyId := gjson.Get(body, "keyId").String()

	if publicKey == "" || keyId == "" {
		return "", "", errors.New("failed to get public key")
	}

	return publicKey, keyId, nil
}

func (l *LoginClient) getLoginData() (BestbuyLoginData, error) {
	body, err := l.reqLoginData()
	if err != nil {
		return BestbuyLoginData{}, err
	}

	initData := shr.ParseV3(body, "var initData = ", "; </script>")

	var ld BestbuyLoginData
	ld.VerificationCodeFieldName = gjson.Get(initData, "verificationCodeFieldName").String()

	passwordArray := gjson.Get(initData, "codeList")
	for _, passwordField := range passwordArray.Array() {
		decodedString, _ := base64.URLEncoding.DecodeString(passwordField.String())
		if strings.Contains(string(decodedString), "_X_") {
			ld.EncryptedPasswordField = passwordField.String()
			break
		}
	}

	alphaArray := gjson.Get(initData, "alpha")
	for _, alpha := range alphaArray.Array() {
		decodedString, _ := base64.URLEncoding.DecodeString(shr.Reverse(alpha.String()))
		if strings.Contains(string(decodedString), "_A_") {
			ld.EncryptedAlpha = alpha.String()
			break
		}
	}

	ld.EmailField = gjson.Get(initData, "emailFieldName").String()
	ld.Salmon = gjson.Get(initData, "Salmon").String()
	ld.Token = gjson.Get(initData, "token").String()

	return ld, nil
}

func (l *LoginClient) getBestbuyLogin(loginJson string) (bool, string, string, error) {
	body, err := l.reqBestbuyLogin(loginJson)
	if err != nil {
		return false, "", "", err
	}

	status := gjson.Get(body, "status").String()
	switch status {
	case "success":
		return true, "", "", nil
	case "stepUpRequired":
		flowOptions := gjson.Get(body, "flowOptions").String()
		challengeType := gjson.Get(body, "challengeType").String()

		return false, flowOptions, challengeType, errors.New("step up required")
	case "failure":
		return false, "", "", errors.New("login failed")
	default:
		return false, "", "", errors.New("failed login, status: " + status)
	}
}

func (l *LoginClient) getSubmitAuthCode(ld BestbuyLoginData, authCode, flowOptions string) (bool, error) {
	body, err := l.reqSubmitAuthCode(ld, authCode, flowOptions)
	if err != nil {
		return false, err
	}

	if strings.Contains(body, "temporaryAccessToken") {
		l.logMessage("Bestbuy", "Account flagged, new password required")
		return false, nil
	}

	if strings.Contains(body, "success") {
		return true, nil
	}

	return false, errors.New("failed to submit auth code")
}

func (l *LoginClient) verifyWithEmail(ld BestbuyLoginData, flowOptions string, challengeType string) error {
	var data = strings.NewReader(fmt.Sprintf(`{"token":"%s","recoveryOptionType":"email","email":"%s","smsDigits":"","isResetFlow":false,"challengeType":"%s","flowOptions":"%s"}`, ld.Token, l.opts.Email, challengeType, flowOptions))

	req, err := http.NewRequest("POST", BestbuyVerifyEmailUrl, data)
	if err != nil {
		return err
	}

	req.Header = l.defaultHeaders(req.Header.Clone())

	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (l *LoginClient) reqPublicKey(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header = l.defaultHeaders(req.Header.Clone())

	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (l *LoginClient) reqLoginData() (string, error) {
	req, err := http.NewRequest("GET", BestbuySigninUrl, nil)
	if err != nil {
		return "", err
	}

	req.Header = l.defaultHeaders(req.Header.Clone())

	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (l *LoginClient) reqBestbuyLogin(loginJson string) (string, error) {
	var data = strings.NewReader(loginJson)

	req, err := http.NewRequest("POST", BestbuyAuthUrl, data)
	if err != nil {
		return "", err
	}

	req.Header = l.defaultHeaders(req.Header.Clone())

	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (l *LoginClient) reqSubmitAuthCode(ld BestbuyLoginData, authCode, flowOptions string) (string, error) {
	var data = strings.NewReader(fmt.Sprintf(`{"token":"%s","isResetFlow":false,"challengeType":"2","smsDigits":"","flowOptions":"%s","%s":"%s","%s":"%s"}`, ld.Token, flowOptions, ld.EmailField, l.opts.Email, ld.VerificationCodeFieldName, authCode))

	req, err := http.NewRequest("POST", BestbuySubmitAuthUrl, data)
	if err != nil {
		return "", err
	}

	req.Header = l.defaultHeaders(req.Header.Clone())

	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

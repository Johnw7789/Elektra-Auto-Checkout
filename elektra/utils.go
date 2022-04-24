package elektra

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/obito/cclient"
	utls "github.com/refraction-networking/utls"
	"github.com/xxtea/xxtea-go/xxtea"
	"log"
	"net/http"
	"strings"
)

func Parse(value string, a string, b string) string {
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func CreateClient(proxy string) (*http.Client, error) {
	if proxy != "" {
		proxyUrl := "http://" + proxy //Only works with IP authenticated proxies atm (IP:Port), not yet with User:Pass:IP:Port proxies

		client, err := cclient.NewClient(utls.HelloChrome_Auto, true, proxyUrl) //Create an http client with a Firefox TLS fingerprint, set automatic storage of cookies to true, and use a proxy
		if err != nil {
			return nil, err
		}

		log.Println("Created Client")

		return &client, nil
	} else {
		client, err := cclient.NewClient(utls.HelloChrome_Auto, true) //Create an http client with a Firefox TLS fingerprint, set automatic storage of cookies to true
		if err != nil {
			return nil, err
		}

		log.Println("Created Client")

		return &client, nil
	}
}

func XxteaEncrypt(data string) string {
	key := "a\x03\x8fp4\x18\x97\x99:\xeb\xe7\x8b\x85\x97$4"
	encryptedData := xxtea.EncryptString(data, key)

	return encryptedData
}

func BestbuyEncrypt(s string, publicKey string, keyId string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	parsedKey, _ := x509.ParsePKIXPublicKey(block.Bytes)

	var pubkey *rsa.PublicKey
	pubkey, _ = parsedKey.(*rsa.PublicKey)

	rng := rand.Reader
	ciphertext, _ := rsa.EncryptOAEP(sha1.New(), rng, pubkey, []byte(s), nil)

	return strings.Join([]string{"1", keyId, base64.StdEncoding.EncodeToString(ciphertext)}, ":"), nil
}

package login

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"

	xxtea "github.com/xxtea/xxtea-go/xxtea"
)

func AmazonEncrypt(data []byte) []byte {
	key := []byte("a\x03\x8fp4\x18\x97\x99:\xeb\xe7\x8b\x85\x97$4")
	encryptedData := xxtea.Encrypt(data, key)

	return encryptedData
}

func BestbuyEncrypt(s string, publicKey string, keyId string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	var pubkey *rsa.PublicKey
	pubkey, _ = parsedKey.(*rsa.PublicKey)

	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha1.New(), rng, pubkey, []byte(s), nil)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{"1", keyId, base64.StdEncoding.EncodeToString(ciphertext)}, ":"), nil
}

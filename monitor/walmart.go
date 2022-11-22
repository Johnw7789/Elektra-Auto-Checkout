package monitor 

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)



type HashList struct {
	Id          string `json:"WM_CONSUMER.ID"`
	TimeStamp   string `json:"WM_CONSUMER.INTIMESTAMP"`
	KeyVer      string `json:"WM_SEC.KEY_VERSION"`
}

func getRequestHashSum(headersString string) []byte {
	var data = []byte(headersString)
	msgHash := sha256.New()
	_, err := msgHash.Write(data)
	if err != nil {
		panic(err)
	}
	return msgHash.Sum(nil)
}

func signPkcs1(sortedHashString string) string {
  hashSum := getRequestHashSum(sortedHashString)
	block, err := pem.Decode([]byte(privateKey))
  if err != nil {
    panic(err) 
  }

	parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
  
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, requestHashSum)
	if err != nil {
		panic(err)
	}
  
  signatureEnc := base64.StdEncoding.EncodeToString(signature)
	return signatureEnc
}

func WalmartCheckStockV1() {
}

func WalmartMonitorTask() {
	productId := "931920073"
	apiUrl := "https://developer.api.walmart.com/api-proxy/service/affil/product/v2/items/" + productId

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		panic(err)
	}

	timestamp := strconv.Itoa(int(time.Now().UnixMilli()))
	var hashList HashList = HashList{
		Id:          consumerId,
		TimeStamp:   timestamp,
		KeyVer:      "2",
	}

	sortedHashString := fmt.Sprintf("%s\n%s\n%s\n", hashList.Id, hashList.TimeStamp, hashList.KeyVer)

	signatureEnc := signPkcs1(sortedHashString)

	req.Header.Set("WM_CONSUMER.ID", hashList.Id)
	req.Header.Set("WM_CONSUMER.INTIMESTAMP", hashList.TimeStamp)
	req.Header.Set("WM_SEC.KEY_VERSION", hashList.KeyVer)
	req.Header.Set("WM_SEC.AUTH_SIGNATURE", signatureEnc)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	bodyStr, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bodyStr))
}

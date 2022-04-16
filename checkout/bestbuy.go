package checkout

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
)


type BestBuyOrderData struct {
	OrderId    string
	ItemId     string
}

type BestBuyCheckout struct {
	UserAgent  		 string
	CartId     		 string
	Cookies    		 string
	Proxy      		 string
	RetryDelay 		 int
	MaxRetries 		 int
	Sku        		 string
	OrderNum   		 string
	BestBuyOrderData BestBuyOrderData
}





//successful checkout start bool, ip banned bool, error
func (checkout *BestBuyCheckout) initiateCheckout(client *http.Client) (bool, bool, error) {
	req, err := http.NewRequest("GET", "https://www.bestbuy.com/checkout/r/fufillment", nil)
	if err != nil {
		return false, false, err
	}
	req.Header.Set("authority", "www.bestbuy.com")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"92\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"92\"")
	req.Header.Set("accept", "application/json")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", checkout.UserAgent)
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("origin", "https://www.bestbuy.com")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("accept-language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return false, false, err
	}
	defer resp.Body.Close()
	atcResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, false, err
	}

	s := strings.Split(string(atcResp), "var orderData = ")
	json := strings.Split(s[1], ";")[0]

	checkout.BestBuyOrderData.OrderId = gjson.Get(json, "id").String()
	checkout.BestBuyOrderData.ItemId = gjson.Get(json, "items.0.id").String()

	return true, false, nil
}



//successful cart bool, ip banned bool, error
func (checkout *BestBuyCheckout) addToCart(client *http.Client) (bool, bool, error) {
	var data = strings.NewReader(fmt.Sprintf(`{"items":[{"skuId":"%s"}]}`, checkout.Sku))
	req, err := http.NewRequest("POST", "https://www.bestbuy.com/cart/api/v1/addToCart", data)
	if err != nil {
		return false, false, err
	}
	req.Header.Set("authority", "www.bestbuy.com")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"92\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"92\"")
	req.Header.Set("accept", "application/json")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", checkout.UserAgent)
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("origin", "https://www.bestbuy.com")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("accept-language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return false, false, err
	}
	defer resp.Body.Close()
	atcResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, false, err
	}

	if gjson.Get(string(atcResp), "cartCount").Int() > 0 {
		return true, false, nil
	}

	return false, false, nil
}


func (checkout *BestBuyCheckout) fetchTasKey(client *http.Client) (string, string, error) {
	req, err := http.NewRequest("GET", "https://www.bestbuy.com/api/csiservice/v2/key/tas", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("authority", "www.bestbuy.com")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"92\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"92\"")
	req.Header.Set("accept", "application/json")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", checkout.UserAgent)
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("origin", "https://www.bestbuy.com")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("accept-language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	atcResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	publicKey := gjson.Get(string(atcResp), "publicKey").String()
	keyId := gjson.Get(string(atcResp), "keyId").String()

	return publicKey, keyId, nil
}


func (checkout *BestBuyCheckout) BestBuyCheckoutTask () {

}

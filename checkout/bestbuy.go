package checkout

import (
	"errors"
	"fmt"
	"github.com/ffeathers/Elektra-Auto-Checkout/elektra"
	"github.com/tidwall/gjson"
	ua "github.com/wux1an/fake-useragent"
	"io/ioutil"
	"net/http"
	"strings"
)


type BestBuyOrderData struct {
	OrderId    string
	ItemId     string
	PaymentId  string
}

type BestBuyShipping struct {
	Phone 			string
	Email 			string
	FirstName 		string
	MiddleInitial 	string
	LastName 		string
	Addressln1 		string
	Addressln2 		string
	State 			string
	City 			string
	ZipCode 		string
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
	BestBuyShipping  BestBuyShipping
	BestBuyOrderData BestBuyOrderData
}





func (checkout *BestBuyCheckout) refreshPayment(client *http.Client) (bool, bool, error) {
	var data = strings.NewReader(`{}`)
	req, err := http.NewRequest("POST", "https://www.bestbuy.com/checkout/orders/"+checkout.BestBuyOrderData.OrderId+"/paymentMethods/refreshPayment", data)
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

	return true, false, nil
}


func (checkout *BestBuyCheckout) submitFufillment(client *http.Client) (bool, bool, error) {
	var data = strings.NewReader(fmt.Sprintf(`"{"items":[{"id":"%s","type":"DEFAULT","selectedFulfillment":{"shipping":{"address":{"firstName":"%s","middleInitial":"%s","lastName":"%s","street":"%s","street2":"%s","state":"%s","city":"%s","zipcode":"%s","saveToProfile":false,"country":"US","id":"","dayPhoneNumber":"%s","setAsPrimaryShipping":true,"useAddressAsBilling":false}}},"giftMessageSelected":false}]}`, checkout.BestBuyOrderData.ItemId, checkout.BestBuyShipping.FirstName, checkout.BestBuyShipping.MiddleInitial, checkout.BestBuyShipping.LastName, checkout.BestBuyShipping.Addressln1, checkout.BestBuyShipping.Addressln2, checkout.BestBuyShipping.State, checkout.BestBuyShipping.City, checkout.BestBuyShipping.ZipCode, checkout.BestBuyShipping.Phone))
	req, err := http.NewRequest("PATCH", "https://www.bestbuy.com/checkout/orders/"+checkout.BestBuyOrderData.OrderId, data)
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
	fufillResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, false, err
	}


	checkout.BestBuyOrderData.PaymentId = gjson.Get(string(fufillResp), "payment.id").String()

	return true, false, nil
}

func (checkout *BestBuyCheckout) selectFufillment(client *http.Client) (bool, bool, error) {
	var data = strings.NewReader(fmt.Sprintf(`[{"id":"%s","selectedFulfillment":{"shipping":{}}}]`, checkout.BestBuyOrderData.ItemId))
	req, err := http.NewRequest("PATCH", "https://www.bestbuy.com/checkout/orders/"+checkout.BestBuyOrderData.OrderId+"/items", data)
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

	return true, false, nil
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

func (checkout *BestBuyCheckout) calculateQueueTime(queueId string) int {
	//placeholder
	return 0
}

//successful cart bool, ip banned bool, error
func (checkout *BestBuyCheckout) addToCart(client *http.Client) (bool, bool, int, bool, error) {
	var data = strings.NewReader(fmt.Sprintf(`{"items":[{"skuId":"%s"}]}`, checkout.Sku))
	req, err := http.NewRequest("POST", "https://www.bestbuy.com/cart/api/v1/addToCart", data)
	if err != nil {
		return false, false, 0, false, err
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
		return false, false, 0, false, err
	}
	defer resp.Body.Close()
	atcResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, false, 0, false, err
	}

	if gjson.Get(string(atcResp), "cartCount").Int() > 0 {
		return true, false, 0, false, nil
	} else if strings.Contains(string(atcResp), "CONSTRAINED_ITEM") {
		queueTime := checkout.calculateQueueTime("")
		return true, true, queueTime, false, nil
	}

	return true, false, 0, false, nil
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


func (checkout *BestBuyCheckout) BestBuyCheckoutTask () (bool, bool, error) {
	var queueEnforced bool
	//var queueTime int

	client, err := elektra.CreateClient(checkout.Proxy)
	if err != nil {
		return false, false, err
	}

	if checkout.UserAgent == "" {
		checkout.UserAgent = ua.RandomType(ua.Desktop) //If checkoutData.UserAgent is empty, set it to a randomly generated user agent
	}

	_, queueEnforced, _, _, err = checkout.addToCart(client)
	if err != nil {
		return false, false, err
	}

	if queueEnforced {
		//placeholder, do nothing since func calculateQueueTime is not complete

		/*
		time.Sleep(time.Second * time.Duration(queueTime))

		_, _, _, _, err = checkout.addToCart(client)
		if err != nil {
			return false, false, err
		}
		*/

		return false, false, errors.New("can not yet solve bestbuy queue")
	}

	_, _, err = checkout.initiateCheckout(client)
	if err != nil {
		return false, false, err
	}

	_, _, err = checkout.selectFufillment(client)
	if err != nil {
		return false, false, err
	}

	_, _, err = checkout.submitFufillment(client)
	if err != nil {
		return false, false, err
	}

	_, _, err = checkout.refreshPayment(client)
	if err != nil {
		return false, false, err
	}


	return false, false, nil
}

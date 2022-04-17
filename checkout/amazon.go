package checkout

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/ffeathers/Elektra-Auto-Checkout/elektra"
	ua "github.com/wux1an/fake-useragent"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)


type AmazonCheckout struct {
	UserAgent  string
	SessionId  string
	Cookies    string
	Proxy      string
	UseProxies bool
	RetryDelay int
	MaxRetries int
	Sku        string
	OfferId    string
	OrderNum   string
}


var turboHeaders = []string{
	"Accept: */*",
	"Content-Type: application/x-www-form-urlencoded",
	"x-amz-support-custom-signin: 1",
	"x-requested-with: XMLHttpRequest",
	"accept-language: en-US,en;q=0.9",
	"origin: https://www.amazon.com",
	"referer: https://www.amazon.com",
}


func (checkout *AmazonCheckout) amazonPlaceOrder(client *http.Client, purchaseId string, csrfToken string) (bool, string, bool, error) {
	var data = strings.NewReader(fmt.Sprintf(`x-amz-checkout-csrf-token=%s&ref_=chk_spc_placeOrder&referrer=spc&pid=%s&pipelineType=turbo&clientId=retailwebsite&temporaryAddToCart=1&hostPage=detail&weblab=RCX_CHECKOUT_TURBO_DESKTOP_PRIME_87783&isClientTimeBased=1`, checkout.SessionId, purchaseId))
	req, err := http.NewRequest("POST", "https://www.amazon.com/checkout/spc/place-order?ref_=chk_spc_placeOrder&_srcRID=&clientId=retailwebsite&pipelineType=turbo&cachebuster=&pid="+purchaseId, data)
	if err != nil {
		return false, "", false, err
	}

	req.Header.Set("x-amz-checkout-entry-referer-url", "https://www.amazon.com/gp/product/"+checkout.Sku)
	req.Header.Set("anti-csrftoken-a2z", csrfToken)
	req.Header.Set("User-Agent", checkout.UserAgent)
	req.Header.Set("Referer", "https://www.amazon.com/checkout/spc?pid="+purchaseId+"&pipelineType=turbo&clientId=retailwebsite&temporaryAddToCart=1&hostPage=detail&weblab=RCX_CHECKOUT_TURBO_DESKTOP_PRIME_87783")
	req.Header.Add("cookie", checkout.Cookies)
	resp, err := client.Do(req)
	if err != nil {
		return false, "", false, err
	}

	defer resp.Body.Close()

	for key, value := range resp.Header {
		if strings.Contains(key, "thankyou") || strings.Contains(value[0], "thankyou") {
			return true, "", false, nil
		}
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "", false, err
	}

	if strings.Contains(string(bodyText), "thankyou") {
		return true, "", false, nil
	}

	return false, "", false, nil
}

func (checkout *AmazonCheckout) amazonAddToCart(client *http.Client) (bool, string, string, bool, error) {
	postData := fmt.Sprintf(`isAsync=1&asin.1=%s&offerListing.1=%s&quantity.1=1`, checkout.Sku, checkout.OfferId)

	var data = strings.NewReader(postData)
	req, err := http.NewRequest("POST", "https://www.amazon.com/checkout/turbo-initiate?ref_=dp_start-bbf_1_glance_buyNow_2-1&referrer=detail&pipelineType=turbo&clientId=retailwebsite&weblab=RCX_CHECKOUT_TURBO_DESKTOP_PRIME_87783&temporaryAddToCart=1&asin.1=", data)
	if err != nil {
		return false, "", "", false, err
	}

	req.Header.Set("x-amz-checkout-entry-referer-url", "https://www.amazon.com/dp/"+checkout.Sku)
	req.Header.Set("x-amz-turbo-checkout-dp-url", "https://www.amazon.com/dp/"+checkout.Sku)
	req.Header.Set("x-amz-checkout-csrf-token", checkout.SessionId)
	req.Header.Set("user-agent", checkout.UserAgent)
	req.Header.Set("cookie", checkout.Cookies)
	resp, err := client.Do(req)
	if err != nil {
		return false, "", "", false, err
	}

	defer resp.Body.Close()

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "", "", false, err
	}

	if strings.Contains(string(bodyText), "Place your order") {
		doc := soup.HTMLParse(string(bodyText))
		purchaseId := elektra.Parse(string(bodyText), "currentPurchaseId\":\"", "\",\"pipelineType\"")
		csrfToken := doc.Find("input", "name", "anti-csrftoken-a2z").Attrs()["value"]
		return true, purchaseId, csrfToken, false, nil
	} else {
		return false, "", "", false, nil
	}
}

func (checkout *AmazonCheckout) amazonPrepareCart(client *http.Client) error {
	var data = strings.NewReader(`isAsync=1&addressID=`)
	req, err := http.NewRequest("POST", "https://www.amazon.com/checkout/turbo-initiate?ref_=chk_detail_eligibility_1-0&referrer=detail&pipelineType=turbo&clientId=retailwebsite&weblab=RCX_CHECKOUT_TURBO_DESKTOP_NONPRIME_87784&checkEligibilityOnly=true&temporaryAddToCart=1", data)
	if err != nil {
		return err
	}

	req.Header.Set("x-amz-checkout-entry-referer-url", "https://www.amazon.com/dp/"+checkout.Sku)
	req.Header.Set("x-amz-turbo-checkout-dp-url", "https://www.amazon.com/dp/"+checkout.Sku)
	req.Header.Set("x-amz-checkout-csrf-token", checkout.SessionId)
	req.Header.Set("user-agent", checkout.UserAgent)
	req.Header.Set("cookie", checkout.Cookies)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (checkout *AmazonCheckout) AmazonCheckoutTask() (bool, bool, error) {
	client, err := elektra.CreateClient(checkout.Proxy)
	if err != nil {
		return false, false, err
	}

	if checkout.UserAgent == "" {
		checkout.UserAgent = ua.RandomType(ua.Desktop) //If checkoutData.UserAgent is empty, set it to a randomly generated user agent
	}

	for retries := 0; retries < checkout.MaxRetries; retries++ {
		log.Println("Preparing Cart")
		checkout.amazonPrepareCart(client)
		log.Println("Adding to Cart")
		cartSuccess, purchaseId, csrfToken, isBanned, err := checkout.amazonAddToCart(client)
		if err != nil {
			return false, false, err
		} else if isBanned {
			return false, false, nil
		}

		if cartSuccess {
			log.Println("Submitting Order")
			success, orderNum, isBanned, err := checkout.amazonPlaceOrder(client, purchaseId, csrfToken) //Todo: add ability to fetch order number, currently returns empty string
			if err != nil {
				return false, false, err
			} else if isBanned {
				return false, false, nil
			}

			if success {
				checkout.OrderNum = orderNum
				return true, false, nil
			}
		}
		time.Sleep(time.Second * time.Duration(checkout.RetryDelay))
	}
	return false, false, nil
}

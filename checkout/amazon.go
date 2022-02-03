package checkout

import (
	"github.com/obito/cclient"
	utls "github.com/refraction-networking/utls"
)

var turboHeaders = []string{
	"Accept: */*", 
	"Content-Type: application/x-www-form-urlencoded",
	"x-amz-support-custom-signin: 1",
  	"x-requested-with: XMLHttpRequest",
  	"accept-language: en-US,en;q=0.9",
  	"origin: https://www.amazon.com",
  	"referer: https://www.amazon.com",
}



func Parse(value string, a string, b string) {
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

func amazonPlaceOrder() {
	var data = strings.NewReader(`x-amz-checkout-csrf-token=` + accountSessionId + `&ref_=chk_spc_placeOrder&referrer=spc&pid=` + purchaseId + `&pipelineType=turbo&clientId=retailwebsite&temporaryAddToCart=1&hostPage=detail&weblab=RCX_CHECKOUT_TURBO_DESKTOP_PRIME_87783&isClientTimeBased=1`)
	req, err := http.NewRequest("POST", "https://www.amazon.com/checkout/spc/place-order?ref_=chk_spc_placeOrder&_srcRID=&clientId=retailwebsite&pipelineType=turbo&cachebuster=&pid=" + purchaseId, data)
	if err != nil {
		return false
	}
	
	req.Header.Set("x-amz-checkout-entry-referer-url", "https://" + domain + ".amazon.com/gp/product/" + productId + "/ref=ewc_pr_img_1?smid=AZ5LJ56P0QUDV&psc=1")
	req.Header.Set("anti-csrftoken-a2z", csrfToken)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")
	req.Header.Set("Referer", "https://www.amazon.com/checkout/spc?pid=" + purchaseId + "&pipelineType=turbo&clientId=retailwebsite&temporaryAddToCart=1&hostPage=detail&weblab=RCX_CHECKOUT_TURBO_DESKTOP_PRIME_87783")
	req.Header.Add("cookie", cookies)
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	
	defer resp.Body.Close()
	
	for key, value := range resp.Header {
		if strings.Contains(key, "thankyou") || strings.Contains(value[0], "thankyou") {
			return true
		}
	}
	
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	
	if strings.Contains(string(bodyText), "thankyou") {
		return true
	}
	
	return false
}

func amazonAddToCart(client *http.Client) (bool, string, string) {
	postData := fmt.Sprintf(`isAsync=1&asin.1=%s&offerListing.1=%s&quantity.1=1`, productId, offerId)
  
	var data = strings.NewReader(postData)
	req, err := http.NewRequest("POST", "https://www.amazon.com/checkout/turbo-initiate?ref_=dp_start-bbf_1_glance_buyNow_2-1&referrer=detail&pipelineType=turbo&clientId=retailwebsite&weblab=RCX_CHECKOUT_TURBO_DESKTOP_PRIME_87783&temporaryAddToCart=1&asin.1=", data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("x-amz-checkout-entry-referer-url", "https://" + domain + ".amazon.com/dp/" + productId)
	req.Header.Set("x-amz-turbo-checkout-dp-url", "https://" + domain + ".amazon.com/dp/" + productId)
  	req.Header.Set("x-amz-checkout-csrf-token", accountSessionId)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")
	req.Header.Set("cookie", cookies)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	
	defer resp.Body.Close()
	
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	
	if strings.Contains(string(bodyText), "Place your order") {
		doc := soup.HTMLParse(string(bodyText))
		purchaseId := Parse(string(bodyText), "currentPurchaseId\":\"", "\",\"pipelineType\"")
		csrfToken := doc.Find("input", "name", "anti-csrftoken-a2z").Attrs()["value"]
		return true, purchaseId, csrfToken
	} else {
		return false, "", ""
	}
}

func amazonPrepareCart(checkoutData *AmazonCheckoutData) {
	var data = strings.NewReader(`isAsync=1&addressID=`)
	req, err := http.NewRequest("POST", "https://www.amazon.com/checkout/turbo-initiate?ref_=chk_detail_eligibility_1-0&referrer=detail&pipelineType=turbo&clientId=retailwebsite&weblab=RCX_CHECKOUT_TURBO_DESKTOP_NONPRIME_87784&checkEligibilityOnly=true&temporaryAddToCart=1", data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("x-amz-checkout-entry-referer-url", "https://www.amazon.com/dp/" + checkoutData.Sku)
	req.Header.Set("x-amz-turbo-checkout-dp-url", "https://www.amazon.com/dp/" + checkoutData.Sku)
	req.Header.Set("x-amz-checkout-csrf-token", checkoutData.SessionId)
	req.Header.Set("user-agent", checkoutData.UserAgent)
	req.Header.Set("cookie", checkoutData.Cookies)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}


func AmazonCheckoutTask() {
	var client *http.Client
	if checkoutData.UseProxies {
		rand.Seed(time.Now().Unix())
		proxy := "http://" + checkoutData.Proxies[rand.Intn(len(checkoutData.Proxies))]
		
		client, err = cclient.NewClient(utls.HelloFirefox_Auto, true, proxy) //Create an http client with a Firefox TLS fingerprint, set automatic storage of cookies to true, and use a proxy
		if err != nil {
			log.Fatal(err)
		}
	} else {
		client, err = cclient.NewClient(utls.HelloFirefox_Auto, true) //Create an http client with a Firefox TLS fingerprint, set automatic storage of cookies to true
		if err != nil {
			log.Fatal(err)
		}
	}

	
	for retries := 0, retries < checkoutData.MaxRetries, retries++ {
		amazonPrepareCart(&client, &checkoutData)
		cartSuccess, purchaseId, csrfToken := amazonAddToCart(&client, checkoutData)
		
		if cartSuccess {
			success := amazonPlaceOrder(&client, &checkoutData)
			if success {
				return true
			}
		}
		time.Sleep(time.Second * time.Interval(checkoutData.RetryDelay))
	}
	return false
}


package monitor

var turboHeaders = []string{
	"Accept: */*", 
	"Content-Type: application/x-www-form-urlencoded",
	"x-amz-support-custom-signin: 1",
  	"x-requested-with: XMLHttpRequest",
  	"accept-language: en-US,en;q=0.9",
  	"origin: https://www.amazon.com",
  	"referer: https://www.amazon.com",
}

func checkStock() (bool, bool) {
	acceptheader := "application/vnd.com.amazon.api+json; type=\"cart.add-items/v1\""
	contentheader := "application/vnd.com.amazon.api+json; type=\"cart.add-items.request/v1\""
	
	var data = strings.NewReader(`{"items":[{"asin":"` + productId + `","offerListingId":"` + offerId + `","quantity":1}]}`)
	req, err := http.NewRequest("POST", "https://data.amazon.com/api/marketplaces/ATVPDKIKX0DER/cart/carts/retail/items", data)
	if err != nil {
		log.Fatal(err)
	}
	
	req.Header.Set("x-api-csrf-token", apitoken)
	req.Header.Set("Content-Type", contentheader)
	req.Header.Set("Accept", acceptheader)
	req.Header.Set("User-Agent", "Bestbuy-mApp/202104201730 CFNetwork/1209 Darwin/20.2.0")
	
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	
	if resp.StatusCode == 200 {
		return true, false //In stock
	} else if resp.StatusCode == 404 {
		return false, true //Out of stock, but an api token refresh is required 
	} 
	
	return false, false //Out of stock and an api token refresh is not required
}

func getApiToken() {
  	url := "https://www.amazon.com/gp/aw/d/B00M382RJO" //One of many Amazon product pages that contains an embedded api token
  
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	req.Header.Set("Cookie", "session-id=" + sessionid)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	apitoken := Parse(string(body), "\"csrfToken\":\"", "\",\"baseAsin\"")
	return apitoken
}

func getSessionId(client *http.Client) {
	url := "https://www.amazon.com/gp/aws/cart/add-res.html?Quantity.1=1&OfferListingId.1="

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()


	body, _ := ioutil.ReadAll(resp.Body)
	sessionid := (Parse(string(body), "ue_sid='", "',\nue_mid='"))
	return string(sessionid)
}

func Amazon() {
	
}

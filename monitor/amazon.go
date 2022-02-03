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

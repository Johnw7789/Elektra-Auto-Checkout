package checkout

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/google/uuid"
	http "github.com/saucesteals/fhttp"
	"github.com/saucesteals/mimic"
	"golang.design/x/clipboard"
	"github.com/ffeathers/Elektra-Auto-Checkout/elektra"
)

var (
	latestVersion = mimic.MustGetLatestVersion(mimic.PlatformWindows)
)

type AmazonCheckout struct {
	Id              	string
	UserAgent  		string
	SessionId 		string
	Cookies    		string
	Proxy      		string
	RetryDelay 		int
	MaxRetries 		int
	Sku        		string
	OfferId    		string
	OrderNum   		string
	LoggingDisabled 	bool
	Active     	    	bool
	Mimicer         	mimic.ClientSpec
}


func (checkout *AmazonCheckout) logMessage(msg string) {
	if !checkout.LoggingDisabled {
		log.Println(fmt.Sprintf("[Checkout %s] [Amazon] %s", checkout.Id, msg))
	}
}

func (checkout *AmazonCheckout) Cancel() {
	checkout.Active = false
	checkout.logMessage(fmt.Sprintf("[Checkout %s] [Amazon] Checkout Canceled", checkout.Id))
	//add exit code
}

func cookieHeader(rawCookies string) []*http.Cookie {
	header := http.Header{}
	header.Add("Cookie", rawCookies)
	req := http.Request{Header: header}
	return req.Cookies()
}

func (checkout *AmazonCheckout) AmazonPlaceOrderV1(client *http.Client, purchaseId string, csrfToken string) (bool, string, bool, error) {
	var data = strings.NewReader(fmt.Sprintf(`x-amz-checkout-csrf-token=%s&ref_=chk_spc_placeOrder&referrer=spc&pid=%s&pipelineType=turbo&clientId=retailwebsite&hostPage=detail&isClientTimeBased=1`, checkout.SessionId, purchaseId))
	req, err := http.NewRequest("POST", "https://www.amazon.com/checkout/spc/place-order?pipelineType=turbo", data)
	if err != nil {
		return false, "", false, err
	}

	req.Header = http.Header{
		"accept":                           {"*/*"},
		"anti-csrftoken-a2z":               {csrfToken},
		"accept-encoding":                  {"gzip, deflate, br"},
		"accept-language":                  {"en,en_US;q=0.9"},
		"content-length":                   {"184"},
		"content-type":                     {"application/x-www-form-urlencoded"},
		"cookie":                           {checkout.Cookies},
		"device-memory":                    {"8"},
		"downlink":                         {"3.9"},
		"dpr":                              {"1"},
		"ect":                              {"4g"},
		"rtt":                              {"50"},
		"sec-ch-device-memory":             {"8"},
		"sec-ch-dpr":                       {"1"},
		"sec-ch-ua":                        {checkout.Mimicer.ClientHintUA()},
		"sec-ch-ua-mobile":                 {"?0"},
		"sec-ch-ua-platform":               {`"Windows"`},
		"sec-ch-viewport-width":            {`"988"`},
		"sec-fetch-dest":                   {"empty"},
		"sec-fetch-mode":                   {"cors"},
		"sec-fetch-site":                   {"same-origin"},
		"user-agent":                       {checkout.UserAgent},
		"viewport-width":                   {"988"},
		"x-amz-checkout-entry-referer-url": {"https://www.amazon.com/dp/" + checkout.Sku},
		"x-requested-with":                 {"XMLHttpRequest"},

		http.HeaderOrderKey: {
			"accept", "anti-csrftoken-a2z", "accept-encoding", "accept-language",
			"content-length", "content-type", "cookie",
			"device-memory", "downlink", "dpr",
			"ect", "rtt", "sec-ch-device-memory",
			"sec-ch-dpr", "sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-platform",
			"sec-ch-viewport-width", "sec-fetch-dest", "sec-fetch-mode",
			"sec-fetch-site", "user-agent", "viewport-width", "x-amz-checkout-entry-referer-url",
			"x-amz-support-custom-signin", "x-amz-turbo-checkout-dp-url",
			"x-requested-with",
		},
		http.PHeaderOrderKey: checkout.Mimicer.PseudoHeaderOrder(),
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, "", false, err
	}

	defer resp.Body.Close()

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "", false, err
	}

	_ = clipboard.Init()

	clipboard.Write(clipboard.FmtText, []byte(string(bodyText)))

	for key, value := range resp.Header {
		if strings.Contains(key, "thankyou") || strings.Contains(value[0], "thankyou") {
			return true, "", false, nil
		}
	}

	if strings.Contains(string(bodyText), "thankyou") || strings.Contains(string(bodyText), "please check to confirm that the order was placed") {
		return true, "", false, nil
	}

	return false, "", false, nil
}

func (checkout *AmazonCheckout) AmazonAddToCartV1(client *http.Client) (bool, string, string, bool, error) {
	postData := fmt.Sprintf(`isAsync=1&asin.1=%s&quantity.1=1`, checkout.Sku)

	var data = strings.NewReader(postData)
	req, err := http.NewRequest("POST", "https://www.amazon.com/checkout/turbo-initiate?pipelineType=turbo", data)
	if err != nil {
		return false, "", "", false, err
	}

	req.Header = http.Header{
		"accept":                           {"*/*"},
		"accept-encoding":                  {"gzip, deflate, br"},
		"accept-language":                  {"en,en_US;q=0.9"},
		"content-length":                   {"40"},
		"content-type":                     {"application/x-www-form-urlencoded"},
		"cookie":                           {checkout.Cookies},
		"device-memory":                    {"8"},
		"downlink":                         {"3.9"},
		"dpr":                              {"1"},
		"ect":                              {"4g"},
		"rtt":                              {"50"},
		"sec-ch-device-memory":             {"8"},
		"sec-ch-dpr":                       {"1"},
		"sec-ch-ua":                        {checkout.Mimicer.ClientHintUA()},
		"sec-ch-ua-mobile":                 {"?0"},
		"sec-ch-ua-platform":               {`"Windows"`},
		"sec-ch-viewport-width":            {`"988"`},
		"sec-fetch-dest":                   {"empty"},
		"sec-fetch-mode":                   {"cors"},
		"sec-fetch-site":                   {"same-origin"},
		"user-agent":                       {checkout.UserAgent},
		"viewport-width":                   {"988"},
		"x-amz-checkout-csrf-token":        {checkout.SessionId},
		"x-amz-checkout-entry-referer-url": {"https://www.amazon.com/dp/" + checkout.Sku},
		"x-amz-support-custom-signin":      {"1"},
		"x-amz-turbo-checkout-dp-url":      {"https://www.amazon.com/dp/" + checkout.Sku},
		"x-requested-with":                 {"XMLHttpRequest"},

		http.HeaderOrderKey: {
			"accept", "accept-encoding", "accept-language",
			"content-length", "content-type", "cookie",
			"device-memory", "downlink", "dpr",
			"ect", "rtt", "sec-ch-device-memory",
			"sec-ch-dpr", "sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-platform",
			"sec-ch-viewport-width", "sec-fetch-dest", "sec-fetch-mode",
			"sec-fetch-site", "user-agent", "viewport-width",
			"x-amz-checkout-csrf-token", "x-amz-checkout-entry-referer-url",
			"x-amz-support-custom-signin", "x-amz-turbo-checkout-dp-url",
			"x-requested-with",
		},
		http.PHeaderOrderKey: checkout.Mimicer.PseudoHeaderOrder(),
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, "", "", false, err
	}

	defer resp.Body.Close()

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "", "", false, err
	}

	log.Println(string(bodyText))

	if strings.Contains(string(bodyText), "Place your order") {
		log.Println("Generated Checkout")
		doc := soup.HTMLParse(string(bodyText))
		purchaseId := elektra.Parse(string(bodyText), "currentPurchaseId\":\"", "\",\"pipelineType\"")
		csrfToken := doc.Find("input", "name", "anti-csrftoken-a2z").Attrs()["value"]
		return true, purchaseId, csrfToken, false, nil
	} else {
		return false, "", "", false, nil
	}
}


func (checkout *AmazonCheckout) AmazonCheckoutTask() (bool, bool, error) {
	checkout.Active = true
	checkout.Id = uuid.New().String()
	
	m, _ := mimic.Chromium(mimic.BrandChrome, latestVersion)
	
	checkout.Mimicer = *m

	parsedCookies := cookieHeader(checkout.Cookies)

	for _, cookie := range parsedCookies {
		if cookie.Name == "session-id" {
			checkout.SessionId = cookie.Value
			break
		}
	}

	client := &http.Client{Transport: m.ConfigureTransport(&http.Transport{})}

	for retries := 0; retries < checkout.MaxRetries; retries++ {
		checkout.logMessage("Generating Checkout")
		if !checkout.Active {return false, false, nil}
		cartSuccess, purchaseId, csrfToken, isBanned, err := checkout.AmazonAddToCartV1(client)
		if err != nil {
			return false, false, err
		} else if isBanned {
			if err != nil {
				return false, false, nil
			}
		}

		if cartSuccess {
			checkout.logMessage("Placing Order")
			if !checkout.Active {return false, false, nil}
			success, orderNum, isBanned, err := checkout.AmazonPlaceOrderV1(client, purchaseId, csrfToken) //Todo: add ability to fetch order number, currently returns empty string
			if err != nil {
				return false, false, err
			} else if isBanned {
				if err != nil {
					return false, false, nil
				}
			}

			if success {
				checkout.OrderNum = orderNum
				checkout.logMessage("Order Placed")
				return true, false, nil
			}
		}
		time.Sleep(time.Second * time.Duration(checkout.RetryDelay))
	}
	return false, false, nil
}

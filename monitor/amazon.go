package monitor

import (
	"fmt"
	"github.com/ffeathers/Elektra-Auto-Checkout/elektra"
	ua "github.com/wux1an/fake-useragent"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func amazonCheckStock(client *http.Client, monitorData *elektra.AmazonMonitorData, apiToken string) (bool, bool) {
	acceptheader := "application/vnd.com.amazon.api+json; type=\"cart.add-items/v1\""
	contentheader := "application/vnd.com.amazon.api+json; type=\"cart.add-items.request/v1\""

	var data = strings.NewReader(fmt.Sprintf(`{"items":[{"asin":"%s","offerListingId":"%s","quantity":1}]}`, monitorData.Sku, monitorData.OfferId))
	req, err := http.NewRequest("POST", "https://data.amazon.com/api/marketplaces/ATVPDKIKX0DER/cart/carts/retail/items", data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("x-api-csrf-token", apiToken)
	req.Header.Set("Content-Type", contentheader)
	req.Header.Set("Accept", acceptheader)
	req.Header.Set("User-Agent", monitorData.UserAgent)
	req.Header.Set("accept-language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		return true, false //In stock, api token refresh is not required
	} else if resp.StatusCode == 404 {
		return false, true //Out of stock, but an api token refresh is required
	}

	return false, false //Usually status code 422 (out of stock) but an api token refresh is not required
}

func getApiToken(client *http.Client, monitorData *elektra.AmazonMonitorData) string {
	url := "https://www.amazon.com/gp/aw/d/B00M382RJO" //One of many Amazon product pages that contains an embedded api token

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", monitorData.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	apiToken := elektra.Parse(string(body), "\"csrfToken\":\"", "\",\"baseAsin\"")
	return apiToken
}

func createSession(client *http.Client, monitorData *elektra.AmazonMonitorData) {
	url := "https://www.amazon.com/gp/aws/cart/add-res.html?Quantity.1=1&OfferListingId.1="

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", monitorData.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func AmazonMonitorTask(monitorData *elektra.AmazonMonitorData) {
	client := elektra.CreateClient(monitorData.UseProxies, monitorData.Proxies)

	if monitorData.UserAgent == "" {
		monitorData.UserAgent = ua.RandomType(ua.Desktop)
	}

	log.Println("Getting Session")
	createSession(client, monitorData)

	log.Println("Getting API Token")
	apiToken := getApiToken(client, monitorData)

	for {
		log.Println("Checking Stock")
		inStock, refreshRequired := amazonCheckStock(client, monitorData, apiToken)
		if inStock {
			return
		} else {
			if refreshRequired {
				apiToken = getApiToken(client, monitorData)
			}
		}

		time.Sleep(time.Second * time.Duration(monitorData.PollingInterval))
	}
}

package monitor

import (
	"elektra/checkout"
	"fmt"
	"github.com/wux1an/fake-useragent"
	"github.com/obito/cclient"
	utls "github.com/refraction-networking/utls"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"elektra"
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
	
	apiToken := checkout.Parse(string(body), "\"csrfToken\":\"", "\",\"baseAsin\"")
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


	/*body, _ := ioutil.ReadAll(resp.Body)
	sessionid := (Parse(string(body), "ue_sid='", "',\nue_mid='"))
	return string(sessionid)*/
}

func AmazonMonitorTask(monitorData *elektra.AmazonMonitorData) {
	var client *http.Client
	if monitorData.UseProxies {
		rand.Seed(time.Now().Unix())
		proxy := "http://" + monitorData.Proxies[rand.Intn(len(monitorData.Proxies))] //Only works with IP authenticated proxies atm (IP:Port), not yet with User:Pass:IP:Port proxies
		
		client, err := cclient.NewClient(utls.HelloFirefox_Auto, true, proxy) //Create an http client with a Firefox TLS fingerprint, set automatic storage of cookies to true, and use a proxy
		_ = client
		if err != nil {
			log.Fatal(err)
		}
	} else {
		client, err := cclient.NewClient(utls.HelloFirefox_Auto, true) //Create an http client with a Firefox TLS fingerprint, set automatic storage of cookies to true
		_ = client
		if err != nil {
			log.Fatal(err)
		}
	}

	if monitorData.UserAgent == "" {
		monitorData.UserAgent = ua.RandomType(ua.Desktop)
	}
	
	createSession(client, monitorData)
	apiToken := getApiToken(client, monitorData)
	
	for {
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

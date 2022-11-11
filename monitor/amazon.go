package monitor

import (
	"fmt"
	"github.com/ffeathers/Elektra-Auto-Checkout/elektra"
	"github.com/google/uuid"
	ua "github.com/wux1an/fake-useragent"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/tidwall/gjson"
)

type AmazonMonitor struct {
	Id              string
	UserAgent       string
	Proxy           string
	UseProxy        bool
	Price           float64
	PollingInterval int
	Sku             string
	OfferId         string
	LoggingDisabled bool
	Active          bool
}

func (monitor *AmazonMonitor) logMessage(msg string) {
	if !monitor.LoggingDisabled {
		log.Println(fmt.Sprintf("[Task %s] [Amazon] %s", monitor.Id, msg))
	}
}

func (monitor *AmazonMonitor) Cancel() {
	monitor.Active = false
	monitor.logMessage("Task canceled")
	//add exit code
}

func ParseV2(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	e += s + e - 1
	return str[s:e]
}



func (monitor *AmazonMonitor) AmazonCheckStock(client *http.Client, apiToken string) (bool, bool, bool, error) {
	acceptheader := "application/vnd.com.amazon.api+json; type=\"cart.add-items/v1\""
	contentheader := "application/vnd.com.amazon.api+json; type=\"cart.add-items.request/v1\""

	var data = strings.NewReader(fmt.Sprintf(`{"items":[{"asin":"%s","offerListingId":"%s","quantity":1}]}`, monitor.Sku, monitor.OfferId))
	req, err := http.NewRequest("POST", "https://data.amazon.com/api/marketplaces/ATVPDKIKX0DER/cart/carts/retail/items", data)
	if err != nil {
		return false, false, false, err
	}

	req.Header.Set("x-api-csrf-token", apiToken)
	req.Header.Set("Content-Type", contentheader)
	req.Header.Set("Accept", acceptheader)
	req.Header.Set("User-Agent", monitor.UserAgent)
	req.Header.Set("accept-language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return false, false, false, err
	}

	if resp.StatusCode == 200 {
		return true, false, false, nil //In stock, api token refresh is not required
	} else if resp.StatusCode == 404 {
		return false, true, false, nil //Out of stock, but an api token refresh is required
	} else if resp.StatusCode == 422 { //Usually status code 422 (out of stock) but an api token refresh is not required
		return false, false, false, nil
	}

	return false, false, true, nil
}

func (monitor *AmazonMonitor) AmazonCheckStockV2(client *http.Client, apiToken string) (bool, bool, bool, error) {
	acceptheader := "application/vnd.com.amazon.api+json; type=\"collection(product/v2)/v1\"; expand=\"buyingOptions[].price(product.price/v1),productImages(product.product-images/v2)\""

	req, err := http.NewRequest("GET", "https://data.amazon.com/api/marketplaces/ATVPDKIKX0DER/products/" + monitor.Sku, nil)
	if err != nil {
		return false, false, false, err
	}

	req.Header.Set("x-api-csrf-token", apiToken)
	req.Header.Set("Accept", acceptheader)
	req.Header.Set("User-Agent", monitor.UserAgent)
	req.Header.Set("accept-language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return false, false, false, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	result := gjson.Get(string(body), "entities.0")

	if strings.Contains(result, "amount") {
		price := gjson.Get(result, "entity.buyingOptions.0.price.entity.priceToPay.moneyValueOrRange.value.amount").Float()
		if price <= monitor.Price {
			//in stock because price shows up and is equal to or lower than limit
			return true, false, false, nil
		}
	}

	return false, false, true, nil
}

func (monitor *AmazonMonitor) GetApiToken(client *http.Client) (string, error) {
	url := "https://www.amazon.com/Ring-Video-Doorbell-3/dp/B0849J7W5X" //One of many Amazon product pages that contains an embedded api token

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", monitor.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	tokenStr := Parse(string(body), `","aapiAjaxEndpoint":"data.amazon.com","csrfToken":"`, `"}</script>`)
	tokenSplit := strings.Split(tokenStr, "\"}<")
	apiToken := tokenSplit[0]
	return apiToken, nil
}

func (monitor *AmazonMonitor) CreateSession(client *http.Client) error {
	url := "https://www.amazon.com"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", monitor.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (monitor *AmazonMonitor) AmazonMonitorTaskV2() (bool, error) {
	var inStock, refreshRequired, isBanned bool
	var apiToken string

	monitor.Active = true
	monitor.Id = uuid.New().String()

	client, err := elektra.CreateClient(monitor.Proxy)
	if err != nil {
		return false, err
	}

	monitor.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36"

	monitor.logMessage("Getting session")
	if !monitor.Active {return false, nil}
	monitor.CreateSession(client)
	

	monitor.logMessage("Getting API token")
	if !monitor.Active {return false, nil}
	apiToken, err = monitor.GetApiToken(client)
	if err != nil {
		return false, err
	}

	for monitor.Active {
		monitor.logMessage("Checking stock")
		inStock, refreshRequired, isBanned, err = monitor.AmazonCheckStockV2(client, apiToken)
		if err != nil {
			return isBanned, err
		}
		if inStock {
			return isBanned, nil
		} else {
			if refreshRequired {
				if !monitor.Active {return false, nil}
				apiToken, err = monitor.GetApiToken(client)
				if err != nil {
					return isBanned, err
				}
			}
		}

		time.Sleep(time.Second * time.Duration(monitor.PollingInterval))
	}

	return false, nil
}

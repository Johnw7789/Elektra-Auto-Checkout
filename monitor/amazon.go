package monitor

import (
	"errors"
	"io"

	"github.com/tidwall/gjson"

	"github.com/Johnw7789/Elektra-Auto-Checkout/shr"

	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

// * AmazonTask continuously refreshes the api token and checks stock until the cancel channel is triggered
func (m *MonitorClient) AmazonTask(priceLimit float64) error {
	err := m.reqSession()
	if err != nil {
		return err
	}

	apiToken, err := m.getApiToken()
	if err != nil {
		return err
	}

	tokenUpdated := time.Now()

	errCount := 5

	for {
		select {
		case <-m.cancelChannel:
			return nil
		default:
			if time.Since(tokenUpdated).Hours() > 1 {
				apiToken, err = m.getApiToken()
				if err != nil {
					errCount++
				}
			}

			m.logMessage("Amazon", "Checking stock")
			inStock, err := m.getAmazonStock(priceLimit, apiToken)
			if err != nil && errCount > 5 {
				errCount++
			}

			if errCount > 5 {
				return err
			}

			if inStock {
				m.logMessage("Amazon", "In stock")
				m.AlertChannel <- true
				return nil
			} else {
				m.logMessage("Amazon", "Out of stock")
			}

			time.Sleep(time.Millisecond * time.Duration(m.opts.Delay))
		}
	}
}

// * Evaluate whether in stock based off resp json
func (m *MonitorClient) getAmazonStock(priceLimit float64, apiToken string) (bool, error) {
	body, err := m.reqAmazonStock(apiToken)
	if err != nil {
		return false, err
	}

	entity := gjson.Get(string(body), "entities.0").String()

	if strings.Contains(entity, "amount") {
		price := gjson.Get(entity, "entity.buyingOptions.0.price.entity.priceToPay.moneyValueOrRange.value.amount").Float()
		if price <= priceLimit {
			// * In stock because price shows up and is equal to or lower than limit
			return true, nil
		}
	}

	return false, nil
}

// * Get the api token from the Amazon product page
func (m *MonitorClient) getApiToken() (string, error) {
	body, err := m.reqApiToken()
	if err != nil {
		return "", err
	}

	tokenStr := shr.ParseV2(string(body), `","aapiAjaxEndpoint":"data.amazon.com","csrfToken":"`, `"}</script>`)
	tokenSplit := strings.Split(tokenStr, "\"}<")

	if len(tokenSplit) < 1 {
		return "", errors.New("could not parse api token")
	}

	return tokenSplit[0], nil
}

func (m *MonitorClient) reqAmazonStock(apiToken string) (string, error) {
	req, err := http.NewRequest("GET", AmazonBaseUrl+m.opts.Sku, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("x-api-csrf-token", apiToken)
	req.Header.Set("Accept", AmazonAcceptHeader)
	req.Header.Set("origin", "https://www.amazon.com")

	req.Header = m.defaultHeaders(req.Header.Clone())

	resp, err := m.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (m *MonitorClient) reqApiToken() (string, error) {
	req, err := http.NewRequest("GET", AmazonTokenUrl, nil)

	req.Header = m.defaultHeaders(req.Header.Clone())

	resp, err := m.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (m *MonitorClient) reqSession() error {
	req, err := http.NewRequest("GET", "https://www.amazon.com", nil)

	req.Header = m.defaultHeaders(req.Header.Clone())

	resp, err := m.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

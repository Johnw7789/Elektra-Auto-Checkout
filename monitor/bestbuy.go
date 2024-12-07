package monitor

import (
	"errors"
	"io"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

// * BestbuyTask continuously checks stock until the cancel channel is triggered
func (m *MonitorClient) BestbuyTask() error {
	errCount := 5

	for {
		select {
		case <-m.cancelChannel:
			return nil
		default:
			m.logMessage("Bestbuy", "Checking stock")
			inStock, err := m.getBestbuyStock()
			if err != nil && errCount > 5 {
				return err
			}

			if inStock {
				m.logMessage("Bestbuy", "In stock")
				m.AlertChannel <- true
				return nil
			} else {
				m.logMessage("Bestbuy", "Out of stock")
			}

			time.Sleep(time.Millisecond * time.Duration(m.opts.Delay))
		}
	}
}

func (m *MonitorClient) getBestbuyStock() (bool, error) {
	body, err := m.reqBestbuyStock()
	if err != nil {
		return false, err
	}

	// * ADD_TO_CART is the only 100% reliable indicator of stock, but CHECK_STORES is probably still worth including as it means at least 1 store has stock somewhere
	if strings.Contains(string(body), "CHECK_STORES") || strings.Contains(string(body), "ADD_TO_CART") {
		return true, nil
	}

	return false, nil
}

func (m *MonitorClient) reqBestbuyStock() (string, error) {
	req, err := http.NewRequest("GET", BestbuyBaseUrl+m.opts.Sku, nil)
	if err != nil {
		return "", nil
	}

	req.Header.Set("authority", "www.bestbuy.com")
	req.Header.Set("host", "www.bestbuy.com")
	req.Header.Set("accept", "*/*")
	req.Header.Set("x-client-id", "FRV")
	req.Header.Set("Connection", "keep-alive")

	req.Header = m.defaultHeaders(req.Header.Clone())

	resp, err := m.HttpClient.Do(req)
	if err != nil {
		return "", nil
	}

	if resp.StatusCode != 200 {
		return "", errors.New("non 200 status code: " + resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

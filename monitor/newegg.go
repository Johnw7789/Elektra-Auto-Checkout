package monitor

import (
	"errors"
	"io"
	"time"

	http "github.com/bogdanfinn/fhttp"
	"github.com/tidwall/gjson"
)

// * NeweggTask continuously checks stock until the cancel channel is triggered
func (m *MonitorClient) NeweggTask() error {
	errCount := 5

	for {
		select {
		case <-m.cancelChannel:
			return nil
		default:
			m.logMessage("Newegg", "Checking stock")
			inStock, err := m.getNeweggStock()
			if err != nil && errCount > 5 {
				return err
			}

			if inStock {
				m.logMessage("Newegg", "In stock")
				m.AlertChannel <- true
				return nil
			} else {
				m.logMessage("Newegg", "Out of stock")
			}

			time.Sleep(time.Millisecond * time.Duration(m.opts.Delay))
		}
	}
}

// * Evaluate whether in stock based off resp json
func (m *MonitorClient) getNeweggStock() (bool, error) {
	body, err := m.reqNeweggStock()
	if err != nil {
		return false, err
	}

	return gjson.Get(body, "MainItem.Instock").Bool(), nil
}

func (m *MonitorClient) reqNeweggStock() (string, error) {
	req, err := http.NewRequest("GET", NeweggBaseUrl+m.opts.Sku, nil)
	if err != nil {
		return "", nil
	}

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

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
)

type BestbuyMonitor struct {
	Id		string
	UserAgent       string
	Proxy           string
	PollingInterval int
	Sku             string
	LoggingDisabled bool
	Active		bool
}
func (monitor *BestbuyMonitor) logMessage(msg string) {
	if !monitor.LoggingDisabled {
		log.Println(fmt.Sprintf("[Task %s] [BestBuy] %s", monitor.Id, msg))
	}
}

func (monitor *BestbuyMonitor) Cancel() {
	monitor.Active = false
	monitor.logMessage("Task canceled")
	//add exit code
}

func (monitor *BestbuyMonitor) BestbuyCheckStock(client *http.Client) (bool, bool, error) {
  	req, err := http.NewRequest("GET", "https://www.bestbuy.com/button-state/api/v5/button-state?skus=" + monitor.Sku + "&context=pdp&source=buttonView", nil)
	if err != nil {
		return false, false, nil
	}
	req.Header.Set("authority", "www.bestbuy.com")
	req.Header.Set("host", "www.bestbuy.com")
	req.Header.Set("user-agent", monitor.UserAgent)
	req.Header.Set("accept", "*/*")
	req.Header.Set("x-client-id", "FRV")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		return false, false, nil
	}
  
	if resp.StatusCode == 200 {
    		body, _ := ioutil.ReadAll(resp.Body)
    		if strings.Contains(string(body), "CHECK_STORES") || strings.Contains(string(body), "ADD_TO_CART") {
      			return false, true, nil
		}
  	} else {
   		 monitor.logMessage("Status: " + resp.Status)
   		 return true, false, nil
  	}
  
  	return false, false, nil
}

func (monitor *BestbuyMonitor) BestbuyMonitorTask() (bool, error) {
	monitor.Active = true
	monitor.Id = uuid.New().String()

	client, err := elektra.CreateClient(monitor.Proxy)
	if err != nil {
		return false, err
	}
 
	if monitor.UserAgent == "" {
		monitor.UserAgent = ua.RandomType(ua.Desktop)
	}
  
	for monitor.Active {
		monitor.logMessage("Checking stock")

		isBanned, inStock, err := monitor.BestbuyCheckStock(client)
		if err != nil {
			return isBanned, err
		}

		if inStock {
			return isBanned, nil
		}

		time.Sleep(time.Second * time.Duration(monitor.PollingInterval))
	}

	return false, nil
}

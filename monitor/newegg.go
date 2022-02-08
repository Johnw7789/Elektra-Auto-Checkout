package monitor

import (
	//"fmt"
	"github.com/ffeathers/Elektra-Auto-Checkout/elektra"
	ua "github.com/wux1an/fake-useragent"
	//"io/ioutil"
	"log"
	"net/http"
	//"strings"
	"time"
)


type NeweggMonitor struct {
	UserAgent       string
	Proxy           string
	PollingInterval int
	Sku             string
}

func (monitor *NeweggMonitor) neweggCheckStock (client *http.Client,) (bool, bool, error) {
  	req , err := http.NewRequest("GET", "https://www.newegg.com/product/api/ProductRealtime?ItemNumber=" + monitor.Sku, nil)
	if err != nil {
		return false, false, nil
	}
	
	req.Header.Set("user-agent", monitor.UserAgent)
	req.Header.Set("accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return false, false, nil
	}

	if resp.StatusCode == 200 {
		//body, _ := ioutil.ReadAll(resp.Body)
	} else {
		return false, true, nil
	}

  	return false, false, nil
}

func (monitor *NeweggMonitor) NeweggMonitorTask() (bool, error) {
	client, err := elektra.CreateClient(monitor.Proxy)
	if err != nil {
		return false, err
	}
 
	if monitor.UserAgent == "" {
		monitor.UserAgent = ua.RandomType(ua.Desktop)
	}
  
	for {
		log.Println("Checking Stock")
		inStock, isBanned, err := monitor.neweggCheckStock(client)
		if err != nil {
			return false, err
		}

		if inStock {
			return false, nil
		} else if isBanned {
			return true, nil
		}

		time.Sleep(time.Second * time.Duration(monitor.PollingInterval))
	}
}



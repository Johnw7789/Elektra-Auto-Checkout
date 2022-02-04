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

func neweggCheckStock(client *http.Client, monitorData *elektra.NeweggMonitorData) bool {
  	req , err := http.NewRequest("GET", "https://www.newegg.com/product/api/ProductRealtime?ItemNumber=" + monitorData.Sku, nil)
	if err != nil {
		log.Fatal(err)
	}
	
	req.Header.Set("user-agent", monitorData.UserAgent)
	req.Header.Set("accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
  
	//Todo: check if in stock
  
  	return false
}

func NeweggMonitorTask(monitorData *elektra.NeweggMonitorData) {
	client := elektra.CreateClient(monitorData.UseProxies, monitorData.Proxies)
 
	if monitorData.UserAgent == "" {
		monitorData.UserAgent = ua.RandomType(ua.Desktop)
	}
  
	for {
		log.Println("Checking Stock")
		inStock := neweggCheckStock(client, monitorData)
		if inStock {
			return
		}

		time.Sleep(time.Second * time.Duration(monitorData.PollingInterval))
	}
}



package monitor

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/ffeathers/Elektra-Auto-Checkout/elektra"
	ua "github.com/wux1an/fake-useragent"
	"log"
	"net/http"
	"time"
)


type NeweggMonitor struct {
	Id         	string
	UserAgent       string
	Proxy           string
	PollingInterval int
	Sku             string
	LoggingDisabled bool
	Active          bool
}

func (monitor *NeweggMonitor) logMessage(msg string) {
	if !monitor.LoggingDisabled {
		log.Println(fmt.Sprintf("[Task %s] [Newegg] %s", monitor.Id, msg))
	}
}

func (monitor *NeweggMonitor) Cancel() {
	monitor.Active = false
	log.Println(fmt.Sprintf("[Task %s] Task canceled", monitor.Id))
	//add exit code
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
	monitor.Active = true
	monitor.Id = uuid.New().String()

	client, err := elektra.CreateClient(monitor.Proxy)
	if err != nil {
		return false, err
	}
 
	if monitor.UserAgent == "" {
		monitor.UserAgent = ua.RandomType(ua.Desktop)
	}
  
	for {
		monitor.logMessage("Checking Stock")
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



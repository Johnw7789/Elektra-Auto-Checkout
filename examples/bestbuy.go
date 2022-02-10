package examples

import (
	"fmt"
	"github.com/ffeathers/Elektra-Auto-Checkout/elektra"
	"github.com/ffeathers/Elektra-Auto-Checkout/monitor"
	"log"
)

func main() {
	monitor := monitor.BestbuyMonitor{
		UserAgent:       "",
		Proxy:           "",
		PollingInterval: 3,
		Sku:             "5457800",
	}

	banned, err := monitor.BestbuyMonitorTask()
	if err != nil {
		log.Fatal(err)
	}
	
	if banned {
		log.Println(fmt.Sprintf("Your IP is flagged", monitor.Sku))
	} else {
		log.Println(fmt.Sprintf("SKU %s: In Stock", monitor.Sku))
	}
}

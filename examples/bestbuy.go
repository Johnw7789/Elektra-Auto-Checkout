package examples

import (
	"fmt"
	"github.com/ffeathers/Elektra-Auto-Checkout/elektra"
	"github.com/ffeathers/Elektra-Auto-Checkout/monitor"
	"log"
)

func main() {
	monitorData := elektra.BestbuyMonitorData{
		UserAgent:       "",
		UseProxies:      false,
		PollingInterval: 3,
		Sku:             "5457800",
	}

	monitor.BestbuyMonitorTask(&monitorData)

	log.Println(fmt.Sprintf("SKU %s: In Stock", monitorData.Sku))
}

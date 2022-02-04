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
		Sku:             "B0873C4C67",
	}

	monitor.BestbuyMonitorTask(&monitorData)

	log.Println(fmt.Sprintf("SKU %s: In Stock", monitorData.Sku))
}

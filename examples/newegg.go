package examples

import (
	"fmt"
	"github.com/ffeathers/Elektra-Auto-Checkout/elektra"
	"github.com/ffeathers/Elektra-Auto-Checkout/monitor"
	"log"
)

func main() {
	monitorData := elektra.NeweggMonitorData{
		UserAgent:       "",
		UseProxies:      false,
		PollingInterval: 3,
		Sku:             "N82E16835181166",
	}

	monitor.NeweggMonitorTask(&monitorData)

	log.Println(fmt.Sprintf("SKU %s: In Stock", monitorData.Sku))
}

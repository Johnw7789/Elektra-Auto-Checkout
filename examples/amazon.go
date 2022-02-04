package examples

import (
	"fmt"
	"github.com/ffeathers/Elektra-Auto-Checkout/elektra"
	"github.com/ffeathers/Elektra-Auto-Checkout/monitor"
	"log"
)

func main() {
	monitorData := elektra.AmazonMonitorData{
		UserAgent:       "",
		UseProxies:      false,
		PollingInterval: 3,
		Sku:             "B0873C4C67",
		OfferId:         "5%2BU3RbI4MrLxJP1riew3ktYPNAEuKmceCPF1BTaKdwF9bGnxPX3cfIChUFRKBusiTPTd3gJEB9Az0V3TlZw0po6Mob%2BYvq37tir2AWHORVYNxN9kBTPxMuvTkuiELMuz3q9BWdzZKsylbBhRmq7cAHQgq7p9VSdR5e6J%2BWxORLR95D2He%2BodtT4wtctu24wt",
	}

	monitor.AmazonMonitorTask(&monitorData)

	log.Println(fmt.Sprintf("SKU %s: In Stock", monitorData.Sku))
}

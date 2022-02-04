package examples

import (
	"fmt"
	"github.com/ffeathers/Elektra-Auto-Checkout/elektra"
	"github.com/ffeathers/Elektra-Auto-Checkout/monitor"
	"log"
)

func main() {
	cookieString := ""
	sku := "B0873C4C67"
	offerId := "5%2BU3RbI4MrLxJP1riew3ktYPNAEuKmceCPF1BTaKdwF9bGnxPX3cfIChUFRKBusiTPTd3gJEB9Az0V3TlZw0po6Mob%2BYvq37tir2AWHORVYNxN9kBTPxMuvTkuiELMuz3q9BWdzZKsylbBhRmq7cAHQgq7p9VSdR5e6J%2BWxORLR95D2He%2BodtT4wtctu24wt"
	
	monitorData := elektra.AmazonMonitorData{
		UserAgent:       "",
		UseProxies:      false,
		PollingInterval: 3,
		Sku:             sku,
		OfferId:         offerid,
	}

	monitor.AmazonMonitorTask(&monitorData)

	log.Println(fmt.Sprintf("SKU %s: In Stock, Initiating Checkout", monitorData.Sku))
	
	
	checkoutData := elektra.AmazonCheckoutData{
  		UserAgent: "",
  		UseProxies: true,
 		Proxies: []string{"IP:Port", "IP:Port"},
  		Cookies: cookieString,
  		MaxRetries: 5,
  		RetryDelay: 3,
  		Sku: sku,
 		OfferId: offerId,
	}
  
	orderSuccess := checkout.AmazonCheckoutTask(&checkoutData) 
	if orderSuccess {
  		log.Println("Checkout successful | order number: " + checkoutData.OrderNum)
	}
}

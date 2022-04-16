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
		
	monitor := monitor.AmazonMonitor{
		UserAgent: "",
		Proxy: "",
		PollingInterval: 3,
		Sku: "ASIN",
		OfferId: "OfferId",
	}

	banned, err := monitor.AmazonMonitorTask()
	if err != nil {
		log.Fatal(err)
	}

	if banned {
		log.Println("Your IP was flagged")
	} else {
		log.Println(fmt.Sprintf("SKU %s: In Stock, Initiating Checkout", monitorData.Sku))
		
			
		amazonCheckout := checkout.AmazonCheckout{
		  UserAgent: "",
		  Proxy: "",
		  Cookies: "exampleCookie=exampleValue",
		  MaxRetries: 5,
		  RetryDelay: 3,
		  Sku: "ASIN",
		  OfferId: "OfferId",
		}

		orderSuccess, isBanned, err := amazonCheckout.AmazonCheckoutTask() 
		if err != nil {
		  log.Fatal(err)
		}

		if isBanned {
		  //ip banned
		} else if orderSuccess {
		  log.Println("Checkout successful | order number: " + amazonCheckout.OrderNum)
		}
	}
}

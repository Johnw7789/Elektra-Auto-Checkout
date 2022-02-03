package elektra

import (
	"log"
	"orderapi/checkout"
	"orderapi/monitor"
)

func NewCheckoutTask(checkoutData *CheckoutData) {
	switch checkoutData.Site {
	case "amazon":
		checkout.Amazon()
	case "bestbuy":
		log.Println("Site not yet supported")
	case "newegg":
		log.Println("Site not yet supported")
	}
}

func NewMonitorTask(monitorData *MonitorData) {
	switch monitorData.Site {
	case "amazon":
		monitor.Amazon()
	case "bestbuy":
		monitor.Bestbuy()
	case "newegg":
		monitor.Newegg()
	}
}

func NewLoginSession(accountData *AccountData) {
	switch accountData.Site {
	case "amazon":
		checkout.Amazon()
	case "bestbuy":
		log.Println("Site not yet supported")
	case "newegg":
		log.Println("Site not yet supported")
	}
}

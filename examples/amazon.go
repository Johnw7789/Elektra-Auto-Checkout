package examples

import (
	"log"

	"github.com/Johnw7789/Elektra-Auto-Checkout/monitor"
)

func TestAmazonMonitor() {
	opts := monitor.MonitorOpts{
		Sku:   "B071JM699B",
		Delay: 3000,
		Proxy:    "http://localhost:8888", // * Sniff using local proxy
		Logging: true,
	}

	monitor, err := monitor.NewMonitorClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	priceLimit := 6.17

	go monitor.AmazonTask(priceLimit)

	inStock := <-monitor.AlertChannel

	log.Println("In stock:", inStock)
}

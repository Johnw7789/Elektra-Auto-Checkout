package examples

import (
	"log"

	"github.com/Johnw7789/Elektra-Auto-Checkout/monitor"
)

func TestNeweggMonitor() {
	opts := monitor.MonitorOpts{
		Sku:     "N82E16824012083",
		Delay:   3000,
		Proxy:    "http://localhost:8888", // * Sniff using local proxy
		Logging: true,
	}

	monitor, err := monitor.NewMonitorClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	go monitor.NeweggTask()

	inStock := <-monitor.AlertChannel

	log.Println("In stock:", inStock)
}

package examples

import (
	"log"

	"github.com/Johnw7789/Elektra-Auto-Checkout/login"
	"github.com/Johnw7789/Elektra-Auto-Checkout/monitor"
)

func TestBestbuyMonitor() {
	opts := monitor.MonitorOpts{
		Sku:     "6473498",
		Delay:   3000,
		Proxy:   "http://localhost:8888", // * Sniff using local proxy
		Logging: true,
	}

	monitor, err := monitor.NewMonitorClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	go monitor.BestbuyTask()

	inStock := <-monitor.AlertChannel

	log.Println("In stock:", inStock)
}

// * Deprecated
func TestBestbuyLogin() {
	opts := login.LoginOpts{
		Email:    "",
		Password: "",
		Proxy:    "http://localhost:8888", // * Sniff using local proxy
		Logging:  true,
	}

	lc, err := login.NewLoginClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	success, err := lc.BestbuyTask()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Login success:", success)

	cookies := lc.GetCookieStr()
	log.Println("Cookies:", cookies)
}

package elektra

import (
	"github.com/obito/cclient"
	utls "github.com/refraction-networking/utls"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type BillingData struct {
	Email     string
	Phone     string
	FirstName string
	LastName  string
	State     string
	City      string
	ZipCode   string
	Address1  string
	Address2  string
	CardType  string
	CardNum   string
	ExpMonth  string
	ExpYear   string
	CardName  string
	Cvv       string
}

type AccountData struct {
	Username string
	Password string
	Email    string
	Phone    string
}

type AmazonMonitorData struct {
	UserAgent       string
	Proxies         []string
	UseProxies      bool
	PollingInterval int
	Sku             string
	OfferId         string
}

type AmazonCheckoutData struct {
	UserAgent  string
	SessionId  string
	Cookies    string
	Proxies    []string
	UseProxies bool
	RetryDelay int
	MaxRetries int
	Sku        string
	OfferId    string
	OrderNum   string
}

type BestbuyMonitorData struct {
	UserAgent       string
	Proxies         []string
	UseProxies      bool
	PollingInterval int
	Sku             string
}

type NeweggMonitorData struct {
	Delay int
	Sku   string
}

func Parse(value string, a string, b string) string {
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func CreateClient(useProxies bool, proxies []string) *http.Client {
	if useProxies {
		rand.Seed(time.Now().Unix())
		proxy := "http://" + proxies[rand.Intn(len(proxies))] //Only works with IP authenticated proxies atm (IP:Port), not yet with User:Pass:IP:Port proxies

		client, err := cclient.NewClient(utls.HelloFirefox_Auto, true, proxy) //Create an http client with a Firefox TLS fingerprint, set automatic storage of cookies to true, and use a proxy
		if err != nil {
			log.Fatal(err)
		}
		
		log.Println("Created client")

		return &client
	} else {
		client, err := cclient.NewClient(utls.HelloFirefox_Auto, true) //Create an http client with a Firefox TLS fingerprint, set automatic storage of cookies to true
		if err != nil {
			log.Fatal(err)
		}
		
		log.Println("Created client")

		return &client
	}
}

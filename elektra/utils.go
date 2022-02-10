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


func CreateClient(proxy string) (*http.Client, error) {
	if proxy != "" {
		proxyUrl := "http://" + proxy //Only works with IP authenticated proxies atm (IP:Port), not yet with User:Pass:IP:Port proxies

		client, err := cclient.NewClient(utls.HelloFirefox_Auto, true, proxyUrl) //Create an http client with a Firefox TLS fingerprint, set automatic storage of cookies to true, and use a proxy
		if err != nil {
			return nil, err
		}
		
		log.Println("Created Client")

		return &client, nil
	} else {
		client, err := cclient.NewClient(utls.HelloFirefox_Auto, true) //Create an http client with a Firefox TLS fingerprint, set automatic storage of cookies to true
		if err != nil {
			return nil, err
		}
		
		log.Println("Created Client")

		return &client, nil
	}
}

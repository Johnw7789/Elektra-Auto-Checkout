package monitor

func bestbuyCheckStock(client *http.Client, monitorData *elektra.BestbuyMonitorData) bool {
  	req, err := http.NewRequest("GET", "https://www.bestbuy.com/button-state/api/v5/button-state?skus=" + monitorData.Sku + "&context=pdp&source=buttonView", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "www.bestbuy.com")
	req.Header.Set("host", "www.bestbuy.com")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.64")
	req.Header.Set("accept", "*/*")
	req.Header.Set("x-client-id", "FRV")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
  
	if resp.StatusCode == 200 {
    		bodyText, _ := ioutil.ReadAll(resp.Body)
    		if strings.Contains(string(bodyText), "CHECK_STORES") || strings.Contains(string(bodyText), "ADD_TO_CART") {
      			return true
		}
  	} else {
   		 log.Println(fmt.Sprintf("Status Code: %d", resp.StatusCode))
  	}
  
  	return false
}

func BestbuyMonitorTask((monitorData *elektra.BestbuyMonitorData) {
	client := elektra.CreateClient(monitorData.UseProxies, monitorData.Proxies)
 
	if monitorData.UserAgent == "" {
		monitorData.UserAgent = ua.RandomType(ua.Desktop)
	}
  
	for {
		log.Println("Checking Stock")
		inStock := bestbuyCheckStock(client, monitorData)
		if inStock {
			return
		}

		time.Sleep(time.Second * time.Duration(monitorData.PollingInterval))
	}
}

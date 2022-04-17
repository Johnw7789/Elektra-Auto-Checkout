# Elektra
## About This Project
Elektra is designed to automate the process of inventory checking, purchasing items, and generating account login sessions for various commercial and retail sites.

Amazon is the first of many sites to come. Expect weekly additions, though some sites may only receive login and monitor modules and not full fledged checkout.

## Note
This project is not intended for resellers. This is a project for educational purposes / experimentation / to help others get an item which they may need.

## Progress / Roadmap

| **Site** | **Login** | **Monitor** | **Checkout** |
|:---:|:---:|:---:|:---:|
| amazon.com |WIP|`✔`|`✔`|
| bestbuy.com |`✔`|`✔`|WIP|
| newegg.com ||`✔`| |
| evga.com ||||
| target.com ||WIP||
| walmart.com ||WIP||

* Add notifications module (Discord, Slack, Twilio)
* ~~Add auth code fetcher (imap + Gmail)~~ - Complete
* Add account generators

## Installation
##### Requires Go 1.7
``go get github.com/ffeathers/Elektra-Auto-Checkout``

Use ``go mod tidy`` if issues arise with some of Elektra's imported modules.

## Getting Started
Below is some example usage of the Amazon module. You can find additional examples for other sites in the [examples](https://github.com/ffeathers/Elektra-Auto-Checkout/tree/main/examples) folder.

###### Checking stock
If ``UserAgent`` is left empty, a user-agent will be automatically generated for you. ``PollingInterval`` is the delay in seconds for which a monitor will sleep after every stock check. Once a monitor task is started, it will continue to monitor indefinitely until stock is detected, then it will return.

```  
amazonMonitor := monitor.AmazonMonitor{
  UserAgent: "", 
  Proxy: "",
  PollingInterval: 3,
  Sku: "ASIN",
  OfferId: "OfferId",
}
  
amazonMonitor.AmazonMonitorTask() 
  
log.Println(fmt.Sprintf("SKU %s: In Stock", amazonMonitor.Sku))
```
###### Starting a checkout instance
Account ``Cookies`` are needed in order to complete a checkout. You can use cookies from your browser or you can create a session using the Amazon login module (not yet implemented).``RetryDelay`` is the amount of time that a checkout task will sleep if there is an error in the checkout flow, before restarting. ``MaxRetries`` is the maximum amount of checkout attempts a checkout task will make before it returns. If the return value is false, the task was unable to complete a successful checkout after every attempt made. If it is true, then the checkout was succesful and ``OrderNum`` should now be populated with the order number. 

If you would like to use your local IP, you can exclude the Proxy param or leave it as an empty string. 

```
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
```

## License
[MIT](https://choosealicense.com/licenses/mit)

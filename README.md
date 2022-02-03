## About This Project
Elektra is designed to automate the process of inventory checking, purchase automation, and automated logins for various commercial and retail sites.

## Getting Started
###### Checking stock
If ``UserAgent`` is left empty, a user-agent will be automatically generated for you. ``PollingInterval`` is the delay in seconds for which a monitor will sleep after every stock check.

```  
monitorData := elektra.AmazonMonitorData{
  UserAgent: "", 
  PollingInterval: 3,
  Sku: "ASIN",
  OfferId: "OfferId",
}
  
AmazonMonitorTask(monitorData) 
  
//Do something when in stock
```
###### Starting a checkout instance

```
checkoutData := elektra.AmazonCheckoutData{
  UserAgent: "", //If left empty, a user-agent will be randomly generated for you
  PollingInterval: 3,
  Sku: "ASIN",
  OfferId: "OfferId",
}
  
AmazonCheckoutTask(checkoutData) 
  
//Do something when in stock
```

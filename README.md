## About This Project
Elektra is designed to automate the process of inventory checking, purchase automation, and automated logins for various commercial and retail sites.

## Getting Started
###### Checking inventory
If UserAgent is left empty, a user-agent will be automatically generated for you.

```  
monitorData := elektra.AmazonMonitorData{
  UserAgent: "", 
  PollingInterval: 3,
  Sku: "ASIN",
  OfferId: "OfferId",
}
  
AmazonMonitorTask(monitorData) //Checks stock using the designated PollingInterval delay, returns once in stock
  
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
  
AmazonCheckoutTask(checkoutData) //Checks stock using the designated PollingInterval delay, returns once in stock
  
//Do something when in stock
```

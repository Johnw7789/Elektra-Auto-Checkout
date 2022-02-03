## About This Project
Elektra is designed to automate the process of inventory checking, purchasing items, and generating account login sessions for various commercial and retail sites.

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
  
log.Println(fmt.Sprintf("SKU %s: In Stock", monitorData.Sku))
```
###### Starting a checkout instance

```
checkoutData := elektra.AmazonCheckoutData{
  UserAgent: "", //If left empty, a user-agent will be randomly generated for you
  MaxRetries: 5,
  RetryDelay: 3,
  Sku: "ASIN",
  OfferId: "OfferId",
  OrderNum: "",
}
  
orderSuccess := AmazonCheckoutTask(checkoutData) 
if orderSuccess {
  log.Println("Checkout successful | order number: " + checkoutData.OrderNum)
}
```

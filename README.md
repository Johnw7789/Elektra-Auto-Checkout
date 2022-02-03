## About This Project
Elektra is designed to automate the process of inventory checking, purchasing items, and generating account login sessions for various commercial and retail sites.

## Getting Started
###### Checking stock
If ``UserAgent`` is left empty, a user-agent will be automatically generated for you. ``PollingInterval`` is the delay in seconds for which a monitor will sleep after every stock check. Once a monitor task is started, it will continue to monitor indefinitely until stock is detected.

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
``RetryDelay`` is the amount of time that a checkout task will sleep if there is an error in the checkout flow, before restarting. ``MaxRetries`` is the maximum amount of attempts a checkout task will make before it returns. If the return value is false, the task was unable to complete a successful checkout after every attempt made. If it is true, then the checkout was succesful and ``OrderNum`` should now be populated with the order number. 

```
checkoutData := elektra.AmazonCheckoutData{
  UserAgent: "",
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

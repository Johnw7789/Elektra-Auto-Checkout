## About This Project
Elektra is designed to automate the process of inventory checking, purchase automation, and automated logins for various commercial and retail sites.

## Getting Started
###### Checking inventory

```
func main() {
  monitorData := elektra.AmazonMonitorData{
    PollingInterval: 3,
    Sku: "ASIN",
    OfferId: "OfferId",
  }
  
  AmazonMonitorTask(monitorData) //Checks stock using the designated PollingInterval delay, returns once in stock
  
  //Do something when in stock
}
```


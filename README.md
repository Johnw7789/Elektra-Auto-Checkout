# Elektra
## About This Project
Note: Checkout modules and Bestbuy login are now deprecated. Bestbuy login can be made to work by solving the x-grid and x-grid-b encrypted headers. I may update with this functionality in the future but do not have the time currently. 

Elektra is designed to automate the process of inventory checking ~~and in the case of certain sites, login and automated checkout~~.

Elektra does not use any automated browser, all automation is executed through the use of unofficial, public facing APIs.

## Note
This project is not intended for resellers. This is a project for educational purposes / experimentation / to help others get an item which they may need.

## Installation
``go get github.com/Johnw7789/Elektra-Auto-Checkout``

Use ``go mod tidy`` if issues arise with some of Elektra's imported modules.

## Getting Started
Below is example usage of the Amazon, Bestbuy, and Newegg modules. You can find examples in the [examples](https://github.com/Johnw7789/Elektra-Auto-Checkout/tree/main/examples) folder.

###### Checking stock
Please specify the ```Delay``` in milliseconds. The primary method of receiving stock status is done through a Go channel. If there are a number of errors that exceed the limit, the process will be cancelled and the channel will send back false for the stock status. The err can be read as long as it is done so in its own go routine. That is to say, the monitor task should be fired in a goroutine so that the alert channel can be waited for outside of it. If an err is returned then it can be read in the previous goroutine. 

### Amazon
```go
opts := monitor.MonitorOpts{
	Sku:     "B071JM699B", // * ASIN, or the sku of the amazon product
	Delay:   3000, // * Delay of 3 seconds
	Proxy:   "", 
	Logging: true,
}

monitor, err := monitor.NewMonitorClient(opts)
if err != nil {
	log.Fatal(err)
}

// * A price limit must be set as it will be checked to determine if the product is in stock and sold by the correct merchant
priceLimit := 6.17 

go func() {
	// * Optional if wanting to read a potential error, should otherwise just fire like this: go monitor.AmazonTask(priceLimit)
	err := monitor.AmazonTask(priceLimit)
	if err != nil {
		log.Fatal(err)
	}
}()

// * Wait for in stock
inStock := <-monitor.AlertChannel
```

### Bestbuy
```go
opts := monitor.MonitorOpts{
	Sku:     "6473498",
	Delay:   3000,
	Proxy:   "", 
	Logging: true,
}

monitor, err := monitor.NewMonitorClient(opts)
if err != nil {
	log.Fatal(err)
}

go monitor.BestbuyTask()

inStock := <-monitor.AlertChannel
```

### Newegg
```go
opts := monitor.MonitorOpts{
	Sku:     "N82E16824012083",
	Delay:   3000,
	Proxy:    "",
	Logging: true,
}

monitor, err := monitor.NewMonitorClient(opts)
if err != nil {
	log.Fatal(err)
}

go monitor.NeweggTask()

inStock := <-monitor.AlertChannel
```


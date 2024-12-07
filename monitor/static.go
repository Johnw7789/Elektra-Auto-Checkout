package monitor

const (
	DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

	AmazonAcceptHeader = "application/vnd.com.amazon.api+json; type=\"collection(product/v2)/v1\"; expand=\"buyingOptions[].price(product.price/v1),productImages(product.product-images/v2)\""
	AmazonTokenUrl = "https://www.amazon.com/dp/B0BZWRSRWV" // * One of many Amazon product pages that contains an embedded api token

	AmazonBaseUrl  = "https://data.amazon.com/api/marketplaces/ATVPDKIKX0DER/products/"
	BestbuyBaseUrl = "https://www.bestbuy.com/button-state/api/v5/button-state?context=pdp&source=buttonView&skus="
	NeweggBaseUrl  = "https://www.newegg.com/product/api/ProductRealtime?ItemNumber="
)

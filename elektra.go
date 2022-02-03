package elektra

type BillingData struct {
	Email     string
	Phone     string
	FirstName string
	LastName  string
	State     string
	City      string
	ZipCode   string
	Address1  string
	Address2  string
	CardType  string
	CardNum   string
	ExpMonth  string
	ExpYear   string
	CardName  string
	Cvv       string
}

type AccountData struct {
	Username string
	Password string
	Email    string
	Phone    string
}

type AmazonMonitorData struct {
	UserAgent string
	Delay     int
	Sku       string
	OfferId   string
}

type AmazonCheckoutData struct {
	UserAgent  string
	RetryDelay int
	SessionId  string
	Cookies    string
	Sku        string
	OfferId    string
	OrderNum   string
}

type BestbuyMonitorData struct {
	Delay   int
	Sku     string
}

type NeweggMonitorData struct {
	Delay   int
	Sku     string
}


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
	Site     string
	Username string
	Password string
	Email    string
	Phone    string
}

type Account struct {
	Site    string
	Cookies string
}

type MonitorData struct {
	Delay   int
	Site    string
	Sku     string
	Offerid string
}

type CheckoutData struct {
	RetryDelay int
	Site       string
	Sku        string
	Offerid    string
	OrderNum   string
}

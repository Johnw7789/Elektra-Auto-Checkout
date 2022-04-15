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

type AccountLoginData struct {
	Username string
	Password string
	Email    string
	Phone    string
}

type AccountCreationData struct {
	Username string
	Password string
	Email    string
	Phone    string
}





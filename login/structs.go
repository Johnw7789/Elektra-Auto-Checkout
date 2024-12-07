package login

type BestbuyEncryptionData struct {
	EncryptedEmail    string
	EncryptedAgent    string
	EncryptedActivity string
}

type BestbuyLoginData struct {
	VerificationCodeFieldName string
	EncryptedPasswordField    string
	EncryptedAlpha            string
	EmailField                string
	Salmon                    string
	Token                     string
}

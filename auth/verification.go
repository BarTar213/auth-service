package auth

func GenerateVerificationCode() string {
	return randString(20)
}

package model

type ResetPasswordEmail struct {
	Email string
	URL   string
}

type OTPEmail struct {
	Email      string
	DigitOne   string
	DigitTwo   string
	DigitThree string
	DigitFour  string
	DigitFive  string
	DigitSix   string
}

package money

var currencies = []string{
	"RON",
}

func IsCurrencyValid(currency string) bool {
	for _, c := range currencies {
		if c == currency {
			return true
		}
	}
	return false
}

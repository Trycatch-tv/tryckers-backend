package enums

type Country string

const (
	Colombia Country = "colombia"
	Chile    Country = "chile"
)

var ValidCountries = map[Country]struct{}{
	Colombia: {},
	Chile:    {},
}

func IsValidCountry(country string) bool {
	_, ok := ValidCountries[Country(country)]
	return ok
}

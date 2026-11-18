package enums

type Country string

const (
	Colombia Country = "colombia"
	Chile    Country = "chile"
	Honduras Country = "honduras"
)

var ValidCountries = map[Country]struct{}{
	Colombia: {},
	Chile:    {},
	Honduras: {},
}

func IsValidCountry(country string) bool {
	_, ok := ValidCountries[Country(country)]
	return ok
}

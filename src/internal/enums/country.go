package enums

type Country string

const (
	Argentina           Country = "argentina"
	Bolivia             Country = "bolivia"
	Brasil              Country = "brasil"
	Chile               Country = "chile"
	Colombia            Country = "colombia"
	CostaRica           Country = "costa rica"
	Cuba                Country = "cuba"
	Ecuador             Country = "ecuador"
	ElSalvador          Country = "el salvador"
	Guatemala           Country = "guatemala"
	Haiti               Country = "haiti"
	Honduras            Country = "honduras"
	Mexico              Country = "mexico"
	Nicaragua           Country = "nicaragua"
	Panama              Country = "panama"
	Paraguay            Country = "paraguay"
	Peru                Country = "peru"
	PuertoRico          Country = "puerto rico"
	RepublicaDominicana Country = "republica dominicana"
	Spain               Country = "spain"
	Uruguay             Country = "uruguay"
	Venezuela           Country = "venezuela"
)

var ValidCountries = map[Country]struct{}{
	Argentina:           {},
	Bolivia:             {},
	Brasil:              {},
	Chile:               {},
	Colombia:            {},
	CostaRica:           {},
	Cuba:                {},
	Ecuador:             {},
	ElSalvador:          {},
	Guatemala:           {},
	Haiti:               {},
	Honduras:            {},
	Mexico:              {},
	Nicaragua:           {},
	Panama:              {},
	Paraguay:            {},
	Peru:                {},
	PuertoRico:          {},
	RepublicaDominicana: {},
	Spain:               {},
	Uruguay:             {},
	Venezuela:           {},
}

func IsValidCountry(country string) bool {
	_, ok := ValidCountries[Country(country)]
	return ok
}

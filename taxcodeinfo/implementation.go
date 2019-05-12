package taxcodeinfo

// this file implements the functions for calculating tax
func entertainmentTax(price float64) float64 {
	if price < 100 && price > 0 {
		return 0
	}
	if price >= 100 {
		return float64(0.01) * (price - 100)
	}

	return -1
}

func tobaccoTax(price float64) float64 {
	return 10 + float64(0.02)*price
}

func fnbTax(price float64) float64 {
	return price * float64(0.1)
}

// Map stores the statically assigned (just like what is specified in the problem statement) information for tax codes. Key of the map = tax code.
var Map = map[int]Fetcher{
	1: NewFetcher("Food and Beverage", true, fnbTax),
	2: NewFetcher("Tobacco", false, tobaccoTax),
	3: NewFetcher("Entertainment", false, entertainmentTax),
}

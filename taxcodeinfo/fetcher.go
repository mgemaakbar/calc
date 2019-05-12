package taxcodeinfo

// Fetcher is an abstraction which returns the information of a tax code.
type Fetcher interface {
	CalculateTaxAmount(itemPrice float64) float64
	IsItemRefundable() bool
	TypeName() string
}

type fetcher struct {
	name              string
	refundable        bool
	taxCalculatorFunc func(float64) float64
}

func NewFetcher(name string, refundable bool, taxCalculatorFunc func(float64) float64) Fetcher {
	return &fetcher{name: name, refundable: refundable, taxCalculatorFunc: taxCalculatorFunc}
}

// CalculateTaxAmount returns the tax amount of the given item price and tax code.
func (t *fetcher) CalculateTaxAmount(itemPrice float64) float64 {
	return t.taxCalculatorFunc(itemPrice)
}

// IsItemRefundable returns 'true' if the item is refundable.
func (t *fetcher) IsItemRefundable() bool {
	return t.refundable
}

// Name is a method that returns the item type name.
func (t *fetcher) TypeName() string {
	return t.name
}

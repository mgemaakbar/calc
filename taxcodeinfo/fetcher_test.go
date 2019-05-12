package taxcodeinfo

import "testing"

func TestFetcher(t *testing.T) {
	f := NewFetcher("name", true, func(price float64) float64 {
		return 1
	})
	if f.CalculateTaxAmount(100) != float64(1) {
		t.Errorf("Expected 1 got %f", f.CalculateTaxAmount(100))
	}
	if f.TypeName() != "name" {
		t.Errorf("Expected 'name' got %s", f.TypeName())
	}
	if !f.IsItemRefundable() {
		t.Error("Should return true")
	}
}

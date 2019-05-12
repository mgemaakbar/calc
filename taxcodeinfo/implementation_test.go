package taxcodeinfo

import "testing"

func TestEntertainmentTax(t *testing.T) {
	if entertainmentTax(150) != float64(0.5) {
		t.Errorf("Expected 10.0 got %f", entertainmentTax(1100))
	}
	if entertainmentTax(10) != float64(0) {
		t.Errorf("Expected 0 got %f", entertainmentTax(10))
	}
	if entertainmentTax(-1) != float64(-1) {
		t.Errorf("Expected -1 got %f", entertainmentTax(-1))
	}
}

func TestTobaccoTax(t *testing.T) {
	if tobaccoTax(100) != 12 {
		t.Errorf("Expected 12 got %f", tobaccoTax(100))
	}
}

func TestFnbTax(t *testing.T) {
	if fnbTax(100) != 10 {
		t.Errorf("Expected 20 got %f", fnbTax(100))
	}
}

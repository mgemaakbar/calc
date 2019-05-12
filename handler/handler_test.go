package handler_test

import (
	"calc/handler"
	"calc/model"
	"calc/taxcodeinfo"
	"errors"
	"testing"
)

type mockQuerierReturnError struct {
	model.Querier
}

func (m *mockQuerierReturnError) GetItems() ([]*model.Item, error) {
	return nil, errors.New(`error`)
}

type mockQuerierReturnItemTaxCodeInvalid struct {
	model.Querier
}

func (m *mockQuerierReturnItemTaxCodeInvalid) GetItems() ([]*model.Item, error) {
	return []*model.Item{
		&model.Item{
			Name:    `name1`,
			TaxCode: 123,
			Price:   1000,
		},
	}, nil
}

type mockQuerierReturnValidItems struct {
	model.Querier
}

func (m *mockQuerierReturnValidItems) GetItems() ([]*model.Item, error) {
	return []*model.Item{
		&model.Item{
			Name:    `name1`,
			TaxCode: 1,
			Price:   1000,
		},
	}, nil
}

type mockFetcher struct {
	taxcodeinfo.Fetcher
}

func (m *mockFetcher) TypeName() string {
	return "typeName"
}

func (m *mockFetcher) IsItemRefundable() bool {
	return true
}

func (m *mockFetcher) CalculateTaxAmount(price float64) float64 {
	return 123
}

func TestGetBill_QueryError(t *testing.T) {
	h := handler.NewHandler(&mockQuerierReturnError{}, nil)

	res, err := h.GetBill()
	if res != nil {
		t.Error("Result should be nil")
	}
	if err.Error() != `error` {
		t.Error("An error should be returned, with string 'error'")
	}
}

func TestGetBill_TaxCodeNotFoundInMap(t *testing.T) {
	h := handler.NewHandler(&mockQuerierReturnItemTaxCodeInvalid{}, map[int]taxcodeinfo.Fetcher{
		1: nil,
	})

	res, err := h.GetBill()
	if res != nil {
		t.Error("Result should be nil")
	}
	if err.Error() != `tax code 123 not found` {
		t.Error("An error should be returned, with string 'tax code 123 not found'")
	}
}

func TestGetBill_noErrors(t *testing.T) {
	h := handler.NewHandler(&mockQuerierReturnValidItems{}, map[int]taxcodeinfo.Fetcher{
		1: &mockFetcher{},
	})

	res, err := h.GetBill()
	if res == nil {
		t.Error("Result should not be nil")
	}
	if err != nil {
		t.Error("Error should be nil")
	}
	if len(res) != 1 {
		t.Error("bills length should be 1")
	}

	if res[0].Name != "name1" {
		t.Errorf("Expected name1 got %s", res[0].Name)
	}
	if res[0].TaxCode != 1 {
		t.Errorf("Expected 1 got %d", res[0].TaxCode)
	}
	if !res[0].Refundable {
		t.Error("Refundable should be true")
	}
	if res[0].Tax != float64(123) {
		t.Errorf("Expected %f got %f", float64(123), res[0].Tax)
	}
	if res[0].Price != float64(1000) {
		t.Errorf("Expected 1000 got %f", res[0].Price)
	}
	if res[0].Amount != float64(1123) {
		t.Errorf("Expected 1123 got %f", res[0].Amount)
	}
}

type mockCreateItemReturnError struct {
	model.Querier
}

func (m *mockCreateItemReturnError) CreateItem(name string, taxCode int, price float64) error {
	return errors.New(`error`)
}

func TestCreateItem(t *testing.T) {
	h := handler.NewHandler(&mockCreateItemReturnError{}, nil)
	err := h.CreateItem("name", 1, 100)
	if err == nil {
		t.Error(`Expecting an error with 'error' string to be returned`)
	}
}

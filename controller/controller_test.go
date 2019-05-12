package controller_test

import (
	"calc/controller"
	"calc/handler"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockHandlerGetBillError struct {
	handler.Handler
}

func (h *mockHandlerGetBillError) GetBill() ([]*handler.Bill, error) {
	return nil, errors.New(`error`)
}

func TestController_500status(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(controller.NewController(&mockHandlerGetBillError{}).GetBill)

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	if rr.Body.String() != `{"message":"error"}` {
		t.Errorf(`Expected {"message":"error"} got %s`, rr.Body.String())
	}
}

type mockHandlerGetBillNoError struct {
	handler.Handler
}

func (h *mockHandlerGetBillNoError) GetBill() ([]*handler.Bill, error) {
	return []*handler.Bill{
		&handler.Bill{},
	}, nil
}

func TestController_200Status(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(controller.NewController(&mockHandlerGetBillNoError{}).GetBill)

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != `{"bills":[{"name":"","tax_code":0,"type":"","refundable":false,"price":0,"tax":0,"amount":0}]}` {
		t.Errorf(`Expected {"bills":[{"name":"","tax_code":0,"type":"","refundable":false,"price":0,"tax":0,"amount":0}]} got %s`, rr.Body.String())
	}
}

type mockHandlerCreateNoError struct {
	handler.Handler
}

func (m *mockHandlerCreateNoError) CreateItem(name string, taxCode int, price float64) error {
	return nil
}

func TestCreate_200Status(t *testing.T) {
	req, err := http.NewRequest("POST", "/", strings.NewReader(`{"name": "nama", "price": 100, "tax_code":1}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(controller.NewController(&mockHandlerCreateNoError{}).Create)

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

type mockHandlerCreateError struct {
	handler.Handler
}

func (m *mockHandlerCreateError) CreateItem(name string, taxCode int, price float64) error {
	return errors.New(`error`)
}

func TestCreate_500Status(t *testing.T) {
	req, err := http.NewRequest("POST", "/", strings.NewReader(`{"name": "nama", "price": 100, "tax_code":1}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(controller.NewController(&mockHandlerCreateError{}).Create)

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	if rr.Body.String() != `{"message":"error"}` {
		t.Errorf(`Expected {"message":"error"} got %s`, rr.Body.String())
	}
}

func TestCreate_400Status(t *testing.T) {
	req, err := http.NewRequest("POST", "/", strings.NewReader(`{"name": 400, "price": 100, "tax_code":1}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(controller.NewController(&mockHandlerCreateNoError{}).Create)

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != 400 {
		t.Errorf("returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), "cannot unmarshal") {
		t.Errorf(`Expected the string to contain 'cannot unmarshal'`)
	}
}

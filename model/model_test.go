package model_test

import (
	"calc/model"
	"errors"
	"fmt"
	"regexp"
	"testing"

	_ "github.com/lib/pq"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func FixedFullRe(s string) string {
	return fmt.Sprintf("^%s$", regexp.QuoteMeta(s))
}

func TestCreateItem_error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	h := model.NewQuerier(db)

	mock.ExpectExec(FixedFullRe(`INSERT INTO item ("name", "tax_code", "price") VALUES ($1,$2,$3)`)).WillReturnError(errors.New(`error`))

	err := h.CreateItem(`name`, 1, 100)

	if err == nil {
		t.Error("Error should not be nil")
	}
	if err.Error() != `error` {
		t.Error("Error string should be 'error'")
	}

}

func TestCreateItem_noError(t *testing.T) {
	db, mock, _ := sqlmock.New()

	defer db.Close()

	h := model.NewQuerier(db)

	mock.ExpectExec(FixedFullRe(`INSERT INTO item ("name", "tax_code", "price") VALUES ($1,$2,$3)`)).WillReturnResult(sqlmock.NewResult(1, 2))

	err := h.CreateItem(`name`, 1, 100)

	if err != nil {
		t.Error("Error should be nil")
	}
}

func TestGetItems_queryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	q := model.NewQuerier(db)

	mock.ExpectQuery(FixedFullRe(`select name,tax_code,price from item`)).WillReturnError(errors.New(`error`))

	res, err := q.GetItems()

	if res != nil {
		t.Error(`bill should be nil`)
	}

	if err == nil {
		t.Error("Error should not be nil")
	}

	if err.Error() != `error` {
		t.Error("Error string should be 'error'")
	}
}

func TestGetItems_scanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	q := model.NewQuerier(db)

	rows := sqlmock.NewRows([]string{
		"name",
		"tax_code",
		"price",
	}).AddRow(
		"123",
		"salah",
		float64(123.0),
	)

	mock.ExpectQuery(FixedFullRe(`select name,tax_code,price from item`)).WillReturnRows(rows)

	res, err := q.GetItems()

	if res != nil {
		t.Error(`bill should be nil`)
	}

	if err == nil {
		t.Error("Error should not be nil")
	}
}

func TestGetItems_allFine(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	q := model.NewQuerier(db)

	rows := sqlmock.NewRows([]string{
		"name",
		"tax_code",
		"price",
	}).AddRow(
		"name",
		1,
		float64(123.0),
	)

	mock.ExpectQuery(FixedFullRe(`select name,tax_code,price from item`)).WillReturnRows(rows)

	res, err := q.GetItems()

	if err != nil {
		t.Error("Error should be nil")
	}

	if res == nil {
		t.Error(`bill should not be nil`)
	}

	if res[0].Name != "name" {
		t.Errorf("Expecting 'name got %s", res[0].Name)
	}
	if res[0].TaxCode != 1 {
		t.Errorf("Expecting 1 got %d", res[0].TaxCode)
	}
	if res[0].Price != float64(123) {
		t.Errorf("Expecting 123 got %f", res[0].Price)
	}

}

package model

import "database/sql"

type Querier interface {
	GetItems() ([]*Item, error)
	CreateItem(name string, taxCode int, price float64) error
}

type querier struct {
	db *sql.DB
}

func NewQuerier(db *sql.DB) Querier {
	return &querier{db: db}
}

func (q *querier) CreateItem(name string, taxCode int, price float64) error {
	_, err := q.db.Exec(`INSERT INTO item ("name", "tax_code", "price") VALUES ($1,$2,$3)`, name, taxCode, price)
	if err != nil {
		return err
	}
	return nil
}

func (q *querier) GetItems() ([]*Item, error) {
	query := `select name,tax_code,price from item`

	rows, err := q.db.Query(query)
	if err != nil {
		return nil, err
	}

	ret := []*Item{}

	for rows.Next() {
		item := Item{}
		err := rows.Scan(&item.Name, &item.TaxCode, &item.Price)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &item)
	}

	return ret, nil
}

type Item struct {
	Name    string
	TaxCode int
	Price   float64
}

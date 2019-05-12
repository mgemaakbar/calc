package handler

import (
	"calc/model"
	"calc/taxcodeinfo"
	"fmt"
)

type handler struct {
	querier               model.Querier
	taxCodeInfoFetcherMap map[int]taxcodeinfo.Fetcher
}

type Handler interface {
	GetBill() ([]*Bill, error)
	CreateItem(name string, taxCode int, price float64) error
}

type Bill struct {
	Name       string  `json:"name"`
	TaxCode    int     `json:"tax_code"`
	Type       string  `json:"type"`
	Refundable bool    `json:"refundable"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	Amount     float64 `json:"amount"`
}

func NewHandler(querier model.Querier, taxCodeInfoFetcherMap map[int]taxcodeinfo.Fetcher) Handler {
	return &handler{querier: querier, taxCodeInfoFetcherMap: taxCodeInfoFetcherMap}
}

func (h *handler) CreateItem(name string, taxCode int, price float64) error {
	return h.querier.CreateItem(name, taxCode, price)
}

func (h *handler) GetBill() ([]*Bill, error) {
	items, err := h.querier.GetItems()
	if err != nil {
		return nil, err
	}

	ret := []*Bill{}

	for _, item := range items {
		var ok bool
		var taxCodeInfoFetcher taxcodeinfo.Fetcher

		if taxCodeInfoFetcher, ok = h.taxCodeInfoFetcherMap[item.TaxCode]; !ok {
			return nil, fmt.Errorf(`tax code %d not found`, item.TaxCode)
		}

		typeName := taxCodeInfoFetcher.TypeName()
		refundable := taxCodeInfoFetcher.IsItemRefundable()
		tax := taxCodeInfoFetcher.CalculateTaxAmount(item.Price)

		ret = append(ret, &Bill{
			Name:       item.Name,
			TaxCode:    item.TaxCode,
			Type:       typeName,
			Refundable: refundable,
			Price:      item.Price,
			Tax:        tax,
			Amount:     item.Price + tax,
		})
	}

	return ret, nil
}

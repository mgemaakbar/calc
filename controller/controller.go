package controller

import (
	"calc/handler"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type controller struct {
	handler handler.Handler
}

func NewController(handler handler.Handler) *controller {
	return &controller{handler: handler}
}

func (c *controller) GetBill(w http.ResponseWriter, r *http.Request) {
	bills, err := c.handler.GetBill()

	if err != nil {
		resp := map[string]interface{}{
			"message": err.Error(),
		}

		respByte, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respByte)
		return
	}

	resp := map[string]interface{}{
		"bills": bills,
	}

	respByte, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respByte)
}

type CreateRequestBody struct {
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	TaxCode int     `json:"tax_code"`
}

func (c *controller) Create(w http.ResponseWriter, r *http.Request) {
	var reqBody CreateRequestBody
	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request

	defer r.Body.Close()

	if err := json.Unmarshal(body, &reqBody); err != nil {
		resp := map[string]interface{}{
			"message": err.Error(),
		}
		respByte, _ := json.Marshal(resp)
		w.WriteHeader(400)
		w.Write(respByte)
		return
	}

	err := c.handler.CreateItem(reqBody.Name, reqBody.TaxCode, reqBody.Price)
	if err != nil {
		resp := map[string]interface{}{
			"message": err.Error(),
		}
		respByte, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respByte)
		return
	}

	w.WriteHeader(http.StatusOK)
}

package go_certcentral

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const productURL = "/product"

type Product struct {
	GroupName             string `json:"group_name,omitempty"`
	NameID                string `json:"name_id,omitempty"`
	Name                  string `json:"name,omitempty"`
	Type                  string `json:"type,omitempty"`
	ValidationType        string `json:"validation_type,omitempty"`
	ValidationName        string `json:"validation_name,omitempty"`
	ValidationDescription string `json:"validation_description,omitempty"`
}

func (c *Client) ListProducts() ([]Product, error) {
	req, err := http.NewRequest(http.MethodGet, makeURL(productURL), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	type data struct {
		Products []Product `json:"products"`
	}

	var d data
	err = json.Unmarshal(resBody, &d)
	return d.Products, err
}

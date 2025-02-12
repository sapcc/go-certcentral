package go_certcentral

import (
	"encoding/json"
	"io"
	"net/http"
)

const productURL = "/product"

type ProductRsp struct {
	Products []Product `json:"products"`
}

type Product struct {
	GroupName                 string             `json:"group_name,omitempty"`
	NameID                    string             `json:"name_id,omitempty"`
	Name                      string             `json:"name,omitempty"`
	Type                      string             `json:"type,omitempty"`
	ValidationType            string             `json:"validation_type,omitempty"`
	ValidationName            string             `json:"validation_name,omitempty"`
	ValidationDescription     string             `json:"validation_description,omitempty"`
	AllowedContainerIDs       []int              `json:"allowed_container_ids,omitempty"`
	AllowedValidityYears      []int              `json:"allowed_validity_years,omitempty"`
	AllowedOrderValidityYears []int              `json:"allowed_order_validity_years,omitempty"`
	SignatureHashTypes        SignatureHashTypes `json:"signature_hash_types,omitempty"`
	AllowedCACerts            []AllowedCACert    `json:"allowed_ca_certs,omitempty"`
	CSRRequired               bool               `json:"csr_required,omitempty"`
	ReplacementProductNameID  string             `json:"replacement_product_name_id,omitempty"`
}

type SignatureHashTypes struct {
	AllowedHashTypes  []AllowedHashType `json:"allowed_hash_types,omitempty"`
	DefaultHashTypeID string            `json:"default_hash_type_id,omitempty"`
}

type AllowedHashType struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type AllowedCACert struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
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

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var productRsp ProductRsp
	err = json.Unmarshal(resBody, &productRsp)
	return productRsp.Products, err
}

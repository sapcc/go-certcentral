package go_certcentral

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const domainURL = "/domain"

type Domain struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	IsActive     bool          `json:"is_active,omitempty"`
	DateCreated  *time.Time    `json:"date_created,omitempty"`
	Organization *Organization `json:"organization,omitempty"`
	Validations  []Validation  `json:"validations,omitempty"`
	DCV          *DCV          `json:"dcv,omitempty"`
	Container    *Container    `json:"container,omitempty"`
}

func (c *Client) ListDomains(containerID string) ([]Domain, error) {
	if containerID == "" {
		return nil, errors.New("cannot list domains without container ID")
	}

	req, err := http.NewRequest(http.MethodGet, makeURL(domainURL, "?container_id", containerID), nil)
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
		Domains []Domain `json:"domains"`
	}

	var d data
	err = json.Unmarshal(resBody, &d)
	return d.Domains, err
}

func (c *Client) GetDomain(domainID string) (*Domain, error) {
	req, err := http.NewRequest(
		http.MethodGet, makeURL(domainURL, domainID, "include_dcv=true&include_validation=true"), nil,
	)
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

	var d Domain
	err = json.Unmarshal(resBody, &d)
	return &d, err
}

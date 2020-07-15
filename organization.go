package go_certcentral

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const organizationURL = "/organization"

type Organization struct {
	ID                  int          `json:"id"`
	Status              string       `json:"status"`
	Name                string       `json:"name"`
	AssumedName         string       `json:"assumed_name,omitempty"`
	DisplayName         string       `json:"display_name,omitempty"`
	IsActive            bool         `json:"is_active,omitempty"`
	Address             string       `json:"address"`
	Address2            string       `json:"address2,omitempty"`
	City                string       `json:"city"`
	State               string       `json:"state"`
	Zip                 string       `json:"zip"`
	Country             string       `json:"country"`
	Telephone           string       `json:"telephone,omitempty"`
	Container           *Container   `json:"container,omitempty"`
	Validations         []Validation `json:"validations,omitempty"`
	EvApprovers         []User       `json:"ev_approvers,omitempty"`
	OrganizationContact *User        `json:"organization_contact,omitempty"`
	Contacts            []User       `json:"contacts,omitempty"`
}

func (c *Client) ListOrganizations() ([]Organization, error) {
	req, err := http.NewRequest(http.MethodGet, makeURL(organizationURL), nil)
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
		Organization []Organization `json:"organizations"`
	}

	var d data
	err = json.Unmarshal(resBody, &d)
	return d.Organization, err
}

func (c *Client) GetOrganizationByName(organizationName string) (*Organization, error) {
	orgList, err := c.ListOrganizations()
	if err != nil {
		return nil, err
	}

	organizationName = strings.ToLower(organizationName)
	for _, org := range orgList {
		if strings.ToLower(org.Name) == organizationName {
			return &org, nil
		}
	}

	return nil, fmt.Errorf("no organization found for name: %s", organizationName)
}

func (c *Client) GetOrganization(organizationID string) (*Organization, error) {
	req, err := http.NewRequest(http.MethodGet, makeURL(organizationURL, organizationID), nil)
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

	var org Organization
	err = json.Unmarshal(resBody, &org)
	return &org, err
}

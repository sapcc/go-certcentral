package go_certcentral

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const orderURL = "/order/certificate/"

type OrderType string

func (o OrderType) String() string {
	return string(o)
}

var OrderTypes = struct {
	SSLPlus,
	SSLMultiDomain,
	SSLWildcard,
	PrivateSSLPlus,
	PrivateSSLWildcard,
	SecureSiteProSSL,
	SecureSiteOV,
	SecureSiteProEVSSL,
	SecureSiteEV OrderType
}{
	"ssl_plus",
	"ssl_multi_domain",
	"ssl_wildcard",
	"private_ssl_plus",
	"private_ssl_wildcard",
	"ssl_securesite_pro",
	"ssl_securesite_flex",
	"ssl_ev_securesite_pro",
	"ssl_ev_securesite_flex",
}

type PaymentMethod string

func (p PaymentMethod) String() string {
	return string(p)
}

var PaymentMethods = struct {
	Balance,
	Card,
	Profile,
	WireTransfer PaymentMethod
}{
	"balance",
	"card",
	"profile",
	"wire_transfer",
}

type (
	Order struct {
		Certificate                 Certificate        `json:"certificate,omitempty"`
		Organization                *Organization      `json:"organization,omitempty"`
		OrderValidity               OrderValidity      `json:"order_validity,omitempty"`
		CustomExpirationDate        string             `json:"custom_expiration_date,omitempty"`
		Comments                    string             `json:"comments,omitempty"`
		ProcessorComment            string             `json:"processor_comment,omitempty"`
		DisableRenewalNotifications bool               `json:"disable_renewal_notifications,omitempty"`
		RenewalOfOrderID            int                `json:"renewal_of_order_id,omitempty"`
		PaymentMethod               PaymentMethod      `json:"payment_method,omitempty"`
		SkipApproval                bool               `json:"skip_approval,omitempty"`
		Product                     *Product           `json:"product,omitempty"`
		OrganizationContact         *User              `json:"organization_contact,omitempty"`
		TechnicalContact            *User              `json:"technical_contact,omitempty"`
		User                        *User              `json:"user,omitempty"`
		CsProvisioningMethod        string             `json:"cs_provisioning_method,omitempty"`
		DisableCT                   bool               `json:"disable_ct,omitempty"`
		Requests                    []OrderRequest     `json:"requests,omitempty"`
		ID                          int                `json:"id"`
		Domains                     []Domain           `json:"domains,omitempty"`
		CertificateID               int                `json:"certificate_id,omitempty"`
		CertificateChain            []CertificateChain `json:"certificate_chain,omitempty"`
		Container                   *Container         `json:"container,omitempty"`
	}

	OrderRequest struct {
		ID       int        `json:"id"`
		Date     *time.Time `json:"date,omitempty"`
		Type     string     `json:"type,omitempty"`
		Status   Status     `json:"status,omitempty"`
		Comments string     `json:"comments,omitempty"`
	}

	OrderValidity struct {
		Years int `json:"years,omitempty"`
		Days  int `json:"days,omitempty"`
	}
)

func (o Order) DecodeCertificateChain() ([]*x509.Certificate, error) {
	crtChain := make([]*x509.Certificate, 0)
	for _, crt := range o.CertificateChain {
		decodedCrt, err := crt.DecodePEM()
		if err != nil {
			return crtChain, err
		}
		crtChain = append(crtChain, decodedCrt...)
	}
	return crtChain, nil
}

func (c *Client) SubmitOrder(order Order, orderType OrderType) (*Order, error) {
	body, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, makeURL(orderURL, orderType.String()), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var orderRes Order
	err = json.Unmarshal(body, &orderRes)
	return &orderRes, err
}

func (c *Client) GetOrder(orderID string) (*Order, error) {
	if orderID == "" {
		return nil, errors.New("cannot get order without ID")
	}

	req, err := http.NewRequest(http.MethodGet, makeURL(orderURL, orderID), nil)
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

	var o Order
	err = json.Unmarshal(resBody, &o)
	return &o, err
}

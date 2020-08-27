package go_certcentral

import (
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const certificateURL = "/certificate"

type CertificateFormat string

func (cf CertificateFormat) String() string {
	return string(cf)
}

// CertificateFormats is the set of formats for a certificate.
// Additional documentation can be found here:
// https://dev.digicert.com/glossary/#certificate-formats
var CertificateFormats = struct {
	Default,
	PEM,
	DefaultPEM,
	PEMAll,
	PEMNoIntermediate,
	PEMNoRoot,
	P7B,
	CRT CertificateFormat
}{
	"default",
	"pem",
	"default_pem",
	"pem_all",
	"pem_nointermediate",
	"pem_noroot",
	"p7b",
	"crt",
}

type SignatureHash string

func (s SignatureHash) String() string {
	return string(s)
}

var SignatureHashes = struct {
	SHA256,
	SHA384,
	SHA512,
	SHA1 SignatureHash
}{
	"sha256",
	"sha384",
	"sha512",
	"sha1",
}

type (
	Certificate struct {
		CommonName        string         `json:"common_name"`
		DNSNames          []string       `json:"dns_names"`
		CSR               string         `json:"csr,omitempty"`
		ServerPlatform    ServerPlatform `json:"server_platform"`
		SignatureHash     SignatureHash  `json:"signature_hash"`
		CaCertID          string         `json:"ca_cert_id,omitempty"`
		OrganizationUnits []string       `json:"organization_units,omitempty"`
		Organization      *Organization  `json:"organization,omitempty"`
		ProfileOption     string         `json:"profile_option,omitempty"`
		ID                int            `json:"id,omitempty"`
		Thumbprint        string         `json:"thumbprint,omitempty"`
		SerialNumber      string         `json:"serial_number,omitempty"`
		DateCreated       *time.Time     `json:"date_created,omitempty"`
		ValidFrom         string         `json:"valid_from,omitempty"`
		ValidTill         string         `json:"valid_till,omitempty"`
		KeySize           int            `json:"key_size,omitempty"`
		CACert            *CACert        `json:"ca_cert,omitempty"`
	}

	CACert struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	CertificateChain struct {
		SubjectCommonName string `json:"subject_common_name"`
		Pem               string `json:"pem"`
	}

	CertificateRevokeResponse struct {
		ID        int        `json:"id"`
		Date      *time.Time `json:"date,omitempty"`
		Type      string     `json:"type,omitempty"`
		Status    Status     `json:"status,omitempty"`
		Requester *User      `json:"requester,omitempty"`
		Comments  string     `json:"comments,omitempty"`
	}
)

func (cc CertificateChain) DecodePEM() ([]*x509.Certificate, error) {
	return decodePEM([]byte(cc.Pem))
}

func (c *Client) DownloadCertificateForOrder(orderID string, certFormat CertificateFormat) ([]*x509.Certificate, error) {
	if orderID == "" {
		return nil, errors.New("cannot download certificate without orderID")
	}

	req, err := http.NewRequest(
		http.MethodGet, makeURL(certificateURL, "download/order", orderID, "format", certFormat.String()), nil,
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

	return decodePEM(resBody)
}

func (c *Client) DownloadCertificate(certificateID string, certFormat CertificateFormat) ([]*x509.Certificate, error) {
	if certificateID == "" {
		return nil, errors.New("cannot download certificate without its ID")
	}

	req, err := http.NewRequest(http.MethodGet, makeURL(certificateURL, certificateID, "download/format", certFormat.String()), nil)
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

	return decodePEM(resBody)
}

func (c *Client) RevokeCertificate(certificateID string) (*CertificateRevokeResponse, error) {
	if certificateID == "" {
		return nil, errors.New("cannot revoke certificate without its ID")
	}

	req, err := http.NewRequest(
		http.MethodPut, makeURL(certificateURL, certificateID, "revoke"),
		strings.NewReader("{\n  \"skip_approval\":true\n}"),
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

	var r CertificateRevokeResponse
	err = json.Unmarshal(resBody, &r)
	return &r, err
}

func (c *Client) GetCertificateChain(certID string) ([]CertificateChain, error) {
	if certID == "" {
		return nil, errors.New("cannot get certificate chain without certificate ID")
	}

	req, err := http.NewRequest(http.MethodGet, makeURL("certificate", certID, "chain"), nil)
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

	var chain []CertificateChain
	err = json.Unmarshal(resBody, &chain)
	return chain, err
}

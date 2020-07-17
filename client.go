package go_certcentral

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

const (
	baseURL = "https://www.digicert.com/services/v2"
	contentTypeJson = "application/json"
)

// Client is the client for the DigiCert cert-central API.
type Client struct {
	*Options
	httpClient *http.Client
}

func New(opts *Options) (*Client, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}

	tr := &http.Transport{
		DisableCompression: false,
		DisableKeepAlives:  false,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			InsecureSkipVerify: false,
		},
		Proxy: http.ProxyFromEnvironment,
	}

	return &Client{
		Options:    opts,
		httpClient: &http.Client{Transport: tr},
	}, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "sapcc/go-certcentral")
	req.Header.Set("X-DC-DEVKEY", c.Token)
	req.Header.Set("Accept", contentTypeJson)
	req.Header.Set("Content-Type", contentTypeJson)

	if c.IsDebug {
		if reqDump, err := httputil.DumpRequest(req, true); err != nil {
			fmt.Println("err: ", err)
		} else {
			fmt.Println("sending request: ", string(reqDump))
		}
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if c.IsDebug {
		if resDump, err := httputil.DumpResponse(res, true); err != nil {
			fmt.Println("err: ", err)
		} else {
			fmt.Println("got response: ", string(resDump))
		}
	}

	err = wrapErrorIfAny(res)
	return res, err
}

func wrapErrorIfAny(res *http.Response) error {
	if res.StatusCode >= 400 {
		if res.Header.Get("Content-Type") == contentTypeJson {
			return parseJsonError(res)
		}

		return &Error{
			Code:    res.StatusCode,
			Status:  res.Status,
			Message: "unknown error",
		}
	}

	return nil
}

func parseJsonError(res *http.Response) error {
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	type errs struct {
		Errors []Error `json:"errors"`
	}

	var e errs
	if err := json.Unmarshal(resBody, &e); err != nil {
		return err
	}

	resErr := &Error{
		Code:    res.StatusCode,
		Status:  res.Status,
		Message: "unkown",
	}

	if len(e.Errors) > 0 {
		resErr.Status = e.Errors[0].Status
		resErr.Message = e.Errors[0].Message
	}

	return resErr
}

package go_certcentral

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

const (
	separator       = "/"
	certificateType = "CERTIFICATE"
)

func makeURL(urlParts ...string) string {
	normalizedParts := normalizeURLParts(urlParts, separator)
	url := ensureSuffix(baseURL, separator)
	url += strings.Join(normalizedParts, separator)
	return url
}

func normalizeURLParts(urlParts []string, separator string) []string {
	normalizedParts := make([]string, len(urlParts))
	for idx, part := range urlParts {
		part = strings.TrimPrefix(part, separator)
		part = strings.TrimSuffix(part, separator)
		normalizedParts[idx] = part
	}
	return normalizedParts
}

func ensureSuffix(s, suffix string) string {
	if !strings.HasSuffix(s, suffix) {
		return s + suffix
	}
	return s
}

func decodePEM(data []byte) ([]*x509.Certificate, error) {
	crtList := make([]*x509.Certificate, 0)

	for {
		block, rest := pem.Decode(data)
		if block == nil {
			return crtList, errors.New("couldn't decode certificate from PEM block")
		}

		if block.Type != certificateType {
			return crtList, errors.New("certificate contains invalid date")
		}

		crt, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return crtList,  err
		}

		crtList = append(crtList, crt)
		data = rest

		if data == nil || len(data) == 0 {
			break
		}
	}

	return crtList, nil
}

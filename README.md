go-certcentral
--------------

[![REUSE status](https://api.reuse.software/badge/github.com/sapcc/go-certcentral)](https://api.reuse.software/info/github.com/sapcc/go-certcentral)

GoLang client for the [DigiCert cert-central services API](https://dev.digicert.com/services-api).


# Usage

```go
import cc "github.com/sapcc/go-certcentral"

client, err := cc.New(&cc.Options{
  Token: "DIGICERT_API_TOKEN",
  IsDebug: false,
})
handleError(err)

// Submit a certificate order.
orderResponse, err := cli.SubmitOrder(
  cc.Order{
    Certificate: cc.Certificate{
      CommonName:     csr.Subject.CommonName,
      DNSNames:       csr.DNSNames,
      CSR:            csr.PEM,
      ServerPlatform: cc.ServerPlatformForType(cc.ServerPlatformTypes.Nginx),
      SignatureHash:  cc.SignatureHashes.SHA256,
      CaCertID:       "CACertID",
      OrganizationUnits: []string{
        "SomeOrganization ",
      },
    },
    OrderValidity:               cc.OrderValidity{Years: 1},
    DisableRenewalNotifications: true,
    PaymentMethod:               cc.PaymentMethods.Balance,
    SkipApproval:                true,
    Organization:                &cc.Organization{ID: 123456},
}, cc.OrderTypes.PrivateSSLPlus)
handleError(err)

// If auto-approval is allowed the response contains the full chain of certificates in PEM format. 
if len(orderResponse.CertificateChain) > 0 {
  crtChain, err := orderResponse.DecodeCertificateChain()
  handleError(err)

  for _, crt := range crtChain {
    fmt.Println(crt.Subject.CommonName)
  }
}

// Download the certificate(s) for an order.
certList, err := client.DownloadCertificateForOrder("123456", cc.CertificateFormats.PEMAll)
handlerError(err)
for _, cert := range certList {
  fmt.Println(cert.Subject.CommonName)
}

```


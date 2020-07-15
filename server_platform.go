package go_certcentral

// ServerPlatform ...
type ServerPlatform struct {
	ID         int    `json:"id"`
	Name       string `json:"name,omitempty"`
	InstallURL string `json:"install_url,omitempty"`
	CsrURL     string `json:"csr_url,omitempty"`
}

func ServerPlatformForType(platform ServerPlatformType) ServerPlatform {
	return ServerPlatform{
		ID: platform.Int(),
	}
}

type ServerPlatformType int

func (s ServerPlatformType) Int() int {
	return int(s)
}

// ServerPlatformTypes is the collection of available ServerPlatformType.
// The full list of supported server platforms can be found here:
// https://www.digicert.com/services/v2/documentation/appendix-server-platforms
var ServerPlatformTypes = struct {
	Nginx,
	Other ServerPlatformType
}{
	45,
	57,
}

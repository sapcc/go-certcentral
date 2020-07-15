package go_certcentral

type Status string

func (o Status) String() string {
	return string(o)
}

var Stati = struct {
	Pending,
	Approved,
	Rejected Status
}{
	"pending",
	"approved",
	"rejected",
}

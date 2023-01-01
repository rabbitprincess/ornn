package ornn

import (
	"github.com/gokch/ornn/db"
)

func NewVendor(vendor db.Vendor) *Vendor {
	return &Vendor{dbVendor: vendor}
}

// vendor config by schema
type Vendor struct {
	dbVendor db.Vendor
}

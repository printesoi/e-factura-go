package efactura

import (
	"encoding/xml"

	"github.com/shopspring/decimal"
)

// Decimal is a small wrapper of the github.com/shopspring/decimal.Decimal
// type.
type Decimal struct {
	decimal.Decimal
}

func NewFromDecimal(d decimal.Decimal) Decimal {
	return Decimal{Decimal: d}
}

func (d Decimal) Ptr() *Decimal {
	return &d
}

func (d Decimal) String() string {
	return d.Decimal.StringFixed(2)
}

func (d Decimal) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(d.String(), start)
}

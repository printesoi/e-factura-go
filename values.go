package efactura

import (
	"encoding/xml"
)

type AmountWithCurrency struct {
	Amount Decimal
	// Term: Codul monedei
	CurrencyID CurrencyCodeType
}

func (a AmountWithCurrency) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type xmlAmount struct {
		Amount     string           `xml:",chardata"`
		CurrencyID CurrencyCodeType `xml:"currencyID,attr,omitempty"`
	}
	xa := xmlAmount{
		Amount:     a.Amount.String(),
		CurrencyID: a.CurrencyID,
	}
	return e.EncodeElement(xa, start)
}

type ValueWithAttrs struct {
	Value      string     `xml:",chardata"`
	Attributes []xml.Attr `xml:",any,attr,omitempty"`
}

func (v ValueWithAttrs) Ptr() *ValueWithAttrs {
	return &v
}

func MakeValueWithAttrs(value string, attrs ...xml.Attr) ValueWithAttrs {
	return ValueWithAttrs{
		Value:      value,
		Attributes: attrs,
	}
}

func MakeValueWithScheme(value string, schemeID string) ValueWithAttrs {
	return MakeValueWithAttrs(value, xml.Attr{
		Name:  xml.Name{Local: "schemeID"},
		Value: schemeID,
	})
}

func NewValueWithAttrs(value string, attrs ...xml.Attr) *ValueWithAttrs {
	return MakeValueWithAttrs(value, attrs...).Ptr()
}

package efactura

import (
	"encoding/xml"
)

type ValueWithScheme struct {
	Value string `xml:",chardata"`
	// Term: Identificatorul schemei
	// Cardinality: 0..1
	SchemeID string `xml:"schemeID,attr,omitempty"`
}

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

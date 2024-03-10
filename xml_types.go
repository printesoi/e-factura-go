// Copyright 2024 Victor Dodon
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package efactura

import (
	"time"

	"github.com/m29h/xml"
	"github.com/shopspring/decimal"
)

// Date is a wrapper of the time.Time type which marshals to XML in the
// YYYY-MM-DDD.
type Date struct {
	time.Time
}

func MakeDateLocal(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.Local)}
}

func MakeDateUTC(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

func (d Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := d.Format(time.DateOnly)
	return e.EncodeElement(v, start)
}

func (dt *Date) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var sd string
	if err := d.DecodeElement(&sd, &start); err != nil {
		return err
	}

	// TODO: always use Romanian time zone for date in efactura
	t, err := time.ParseInLocation(time.DateOnly, sd, time.Local)
	if err != nil {
		return err
	}

	*dt = Date{Time: t}
	return nil
}

func (d Date) Ptr() *Date {
	return &d
}

// IsInitialized if the Date is initialized (ie is created explicitly with a
// constructor or initialized by setting the Time, not implicitly via var
// declaration with no initialization).
func (d Date) IsInitialized() bool {
	return d != Date{}
}

// AmountWithCurrency represents an embeddable type that stores an amount as
// chardata and the currency ID as the currencyID attribute. The name of the
// node must be controlled by the parent type.
type AmountWithCurrency struct {
	Amount     Decimal
	CurrencyID CurrencyCodeType
}

// this type is a hack for a limitation of the encoding/xml package: it only
// supports []byte and string for a chardata.
type xmlAmountWithCurrency struct {
	Amount     string           `xml:",chardata"`
	CurrencyID CurrencyCodeType `xml:"currencyID,attr,omitempty"`
}

func (a AmountWithCurrency) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	xa := xmlAmountWithCurrency{
		Amount:     a.Amount.String(),
		CurrencyID: a.CurrencyID,
	}
	return e.EncodeElement(xa, start)
}

func (a *AmountWithCurrency) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var xa xmlAmountWithCurrency
	if err := d.DecodeElement(&xa, &start); err != nil {
		return err
	}

	amount, err := decimal.NewFromString(xa.Amount)
	if err != nil {
		return err
	}

	a.Amount = DD(amount)
	a.CurrencyID = xa.CurrencyID
	return nil
}

// ValueWithAttrs represents and embeddable type that stores a string as
// chardata and a list of attributes. The name of the XML node must be
// controlled by the parent type.
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

func (v *ValueWithAttrs) GetAttrByName(name string) (attr xml.Attr) {
	if v == nil {
		return
	}
	for _, a := range v.Attributes {
		if a.Name.Local == name {
			return a
		}
	}
	return
}

// InvoicedQuantity represents the quantity (of items) on an invoice line.
type InvoicedQuantity struct {
	Quantity Decimal
	// The unit of the quantity.
	UnitCode UnitCodeType
	// The quantity unit code list.
	UnitCodeListID string
	// The identification of the agency that maintains the quantity unit code
	// list.
	UnitCodeListAgencyID string
	// The name of the agency which maintains the quantity unit code list.
	UnitCodeListAgencyName string
}

// this type is a hack for a limitation of the encoding/xml package: it only
// supports []byte and string for a chardata.
type xmlInvoicedQuantity struct {
	Quantity               string       `xml:",chardata"`
	UnitCode               UnitCodeType `xml:"unitCode,attr"`
	UnitCodeListID         string       `xml:"unitCodeListID,attr,omitempty"`
	UnitCodeListAgencyID   string       `xml:"unitCodeListAgencyID,attr,omitempty"`
	UnitCodeListAgencyName string       `xml:"unitCodeListAgencyName,attr,omitempty"`
}

func (q InvoicedQuantity) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	xq := xmlInvoicedQuantity{
		Quantity:               q.Quantity.String(),
		UnitCode:               q.UnitCode,
		UnitCodeListID:         q.UnitCodeListID,
		UnitCodeListAgencyID:   q.UnitCodeListAgencyID,
		UnitCodeListAgencyName: q.UnitCodeListAgencyName,
	}
	return e.EncodeElement(xq, start)
}

func (q *InvoicedQuantity) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var xq xmlInvoicedQuantity
	if err := d.DecodeElement(&xq, &start); err != nil {
		return err
	}

	quantity, err := decimal.NewFromString(xq.Quantity)
	if err != nil {
		return err
	}

	q.Quantity = DD(quantity)
	q.UnitCode = xq.UnitCode
	q.UnitCodeListID = xq.UnitCodeListID
	q.UnitCodeListAgencyID = xq.UnitCodeListAgencyID
	q.UnitCodeListAgencyName = xq.UnitCodeListAgencyName
	return nil
}

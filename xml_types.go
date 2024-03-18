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

	"github.com/printesoi/xml-go"
	"github.com/shopspring/decimal"
)

var (
	// RoZoneLocation is the Romanian timezone location loaded in the init
	// function. This library does not load the time/tzdata package for the
	// embedded timezone database, so the user of this library is responsible
	// to ensure the Europe/Bucharest location is available.
	RoZoneLocation *time.Location

	// Allow mocking and testing
	timeNow = time.Now
)

func init() {
	if loc, err := time.LoadLocation("Europe/Bucharest"); err == nil {
		RoZoneLocation = loc
	} else {
		// If we could not load the Europe/Bucharest location, fallback to
		// time.UTC
		RoZoneLocation = time.UTC
	}
}

// Date is a wrapper of the time.Time type which marshals to XML in the
// YYYY-MM-DD format.
type Date struct {
	time.Time
}

// MakeDateLocal creates a date with the provided year, month and day in the
// Local time zone location.
func MakeDate(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, RoZoneLocation)}
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

	t, err := time.ParseInLocation(time.DateOnly, sd, RoZoneLocation)
	if err != nil {
		return err
	}

	*dt = Date{Time: t}
	return nil
}

// Ptr is a helper to return a *Date in contexts where a pointer is needed.
func (d Date) Ptr() *Date {
	return &d
}

// IsInitialized checks if the Date is initialized (ie is created explicitly
// with a constructor or initialized by setting the Time, not implicitly via
// var declaration with no initialization).
func (d Date) IsInitialized() bool {
	return d != Date{}
}

// Now returns time.Now() in Romanian zone location (Europe/Bucharest)
func Now() time.Time {
	return timeNow().In(RoZoneLocation)
}

// TimeInRomania returns the time t in Romanian zone location
// (Europe/Bucharest).
func TimeInRomania(t time.Time) time.Time {
	return t.In(RoZoneLocation)
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

// Ptr is a helper method to return a *ValueWithAttrs from the receiver in
// contexts where a pointer is needed.
func (v ValueWithAttrs) Ptr() *ValueWithAttrs {
	return &v
}

// MakeValueWithAttrs create a ValueWithAttrs using the provided chardata value
// and attributes.
func MakeValueWithAttrs(value string, attrs ...xml.Attr) ValueWithAttrs {
	return ValueWithAttrs{
		Value:      value,
		Attributes: attrs,
	}
}

// NewValueWithAttrs same as [MakeValueWithAttrs] but a pointer is returned.
func NewValueWithAttrs(value string, attrs ...xml.Attr) *ValueWithAttrs {
	return MakeValueWithAttrs(value, attrs...).Ptr()
}

// MakeValueWithScheme creates ValueWithAttrs with the provided chardata value
// and an attribute named `schemeID` with the given scheme ID.
func MakeValueWithScheme(value string, schemeID string) ValueWithAttrs {
	return MakeValueWithAttrs(value, xml.Attr{
		Name:  xml.Name{Local: "schemeID"},
		Value: schemeID,
	})
}

// GetAttrByName returns the attribute by local name. If no attribute with the
// given name exists, an empty xml.Attr is returned.
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

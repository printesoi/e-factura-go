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
	"encoding/xml"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Date is a wrapper of the time.Time type which marshals to XML in the
// YYYY-MM-DDD.
type Date struct {
	time.Time
}

func MakeDateLocal(year int, month time.Month, day int) Date {
	return Date{Time: time.Date(year, month, day, 0, 0, 0, 0, time.Local)}
}

func MakeDateUTC(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

func (d Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	dy, dm, dd := d.Time.Date()
	v := fmt.Sprintf("%04d-%02d-%02d", dy, dm, dd)
	return e.EncodeElement(v, start)
}

func (d Date) Ptr() *Date {
	return &d
}

// Decimal is a wrapper of the github.com/shopspring/decimal.Decimal type in
// order to ensure type safety and lossless computation.
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

// AmountWithCurrency represents an embeddable type that stores an amount as
// chardata and the currency ID as the currencyID attribute. The name of the
// node must be controlled by the parent type.
type AmountWithCurrency struct {
	Amount     Decimal
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

// ValueWithAttrs represents and embeddable type that stores a string as
// chardata and a list of attributes. The name of the XML name must be
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

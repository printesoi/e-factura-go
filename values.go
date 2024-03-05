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

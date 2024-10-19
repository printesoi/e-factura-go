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
	"github.com/printesoi/e-factura-go/pkg/types"
	"github.com/printesoi/xml-go"
)

// AmountWithCurrency represents an embeddable type that stores an amount as
// chardata and the currency ID as the currencyID attribute. The name of the
// node must be controlled by the parent type.
type AmountWithCurrency struct {
	Amount     types.Decimal    `xml:",chardata"`
	CurrencyID CurrencyCodeType `xml:"currencyID,attr,omitempty"`
}

// MarshalXML implements the xml.Marshaler interface. We use a custom
// marshaling function for AmountWithCurrency to ensure two digits after the
// decimal point.
func (a AmountWithCurrency) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type amountWithCurrency struct {
		Amount     string           `xml:",chardata"`
		CurrencyID CurrencyCodeType `xml:"currencyID,attr,omitempty"`
	}
	xmlAmount := amountWithCurrency{
		Amount:     a.Amount.StringFixed(2),
		CurrencyID: a.CurrencyID,
	}
	return e.EncodeElement(xmlAmount, start)
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

// IDNode is a struct that encodes a node that only has a cbc:ID property.
type IDNode struct {
	ID string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID"`
}

// MakeIDNode creates a IDNode with the given id.
func MakeIDNode(id string) IDNode {
	return IDNode{ID: id}
}

// NewIDNode creates a *IDNode with the given id.
func NewIDNode(id string) *IDNode {
	return &IDNode{ID: id}
}

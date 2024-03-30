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
)

var (
	// RoZoneLocation is the Romanian timezone location loaded in the init
	// function. This library does NOT load the time/tzdata package for the
	// embedded timezone database, so the user of this library is responsible
	// to ensure the Europe/Bucharest location is available, otherwise UTC is
	// used and may lead to unexpected results.
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
// YYYY-MM-DD format and is assumed to be in the Romanian timezone location.
type Date struct {
	time.Time
}

// MakeDate creates a date with the provided year, month and day in the
// Romanian time zone location.
func MakeDate(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, RoZoneLocation)}
}

// NewDate same as MakeDate, but returns a pointer to Date.
func NewDate(year int, month time.Month, day int) *Date {
	return MakeDate(year, month, day).Ptr()
}

// MakeDateFromTime creates a Date in Romanian time zone location from the
// given time.Time.
func MakeDateFromTime(t time.Time) Date {
	return MakeDate(t.In(RoZoneLocation).Date())
}

// NewDateFromTime same as MakeDateFromTime, but returns a pointer to Date.
func NewDateFromTime(t time.Time) *Date {
	return MakeDate(t.In(RoZoneLocation).Date()).Ptr()
}

// MarshalXML implements the xml.Marshaler interface.
func (d Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := d.Format(time.DateOnly)
	return e.EncodeElement(v, start)
}

// UnmarshalXML implements the xml.Unmarshaler interface.
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
	Amount     Decimal          `xml:",chardata"`
	CurrencyID CurrencyCodeType `xml:"currencyID,attr,omitempty"`
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

// IDNote is a struct that encodes a node that only has a cbc:ID property.
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

// MarshalXML returns the XML encoding of v in Canonical XML form [XML-C14N].
// This method must be used for marshaling objects from this library, instead
// of encoding/xml. This method does NOT include the XML header declaration.
func MarshalXML(v any) ([]byte, error) {
	return xml.Marshal(v)
}

// MarshalXMLWithHeader same as MarshalXML, but also add the XML header
// declaration.
func MarshalXMLWithHeader(v any) ([]byte, error) {
	data, err := MarshalXML(v)
	if err != nil {
		return nil, err
	}
	return concatBytes([]byte(xml.Header), data), nil
}

// MarshalIndentXML works like MarshalXML, but each XML element begins on a new
// indented line that starts with prefix and is followed by one or more
// copies of indent according to the nesting depth. This method does NOT
// include the XML header declaration.
func MarshalIndentXML(v any, prefix, indent string) ([]byte, error) {
	return xml.MarshalIndent(v, prefix, indent)
}

// MarshalIndentXMLWithHeader same as MarshalIndentXML, but also add the XML
// header declaration.
func MarshalIndentXMLWithHeader(v any, prefix, indent string) ([]byte, error) {
	data, err := MarshalIndentXML(v, prefix, indent)
	if err != nil {
		return nil, err
	}
	return concatBytes([]byte(xml.Header), data), nil
}

// Unmarshal parses the XML-encoded data and stores the result in
// the value pointed to by v, which must be an arbitrary struct,
// slice, or string. Well-formed data that does not fit into v is
// discarded. This method must be used for unmarshaling objects from this
// library, instead of encoding/xml.
func UnmarshalXML(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}

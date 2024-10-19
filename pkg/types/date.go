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

package types

import (
	"time"

	"github.com/printesoi/xml-go"

	itime "github.com/printesoi/e-factura-go/pkg/time"
)

// Date is a wrapper of the time.Time type which marshals to XML in the
// YYYY-MM-DD format and is assumed to be in the Romanian timezone location.
type Date struct {
	time.Time
}

// MakeDate creates a date with the provided year, month and day in the
// Romanian time zone location.
func MakeDate(year int, month time.Month, day int) Date {
	return Date{itime.Date(year, month, day, 0, 0, 0, 0)}
}

// NewDate same as MakeDate, but returns a pointer to Date.
func NewDate(year int, month time.Month, day int) *Date {
	return MakeDate(year, month, day).Ptr()
}

// MakeDateFromTime creates a Date in Romanian time zone location from the
// given time.Time.
func MakeDateFromTime(t time.Time) Date {
	return MakeDate(itime.TimeInRomania(t).Date())
}

// NewDateFromTime same as MakeDateFromTime, but returns a pointer to Date.
func NewDateFromTime(t time.Time) *Date {
	return MakeDate(itime.TimeInRomania(t).Date()).Ptr()
}

// MakeDateFromString creates a Date in Romanian time zone from a string in the
// YYYY-MM-DD format.
func MakeDateFromString(str string) (Date, error) {
	t, err := itime.ParseInRomania(time.DateOnly, str)
	if err != nil {
		return Date{}, err
	}
	return MakeDate(t.Date()), nil
}

// NewDateFromString same as MakeDateFromString, but returns a pointer to Date.
func NewDateFromString(str string) (*Date, error) {
	d, err := MakeDateFromString(str)
	if err != nil {
		return nil, err
	}
	return d.Ptr(), nil
}

// MarshalXML implements the xml.Marshaler interface.
func (d Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := d.Format(time.DateOnly)
	return e.EncodeElement(v, start)
}

// MarshalXMLAttr implements the xml.MarshalerAttr interface.
func (d Date) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	v := d.Format(time.DateOnly)
	return xml.Attr{
		Name:  name,
		Value: v,
	}, nil
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (dt *Date) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var sd string
	if err := d.DecodeElement(&sd, &start); err != nil {
		return err
	}

	t, err := itime.ParseInRomania(time.DateOnly, sd)
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

const (
	xsDateTimeFmt = "2006-01-02T15:04:05"
)

// DateTime is a wrapper of the time.Time type which marshals to XML in the
// YYYY-MM-DDTHH:MM:SS format and is assumed to be in the Romanian timezone
// location.
type DateTime struct {
	time.Time
}

// MakeDateTime creates a DateTime in RoZoneLocation.
func MakeDateTime(year int, month time.Month, day, hour, min, sec, nsec int) DateTime {
	return DateTime{Time: itime.Date(year, month, day, hour, min, sec, nsec)}
}

// MarshalXML implements the xml.Marshaler interface.
func (d DateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := d.Format(xsDateTimeFmt)
	return e.EncodeElement(v, start)
}

// MarshalXMLAttr implements the xml.MarshalerAttr interface.
func (d DateTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	v := d.Format(xsDateTimeFmt)
	return xml.Attr{
		Name:  name,
		Value: v,
	}, nil
}

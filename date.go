package efactura

import (
	"encoding/xml"
	"fmt"
	"time"
)

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

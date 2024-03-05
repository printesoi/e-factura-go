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

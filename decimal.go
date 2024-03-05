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

	"github.com/shopspring/decimal"
)

// Decimal is a small wrapper of the github.com/shopspring/decimal.Decimal
// type.
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

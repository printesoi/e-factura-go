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

package xml

import (
	"github.com/printesoi/xml-go"
)

// Constants for namespaces and versions
const (
	XMLNSUBLInvoice2 = "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2"
	XMLNSUBLcac      = "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2"
	XMLNSUBLcbc      = "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2"
)

// SetupUBLXMLEncoder will configure the xml.Encoder to make it suitable for
// marshaling UBL objects to XML.
func SetupUBLXMLEncoder(enc *xml.Encoder) *xml.Encoder {
	enc.AddNamespaceBinding(XMLNSUBLcac, "cac")
	enc.AddSkipNamespaceAttrForPrefix(XMLNSUBLcac, "cac")
	enc.AddNamespaceBinding(XMLNSUBLcbc, "cbc")
	enc.AddSkipNamespaceAttrForPrefix(XMLNSUBLcbc, "cbc")
	return enc
}

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
	xml "github.com/m29h/xml"
)

// Constants for namespaces and versions
const (
	XMLNSInvoice2 = "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2"
	XMLNSUBLcac   = "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2"
	XMLNSUBLcbc   = "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2"

	// e-factura: Customization ID implemented  CIUS-RO v1.0.1
	CIUSRO_v101 = "urn:cen.eu:en16931:2017#compliant#urn:efactura.mfinante.ro:CIUS-RO:1.0.1"
	// e-factura: UBL Version implemented
	UBLVersionID = "2.1"
)

func init() {
	xml.NameSpaceBinding.Add("urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2", "cac")
	xml.NameSpaceBinding.Add("urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2", "cbc")
}

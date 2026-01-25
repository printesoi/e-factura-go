// Copyright 2026 Victor Dodon
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

// Package vat provides a client for using the ANAF VAT v9 API for checking
// registered taxpayers for VAT purposes (Art. 316 of Fiscal Code), the
// register of taxpayers for cash-accounting VAT scheme, the register of
// inactive or reactivated taxpayers, the register for split VAT and RO
// e-factura register.
// https://static.anaf.ro/static/10/Anaf/Informatii_R/Servicii_web/doc_WS_V9.txt
//
// Usage:
//
//	import (
//		"log"
//		"github.com/printesoi/e-factura-go/pkg/vat"
//		"github.com/printesoi/e-factura-go/pkg/types"
//	)
//	func main() {
//		client, err := vat.NewClient()
//		// Query current status (current date in Romania)
//		res1, err := client.CheckVatV9(context.TODO(), MakeCheckVatRequest(CIF(12345678)))
//		if err != nil {
//			// Handle error
//		}
//		// Query status for multiple companies for specific dates
//		res2, err := client.CheckVatV9(context.TODO(), MakeCheckVatRequestFromItems(
//			vat.MakeCheckVatRequestItem(CIF(12345678), types.MakeDate(2025, 12, 31)),
//			vat.MakeCheckVatRequestItem(CIF(98765432), types.MakeDate(2025, 1, 1)),
//		))
//		if err != nil {
//			// Handle error
//		}
//		if len(res2.Found) == 0 {
//			// None of the provides CIFs are valid.
//		} else {
//			for _, data := range res2.Found {
//				if hasVat := data.HasVat(); hasVat {
//					lastVatPeriod, _ := data.GetLastVatPeriod()
//					log.Printf("CIF: %s, Name: %s, VAT enabled since %s\n",
//						data.GeneralData.CIF, data.GeneralData.Name,
//							lastVatPeriod.StartDate)
//				} else {
//					log.Printf("CIF: %s, Name: %s, no VAT\n",
//						data.GeneralData.CIF, data.GeneralData.Name)
//				}
//			}
//		}
//	}
package vat

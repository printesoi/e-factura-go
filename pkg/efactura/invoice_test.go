// Copyright 2024-2026 Victor Dodon
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
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/printesoi/e-factura-go/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type assersions struct {
	*assert.Assertions
	t *testing.T
}

func newAssert(t *testing.T) *assersions {
	return &assersions{
		Assertions: assert.New(t),
		t:          t,
	}
}

func assertBuild[T any](assert *assersions, builder Builder[T]) T {
	assert.t.Helper()
	result, err := builder.Build()
	require.NoError(assert.t, err)
	return result
}

func TestInvoiceRefBuild(t *testing.T) {
	assert := newAssert(t)

	ron := CurrencyRON

	invoiceNote := assertBuild(assert, NewInvoiceNoteBuilder("some text for invoice note"))
	invoicePeriod := assertBuild(assert, NewInvoicePeriodBuilder().WithEndDate(types.MakeDate(2022, 5, 31)))

	invoiceSupplier := InvoiceSupplierParty{
		CommercialName: &InvoicePartyName{
			Name: "ADMINISTRATIA SECTOR 2 A FINANTELOR PUBLICE",
		},
		PostalAddress: MakeInvoiceSupplierPostalAddress(PostalAddress{
			Line1:            "C. A. Rosetti, nr. 39",
			CityName:         CityNameROBSector2,
			CountrySubentity: CountrySubentityRO_B,
			PostalZone:       "013329",
			Country:          CountryRO,
		}),
		TaxScheme: &InvoicePartyTaxScheme{
			TaxScheme: TaxSchemeVAT,
			CompanyID: "RO4266367",
		},
		LegalEntity: InvoiceSupplierLegalEntity{
			Name:             "Seller SRL",
			CompanyLegalForm: "J40/12345/1998",
		},
		Contact: &InvoiceSupplierContact{
			Email: "mail@seller.com",
		},
	}

	invoiceCustomer := InvoiceCustomerParty{
		Identifications: []InvoicePartyIdentification{
			{
				ID: MakeValueWithAttrs("4340188"),
			},
		},
		CommercialName: &InvoicePartyName{
			Name: "Directia Generala Regionala a Finantelor Publice Bucuresti - Administratia Sector 3 a Finantelor Publice",
		},
		PostalAddress: MakeInvoiceCustomerPostalAddress(PostalAddress{
			Line1:            "BD DECEBAL NR 1 ET1",
			CityName:         CityNameROBSector2,
			CountrySubentity: CountrySubentityRO_B,
			PostalZone:       "123456",
			Country:          CountryRO,
		}),
		TaxScheme: &InvoicePartyTaxScheme{
			TaxScheme: TaxSchemeVAT,
			CompanyID: "RO4340188",
		},
		LegalEntity: InvoiceCustomerLegalEntity{
			Name:      "Administratia Sector 3 a Finantelor Publice",
			CompanyID: NewValueWithAttrs("J02/321/2010"),
		},
	}

	invoicePaymentMeans := InvoicePaymentMeans{
		PaymentMeansCode: PaymentMeansCode{
			Code: PaymentMeansDebitTransfer,
		},
		PayeeFinancialAccounts: []PayeeFinancialAccount{
			{
				ID: "RO80RNCB0067054355123456",
			},
		},
	}

	standardTaxCategory := InvoiceLineTaxCategory{
		TaxScheme: TaxSchemeVAT,
		ID:        TaxCategoryVATStandardRate,
		Percent:   types.D(21),
	}

	standardTaxCategory9 := InvoiceLineTaxCategory{
		TaxScheme: TaxSchemeVAT,
		ID:        TaxCategoryVATStandardRate,
		Percent:   types.D(9),
	}

	standardTaxCategory5 := InvoiceLineTaxCategory{
		TaxScheme: TaxSchemeVAT,
		ID:        TaxCategoryVATStandardRate,
		Percent:   types.D(5),
	}

	line1 := assertBuild(assert, NewInvoiceLineBuilder("1", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(46396.67)).
		WithGrossPriceAmount(types.NewPrice(types.D(7.6453))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0102").
		WithItemCommodityClassification(ItemCommodityClassification{
			ItemClassificationCode: ItemClassificationCode{
				Code:   "03222000-3",
				ListID: "STI",
			},
		}).
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineAllowanceBuilder(ron, types.A(801.98)).
				WithAllowanceChargeReasonCode("95").
				WithAllowanceChargeReason("Discount"))).
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineChargeBuilder(ron, types.A(-19272.48)).
				WithBaseAmount(types.A(354715.84)).
				WithAllowanceChargeReasonCode("ZZZ").
				WithAllowanceChargeReason("Mutually defined"))).
		WithItemTaxCategory(standardTaxCategory))

	line2 := assertBuild(assert, NewInvoiceLineBuilder("2", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(622078.28)).
		WithGrossPriceAmount(types.NewPrice(types.D(7.0987))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0104").
		WithItemCommodityClassification(ItemCommodityClassification{
			ItemClassificationCode: ItemClassificationCode{
				Code:   "08055010",
				ListID: "TSP",
			},
		}).
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineAllowanceBuilder(ron, types.A(10454.98)).
				WithAllowanceChargeReasonCode("95").
				WithAllowanceChargeReason("Discount"))).
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineChargeBuilder(ron, types.A(-116445.65)).
				WithBaseAmount(types.A(4415931.88)).
				WithAllowanceChargeReasonCode("ZZZ").
				WithAllowanceChargeReason("Mutually defined"))).
		WithItemTaxCategory(standardTaxCategory))

	line3 := assertBuild(assert, NewInvoiceLineBuilder("3", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(94104.55)).
		WithGrossPriceAmount(types.NewPrice(types.D(7.2813))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0106").
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineAllowanceBuilder(ron, types.A(3589.66)).
				WithAllowanceChargeReasonCode("95").
				WithAllowanceChargeReason("Discount"))).
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineChargeBuilder(ron, types.A(-19458.05)).
				WithBaseAmount(types.A(685199.15)).
				WithAllowanceChargeReasonCode("ZZZ").
				WithAllowanceChargeReason("Mutually defined"))).
		WithItemTaxCategory(standardTaxCategory))

	line4 := assertBuild(assert, NewInvoiceLineBuilder("4", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(3764335.51)).
		WithGrossPriceAmount(types.NewPrice(types.D(6.9490))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0107").
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineAllowanceBuilder(ron, types.A(63265.49)).
				WithAllowanceChargeReasonCode("95").
				WithAllowanceChargeReason("Discount"))).
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineChargeBuilder(ron, types.A(-650523.32)).
				WithBaseAmount(types.A(26158294.04)).
				WithAllowanceChargeReasonCode("ZZZ").
				WithAllowanceChargeReason("Mutually defined"))).
		WithItemTaxCategory(standardTaxCategory))

	line5 := assertBuild(assert, NewInvoiceLineBuilder("5", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(51772.34)).
		WithGrossPriceAmount(types.NewPrice(types.D(7.3995))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0108").
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineAllowanceBuilder(ron, types.A(2980.02)).
				WithAllowanceChargeReasonCode("95").
				WithAllowanceChargeReason("Discount"))).
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineChargeBuilder(ron, types.A(-654.54)).
				WithBaseAmount(types.A(383091.04)).
				WithAllowanceChargeReasonCode("ZZZ").
				WithAllowanceChargeReason("Mutually defined"))).
		WithItemTaxCategory(standardTaxCategory))

	line6 := assertBuild(assert, NewInvoiceLineBuilder("6", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(20807.57)).
		WithGrossPriceAmount(types.NewPrice(types.D(6.8777))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0201").
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineAllowanceBuilder(ron, types.A(757.66)).
				WithAllowanceChargeReasonCode("95").
				WithAllowanceChargeReason("Discount"))).
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineChargeBuilder(ron, types.A(-4664.27)).
				WithBaseAmount(types.A(143107.65)).
				WithAllowanceChargeReasonCode("ZZZ").
				WithAllowanceChargeReason("Mutually defined"))).
		WithItemTaxCategory(standardTaxCategory))

	line7 := assertBuild(assert, NewInvoiceLineBuilder("7", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(217932.24)).
		WithGrossPriceAmount(types.NewPrice(types.D(6.4995))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0203").
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineAllowanceBuilder(
				ron, types.A(3662.64)).
				WithAllowanceChargeReasonCode("95").
				WithAllowanceChargeReason("Discount"))).
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineChargeBuilder(
				ron, types.A(-41460.64)).
				WithBaseAmount(types.A(1416445.96)).
				WithAllowanceChargeReasonCode("ZZZ").
				WithAllowanceChargeReason("Mutually defined"))).
		WithItemTaxCategory(standardTaxCategory))

	line8 := assertBuild(assert, NewInvoiceLineBuilder("8", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(137142.39)).
		WithGrossPriceAmount(types.NewPrice(types.D(6.5974))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0204").
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineAllowanceBuilder(
				ron, types.A(2305.01)).
				WithAllowanceChargeReasonCode("95").
				WithAllowanceChargeReason("Discount"))).
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineChargeBuilder(
				ron, types.A(-26705.29)).
				WithBaseAmount(types.A(904782.33)).
				WithAllowanceChargeReasonCode("ZZZ").
				WithAllowanceChargeReason("Mutually defined"))).
		WithItemTaxCategory(standardTaxCategory))

	line9 := assertBuild(assert, NewInvoiceLineBuilder("9", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(40993.25)).
		WithGrossPriceAmount(types.NewPrice(types.D(7.1266))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0205").
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineAllowanceBuilder(
				ron, types.A(1568.00)).
				WithAllowanceChargeReasonCode("95").
				WithAllowanceChargeReason("Discount"))).
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineChargeBuilder(
				ron, types.A(-8434.52)).
				WithBaseAmount(types.A(292142.98)).
				WithAllowanceChargeReasonCode("ZZZ").
				WithAllowanceChargeReason("Mutually defined"))).
		WithItemTaxCategory(standardTaxCategory))

	line10 := assertBuild(assert, NewInvoiceLineBuilder("10", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(32676.41)).
		WithGrossPriceAmount(types.NewPrice(types.D(3.4328))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0330").
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineAllowanceBuilder(
				ron, types.A(548.95)).
				WithAllowanceChargeReasonCode("95").
				WithAllowanceChargeReason("Discount"))).
		AppendAllowanceCharge(assertBuild(assert,
			NewInvoiceLineChargeBuilder(
				ron, types.A(-6057.84)).
				WithBaseAmount(types.A(112173.07)).
				WithAllowanceChargeReasonCode("ZZZ").
				WithAllowanceChargeReason("Mutually defined"))).
		WithItemTaxCategory(standardTaxCategory))

	line11 := assertBuild(assert, NewInvoiceLineBuilder("11", ron).
		WithItemName("Vignieta").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(2730.01)).
		WithGrossPriceAmount(types.NewPrice(types.D(116.6378))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0452").
		WithItemTaxCategory(standardTaxCategory))

	line12 := assertBuild(assert, NewInvoiceLineBuilder("12", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(958.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(120.6842))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0454").
		WithItemTaxCategory(standardTaxCategory))

	line13 := assertBuild(assert, NewInvoiceLineBuilder("13", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(125.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(24.0754))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0501").
		WithItemTaxCategory(standardTaxCategory))

	line14 := assertBuild(assert, NewInvoiceLineBuilder("14", ron).
		WithItemName("Taxa Ulei").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(35.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(0.3857))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0520").
		WithItemTaxCategory(standardTaxCategory))

	line15 := assertBuild(assert, NewInvoiceLineBuilder("15", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(8875.75)).
		WithGrossPriceAmount(types.NewPrice(types.D(5.9405))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0540").
		WithItemTaxCategory(standardTaxCategory))

	line16 := assertBuild(assert, NewInvoiceLineBuilder("16", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(538.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(119.8705))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0541").
		WithItemTaxCategory(standardTaxCategory))

	line17 := assertBuild(assert, NewInvoiceLineBuilder("17", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(17.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(16.5371))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0550").
		WithItemTaxCategory(standardTaxCategory))

	line18 := assertBuild(assert, NewInvoiceLineBuilder("18", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(639.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(20.8761))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0552").
		WithItemTaxCategory(standardTaxCategory))

	line19 := assertBuild(assert, NewInvoiceLineBuilder("19", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(1084.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(10.8295))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0632").
		WithItemTaxCategory(standardTaxCategory))

	line20 := assertBuild(assert, NewInvoiceLineBuilder("20", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(5.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(4.1920))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0640").
		WithItemTaxCategory(standardTaxCategory))

	line21 := assertBuild(assert, NewInvoiceLineBuilder("21", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(9.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(63.3744))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0702").
		WithItemTaxCategory(standardTaxCategory))

	line22 := assertBuild(assert, NewInvoiceLineBuilder("22", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(198.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(13.4467))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0710").
		WithItemTaxCategory(standardTaxCategory))

	line23 := assertBuild(assert, NewInvoiceLineBuilder("23", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(36.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(24.3350))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0724").
		WithItemTaxCategory(standardTaxCategory))

	line24 := assertBuild(assert, NewInvoiceLineBuilder("24", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(382.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(16.5029))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0810").
		WithItemTaxCategory(standardTaxCategory))

	line25 := assertBuild(assert, NewInvoiceLineBuilder("25", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(18.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(14.1472))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0812").
		WithItemTaxCategory(standardTaxCategory9))

	line26 := assertBuild(assert, NewInvoiceLineBuilder("26", ron).
		WithItemName("Mancare").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(1.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(14.6100))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0820").
		WithItemTaxCategory(standardTaxCategory))

	line27 := assertBuild(assert, NewInvoiceLineBuilder("27", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(2228.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(7.0712))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0824").
		WithItemTaxCategory(standardTaxCategory9))

	line28 := assertBuild(assert, NewInvoiceLineBuilder("28", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(40.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(15.1573))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0830").
		WithItemTaxCategory(standardTaxCategory))

	line29 := assertBuild(assert, NewInvoiceLineBuilder("29", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(1242.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(5.5284))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0832").
		WithItemTaxCategory(standardTaxCategory9))

	line30 := assertBuild(assert, NewInvoiceLineBuilder("30", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(64.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(7.0114))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0851").
		WithItemTaxCategory(standardTaxCategory5))

	line31 := assertBuild(assert, NewInvoiceLineBuilder("31", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(1359.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(18.4958))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0854").
		WithItemTaxCategory(standardTaxCategory))

	line32 := assertBuild(assert, NewInvoiceLineBuilder("32", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(6.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(41.2317))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("0856").
		WithItemTaxCategory(standardTaxCategory5))

	line33 := assertBuild(assert, NewInvoiceLineBuilder("33", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(2315.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(4.50))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("9000").
		WithItemTaxCategory(standardTaxCategory))

	line34 := assertBuild(assert, NewInvoiceLineBuilder("34", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(1.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(1.1200))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("9008").
		WithItemTaxCategory(standardTaxCategory))

	line35 := assertBuild(assert, NewInvoiceLineBuilder("35", ron).
		WithItemName("item name").
		WithUnitCode("C62").
		WithInvoicedQuantity(types.D(15629.00)).
		WithGrossPriceAmount(types.NewPrice(types.D(13.4297))).
		WithBaseQuantity(types.D(1)).
		WithItemSellerID("9012").
		WithItemTaxCategory(standardTaxCategory))

	invoiceBuilder := NewInvoiceBuilder("6422451356").
		WithIssueDate(types.MakeDate(2022, 5, 31)).
		WithDueDate(types.MakeDate(2022, 5, 31)).
		WithInvoiceTypeCode(InvoiceTypeCommercialInvoice).
		WithNote(invoiceNote).
		WithDocumentCurrencyCode(ron).
		WithInvoicePeriod(invoicePeriod).
		WithSupplier(invoiceSupplier).
		WithCustomer(invoiceCustomer).
		WithPaymentMeans(invoicePaymentMeans).
		WithInvoiceLines([]InvoiceLine{
			line1,
			line2,
			line3,
			line4,
			line5,
			line6,
			line7,
			line8,
			line9,
			line10,
			line11,
			line12,
			line13,
			line14,
			line15,
			line16,
			line17,
			line18,
			line19,
			line20,
			line21,
			line22,
			line23,
			line24,
			line25,
			line26,
			line27,
			line28,
			line29,
			line30,
			line31,
			line32,
			line33,
			line34,
			line35,
		})
	invoice := assertBuild(assert, invoiceBuilder)

	invoiceXML, err := invoice.XMLIndent("", "  ")
	assert.NoError(err)

	refData, err := os.ReadFile("./resourses/tests/invoice2.xml")
	assert.NoError(err)

	assert.Equal(strings.TrimSpace(string(refData)), string(invoiceXML))

	fmt.Printf("%s\n", string(invoiceXML))
	/*
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(invoiceXML), string(refData), false)
		fmt.Println(dmp.DiffPrettyText(diffs))
	*/
}

func TestInvoiceRefUnmarshal(t *testing.T) {
	assert := newAssert(t)

	refData, err := os.ReadFile("./resourses/tests/invoice2.xml")
	assert.NoError(err)

	var invoice Invoice
	assert.NoError(UnmarshalInvoice(refData, &invoice))

	marshalXML, err := invoice.XMLIndent("", "  ")
	assert.NoError(err)

	assert.Equal(strings.TrimSpace(string(refData)), string(marshalXML))
}

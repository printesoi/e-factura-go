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
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/printesoi/e-factura-go/pkg/text"
	"github.com/printesoi/e-factura-go/pkg/types"
)

const (
	defaultTestSupplierCompanyID               = "RO1234567890"
	defaultTestSupplierLegalName               = "Seller SRL"
	defaultTestSupplierLegalForm               = "J40/12345/1998"
	defaultTestSupplierAddressLine1            = "Piata Victoriei 1"
	defaultTestSupplierAddressCityName         = CityNameROBSector1
	defaultTestSupplierAddressCountrySubentity = CountrySubentityRO_B

	defaultTestCustomerCompanyID               = "RO987456123"
	defaultTestCustomerLegalName               = "Buyer SRL"
	defaultTestCustomerAddressLine1            = "Piata Victoriei 1"
	defaultTestCustomerAddressCityName         = CityNameROBSector1
	defaultTestCustomerAddressCountrySubentity = CountrySubentityRO_B
)

func getInvoiceSupplierParty() InvoiceSupplierParty {
	var (
		legalName    = defaultTestSupplierLegalName
		companyID    = defaultTestSupplierCompanyID
		regComNo     = defaultTestSupplierLegalForm
		addressLine1 = defaultTestSupplierAddressLine1
		cityName     = defaultTestSupplierAddressCityName
		subentity    = defaultTestSupplierAddressCountrySubentity
	)
	if val := os.Getenv("EFACTURA_TEST_INVOICE_SUPPLIER_LEGAL_NAME"); val != "" {
		legalName = val
	}
	if val := os.Getenv("EFACTURA_TEST_INVOICE_SUPPLIER_COMPANY_ID"); val != "" {
		companyID = val
	}
	if val := os.Getenv("EFACTURA_TEST_INVOICE_SUPPLIER_LEGAL_FORM"); val != "" {
		regComNo = val
	}
	if val := os.Getenv("EFACTURA_TEST_INVOICE_SUPPLIER_ADDRESS_LINE1"); val != "" {
		addressLine1 = val
	}
	if val := os.Getenv("EFACTURA_TEST_INVOICE_SUPPLIER_ADDRESS_CITY_NAME"); val != "" {
		cityName = val
	}
	if val := os.Getenv("EFACTURA_TEST_INVOICE_SUPPLIER_ADDRESS_COUNTRY_SUBENTITY"); val != "" {
		subentity = CountrySubentityType(val)
	}
	return InvoiceSupplierParty{
		PostalAddress: MakeInvoiceSupplierPostalAddress(PostalAddress{
			Country:          CountryRO,
			CountrySubentity: subentity,
			CityName:         cityName,
			Line1:            addressLine1,
		}),
		TaxScheme: &InvoicePartyTaxScheme{
			TaxScheme: TaxSchemeVAT,
			CompanyID: companyID,
		},
		LegalEntity: InvoiceSupplierLegalEntity{
			Name:             legalName,
			CompanyLegalForm: regComNo,
		},
	}
}

func getInvoiceCustomerParty() InvoiceCustomerParty {
	var (
		legalName    = defaultTestCustomerLegalName
		companyID    = defaultTestCustomerCompanyID
		addressLine1 = defaultTestCustomerAddressLine1
		cityName     = defaultTestCustomerAddressCityName
		subentity    = defaultTestCustomerAddressCountrySubentity
	)
	if val := os.Getenv("EFACTURA_TEST_INVOICE_CUSTOMER_LEGAL_NAME"); val != "" {
		legalName = val
	}
	if val := os.Getenv("EFACTURA_TEST_INVOICE_CUSTOMER_COMPANY_ID"); val != "" {
		companyID = val
	}
	if val := os.Getenv("EFACTURA_TEST_INVOICE_CUSTOMER_ADDRESS_LINE1"); val != "" {
		addressLine1 = val
	}
	if val := os.Getenv("EFACTURA_TEST_INVOICE_CUSTOMER_ADDRESS_CITY_NAME"); val != "" {
		cityName = val
	}
	if val := os.Getenv("EFACTURA_TEST_INVOICE_CUSTOMER_ADDRESS_COUNTRY_SUBENTITY"); val != "" {
		subentity = CountrySubentityType(val)
	}
	return InvoiceCustomerParty{
		PostalAddress: MakeInvoiceCustomerPostalAddress(PostalAddress{
			Country:          CountryRO,
			CountrySubentity: subentity,
			CityName:         cityName,
			Line1:            addressLine1,
		}),
		TaxScheme: &InvoicePartyTaxScheme{
			TaxScheme: TaxSchemeVAT,
			CompanyID: companyID,
		},
		LegalEntity: InvoiceCustomerLegalEntity{
			Name: legalName,
		},
	}
}

func TestInvoiceLineBuilder(t *testing.T) {
	assert := assert.New(t)

	a := func(a types.Amount) string {
		return a.String()
	}
	d := func(d types.Decimal) string {
		return d.String()
	}

	{
		b := NewInvoiceLineBuilder("1", CurrencyRON)
		_, err := b.Build()
		assert.Error(err, "should not build if required fields are missing")
	}
	type lineTest struct {
		ID             string
		CurrencyID     CurrencyCodeType
		UnitCode       UnitCodeType
		Quantity       types.Decimal
		BaseQuantity   types.Decimal
		GrossPrice     types.Amount
		PriceDeduction types.Amount
		Allowances     []types.Amount
		Charges        []types.Amount
		ItemName       string
		TaxCategory    InvoiceLineTaxCategory

		ExpectedLineAmount types.Amount
	}
	tests := []lineTest{
		{
			// A.1.4 Exemplul 1. Linia 1
			ID:         "01.1",
			CurrencyID: CurrencyEUR,
			UnitCode:   "XGR",
			Quantity:   types.D(5),
			GrossPrice: types.A(12),
			ItemName:   text.Transliterate("Sticle cu vin"),
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   types.D(25),
			},
			ExpectedLineAmount: types.A(60),
		},
		{
			// A.1.4 Exemplul 1. Linia 2
			ID:         "01.2",
			CurrencyID: CurrencyEUR,
			UnitCode:   "XBX",
			Quantity:   types.D(1),
			GrossPrice: types.A(90),
			ItemName:   text.Transliterate("Vin - cutie de 6"),
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   types.D(25),
			},
			ExpectedLineAmount: types.A(90),
		},
		{
			// A.1.4 Exemplul 2
			ID:           "02",
			CurrencyID:   CurrencyEUR,
			UnitCode:     "C62",
			Quantity:     types.D(10_000),
			BaseQuantity: types.D(1_000),
			GrossPrice:   types.A(4.5),
			ItemName:     text.Transliterate("Șurub"),
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   types.D(25),
			},
			ExpectedLineAmount: types.A(45),
		},
		{
			// A.1.4 Exemplul 3
			ID:             "03",
			CurrencyID:     CurrencyEUR,
			UnitCode:       "58",
			Quantity:       types.D(1.3),
			GrossPrice:     types.A(10),
			PriceDeduction: types.A(0.5),
			ItemName:       text.Transliterate("Pui"),
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   types.D(12.5),
			},
			ExpectedLineAmount: types.A(12.35),
		},
		{
			// A.1.5 Exemplul 4 (Reduceri, deduceri şi taxe suplimentare). Line 1
			ID:             "04.1",
			CurrencyID:     CurrencyEUR,
			UnitCode:       "XBX",
			Quantity:       types.D(25),
			GrossPrice:     types.A(9.5),
			PriceDeduction: types.A(1),
			ItemName:       text.Transliterate("Stilou"),
			Charges: []types.Amount{
				types.A(10),
			},
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   types.D(25),
			},
			ExpectedLineAmount: types.A(222.50),
		},
		{
			// A.1.5 Exemplul 4 (Reduceri, deduceri şi taxe suplimentare). Line 2
			ID:         "04.2",
			CurrencyID: CurrencyEUR,
			UnitCode:   "RM",
			Quantity:   types.D(15),
			GrossPrice: types.A(4.5),
			ItemName:   text.Transliterate("Hârtie"),
			Allowances: []types.Amount{
				types.A(3.38),
			},
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   types.D(25),
			},
			ExpectedLineAmount: types.A(64.12),
		},
		{
			// A.1.6 Exemplul 5 (Linie a facturii negativă). Linia 1
			ID:             "05.1",
			CurrencyID:     CurrencyEUR,
			UnitCode:       "XBX",
			Quantity:       types.D(25),
			GrossPrice:     types.A(9.5),
			PriceDeduction: types.A(1),
			ItemName:       text.Transliterate("Stilou"),
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   types.D(25),
			},
			ExpectedLineAmount: types.A(212.50),
		},
		{
			// A.1.6 Exemplul 5 (Linie a facturii negativă). Linia 2
			ID:             "05.2",
			CurrencyID:     CurrencyEUR,
			UnitCode:       "XBX",
			Quantity:       types.D(-10),
			GrossPrice:     types.A(9.5),
			PriceDeduction: types.A(1),
			ItemName:       text.Transliterate("Stilou"),
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   types.D(25),
			},
			ExpectedLineAmount: types.A(-85),
		},
	}
	for _, t := range tests {
		b := NewInvoiceLineBuilder(t.ID, t.CurrencyID).
			WithUnitCode(t.UnitCode).WithInvoicedQuantity(t.Quantity).
			WithGrossPriceAmount(t.GrossPrice).
			WithItemName(t.ItemName).
			WithItemTaxCategory(t.TaxCategory)
		if !t.BaseQuantity.IsZero() {
			b.WithBaseQuantity(t.BaseQuantity)
		}
		if !t.PriceDeduction.IsZero() {
			b.WithPriceDeduction(t.PriceDeduction)
		}
		for _, allowance := range t.Allowances {
			if !allowance.IsZero() {
				lineAllowance, err := NewInvoiceLineAllowanceBuilder(t.CurrencyID, allowance).Build()
				if assert.NoError(err) {
					b.AppendAllowanceCharge(lineAllowance)
				}
			}
		}
		for _, charge := range t.Charges {
			if !charge.IsZero() {
				lineCharge, err := NewInvoiceLineChargeBuilder(t.CurrencyID, charge).Build()
				if assert.NoError(err) {
					b.AppendAllowanceCharge(lineCharge)
				}
			}
		}

		line, err := b.Build()
		if assert.NoError(err) {
			assert.Equal(a(t.ExpectedLineAmount), a(line.LineExtensionAmount.Amount))

			assert.Equal(t.ID, line.ID)
			assert.Equal(d(t.Quantity), d(line.InvoicedQuantity.Quantity))
			assert.Equal(t.ItemName, line.Item.Name)
			assert.Equal(t.TaxCategory, line.Item.TaxCategory)

			// TODO: compare all fields
		}
	}
}

func TestInvoiceBuilder(t *testing.T) {
	assert := assert.New(t)

	a := func(d types.Amount) string {
		return d.String()
	}
	d := func(d types.Decimal) string {
		return d.String()
	}

	{
		b := NewInvoiceBuilder("1")
		_, err := b.Build()
		assert.Error(err, "should not build if required fields are missing")
	}
	{
		// A.1.6 Exemplul 5 (Linie a facturii negativă)
		documentCurrencyID := CurrencyEUR

		var lines []InvoiceLine

		standardTaxCategory := InvoiceLineTaxCategory{
			TaxScheme: TaxSchemeVAT,
			ID:        TaxCategoryVATStandardRate,
			Percent:   types.D(25),
		}

		line1, err := NewInvoiceLineBuilder("1", documentCurrencyID).
			WithUnitCode("XBX").
			WithInvoicedQuantity(types.D(25)).
			WithGrossPriceAmount(types.A(9.5)).
			WithPriceDeduction(types.A(1)).
			WithItemName("Stilouri").
			WithItemTaxCategory(standardTaxCategory).
			Build()
		if assert.NoError(err) {
			lines = append(lines, line1)
		}

		line2, err := NewInvoiceLineBuilder("2", documentCurrencyID).
			WithUnitCode("XBX").
			WithInvoicedQuantity(types.D(-10)).
			WithGrossPriceAmount(types.A(9.5)).
			WithPriceDeduction(types.A(1)).
			WithItemName("Stilouri").
			WithItemTaxCategory(standardTaxCategory).
			Build()
		if assert.NoError(err) {
			lines = append(lines, line2)
		}

		invoiceBuilder := NewInvoiceBuilder("test.example.05").
			WithIssueDate(types.MakeDate(2024, 3, 1)).
			WithDueDate(types.MakeDate(2024, 3, 31)).
			WithInvoiceTypeCode(InvoiceTypeCommercialInvoice).
			WithDocumentCurrencyCode(documentCurrencyID).
			WithSupplier(getInvoiceSupplierParty()).
			WithCustomer(getInvoiceCustomerParty()).
			WithInvoiceLines(lines)

		invoice, err := invoiceBuilder.Build()
		if assert.NoError(err) {
			// Invoice lines
			if assert.Equal(2, len(invoice.InvoiceLines), "should have correct number of lines") {
				line1 := invoice.InvoiceLines[0]
				assert.Equal(a(types.A(212.5)), a(line1.LineExtensionAmount.Amount))

				line2 := invoice.InvoiceLines[1]
				assert.Equal(a(types.A(-85)), a(line2.LineExtensionAmount.Amount))
			}

			// VAT details (BG-23)
			if assert.Equal(1, len(invoice.TaxTotal), "Must have only one TaxTotal") {
				if assert.Equal(1, len(invoice.TaxTotal[0].TaxSubtotals)) {
					subtotal := invoice.TaxTotal[0].TaxSubtotals[0]
					assert.Equal(TaxCategoryVATStandardRate, subtotal.TaxCategory.ID)
					assert.Equal(d(types.D(25)), d(subtotal.TaxCategory.Percent))
					assert.Equal(a(types.A(127.5)), a(subtotal.TaxableAmount.Amount))
					assert.Equal(a(types.A(31.88)), a(subtotal.TaxAmount.Amount))
				}

				// BT-110
				if assert.NotNil(invoice.TaxTotal[0].TaxAmount, "BT-110 must exist") {
					assert.Equal(a(types.A(31.88)), a(invoice.TaxTotal[0].TaxAmount.Amount), "BT-110 incorrect value")
				}
			}

			// Document totals (BG-22)
			// BT-106
			assert.Equal(a(types.A(127.5)), a(invoice.LegalMonetaryTotal.LineExtensionAmount.Amount), "BT-106 incorrect value")
			// BT-109
			assert.Equal(a(types.A(127.5)), a(invoice.LegalMonetaryTotal.TaxExclusiveAmount.Amount), "BT-109 incorrect value")
			// BT-112
			assert.Equal(a(types.A(159.38)), a(invoice.LegalMonetaryTotal.TaxInclusiveAmount.Amount), "BT-112 incorrect value")
			// BT-115
			assert.Equal(a(types.A(159.38)), a(invoice.LegalMonetaryTotal.PayableAmount.Amount), "BT-115 incorrect value")
		}
	}
	{
		// A.1.8 Exemplul 7 (Cota normală de TVA cu linii scutite de TVA)
		buildInvoice := func(documentCurrencyID CurrencyCodeType) (Invoice, error) {
			var lines []InvoiceLine

			line1, err := NewInvoiceLineBuilder("1", documentCurrencyID).
				WithUnitCode("H87").
				WithInvoicedQuantity(types.D(5)).
				WithGrossPriceAmount(types.A(25.0)).
				WithItemName(text.Transliterate("Cerneală pentru imprimantă")).
				WithItemTaxCategory(InvoiceLineTaxCategory{
					TaxScheme: TaxSchemeVAT,
					ID:        TaxCategoryVATStandardRate,
					Percent:   types.D(25),
				}).
				Build()
			if assert.NoError(err) {
				lines = append(lines, line1)
			}

			line2, err := NewInvoiceLineBuilder("2", documentCurrencyID).
				WithUnitCode("H87").
				WithInvoicedQuantity(types.D(1)).
				WithGrossPriceAmount(types.A(24.0)).
				WithItemName(text.Transliterate("Imprimare afiș")).
				WithItemTaxCategory(InvoiceLineTaxCategory{
					TaxScheme: TaxSchemeVAT,
					ID:        TaxCategoryVATStandardRate,
					Percent:   types.D(10),
				}).
				Build()
			if assert.NoError(err) {
				lines = append(lines, line2)
			}

			line3, err := NewInvoiceLineBuilder("3", documentCurrencyID).
				WithUnitCode("H87").
				WithInvoicedQuantity(types.D(1)).
				WithGrossPriceAmount(types.A(136.0)).
				WithItemName(text.Transliterate("Scaun de birou")).
				WithItemTaxCategory(InvoiceLineTaxCategory{
					TaxScheme: TaxSchemeVAT,
					ID:        TaxCategoryVATStandardRate,
					Percent:   types.D(25),
				}).
				Build()
			if assert.NoError(err) {
				lines = append(lines, line3)
			}

			line4, err := NewInvoiceLineBuilder("4", documentCurrencyID).
				WithUnitCode("H87").
				WithInvoicedQuantity(types.D(1)).
				WithGrossPriceAmount(types.A(95.0)).
				WithItemName(text.Transliterate("Tastatură fără fir")).
				WithItemTaxCategory(InvoiceLineTaxCategory{
					TaxScheme: TaxSchemeVAT,
					ID:        TaxCategoryVATExempt,
				}).
				Build()
			if assert.NoError(err) {
				lines = append(lines, line4)
			}

			line5, err := NewInvoiceLineBuilder("5", documentCurrencyID).
				WithUnitCode("H87").
				WithInvoicedQuantity(types.D(1)).
				WithGrossPriceAmount(types.A(53.0)).
				WithItemName(text.Transliterate("Cablu de adaptare")).
				WithItemTaxCategory(InvoiceLineTaxCategory{
					TaxScheme: TaxSchemeVAT,
					ID:        TaxCategoryVATExempt,
				}).
				Build()
			if assert.NoError(err) {
				lines = append(lines, line5)
			}

			invoiceBuilder := NewInvoiceBuilder("test.example.07").
				WithIssueDate(types.MakeDate(2024, 3, 1)).
				WithDueDate(types.MakeDate(2024, 4, 1)).
				WithInvoiceTypeCode(InvoiceTypeSelfBilledInvoice).
				WithDocumentCurrencyCode(documentCurrencyID).
				WithSupplier(getInvoiceSupplierParty()).
				WithCustomer(getInvoiceCustomerParty()).
				WithInvoiceLines(lines).
				AddTaxExemptionReason(TaxCategoryVATExempt, "MOTIVUL A", "")

			documentAllowance, err := NewInvoiceDocumentAllowanceBuilder(
				documentCurrencyID,
				types.A(15),
				InvoiceTaxCategory{
					TaxScheme: TaxSchemeVAT,
					ID:        TaxCategoryVATStandardRate,
					Percent:   types.D(25),
				},
			).WithAllowanceChargeReason("Motivul C").Build()
			if assert.NoError(err) {
				invoiceBuilder.AppendAllowanceCharge(documentAllowance)
			}

			documentCharge, err := NewInvoiceDocumentChargeBuilder(
				documentCurrencyID,
				types.A(35),
				InvoiceTaxCategory{
					TaxScheme: TaxSchemeVAT,
					ID:        TaxCategoryVATStandardRate,
					Percent:   types.D(25),
				},
			).WithAllowanceChargeReason("Motivul B").Build()
			if assert.NoError(err) {
				invoiceBuilder.AppendAllowanceCharge(documentCharge)
			}

			if documentCurrencyID != CurrencyRON {
				invoiceBuilder.WithTaxCurrencyCode(CurrencyRON)
				invoiceBuilder.WithDocumentToTaxCurrencyExchangeRate(types.D(4.9691))
			}
			return invoiceBuilder.Build()
		}

		{
			// BT-5 is RON
			invoice, err := buildInvoice(CurrencyRON)
			{
				ivxml, _ := invoice.XML()
				fmt.Println(string(ivxml))
			}
			if assert.NoError(err) {
				// Invoice lines
				if assert.Equal(5, len(invoice.InvoiceLines)) {
					line1 := invoice.InvoiceLines[0]
					assert.Equal(a(types.A(125.0)), a(line1.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATStandardRate, line1.Item.TaxCategory.ID)
					assert.Equal(d(types.D(25.0)), d(line1.Item.TaxCategory.Percent))

					line2 := invoice.InvoiceLines[1]
					assert.Equal(a(types.A(24.0)), a(line2.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATStandardRate, line2.Item.TaxCategory.ID)
					assert.Equal(d(types.D(10.0)), d(line2.Item.TaxCategory.Percent))

					line3 := invoice.InvoiceLines[2]
					assert.Equal(a(types.A(136.0)), a(line3.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATStandardRate, line3.Item.TaxCategory.ID)
					assert.Equal(d(types.D(25.0)), d(line3.Item.TaxCategory.Percent))

					line4 := invoice.InvoiceLines[3]
					assert.Equal(a(types.A(95.0)), a(line4.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATExempt, line4.Item.TaxCategory.ID)
					assert.Equal(d(types.D(0.0)), d(line4.Item.TaxCategory.Percent))

					line5 := invoice.InvoiceLines[4]
					assert.Equal(a(types.A(53.0)), a(line5.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATExempt, line5.Item.TaxCategory.ID)
					assert.Equal(d(types.D(0.0)), d(line5.Item.TaxCategory.Percent))
				}

				// Document totals (BG-22)
				// BT-106
				assert.Equal(a(types.A(433.0)), a(invoice.LegalMonetaryTotal.LineExtensionAmount.Amount), "BT-106 incorrect value")
				// BT-107
				if allowanceTotalAmount := invoice.LegalMonetaryTotal.AllowanceTotalAmount; assert.NotNil(allowanceTotalAmount, "BT-107 must be non-nil") {
					assert.Equal(a(types.A(15)), a(allowanceTotalAmount.Amount), "BT-107 incorrect value")
				}
				// BT-108
				if chargeTotalAmount := invoice.LegalMonetaryTotal.ChargeTotalAmount; assert.NotNil(chargeTotalAmount, "BT-108 must be non-nil") {
					assert.Equal(a(types.A(35)), a(chargeTotalAmount.Amount), "BT-108 incorrect value")
				}
				// BT-109
				assert.Equal(a(types.A(453.0)), a(invoice.LegalMonetaryTotal.TaxExclusiveAmount.Amount), "BT-109 incorrect value")
				// BT-110
				if assert.Equal(1, len(invoice.TaxTotal), "Must have only one TaxTotal") &&
					assert.NotNil(invoice.TaxTotal[0].TaxAmount, "BT-110 must exist") {

					assert.Equal(a(types.A(72.65)), a(invoice.TaxTotal[0].TaxAmount.Amount), "BT-110 incorrect value")
				}
				// BT-112
				assert.Equal(a(types.A(525.65)), a(invoice.LegalMonetaryTotal.TaxInclusiveAmount.Amount), "BT-112 incorrect value")
				// BT-115
				assert.Equal(a(types.A(525.65)), a(invoice.LegalMonetaryTotal.PayableAmount.Amount), "BT-115 incorrect value")
			}
		}
		{
			// BT-5 is EUR, BT-6 is RON
			documentCurrencyID := CurrencyEUR
			invoice, err := buildInvoice(documentCurrencyID)
			if assert.NoError(err) {
				// Invoice lines
				if assert.Equal(5, len(invoice.InvoiceLines)) {
					line1 := invoice.InvoiceLines[0]
					assert.Equal(a(types.A(125.0)), a(line1.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATStandardRate, line1.Item.TaxCategory.ID)
					assert.Equal(d(types.D(25.0)), d(line1.Item.TaxCategory.Percent))

					line2 := invoice.InvoiceLines[1]
					assert.Equal(a(types.A(24.0)), a(line2.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATStandardRate, line2.Item.TaxCategory.ID)
					assert.Equal(d(types.D(10.0)), d(line2.Item.TaxCategory.Percent))

					line3 := invoice.InvoiceLines[2]
					assert.Equal(a(types.A(136.0)), a(line3.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATStandardRate, line3.Item.TaxCategory.ID)
					assert.Equal(d(types.D(25.0)), d(line3.Item.TaxCategory.Percent))

					line4 := invoice.InvoiceLines[3]
					assert.Equal(a(types.A(95.0)), a(line4.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATExempt, line4.Item.TaxCategory.ID)
					assert.Equal(d(types.D(0.0)), d(line4.Item.TaxCategory.Percent))

					line5 := invoice.InvoiceLines[4]
					assert.Equal(a(types.A(53.0)), a(line5.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATExempt, line5.Item.TaxCategory.ID)
					assert.Equal(d(types.D(0.0)), d(line5.Item.TaxCategory.Percent))
				}

				// Document totals (BG-22)
				// BT-106
				assert.Equal(a(types.A(433.0)), a(invoice.LegalMonetaryTotal.LineExtensionAmount.Amount), "BT-106 incorrect value")
				// BT-107
				if allowanceTotalAmount := invoice.LegalMonetaryTotal.AllowanceTotalAmount; assert.NotNil(allowanceTotalAmount, "BT-107 must be non-nil") {
					assert.Equal(a(types.A(15)), a(allowanceTotalAmount.Amount), "BT-107 incorrect value")
				}
				// BT-108
				if chargeTotalAmount := invoice.LegalMonetaryTotal.ChargeTotalAmount; assert.NotNil(chargeTotalAmount, "BT-108 must be non-nil") {
					assert.Equal(a(types.A(35)), a(chargeTotalAmount.Amount), "BT-108 incorrect value")
				}
				// BT-109
				assert.Equal(a(types.A(453.0)), a(invoice.LegalMonetaryTotal.TaxExclusiveAmount.Amount), "BT-109 incorrect value")
				// BT-110
				if assert.Equal(2, len(invoice.TaxTotal), "Must have a TaxTotal for each currency") {
					taxTotalID := findTaxTotalByCurrency(invoice.TaxTotal, documentCurrencyID)
					if assert.True(taxTotalID >= 0, "TaxTotal for document currency code(BT-5) not set") &&
						assert.NotNil(invoice.TaxTotal[taxTotalID].TaxAmount, "BT-110 must be non-nil") {
						assert.Equal(a(types.A(72.65)), a(invoice.TaxTotal[taxTotalID].TaxAmount.Amount), "BT-110 incorrect value")
					}

					taxTotalTaxCurrencyID := findTaxTotalByCurrency(invoice.TaxTotal, CurrencyRON)
					if assert.True(taxTotalTaxCurrencyID >= 0, "TaxTotal for tax currency code(BT-6) not set") &&
						assert.NotNil(invoice.TaxTotal[taxTotalTaxCurrencyID].TaxAmount, "BT-111 must be non-nil") {
						assert.Equal(a(types.A(361.01)), a(invoice.TaxTotal[taxTotalTaxCurrencyID].TaxAmount.Amount), "BT-110 incorrect value")
					}
				}
				// BT-112
				assert.Equal(a(types.A(525.65)), a(invoice.LegalMonetaryTotal.TaxInclusiveAmount.Amount), "BT-112 incorrect value")
				// BT-115
				assert.Equal(a(types.A(525.65)), a(invoice.LegalMonetaryTotal.PayableAmount.Amount), "BT-115 incorrect value")
			}
		}
	}
}

func findTaxTotalByCurrency(totals []InvoiceTaxTotal, currencyID CurrencyCodeType) int {
	for i, t := range totals {
		if t.TaxAmount != nil && t.TaxAmount.CurrencyID == currencyID {
			return i
		}
	}
	return -1
}

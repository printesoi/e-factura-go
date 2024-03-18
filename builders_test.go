package efactura

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

	a := func(d Decimal) string {
		return d.StringFixed(2)
	}
	d := func(d Decimal) string {
		return d.String()
	}

	{
		b := NewInvoiceLineBuilder("1", CurrencyRON)
		_, ok := b.Build()
		assert.False(ok, "should not build if required fields are missing")
	}
	type lineTest struct {
		ID             string
		CurrencyID     CurrencyCodeType
		UnitCode       UnitCodeType
		Quantity       Decimal
		BaseQuantity   Decimal
		GrossPrice     Decimal
		PriceDeduction Decimal
		Allowances     []Decimal
		Charges        []Decimal
		ItemName       string
		TaxCategory    InvoiceLineTaxCategory

		ExpectedLineAmount Decimal
	}
	tests := []lineTest{
		{
			// A.1.4 Exemplul 1. Linia 1
			ID:         "01.1",
			CurrencyID: CurrencyEUR,
			UnitCode:   "XGR",
			Quantity:   D(5),
			GrossPrice: D(12),
			ItemName:   Transliterate("Sticle cu vin"),
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   D(25),
			},
			ExpectedLineAmount: D(60),
		},
		{
			// A.1.4 Exemplul 1. Linia 2
			ID:         "01.2",
			CurrencyID: CurrencyEUR,
			UnitCode:   "XBX",
			Quantity:   D(1),
			GrossPrice: D(90),
			ItemName:   Transliterate("Vin - cutie de 6"),
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   D(25),
			},
			ExpectedLineAmount: D(90),
		},
		{
			// A.1.4 Exemplul 2
			ID:           "02",
			CurrencyID:   CurrencyEUR,
			UnitCode:     "C62",
			Quantity:     D(10_000),
			BaseQuantity: D(1_000),
			GrossPrice:   D(4.5),
			ItemName:     Transliterate("Șurub"),
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   D(25),
			},
			ExpectedLineAmount: D(45),
		},
		{
			// A.1.4 Exemplul 3
			ID:             "03",
			CurrencyID:     CurrencyEUR,
			UnitCode:       "58",
			Quantity:       D(1.3),
			GrossPrice:     D(10),
			PriceDeduction: D(0.5),
			ItemName:       Transliterate("Pui"),
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   D(12.5),
			},
			ExpectedLineAmount: D(12.35),
		},
		{
			// A.1.5 Exemplul 4 (Reduceri, deduceri şi taxe suplimentare). Line 1
			ID:             "04.1",
			CurrencyID:     CurrencyEUR,
			UnitCode:       "XBX",
			Quantity:       D(25),
			GrossPrice:     D(9.5),
			PriceDeduction: D(1),
			ItemName:       Transliterate("Stilou"),
			Charges: []Decimal{
				D(10),
			},
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   D(25),
			},
			ExpectedLineAmount: D(222.50),
		},
		{
			// A.1.5 Exemplul 4 (Reduceri, deduceri şi taxe suplimentare). Line 2
			ID:         "04.2",
			CurrencyID: CurrencyEUR,
			UnitCode:   "RM",
			Quantity:   D(15),
			GrossPrice: D(4.5),
			ItemName:   Transliterate("Hârtie"),
			Allowances: []Decimal{
				D(3.38),
			},
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   D(25),
			},
			ExpectedLineAmount: D(64.12),
		},
		{
			// A.1.6 Exemplul 5 (Linie a facturii negativă). Linia 1
			ID:             "05.1",
			CurrencyID:     CurrencyEUR,
			UnitCode:       "XBX",
			Quantity:       D(25),
			GrossPrice:     D(9.5),
			PriceDeduction: D(1),
			ItemName:       Transliterate("Stilou"),
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   D(25),
			},
			ExpectedLineAmount: D(212.50),
		},
		{
			// A.1.6 Exemplul 5 (Linie a facturii negativă). Linia 2
			ID:             "05.2",
			CurrencyID:     CurrencyEUR,
			UnitCode:       "XBX",
			Quantity:       D(-10),
			GrossPrice:     D(9.5),
			PriceDeduction: D(1),
			ItemName:       Transliterate("Stilou"),
			TaxCategory: InvoiceLineTaxCategory{
				TaxScheme: TaxSchemeVAT,
				ID:        TaxCategoryVATStandardRate,
				Percent:   D(25),
			},
			ExpectedLineAmount: D(-85),
		},
	}
	for _, t := range tests {
		b := NewInvoiceLineBuilder(t.ID, t.CurrencyID).
			WithUnitCode(t.UnitCode).WithInvoicedQuantity(t.Quantity).
			WithGrossPriceAmount(t.GrossPrice).
			WithItem(InvoiceLineItem{
				Name:        t.ItemName,
				TaxCategory: t.TaxCategory,
			})
		if !t.BaseQuantity.IsZero() {
			b.WithBaseQuantity(t.BaseQuantity)
		}
		if !t.PriceDeduction.IsZero() {
			b.WithPriceDeduction(t.PriceDeduction)
		}
		for _, allowance := range t.Allowances {
			if !allowance.IsZero() {
				lineAllowance, ok := NewInvoiceLineAllowanceBuilder(t.CurrencyID, allowance).Build()
				if assert.True(ok) {
					b.AppendAllowanceCharge(lineAllowance)
				}
			}
		}
		for _, charge := range t.Charges {
			if !charge.IsZero() {
				lineCharge, ok := NewInvoiceLineChargeBuilder(t.CurrencyID, charge).Build()
				if assert.True(ok) {
					b.AppendAllowanceCharge(lineCharge)
				}
			}
		}
		line, ok := b.Build()

		if assert.True(ok) {
			assert.Equal(a(t.ExpectedLineAmount), a(line.LineExtensionAmount.Amount))

			assert.Equal(t.ID, line.ID)
			assert.Equal(d(t.Quantity), d(line.InvoicedQuantity.Quantity))
			assert.Equal(t.ItemName, line.Item.Name)
			// TODO: compare all fields
		}
	}
}

func TestInvoiceBuilder(t *testing.T) {
	assert := assert.New(t)

	a := func(d Decimal) string {
		return d.StringFixed(2)
	}
	d := func(d Decimal) string {
		return d.String()
	}

	{
		b := NewInvoiceBuilder("1")
		_, ok := b.Build()
		assert.False(ok, "should not build if required fields are missing")
	}
	{
		// A.1.6 Exemplul 5 (Linie a facturii negativă)
		documentCurrencyID := CurrencyEUR

		var lines []InvoiceLine

		line1, ok := NewInvoiceLineBuilder("1", documentCurrencyID).
			WithUnitCode("XBX").
			WithInvoicedQuantity(D(25)).
			WithGrossPriceAmount(D(9.5)).
			WithPriceDeduction(D(1)).
			WithItem(InvoiceLineItem{
				Name: "Stilouri",
				TaxCategory: InvoiceLineTaxCategory{
					TaxScheme: TaxSchemeVAT,
					ID:        TaxCategoryVATStandardRate,
					Percent:   D(25),
				},
			}).
			Build()
		if assert.True(ok) {
			lines = append(lines, line1)
		}

		line2, ok := NewInvoiceLineBuilder("2", documentCurrencyID).
			WithUnitCode("XBX").
			WithInvoicedQuantity(D(-10)).
			WithGrossPriceAmount(D(9.5)).
			WithPriceDeduction(D(1)).
			WithItem(InvoiceLineItem{
				Name: "Stilouri",
				TaxCategory: InvoiceLineTaxCategory{
					TaxScheme: TaxSchemeVAT,
					ID:        TaxCategoryVATStandardRate,
					Percent:   D(25),
				},
			}).
			Build()
		if assert.True(ok) {
			lines = append(lines, line2)
		}

		invoiceBuilder := NewInvoiceBuilder("test.example.05").
			WithIssueDate(MakeDate(2024, 3, 1)).
			WithDueDate(MakeDate(2024, 3, 31)).
			WithInvoiceTypeCode(InvoiceTypeCommercialInvoice).
			WithDocumentCurrencyCode(documentCurrencyID).
			WithSupplier(getInvoiceSupplierParty()).
			WithCustomer(getInvoiceCustomerParty()).
			WithInvoiceLines(lines)

		invoice, ok := invoiceBuilder.Build()
		if assert.True(ok) {
			// Invoice lines
			if assert.Equal(2, len(invoice.InvoiceLines), "should have correct number of lines") {
				line1 := invoice.InvoiceLines[0]
				assert.Equal(a(D(212.5)), a(line1.LineExtensionAmount.Amount))

				line2 := invoice.InvoiceLines[1]
				assert.Equal(a(D(-85)), a(line2.LineExtensionAmount.Amount))
			}

			// VAT details (BG-23)
			if assert.Equal(1, len(invoice.TaxTotal), "Must have only one TaxTotal") {
				if assert.Equal(1, len(invoice.TaxTotal[0].TaxSubtotals)) {
					subtotal := invoice.TaxTotal[0].TaxSubtotals[0]
					assert.Equal(TaxCategoryVATStandardRate, subtotal.TaxCategory.ID)
					assert.Equal(a(D(25)), a(subtotal.TaxCategory.Percent))
					assert.Equal(a(D(127.5)), a(subtotal.TaxableAmount.Amount))
					assert.Equal(a(D(31.88)), a(subtotal.TaxAmount.Amount))
				}

				// BT-110
				if assert.NotNil(invoice.TaxTotal[0].TaxAmount, "BT-110 must exist") {
					assert.Equal(a(D(31.88)), a(invoice.TaxTotal[0].TaxAmount.Amount), "BT-110 incorrect value")
				}
			}

			// Document totals (BG-22)
			// BT-106
			assert.Equal(a(D(127.5)), a(invoice.LegalMonetaryTotal.LineExtensionAmount.Amount), "BT-106 incorrect value")
			// BT-109
			assert.Equal(a(D(127.5)), a(invoice.LegalMonetaryTotal.TaxExclusiveAmount.Amount), "BT-109 incorrect value")
			// BT-112
			assert.Equal(a(D(159.38)), a(invoice.LegalMonetaryTotal.TaxInclusiveAmount.Amount), "BT-112 incorrect value")
			// BT-115
			assert.Equal(a(D(159.38)), a(invoice.LegalMonetaryTotal.PayableAmount.Amount), "BT-115 incorrect value")
		}
	}
	{
		// A.1.8 Exemplul 7 (Cota normală de TVA cu linii scutite de TVA)
		buildInvoice := func(documentCurrencyID CurrencyCodeType) (Invoice, bool) {
			var lines []InvoiceLine

			line1, ok := NewInvoiceLineBuilder("1", documentCurrencyID).
				WithUnitCode("H87").
				WithInvoicedQuantity(D(5)).
				WithGrossPriceAmount(D(25.0)).
				WithItem(InvoiceLineItem{
					Name: Transliterate("Cerneală pentru imprimantă"),
					TaxCategory: InvoiceLineTaxCategory{
						TaxScheme: TaxSchemeVAT,
						ID:        TaxCategoryVATStandardRate,
						Percent:   D(25),
					},
				}).
				Build()
			if assert.True(ok) {
				lines = append(lines, line1)
			}

			line2, ok := NewInvoiceLineBuilder("2", documentCurrencyID).
				WithUnitCode("H87").
				WithInvoicedQuantity(D(1)).
				WithGrossPriceAmount(D(24.0)).
				WithItem(InvoiceLineItem{
					Name: Transliterate("Imprimare afiș"),
					TaxCategory: InvoiceLineTaxCategory{
						TaxScheme: TaxSchemeVAT,
						ID:        TaxCategoryVATStandardRate,
						Percent:   D(10),
					},
				}).
				Build()
			if assert.True(ok) {
				lines = append(lines, line2)
			}

			line3, ok := NewInvoiceLineBuilder("3", documentCurrencyID).
				WithUnitCode("H87").
				WithInvoicedQuantity(D(1)).
				WithGrossPriceAmount(D(136.0)).
				WithItem(InvoiceLineItem{
					Name: Transliterate("Scaun de birou"),
					TaxCategory: InvoiceLineTaxCategory{
						TaxScheme: TaxSchemeVAT,
						ID:        TaxCategoryVATStandardRate,
						Percent:   D(25),
					},
				}).
				Build()
			if assert.True(ok) {
				lines = append(lines, line3)
			}

			line4, ok := NewInvoiceLineBuilder("4", documentCurrencyID).
				WithUnitCode("H87").
				WithInvoicedQuantity(D(1)).
				WithGrossPriceAmount(D(95.0)).
				WithItem(InvoiceLineItem{
					Name: Transliterate("Tastatură fără fir"),
					TaxCategory: InvoiceLineTaxCategory{
						TaxScheme: TaxSchemeVAT,
						ID:        TaxCategoryVATExempt,
					},
				}).
				Build()
			if assert.True(ok) {
				lines = append(lines, line4)
			}

			line5, ok := NewInvoiceLineBuilder("5", documentCurrencyID).
				WithUnitCode("H87").
				WithInvoicedQuantity(D(1)).
				WithGrossPriceAmount(D(53.0)).
				WithItem(InvoiceLineItem{
					Name: Transliterate("Cablu de adaptare"),
					TaxCategory: InvoiceLineTaxCategory{
						TaxScheme: TaxSchemeVAT,
						ID:        TaxCategoryVATExempt,
					},
				}).
				Build()
			if assert.True(ok) {
				lines = append(lines, line5)
			}

			invoiceBuilder := NewInvoiceBuilder("test.example.07").
				WithIssueDate(MakeDate(2024, 3, 1)).
				WithDueDate(MakeDate(2024, 4, 1)).
				WithInvoiceTypeCode(InvoiceTypeCommercialInvoice).
				WithDocumentCurrencyCode(documentCurrencyID).
				WithSupplier(getInvoiceSupplierParty()).
				WithCustomer(getInvoiceCustomerParty()).
				WithInvoiceLines(lines).
				AddTaxExemptionReason(TaxCategoryVATExempt, "MOTIVUL A", "")

			documentAllowance, ok := NewInvoiceDocumentAllowanceBuilder(
				documentCurrencyID,
				D(15),
				InvoiceTaxCategory{
					TaxScheme: TaxSchemeVAT,
					ID:        TaxCategoryVATStandardRate,
					Percent:   D(25),
				},
			).WithAllowanceChargeReason("Motivul C").Build()
			if assert.True(ok) {
				invoiceBuilder.AppendAllowanceCharge(documentAllowance)
			}

			documentCharge, ok := NewInvoiceDocumentChargeBuilder(
				documentCurrencyID,
				D(35),
				InvoiceTaxCategory{
					TaxScheme: TaxSchemeVAT,
					ID:        TaxCategoryVATStandardRate,
					Percent:   D(25),
				},
			).WithAllowanceChargeReason("Motivul B").Build()
			if assert.True(ok) {
				invoiceBuilder.AppendAllowanceCharge(documentCharge)
			}

			if documentCurrencyID != CurrencyRON {
				invoiceBuilder.WithTaxCurrencyCode(CurrencyRON)
				invoiceBuilder.WithDocumentToTaxCurrencyExchangeRate(D(4.9691))
			}

			return invoiceBuilder.Build()
		}

		{
			// BT-5 is RON
			invoice, ok := buildInvoice(CurrencyRON)
			if assert.True(ok) {
				// Invoice lines
				if assert.Equal(5, len(invoice.InvoiceLines)) {
					line1 := invoice.InvoiceLines[0]
					assert.Equal(a(D(125.0)), a(line1.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATStandardRate, line1.Item.TaxCategory.ID)
					assert.Equal(d(D(25.0)), d(line1.Item.TaxCategory.Percent))

					line2 := invoice.InvoiceLines[1]
					assert.Equal(a(D(24.0)), a(line2.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATStandardRate, line2.Item.TaxCategory.ID)
					assert.Equal(d(D(10.0)), d(line2.Item.TaxCategory.Percent))

					line3 := invoice.InvoiceLines[2]
					assert.Equal(a(D(136.0)), a(line3.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATStandardRate, line3.Item.TaxCategory.ID)
					assert.Equal(d(D(25.0)), d(line3.Item.TaxCategory.Percent))

					line4 := invoice.InvoiceLines[3]
					assert.Equal(a(D(95.0)), a(line4.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATExempt, line4.Item.TaxCategory.ID)
					assert.Equal(d(D(0.0)), d(line4.Item.TaxCategory.Percent))

					line5 := invoice.InvoiceLines[4]
					assert.Equal(a(D(53.0)), a(line5.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATExempt, line5.Item.TaxCategory.ID)
					assert.Equal(d(D(0.0)), d(line5.Item.TaxCategory.Percent))
				}

				// Document totals (BG-22)
				// BT-106
				assert.Equal(a(D(433.0)), a(invoice.LegalMonetaryTotal.LineExtensionAmount.Amount), "BT-106 incorrect value")
				// BT-107
				if allowanceTotalAmount := invoice.LegalMonetaryTotal.AllowanceTotalAmount; assert.NotNil(allowanceTotalAmount, "BT-107 must be non-nil") {
					assert.Equal(a(D(15)), a(allowanceTotalAmount.Amount), "BT-107 incorrect value")
				}
				// BT-108
				if chargeTotalAmount := invoice.LegalMonetaryTotal.ChargeTotalAmount; assert.NotNil(chargeTotalAmount, "BT-108 must be non-nil") {
					assert.Equal(a(D(35)), a(chargeTotalAmount.Amount), "BT-108 incorrect value")
				}
				// BT-109
				assert.Equal(a(D(453.0)), a(invoice.LegalMonetaryTotal.TaxExclusiveAmount.Amount), "BT-109 incorrect value")
				// BT-110
				if assert.Equal(1, len(invoice.TaxTotal), "Must have only one TaxTotal") &&
					assert.NotNil(invoice.TaxTotal[0].TaxAmount, "BT-110 must exist") {

					assert.Equal(a(D(72.65)), a(invoice.TaxTotal[0].TaxAmount.Amount), "BT-110 incorrect value")
				}
				// BT-112
				assert.Equal(a(D(525.65)), a(invoice.LegalMonetaryTotal.TaxInclusiveAmount.Amount), "BT-112 incorrect value")
				// BT-115
				assert.Equal(a(D(525.65)), a(invoice.LegalMonetaryTotal.PayableAmount.Amount), "BT-115 incorrect value")
			}
		}
		{
			// BT-5 is EUR, BT-6 is RON
			documentCurrencyID := CurrencyEUR
			invoice, ok := buildInvoice(documentCurrencyID)
			if assert.True(ok) {
				// Invoice lines
				if assert.Equal(5, len(invoice.InvoiceLines)) {
					line1 := invoice.InvoiceLines[0]
					assert.Equal(a(D(125.0)), a(line1.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATStandardRate, line1.Item.TaxCategory.ID)
					assert.Equal(d(D(25.0)), d(line1.Item.TaxCategory.Percent))

					line2 := invoice.InvoiceLines[1]
					assert.Equal(a(D(24.0)), a(line2.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATStandardRate, line2.Item.TaxCategory.ID)
					assert.Equal(d(D(10.0)), d(line2.Item.TaxCategory.Percent))

					line3 := invoice.InvoiceLines[2]
					assert.Equal(a(D(136.0)), a(line3.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATStandardRate, line3.Item.TaxCategory.ID)
					assert.Equal(d(D(25.0)), d(line3.Item.TaxCategory.Percent))

					line4 := invoice.InvoiceLines[3]
					assert.Equal(a(D(95.0)), a(line4.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATExempt, line4.Item.TaxCategory.ID)
					assert.Equal(d(D(0.0)), d(line4.Item.TaxCategory.Percent))

					line5 := invoice.InvoiceLines[4]
					assert.Equal(a(D(53.0)), a(line5.LineExtensionAmount.Amount))
					assert.Equal(TaxCategoryVATExempt, line5.Item.TaxCategory.ID)
					assert.Equal(d(D(0.0)), d(line5.Item.TaxCategory.Percent))
				}

				// Document totals (BG-22)
				// BT-106
				assert.Equal(a(D(433.0)), a(invoice.LegalMonetaryTotal.LineExtensionAmount.Amount), "BT-106 incorrect value")
				// BT-107
				if allowanceTotalAmount := invoice.LegalMonetaryTotal.AllowanceTotalAmount; assert.NotNil(allowanceTotalAmount, "BT-107 must be non-nil") {
					assert.Equal(a(D(15)), a(allowanceTotalAmount.Amount), "BT-107 incorrect value")
				}
				// BT-108
				if chargeTotalAmount := invoice.LegalMonetaryTotal.ChargeTotalAmount; assert.NotNil(chargeTotalAmount, "BT-108 must be non-nil") {
					assert.Equal(a(D(35)), a(chargeTotalAmount.Amount), "BT-108 incorrect value")
				}
				// BT-109
				assert.Equal(a(D(453.0)), a(invoice.LegalMonetaryTotal.TaxExclusiveAmount.Amount), "BT-109 incorrect value")
				// BT-110
				if assert.Equal(2, len(invoice.TaxTotal), "Must have a TaxTotal for each currency") {
					taxTotalID := findTaxTotalByCurrency(invoice.TaxTotal, documentCurrencyID)
					if assert.True(taxTotalID >= 0, "TaxTotal for document currency code(BT-5) not set") &&
						assert.NotNil(invoice.TaxTotal[taxTotalID].TaxAmount, "BT-110 must be non-nil") {
						assert.Equal(a(D(72.65)), a(invoice.TaxTotal[taxTotalID].TaxAmount.Amount), "BT-110 incorrect value")
					}

					taxTotalTaxCurrencyID := findTaxTotalByCurrency(invoice.TaxTotal, CurrencyRON)
					if assert.True(taxTotalTaxCurrencyID >= 0, "TaxTotal for tax currency code(BT-6) not set") &&
						assert.NotNil(invoice.TaxTotal[taxTotalTaxCurrencyID].TaxAmount, "BT-111 must be non-nil") {
						assert.Equal(a(D(361.01)), a(invoice.TaxTotal[taxTotalTaxCurrencyID].TaxAmount.Amount), "BT-110 incorrect value")
					}
				}
				// BT-112
				assert.Equal(a(D(525.65)), a(invoice.LegalMonetaryTotal.TaxInclusiveAmount.Amount), "BT-112 incorrect value")
				// BT-115
				assert.Equal(a(D(525.65)), a(invoice.LegalMonetaryTotal.PayableAmount.Amount), "BT-115 incorrect value")
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

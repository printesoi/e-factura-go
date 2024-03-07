package efactura

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvoiceLineBuilder(t *testing.T) {
	assert := assert.New(t)

	{
		b := NewInvoiceLineBuilder("1", CurrencyRON)
		_, ok := b.Build()
		assert.False(ok, "should not build if required fields are missing")
	}
	{
		// A.1.4 Exemplul 1
		{
			// Linia 1
			b := NewInvoiceLineBuilder("1", CurrencyRON).
				WithUnitCode("XGR").
				WithInvoicedQuantity(D(5)).
				WithGrossPriceAmount(D(12)).
				WithItem(InvoiceLineItem{
					Name: "Sticle cu vin",
					TaxCategory: InvoiceLineTaxCategory{
						ID:          TaxCategoryTVACotaNormalaRedusa,
						Percent:     D(25),
						TaxSchemeID: TaxSchemeVAT,
					},
				})
			line, ok := b.Build()
			if assert.True(ok) {
				assert.Equal(D(60).String(), line.LineExtensionAmount.Amount.String())
			}
		}
		{
			// Linia 2
			b := NewInvoiceLineBuilder("1", CurrencyRON).
				WithUnitCode("XBX").
				WithInvoicedQuantity(D(1)).
				WithGrossPriceAmount(D(90)).
				WithItem(InvoiceLineItem{
					Name: "Vin - cutie de 6",
					TaxCategory: InvoiceLineTaxCategory{
						ID:          TaxCategoryTVACotaNormalaRedusa,
						Percent:     D(12),
						TaxSchemeID: TaxSchemeVAT,
					},
				})
			line, ok := b.Build()
			if assert.True(ok) {
				assert.Equal(D(90).String(), line.LineExtensionAmount.Amount.String())
			}
		}
	}
	{
		// A.1.4 Exemplul 2
		b := NewInvoiceLineBuilder("1", CurrencyRON).
			WithUnitCode("C62").
			WithInvoicedQuantity(D(10_000)).
			WithGrossPriceAmount(D(4.5)).
			WithBaseQuantity(D(1_000)).
			WithItem(InvoiceLineItem{
				Name: "șurub",
				TaxCategory: InvoiceLineTaxCategory{
					ID:          TaxCategoryTVACotaNormalaRedusa,
					Percent:     D(25),
					TaxSchemeID: TaxSchemeVAT,
				},
			})
		line, ok := b.Build()
		if assert.True(ok) {
			assert.Equal(D(45.0).String(), line.LineExtensionAmount.Amount.String())
		}
	}
	{
		// A.1.4 Exemplul 3
		b := NewInvoiceLineBuilder("1", CurrencyRON).
			WithUnitCode("58").
			WithInvoicedQuantity(D(1.3)).
			WithGrossPriceAmount(D(10)).
			WithPriceDeduction(D(0.5)).
			WithItem(InvoiceLineItem{
				Name: "Pui",
				TaxCategory: InvoiceLineTaxCategory{
					ID:          TaxCategoryTVACotaNormalaRedusa,
					Percent:     D(12.5),
					TaxSchemeID: TaxSchemeVAT,
				},
			})
		line, ok := b.Build()
		if assert.True(ok) {
			assert.Equal(D(12.35).String(), line.LineExtensionAmount.Amount.String())
		}
	}
	{
		// A.1.5 Exemplul 4 (Reduceri, deduceri şi taxe suplimentare)
		{
			// Line 1
			currencyID := CurrencyRON
			b := NewInvoiceLineBuilder("1", currencyID).
				WithUnitCode("XBX").
				WithInvoicedQuantity(D(25)).
				WithGrossPriceAmount(D(9.5)).
				WithPriceDeduction(D(1.0)).
				WithItem(InvoiceLineItem{
					Name: "Stilou",
					TaxCategory: InvoiceLineTaxCategory{
						ID:          TaxCategoryTVACotaNormalaRedusa,
						Percent:     D(25),
						TaxSchemeID: TaxSchemeVAT,
					},
				})
			lineCharge, ok := NewInvoiceLineChargeBuilder(currencyID, D(10)).Build()
			assert.True(ok)
			b.AppendAllowanceCharge(lineCharge)
			line, ok := b.Build()
			if assert.True(ok) {
				assert.Equal(D(222.50).String(), line.LineExtensionAmount.Amount.String())
			}
		}
		{
			// Line 2
			currencyID := CurrencyRON
			b := NewInvoiceLineBuilder("2", currencyID).
				WithUnitCode("RM").
				WithInvoicedQuantity(D(15)).
				WithGrossPriceAmount(D(4.50)).
				WithItem(InvoiceLineItem{
					Name: "Hârtie",
					TaxCategory: InvoiceLineTaxCategory{
						ID:          TaxCategoryTVACotaNormalaRedusa,
						Percent:     D(25),
						TaxSchemeID: TaxSchemeVAT,
					},
				})
			lineCharge, ok := NewInvoiceLineAllowanceBuilder(currencyID, D(3.38)).Build()
			assert.True(ok)
			b.AppendAllowanceCharge(lineCharge)
			line, ok := b.Build()
			if assert.True(ok) {
				assert.Equal(D(64.12).String(), line.LineExtensionAmount.Amount.String())
			}
		}
	}
}

func TestInvoiceBuilder(t *testing.T) {
	assert := assert.New(t)

	{
		b := NewInvoiceBuilder("1")
		_, ok := b.Build()
		assert.False(ok, "should not build if required fields are missing")
	}
	{
		// A.1.8 Exemplul 7 (Cota normală de TVA cu linii scutite de TVA)
		documentCurrencyID := CurrencyRON

		var lines []InvoiceLine

		line1, ok := NewInvoiceLineBuilder("1", documentCurrencyID).
			WithUnitCode("H87").
			WithInvoicedQuantity(D(1)).
			WithGrossPriceAmount(D(125.0)).
			WithItem(InvoiceLineItem{
				Name: "Cerneală pentru imprimantă",
				TaxCategory: InvoiceLineTaxCategory{
					ID:          TaxCategoryTVACotaNormalaRedusa,
					Percent:     D(25),
					TaxSchemeID: TaxSchemeVAT,
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
				Name: "Imprimare afiș",
				TaxCategory: InvoiceLineTaxCategory{
					ID:          TaxCategoryTVACotaNormalaRedusa,
					Percent:     D(10),
					TaxSchemeID: TaxSchemeVAT,
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
				Name: "Scaun de birou",
				TaxCategory: InvoiceLineTaxCategory{
					ID:          TaxCategoryTVACotaNormalaRedusa,
					Percent:     D(25),
					TaxSchemeID: TaxSchemeVAT,
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
				Name: "Tastatură fără fir",
				TaxCategory: InvoiceLineTaxCategory{
					ID:          TaxCategoryScutireTVA,
					TaxSchemeID: TaxSchemeVAT,
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
				Name: "Cablu de adaptare",
				TaxCategory: InvoiceLineTaxCategory{
					ID:          TaxCategoryScutireTVA,
					TaxSchemeID: TaxSchemeVAT,
				},
			}).
			Build()
		if assert.True(ok) {
			lines = append(lines, line5)
		}

		invoiceBuilder := NewInvoiceBuilder("test01").
			WithIssueDate(MakeDateLocal(2024, 3, 7)).
			WithDueDate(MakeDateLocal(2024, 4, 7)).
			WithInvoiceTypeCode(InvoiceTypeCommercialInvoice).
			WithDocumentCurrencyCode(documentCurrencyID).
			WithSupplier(InvoiceSupplier{
				PostalAddress: MakeInvoiceSupplierPostalAddress(PostalAddress{
					CountryIdentificationCode: CountryCodeRO,
					CountrySubentity:          CountrySubentityRO_B,
					CityName:                  CityNameROBSector6,
					Line1:                     "B-DUL IULIU MANIU, NR.6E, PARTER,CAMERA 1, SC.1, AP.3, SECTOR 6",
				}),
				TaxScheme: &InvoicePartyTaxScheme{
					CompanyID:   "RO34283300",
					TaxSchemeID: TaxSchemeVAT,
				},
				LegalEntity: InvoiceSupplierLegalEntity{
					Name:      "FACTURIS ONLINE SRL",
					CompanyID: MakeValueWithAttrs("34283300").Ptr(),
				},
			}).
			WithCustomer(InvoiceCustomer{
				PostalAddress: MakeInvoiceCustomerPostalAddress(PostalAddress{
					CountryIdentificationCode: CountryCodeRO,
					CountrySubentity:          CountrySubentityRO_B,
					CityName:                  CityNameROBSector6,
					Line1:                     "B-DUL IULIU MANIU, NR.6E, PARTER,CAMERA 1, SC.1, AP.3, SECTOR 6",
				}),
				TaxScheme: &InvoicePartyTaxScheme{
					CompanyID:   "RO19211548",
					TaxSchemeID: TaxSchemeVAT,
				},
				LegalEntity: InvoiceCustomerLegalEntity{
					Name:      "MIDSOFT IT GROUP SRL",
					CompanyID: MakeValueWithAttrs("19211548").Ptr(),
				},
			}).
			WithInvoiceLines(lines).
			AddTaxExemptionReason(TaxCategoryScutireTVA, "MOTIVUL A", "")

		documentAllowance, ok := NewInvoiceDocumentAllowanceBuilder(documentCurrencyID,
			D(15), InvoiceTaxCategory{
				ID:          TaxCategoryTVACotaNormalaRedusa,
				Percent:     D(25),
				TaxSchemeID: TaxSchemeVAT,
			}).
			WithAllowanceChargeReason("Motivul C").
			Build()
		if assert.True(ok) {
			invoiceBuilder.AppendAllowanceCharge(documentAllowance)
		}

		documentCharge, ok := NewInvoiceDocumentChargeBuilder(documentCurrencyID,
			D(35), InvoiceTaxCategory{
				ID:          TaxCategoryTVACotaNormalaRedusa,
				Percent:     D(25),
				TaxSchemeID: TaxSchemeVAT,
			}).
			WithAllowanceChargeReason("Motivul B").
			Build()
		if assert.True(ok) {
			invoiceBuilder.AppendAllowanceCharge(documentCharge)
		}

		invoice, ok := invoiceBuilder.Build()
		if assert.True(ok) {
			assert.Equal(5, len(invoice.InvoiceLines))
		}
	}
}

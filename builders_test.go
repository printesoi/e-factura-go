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
					ClassifiedTaxCategory: InvoiceClassifiedTaxCategory{
						ID:          TaxCategoryTVACotaNormalaRedusa,
						Percent:     D(25).Ptr(),
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
					ClassifiedTaxCategory: InvoiceClassifiedTaxCategory{
						ID:          TaxCategoryTVACotaNormalaRedusa,
						Percent:     D(12).Ptr(),
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
				ClassifiedTaxCategory: InvoiceClassifiedTaxCategory{
					ID:          TaxCategoryTVACotaNormalaRedusa,
					Percent:     D(25).Ptr(),
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
				ClassifiedTaxCategory: InvoiceClassifiedTaxCategory{
					ID:          TaxCategoryTVACotaNormalaRedusa,
					Percent:     D(12.5).Ptr(),
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
					ClassifiedTaxCategory: InvoiceClassifiedTaxCategory{
						ID:          TaxCategoryTVACotaNormalaRedusa,
						Percent:     D(25).Ptr(),
						TaxSchemeID: TaxSchemeVAT,
					},
				})
			lineCharge, ok := NewInvoiceChargeBuilder(currencyID, D(10)).Build()
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
					ClassifiedTaxCategory: InvoiceClassifiedTaxCategory{
						ID:          TaxCategoryTVACotaNormalaRedusa,
						Percent:     D(25).Ptr(),
						TaxSchemeID: TaxSchemeVAT,
					},
				})
			lineCharge, ok := NewInvoiceAllowanceBuilder(currencyID, D(3.38)).Build()
			assert.True(ok)
			b.AppendAllowanceCharge(lineCharge)
			line, ok := b.Build()
			if assert.True(ok) {
				assert.Equal(D(64.12).String(), line.LineExtensionAmount.Amount.String())
			}
		}
	}
}

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

// InvoiceAllowanceChargeBuilder builds an InvoiceAllowanceCharge object
type InvoiceAllowanceChargeBuilder struct {
	chargeIndicator           bool
	currencyID                CurrencyCodeType
	amount                    Decimal
	baseAmount                *Decimal
	allowanceChargeReasonCode *string
	allowanceChargeReason     *string
}

func NewInvoiceAllowanceChargeBuilder(chargeIndicator bool, currencyID CurrencyCodeType, amount Decimal) *InvoiceAllowanceChargeBuilder {
	f := new(InvoiceAllowanceChargeBuilder)
	return f.WithChargeIndicator(chargeIndicator).
		WithCurrencyID(currencyID).WithAmount(amount)
}

func NewInvoiceAllowanceBuilder(currencyID CurrencyCodeType, amount Decimal) *InvoiceAllowanceChargeBuilder {
	return NewInvoiceAllowanceChargeBuilder(false, currencyID, amount)
}

func NewInvoiceChargeBuilder(currencyID CurrencyCodeType, amount Decimal) *InvoiceAllowanceChargeBuilder {
	return NewInvoiceAllowanceChargeBuilder(true, currencyID, amount)
}

func (f *InvoiceAllowanceChargeBuilder) WithChargeIndicator(charge bool) *InvoiceAllowanceChargeBuilder {
	f.chargeIndicator = charge
	return f
}

func (f *InvoiceAllowanceChargeBuilder) WithCurrencyID(currencyID CurrencyCodeType) *InvoiceAllowanceChargeBuilder {
	f.currencyID = currencyID
	return f
}

func (f *InvoiceAllowanceChargeBuilder) WithAmount(amount Decimal) *InvoiceAllowanceChargeBuilder {
	f.amount = amount
	return f
}

func (f *InvoiceAllowanceChargeBuilder) WithBaseAmount(amount Decimal) *InvoiceAllowanceChargeBuilder {
	f.baseAmount = amount.Ptr()
	return f
}

func (f *InvoiceAllowanceChargeBuilder) WithAllowanceChargeReasonCode(allowanceChargeReasonCode string) *InvoiceAllowanceChargeBuilder {
	f.allowanceChargeReasonCode = ptrfyString(allowanceChargeReasonCode)
	return f
}

func (f *InvoiceAllowanceChargeBuilder) WithAllowanceChargeReason(allowanceChargeReason string) *InvoiceAllowanceChargeBuilder {
	f.allowanceChargeReason = ptrfyString(allowanceChargeReason)
	return f
}

func (f *InvoiceAllowanceChargeBuilder) Build() (InvoiceAllowanceCharge, bool) {
	m := InvoiceAllowanceCharge{
		ChargeIndicator: f.chargeIndicator,
		Amount: AmountWithCurrency{
			Amount:     f.amount,
			CurrencyID: f.currencyID,
		},
	}
	if f.baseAmount != nil {
		m.BaseAmount = &AmountWithCurrency{
			Amount:     *f.baseAmount,
			CurrencyID: f.currencyID,
		}
	}
	if f.allowanceChargeReasonCode != nil {
		m.AllowanceChargeReasonCode = *f.allowanceChargeReasonCode
	}
	if f.allowanceChargeReason != nil {
		m.AllowanceChargeReason = *f.allowanceChargeReason
	}
	return m, true
}

// InvoiceLineBuilder builds an InvoiceLine object
type InvoiceLineBuilder struct {
	id               string
	note             string
	currencyID       CurrencyCodeType
	unitCode         UnitCodeType
	invoicedQuantity Decimal
	baseQuantity     *Decimal

	grossPriceAmount Decimal
	itemAllowance    Decimal

	invoicePeriod *InvoiceLinePeriod
	lineAllowance *InvoiceAllowanceCharge
	lineCharge    *InvoiceAllowanceCharge
	item          InvoiceLineItem
}

func NewInvoiceLineBuilder(id string, currencyID CurrencyCodeType) (b *InvoiceLineBuilder) {
	b = new(InvoiceLineBuilder)
	return b.WithID(id)
}

func (b *InvoiceLineBuilder) WithID(id string) *InvoiceLineBuilder {
	b.id = id
	return b
}

func (b *InvoiceLineBuilder) WithCurrencyID(currencyID CurrencyCodeType) *InvoiceLineBuilder {
	b.currencyID = currencyID
	return b
}

func (b *InvoiceLineBuilder) WithNote(note string) *InvoiceLineBuilder {
	b.note = note
	return b
}

func (b *InvoiceLineBuilder) WithUnitCode(unitCode UnitCodeType) *InvoiceLineBuilder {
	b.unitCode = unitCode
	return b
}

func (b *InvoiceLineBuilder) WithInvoicedQuantity(quantity Decimal) *InvoiceLineBuilder {
	b.invoicedQuantity = quantity
	return b
}

func (b *InvoiceLineBuilder) WithBaseQuantity(quantity Decimal) *InvoiceLineBuilder {
	b.baseQuantity = &quantity
	return b
}

func (b *InvoiceLineBuilder) WithGrossPriceAmount(priceAmount Decimal) *InvoiceLineBuilder {
	b.grossPriceAmount = priceAmount
	return b
}

func (b *InvoiceLineBuilder) WithPriceDeduction(allowance Decimal) *InvoiceLineBuilder {
	b.itemAllowance = allowance
	return b
}

func (b *InvoiceLineBuilder) WithInvoicePeriod(invoicePeriod *InvoiceLinePeriod) *InvoiceLineBuilder {
	b.invoicePeriod = invoicePeriod
	return b
}

func (b *InvoiceLineBuilder) WithAllowance(allowance InvoiceAllowanceCharge) *InvoiceLineBuilder {
	b.lineAllowance = &allowance
	return b
}

func (b *InvoiceLineBuilder) WithCharge(charge InvoiceAllowanceCharge) *InvoiceLineBuilder {
	b.lineCharge = &charge
	return b
}

func (b *InvoiceLineBuilder) WithItem(item InvoiceLineItem) *InvoiceLineBuilder {
	b.item = item
	return b
}

func (b *InvoiceLineBuilder) Build() (line InvoiceLine, ok bool) {
	if b.id == "" || !b.invoicedQuantity.IsInitialized() ||
		b.unitCode == "" || !b.grossPriceAmount.IsInitialized() ||
		b.item.Name == "" || b.item.ClassifiedTaxCategory.ID == "" ||
		b.item.ClassifiedTaxCategory.TaxSchemeID == "" {
		return
	}

	line.ID = b.id
	line.Note = b.note
	line.InvoicedQuantity = InvoicedQuantity{
		Quantity: b.invoicedQuantity,
		UnitCode: b.unitCode,
	}
	var netPriceAmount Decimal
	if b.itemAllowance.IsZero() {
		netPriceAmount = b.grossPriceAmount
	} else {
		netPriceAmount = b.grossPriceAmount.Sub(b.itemAllowance)
		line.Price.PriceAmount = AmountWithCurrency{
			Amount:     netPriceAmount,
			CurrencyID: b.currencyID,
		}
		line.Price.AllowanceCharge = &InvoiceLinePriceAllowanceCharge{
			ChargeIndicator: false,
			Amount: AmountWithCurrency{
				Amount:     b.itemAllowance,
				CurrencyID: b.currencyID,
			},
			BaseAmount: AmountWithCurrency{
				Amount:     b.grossPriceAmount,
				CurrencyID: b.currencyID,
			},
		}
	}
	line.Price.PriceAmount = AmountWithCurrency{
		Amount:     netPriceAmount,
		CurrencyID: b.currencyID,
	}
	if b.baseQuantity != nil {
		line.Price.BaseQuantity = &InvoicedQuantity{
			Quantity: *b.baseQuantity,
			UnitCode: b.unitCode,
		}
	}
	line.InvoicePeriod = b.invoicePeriod
	if b.lineAllowance != nil {
		line.AllowanceCharges = append(line.AllowanceCharges, *b.lineAllowance)
	}
	if b.lineCharge != nil {
		line.AllowanceCharges = append(line.AllowanceCharges, *b.lineCharge)
	}
	line.Item = b.item

	// Invoiced quantity * (Item net price / item price base quantity)
	//  + Sum of invoice line charge amount
	//  - Sum of invoice line allowance amount
	baseQuantity := D(1)
	if b.baseQuantity != nil {
		baseQuantity = *b.baseQuantity
	}
	if baseQuantity.IsZero() {
		return line, false
	}
	netAmount := b.invoicedQuantity.Mul(netPriceAmount).Div(baseQuantity)
	for _, charge := range line.AllowanceCharges {
		if charge.ChargeIndicator {
			netAmount = netAmount.Add(charge.Amount.Amount)
		} else {
			netAmount = netAmount.Sub(charge.Amount.Amount)
		}
	}

	line.LineExtensionAmount = AmountWithCurrency{
		Amount:     netAmount.AsAmount(),
		CurrencyID: b.currencyID,
	}
	ok = true
	return
}

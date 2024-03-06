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
		WithCurrencyID(currencyID)
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
	priceAmount      Decimal
	invoicePeriod    *InvoiceLinePeriod
	allowanceCharges []InvoiceAllowanceCharge
	item             InvoiceLineItem
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

func (b *InvoiceLineBuilder) WithPriceAmount(priceAmount Decimal) *InvoiceLineBuilder {
	b.priceAmount = priceAmount
	return b
}

func (b *InvoiceLineBuilder) WithInvoicePeriod(invoicePeriod *InvoiceLinePeriod) *InvoiceLineBuilder {
	b.invoicePeriod = invoicePeriod
	return b
}

func (b *InvoiceLineBuilder) WithAllowanceCharges(allowanceCharges []InvoiceAllowanceCharge) *InvoiceLineBuilder {
	b.allowanceCharges = allowanceCharges
	return b
}

func (b *InvoiceLineBuilder) AppendAllowanceCharges(allowanceCharge InvoiceAllowanceCharge) *InvoiceLineBuilder {
	return b.WithAllowanceCharges(append(b.allowanceCharges, allowanceCharge))
}

func (b *InvoiceLineBuilder) WithItem(item InvoiceLineItem) *InvoiceLineBuilder {
	b.item = item
	return b
}

func (b *InvoiceLineBuilder) Build() (line InvoiceLine, ok bool) {
	line.ID = b.id
	line.Note = b.note
	line.InvoicedQuantity = InvoicedQuantity{
		Quantity: b.invoicedQuantity,
		UnitCode: b.unitCode,
	}
	line.Price.PriceAmount = AmountWithCurrency{
		Amount:     b.priceAmount,
		CurrencyID: b.currencyID,
	}
	if b.baseQuantity != nil {
		line.Price.BaseQuantity = &InvoicedQuantity{
			Quantity: *b.baseQuantity,
			UnitCode: b.unitCode,
		}
	}
	line.InvoicePeriod = b.invoicePeriod
	line.AllowanceCharge = b.allowanceCharges
	line.Item = b.item

	// TODO: compute net amount
	// (Invoiced quantity * (Item net price/item price base quantity) + Sum of invoice line charge amount - sum of invoice line allowance amount
	netAmount := line.InvoicedQuantity.Quantity
	line.LineExtensionAmount = AmountWithCurrency{
		Amount:     netAmount,
		CurrencyID: b.currencyID,
	}
	return
}

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

// InvoiceLineAllowanceChargeBuilder builds an InvoiceLineAllowanceCharge object
type InvoiceLineAllowanceChargeBuilder struct {
	chargeIndicator           bool
	currencyID                CurrencyCodeType
	amount                    Decimal
	baseAmount                *Decimal
	allowanceChargeReasonCode *string
	allowanceChargeReason     *string
}

// NewInvoiceLineAllowanceChargeBuilder creates a new generic
// InvoiceLineAllowanceChargeBuilder.
func NewInvoiceLineAllowanceChargeBuilder(chargeIndicator bool, currencyID CurrencyCodeType, amount Decimal) *InvoiceLineAllowanceChargeBuilder {
	b := new(InvoiceLineAllowanceChargeBuilder)
	return b.WithChargeIndicator(chargeIndicator).
		WithCurrencyID(currencyID).WithAmount(amount)
}

// NewInvoiceAllowanceBuilder creates a new InvoiceLineAllowanceChargeBuilder
// builder that will build InvoiceLineAllowanceCharge object correspoding to an
// allowance (ChargeIndicator = false)
func NewInvoiceAllowanceBuilder(currencyID CurrencyCodeType, amount Decimal) *InvoiceLineAllowanceChargeBuilder {
	return NewInvoiceLineAllowanceChargeBuilder(false, currencyID, amount)
}

// NewInvoiceChargeBuilder creates a new InvoiceLineAllowanceChargeBuilder
// builder that will build InvoiceLineAllowanceCharge object correspoding to a
// charge (ChargeIndicator = true)
func NewInvoiceChargeBuilder(currencyID CurrencyCodeType, amount Decimal) *InvoiceLineAllowanceChargeBuilder {
	return NewInvoiceLineAllowanceChargeBuilder(true, currencyID, amount)
}

func (b *InvoiceLineAllowanceChargeBuilder) WithChargeIndicator(charge bool) *InvoiceLineAllowanceChargeBuilder {
	b.chargeIndicator = charge
	return b
}

func (b *InvoiceLineAllowanceChargeBuilder) WithCurrencyID(currencyID CurrencyCodeType) *InvoiceLineAllowanceChargeBuilder {
	b.currencyID = currencyID
	return b
}

func (b *InvoiceLineAllowanceChargeBuilder) WithAmount(amount Decimal) *InvoiceLineAllowanceChargeBuilder {
	b.amount = amount
	return b
}

func (b *InvoiceLineAllowanceChargeBuilder) WithBaseAmount(amount Decimal) *InvoiceLineAllowanceChargeBuilder {
	b.baseAmount = amount.Ptr()
	return b
}

func (b *InvoiceLineAllowanceChargeBuilder) WithAllowanceChargeReasonCode(allowanceChargeReasonCode string) *InvoiceLineAllowanceChargeBuilder {
	b.allowanceChargeReasonCode = ptrfyString(allowanceChargeReasonCode)
	return b
}

func (b *InvoiceLineAllowanceChargeBuilder) WithAllowanceChargeReason(allowanceChargeReason string) *InvoiceLineAllowanceChargeBuilder {
	b.allowanceChargeReason = ptrfyString(allowanceChargeReason)
	return b
}

func (b *InvoiceLineAllowanceChargeBuilder) Build() (allowanceCharge InvoiceLineAllowanceCharge, ok bool) {
	if !b.amount.IsInitialized() || b.currencyID == "" {
		return
	}
	allowanceCharge.ChargeIndicator = b.chargeIndicator
	allowanceCharge.Amount = AmountWithCurrency{
		Amount:     b.amount,
		CurrencyID: b.currencyID,
	}
	if b.baseAmount != nil {
		allowanceCharge.BaseAmount = &AmountWithCurrency{
			Amount:     *b.baseAmount,
			CurrencyID: b.currencyID,
		}
	}
	if b.allowanceChargeReasonCode != nil {
		allowanceCharge.AllowanceChargeReasonCode = *b.allowanceChargeReasonCode
	}
	if b.allowanceChargeReason != nil {
		allowanceCharge.AllowanceChargeReason = *b.allowanceChargeReason
	}
	ok = true
	return
}

// InvoiceLineBuilder builds an InvoiceLine object. The only (useful) role of
// this builder is to help build a complex InvoiceLine object while ensuring
// the amounts are calculated correctly.
type InvoiceLineBuilder struct {
	id               string
	note             string
	currencyID       CurrencyCodeType
	unitCode         UnitCodeType
	invoicedQuantity Decimal
	baseQuantity     *Decimal

	grossPriceAmount Decimal
	itemAllowance    Decimal

	invoicePeriod     *InvoiceLinePeriod
	allowancesCharges []InvoiceLineAllowanceCharge
	item              InvoiceLineItem
}

// NewInvoiceLineBuilder creates a new InvoiceLineBuilder
func NewInvoiceLineBuilder(id string, currencyID CurrencyCodeType) (b *InvoiceLineBuilder) {
	b = new(InvoiceLineBuilder)
	return b.WithID(id).WithCurrencyID(currencyID)
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

func (b *InvoiceLineBuilder) WithAllowancesCharges(allowancesCharges []InvoiceLineAllowanceCharge) *InvoiceLineBuilder {
	b.allowancesCharges = allowancesCharges
	return b
}

func (b *InvoiceLineBuilder) AppendAllowanceCharge(allowanceCharge InvoiceLineAllowanceCharge) *InvoiceLineBuilder {
	return b.WithAllowancesCharges(append(b.allowancesCharges, allowanceCharge))
}

func (b *InvoiceLineBuilder) WithItem(item InvoiceLineItem) *InvoiceLineBuilder {
	b.item = item
	return b
}

func (b *InvoiceLineBuilder) Build() (line InvoiceLine, ok bool) {
	if b.id == "" || b.currencyID == "" ||
		!b.invoicedQuantity.IsInitialized() ||
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
	line.Item = b.item
	line.AllowanceCharges = b.allowancesCharges
	line.InvoicePeriod = b.invoicePeriod

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

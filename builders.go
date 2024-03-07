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

// NewInvoiceLineAllowanceBuilder creates a new InvoiceLineAllowanceChargeBuilder
// builder that will build InvoiceLineAllowanceCharge object correspoding to an
// allowance (ChargeIndicator = false)
func NewInvoiceLineAllowanceBuilder(currencyID CurrencyCodeType, amount Decimal) *InvoiceLineAllowanceChargeBuilder {
	return NewInvoiceLineAllowanceChargeBuilder(false, currencyID, amount)
}

// NewInvoiceLineChargeBuilder creates a new InvoiceLineAllowanceChargeBuilder
// builder that will build InvoiceLineAllowanceCharge object correspoding to a
// charge (ChargeIndicator = true)
func NewInvoiceLineChargeBuilder(currencyID CurrencyCodeType, amount Decimal) *InvoiceLineAllowanceChargeBuilder {
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
		b.item.Name == "" || b.item.TaxCategory.ID == "" ||
		b.item.TaxCategory.TaxSchemeID == "" {
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

// InvoiceDocumentAllowanceChargeBuilder builds an InvoiceDocumentAllowanceCharge object
type InvoiceDocumentAllowanceChargeBuilder struct {
	chargeIndicator           bool
	currencyID                CurrencyCodeType
	amount                    Decimal
	taxCategory               InvoiceTaxCategory
	baseAmount                *Decimal
	allowanceChargeReasonCode *string
	allowanceChargeReason     *string
}

// NewInvoiceDocumentAllowanceChargeBuilder creates a new generic
// InvoiceDocumentAllowanceChargeBuilder.
func NewInvoiceDocumentAllowanceChargeBuilder(chargeIndicator bool, currencyID CurrencyCodeType, amount Decimal, taxCategory InvoiceTaxCategory) *InvoiceDocumentAllowanceChargeBuilder {
	b := new(InvoiceDocumentAllowanceChargeBuilder)
	return b.WithChargeIndicator(chargeIndicator).WithCurrencyID(currencyID).
		WithAmount(amount).WithTaxCategory(taxCategory)
}

// NewInvoiceDocumentAllowanceBuilder creates a new InvoiceDocumentAllowanceChargeBuilder
// builder that will build InvoiceDocumentAllowanceCharge object correspoding to an
// allowance (ChargeIndicator = false)
func NewInvoiceDocumentAllowanceBuilder(currencyID CurrencyCodeType, amount Decimal, taxCategory InvoiceTaxCategory) *InvoiceDocumentAllowanceChargeBuilder {
	return NewInvoiceDocumentAllowanceChargeBuilder(false, currencyID, amount, taxCategory)
}

// NewInvoiceDocumentChargeBuilder creates a new InvoiceDocumentAllowanceChargeBuilder
// builder that will build InvoiceDocumentAllowanceCharge object correspoding to a
// charge (ChargeIndicator = true)
func NewInvoiceDocumentChargeBuilder(currencyID CurrencyCodeType, amount Decimal, taxCategory InvoiceTaxCategory) *InvoiceDocumentAllowanceChargeBuilder {
	return NewInvoiceDocumentAllowanceChargeBuilder(true, currencyID, amount, taxCategory)
}

func (b *InvoiceDocumentAllowanceChargeBuilder) WithChargeIndicator(charge bool) *InvoiceDocumentAllowanceChargeBuilder {
	b.chargeIndicator = charge
	return b
}

func (b *InvoiceDocumentAllowanceChargeBuilder) WithCurrencyID(currencyID CurrencyCodeType) *InvoiceDocumentAllowanceChargeBuilder {
	b.currencyID = currencyID
	return b
}

func (b *InvoiceDocumentAllowanceChargeBuilder) WithAmount(amount Decimal) *InvoiceDocumentAllowanceChargeBuilder {
	b.amount = amount
	return b
}

func (b *InvoiceDocumentAllowanceChargeBuilder) WithTaxCategory(taxCategory InvoiceTaxCategory) *InvoiceDocumentAllowanceChargeBuilder {
	b.taxCategory = taxCategory
	return b
}

func (b *InvoiceDocumentAllowanceChargeBuilder) WithBaseAmount(amount Decimal) *InvoiceDocumentAllowanceChargeBuilder {
	b.baseAmount = amount.Ptr()
	return b
}

func (b *InvoiceDocumentAllowanceChargeBuilder) WithAllowanceChargeReasonCode(allowanceChargeReasonCode string) *InvoiceDocumentAllowanceChargeBuilder {
	b.allowanceChargeReasonCode = ptrfyString(allowanceChargeReasonCode)
	return b
}

func (b *InvoiceDocumentAllowanceChargeBuilder) WithAllowanceChargeReason(allowanceChargeReason string) *InvoiceDocumentAllowanceChargeBuilder {
	b.allowanceChargeReason = ptrfyString(allowanceChargeReason)
	return b
}

func (b *InvoiceDocumentAllowanceChargeBuilder) Build() (allowanceCharge InvoiceDocumentAllowanceCharge, ok bool) {
	if !b.amount.IsInitialized() || b.currencyID == "" ||
		b.taxCategory.ID == "" || b.taxCategory.TaxSchemeID == "" {
		return
	}
	allowanceCharge.ChargeIndicator = b.chargeIndicator
	allowanceCharge.Amount = AmountWithCurrency{
		Amount:     b.amount,
		CurrencyID: b.currencyID,
	}
	allowanceCharge.TaxCategory = b.taxCategory
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

type taxExemptionReason struct {
	reason string
	code   TaxExemptionReasonCodeType
}

// InvoiceBuilder builds an Invoice object
type InvoiceBuilder struct {
	id          string
	issueDate   Date
	dueDate     *Date
	invoiceType InvoiceTypeCodeType

	documentCurrencyID      CurrencyCodeType
	taxCurrencyID           CurrencyCodeType
	taxCurrencyExchangeRate Decimal

	taxExeptionReasons map[TaxCategoryCodeType]taxExemptionReason

	billingReferences []InvoiceBillingReference
	supplier          InvoiceSupplier
	customer          InvoiceCustomer

	allowancesCharges []InvoiceDocumentAllowanceCharge
	invoiceLines      []InvoiceLine
}

func NewInvoiceBuilder(id string) (b *InvoiceBuilder) {
	b = new(InvoiceBuilder)
	return b.WithID(id)
}

func (b *InvoiceBuilder) WithID(id string) *InvoiceBuilder {
	b.id = id
	return b
}

func (b *InvoiceBuilder) WithIssueDate(date Date) *InvoiceBuilder {
	b.issueDate = date
	return b
}

func (b *InvoiceBuilder) WithDueDate(date Date) *InvoiceBuilder {
	b.dueDate = &date
	return b
}

func (b *InvoiceBuilder) WithInvoiceTypeCode(invoiceType InvoiceTypeCodeType) *InvoiceBuilder {
	b.invoiceType = invoiceType
	return b
}

func (b *InvoiceBuilder) WithDocumentCurrencyCode(currencyID CurrencyCodeType) *InvoiceBuilder {
	b.documentCurrencyID = currencyID
	return b
}

func (b *InvoiceBuilder) WithDocumentToTaxCurrencyExchangeRate(rate Decimal) *InvoiceBuilder {
	b.taxCurrencyExchangeRate = rate
	return b
}

func (b *InvoiceBuilder) WithTaxCurrencyCode(currencyID CurrencyCodeType) *InvoiceBuilder {
	b.taxCurrencyID = currencyID
	return b
}

func (b *InvoiceBuilder) WithBillingReferences(billingReferences []InvoiceBillingReference) *InvoiceBuilder {
	b.billingReferences = billingReferences
	return b
}

func (b *InvoiceBuilder) AppendBillingReferences(billingReference InvoiceBillingReference) *InvoiceBuilder {
	return b.WithBillingReferences(append(b.billingReferences, billingReference))
}

func (b *InvoiceBuilder) WithSupplier(supplier InvoiceSupplier) *InvoiceBuilder {
	b.supplier = supplier
	return b
}

func (b *InvoiceBuilder) WithCustomer(customer InvoiceCustomer) *InvoiceBuilder {
	b.customer = customer
	return b
}

func (b *InvoiceBuilder) WithAllowancesCharges(allowancesCharges []InvoiceDocumentAllowanceCharge) *InvoiceBuilder {
	b.allowancesCharges = allowancesCharges
	return b
}

func (b *InvoiceBuilder) AppendAllowanceCharge(allowanceCharge InvoiceDocumentAllowanceCharge) *InvoiceBuilder {
	return b.WithAllowancesCharges(append(b.allowancesCharges, allowanceCharge))
}

func (b *InvoiceBuilder) WithInvoiceLines(invoiceLines []InvoiceLine) *InvoiceBuilder {
	b.invoiceLines = invoiceLines
	return b
}

func (b *InvoiceBuilder) AddTaxExemptionReason(taxCategoryCode TaxCategoryCodeType, reason string, exemptionCode TaxExemptionReasonCodeType) *InvoiceBuilder {
	if b.taxExeptionReasons == nil {
		b.taxExeptionReasons = make(map[TaxCategoryCodeType]taxExemptionReason)
	}
	b.taxExeptionReasons[taxCategoryCode] = taxExemptionReason{
		reason: reason,
		code:   exemptionCode,
	}
	return b
}

func (b *InvoiceBuilder) Build() (invoice Invoice, ok bool) {
	if b.id == "" || !b.issueDate.IsInitialized() ||
		b.documentCurrencyID == "" ||
		(b.taxCurrencyID != "" && b.taxCurrencyID != b.documentCurrencyID && !b.taxCurrencyExchangeRate.IsInitialized()) {
		return
	}

	taxCurrencyID := b.taxCurrencyID
	if taxCurrencyID == "" {
		taxCurrencyID = b.documentCurrencyID
	}

	invoice.ID = b.id
	invoice.IssueDate = b.issueDate
	invoice.DueDate = b.dueDate
	invoice.InvoiceTypeCode = b.invoiceType
	invoice.DocumentCurrencyCode = b.documentCurrencyID
	invoice.TaxCurrencyCode = b.taxCurrencyID

	// TODO:

	invoice.BillingReferences = b.billingReferences

	// TODO:

	invoice.Supplier = b.supplier
	invoice.Customer = b.customer

	// amountToTaxAmount converts an Amount assumed to be in the
	// DocumentCurrencyCode to an amount in TaxCurrencyCode
	amountToTaxAmount := func(a Decimal) Decimal {
		if taxCurrencyID == invoice.DocumentCurrencyCode {
			return a
		}
		return a.Mul(b.taxCurrencyExchangeRate)
	}

	invoice.AllowanceCharges = b.allowancesCharges
	invoice.InvoiceLines = b.invoiceLines

	var (
		lineExtensionAmount   = Zero
		allowanceTotalAmount  = Zero
		chargeTotalAmount     = Zero
		taxExclusiveAmount    = Zero
		taxInclusiveAmount    = Zero
		prepaidAmount         = Zero
		payableRoundingAmount = Zero
		payableAmount         = Zero
	)
	taxCategoryMap := make(taxCategoryMap)

	for _, line := range invoice.InvoiceLines {
		if line.LineExtensionAmount.CurrencyID != invoice.DocumentCurrencyCode {
			return
		}

		lineAmount := line.LineExtensionAmount.Amount
		lineExtensionAmount = lineExtensionAmount.Add(lineAmount)
		taxCategoryMap.addLineTaxCategory(line.Item.TaxCategory, amountToTaxAmount(lineAmount))
	}
	for _, allowanceCharge := range invoice.AllowanceCharges {
		var amount Decimal
		if allowanceCharge.ChargeIndicator {
			amount = allowanceCharge.Amount.Amount
			chargeTotalAmount = chargeTotalAmount.Add(allowanceCharge.Amount.Amount)
		} else {
			amount = allowanceCharge.Amount.Amount.Neg()
			allowanceTotalAmount = allowanceTotalAmount.Add(allowanceCharge.Amount.Amount)
		}
		taxCategoryMap.addDocumentTaxCategory(allowanceCharge.TaxCategory, amountToTaxAmount(amount))
	}

	taxTotal := Zero
	for _, taxCategorySummary := range taxCategoryMap.getSummaries() {
		taxAmount := taxCategorySummary.getTaxAmount()
		taxTotal = taxTotal.Add(taxAmount)

		subtotal := InvoiceTaxSubtotal{
			TaxableAmount: AmountWithCurrency{
				Amount:     taxCategorySummary.baseAmount,
				CurrencyID: taxCurrencyID,
			},
			TaxAmount: AmountWithCurrency{
				Amount:     taxAmount,
				CurrencyID: taxCurrencyID,
			},
			TaxCategory: taxCategorySummary.category,
		}
		if categoryCode := subtotal.TaxCategory.ID; categoryCode.TaxRateExempted() {
			if reason, rok := b.taxExeptionReasons[categoryCode]; !rok {
				return
			} else {
				subtotal.TaxCategory.TaxExemptionReason = reason.reason
				subtotal.TaxCategory.TaxExemptionReasonCode = reason.code
			}
		}
		invoice.TaxTotal.TaxSubtotals = append(invoice.TaxTotal.TaxSubtotals, subtotal)
	}

	taxExclusiveAmount = lineExtensionAmount.Add(chargeTotalAmount).Sub(allowanceTotalAmount)
	taxInclusiveAmount = taxExclusiveAmount.Add(taxTotal)
	payableAmount = taxInclusiveAmount.Sub(prepaidAmount)

	if taxTotal.IsPositive() {
		invoice.TaxTotal.TaxAmount = &AmountWithCurrency{
			Amount:     taxTotal,
			CurrencyID: taxCurrencyID,
		}
	}

	invoice.LegalMonetaryTotal.LineExtensionAmount = AmountWithCurrency{
		Amount:     lineExtensionAmount,
		CurrencyID: b.documentCurrencyID,
	}
	invoice.LegalMonetaryTotal.AllowanceTotalAmount = &AmountWithCurrency{
		Amount:     allowanceTotalAmount,
		CurrencyID: b.documentCurrencyID,
	}
	invoice.LegalMonetaryTotal.ChargeTotalAmount = &AmountWithCurrency{
		Amount:     chargeTotalAmount,
		CurrencyID: b.documentCurrencyID,
	}
	invoice.LegalMonetaryTotal.PrepaidAmount = &AmountWithCurrency{
		Amount:     prepaidAmount,
		CurrencyID: b.documentCurrencyID,
	}
	invoice.LegalMonetaryTotal.PayableRoundingAmount = &AmountWithCurrency{
		Amount:     payableRoundingAmount,
		CurrencyID: b.documentCurrencyID,
	}
	invoice.LegalMonetaryTotal.TaxExclusiveAmount = AmountWithCurrency{
		Amount:     taxExclusiveAmount,
		CurrencyID: b.documentCurrencyID,
	}
	invoice.LegalMonetaryTotal.TaxInclusiveAmount = AmountWithCurrency{
		Amount:     taxInclusiveAmount,
		CurrencyID: b.documentCurrencyID,
	}
	invoice.LegalMonetaryTotal.PayableAmount = AmountWithCurrency{
		Amount:     payableAmount,
		CurrencyID: b.documentCurrencyID,
	}

	ok = true
	return
}

type taxCategoryKey struct {
	id          TaxCategoryCodeType
	percent     float64
	taxSchemeID TaxSchemeIDType
}

func makeTaxCategoryKey(category InvoiceTaxCategory) taxCategoryKey {
	percent, _ := category.Percent.Value().Float64()
	return taxCategoryKey{
		id:          category.ID,
		percent:     percent,
		taxSchemeID: category.TaxSchemeID,
	}
}

func makeTaxCategoryKeyLine(category InvoiceLineTaxCategory) taxCategoryKey {
	percent, _ := category.Percent.Value().Float64()
	return taxCategoryKey{
		id:          category.ID,
		percent:     percent,
		taxSchemeID: category.TaxSchemeID,
	}
}

type taxCategorySummary struct {
	category   InvoiceTaxCategory
	baseAmount Decimal
}

func (s taxCategorySummary) getTaxAmount() Decimal {
	return s.baseAmount.Mul(s.category.Percent.Value()).Div(D(100)).Round(2)
}

// taxCategoryMap is not concurency-safe
type taxCategoryMap map[taxCategoryKey]*taxCategorySummary

func (m *taxCategoryMap) add(k taxCategoryKey, category InvoiceTaxCategory, amount Decimal) bool {
	if category.TaxSchemeID == TaxSchemeVAT {
		percent := category.Percent.Value()
		if !category.ID.TaxRateExempted() {
			if !percent.IsPositive() {
				return false
			}
		} else if !percent.IsZero() {
			return false
		}
	}

	val, ok := (*m)[k]
	if !ok {
		(*m)[k] = &taxCategorySummary{
			category:   category,
			baseAmount: amount,
		}
	} else {
		val.baseAmount = val.baseAmount.Add(amount)
	}

	return true
}

func (m *taxCategoryMap) addDocumentTaxCategory(category InvoiceTaxCategory, amount Decimal) bool {
	if m == nil {
		return false
	}

	k := makeTaxCategoryKey(category)
	return m.add(k, category, amount)
}

func (m *taxCategoryMap) addLineTaxCategory(category InvoiceLineTaxCategory, amount Decimal) bool {
	if m == nil {
		return false
	}

	k := makeTaxCategoryKeyLine(category)
	documentCategory := InvoiceTaxCategory{
		ID:          category.ID,
		Percent:     category.Percent,
		TaxSchemeID: category.TaxSchemeID,
	}
	return m.add(k, documentCategory, amount)
}

func (m taxCategoryMap) getSummaries() (summaries []taxCategorySummary) {
	for _, v := range m {
		summaries = append(summaries, taxCategorySummary{
			category:   v.category,
			baseAmount: v.baseAmount,
		})
	}
	return
}

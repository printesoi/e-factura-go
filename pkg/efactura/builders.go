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
	ierrors "github.com/printesoi/e-factura-go/internal/errors"
	"github.com/printesoi/e-factura-go/internal/ptr"
	"github.com/printesoi/e-factura-go/pkg/types"
)

// InvoiceLineAllowanceChargeBuilder builds an InvoiceLineAllowanceCharge object
type InvoiceLineAllowanceChargeBuilder struct {
	chargeIndicator           bool
	currencyID                CurrencyCodeType
	amount                    types.Decimal
	baseAmount                *types.Decimal
	allowanceChargeReasonCode *string
	allowanceChargeReason     *string
}

// NewInvoiceLineAllowanceChargeBuilder creates a new generic
// InvoiceLineAllowanceChargeBuilder.
func NewInvoiceLineAllowanceChargeBuilder(chargeIndicator bool, currencyID CurrencyCodeType, amount types.Decimal) *InvoiceLineAllowanceChargeBuilder {
	b := new(InvoiceLineAllowanceChargeBuilder)
	return b.WithChargeIndicator(chargeIndicator).
		WithCurrencyID(currencyID).WithAmount(amount)
}

// NewInvoiceLineAllowanceBuilder creates a new InvoiceLineAllowanceChargeBuilder
// builder that will build InvoiceLineAllowanceCharge object correspoding to an
// allowance (ChargeIndicator = false)
func NewInvoiceLineAllowanceBuilder(currencyID CurrencyCodeType, amount types.Decimal) *InvoiceLineAllowanceChargeBuilder {
	return NewInvoiceLineAllowanceChargeBuilder(false, currencyID, amount)
}

// NewInvoiceLineChargeBuilder creates a new InvoiceLineAllowanceChargeBuilder
// builder that will build InvoiceLineAllowanceCharge object correspoding to a
// charge (ChargeIndicator = true)
func NewInvoiceLineChargeBuilder(currencyID CurrencyCodeType, amount types.Decimal) *InvoiceLineAllowanceChargeBuilder {
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

func (b *InvoiceLineAllowanceChargeBuilder) WithAmount(amount types.Decimal) *InvoiceLineAllowanceChargeBuilder {
	b.amount = amount
	return b
}

func (b *InvoiceLineAllowanceChargeBuilder) WithBaseAmount(amount types.Decimal) *InvoiceLineAllowanceChargeBuilder {
	b.baseAmount = amount.Ptr()
	return b
}

func (b *InvoiceLineAllowanceChargeBuilder) WithAllowanceChargeReasonCode(allowanceChargeReasonCode string) *InvoiceLineAllowanceChargeBuilder {
	b.allowanceChargeReasonCode = ptr.String(allowanceChargeReasonCode)
	return b
}

func (b *InvoiceLineAllowanceChargeBuilder) WithAllowanceChargeReason(allowanceChargeReason string) *InvoiceLineAllowanceChargeBuilder {
	b.allowanceChargeReason = ptr.String(allowanceChargeReason)
	return b
}

func (b InvoiceLineAllowanceChargeBuilder) Build() (allowanceCharge InvoiceLineAllowanceCharge, err error) {
	if !b.amount.IsInitialized() {
		err = ierrors.NewBuilderErrorf(b, "", "amount not set")
		return
	}
	if b.currencyID == "" {
		err = ierrors.NewBuilderErrorf(b, "", "currency not set")
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
	invoicedQuantity types.Decimal
	baseQuantity     *types.Decimal

	grossPriceAmount types.Decimal
	priceDeduction   types.Decimal

	invoicePeriod     *InvoiceLinePeriod
	allowancesCharges []InvoiceLineAllowanceCharge

	itemName                       string
	itemDescription                string
	itemSellerID                   *string
	itemStandardItemIdentification *ItemStandardIdentificationCode
	itemCommodityClassification    *ItemCommodityClassification
	itemTaxCategory                InvoiceLineTaxCategory
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

func (b *InvoiceLineBuilder) WithInvoicedQuantity(quantity types.Decimal) *InvoiceLineBuilder {
	b.invoicedQuantity = quantity
	return b
}

func (b *InvoiceLineBuilder) WithBaseQuantity(quantity types.Decimal) *InvoiceLineBuilder {
	b.baseQuantity = &quantity
	return b
}

func (b *InvoiceLineBuilder) WithGrossPriceAmount(priceAmount types.Decimal) *InvoiceLineBuilder {
	b.grossPriceAmount = priceAmount
	return b
}

func (b *InvoiceLineBuilder) WithPriceDeduction(deduction types.Decimal) *InvoiceLineBuilder {
	b.priceDeduction = deduction
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

func (b *InvoiceLineBuilder) WithItemName(name string) *InvoiceLineBuilder {
	b.itemName = name
	return b
}

func (b *InvoiceLineBuilder) WithItemDescription(description string) *InvoiceLineBuilder {
	b.itemDescription = description
	return b
}

func (b *InvoiceLineBuilder) WithItemSellerID(id string) *InvoiceLineBuilder {
	b.itemSellerID = &id
	return b
}

func (b *InvoiceLineBuilder) WithItemStandardItemIdentification(identification ItemStandardIdentificationCode) *InvoiceLineBuilder {
	b.itemStandardItemIdentification = &identification
	return b
}

func (b *InvoiceLineBuilder) WithItemCommodityClassification(classification ItemCommodityClassification) *InvoiceLineBuilder {
	b.itemCommodityClassification = &classification
	return b
}

func (b *InvoiceLineBuilder) WithItemTaxCategory(taxCategory InvoiceLineTaxCategory) *InvoiceLineBuilder {
	b.itemTaxCategory = taxCategory
	return b
}

func (b InvoiceLineBuilder) Build() (line InvoiceLine, err error) {
	if b.id == "" {
		err = ierrors.NewBuilderErrorf(b, "", "id not set")
		return
	}
	if b.currencyID == "" {
		err = ierrors.NewBuilderErrorf(b, "", "id currency id not set")
		return
	}
	if !b.invoicedQuantity.IsInitialized() {
		err = ierrors.NewBuilderErrorf(b, "", "invoiced quantity not set")
		return
	}
	if b.unitCode == "" {
		err = ierrors.NewBuilderErrorf(b, "", "unit code not set")
		return
	}
	if !b.grossPriceAmount.IsInitialized() {
		err = ierrors.NewBuilderErrorf(b, "", "gross price amount not set")
		return
	}
	if b.itemName == "" {
		err = ierrors.NewBuilderErrorf(b, "", "item name not set")
		return
	}
	if b.itemTaxCategory.ID == "" || b.itemTaxCategory.TaxScheme.ID == "" {
		err = ierrors.NewBuilderErrorf(b, "", "item tax category not set")
		return
	}

	line.ID = b.id
	line.Note = b.note
	line.InvoicedQuantity = InvoicedQuantity{
		Quantity: b.invoicedQuantity,
		UnitCode: b.unitCode,
	}
	var netPriceAmount types.Decimal
	if b.priceDeduction.IsZero() {
		netPriceAmount = b.grossPriceAmount
	} else {
		netPriceAmount = b.grossPriceAmount.Sub(b.priceDeduction)
		line.Price.PriceAmount = AmountWithCurrency{
			Amount:     netPriceAmount,
			CurrencyID: b.currencyID,
		}
		line.Price.AllowanceCharge = &InvoiceLinePriceAllowanceCharge{
			ChargeIndicator: false,
			Amount: AmountWithCurrency{
				Amount:     b.priceDeduction,
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

	line.Item.Name = b.itemName
	line.Item.Description = b.itemDescription
	if b.itemSellerID != nil {
		line.Item.SellerItemID = NewIDNode(*b.itemSellerID)
	}
	line.Item.StandardItemIdentification = b.itemStandardItemIdentification
	line.Item.CommodityClassification = b.itemCommodityClassification
	line.Item.TaxCategory = b.itemTaxCategory

	line.AllowanceCharges = b.allowancesCharges
	line.InvoicePeriod = b.invoicePeriod

	// Invoiced quantity * (Item net price / item price base quantity)
	//  + Sum of invoice line charge amount
	//  - Sum of invoice line allowance amount
	baseQuantity := types.D(1)
	if b.baseQuantity != nil {
		baseQuantity = *b.baseQuantity
	}
	if baseQuantity.IsZero() {
		err = ierrors.NewBuilderErrorf(b, "", "base quantity cannot be zero")
		return
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
	return
}

// InvoiceDocumentAllowanceChargeBuilder builds an InvoiceDocumentAllowanceCharge object
type InvoiceDocumentAllowanceChargeBuilder struct {
	chargeIndicator           bool
	currencyID                CurrencyCodeType
	amount                    types.Decimal
	taxCategory               InvoiceTaxCategory
	baseAmount                *types.Decimal
	allowanceChargeReasonCode *string
	allowanceChargeReason     *string
}

// NewInvoiceDocumentAllowanceChargeBuilder creates a new generic
// InvoiceDocumentAllowanceChargeBuilder.
func NewInvoiceDocumentAllowanceChargeBuilder(chargeIndicator bool, currencyID CurrencyCodeType, amount types.Decimal, taxCategory InvoiceTaxCategory) *InvoiceDocumentAllowanceChargeBuilder {
	b := new(InvoiceDocumentAllowanceChargeBuilder)
	return b.WithChargeIndicator(chargeIndicator).WithCurrencyID(currencyID).
		WithAmount(amount).WithTaxCategory(taxCategory)
}

// NewInvoiceDocumentAllowanceBuilder creates a new InvoiceDocumentAllowanceChargeBuilder
// builder that will build InvoiceDocumentAllowanceCharge object correspoding to an
// allowance (ChargeIndicator = false)
func NewInvoiceDocumentAllowanceBuilder(currencyID CurrencyCodeType, amount types.Decimal, taxCategory InvoiceTaxCategory) *InvoiceDocumentAllowanceChargeBuilder {
	return NewInvoiceDocumentAllowanceChargeBuilder(false, currencyID, amount, taxCategory)
}

// NewInvoiceDocumentChargeBuilder creates a new InvoiceDocumentAllowanceChargeBuilder
// builder that will build InvoiceDocumentAllowanceCharge object correspoding to a
// charge (ChargeIndicator = true)
func NewInvoiceDocumentChargeBuilder(currencyID CurrencyCodeType, amount types.Decimal, taxCategory InvoiceTaxCategory) *InvoiceDocumentAllowanceChargeBuilder {
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

func (b *InvoiceDocumentAllowanceChargeBuilder) WithAmount(amount types.Decimal) *InvoiceDocumentAllowanceChargeBuilder {
	b.amount = amount
	return b
}

func (b *InvoiceDocumentAllowanceChargeBuilder) WithTaxCategory(taxCategory InvoiceTaxCategory) *InvoiceDocumentAllowanceChargeBuilder {
	b.taxCategory = taxCategory
	return b
}

func (b *InvoiceDocumentAllowanceChargeBuilder) WithBaseAmount(amount types.Decimal) *InvoiceDocumentAllowanceChargeBuilder {
	b.baseAmount = amount.Ptr()
	return b
}

func (b *InvoiceDocumentAllowanceChargeBuilder) WithAllowanceChargeReasonCode(allowanceChargeReasonCode string) *InvoiceDocumentAllowanceChargeBuilder {
	b.allowanceChargeReasonCode = ptr.String(allowanceChargeReasonCode)
	return b
}

func (b *InvoiceDocumentAllowanceChargeBuilder) WithAllowanceChargeReason(allowanceChargeReason string) *InvoiceDocumentAllowanceChargeBuilder {
	b.allowanceChargeReason = ptr.String(allowanceChargeReason)
	return b
}

func (b InvoiceDocumentAllowanceChargeBuilder) Build() (allowanceCharge InvoiceDocumentAllowanceCharge, err error) {
	if !b.amount.IsInitialized() {
		err = ierrors.NewBuilderErrorf(b, "", "amount not set")
		return
	}
	if b.currencyID == "" {
		err = ierrors.NewBuilderErrorf(b, "", "current id not set")
		return
	}
	if b.taxCategory.ID == "" || b.taxCategory.TaxScheme.ID == "" {
		err = ierrors.NewBuilderErrorf(b, "", "item tax category not set")
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
	return
}

type taxExemptionReason struct {
	reason string
	code   TaxExemptionReasonCodeType
}

// InvoiceBuilder builds an Invoice object
type InvoiceBuilder struct {
	id          string
	issueDate   types.Date
	dueDate     *types.Date
	invoiceType InvoiceTypeCodeType

	documentCurrencyID      CurrencyCodeType
	taxCurrencyID           CurrencyCodeType
	taxCurrencyExchangeRate types.Decimal

	taxExeptionReasons map[TaxCategoryCodeType]taxExemptionReason

	accountingCost            string
	buyerReference            string
	orderReference            *InvoiceOrderReference
	notes                     []InvoiceNote
	invoicePeriod             *InvoicePeriod
	billingReferences         []InvoiceDocumentReference
	contractDocumentReference *string
	supplier                  InvoiceSupplierParty
	customer                  InvoiceCustomerParty
	paymentMeans              *InvoicePaymentMeans
	paymentTerms              *InvoicePaymentTerms

	allowancesCharges []InvoiceDocumentAllowanceCharge
	invoiceLines      []InvoiceLine

	expectedTaxInclusiveAmount *types.Decimal
}

func NewInvoiceBuilder(id string) (b *InvoiceBuilder) {
	b = new(InvoiceBuilder)
	return b.WithID(id)
}

func (b *InvoiceBuilder) WithID(id string) *InvoiceBuilder {
	b.id = id
	return b
}

func (b *InvoiceBuilder) WithIssueDate(date types.Date) *InvoiceBuilder {
	b.issueDate = date
	return b
}

func (b *InvoiceBuilder) WithDueDate(date types.Date) *InvoiceBuilder {
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

func (b *InvoiceBuilder) WithDocumentToTaxCurrencyExchangeRate(rate types.Decimal) *InvoiceBuilder {
	b.taxCurrencyExchangeRate = rate
	return b
}

func (b *InvoiceBuilder) WithTaxCurrencyCode(currencyID CurrencyCodeType) *InvoiceBuilder {
	b.taxCurrencyID = currencyID
	return b
}

func (b *InvoiceBuilder) WithBillingReferences(billingReferences []InvoiceDocumentReference) *InvoiceBuilder {
	b.billingReferences = billingReferences
	return b
}

func (b *InvoiceBuilder) AppendBillingReferences(billingReferences ...InvoiceDocumentReference) *InvoiceBuilder {
	return b.WithBillingReferences(append(b.billingReferences, billingReferences...))
}

func (b *InvoiceBuilder) WithSupplier(supplier InvoiceSupplierParty) *InvoiceBuilder {
	b.supplier = supplier
	return b
}

func (b *InvoiceBuilder) WithCustomer(customer InvoiceCustomerParty) *InvoiceBuilder {
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

func (b *InvoiceBuilder) AppendInvoiceLines(lines ...InvoiceLine) *InvoiceBuilder {
	b.invoiceLines = append(b.invoiceLines, lines...)
	return b
}

func (b *InvoiceBuilder) WithAccountingCost(accountingCost string) *InvoiceBuilder {
	b.accountingCost = accountingCost
	return b
}

func (b *InvoiceBuilder) WithBuyerReference(buyerReference string) *InvoiceBuilder {
	b.buyerReference = buyerReference
	return b
}

func (b *InvoiceBuilder) WithOrderReference(orderReference InvoiceOrderReference) *InvoiceBuilder {
	b.orderReference = &orderReference
	return b
}

func (b *InvoiceBuilder) WithNotes(notes []InvoiceNote) *InvoiceBuilder {
	b.notes = notes
	return b
}

func (b *InvoiceBuilder) AppendNotes(notes ...InvoiceNote) *InvoiceBuilder {
	b.notes = append(b.notes, notes...)
	return b
}

func (b *InvoiceBuilder) WithInvoicePeriod(invoicePeriod InvoicePeriod) *InvoiceBuilder {
	b.invoicePeriod = &invoicePeriod
	return b
}

func (b *InvoiceBuilder) WithContractDocumentReference(contractDocumentReference string) *InvoiceBuilder {
	b.contractDocumentReference = &contractDocumentReference
	return b
}

func (b *InvoiceBuilder) WithPaymentMeans(paymentMeans InvoicePaymentMeans) *InvoiceBuilder {
	b.paymentMeans = &paymentMeans
	return b
}

func (b *InvoiceBuilder) WithPaymentTerms(paymentTerms InvoicePaymentTerms) *InvoiceBuilder {
	b.paymentTerms = &paymentTerms
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

// WithExpectedTaxInclusiveAmount sets the expected tax inclusive amount. This
// is useful in cases where the invoice was already generated and the rounding
// algorithm might differ from the way the rounding is done for e-factura. If
// the tax inclusive amount generated is different than the given amount, the
// BT-114 term will be set (Payable rounding amount) and Payable Amount
// (BT-115) is adjusted with the difference.
func (b *InvoiceBuilder) WithExpectedTaxInclusiveAmount(amount types.Decimal) *InvoiceBuilder {
	b.expectedTaxInclusiveAmount = amount.Ptr()
	return b
}

func (b InvoiceBuilder) Build() (retInvoice Invoice, err error) {
	if b.id == "" {
		err = ierrors.NewBuilderErrorf(b, "", "id not set")
		return
	}
	if !b.issueDate.IsInitialized() {
		err = ierrors.NewBuilderErrorf(b, "", "issue date not set")
		return
	}
	if b.documentCurrencyID == "" {
		err = ierrors.NewBuilderErrorf(b, "", "document currency id not set")
		return
	}
	if b.taxCurrencyID != "" && b.taxCurrencyID != b.documentCurrencyID && !b.taxCurrencyExchangeRate.IsInitialized() {
		err = ierrors.NewBuilderErrorf(b, "", "document to tax currency exchange rate not set")
		return
	}

	taxCurrencyID := b.taxCurrencyID
	if taxCurrencyID == "" {
		taxCurrencyID = b.documentCurrencyID
	}

	var invoice Invoice
	invoice.Prefill()

	invoice.ID = b.id
	invoice.IssueDate = b.issueDate
	invoice.DueDate = b.dueDate
	invoice.InvoiceTypeCode = b.invoiceType
	invoice.DocumentCurrencyCode = b.documentCurrencyID
	invoice.TaxCurrencyCode = b.taxCurrencyID
	invoice.AccountingCost = b.accountingCost
	invoice.BuyerReference = b.buyerReference
	invoice.OrderReference = b.orderReference
	invoice.Note = b.notes
	invoice.InvoicePeriod = b.invoicePeriod

	for _, ref := range b.billingReferences {
		invoice.BillingReferences = append(invoice.BillingReferences, InvoiceBillingReference{
			InvoiceDocumentReference: ref,
		})
	}

	if b.contractDocumentReference != nil {
		invoice.ContractDocumentReference = NewIDNode(*b.contractDocumentReference)
	}

	invoice.Supplier.Party = b.supplier
	invoice.Customer.Party = b.customer

	invoice.PaymentMeans = b.paymentMeans
	invoice.PaymentTerms = b.paymentTerms

	// amountToTaxAmount converts an Amount assumed to be in the
	// DocumentCurrencyCode to an amount in TaxCurrencyCode
	amountToTaxAmount := func(a types.Decimal) types.Decimal {
		if taxCurrencyID == invoice.DocumentCurrencyCode {
			return a
		}
		return a.Mul(b.taxCurrencyExchangeRate).AsAmount()
	}

	invoice.AllowanceCharges = b.allowancesCharges
	invoice.InvoiceLines = b.invoiceLines

	var (
		lineExtensionAmount   = types.Zero
		allowanceTotalAmount  = types.Zero
		chargeTotalAmount     = types.Zero
		taxExclusiveAmount    = types.Zero
		taxInclusiveAmount    = types.Zero
		prepaidAmount         = types.Zero
		payableRoundingAmount = types.Zero
		payableAmount         = types.Zero
	)

	taxCategoryMap := make(taxCategoryMap)
	for i, line := range invoice.InvoiceLines {
		if line.LineExtensionAmount.CurrencyID != invoice.DocumentCurrencyCode {
			err = ierrors.NewBuilderErrorf(b, "", "invoice line %d: invalid currency id", i)
			return
		}

		lineAmount := line.LineExtensionAmount.Amount
		lineExtensionAmount = lineExtensionAmount.Add(lineAmount)
		if !taxCategoryMap.addLineTaxCategory(line.Item.TaxCategory, lineAmount) {
			err = ierrors.NewBuilderErrorf(b, "", "invoice line %d: invalid tax category", i)
			return
		}
	}
	for i, allowanceCharge := range invoice.AllowanceCharges {
		var amount types.Decimal
		if allowanceCharge.ChargeIndicator {
			amount = allowanceCharge.Amount.Amount
			chargeTotalAmount = chargeTotalAmount.Add(allowanceCharge.Amount.Amount)
		} else {
			amount = allowanceCharge.Amount.Amount.Neg()
			allowanceTotalAmount = allowanceTotalAmount.Add(allowanceCharge.Amount.Amount)
		}
		if !taxCategoryMap.addDocumentTaxCategory(allowanceCharge.TaxCategory, amount) {
			err = ierrors.NewBuilderErrorf(b, "", "invoice allowance/charge %d: invalid tax category", i)
			return
		}
	}

	taxTotal, taxTotalTaxCurrency := types.Zero, types.Zero
	var taxSubtotals []InvoiceTaxSubtotal

	for _, taxCategorySummary := range taxCategoryMap.getSummaries() {
		taxAmount := taxCategorySummary.getTaxAmount()
		taxAmountTaxCurrency := amountToTaxAmount(taxAmount)

		taxTotal = taxTotal.Add(taxAmount)
		taxTotalTaxCurrency = taxTotalTaxCurrency.Add(taxAmountTaxCurrency)

		subtotal := InvoiceTaxSubtotal{
			TaxableAmount: AmountWithCurrency{
				Amount:     taxCategorySummary.baseAmount,
				CurrencyID: invoice.DocumentCurrencyCode,
			},
			TaxAmount: AmountWithCurrency{
				Amount:     taxAmount,
				CurrencyID: invoice.DocumentCurrencyCode,
			},
			TaxCategory: taxCategorySummary.category,
		}

		if subtotal.TaxCategory.ID.TaxRateExempted() && subtotal.TaxCategory.ID.ExemptionReasonRequired() {
			if reason, rok := b.taxExeptionReasons[subtotal.TaxCategory.ID]; !rok {
				err = ierrors.NewBuilderErrorf(b, "", "tax category %s/%s: no exemption reason",
					subtotal.TaxCategory.ID, subtotal.TaxCategory.Percent.String())
				return
			} else {
				subtotal.TaxCategory.TaxExemptionReason = reason.reason
				subtotal.TaxCategory.TaxExemptionReasonCode = reason.code
			}
		}
		taxSubtotals = append(taxSubtotals, subtotal)
	}

	taxExclusiveAmount = lineExtensionAmount.Add(chargeTotalAmount).Sub(allowanceTotalAmount)
	taxInclusiveAmount = taxExclusiveAmount.Add(taxTotal)
	if b.expectedTaxInclusiveAmount != nil && !b.expectedTaxInclusiveAmount.Equal(taxInclusiveAmount) {
		payableRoundingAmount = b.expectedTaxInclusiveAmount.Sub(taxInclusiveAmount)
	}
	payableAmount = taxInclusiveAmount.Sub(prepaidAmount).Add(payableRoundingAmount)

	if len(taxSubtotals) > 0 {
		taxTotalNode := InvoiceTaxTotal{
			TaxAmount: &AmountWithCurrency{
				Amount:     taxTotal,
				CurrencyID: invoice.DocumentCurrencyCode,
			},
			TaxSubtotals: taxSubtotals,
		}
		invoice.TaxTotal = append(invoice.TaxTotal, taxTotalNode)
	}
	if taxCurrencyID != invoice.DocumentCurrencyCode {
		taxTotalNode := InvoiceTaxTotal{
			TaxAmount: &AmountWithCurrency{
				Amount:     taxTotalTaxCurrency,
				CurrencyID: taxCurrencyID,
			},
		}
		invoice.TaxTotal = append(invoice.TaxTotal, taxTotalNode)
	}

	invoice.LegalMonetaryTotal.LineExtensionAmount = AmountWithCurrency{
		Amount:     lineExtensionAmount,
		CurrencyID: b.documentCurrencyID,
	}
	if !allowanceTotalAmount.IsZero() {
		invoice.LegalMonetaryTotal.AllowanceTotalAmount = &AmountWithCurrency{
			Amount:     allowanceTotalAmount,
			CurrencyID: b.documentCurrencyID,
		}
	}
	if !chargeTotalAmount.IsZero() {
		invoice.LegalMonetaryTotal.ChargeTotalAmount = &AmountWithCurrency{
			Amount:     chargeTotalAmount,
			CurrencyID: b.documentCurrencyID,
		}
	}
	invoice.LegalMonetaryTotal.TaxExclusiveAmount = AmountWithCurrency{
		Amount:     taxExclusiveAmount,
		CurrencyID: b.documentCurrencyID,
	}
	invoice.LegalMonetaryTotal.TaxInclusiveAmount = AmountWithCurrency{
		Amount:     taxInclusiveAmount,
		CurrencyID: b.documentCurrencyID,
	}
	if !prepaidAmount.IsZero() {
		invoice.LegalMonetaryTotal.PrepaidAmount = &AmountWithCurrency{
			Amount:     prepaidAmount,
			CurrencyID: b.documentCurrencyID,
		}
	}
	if !payableRoundingAmount.IsZero() {
		invoice.LegalMonetaryTotal.PayableRoundingAmount = &AmountWithCurrency{
			Amount:     payableRoundingAmount,
			CurrencyID: b.documentCurrencyID,
		}
	}
	invoice.LegalMonetaryTotal.PayableAmount = AmountWithCurrency{
		Amount:     payableAmount,
		CurrencyID: b.documentCurrencyID,
	}

	retInvoice = invoice
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
		taxSchemeID: category.TaxScheme.ID,
	}
}

func makeTaxCategoryKeyLine(category InvoiceLineTaxCategory) taxCategoryKey {
	percent, _ := category.Percent.Value().Float64()
	return taxCategoryKey{
		id:          category.ID,
		percent:     percent,
		taxSchemeID: category.TaxScheme.ID,
	}
}

type taxCategorySummary struct {
	category   InvoiceTaxCategory
	baseAmount types.Decimal
}

func (s taxCategorySummary) getTaxAmount() types.Decimal {
	return s.baseAmount.Mul(s.category.Percent.Value()).Div(types.D(100)).Round(2)
}

// taxCategoryMap is not concurency-safe
type taxCategoryMap map[taxCategoryKey]*taxCategorySummary

func (m *taxCategoryMap) add(k taxCategoryKey, category InvoiceTaxCategory, amount types.Decimal) bool {
	if category.TaxScheme.ID == TaxSchemeIDVAT {
		percent := category.Percent.Value()
		if !category.ID.TaxRateExempted() {
			if percent.IsZero() {
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

func (m *taxCategoryMap) addDocumentTaxCategory(category InvoiceTaxCategory, amount types.Decimal) bool {
	if m == nil {
		return false
	}

	k := makeTaxCategoryKey(category)
	return m.add(k, category, amount)
}

func (m *taxCategoryMap) addLineTaxCategory(category InvoiceLineTaxCategory, amount types.Decimal) bool {
	if m == nil {
		return false
	}

	k := makeTaxCategoryKeyLine(category)
	documentCategory := InvoiceTaxCategory{
		ID:        category.ID,
		Percent:   category.Percent,
		TaxScheme: category.TaxScheme,
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

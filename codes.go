package efactura

type InvoiceTypeCodeType string

const (
	InvoiceTypeCodeFactura                    InvoiceTypeCodeType = "380"
	InvoiceTypeCodeFacturaCorectata           InvoiceTypeCodeType = "381"
	InvoiceTypeCodeAutoFactura                InvoiceTypeCodeType = "389"
	InvoiceTypeCodeFacturaInformatiiContabile InvoiceTypeCodeType = "751"
)

type CurrencyCodeType string

const (
	CurrencyRON CurrencyCodeType = "RON"
)

type TaxExemptionReasonCodeType string

// TODO: values

type TaxCategoryCodeType string

// TODO: values

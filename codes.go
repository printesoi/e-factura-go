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

type TaxCategoryCodeType string

const (
	TaxCategoryTVACotaNormalaRedusa      TaxCategoryCodeType = "S"
	TaxCategoryTVACotaZero               TaxCategoryCodeType = "S"
	TaxCategoryScutireTVA                TaxCategoryCodeType = "E"
	TaxCategoryTVATaxareInversa          TaxCategoryCodeType = "AE"
	TaxCategoryTVALivrariIntracomunitare TaxCategoryCodeType = "K"
	TaxCategoryTVAExporturi              TaxCategoryCodeType = "G"
	TaxCategoryNuFaceObiectulTVA         TaxCategoryCodeType = "O"
	TaxCategoryTaxeInsuleCanare          TaxCategoryCodeType = "L"
	TaxCategoryTaxeCeutaMelilla          TaxCategoryCodeType = "M"
)

type TaxExemptionReasonCodeType string

const (
	// VATEX-EU-79-C - Exceptie cf. Art. 79, lit c din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_79_C = "VATEX-EU-79-C"
	// VATEX-EU-132 - Exceptie cf. Art. 132 din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132 = "VATEX-EU-132"
	// VATEX-EU-132-1A - Exceptie cf. Art. 132 , alin. 1, lit (a) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1A = "VATEX-EU-132-1A"
	// VATEX-EU-132-1B - Exceptie cf. Art. 132 , alin. 1, lit (b) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1B = "VATEX-EU-132-1B"
	// VATEX-EU-132-1C - Exceptie cf. Art. 132 , alin. 1, lit (c) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1C = "VATEX-EU-132-1C"
	// VATEX-EU-132-1D - Exceptie cf. Art. 132 , alin. 1, lit (d) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1D = "VATEX-EU-132-1D"
	// VATEX-EU-132-1E - Exceptie cf. Art. 132 , alin. 1, lit (e) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1E = "VATEX-EU-132-1E"
	// VATEX-EU-132-1F - Exceptie cf. Art. 132 , alin. 1, lit (f) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1F = "VATEX-EU-132-1F"
	// VATEX-EU-132-1G - Exceptie cf. Art. 132 , alin. 1, lit (g) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1G = "VATEX-EU-132-1G"
	// VATEX-EU-132-1H - Exceptie cf. Art. 132 , alin. 1, lit (h) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1H = "VATEX-EU-132-1H"
	// VATEX-EU-132-1I - Exceptie cf. Art. 132 , alin. 1, lit (i) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1I = "VATEX-EU-132-1I"
	// VATEX-EU-132-1J - Exceptie cf. Art. 132 , alin. 1, lit (j) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1J = "VATEX-EU-132-1J"
	// VATEX-EU-132-1K - Exceptie cf. Art. 132 , alin. 1, lit (k) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1K = "VATEX-EU-132-1K"
	// VATEX-EU-132-1L - Exceptie cf. Art. 132 , alin. 1, lit (l) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1L = "VATEX-EU-132-1L"
	// VATEX-EU-132-1M - Exceptie cf. Art. 132 , alin. 1, lit (m) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1M = "VATEX-EU-132-1M"
	// VATEX-EU-132-1N - Exceptie cf. Art. 132 , alin. 1, lit (n) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1N = "VATEX-EU-132-1N"
	// VATEX-EU-132-1O - Exceptie cf. Art. 132 , alin. 1, lit (o) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1O = "VATEX-EU-132-1O"
	// VATEX-EU-132-1P - Exceptie cf. Art. 132 , alin. 1, lit (p) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1P = "VATEX-EU-132-1P"
	// VATEX-EU-132-1Q - Exceptie cf. Art. 132 , alin. 1, lit (q) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_132_1Q = "VATEX-EU-132-1Q"
	// VATEX-EU-143 - Exceptie cf. Art. 143 din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143 = "VATEX-EU-143"
	// VATEX-EU-143-1A - Exceptie cf. Art. 143, alin. 1, lit (a) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143_1A = "VATEX-EU-143-1A"
	// VATEX-EU-143-1B - Exceptie cf. Art. 143, alin. 1, lit (b) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143_1B = "VATEX-EU-143-1B"
	// VATEX-EU-143-1C - Exceptie cf. Art. 143, alin. 1, lit (c) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143_1C = "VATEX-EU-143-1C"
	// VATEX-EU-143-1D - Exceptie cf. Art. 143, alin. 1, lit (d) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143_1D = "VATEX-EU-143-1D"
	// VATEX-EU-143-1E - Exceptie cf. Art. 143, alin. 1, lit (e) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143_1E = "VATEX-EU-143-1E"
	// VATEX-EU-143-1F - Exceptie cf. Art. 143, alin. 1, lit (f) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143_1F = "VATEX-EU-143-1F"
	// VATEX-EU-143-1FA - Exceptie cf. Art. 143, alin. 1, lit (fa) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143_1FA = "VATEX-EU-143-1FA"
	// VATEX-EU-143-1G - Exceptie cf. Art. 143, alin. 1, lit (g) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143_1G = "VATEX-EU-143-1G"
	// VATEX-EU-143-1H - Exceptie cf. Art. 143, alin. 1, lit (h) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143_1H = "VATEX-EU-143-1H"
	// VATEX-EU-143-1I - Exceptie cf. Art. 143, alin. 1, lit (i) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143_1I = "VATEX-EU-143-1I"
	// VATEX-EU-143-1J - Exceptie cf. Art. 143, alin. 1, lit (j) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143_1J = "VATEX-EU-143-1J"
	// VATEX-EU-143-1K - Exceptie cf. Art. 143, alin. 1, lit (k) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143_1K = "VATEX-EU-143-1K"
	// VATEX-EU-143-1L - Exceptie cf. Art. 143, alin. 1, lit (l) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_143_1L = "VATEX-EU-143-1L"
	// VATEX-EU-148 - Exceptie cf. Art. 148 din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_148 = "VATEX-EU-148"
	// VATEX-EU-148-A - Exceptie cf. Art. 148, lit. (a) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_148_A = "VATEX-EU-148-A"
	// VATEX-EU-148-B - Exceptie cf. Art. 148, lit. (b) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_148_B = "VATEX-EU-148-B"
	// VATEX-EU-148-C - Exceptie cf. Art. 148, lit. (c) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_148_C = "VATEX-EU-148-C"
	// VATEX-EU-148-D - Exceptie cf. Art. 148, lit. (d) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_148_D = "VATEX-EU-148-D"
	// VATEX-EU-148-E - Exceptie cf. Art. 148, lit. (e) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_148_E = "VATEX-EU-148-E"
	// VATEX-EU-148-F - Exceptie cf. Art. 148, lit. (f) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_148_F = "VATEX-EU-148-F"
	// VATEX-EU-148-G - Exceptie cf. Art. 148, lit. (g) din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_148_G = "VATEX-EU-148-G"
	// VATEX-EU-151 - Exceptie cf. Art. 151 din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_151 = "VATEX-EU-151"
	// VATEX-EU-151-1A - Exceptie cf. Art. 151, alin. 1, lit (a). din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_151_1A = "VATEX-EU-151-1A"
	// VATEX-EU-151-1AA - Exceptie cf. Art. 151, alin. 1, lit (aa). din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_151_1AA = "VATEX-EU-151-1AA"
	// VATEX-EU-151-1B - Exceptie cf. Art. 151, alin. 1, lit (b). din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_151_1B = "VATEX-EU-151-1B"
	// VATEX-EU-151-1C - Exceptie cf. Art. 151, alin. 1, lit (c). din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_151_1C = "VATEX-EU-151-1C"
	// VATEX-EU-151-1D - Exceptie cf. Art. 151, alin. 1, lit (d). din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_151_1D = "VATEX-EU-151-1D"
	// VATEX-EU-151-1E - Exceptie cf. Art. 151, alin. 1, lit (e). din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_151_1E = "VATEX-EU-151-1E"
	// VATEX-EU-309 - Exceptie cf. Art. 309 din Directiva 2006/112/EC
	TaxExemptionCodeVATEX_EU_309 = "VATEX-EU-309"
	// VATEX-EU-AE - Taxare inversa
	TaxExemptionCodeVATEX_EU_AE = "VATEX-EU-AE"
	// VATEX-EU-D - Intra-Regim special pentru agentiile de turism
	TaxExemptionCodeVATEX_EU_D = "VATEX-EU-D"
	// VATEX-EU-F - Regim special pentru bunuri second hand
	TaxExemptionCodeVATEX_EU_F = "VATEX-EU-F"
	// VATEX-EU-G - Export in afara UE
	TaxExemptionCodeVATEX_EU_G = "VATEX-EU-G"
	// VATEX-EU-I - Regim special pentru obiecte de arta
	TaxExemptionCodeVATEX_EU_I = "VATEX-EU-I"
	// VATEX-EU-IC - Livrare intra-comunitara
	TaxExemptionCodeVATEX_EU_IC = "VATEX-EU-IC"
	// VATEX-EU-J - Regim special pentru obiecte de colectie si antichitati
	TaxExemptionCodeVATEX_EU_J = "VATEX-EU-J"
	// VATEX-EU-O - Nu face obiectul TVA
	TaxExemptionCodeVATEX_EU_O = "VATEX-EU-O"
)

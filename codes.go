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

type CountryCodeType string

const (
	// RO - Romania
	CountryCodeRO CountryCodeType = "RO"
	// AD - Andorra
	CountryCodeAD CountryCodeType = "AD"
	// AE - United Arab Emirates
	CountryCodeAE CountryCodeType = "AE"
	// AF - Afghanistan
	CountryCodeAF CountryCodeType = "AF"
	// AG - Antigua and Barbuda
	CountryCodeAG CountryCodeType = "AG"
	// AI - Anguilla
	CountryCodeAI CountryCodeType = "AI"
	// AL - Albania
	CountryCodeAL CountryCodeType = "AL"
	// AM - Armenia
	CountryCodeAM CountryCodeType = "AM"
	// AO - Angola
	CountryCodeAO CountryCodeType = "AO"
	// AQ - Antarctica
	CountryCodeAQ CountryCodeType = "AQ"
	// AR - Argentina
	CountryCodeAR CountryCodeType = "AR"
	// AS - American Samoa
	CountryCodeAS CountryCodeType = "AS"
	// AT - Austria
	CountryCodeAT CountryCodeType = "AT"
	// AU - Australia
	CountryCodeAU CountryCodeType = "AU"
	// AW - Aruba
	CountryCodeAW CountryCodeType = "AW"
	// AX - Aland Islands
	CountryCodeAX CountryCodeType = "AX"
	// AZ - Azerbaijan
	CountryCodeAZ CountryCodeType = "AZ"
	// BA - Bosnia and Herzegovina
	CountryCodeBA CountryCodeType = "BA"
	// BB - Barbados
	CountryCodeBB CountryCodeType = "BB"
	// BD - Bangladesh
	CountryCodeBD CountryCodeType = "BD"
	// BE - Belgium
	CountryCodeBE CountryCodeType = "BE"
	// BF - Burkina Faso
	CountryCodeBF CountryCodeType = "BF"
	// BG - Bulgaria
	CountryCodeBG CountryCodeType = "BG"
	// BH - Bahrain
	CountryCodeBH CountryCodeType = "BH"
	// BI - Burundi
	CountryCodeBI CountryCodeType = "BI"
	// BJ - Benin
	CountryCodeBJ CountryCodeType = "BJ"
	// BL - Saint Barth&#xE9;lemy
	CountryCodeBL CountryCodeType = "BL"
	// BM - Bermuda
	CountryCodeBM CountryCodeType = "BM"
	// BN - Brunei Darussalam
	CountryCodeBN CountryCodeType = "BN"
	// BO - Bolivia, Plurinational State of
	CountryCodeBO CountryCodeType = "BO"
	// BQ - Bonaire, Sint Eustatius and Saba
	CountryCodeBQ CountryCodeType = "BQ"
	// BR - Brazil
	CountryCodeBR CountryCodeType = "BR"
	// BS - Bahamas
	CountryCodeBS CountryCodeType = "BS"
	// BT - Bhutan
	CountryCodeBT CountryCodeType = "BT"
	// BV - Bouvet Island
	CountryCodeBV CountryCodeType = "BV"
	// BW - Botswana
	CountryCodeBW CountryCodeType = "BW"
	// BY - Belarus
	CountryCodeBY CountryCodeType = "BY"
	// BZ - Belize
	CountryCodeBZ CountryCodeType = "BZ"
	// CA - Canada
	CountryCodeCA CountryCodeType = "CA"
	// CC - Cocos (Keeling) Islands
	CountryCodeCC CountryCodeType = "CC"
	// CD - Congo, the Democratic Republic of the
	CountryCodeCD CountryCodeType = "CD"
	// CF - Central African Republic
	CountryCodeCF CountryCodeType = "CF"
	// CG - Congo
	CountryCodeCG CountryCodeType = "CG"
	// CH - Switzerland
	CountryCodeCH CountryCodeType = "CH"
	// CI - C&#xF4;te d'Ivoire
	CountryCodeCI CountryCodeType = "CI"
	// CK - Cook Islands
	CountryCodeCK CountryCodeType = "CK"
	// CL - Chile
	CountryCodeCL CountryCodeType = "CL"
	// CM - Cameroon
	CountryCodeCM CountryCodeType = "CM"
	// CN - China
	CountryCodeCN CountryCodeType = "CN"
	// CO - Colombia
	CountryCodeCO CountryCodeType = "CO"
	// CR - Costa Rica
	CountryCodeCR CountryCodeType = "CR"
	// CU - Cuba
	CountryCodeCU CountryCodeType = "CU"
	// CV - Cabo Verde
	CountryCodeCV CountryCodeType = "CV"
	// CW - Cura&#xE7;ao
	CountryCodeCW CountryCodeType = "CW"
	// CX - Christmas Island
	CountryCodeCX CountryCodeType = "CX"
	// CY - Cyprus
	CountryCodeCY CountryCodeType = "CY"
	// CZ - Czechia
	CountryCodeCZ CountryCodeType = "CZ"
	// DE - Germany
	CountryCodeDE CountryCodeType = "DE"
	// DJ - Djibouti
	CountryCodeDJ CountryCodeType = "DJ"
	// DK - Denmark
	CountryCodeDK CountryCodeType = "DK"
	// DM - Dominica
	CountryCodeDM CountryCodeType = "DM"
	// DO - Dominican Republic
	CountryCodeDO CountryCodeType = "DO"
	// DZ - Algeria
	CountryCodeDZ CountryCodeType = "DZ"
	// EC - Ecuador
	CountryCodeEC CountryCodeType = "EC"
	// EE - Estonia
	CountryCodeEE CountryCodeType = "EE"
	// EG - Egypt
	CountryCodeEG CountryCodeType = "EG"
	// EH - Western Sahara
	CountryCodeEH CountryCodeType = "EH"
	// ER - Eritrea
	CountryCodeER CountryCodeType = "ER"
	// ES - Spain
	CountryCodeES CountryCodeType = "ES"
	// ET - Ethiopia
	CountryCodeET CountryCodeType = "ET"
	// FI - Finland
	CountryCodeFI CountryCodeType = "FI"
	// FJ - Fiji
	CountryCodeFJ CountryCodeType = "FJ"
	// FK - Falkland Islands (Malvinas)
	CountryCodeFK CountryCodeType = "FK"
	// FM - Micronesia, Federated States of
	CountryCodeFM CountryCodeType = "FM"
	// FO - Faroe Islands
	CountryCodeFO CountryCodeType = "FO"
	// FR - France
	CountryCodeFR CountryCodeType = "FR"
	// GA - Gabon
	CountryCodeGA CountryCodeType = "GA"
	// GB - United Kingdom of Great Britain and Northern Ireland
	CountryCodeGB CountryCodeType = "GB"
	// GD - Grenada
	CountryCodeGD CountryCodeType = "GD"
	// GE - Georgia
	CountryCodeGE CountryCodeType = "GE"
	// GF - French Guiana
	CountryCodeGF CountryCodeType = "GF"
	// GG - Guernsey
	CountryCodeGG CountryCodeType = "GG"
	// GH - Ghana
	CountryCodeGH CountryCodeType = "GH"
	// GI - Gibraltar
	CountryCodeGI CountryCodeType = "GI"
	// GL - Greenland
	CountryCodeGL CountryCodeType = "GL"
	// GM - Gambia
	CountryCodeGM CountryCodeType = "GM"
	// GN - Guinea
	CountryCodeGN CountryCodeType = "GN"
	// GP - Guadeloupe
	CountryCodeGP CountryCodeType = "GP"
	// GQ - Equatorial Guinea
	CountryCodeGQ CountryCodeType = "GQ"
	// GR - Greece
	CountryCodeGR CountryCodeType = "GR"
	// GS - South Georgia and the South Sandwich Islands
	CountryCodeGS CountryCodeType = "GS"
	// GT - Guatemala
	CountryCodeGT CountryCodeType = "GT"
	// GU - Guam
	CountryCodeGU CountryCodeType = "GU"
	// GW - Gu inea-Bissau
	CountryCodeGW CountryCodeType = "GW"
	// GY - Guyana
	CountryCodeGY CountryCodeType = "GY"
	// HK - Hong Kong
	CountryCodeHK CountryCodeType = "HK"
	// HM - Heard Island and McDonald Islands
	CountryCodeHM CountryCodeType = "HM"
	// HN - Honduras
	CountryCodeHN CountryCodeType = "HN"
	// HR - Croatia
	CountryCodeHR CountryCodeType = "HR"
	// HT - Haiti
	CountryCodeHT CountryCodeType = "HT"
	// HU - Hungary
	CountryCodeHU CountryCodeType = "HU"
	// ID - Indonesia
	CountryCodeID CountryCodeType = "ID"
	// IE - Ireland
	CountryCodeIE CountryCodeType = "IE"
	// IL - Israel
	CountryCodeIL CountryCodeType = "IL"
	// IM - Isle of Man
	CountryCodeIM CountryCodeType = "IM"
	// IN - India
	CountryCodeIN CountryCodeType = "IN"
	// IO - British Indian Ocean Territory
	CountryCodeIO CountryCodeType = "IO"
	// IQ - Iraq
	CountryCodeIQ CountryCodeType = "IQ"
	// IR - Iran, Islamic Republic of
	CountryCodeIR CountryCodeType = "IR"
	// IS - Iceland
	CountryCodeIS CountryCodeType = "IS"
	// IT - Italy
	CountryCodeIT CountryCodeType = "IT"
	// JE - Jersey
	CountryCodeJE CountryCodeType = "JE"
	// JM - Jamaica
	CountryCodeJM CountryCodeType = "JM"
	// JO - Jordan
	CountryCodeJO CountryCodeType = "JO"
	// JP - Japan
	CountryCodeJP CountryCodeType = "JP"
	// KE - Kenya
	CountryCodeKE CountryCodeType = "KE"
	// KG - Kyrgyzstan
	CountryCodeKG CountryCodeType = "KG"
	// KH - Cambodia
	CountryCodeKH CountryCodeType = "KH"
	// KI - Kiribati
	CountryCodeKI CountryCodeType = "KI"
	// KM - Comoros
	CountryCodeKM CountryCodeType = "KM"
	// KN - Saint Kitts and Nevis
	CountryCodeKN CountryCodeType = "KN"
	// KP - Korea, Democratic People's Republic of
	CountryCodeKP CountryCodeType = "KP"
	// KR - Korea, Republic of
	CountryCodeKR CountryCodeType = "KR"
	// KW - Kuwait
	CountryCodeKW CountryCodeType = "KW"
	// KY - Cayman Islands
	CountryCodeKY CountryCodeType = "KY"
	// KZ - Kazakhstan
	CountryCodeKZ CountryCodeType = "KZ"
	// LA - Lao People's Democratic Republic
	CountryCodeLA CountryCodeType = "LA"
	// LB - Lebanon
	CountryCodeLB CountryCodeType = "LB"
	// LC - Saint Lucia
	CountryCodeLC CountryCodeType = "LC"
	// LI - Liechtenstein
	CountryCodeLI CountryCodeType = "LI"
	// LK - Sri Lanka
	CountryCodeLK CountryCodeType = "LK"
	// LR - Liberia
	CountryCodeLR CountryCodeType = "LR"
	// LS - Lesotho
	CountryCodeLS CountryCodeType = "LS"
	// LT - Lithuania
	CountryCodeLT CountryCodeType = "LT"
	// LU - Luxembourg
	CountryCodeLU CountryCodeType = "LU"
	// LV - Latvia
	CountryCodeLV CountryCodeType = "LV"
	// LY - Libya
	CountryCodeLY CountryCodeType = "LY"
	// MA - Morocco
	CountryCodeMA CountryCodeType = "MA"
	// MC - Monaco
	CountryCodeMC CountryCodeType = "MC"
	// MD - Moldova, Republic of
	CountryCodeMD CountryCodeType = "MD"
	// ME - Montenegro
	CountryCodeME CountryCodeType = "ME"
	// MF - Saint Martin (French part)
	CountryCodeMF CountryCodeType = "MF"
	// MG - Madagascar
	CountryCodeMG CountryCodeType = "MG"
	// MH - Marshall Islands
	CountryCodeMH CountryCodeType = "MH"
	// MK - Macedonia, the former Yugoslav Republic of
	CountryCodeMK CountryCodeType = "MK"
	// ML - Mali
	CountryCodeML CountryCodeType = "ML"
	// MM - Myanmar
	CountryCodeMM CountryCodeType = "MM"
	// MN - Mongolia
	CountryCodeMN CountryCodeType = "MN"
	// MO - Macao
	CountryCodeMO CountryCodeType = "MO"
	// MP - Northern Mariana Islands
	CountryCodeMP CountryCodeType = "MP"
	// MQ - Martinique
	CountryCodeMQ CountryCodeType = "MQ"
	// MR - Mauritania
	CountryCodeMR CountryCodeType = "MR"
	// MS - Montserrat
	CountryCodeMS CountryCodeType = "MS"
	// MT - Malta
	CountryCodeMT CountryCodeType = "MT"
	// MU - Mauritius
	CountryCodeMU CountryCodeType = "MU"
	// MV - Maldives
	CountryCodeMV CountryCodeType = "MV"
	// MW - Malawi
	CountryCodeMW CountryCodeType = "MW"
	// MX - Mexico
	CountryCodeMX CountryCodeType = "MX"
	// MY - Malaysia
	CountryCodeMY CountryCodeType = "MY"
	// MZ - Mozambique
	CountryCodeMZ CountryCodeType = "MZ"
	// NA - Namibia
	CountryCodeNA CountryCodeType = "NA"
	// NC - New Caledonia
	CountryCodeNC CountryCodeType = "NC"
	// NE - Niger
	CountryCodeNE CountryCodeType = "NE"
	// NF - Norfolk Island
	CountryCodeNF CountryCodeType = "NF"
	// NG - Nigeria
	CountryCodeNG CountryCodeType = "NG"
	// NI - Nicaragua
	CountryCodeNI CountryCodeType = "NI"
	// NL - Netherlands
	CountryCodeNL CountryCodeType = "NL"
	// NO - Norway
	CountryCodeNO CountryCodeType = "NO"
	// NP - Nepal
	CountryCodeNP CountryCodeType = "NP"
	// NR - Nauru
	CountryCodeNR CountryCodeType = "NR"
	// NU - Niue
	CountryCodeNU CountryCodeType = "NU"
	// NZ - New Zealand
	CountryCodeNZ CountryCodeType = "NZ"
	// OM - Oman
	CountryCodeOM CountryCodeType = "OM"
	// PA - Panama
	CountryCodePA CountryCodeType = "PA"
	// PE - Peru
	CountryCodePE CountryCodeType = "PE"
	// PF - French Polynesia
	CountryCodePF CountryCodeType = "PF"
	// PG - Papua New Guinea
	CountryCodePG CountryCodeType = "PG"
	// PH - Philippines
	CountryCodePH CountryCodeType = "PH"
	// PK - Pakistan
	CountryCodePK CountryCodeType = "PK"
	// PL - Poland
	CountryCodePL CountryCodeType = "PL"
	// PM - Saint Pierre and Miquelon
	CountryCodePM CountryCodeType = "PM"
	// PN - Pitcairn
	CountryCodePN CountryCodeType = "PN"
	// PR - Puerto Rico
	CountryCodePR CountryCodeType = "PR"
	// PS - Palestine, State of
	CountryCodePS CountryCodeType = "PS"
	// PT - Portugal
	CountryCodePT CountryCodeType = "PT"
	// PW - Palau
	CountryCodePW CountryCodeType = "PW"
	// PY - Paraguay
	CountryCodePY CountryCodeType = "PY"
	// QA - Qatar
	CountryCodeQA CountryCodeType = "QA"
	// RE - R&#xE9;union
	CountryCodeRE CountryCodeType = "RE"
	// RS - Serbia
	CountryCodeRS CountryCodeType = "RS"
	// RU - Russian Federation
	CountryCodeRU CountryCodeType = "RU"
	// RW - Rwanda
	CountryCodeRW CountryCodeType = "RW"
	// SA - Saudi Arabia
	CountryCodeSA CountryCodeType = "SA"
	// SB - Solomon Islands
	CountryCodeSB CountryCodeType = "SB"
	// SC - Seychelles
	CountryCodeSC CountryCodeType = "SC"
	// SD - Sudan
	CountryCodeSD CountryCodeType = "SD"
	// SE - Sweden
	CountryCodeSE CountryCodeType = "SE"
	// SG - Singapore
	CountryCodeSG CountryCodeType = "SG"
	// SH - Saint Helena, Ascension and Tristan da Cunha
	CountryCodeSH CountryCodeType = "SH"
	// SI - Slovenia
	CountryCodeSI CountryCodeType = "SI"
	// SJ - Svalbard and Jan Mayen
	CountryCodeSJ CountryCodeType = "SJ"
	// SK - Slovakia
	CountryCodeSK CountryCodeType = "SK"
	// SL - Sierra Leone
	CountryCodeSL CountryCodeType = "SL"
	// SM - San Marino
	CountryCodeSM CountryCodeType = "SM"
	// SN - Senegal
	CountryCodeSN CountryCodeType = "SN"
	// SO - Somalia
	CountryCodeSO CountryCodeType = "SO"
	// SR - Suriname
	CountryCodeSR CountryCodeType = "SR"
	// SS - South Sudan
	CountryCodeSS CountryCodeType = "SS"
	// ST - Sao Tome and Principe
	CountryCodeST CountryCodeType = "ST"
	// SV - El Salvador
	CountryCodeSV CountryCodeType = "SV"
	// SX - Sint Maarten (Dutch part)
	CountryCodeSX CountryCodeType = "SX"
	// SY - Syrian Arab Republic
	CountryCodeSY CountryCodeType = "SY"
	// SZ - Swaziland
	CountryCodeSZ CountryCodeType = "SZ"
	// TC - Turks and Caicos Islands
	CountryCodeTC CountryCodeType = "TC"
	// TD - Chad
	CountryCodeTD CountryCodeType = "TD"
	// TF - French Southern Territories
	CountryCodeTF CountryCodeType = "TF"
	// TG - Togo
	CountryCodeTG CountryCodeType = "TG"
	// TH - Thailand
	CountryCodeTH CountryCodeType = "TH"
	// TJ - Tajikistan
	CountryCodeTJ CountryCodeType = "TJ"
	// TK - Tokelau
	CountryCodeTK CountryCodeType = "TK"
	// TL - Timor-Leste
	CountryCodeTL CountryCodeType = "TL"
	// TM - Turkmenistan
	CountryCodeTM CountryCodeType = "TM"
	// TN - Tunisia
	CountryCodeTN CountryCodeType = "TN"
	// TO - Tonga
	CountryCodeTO CountryCodeType = "TO"
	// TR - Turkey
	CountryCodeTR CountryCodeType = "TR"
	// TT - Trinidad and Tobago
	CountryCodeTT CountryCodeType = "TT"
	// TV - Tuvalu
	CountryCodeTV CountryCodeType = "TV"
	// TW - Taiwan, Province of China
	CountryCodeTW CountryCodeType = "TW"
	// TZ - Tanzania, United Republic of
	CountryCodeTZ CountryCodeType = "TZ"
	// UA - Ukraine
	CountryCodeUA CountryCodeType = "UA"
	// UG - Uganda
	CountryCodeUG CountryCodeType = "UG"
	// UM - United States Minor Outlying Islands
	CountryCodeUM CountryCodeType = "UM"
	// US - United States of America
	CountryCodeUS CountryCodeType = "US"
	// UY - Uruguay
	CountryCodeUY CountryCodeType = "UY"
	// UZ - Uzbekistan
	CountryCodeUZ CountryCodeType = "UZ"
	// VA - Holy See
	CountryCodeVA CountryCodeType = "VA"
	// VC - Saint Vincent and the Grenadines
	CountryCodeVC CountryCodeType = "VC"
	// VE - Venezuela, Bolivarian Republic of
	CountryCodeVE CountryCodeType = "VE"
	// VG - Virgin Islands, British
	CountryCodeVG CountryCodeType = "VG"
	// VI - Virgin Islands, U.S.
	CountryCodeVI CountryCodeType = "VI"
	// VN - Viet Nam
	CountryCodeVN CountryCodeType = "VN"
	// VU - Vanuatu
	CountryCodeVU CountryCodeType = "VU"
	// WF - Wallis and Futuna
	CountryCodeWF CountryCodeType = "WF"
	// WS - Samoa
	CountryCodeWS CountryCodeType = "WS"
	// YE - Yemen
	CountryCodeYE CountryCodeType = "YE"
	// YT - Mayotte
	CountryCodeYT CountryCodeType = "YT"
	// ZA - South Africa
	CountryCodeZA CountryCodeType = "ZA"
	// ZM - Zambia
	CountryCodeZM CountryCodeType = "ZM"
	// ZW - Zimbabwe
	CountryCodeZW CountryCodeType = "ZW"
	// 1A - Kosovo
	CountryCode1A CountryCodeType = "1A"
)

type CountrySubentityType string

const (
	// B - București
	CountrySubentityRO_B CountrySubentityType = "RO-B"
	// AB - Alba
	CountrySubentityRO_AB CountrySubentityType = "RO-AB"
	// AR - Arad
	CountrySubentityRO_AR CountrySubentityType = "RO-AR"
	// AG - Argeș
	CountrySubentityRO_AG CountrySubentityType = "RO-AG"
	// BC - Bacău
	CountrySubentityRO_BC CountrySubentityType = "RO-BC"
	// BH - Bihor
	CountrySubentityRO_BH CountrySubentityType = "RO-BH"
	// BN - Bistrița-Năsăud
	CountrySubentityRO_BN CountrySubentityType = "RO-BN"
	// BT - Botoșani
	CountrySubentityRO_BT CountrySubentityType = "RO-BT"
	// BR - Brăila
	CountrySubentityRO_BR CountrySubentityType = "RO-BR"
	// BV - Brașov
	CountrySubentityRO_BV CountrySubentityType = "RO-BV"
	// BZ - Buzău
	CountrySubentityRO_BZ CountrySubentityType = "RO-BZ"
	// CL - Călărași
	CountrySubentityRO_CL CountrySubentityType = "RO-CL"
	// CS - Caraș-Severin
	CountrySubentityRO_CS CountrySubentityType = "RO-CS"
	// CJ - Cluj
	CountrySubentityRO_CJ CountrySubentityType = "RO-CJ"
	// CT - Constanța
	CountrySubentityRO_CT CountrySubentityType = "RO-CT"
	// CV - Covasna
	CountrySubentityRO_CV CountrySubentityType = "RO-CV"
	// DB - Dâmbovița
	CountrySubentityRO_DB CountrySubentityType = "RO-DB"
	// DJ - Dolj
	CountrySubentityRO_DJ CountrySubentityType = "RO-DJ"
	// GL - Galați
	CountrySubentityRO_GL CountrySubentityType = "RO-GL"
	// GR - Giurgiu
	CountrySubentityRO_GR CountrySubentityType = "RO-GR"
	// GJ - Gorj
	CountrySubentityRO_GJ CountrySubentityType = "RO-GJ"
	// HR - Harghita
	CountrySubentityRO_HR CountrySubentityType = "RO-HR"
	// HD - Hunedoara
	CountrySubentityRO_HD CountrySubentityType = "RO-HD"
	// IL - Ialomița
	CountrySubentityRO_IL CountrySubentityType = "RO-IL"
	// IS - Iași
	CountrySubentityRO_IS CountrySubentityType = "RO-IS"
	// IF - Ilfov
	CountrySubentityRO_IF CountrySubentityType = "RO-IF"
	// MM - Maramureș
	CountrySubentityRO_MM CountrySubentityType = "RO-MM"
	// MH - Mehedinți
	CountrySubentityRO_MH CountrySubentityType = "RO-MH"
	// MS - Mureș
	CountrySubentityRO_MS CountrySubentityType = "RO-MS"
	// NT - Neamț
	CountrySubentityRO_NT CountrySubentityType = "RO-NT"
	// OT - Olt
	CountrySubentityRO_OT CountrySubentityType = "RO-OT"
	// PH - Prahova
	CountrySubentityRO_PH CountrySubentityType = "RO-PH"
	// SJ - Sălaj
	CountrySubentityRO_SJ CountrySubentityType = "RO-SJ"
	// SM - Satu Mare
	CountrySubentityRO_SM CountrySubentityType = "RO-SM"
	// SB - Sibiu
	CountrySubentityRO_SB CountrySubentityType = "RO-SB"
	// SV - Suceava
	CountrySubentityRO_SV CountrySubentityType = "RO-SV"
	// TR - Teleorman
	CountrySubentityRO_TR CountrySubentityType = "RO-TR"
	// TM - Timiș
	CountrySubentityRO_TM CountrySubentityType = "RO-TM"
	// TL - Tulcea
	CountrySubentityRO_TL CountrySubentityType = "RO-TL"
	// VS - Vaslui
	CountrySubentityRO_VS CountrySubentityType = "RO-VS"
	// VL - Vâlcea
	CountrySubentityRO_VL CountrySubentityType = "RO-VL"
	// VN - Vrancea
	CountrySubentityRO_VN CountrySubentityType = "RO-VN"
)

// If the country code for a postal address is RO-B, then the City name must be
// one of the following values.
const (
	CityNameROBSector1 = "SECTOR1"
	CityNameROBSector2 = "SECTOR2"
	CityNameROBSector3 = "SECTOR3"
	CityNameROBSector4 = "SECTOR4"
	CityNameROBSector5 = "SECTOR5"
	CityNameROBSector6 = "SECTOR6"
)

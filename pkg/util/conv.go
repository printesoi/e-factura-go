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

// Package util provides utilities functions for allowing efactura <>
// etransport interoperability.
package util

import (
	"github.com/printesoi/e-factura-go/pkg/efactura"
	"github.com/printesoi/e-factura-go/pkg/etransport"
)

// EfacturaRoCountrySubentityToEtransportCountyCode convert the given
// efactura.CountrySubentityType to an etransport.CountyCode. If
// efacturaSubentity is a valid subentity code, the second returned value will
// be true.
func EfacturaRoCountrySubentityToEtransportCountyCode(efacturaSubentity efactura.CountrySubentityType) (_ etransport.CountyCodeType, ok bool) {
	switch efacturaSubentity {
	case efactura.CountrySubentityRO_B:
		return etransport.CountyCodeB, true
	case efactura.CountrySubentityRO_AB:
		return etransport.CountyCodeAB, true
	case efactura.CountrySubentityRO_AR:
		return etransport.CountyCodeAR, true
	case efactura.CountrySubentityRO_AG:
		return etransport.CountyCodeAG, true
	case efactura.CountrySubentityRO_BC:
		return etransport.CountyCodeBC, true
	case efactura.CountrySubentityRO_BH:
		return etransport.CountyCodeBH, true
	case efactura.CountrySubentityRO_BN:
		return etransport.CountyCodeBN, true
	case efactura.CountrySubentityRO_BT:
		return etransport.CountyCodeBT, true
	case efactura.CountrySubentityRO_BR:
		return etransport.CountyCodeBR, true
	case efactura.CountrySubentityRO_BV:
		return etransport.CountyCodeBV, true
	case efactura.CountrySubentityRO_BZ:
		return etransport.CountyCodeBZ, true
	case efactura.CountrySubentityRO_CL:
		return etransport.CountyCodeCL, true
	case efactura.CountrySubentityRO_CS:
		return etransport.CountyCodeCS, true
	case efactura.CountrySubentityRO_CJ:
		return etransport.CountyCodeCJ, true
	case efactura.CountrySubentityRO_CT:
		return etransport.CountyCodeCT, true
	case efactura.CountrySubentityRO_CV:
		return etransport.CountyCodeCV, true
	case efactura.CountrySubentityRO_DB:
		return etransport.CountyCodeDB, true
	case efactura.CountrySubentityRO_DJ:
		return etransport.CountyCodeDJ, true
	case efactura.CountrySubentityRO_GL:
		return etransport.CountyCodeGL, true
	case efactura.CountrySubentityRO_GR:
		return etransport.CountyCodeGR, true
	case efactura.CountrySubentityRO_GJ:
		return etransport.CountyCodeGJ, true
	case efactura.CountrySubentityRO_HR:
		return etransport.CountyCodeHR, true
	case efactura.CountrySubentityRO_HD:
		return etransport.CountyCodeHD, true
	case efactura.CountrySubentityRO_IL:
		return etransport.CountyCodeIL, true
	case efactura.CountrySubentityRO_IS:
		return etransport.CountyCodeIS, true
	case efactura.CountrySubentityRO_IF:
		return etransport.CountyCodeIF, true
	case efactura.CountrySubentityRO_MM:
		return etransport.CountyCodeMM, true
	case efactura.CountrySubentityRO_MH:
		return etransport.CountyCodeMH, true
	case efactura.CountrySubentityRO_MS:
		return etransport.CountyCodeMS, true
	case efactura.CountrySubentityRO_NT:
		return etransport.CountyCodeNT, true
	case efactura.CountrySubentityRO_OT:
		return etransport.CountyCodeOT, true
	case efactura.CountrySubentityRO_PH:
		return etransport.CountyCodePH, true
	case efactura.CountrySubentityRO_SJ:
		return etransport.CountyCodeSJ, true
	case efactura.CountrySubentityRO_SM:
		return etransport.CountyCodeSM, true
	case efactura.CountrySubentityRO_SB:
		return etransport.CountyCodeSB, true
	case efactura.CountrySubentityRO_SV:
		return etransport.CountyCodeSV, true
	case efactura.CountrySubentityRO_TR:
		return etransport.CountyCodeTR, true
	case efactura.CountrySubentityRO_TM:
		return etransport.CountyCodeTM, true
	case efactura.CountrySubentityRO_TL:
		return etransport.CountyCodeTL, true
	case efactura.CountrySubentityRO_VS:
		return etransport.CountyCodeVS, true
	case efactura.CountrySubentityRO_VL:
		return etransport.CountyCodeVL, true
	case efactura.CountrySubentityRO_VN:
		return etransport.CountyCodeVN, true
	}
	return "", false
}

// EtransportRoCountyCodeToEfacturaCountrySubentity convert the given
// etransport.CountyCodeType to an efactura.CountrySubentityType. If the given
// county code is valid, the second returned value will be true.
func EtransportRoCountyCodeToEfacturaCountrySubentity(etransportCountyCode etransport.CountyCodeType) (_ efactura.CountrySubentityType, ok bool) {
	switch etransportCountyCode {
	case etransport.CountyCodeB:
		return efactura.CountrySubentityRO_B, true
	case etransport.CountyCodeAB:
		return efactura.CountrySubentityRO_AB, true
	case etransport.CountyCodeAR:
		return efactura.CountrySubentityRO_AR, true
	case etransport.CountyCodeAG:
		return efactura.CountrySubentityRO_AG, true
	case etransport.CountyCodeBC:
		return efactura.CountrySubentityRO_BC, true
	case etransport.CountyCodeBH:
		return efactura.CountrySubentityRO_BH, true
	case etransport.CountyCodeBN:
		return efactura.CountrySubentityRO_BN, true
	case etransport.CountyCodeBT:
		return efactura.CountrySubentityRO_BT, true
	case etransport.CountyCodeBR:
		return efactura.CountrySubentityRO_BR, true
	case etransport.CountyCodeBV:
		return efactura.CountrySubentityRO_BV, true
	case etransport.CountyCodeBZ:
		return efactura.CountrySubentityRO_BZ, true
	case etransport.CountyCodeCL:
		return efactura.CountrySubentityRO_CL, true
	case etransport.CountyCodeCS:
		return efactura.CountrySubentityRO_CS, true
	case etransport.CountyCodeCJ:
		return efactura.CountrySubentityRO_CJ, true
	case etransport.CountyCodeCT:
		return efactura.CountrySubentityRO_CT, true
	case etransport.CountyCodeCV:
		return efactura.CountrySubentityRO_CV, true
	case etransport.CountyCodeDB:
		return efactura.CountrySubentityRO_DB, true
	case etransport.CountyCodeDJ:
		return efactura.CountrySubentityRO_DJ, true
	case etransport.CountyCodeGL:
		return efactura.CountrySubentityRO_GL, true
	case etransport.CountyCodeGR:
		return efactura.CountrySubentityRO_GR, true
	case etransport.CountyCodeGJ:
		return efactura.CountrySubentityRO_GJ, true
	case etransport.CountyCodeHR:
		return efactura.CountrySubentityRO_HR, true
	case etransport.CountyCodeHD:
		return efactura.CountrySubentityRO_HD, true
	case etransport.CountyCodeIL:
		return efactura.CountrySubentityRO_IL, true
	case etransport.CountyCodeIS:
		return efactura.CountrySubentityRO_IS, true
	case etransport.CountyCodeIF:
		return efactura.CountrySubentityRO_IF, true
	case etransport.CountyCodeMM:
		return efactura.CountrySubentityRO_MM, true
	case etransport.CountyCodeMH:
		return efactura.CountrySubentityRO_MH, true
	case etransport.CountyCodeMS:
		return efactura.CountrySubentityRO_MS, true
	case etransport.CountyCodeNT:
		return efactura.CountrySubentityRO_NT, true
	case etransport.CountyCodeOT:
		return efactura.CountrySubentityRO_OT, true
	case etransport.CountyCodePH:
		return efactura.CountrySubentityRO_PH, true
	case etransport.CountyCodeSJ:
		return efactura.CountrySubentityRO_SJ, true
	case etransport.CountyCodeSM:
		return efactura.CountrySubentityRO_SM, true
	case etransport.CountyCodeSB:
		return efactura.CountrySubentityRO_SB, true
	case etransport.CountyCodeSV:
		return efactura.CountrySubentityRO_SV, true
	case etransport.CountyCodeTR:
		return efactura.CountrySubentityRO_TR, true
	case etransport.CountyCodeTM:
		return efactura.CountrySubentityRO_TM, true
	case etransport.CountyCodeTL:
		return efactura.CountrySubentityRO_TL, true
	case etransport.CountyCodeVS:
		return efactura.CountrySubentityRO_VS, true
	case etransport.CountyCodeVL:
		return efactura.CountrySubentityRO_VL, true
	case etransport.CountyCodeVN:
		return efactura.CountrySubentityRO_VN, true
	}
	return "", false
}

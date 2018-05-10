// Code generated by "stringer -type=OpCode isainfo.go"; DO NOT EDIT.

package straightISAv2Info

import "fmt"

const _OpCode_name = "OpNOPOpSYSCALLOpSYSRETOpJOpJROpJALOpJRALOpBEZOpBNZOpADDOpADDiOpSUBOpSUBiOpMULOpMULiOpDIVOpDIViOpDIVUOpDIVUiOpMODOpMODiOpMODUOpMODUiOpSLTOpSLTiOpSLTUOpSLTUiOpFTOIOpITOFOpFADDOpFSUBOpFMULOpFDIVOpFSLTOpSHLOpSHLiOpSHROpSHRiOpSHRAOpSHRAiOpANDOpANDiOpOROpORiOpXOROpXORiOpLUiOpSPADDiOpRPINCOpRMOVOpLDOpLDHOpLDHUOpLDBOpLDBUOpSTOpSTHOpSTBOpSEXT16TO32OpSEXT8TO32OpZEXT16TO32OpZEXT8TO32Op_MAX"

var _OpCode_index = [...]uint16{0, 5, 14, 22, 25, 29, 34, 40, 45, 50, 55, 61, 66, 72, 77, 83, 88, 94, 100, 107, 112, 118, 124, 131, 136, 142, 148, 155, 161, 167, 173, 179, 185, 191, 197, 202, 208, 213, 219, 225, 232, 237, 243, 247, 252, 257, 263, 268, 276, 283, 289, 293, 298, 304, 309, 315, 319, 324, 329, 341, 352, 364, 375, 381}

func (i OpCode) String() string {
	if i < 0 || i >= OpCode(len(_OpCode_index)-1) {
		return fmt.Sprintf("OpCode(%d)", i)
	}
	return _OpCode_name[_OpCode_index[i]:_OpCode_index[i+1]]
}

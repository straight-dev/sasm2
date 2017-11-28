package straightISAv1Info

import (
	"errors"
	"strconv"
)

//go:generate stringer -type=InstType isainfo.go
//go:generate stringer -type=OpCode isainfo.go

// InstType : ZeroReg/OneReg/TwoReg
type InstType int

// InstType
const (
	ZeroReg InstType = iota
	OneReg
	TwoReg
)

// OpCode : オペコード
type OpCode int

// OpCode
const (
	OpNOP OpCode = iota
	OpJ
	OpJR
	OpJAL
	OpBEZ
	OpBNZ
	OpFPBEZ
	OpFPBNZ
	OpADD
	OpADDi
	OpSUB
	OpSUBi
	OpFTOI
	OpITOF
	OpFADD
	OpFADDi // deprecated
	OpFSUB
	OpFSUBi // deprecated
	OpFMUL
	OpFMULi // deprecated
	OpFDIV
	OpFDIVi // deprecated
	OpSHL
	OpSHLi
	OpSHR
	OpSHRi
	OpAND
	OpANDi
	OpOR
	OpORi
	OpXOR
	OpXORi
	OpSPADD // deprecated
	OpSPADDi
	OpFPADD
	OpFPADDi
	OpLD
	OpSPLD
	OpFPLD
	OpST
	OpSPST
	OpFPST
	OpRPINC
	OpSLT
	OpRMOV
	OpMUL
	OpMULi
	OpGPADD
	OpGPADDi
	OpGPLD
	OpGPST
	OpLDH
	OpLDHU
	OpLDQ
	OpLDQU
	OpLDB
	OpLDBU
	OpSTH
	OpSTQ
	OpSTB
	OpLUi
	OpFSLT
	OpHALT
	Op_MAX // NOT OP
)

func GetInstType(oc OpCode) (InstType, bool, error) {
	// the second return value: immediate value to be sign-extended
	switch oc {
	case OpNOP:
		return ZeroReg, false, nil
	case OpJ:
		return ZeroReg, true, nil
	case OpJR:
		return OneReg, true, nil
	case OpJAL:
		return ZeroReg, false, nil
	case OpBEZ:
		return OneReg, true, nil
	case OpBNZ:
		return OneReg, true, nil
	case OpFPBEZ:
		return OneReg, true, nil
	case OpFPBNZ:
		return OneReg, true, nil
	case OpADD:
		return TwoReg, false, nil
	case OpADDi:
		return OneReg, true, nil
	case OpSUB:
		return TwoReg, false, nil
	case OpSUBi:
		return OneReg, true, nil
	case OpFTOI:
		return OneReg, false, nil
	case OpITOF:
		return OneReg, false, nil
	case OpFADD:
		return TwoReg, false, nil
	case OpFADDi:
		return OneReg, false, nil
	case OpFSUB:
		return TwoReg, false, nil
	case OpFSUBi:
		return OneReg, false, nil
	case OpFMUL:
		return TwoReg, false, nil
	case OpFMULi:
		return OneReg, false, nil
	case OpFDIV:
		return TwoReg, false, nil
	case OpFDIVi:
		return OneReg, false, nil
	case OpSHL:
		return TwoReg, false, nil
	case OpSHLi:
		return OneReg, false, nil
	case OpSHR:
		return TwoReg, false, nil
	case OpSHRi:
		return OneReg, false, nil
	case OpAND:
		return TwoReg, false, nil
	case OpANDi:
		return OneReg, false, nil
	case OpOR:
		return TwoReg, false, nil
	case OpORi:
		return OneReg, false, nil
	case OpXOR:
		return TwoReg, false, nil
	case OpXORi:
		return OneReg, false, nil
	case OpSPADD:
		return OneReg, false, nil
	case OpSPADDi:
		return ZeroReg, true, nil
	case OpFPADD:
		return OneReg, false, nil
	case OpFPADDi:
		return ZeroReg, true, nil
	case OpLD:
		return OneReg, true, nil
	case OpSPLD:
		return ZeroReg, true, nil
	case OpFPLD:
		return ZeroReg, true, nil
	case OpST:
		return TwoReg, true, nil
	case OpSPST:
		return OneReg, true, nil
	case OpFPST:
		return OneReg, true, nil
	case OpRPINC:
		return ZeroReg, false, nil
	case OpSLT:
		return TwoReg, false, nil
	case OpRMOV:
		return OneReg, false, nil
	case OpMUL:
		return TwoReg, false, nil
	case OpMULi:
		return OneReg, true, nil
	case OpGPADD:
		return OneReg, true, nil
	case OpGPADDi:
		return ZeroReg, true, nil
	case OpGPLD:
		return OneReg, true, nil
	case OpGPST:
		return ZeroReg, true, nil
	case OpLDH:
		return OneReg, true, nil
	case OpLDHU:
		return OneReg, true, nil
	case OpLDQ:
		return OneReg, true, nil
	case OpLDQU:
		return OneReg, true, nil
	case OpLDB:
		return OneReg, true, nil
	case OpLDBU:
		return OneReg, true, nil
	case OpSTH:
		return TwoReg, true, nil
	case OpSTQ:
		return TwoReg, true, nil
	case OpSTB:
		return TwoReg, true, nil
	case OpLUi:
		return ZeroReg, false, nil
	case OpFSLT:
		return TwoReg, false, nil
	case OpHALT:
		return ZeroReg, false, nil
	default:
		return ZeroReg, false, errors.New(strconv.Itoa(int(oc)) + " is invalid opcode")
	}
}

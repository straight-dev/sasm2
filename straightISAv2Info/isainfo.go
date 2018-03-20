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
	OpSYSCALL
	OpSYSRET
	OpJ
	OpJR
	OpJAL
	OpJRAL
	OpBEZ
	OpBNZ
	OpADD
	OpADDi
	OpSUB
	OpSUBi
	OpMUL
	OpMULi
	OpDIV
	OpDIVi
	OpDIVU
	OpDIVUi
	OpMOD
	OpMODi
	OpMODU
	OpMODUi
	OpSLT
	OpSLTi
	OpSLTU
	OpSLTUi
	OpFTOI
	OpITOF
	OpFADD
	OpFSUB
	OpFMUL
	OpFDIV
	OpFSLT
	OpSHL
	OpSHLi
	OpSHR
	OpSHRi
	OpSHRA
	OpSHRAi
	OpAND
	OpANDi
	OpOR
	OpORi
	OpXOR
	OpXORi
	OpLUi
	OpSPADDi
	OpRPINC
	OpRMOV
	OpLD
	OpLDH
	OpLDHU
	OpLDB
	OpLDBU
	OpST
	OpSTH
	OpSTB
	OpSEXT16TO32
	OpSEXT8TO32
	OpZEXT16TO32
	OpZEXT8TO32

	Op_MAX // NOT OP
)

func GetInstType(oc OpCode) (InstType, bool, error) {
	// the second return value: immediate value to be sign-extended
	switch oc {
	case OpNOP:
		return ZeroReg, false, nil
	case OpSYSCALL:
		return ZeroReg, false, nil
	case OpSYSRET:
		return ZeroReg, false, nil
	case OpJ:
		return ZeroReg, true, nil
	case OpJR:
		return OneReg, false, nil
	case OpJAL:
		return ZeroReg, true, nil
	case OpJRAL:
		return OneReg, false, nil
	case OpBEZ:
		return OneReg, true, nil
	case OpBNZ:
		return OneReg, true, nil
	case OpADD:
		return TwoReg, false, nil
	case OpADDi:
		return OneReg, true, nil
	case OpSUB:
		return TwoReg, false, nil
	case OpSUBi:
		return OneReg, true, nil
	case OpMUL:
		return TwoReg, false, nil
	case OpMULi:
		return OneReg, true, nil
	case OpDIV:
		return TwoReg, false, nil
	case OpDIVi:
		return OneReg, true, nil
	case OpDIVU:
		return TwoReg, false, nil
	case OpDIVUi:
		return OneReg, false, nil
	case OpMOD:
		return TwoReg, false, nil
	case OpMODi:
		return OneReg, true, nil
	case OpMODU:
		return TwoReg, false, nil
	case OpMODUi:
		return OneReg, false, nil
	case OpSLT:
		return TwoReg, false, nil
	case OpSLTi:
		return OneReg, true, nil
	case OpSLTU:
		return TwoReg, false, nil
	case OpSLTUi:
		return OneReg, false, nil
	case OpFTOI:
		return OneReg, false, nil
	case OpITOF:
		return OneReg, false, nil
	case OpFADD:
		return TwoReg, false, nil
	case OpFSUB:
		return TwoReg, false, nil
	case OpFMUL:
		return TwoReg, false, nil
	case OpFDIV:
		return TwoReg, false, nil
	case OpFSLT:
		return TwoReg, false, nil
	case OpSHL:
		return TwoReg, false, nil
	case OpSHLi:
		return OneReg, false, nil
	case OpSHR:
		return TwoReg, false, nil
	case OpSHRi:
		return OneReg, false, nil
	case OpSHRA:
		return TwoReg, false, nil
	case OpSHRAi:
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
	case OpLUi:
		return ZeroReg, false, nil
	case OpSPADDi:
		return ZeroReg, true, nil
	case OpRPINC:
		return ZeroReg, false, nil
	case OpRMOV:
		return OneReg, false, nil
	case OpLD:
		return OneReg, true, nil
	case OpLDH:
		return OneReg, true, nil
	case OpLDHU:
		return OneReg, true, nil
	case OpLDB:
		return OneReg, true, nil
	case OpLDBU:
		return OneReg, true, nil
	case OpST:
		return TwoReg, true, nil
	case OpSTH:
		return TwoReg, true, nil
	case OpSTB:
		return TwoReg, true, nil
	case OpSEXT16TO32:
		return OneReg, false, nil
	case OpSEXT8TO32:
		return OneReg, false, nil
	case OpZEXT16TO32:
		return OneReg, false, nil
	case OpZEXT8TO32:
		return OneReg, false, nil

	default:
		return ZeroReg, false, errors.New(strconv.Itoa(int(oc)) + " is invalid opcode")
	}
}

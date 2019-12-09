package main

import (
	"fmt"
	"strconv"
	"strings"
)

type oneRegOperation uint32 // 13 (= 7 + 3 + 3) bit OPCODE

const (
	// opNOP     oneRegOperation = 15   // 000_000_0001111 // NOP = RPINC 0
	opRPINC  oneRegOperation = 15   // 000_000_0001111
	opFENCE  oneRegOperation = 2063 // 010_000_0001111
	opFENCEI oneRegOperation = 3087 // 011_000_0001111
	opJR     oneRegOperation = 143  // 000_001_0001111
	opJALR   oneRegOperation = 1167 // 001_001_0001111
	opECALL  oneRegOperation = 271  // 000_010_0001111
	// opEBREAK  oneRegOperation = 271  // 000_010_0001111
	opCSRRW   oneRegOperation = 1295 // 001_010_0001111
	opCSRRS   oneRegOperation = 2319 // 010_010_0001111
	opCSRRC   oneRegOperation = 3343 // 011_010_0001111
	opCSRRWi  oneRegOperation = 5391 // 101_000_0001111
	opCSRRSi  oneRegOperation = 6415 // 110_010_0001111
	opCSRRCi  oneRegOperation = 7439 // 111_010_0001111
	opSPLD8   oneRegOperation = 527  // 000_100_0001111
	opSPLD16  oneRegOperation = 1551 // 001_100_0001111
	opSPLD32  oneRegOperation = 2575 // 010_100_0001111
	opSPLD64  oneRegOperation = 3599 // 011_100_0001111
	opSPLD8u  oneRegOperation = 4623 // 100_100_0001111
	opSPLD16u oneRegOperation = 5647 // 101_100_0001111
	opSPLD32u oneRegOperation = 6671 // 110_100_0001111
	opSPLD32f oneRegOperation = 7695 // 111_100_0001111
	opSPST8   oneRegOperation = 655  // 000_101_0001111
	opSPST16  oneRegOperation = 1679 // 001_101_0001111
	opSPST32  oneRegOperation = 2703 // 010_101_0001111
	opSPST64  oneRegOperation = 3727 // 011_101_0001111
	opLD8     oneRegOperation = 783  // 000_110_0001111
	opLD16    oneRegOperation = 1807 // 001_110_0001111
	opLD32    oneRegOperation = 2831 // 010_110_0001111
	opLD64    oneRegOperation = 3855 // 011_110_0001111
	opLD8u    oneRegOperation = 4879 // 100_110_0001111
	opLD16u   oneRegOperation = 5903 // 101_110_0001111
	opLD32u   oneRegOperation = 6927 // 110_110_0001111
	opLD32f   oneRegOperation = 7951 // 111_110_0001111

	opADDi32  oneRegOperation = 4175 // 10_000_0_1001111
	opSLTi32  oneRegOperation = 4687 // 10_010_0_1001111
	opSLTiu32 oneRegOperation = 4943 // 10_011_0_1001111
	opXORi32  oneRegOperation = 5199 // 10_100_0_1001111
	opORi32   oneRegOperation = 5711 // 10_110_0_1001111
	opANDi32  oneRegOperation = 5967 // 10_111_0_1001111
	opSLLi32  oneRegOperation = 4431 // 10_001_0_1001111
	opSRLi32  oneRegOperation = 5455 // 10_101_0_1001111
	// opSRAi32  oneRegOperation = 5455 // 10_101_0_1001111
	opADDi64  oneRegOperation = 4303    // 10_000_1_1001111
	opRMOV    oneRegOperation = 9999999 // RMOV [x] = ADDi.64 [x] 0
	opSLTi64  oneRegOperation = 4815    // 10_010_1_1001111
	opSLTiu64 oneRegOperation = 5071    // 10_011_1_1001111
	opXORi64  oneRegOperation = 5327    // 10_100_1_1001111
	opORi64   oneRegOperation = 5839    // 10_110_1_1001111
	opANDi64  oneRegOperation = 6095    // 10_111_1_1001111
	opSLLi64  oneRegOperation = 4559    // 10_001_1_1001111
	opSRLi64  oneRegOperation = 5583    // 10_101_1_1001111
	// opSRAi64  oneRegOperation = 5583 // 10_101_1_1001111

)

type instTypeOneReg struct {
	operation oneRegOperation // 13 bit
	imm12     uint32          // 12 bit (Imm, CSR, 0 or 1)
	srcReg    uint32          // 7 bit
}

var strToOneRegOperation = map[string]oneRegOperation{
	// Specifications: p1
	"NOP":      opRPINC, // NOP = RPINC 0
	"RPINC":    opRPINC,
	"FENCE":    opFENCE,
	"FENCE.I":  opFENCEI,
	"JR":       opJR,
	"JALR":     opJALR,
	"ECALL":    opECALL, // imm = 0
	"EBREAK":   opECALL, // imm = 1
	"CSRRW":    opCSRRW,
	"CSRRS":    opCSRRS,
	"CSRRC":    opCSRRC,
	"CSRRWi":   opCSRRWi,
	"CSRRSi":   opCSRRSi,
	"CSRRCi":   opCSRRCi,
	"SPLD.8":   opSPLD8,
	"SPLD.16":  opSPLD16,
	"SPLD.32":  opSPLD32,
	"SPLD.64":  opSPLD64,
	"SPLD.8u":  opSPLD8u,
	"SPLD.16u": opSPLD16u,
	"SPLD.32u": opSPLD32u,
	"SPLD.f32": opSPLD32f,
	"SPST.8":   opSPST8,
	"SPST.16":  opSPST16,
	"SPST.32":  opSPST32,
	"SPST.64":  opSPST64,
	"LD.8":     opLD8,
	"LD.16":    opLD16,
	"LD.32":    opLD32,
	"LD.64":    opLD64,
	"LD.8u":    opLD8u,
	"LD.16u":   opLD16u,
	"LD.32u":   opLD32u,
	"LD.f32":   opLD32f,

	// Specifications: p4
	"ADDi.32":     opADDi32,
	"SLTi.32":     opSLTi32,
	"SLTiu.32":    opSLTiu32,
	"XORi.32":     opXORi32,
	"ORi.32":      opORi32,
	"ANDi.32":     opANDi32,
	"SLLi.32":     opSLLi32,
	"SRLi.32":     opSRLi32, // imm = 0_xxxxxx_00000
	"SRAi.32":     opSRLi32, // imm = 0_xxxxxx_01000
	"ADDi.64":     opADDi64,
	"RMOV":        opRMOV,
	"BITCASTITOD": opRMOV, //
	"SLTi.64":     opSLTi64,
	"SLTiu.64":    opSLTiu64,
	"XORi.64":     opXORi64,
	"ORi.64":      opORi64,
	"ANDi.64":     opANDi64,
	"SLLi.64":     opSLLi64,
	"SRLi.64":     opSRLi64, // imm = 0_xxxxxx_00000
	"SRAi.64":     opSRLi64, // imm = 0_xxxxxx_01000
}

func (i *instTypeOneReg) toUInt32() uint32 {
	op := i.operation
	if i.operation == opRMOV {
		op = opADDi64
	}
	return uint32(op) | (i.imm12 << 13) | (i.srcReg << 25)
}

func fromStringToInstTypeOneReg(str string) (*instTypeOneReg, error) {
	ss := strings.Fields(str)
	i := instTypeOneReg{}

	op, ok := strToOneRegOperation[ss[0]]
	if !ok {
		return nil, fmt.Errorf("not found '%s' in strToOneRegOperation(%v) :'%s'", ss[0], strToOneRegOperation, str)
	}
	i.operation = op

	switch op {
	case opRPINC:
		if ss[0] == "RPINC" {
			t, err := strconv.ParseUint(ss[1], 10, 7)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s (OneReg instruction '%s' opRPINC): %s", ss[1], str, err)
			}
			i.srcReg = uint32(t)
		} // else -> NOP

	case opFENCE: // FENCE succ(4bit) pred(4bit)
		{
			if len(ss) < 3 {
				return nil, fmt.Errorf("'FENCE succ pred'")
			}
			succ, err := strconv.ParseUint(ss[1], 10, 4)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s (OneReg instruction '%s' opFENCE): %s", ss[1], str, err)
			}
			pred, err := strconv.ParseUint(ss[2], 10, 4)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s (OneReg instruction '%s' opFENCE): %s", ss[2], str, err)
			}
			i.imm12 = uint32(pred | (succ << 4))
		}

	case opFENCEI:
		break

	case opECALL:
		if ss[0] == "EBREAK" {
			i.imm12 = 1
		} // else -> ECALL (i.imm = 1)

	case opJR, opJALR,
		opCSRRW, opCSRRS, opCSRRC, opCSRRWi, opCSRRSi, opCSRRCi,
		opLD8, opLD16, opLD32, opLD64, opLD8u, opLD16u, opLD32u, opLD32f,
		opADDi32, opSLTi32, opSLTiu32, opXORi32, opORi32, opANDi32,
		opADDi64, opSLTi64, opSLTiu64, opXORi64, opORi64, opANDi64,
		opSLLi32, opSRLi32, opSLLi64, opSRLi64,
		opSPST8, opSPST16, opSPST32, opSPST64:
		{
			if len(ss) < 3 {
				return nil, fmt.Errorf("too few arg: %s", str)
			}
			srcReg1, err := strconv.ParseUint(ss[1], 10, 7) // srcReg1 or zImm
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s (OneReg instruction '%s'): %s", ss[1], str, err)
			}
			i.srcReg = uint32(srcReg1)
			imm, err := strconv.ParseInt(ss[2], 10, 12) // Imm or CSR
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s (OneReg instruction '%s'): %s", ss[2], str, err)
			}

			if isShift(op) {
				if imm != imm&0x3f {
					return nil, fmt.Errorf("invalid immediate %d (Shift operation's imm: 6bit)", imm)
				}
				var funct uint32
				if ss[0] == "SRAi.32" || ss[0] == "SRAi.64" {
					funct = 1 << 3
				}
				i.imm12 = uint32(imm&0x3f)<<5 | funct
			} else {
				i.imm12 = uint32(imm) & 0xfff
			}
		}

	case opSPLD8, opSPLD16, opSPLD32, opSPLD64,
		opSPLD8u, opSPLD16u, opSPLD32u, opSPLD32f:
		{
			if len(ss) < 2 {
				return nil, fmt.Errorf("too few arg: %s", str)
			}
			imm, err := strconv.ParseInt(ss[1], 10, 12) // Imm
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s (OneReg instruction '%s' opSPLD): %s", ss[1], str, err)
			}
			i.imm12 = uint32(imm & (1<<12 - 1))
		}

	case opRMOV:
		{
			if len(ss) < 2 {
				return nil, fmt.Errorf("too few arg: %s", str)
			}
			srcReg1, err := strconv.ParseUint(ss[1], 10, 7) // srcReg1 or zImm
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s (OneReg instruction '%s'): %s", ss[1], str, err)
			}
			i.srcReg = uint32(srcReg1)
		}
	default:
		println("op ", op)
		panic("can't reach here...exhaustive switch in fromStringToInstTypeOneReg")
	}
	return &i, nil
}

func isShift(op oneRegOperation) bool {
	switch op {
	case opSLLi32, opSRLi32, opSLLi64, opSRLi64:
		return true
	default:
		return false
	}
}

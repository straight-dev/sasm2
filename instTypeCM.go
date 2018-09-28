package main

import (
	"fmt"
	"strconv"
	"strings"
)

// CM = Control + Memory(SPLD/SPST/LD)
type cmOperation uint32 // 13 (= 7 + 3 + 3) bit OPCODE

const (
	// opNOP     cmOperation = 15   // 000_000_0001111 // NOP = RPINC 0
	opRPINC  cmOperation = 15   // 000_000_0001111
	opFENCE  cmOperation = 2063 // 010_000_0001111
	opFENCEI cmOperation = 3087 // 011_000_0001111
	opJR     cmOperation = 143  // 000_001_0001111
	opJALR   cmOperation = 1167 // 001_001_0001111
	opECALL  cmOperation = 271  // 000_010_0001111
	// opEBREAK  cmOperation = 271  // 000_010_0001111
	opCSRRW   cmOperation = 1295 // 001_010_0001111
	opCSRRS   cmOperation = 2319 // 010_010_0001111
	opCSRRC   cmOperation = 3343 // 011_010_0001111
	opCSRRWi  cmOperation = 5135 // 101_000_0001111
	opCSRRSi  cmOperation = 6415 // 110_010_0001111
	opCSRRCi  cmOperation = 7439 // 111_010_0001111
	opSPLD8   cmOperation = 527  // 000_100_0001111
	opSPLD16  cmOperation = 1551 // 001_100_0001111
	opSPLD32  cmOperation = 2575 // 010_100_0001111
	opSPLD64  cmOperation = 3599 // 011_100_0001111
	opSPLD8u  cmOperation = 4623 // 100_100_0001111
	opSPLD16u cmOperation = 5647 // 101_100_0001111
	opSPLD32u cmOperation = 6671 // 110_100_0001111
	opSPST8   cmOperation = 655  // 000_101_0001111
	opSPST16  cmOperation = 1679 // 001_101_0001111
	opSPST32  cmOperation = 2703 // 010_101_0001111
	opSPST64  cmOperation = 3727 // 011_101_0001111
	opLD8     cmOperation = 783  // 000_110_0001111
	opLD16    cmOperation = 1807 // 001_110_0001111
	opLD32    cmOperation = 2831 // 010_110_0001111
	opLD64    cmOperation = 3855 // 011_110_0001111
	opLD8u    cmOperation = 4879 // 100_110_0001111
	opLD16u   cmOperation = 5903 // 101_110_0001111
	opLD32u   cmOperation = 6927 // 110_110_0001111
)

type instTypeCM struct {
	operation cmOperation // 13 bit
	imm       uint32      // 12 bit (Imm, CSR, 0 or 1)
	srcReg    uint32      // 7 bit
}

var strToCmOperation = map[string]cmOperation{
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
}

func (i *instTypeCM) toUInt32() uint32 {
	return uint32(i.operation) | (i.imm << 13) | (i.srcReg << 25)
}

func fromStringToInstTypeCM(str string) (*instTypeCM, error) {
	ss := strings.Fields(str)
	i := instTypeCM{}

	op, ok := strToCmOperation[ss[0]]
	if !ok {
		return nil, fmt.Errorf("not found '%s' in strToCmOperation(%v) :'%s'", ss[0], strToCmOperation, str)
	}
	i.operation = op

	switch op {
	case opRPINC:
		if ss[0] == "RPINC" {
			t, err := strconv.ParseUint(ss[1], 10, 7)
			if err != nil {
				return nil, err
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
				return nil, err
			}
			pred, err := strconv.ParseUint(ss[2], 10, 4)
			if err != nil {
				return nil, err
			}
			i.imm = uint32(pred | (succ << 4))
		}

	case opFENCEI:
		break

	case opECALL:
		if ss[0] == "EBREAK" {
			i.imm = 1
		} // else -> ECALL (i.imm = 1)

	case opJR, opJALR,
		opCSRRW, opCSRRS, opCSRRC, opCSRRWi, opCSRRSi, opCSRRCi,
		opLD8, opLD16, opLD32, opLD64, opLD8u, opLD16u, opLD32u:
		{
			if len(ss) < 3 {
				return nil, fmt.Errorf("too few arg: %s", str)
			}
			imm, err := strconv.ParseUint(ss[1], 10, 12) // Imm or CSR
			if err != nil {
				return nil, err
			}
			i.imm = uint32(imm)
			srcReg1, err := strconv.ParseUint(ss[2], 10, 7) // srcReg1 or zImm
			if err != nil {
				return nil, err
			}
			i.srcReg = uint32(srcReg1)
		}

	case opSPLD8, opSPLD16, opSPLD32, opSPLD64,
		opSPLD8u, opSPLD16u, opSPLD32u:
		{
			if len(ss) < 2 {
				return nil, fmt.Errorf("too few arg: %s", str)
			}
			imm, err := strconv.ParseUint(ss[1], 10, 12) // Imm
			if err != nil {
				return nil, err
			}
			i.imm = uint32(imm)
		}
	default:
		panic("can't reach here...exhaustive switch in fromStringToInstTypeCM")
	}
	return &i, nil
}

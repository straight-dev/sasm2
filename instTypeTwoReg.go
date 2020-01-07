package main

import (
	"fmt"
	"strconv"
	"strings"
)

type twoRegOperation uint32 // 18 (= 7 + 1 + 3 + 2 + 5) bit OPCODE

const (
	opADD32    twoRegOperation = 6223  // 00000_11_000_0_1001111
	opSUB32    twoRegOperation = 71759 // 01000_11_000_0_1001111
	opSLL32    twoRegOperation = 6479  // 00000_11_001_0_1001111
	opSLT32    twoRegOperation = 6735  // 00000_11_010_0_1001111
	opSLTu32   twoRegOperation = 6991  // 00000_11_011_0_1001111
	opXOR32    twoRegOperation = 7247  // 00000_11_100_0_1001111
	opSRL32    twoRegOperation = 7503  // 00000_11_101_0_1001111
	opSRA32    twoRegOperation = 73039 // 01000_11_101_0_1001111
	opOR32     twoRegOperation = 7759  // 00000_11_110_0_1001111
	opAND32    twoRegOperation = 8015  // 00000_11_111_0_1001111
	opADD64    twoRegOperation = 6351  // 00000_11_000_1_1001111
	opSUB64    twoRegOperation = 71887 // 01000_11_000_1_1001111
	opSLL64    twoRegOperation = 6607  // 00000_11_001_1_1001111
	opSLT64    twoRegOperation = 6863  // 00000_11_010_1_1001111
	opSLTu64   twoRegOperation = 7119  // 00000_11_011_1_1001111
	opXOR64    twoRegOperation = 7375  // 00000_11_100_1_1001111
	opSRL64    twoRegOperation = 7631  // 00000_11_101_1_1001111
	opSRA64    twoRegOperation = 73167 // 01000_11_101_1_1001111
	opOR64     twoRegOperation = 7887  // 00000_11_110_1_1001111
	opAND64    twoRegOperation = 8143  // 00000_11_111_1_1001111
	opMUL32    twoRegOperation = 14415 // 00001_11_000_0_1001111
	opMULH32   twoRegOperation = 14671 // 00001_11_001_0_1001111
	opMULHsu32 twoRegOperation = 14927 // 00001_11_010_0_1001111
	opMULHu32  twoRegOperation = 15183 // 00001_11_011_0_1001111
	opDIV32    twoRegOperation = 15439 // 00001_11_100_0_1001111
	opDIVu32   twoRegOperation = 15695 // 00001_11_101_0_1001111
	opREM32    twoRegOperation = 15951 // 00001_11_110_0_1001111
	opREMu32   twoRegOperation = 16207 // 00001_11_111_0_1001111
	opMUL64    twoRegOperation = 14543 // 00001_11_000_1_1001111
	opMULH64   twoRegOperation = 14799 // 00001_11_001_1_1001111
	opMULHsu64 twoRegOperation = 15055 // 00001_11_010_1_1001111
	opMULHu64  twoRegOperation = 15311 // 00001_11_011_1_1001111
	opDIV64    twoRegOperation = 15567 // 00001_11_100_1_1001111
	opDIVu64   twoRegOperation = 15823 // 00001_11_101_1_1001111
	opREM64    twoRegOperation = 16079 // 00001_11_110_1_1001111
	opREMu64   twoRegOperation = 16335 // 00001_11_111_1_1001111
)

type instTypeTwoReg struct {
	operation twoRegOperation // 18 bit
	srcRegs   [2]uint32       // 7 bit x 2
}

var strToTwoRegOperation = map[string]twoRegOperation{
	"ADD.32":    opADD32,
	"SUB.32":    opSUB32,
	"SLL.32":    opSLL32,
	"SLT.32":    opSLT32,
	"SLTu.32":   opSLTu32,
	"XOR.32":    opXOR32,
	"SRL.32":    opSRL32,
	"SRA.32":    opSRA32,
	"OR.32":     opOR32,
	"AND.32":    opAND32,
	"ADD.64":    opADD64,
	"SUB.64":    opSUB64,
	"SLL.64":    opSLL64,
	"SLT.64":    opSLT64,
	"SLTu.64":   opSLTu64,
	"XOR.64":    opXOR64,
	"SRL.64":    opSRL64,
	"SRA.64":    opSRA64,
	"OR.64":     opOR64,
	"AND.64":    opAND64,
	"MUL.32":    opMUL32,
	"MULH.32":   opMULH32,
	"MULHsu.32": opMULHsu32,
	"MULHu.32":  opMULHu32,
	"DIV.32":    opDIV32,
	"DIVu.32":   opDIVu32,
	"REM.32":    opREM32,
	"REMu.32":   opREMu32,
	"MUL.64":    opMUL64,
	"MULH.64":   opMULH64,
	"MULHsu.64": opMULHsu64,
	"MULHu.64":  opMULHu64,
	"DIV.64":    opDIV64,
	"DIVu.64":   opDIVu64,
	"REM.64":    opREM64,
	"REMu.64":   opREMu64,
}

func (i *instTypeTwoReg) toUInt32() uint32 {
	return uint32(i.operation) | i.srcRegs[0]<<25 | i.srcRegs[1]<<18
}

func fromStringToInstTypeTwoReg(str string) (*instTypeTwoReg, error) {
	ss := strings.Fields(str)
	i := instTypeTwoReg{}

	if len(ss) < 3 {
		return nil, fmt.Errorf("too few arg: %s", str)
	}

	op, ok := strToTwoRegOperation[ss[0]]
	if !ok {
		return nil, fmt.Errorf("not found '%s' in strToTwoRegOperation(%v) :'%s'", ss[0], strToTwoRegOperation, str)
	}
	i.operation = op

	for j := 0; j < 2; j++ {
		srcReg, err := strconv.ParseUint(ss[j+1], 10, 7) // srcReg1 or zImm
		if err != nil {
			return nil, err
		}
		i.srcRegs[j] = uint32(srcReg)
	}

	return &i, nil
}

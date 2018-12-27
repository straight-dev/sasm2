package main

import (
	"fmt"
	"strconv"
	"strings"
)

type sbOperation uint32

const (
	opST8  sbOperation = 7
	opST16 sbOperation = 39
	opST32 sbOperation = 23
	opST64 sbOperation = 55
	opBLT  sbOperation = 3
	opBGE  sbOperation = 35
	opBLTU sbOperation = 19
	opBGEU sbOperation = 51
	opBEQ  sbOperation = 11
	opBNE  sbOperation = 43
)

type instTypeSB struct {
	operation sbOperation
	imm12     uint32
	srcRegs   [2]uint32
}

var strToSBOperation = map[string]sbOperation{
	"ST.8":  opST8,
	"ST.16": opST16,
	"ST.32": opST32,
	"ST.64": opST64,
	"BLT":   opBLT,
	"BGE":   opBGE,
	"BLTU":  opBLTU,
	"BGEU":  opBGEU,
	"BEQ":   opBEQ,
	"BNE":   opBNE,
}

func (i *instTypeSB) toUInt32() uint32 {
	return uint32(i.operation) | (i.imm12 << 6) | (i.srcRegs[1] << 18) | (i.srcRegs[0] << 25)
}

func fromStringToInstTypeSB(str string) (*instTypeSB, error) {
	ss := strings.Fields(str)
	i := instTypeSB{}

	op, ok := strToSBOperation[ss[0]]
	if !ok {
		println(strToSBOperation)
		return nil, fmt.Errorf("not found '%s' in strToSBOperation :'%s'", ss[0], str)
	}
	i.operation = op

	if len(ss) < 4 {
		return nil, fmt.Errorf("invalid inst : few args '%s'", str)
	}

	for j := 0; j < 2; j++ {
		t, err := strconv.ParseUint(ss[j+1], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("failed to ParseUint '%s' in %s: %s", ss[j+1], str, err)
		}
		i.srcRegs[j] = uint32(t)
	}

	t, err := strconv.ParseInt(ss[3], 10, 12)
	if err != nil {
		return nil, fmt.Errorf("failed to ParseUint '%s' in %s: %s", ss[3], str, err)
	}
	i.imm12 = uint32(t) & 0xfff

	return &i, nil
}

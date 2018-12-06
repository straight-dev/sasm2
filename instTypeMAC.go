package main

import (
	"fmt"
	"strconv"
	"strings"
)

type macOperation uint32 // 8 bit OPCODE

const (
	opFMADDs  macOperation = 0x1b
	opFMADDd  macOperation = 0x9b
	opFMSUBs  macOperation = 0x5b
	opFMSUBd  macOperation = 0xdb
	opFNMSUBs macOperation = 0x3b
	opFNMSUBd macOperation = 0xbb
	opFNMADDs macOperation = 0x7b
	opFNMADDd macOperation = 0xfb
)

type instTypeMAC struct {
	operation macOperation // 8 bit
	rm        roundmode    // 3 bit
	srcRegs   [3]uint32    // 7 bit x3
}

var strToMacOperation = map[string]macOperation{
	"FMADD.s":  opFMADDs,
	"FMADD.d":  opFMADDd,
	"FMSUB.s":  opFMSUBs,
	"FMSUB.d":  opFMSUBd,
	"FNMSUB.s": opFNMSUBs,
	"FNMADD.s": opFNMADDs,
	"FNMADD.d": opFNMADDd,
}

func (i *instTypeMAC) toUInt32() uint32 {
	return uint32(i.operation) | (uint32(i.rm) << 8) | (i.srcRegs[2] << 11) | (i.srcRegs[1] << 18) | (i.srcRegs[0] << 25)
}

// fromStringToInstTypeMAC
// example: "FMADD.s 1 2 3 RDN" (RM is optional)
func fromStringToInstTypeMAC(str string) (*instTypeMAC, error) {
	ss := strings.Fields(str)
	i := instTypeMAC{}

	op, ok := strToMacOperation[ss[0]]
	if !ok {
		println(strToMacOperation)
		return nil, fmt.Errorf("not found '%s' in strToMacOperation :'%s'", ss[0], str)
	}
	i.operation = op

	if len(ss) < 4 {
		return nil, fmt.Errorf("invalid inst : few args '%s'", str)
	}

	for j := 0; j < 3; j++ {
		t, err := strconv.ParseUint(ss[j+1], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("failed to ParseUint '%s' in %s: %s", ss[j+1], str, err)
		}
		i.srcRegs[j] = uint32(t)
	}

	if len(ss) >= 5 {
		rm, err := fromStringToRM(ss[4])
		if err != nil {
			return nil, fmt.Errorf("failed to fromStringToRM '%s' in %s: %s", ss[4], str, err)
		}
		i.rm = rm
	}

	return &i, nil
}

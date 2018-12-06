package main

import (
	"fmt"
	"strconv"
	"strings"
)

type noRegOperation uint32

const (
	opJ      noRegOperation = 399  // 00_011_0001111
	opJAL    noRegOperation = 1423 // 01_011_0001111
	opLUi    noRegOperation = 2447 // 10_011_0001111
	opAUiPC  noRegOperation = 3471 // 11_011_0001111
	opSPADDi noRegOperation = 911  // 00_111_0001111
	opAUiSP  noRegOperation = 3983 // 11_111_0001111
)

type instTypeNoReg struct {
	operation noRegOperation
	imm20     uint32
}

var strToNoRegOperation = map[string]noRegOperation{
	"J":      opJ,
	"JAL":    opJAL,
	"LUi":    opLUi,
	"AUiPC":  opAUiPC,
	"SPADDi": opSPADDi,
	"AUiSP":  opAUiSP,
}

func (i *instTypeNoReg) toUInt32() uint32 {
	return uint32(i.operation) | (i.imm20 << 12)
}

func fromStringToInstTypeNoReg(str string) (*instTypeNoReg, error) {
	ss := strings.Fields(str)
	i := instTypeNoReg{}

	op, ok := strToNoRegOperation[ss[0]]
	if !ok {
		println(strToNoRegOperation)
		return nil, fmt.Errorf("not found '%s' in strToNoRegOperation :'%s'", ss[0], str)
	}
	i.operation = op

	t, err := strconv.ParseInt(ss[1], 10, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to ParseUint '%s' in %s: %s", ss[1], str, err)
	}
	i.imm20 = uint32(t) & 0xfffff

	return &i, nil
}

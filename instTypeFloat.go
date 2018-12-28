package main

import (
	"fmt"
	"strconv"
	"strings"
)

type floatOperation uint32 // 8 bit OPCODE

const (
	opFADDs    floatOperation = 0x4f    // 00000_00_000_0_1001111
	opFSUBs    floatOperation = 0x204f  // 00001_00_000_0_1001111
	opFMULs    floatOperation = 0x404f  // 00010_00_000_0_1001111
	opFDIVs    floatOperation = 0x604f  // 00011_00_000_0_1001111
	opFSQRTs   floatOperation = 0x804f  // 00100_00_000_0_1001111
	opFSGNJs   floatOperation = 0x1004f // 01000_00_000_0_1001111
	opFSGNJNs  floatOperation = 0x1204f // 01001_00_000_0_1001111
	opFSGNJXs  floatOperation = 0x1404f // 01010_00_000_0_1001111
	opFMINs    floatOperation = 0x1804f // 01100_00_000_0_1001111
	opFMAXs    floatOperation = 0x1a04f // 01101_00_000_0_1001111
	opFCLASSs  floatOperation = 0x2004f // 10000_00_000_0_1001111
	opFEQs     floatOperation = 0x2204f // 10001_00_000_0_1001111
	opFLTs     floatOperation = 0x2404f // 10010_00_000_0_1001111
	opFLEs     floatOperation = 0x2604f // 10011_00_000_0_1001111
	opFCVTsd   floatOperation = 0x2804f // 10100_00_000_0_1001111
	opFCVT32s  floatOperation = 0x3004f // 11000_00_000_0_1001111
	opFCVT32us floatOperation = 0x3204f // 11001_00_000_0_1001111
	opFCVTs32  floatOperation = 0x3404f // 11010_00_000_0_1001111
	opFCVTs32u floatOperation = 0x3604f // 11011_00_000_0_1001111
	opFCVT64s  floatOperation = 0x3804f // 11100_00_000_0_1001111
	opFCVT64us floatOperation = 0x3a04f // 11101_00_000_0_1001111
	opFCVTs64  floatOperation = 0x3c04f // 11110_00_000_0_1001111
	opFCVTs64u floatOperation = 0x3e04f // 11111_00_000_0_1001111
	opFADDd    floatOperation = 0xcf    // 00000_00_000_1_1001111
	opFSUBd    floatOperation = 0x20cf  // 00001_00_000_1_1001111
	opFMULd    floatOperation = 0x40cf  // 00010_00_000_1_1001111
	opFDIVd    floatOperation = 0x60cf  // 00011_00_000_1_1001111
	opFSQRTd   floatOperation = 0x80cf  // 00100_00_000_1_1001111
	opFSGNJd   floatOperation = 0x100cf // 01000_00_000_1_1001111
	opFSGNJNd  floatOperation = 0x120cf // 01001_00_000_1_1001111
	opFSGNJXd  floatOperation = 0x140cf // 01010_00_000_1_1001111
	opFMINd    floatOperation = 0x180cf // 01100_00_000_1_1001111
	opFMAXd    floatOperation = 0x1a0cf // 01101_00_000_1_1001111
	opFCLASSd  floatOperation = 0x200cf // 10000_00_000_1_1001111
	opFEQd     floatOperation = 0x220cf // 10001_00_000_1_1001111
	opFLTd     floatOperation = 0x240cf // 10010_00_000_1_1001111
	opFLEd     floatOperation = 0x260cf // 10011_00_000_1_1001111
	opFCVTds   floatOperation = 0x280cf // 10100_00_000_1_1001111
	opFCVT32d  floatOperation = 0x300cf // 11000_00_000_1_1001111
	opFCVT32ud floatOperation = 0x320cf // 11001_00_000_1_1001111
	opFCVTd32  floatOperation = 0x340cf // 11010_00_000_1_1001111
	opFCVTd32u floatOperation = 0x360cf // 11011_00_000_1_1001111
	opFCVT64d  floatOperation = 0x380cf // 11100_00_000_1_1001111
	opFCVT64ud floatOperation = 0x3a0cf // 11101_00_000_1_1001111
	opFCVTd64  floatOperation = 0x3c0cf // 11110_00_000_1_1001111
	opFCVTd64u floatOperation = 0x3e0cf // 11111_00_000_1_1001111

)

type instTypeFloat struct {
	operation floatOperation // 5+2+(+3)+1+7 = 15(18) bit
	rm        roundmode      // 3 bit
	srcRegs   [2]uint32      // 7 bit x2
}

var strToFloatOperation = map[string]floatOperation{
	"FADD.s":     opFADDs,
	"FSUB.s":     opFSUBs,
	"FMUL.s":     opFMULs,
	"FDIV.s":     opFDIVs,
	"FSQRT.s":    opFSQRTs,
	"FSGNJ.s":    opFSGNJs,
	"FSGNJN.s":   opFSGNJNs,
	"FSGNJX.s":   opFSGNJXs,
	"FMIN.s":     opFMINs,
	"FMAX.s":     opFMAXs,
	"FCLASS.s":   opFCLASSs,
	"FEQ.s":      opFEQs,
	"FLT.s":      opFLTs,
	"FLE.s":      opFLEs,
	"FCVT.s.d":   opFCVTsd,
	"FCVT.32.s":  opFCVT32s,
	"FCVT.32u.s": opFCVT32us,
	"FCVT.s.32":  opFCVTs32,
	"FCVT.s.32u": opFCVTs32u,
	"FCVT.64.s":  opFCVT64s,
	"FCVT.64u.s": opFCVT64us,
	"FCVT.s.64":  opFCVTs64,
	"FCVT.s.64u": opFCVTs64u,
	"FADD.d":     opFADDd,
	"FSUB.d":     opFSUBd,
	"FMUL.d":     opFMULd,
	"FDIV.d":     opFDIVd,
	"FSQRT.d":    opFSQRTd,
	"FSGNJ.d":    opFSGNJd,
	"FSGNJN.d":   opFSGNJNd,
	"FSGNJX.d":   opFSGNJXd,
	"FMIN.d":     opFMINd,
	"FMAX.d":     opFMAXd,
	"FCLASS.d":   opFCLASSd,
	"FEQ.d":      opFEQd,
	"FLT.d":      opFLTd,
	"FLE.d":      opFLEd,
	"FCVT.d.s":   opFCVTds,
	"FCVT.32.d":  opFCVT32d,
	"FCVT.32u.d": opFCVT32ud,
	"FCVT.d.32":  opFCVTd32,
	"FCVT.d.32u": opFCVTd32u,
	"FCVT.64.d":  opFCVT64d,
	"FCVT.64u.d": opFCVT64ud,
	"FCVT.d.64":  opFCVTd64,
	"FCVT.d.64u": opFCVTd64u,
}

func (i *instTypeFloat) toUInt32() uint32 {
	return uint32(i.operation) | (uint32(i.rm << 8)) | (i.srcRegs[1] << 18) | (i.srcRegs[0] << 25)
}

// fromStringToInstTypeFloat
// example: "FMADD.s 1 2 RDN" (RM is optional)
func fromStringToInstTypeFloat(str string) (*instTypeFloat, error) {
	ss := strings.Fields(str)
	i := instTypeFloat{}

	op, ok := strToFloatOperation[ss[0]]
	if !ok {
		println(strToFloatOperation)
		return nil, fmt.Errorf("not found '%s' in strToFloatOperation :'%s'", ss[0], str)
	}
	i.operation = op

	if len(ss) < 3 {
		return nil, fmt.Errorf("invalid inst : few args '%s'", str)
	}

	for j := 0; j < 2; j++ {
		t, err := strconv.ParseUint(ss[j+1], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("failed to ParseUint '%s' in %s: %s", ss[j+1], str, err)
		}
		i.srcRegs[j] = uint32(t)
	}

	if len(ss) >= 4 {
		rm, err := fromStringToRM(ss[3])
		if err != nil {
			return nil, fmt.Errorf("failed to fromStringToRM '%s' in %s: %s", ss[4], str, err)
		}
		i.rm = rm
	}

	return &i, nil
}

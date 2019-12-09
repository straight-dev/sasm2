package main

import (
	"fmt"
	"strconv"
	"strings"
)

type floatOperation uint32 // 8 bit OPCODE

const (
	opFADD32       floatOperation = 0x4f    // 00000_00_000_0_1001111
	opFSUB32       floatOperation = 0x204f  // 00001_00_000_0_1001111
	opFMUL32       floatOperation = 0x404f  // 00010_00_000_0_1001111
	opFDIV32       floatOperation = 0x604f  // 00011_00_000_0_1001111
	opFSQRTs       floatOperation = 0x804f  // 00100_00_000_0_1001111
	opFSGNJs       floatOperation = 0x1004f // 01000_00_000_0_1001111
	opFSGNJNs      floatOperation = 0x1204f // 01001_00_000_0_1001111
	opFSGNJXs      floatOperation = 0x1404f // 01010_00_000_0_1001111
	opFMINs        floatOperation = 0x1804f // 01100_00_000_0_1001111
	opFMAXs        floatOperation = 0x1a04f // 01101_00_000_0_1001111
	opFCLASSs      floatOperation = 0x2004f // 10000_00_000_0_1001111
	opFEQs         floatOperation = 0x2204f // 10001_00_000_0_1001111
	opFLTs         floatOperation = 0x2404f // 10010_00_000_0_1001111
	opFLEs         floatOperation = 0x2604f // 10011_00_000_0_1001111
	opFCVTf64tof32 floatOperation = 0x2804f // 10100_00_000_0_1001111
	opFCVTf32toi32 floatOperation = 0x3004f // 11000_00_000_0_1001111
	opFCVT32us     floatOperation = 0x3204f // 11001_00_000_0_1001111
	opFCVTs32      floatOperation = 0x3404f // 11010_00_000_0_1001111
	opFCVTs32u     floatOperation = 0x3604f // 11011_00_000_0_1001111
	opFCVT64s      floatOperation = 0x3804f // 11100_00_000_0_1001111
	opFCVT64us     floatOperation = 0x3a04f // 11101_00_000_0_1001111
	opFCVTs64      floatOperation = 0x3c04f // 11110_00_000_0_1001111
	opFCVTs64u     floatOperation = 0x3e04f // 11111_00_000_0_1001111
	opFADDd        floatOperation = 0xcf    // 00000_00_000_1_1001111
	opFSUBd        floatOperation = 0x20cf  // 00001_00_000_1_1001111
	opFMULd        floatOperation = 0x40cf  // 00010_00_000_1_1001111
	opFDIVd        floatOperation = 0x60cf  // 00011_00_000_1_1001111
	opFSQRTd       floatOperation = 0x80cf  // 00100_00_000_1_1001111
	opFSGNJd       floatOperation = 0x100cf // 01000_00_000_1_1001111
	opFSGNJNd      floatOperation = 0x120cf // 01001_00_000_1_1001111
	opFSGNJXd      floatOperation = 0x140cf // 01010_00_000_1_1001111
	opFMINd        floatOperation = 0x180cf // 01100_00_000_1_1001111
	opFMAXd        floatOperation = 0x1a0cf // 01101_00_000_1_1001111
	opFCLASSd      floatOperation = 0x200cf // 10000_00_000_1_1001111
	opFEQd         floatOperation = 0x220cf // 10001_00_000_1_1001111
	opFLTd         floatOperation = 0x240cf // 10010_00_000_1_1001111
	opFLEd         floatOperation = 0x260cf // 10011_00_000_1_1001111
	opFCVTds       floatOperation = 0x280cf // 10100_00_000_1_1001111
	opFCVT32d      floatOperation = 0x300cf // 11000_00_000_1_1001111
	opFCVT32ud     floatOperation = 0x320cf // 11001_00_000_1_1001111
	opFCVTd32      floatOperation = 0x340cf // 11010_00_000_1_1001111
	opFCVTd32u     floatOperation = 0x360cf // 11011_00_000_1_1001111
	opFCVT64d      floatOperation = 0x380cf // 11100_00_000_1_1001111
	opFCVT64ud     floatOperation = 0x3a0cf // 11101_00_000_1_1001111
	opFCVTd64      floatOperation = 0x3c0cf // 11110_00_000_1_1001111
	opFCVTd64u     floatOperation = 0x3e0cf // 11111_00_000_1_1001111

)

type instTypeFloat struct {
	operation floatOperation // 5+2+(+3)+1+7 = 15(18) bit
	rm        roundmode      // 3 bit
	srcRegs   [2]uint32      // 7 bit x2
}

var strToFloatOperation = map[string]floatOperation{
	"FADD.32":         opFADD32,
	"FSUB.32":         opFSUB32,
	"FMUL.32":         opFMUL32,
	"FDIV.32":         opFDIV32,
	"FSQRT.32":        opFSQRTs,
	"FSGNJ.32":        opFSGNJs,
	"FSGNJN.32":       opFSGNJNs,
	"FSGNJX.32":       opFSGNJXs,
	"FMIN.32":         opFMINs,
	"FMAX.32":         opFMAXs,
	"FCLASS.32":       opFCLASSs,
	"FEQ.32":          opFEQs,
	"FLT.32":          opFLTs,
	"FLE.32":          opFLEs,
	"FCVT.f64.to.f32": opFCVTf64tof32,
	"FCVT.32.s":       opFCVTf32toi32,
	"FCVT.f32.to.s32": opFCVTf32toi32,
	"FCVT.32u.s":      opFCVT32us,
	"FCVT.f32.to.u32": opFCVT32us,
	"FCVT.s.32":       opFCVTs32,
	"FCVT.s32.to.f32": opFCVTs32,
	"FCVT.s.32u":      opFCVTs32u,
	"FCVT.u32.to.f32": opFCVTs32u,
	"FCVT.64.s":       opFCVT64s,
	"FCVT.f32.to.s64": opFCVT64s,
	"FCVT.64u.s":      opFCVT64us,
	"FCVT.f32.to.u64": opFCVT64us,
	"FCVT.s.64":       opFCVTs64,
	"FCVT.s64.to.f32": opFCVTs64,
	"FCVT.s.64u":      opFCVTs64u,
	"FCVT.u64.to.f32": opFCVTs64u,
	"FADD.64":         opFADDd,
	"FSUB.64":         opFSUBd,
	"FMUL.64":         opFMULd,
	"FDIV.64":         opFDIVd,
	"FSQRT.64":        opFSQRTd,
	"FSGNJ.64":        opFSGNJd,
	"FSGNJN.64":       opFSGNJNd,
	"FSGNJX.64":       opFSGNJXd,
	"FMIN.64":         opFMINd,
	"FMAX.64":         opFMAXd,
	"FCLASS.64":       opFCLASSd,
	"FEQ.64":          opFEQd,
	"FLT.64":          opFLTd,
	"FLE.64":          opFLEd,
	"FCVT.d.s":        opFCVTds,
	"FCVT.f32.to.f64": opFCVTds,
	"FCVT.32.d":       opFCVT32d,
	"FCVT.f64.to.s32": opFCVT32d,
	"FCVT.32u.d":      opFCVT32ud,
	"FCVT.f64.to.u32": opFCVT32ud,
	"FCVT.d.32":       opFCVTd32,
	"FCVT.s32.to.f64": opFCVTd32,
	"FCVT.d.32u":      opFCVTd32u,
	"FCVT.u32.to.f64": opFCVTd32u,
	"FCVT.f64.to.s64": opFCVT64d,
	"FCVT.64u.d":      opFCVT64ud,
	"FCVT.f64.to.u64": opFCVT64ud,
	"FCVT.s64.to.f64": opFCVTd64,
	"FCVT.d.64u":      opFCVTd64u,
	"FCVT.u64.to.f64": opFCVTd64u,
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

	switch op {
	case opFADD32, opFSUB32, opFMUL32, opFDIV32, opFSQRTs, opFSGNJs, opFSGNJNs, opFSGNJXs, opFMINs, opFMAXs, opFCLASSs, opFEQs, opFLTs, opFLEs,
		opFADDd, opFSUBd, opFMULd, opFDIVd, opFSQRTd, opFSGNJd, opFSGNJNd, opFSGNJXd, opFMINd, opFMAXd, opFCLASSd, opFEQd, opFLTd, opFLEd:
		{
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

		}

	case opFCVTf64tof32, opFCVTf32toi32, opFCVT32us, opFCVTs32, opFCVTs32u, opFCVT64s, opFCVT64us, opFCVTs64, opFCVTs64u,
		opFCVTds, opFCVT32d, opFCVT32ud, opFCVTd32, opFCVTd32u, opFCVT64d, opFCVT64ud, opFCVTd64, opFCVTd64u:
		{

			if len(ss) < 2 {
				return nil, fmt.Errorf("invalid inst : few args '%s'", str)
			}

			t, err := strconv.ParseUint(ss[1], 10, 8)
			if err != nil {
				return nil, fmt.Errorf("failed to ParseUint '%s' in %s: %s", ss[1], str, err)
			}
			i.srcRegs[0] = uint32(t)
		}
	default:
		panic("can't reach here... instTypeFloat.go switch")
	}
	return &i, nil
}

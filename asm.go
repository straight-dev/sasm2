package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"os"
	"strconv"
	"strings"

	. "github.com/clkbug/sasm2/straightISAv1Info"
)

const dataStartAddr = 0x100

type instruction struct {
	opCode   OpCode
	instType InstType
	regs     []uint64
	imm      uint64
}

func (inst *instruction) instTobytes() []byte {
	var i uint32
	b := make([]byte, 4)
	i = (uint32(inst.opCode) & 0x3f) << 26
	imm := uint32(inst.imm)
	switch inst.instType {
	case ZeroReg:
		i += imm & 0x3ffffff
	case OneReg:
		i += (uint32(inst.regs[0]) & 0x3ff) << 16
		i += imm & 0xffff
	case TwoReg:
		i += (uint32(inst.regs[0]) & 0x3ff) << 16
		i += (uint32(inst.regs[1]) & 0x3ff) << 6
		i += imm & 0x3f
	}
	binary.BigEndian.PutUint32(b, i)
	return b
}

var strToOpCodeMap map[string]OpCode

func strToOpCode(s string) (OpCode, error) {
	// make the map of string2OpCode
	if len(strToOpCodeMap) == 0 {
		strToOpCodeMap = make(map[string]OpCode)
		for i := 0; i < int(Op_MAX); i++ {
			oc := OpCode(i)
			strToOpCodeMap[strings.ToLower(oc.String()[2:])] = oc
		}
		strToOpCodeMap["undef"] = OpNOP
	}

	oc, exists := strToOpCodeMap[s]
	if !exists {
		return Op_MAX, errors.New("invalid opname")
	}
	return oc, nil
}

func strToInst(s string) instruction {
	ss := strings.Fields(s)
	inst := instruction{}

	if oc, err := strToOpCode(strings.ToLower(ss[0])); err == nil {
		inst.opCode = oc
	} else {
		println(strToOpCodeMap)
		println(ss[0])
		panic("invalid op")
	}

	for i := 1; i < len(ss); i++ {
		if strings.HasPrefix(ss[i], "$CONST") {
			constIndex, _ := strconv.ParseUint(ss[i][6:], 10, 64)
			ss[i] = strconv.FormatUint(dataStartAddr+constIndex*8, 10)
		}
	}

	var isImmSignExtended bool
	inst.instType, isImmSignExtended, _ = GetInstType(inst.opCode)
	switch inst.instType {
	case ZeroReg:
		if 2 <= len(ss) {
			var imm uint64
			if isImmSignExtended {
				immb, _ := strconv.ParseInt(ss[1], 10, 26)
				imm = uint64(immb)
			} else {
				imm, _ = strconv.ParseUint(ss[1], 10, 26)
			}
			inst.imm = extractBits(imm, 26, isImmSignExtended)
		}
	case OneReg:
		reg, _ := strconv.ParseUint(ss[1], 10, 10)
		inst.regs = append(inst.regs, reg)
		if 3 <= len(ss) {
			var imm uint64
			if isImmSignExtended {
				immb, _ := strconv.ParseInt(ss[2], 10, 16)
				imm = uint64(immb)
			} else {
				imm, _ = strconv.ParseUint(ss[2], 10, 16)
			}
			inst.imm = extractBits(imm, 16, isImmSignExtended)
		}
	case TwoReg:
		reg, _ := strconv.ParseUint(ss[1], 10, 10)
		inst.regs = append(inst.regs, reg)
		reg, _ = strconv.ParseUint(ss[2], 10, 10)
		inst.regs = append(inst.regs, reg)
		if 4 <= len(ss) {
			var imm uint64
			if isImmSignExtended {
				immb, _ := strconv.ParseInt(ss[3], 10, 6)
				imm = uint64(immb)
			} else {
				imm, _ = strconv.ParseUint(ss[3], 10, 6)
			}
			inst.imm = extractBits(imm, 6, isImmSignExtended)
		}
	}
	return inst
}

func assemble(fileName, outputFileName string) error {
	fp, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	var insts []instruction
	var datum []uint64
	for scanner.Scan() {
		t := scanner.Text()
		if d, err := strconv.ParseUint(t, 16, 64); err == nil {
			datum = append(datum, d)
		} else {
			insts = append(insts, strToInst(t))
		}
	}
	insts = append(insts, strToInst("JR 0"))
	if len(insts)%2 == 1 {
		insts = append(insts, strToInst("NOP"))
	}

	elf := NewELFFile()
	prog := make([]byte, len(insts)*4)
	for i, v := range insts {
		t := v.instTobytes()
		copy(prog[4*i:4*i+3], t)
	}

	progHeader := ElfProgHeader{
		ProgType:     ProgTypeLoad,
		ProgFlags:    ProgFlagExecute + ProgFlagRead,
		ProgVAddr:    ProgEntryAddr,
		ProgPAddr:    0,
		ProgFileSize: uint64(len(insts) * 4),
		Prog:         prog,
	}

	secHeader := ElfSecHeader{
		SecType: SecTypeNull,
	}

	elf.AddSegment(&progHeader)
	elf.Sections = append(elf.Sections, &secHeader)
	elf.WriteELFFile(outputFileName)
	return nil
}

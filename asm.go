package main

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

const dataStartAddr = 0x10000
const initialSP = 0x0afffffc
const stackSize = 0x00500000

var entryOffset = 0

func strToInst(s string) (instruction, error) {
	ss := strings.Fields(s)

	if _, ok := strToSBOperation[ss[0]]; ok {
		i, err := fromStringToInstTypeSB(s)
		return instruction(i), err
	}

	if _, ok := strToMacOperation[ss[0]]; ok {
		i, err := fromStringToInstTypeMAC(s)
		return instruction(i), err
	}

	if _, ok := strToOneRegOperation[ss[0]]; ok {
		i, err := fromStringToInstTypeOneReg(s)
		return instruction(i), err
	}

	if _, ok := strToTwoRegOperation[ss[0]]; ok {
		i, err := fromStringToInstTypeTwoReg(s)
		return instruction(i), err
	}

	if _, ok := strToNoRegOperation[ss[0]]; ok {
		i, err := fromStringToInstTypeNoReg(s)
		return instruction(i), err
	}

	if _, ok := strToFloatOperation[ss[0]]; ok {
		i, err := fromStringToInstTypeFloat(s)
		return instruction(i), err
	}

	println(s, ss[0])
	panic("unimplemented yet")
}

func assemble(fileName, outputFileName string) error {
	fp, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	var insts []instruction
	var datum []byte
	for isInst := true; scanner.Scan(); {
		t := scanner.Text()
		if t == "Initialize values" {
			isInst = false
			continue
		} else if []rune(t)[0] == '!' {
			entryOffset = len(insts) * 4
			t = t[1:]
		}

		if isInst {
			i, err := strToInst(t)
			if err != nil {
				return err
			}
			insts = append(insts, i)
		} else {
			t := strings.TrimSpace(t)
			s := strings.Split(t, " ")
			for _, s := range s {
				if d, err := strconv.ParseUint(s, 10, 8); err == nil {
					datum = append(datum, byte(d))
				} else {
					return errors.New("invalid data\n" + err.Error())
				}
			}
		}
	}

	elf := NewELFFile()
	prog := make([]byte, len(insts)*4)
	for i, v := range insts {
		t := instToBytes(v)
		copy(prog[4*i:4*(i+1)], t[:])
	}
	datumbytes := make([]byte, len(datum)+dataStartAddr)
	for i, v := range datum {
		datumbytes[i+dataStartAddr] = v
	}

	progHeader := ElfProgHeader{
		ProgType:     ProgTypeLoad,
		ProgFlags:    ProgFlagExecute + ProgFlagRead,
		ProgVAddr:    ProgEntryAddr,
		ProgPAddr:    0,
		ProgFileSize: uint64(len(insts) * 4), // あとでlegalize
		Prog:         prog,
	}
	elf.AddSegment(&progHeader)

	stackHeader := ElfProgHeader{
		ProgType:     ProgTypeLoad,
		ProgFlags:    ProgFlagWrite + ProgFlagRead,
		ProgVAddr:    initialSP - stackSize,
		ProgPAddr:    0,
		ProgFileSize: 0,
		ProgMemSize:  stackSize,
		Prog:         nil,
	}
	elf.AddSegment(&stackHeader)

	globalDataHeader := ElfProgHeader{
		ProgType:     ProgTypeLoad,
		ProgFlags:    ProgFlagWrite + ProgFlagRead,
		ProgVAddr:    dataStartAddr - dataStartAddr,
		ProgPAddr:    0,
		ProgFileSize: uint64(len(datumbytes)),
		ProgMemSize:  33554432,
		Prog:         datumbytes,
	}
	elf.AddSegment(&globalDataHeader)
	if uint64(len(datumbytes)) > 33554432 {
		panic("too much global data")
	}

	secHeader := ElfSecHeader{
		SecType: SecTypeNull,
	}
	elf.Sections = append(elf.Sections, &secHeader)

	secStrTable := ElfSecHeader{
		SecType: SecTypeStrTab,
		Sec:     make([]byte, 20),
	}
	secStrTable.Sec[0] = 0x0
	copy(secStrTable.Sec[1:], "DummySectionHeader")
	elf.Sections = append(elf.Sections, &secStrTable)

	return elf.WriteELFFile(outputFileName)
}

package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"unsafe"
)

type ElfAddr uint64
type ElfOff uint64
type ElfSection uint16
type ElfVersym uint16

const (
	ElfIdentCLASS      = 4  /* Class of machine. */
	ElfIdentDATA       = 5  /* Data format. */
	ElfIdentVERSION    = 6  /* ELF format version. */
	ElfIdentOSABI      = 7  /* Operating system / ABI identification */
	ElfIdentABIVERSION = 8  /* ABI version */
	ElfIdentPAD        = 9  /* Start of padding (per SVR4 ABI). */
	ElfIdentNIDENT     = 16 /* Size of e_ident array. */
)

const (
	ElfIdentClassNone = 0
	ElfIdentClass32   = 1
	ElfIdentClass64   = 2
)

const (
	ElfIdentDataNone = 0
	ElfIdentData2LSB = 1 // Little Endian
	ElfIdentData2MSB = 2 // Big Endian
)

type ElfType uint16

const (
	ElfTypeNone ElfType = 0 /* Unknown type. */
	ElfTypeRel  ElfType = 1 /* Relocatable. */
	ElfTypeExec ElfType = 2 /* Executable. */
)

type ElfMachine uint16

// EM_STRAIGHT is the machine code of STRAIGHT; found in elf.Header.Machine
const ElfMachineSTRAIGHT ElfMachine = 256

type ElfVersion uint32

const (
	ElfVersionNone    ElfVersion = 0
	ElfVersionCurrent ElfVersion = 1
)

type ElfHeader struct {
	ElfIdent      [16]byte
	ElfType       ElfType
	ElfMachine    ElfMachine
	ElfVersion    ElfVersion
	ElfEntry      ElfAddr
	ElfPHOff      ElfOff
	ElfSHOff      ElfOff
	ElfFlags      uint32
	ElfEHSize     uint16
	ElfPHEntSize  uint16
	ElfPHEntNum   uint16
	ElfSHEntSize  uint16
	ElfSHEntNum   uint16
	ElfSHStrIndex uint16
}

type ProgType uint32

const (
	ProgTypeNull    ProgType = 0 // ignore this entry
	ProgTypeLoad    ProgType = 1 // Loadable
	ProgTypePHeader ProgType = 6 // ProgHeader
)
const (
	ProgFlagExecute = 1
	ProgFlagWrite   = 2
	ProgFlagRead    = 4
)

type ElfSegment []uint64
type ElfProgHeader struct {
	ProgType     ProgType
	ProgFlags    uint32
	ProgOffset   ElfAddr
	ProgVAddr    ElfAddr
	ProgPAddr    ElfAddr
	ProgFileSize uint64
	ProgMemSize  uint64
	ProgAlign    uint64
	Prog         ElfSegment
}

type ElfFile struct {
	Header   ElfHeader
	Programs []*ElfProgHeader
}

func newELFHeader() ElfHeader {
	var ei [16]byte
	ei[0] = 0x7f
	ei[1] = 'E'
	ei[2] = 'L'
	ei[3] = 'F'
	ei[ElfIdentCLASS] = ElfIdentClass64
	ei[ElfIdentDATA] = ElfIdentData2MSB
	ei[ElfIdentVERSION] = byte(ElfVersionCurrent)

	eh := ElfHeader{
		ElfIdent:      ei,
		ElfType:       ElfTypeExec,
		ElfMachine:    ElfMachineSTRAIGHT,
		ElfVersion:    ElfVersionCurrent,
		ElfEntry:      0x40000,
		ElfPHOff:      0,
		ElfSHOff:      0,
		ElfFlags:      0,
		ElfEHSize:     0,
		ElfPHEntSize:  0,
		ElfPHEntNum:   0,
		ElfSHEntSize:  0,
		ElfSHEntNum:   0,
		ElfSHStrIndex: 0,
	}
	eh.ElfEHSize = uint16(unsafe.Sizeof(eh))
	return eh
}

func newELFFile() *ElfFile {
	ef := ElfFile{
		Header: newELFHeader(),
	}
	return &ef
}

func (elf *ElfFile) writeELFFile(fileName string) error {
	var bo binary.ByteOrder
	if elf.Header.ElfIdent[ElfIdentDATA] == ElfIdentData2LSB {
		bo = binary.LittleEndian
	} else {
		bo = binary.BigEndian
	}

	fp, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer fp.Close()

	err = elf.Header.writeELFHeader(fp, bo)
	if err != nil {
		return err
	}
	return nil

}

func (eh *ElfHeader) writeELFHeader(fp *os.File, bo binary.ByteOrder) error {
	ehb := make([]byte, eh.ElfEHSize)
	offset := 0
	for ; offset < ElfIdentNIDENT; offset++ {
		ehb[offset] = eh.ElfIdent[offset]
	}
	bo.PutUint16(ehb[offset:offset+2], uint16(eh.ElfType))
	offset += 2
	bo.PutUint16(ehb[offset:offset+2], uint16(eh.ElfMachine))
	offset += 2
	bo.PutUint32(ehb[offset:offset+4], uint32(eh.ElfVersion))
	offset += 4
	bo.PutUint64(ehb[offset:offset+8], uint64(eh.ElfEntry))
	offset += 8
	bo.PutUint64(ehb[offset:offset+8], uint64(eh.ElfPHOff))
	offset += 8
	bo.PutUint64(ehb[offset:offset+8], uint64(eh.ElfSHOff))
	offset += 8
	bo.PutUint32(ehb[offset:offset+4], eh.ElfFlags)
	offset += 4
	bo.PutUint16(ehb[offset:offset+2], eh.ElfEHSize)
	offset += 2
	bo.PutUint16(ehb[offset:offset+2], eh.ElfPHEntSize)
	offset += 2
	bo.PutUint16(ehb[offset:offset+2], eh.ElfPHEntNum)
	offset += 2
	bo.PutUint16(ehb[offset:offset+2], eh.ElfSHEntSize)
	offset += 2
	bo.PutUint16(ehb[offset:offset+2], eh.ElfSHEntNum)
	offset += 2
	bo.PutUint16(ehb[offset:offset+2], eh.ElfSHStrIndex)
	offset += 2

	n, err := fp.Write(ehb)
	if err != nil {
		return err
	}
	fmt.Printf("write %d bytes (header)\n", n)
	return nil
}

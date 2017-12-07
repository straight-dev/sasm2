package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

// ElfAddr : 64bit address
type ElfAddr uint64

// ElfOff : Offset = 64bit address
type ElfOff uint64

// ElfSection : Section index (16 bit)
type ElfSection uint16

// ElfVersym : Version Number of ELF
type ElfVersym uint16

// ProgEntryAddr : Entry address of this program
const ProgEntryAddr = 0x40000

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

// ElfMachineSTRAIGHT is the machine code of STRAIGHT; found in elf.Header.Machine
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

type ElfSegment []byte
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

const ElfProgHeaderSize = (32*2 + 64*6) / 8 // bytes

type ElfFile struct {
	Header   ElfHeader
	Programs []*ElfProgHeader
}

const ElfHeaderSize = 64

func NewELFHeader() ElfHeader {
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
		ElfEntry:      ProgEntryAddr,
		ElfPHOff:      ElfHeaderSize,
		ElfSHOff:      0,
		ElfFlags:      0,
		ElfEHSize:     ElfHeaderSize,
		ElfPHEntSize:  0,
		ElfPHEntNum:   0,
		ElfSHEntSize:  0,
		ElfSHEntNum:   0,
		ElfSHStrIndex: 0,
	}
	return eh
}

func NewELFFile() *ElfFile {
	ef := ElfFile{
		Header: NewELFHeader(),
	}
	return &ef
}

func (elf *ElfFile) WriteELFFile(fileName string) error {
	var bo binary.ByteOrder
	if elf.Header.ElfIdent[ElfIdentDATA] == ElfIdentData2LSB {
		bo = binary.LittleEndian
	} else {
		bo = binary.BigEndian
	}
	elf.Legalize()
	fp, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer fp.Close()

	err = elf.Header.WriteELFHeader(fp, bo)
	if err != nil {
		return err
	}

	for _, p := range elf.Programs {
		err = p.WriteELFProgHeader(fp, bo)
		if err != nil {
			return err
		}
	}

	for _, p := range elf.Programs {
		n, e := fp.Write(p.Prog)
		if e != nil {
			return e
		} else if n < len(p.Prog) {
			return errors.New("failed to write segments")
		}
	}
	return nil
}

func (eh *ElfHeader) WriteELFHeader(fp *os.File, bo binary.ByteOrder) error {
	var ehb bytes.Buffer
	binary.Write(&ehb, bo, eh)
	_, err := fp.Write(ehb.Bytes())
	return err
}

func (elf *ElfFile) LegalizeHeader() error {
	elf.Header.ElfPHEntNum = uint16(len(elf.Programs))
	elf.Header.ElfPHEntSize = ElfProgHeaderSize
	return nil
}

// PageSize : 4KB (onikiri)
const PageSize = 4096

func (elf *ElfFile) Legalize() error {
	elf.LegalizeHeader()

	var offset uint64 = ElfHeaderSize + ElfProgHeaderSize
	for i := 0; i < len(elf.Programs); i++ {
		elf.Programs[i].ProgFileSize = uint64(len(elf.Programs[i].Prog))
		elf.Programs[i].ProgOffset = ElfAddr(offset)
		offset += elf.Programs[i].ProgFileSize
		elf.Programs[i].ProgAlign = offset % PageSize
		if elf.Programs[i].ProgFlags&ProgFlagExecute == ProgFlagExecute {
			elf.Header.ElfEntry = ElfAddr(ProgEntryAddr + offset%PageSize)
		}
	}
	return nil
}

func (elf *ElfFile) AddSegment(ph *ElfProgHeader) error {
	elf.Programs = append(elf.Programs, ph)
	elf.LegalizeHeader()
	return nil
}

func (ph *ElfProgHeader) WriteELFProgHeader(fp *os.File, bo binary.ByteOrder) error {
	phb := make([]byte, ElfProgHeaderSize)
	offset := 0
	bo.PutUint32(phb[offset:offset+4], uint32(ph.ProgType)) // TODO: use binary.Write
	offset += 4
	bo.PutUint32(phb[offset:offset+4], ph.ProgFlags)
	offset += 4
	bo.PutUint64(phb[offset:offset+8], uint64(ph.ProgOffset))
	offset += 8
	bo.PutUint64(phb[offset:offset+8], uint64(ph.ProgVAddr))
	offset += 8
	bo.PutUint64(phb[offset:offset+8], uint64(ph.ProgPAddr))
	offset += 8
	bo.PutUint64(phb[offset:offset+8], ph.ProgFileSize)
	offset += 8
	bo.PutUint64(phb[offset:offset+8], ph.ProgMemSize)
	offset += 8
	bo.PutUint64(phb[offset:offset+8], ph.ProgAlign)
	offset += 8

	_, err := fp.Write(phb)
	return err
}

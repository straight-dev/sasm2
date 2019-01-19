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

// ElfVersym : Version Number of ELF
type ElfVersym uint16

// ProgEntryAddr : Entry address of this program (precise entry address is ProgEntryAddr + the ELF Header and Segment Headers' offset)
const ProgEntryAddr = 0x020000000

// e_ident
const (
	ElfIdentCLASS      = 4  /* Class of machine. */
	ElfIdentDATA       = 5  /* Data format. */
	ElfIdentVERSION    = 6  /* ELF format version. */
	ElfIdentOSABI      = 7  /* Operating system / ABI identification */
	ElfIdentABIVERSION = 8  /* ABI version */
	ElfIdentPAD        = 9  /* Start of padding (per SVR4 ABI). */
	ElfIdentNIDENT     = 16 /* Size of e_ident array. */
)

// EI_CLASS
const (
	ElfIdentClassNone = 0
	ElfIdentClass32   = 1
	ElfIdentClass64   = 2
)

// EI_DATA
const (
	ElfIdentDataNone = 0
	ElfIdentData2LSB = 1 // Little Endian
	ElfIdentData2MSB = 2 // Big Endian
)

// ElfType : object file type
type ElfType uint16

// e_type
const (
	ElfTypeNone ElfType = 0 /* Unknown type. */
	ElfTypeRel  ElfType = 1 /* Relocatable. */
	ElfTypeExec ElfType = 2 /* Executable. */
)

// ElfMachine : Architecture
type ElfMachine uint16

// ElfMachineSTRAIGHT is the machine code of STRAIGHT; found in elf.Header.Machine
const ElfMachineSTRAIGHT ElfMachine = 256

// ElfVersion : file version
type ElfVersion uint32

// e_version
const (
	ElfVersionNone    ElfVersion = 0
	ElfVersionCurrent ElfVersion = 1
)

// ElfHeader : Header of ELF
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

// ProgType :
type ProgType uint32

// p_type
const (
	ProgTypeNull    ProgType = 0 // ignore this entry
	ProgTypeLoad    ProgType = 1 // Loadable
	ProgTypePHeader ProgType = 6 // ProgHeader
)

// ProgFlag : bitmask
const (
	ProgFlagExecute = 1
	ProgFlagWrite   = 2
	ProgFlagRead    = 4
)

// Elf Segment : segment
type ElfSegment []byte

// ElfProgHeader : segment header
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

// ElfProgHeaderSize = sizeof(ElfProgHeader)
const ElfProgHeaderSize = (32*2 + 64*6) / 8 // bytes

// SecType : Section type
type SecType uint32

// SecTypes
const (
	SecTypeNull SecType = iota
	SecTypeProgBits
	SecTypeSymTab
	SecTypeStrTab
	SecTypeRela
	SecTypeHash
	SecTypeSHLib
	SecTypeDynSym
	SecTypeLoProc SecType = 0x70000000
	SecTypeHiProc SecType = 0x7fffffff
	SecTypeLoUser SecType = 0x80000000
	SecTypeHiUser SecType = 0xf8ffffff
)

// Section Header bitmask
const (
	SHFlagWrite     = 0x1
	SHFlagAlloc     = 0x2
	SHFlagExecInstr = 0x4
	SHFlagMaskProc  = 0xF0000000
)

// ElfSection :
type ElfSection []byte

// ElfSecHeader : section header
type ElfSecHeader struct {
	SecName      uint32
	SecType      SecType
	SecFlags     uint64
	SecAddr      ElfAddr
	SecOffset    ElfOff
	SecSize      uint64
	SecLink      uint32
	SecInfo      uint32
	SecAddrAlign uint64
	SecEntSize   uint64
	Sec          ElfSection
}

// ElfSecHeaderSize = sizeof(ElfSecHeader)
const ElfSecHeaderSize = (32*4 + 64*6) / 8 // bytes

// ElfFile : ElfFile structure
type ElfFile struct {
	Header   ElfHeader
	Programs []*ElfProgHeader
	Sections []*ElfSecHeader
}

// ElfHeaderSize = sizeof(Header)
const ElfHeaderSize = 64

func NewELFHeader() ElfHeader {
	var ei [16]byte
	ei[0] = 0x7f
	ei[1] = 'E'
	ei[2] = 'L'
	ei[3] = 'F'
	ei[ElfIdentCLASS] = ElfIdentClass64
	ei[ElfIdentDATA] = ElfIdentData2LSB
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

	// write strtable
	_, err = fp.Write(elf.Sections[1].Sec)
	if err != nil {
		return err
	}
	for _, s := range elf.Sections {
		err = s.WriteELFSecHeader(fp, bo)
		if err != nil {
			return err
		}
	}
	return err
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
	elf.Header.ElfSHEntNum = uint16(len(elf.Sections))
	elf.Header.ElfSHEntSize = ElfSecHeaderSize
	return nil
}

// PageSize : 4KB (onikiri)
const PageSize = 4096

func (elf *ElfFile) Legalize() error {
	elf.LegalizeHeader()

	// Legalize Segment Header
	var offset uint64 = ElfHeaderSize + ElfProgHeaderSize*3

	// .text
	elf.Programs[0].ProgFileSize = uint64(len(elf.Programs[0].Prog)) + ElfHeaderSize + ElfProgHeaderSize*3 // .text includes ELF Header
	elf.Programs[0].ProgMemSize = elf.Programs[0].ProgFileSize
	elf.Programs[0].ProgOffset = 0
	elf.Programs[0].ProgAlign = 0
	elf.Header.ElfEntry += ElfAddr(offset + uint64(entryOffset))
	offset += uint64(len(elf.Programs[0].Prog))

	// .stack
	elf.Programs[1].ProgFileSize = 0
	elf.Programs[1].ProgMemSize = stackSize
	elf.Programs[1].ProgOffset = ElfAddr(offset)  // maybe useless info
	elf.Programs[1].ProgAlign = offset % PageSize // maybe useless info

	// .global
	elf.Programs[2].ProgFileSize = uint64(len(elf.Programs[2].Prog))
	elf.Programs[2].ProgOffset = ElfAddr(offset)
	elf.Programs[2].ProgAlign = offset % PageSize
	offset += elf.Programs[2].ProgFileSize

	// StrTable
	strTableOffset := offset
	elf.Header.ElfSHStrIndex = 1
	elf.Sections[0].SecName = 1
	elf.Sections[1].SecName = 1
	offset += uint64(len(elf.Sections[1].Sec))
	elf.Header.ElfSHOff = ElfOff(offset)
	// Legalize Section Header
	elf.Sections[0].SecSize = 0
	elf.Sections[0].SecType = SecTypeNull
	elf.Sections[0].SecOffset = ElfOff(offset)

	// Legalize SecHeader of StrTable
	elf.Sections[1].SecSize = uint64(len(elf.Sections[1].Sec))
	elf.Sections[1].SecOffset = ElfOff(strTableOffset)
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

func (sh *ElfSecHeader) WriteELFSecHeader(fp *os.File, bo binary.ByteOrder) error {
	shb := make([]byte, ElfSecHeaderSize)
	offset := 0
	bo.PutUint32(shb[offset:offset+4], sh.SecName)
	offset += 4
	bo.PutUint32(shb[offset:offset+4], uint32(sh.SecType))
	offset += 4
	bo.PutUint64(shb[offset:offset+8], sh.SecFlags)
	offset += 8
	bo.PutUint64(shb[offset:offset+8], uint64(sh.SecAddr))
	offset += 8
	bo.PutUint64(shb[offset:offset+8], uint64(sh.SecOffset))
	offset += 8
	bo.PutUint64(shb[offset:offset+8], sh.SecSize)
	offset += 8
	bo.PutUint32(shb[offset:offset+4], sh.SecLink)
	offset += 4
	bo.PutUint32(shb[offset:offset+4], sh.SecInfo)
	offset += 4
	bo.PutUint64(shb[offset:offset+8], sh.SecAddrAlign)
	offset += 8
	bo.PutUint64(shb[offset:offset+8], sh.SecEntSize)
	offset += 8

	_, err := fp.Write(shb)
	return err
}

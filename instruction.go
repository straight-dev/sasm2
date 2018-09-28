package main

import (
	"fmt"
	"strings"
)

// Instruction interface
type instruction interface {
	toUInt32() uint32
}

func instToBytes(inst instruction) [4]byte {
	word := inst.toUInt32()
	bs := [4]byte{}
	for i := 0; i < 4; i++ {
		bs[i] = byte((word >> uint(8*i)) & 0xff)
	}
	return bs
}

type roundmode uint32

const (
	rmRNE       roundmode = iota // Round to Nearest, ties to Even
	rmRTZ                        // Round towards Zero
	rmRDN                        // Round Down (towards -inf)
	rmRUP                        // Round Up (towards +inf)
	rmRMM                        // Round to Nearest, ties to Max Magnitude
	rmReserved1                  // Invalid
	rmReserved2                  // Invalid
	rmDynamic                    // Dynamic Rounding Mide
)

func fromStringToRM(s string) (roundmode, error) {
	s = strings.TrimSpace(s)
	switch s {
	case "RNE":
		return rmRNE, nil
	case "RTZ":
		return rmRTZ, nil
	case "RDN":
		return rmRDN, nil
	case "RUP":
		return rmRUP, nil
	case "Dynamic":
		return rmDynamic, nil
	default:
		return rmRNE, fmt.Errorf("invalid Rounding Mode: " + s)
	}
}

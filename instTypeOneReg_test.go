package main

import (
	"testing"
)

func TestFromStringToInstTypeOneReg(t *testing.T) {
	compare := func(s string, expected instTypeOneReg) {
		actual, err := fromStringToInstTypeOneReg(s)
		if err != nil {
			t.Error(err.Error())
		}

		if *actual != expected {
			t.Error(actual, expected)
		}

	}

	var table = []struct {
		in       string
		expected instTypeOneReg
	}{
		{
			"NOP",
			instTypeOneReg{
				operation: opRPINC,
				imm:       0,
				srcReg:    0,
			},
		},
		{
			"RPINC 100",
			instTypeOneReg{
				operation: opRPINC,
				imm:       100,
				srcReg:    0,
			},
		},
		{
			"LD.8 	124 30",
			instTypeOneReg{
				operation: opLD8,
				srcReg:    124,
				imm:       30,
			},
		},
		{
			"SRLi.32 12 24",
			instTypeOneReg{
				operation: opSRLi32,
				srcReg:    12,
				imm:       24 << 5,
			},
		},
		{
			"SRAi.64 12 24",
			instTypeOneReg{
				operation: opSRLi64,
				srcReg:    12,
				imm:       24<<5 | 8,
			},
		},
	}
	for _, e := range table {
		compare(e.in, e.expected)
	}
}

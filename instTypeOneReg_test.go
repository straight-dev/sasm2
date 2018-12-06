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

	var oneregtests = []struct {
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
				imm:       30,
				srcReg:    124,
			},
		},
	}
	for _, e := range oneregtests {
		compare(e.in, e.expected)
	}
}

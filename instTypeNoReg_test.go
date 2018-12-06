package main

import (
	"testing"
)

func TestFromStringToInstTypeNoReg(t *testing.T) {
	compare := func(s string, expected instTypeNoReg) {
		actual, err := fromStringToInstTypeNoReg(s)
		if err != nil {
			t.Error(err.Error())
		}

		if *actual != expected {
			t.Error(actual, expected)
		}

	}

	var table = []struct {
		in       string
		expected instTypeNoReg
	}{
		{
			"J 100",
			instTypeNoReg{
				operation: opJ,
				imm20:     100,
			},
		},
		{
			"SPADDi -100",
			instTypeNoReg{
				operation: opSPADDi,
				imm20:     0xfffff + 1 - 100,
			},
		},
	}
	for _, e := range table {
		compare(e.in, e.expected)
	}

}

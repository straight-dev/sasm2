package main

import (
	"testing"
)

func TestFromStringToInstTypeSB(t *testing.T) {
	compare := func(s string, expected instTypeSB) {
		actual, err := fromStringToInstTypeSB(s)
		if err != nil {
			t.Error(err.Error())
		}

		if *actual != expected {
			t.Error(actual, expected)
		}

	}

	var table = []struct {
		in       string
		expected instTypeSB
	}{
		{
			"ST.8 1 2 -10 ",
			instTypeSB{
				operation: opST8,
				imm12:     0xffffffff - 9,
				srcRegs:   [2]uint32{1, 2},
			},
		},
	}
	for _, e := range table {
		compare(e.in, e.expected)
	}

}

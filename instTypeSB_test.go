package main

import (
	"testing"
)

func TestFromStringToInstTypeSB(t *testing.T) {
	actual, err := fromStringToInstTypeSB("ST.8 1 2 -10")
	expected := instTypeSB{
		operation: opST8,
		imm12:     0xffffffff - 9,
		srcRegs:   [2]uint32{1, 2},
	}
	if err != nil {
		t.Error(err.Error())
	}

	if *actual != expected {
		t.Error(actual, expected)
	}

}

package main

import (
	"testing"
)

func TestFromStringToInstTypeMAC(t *testing.T) {
	compare := func(s string, expected instTypeMAC) {
		actual, err := fromStringToInstTypeMAC(s)
		if err != nil {
			t.Error(err.Error())
		}

		if *actual != expected {
			t.Error(actual, expected)
		}

	}
	expected := instTypeMAC{
		operation: opFMADDd,
		rm:        rmRDN,
		srcRegs:   [3]uint32{1, 2, 3},
	}
	compare("FMADD.d 1 2 3 RDN", expected)

}

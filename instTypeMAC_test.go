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
	var table = []struct {
		in       string
		expected instTypeMAC
	}{
		{
			"FMADD.d 1 2 3 RDN",
			instTypeMAC{
				operation: opFMADDd,
				rm:        rmRDN,
				srcRegs:   [3]uint32{1, 2, 3},
			},
		},
	}
	for _, e := range table {
		compare(e.in, e.expected)
	}

}

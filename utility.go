package main

func extractBits(imm uint64, bit uint, isSigned bool) uint64 {
	val := ((1 << bit) - 1) & imm
	if isSigned && ((imm & (1 << (bit - 1))) != 0) {
		val |= ^((1 << bit) - 1)
	}
	return val
}

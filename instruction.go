package main

// Instruction interface
type instruction interface {
	toUInt32() uint32
	toBytes() [4]byte
}

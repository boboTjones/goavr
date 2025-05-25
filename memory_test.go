package main

import (
	"testing"
)

func TestFetch(t *testing.T) {
	mem := Memory{}
	cpu.pc = 0

	// Load some data into memory
	mem[0] = 0x12
	mem[1] = 0x34

	mem.Fetch()
	if cpu.pc != 2 {
		t.Errorf("Expected PC to be 2, got %d", cpu.pc)
	}
	if mem[0] != 0x12 || mem[1] != 0x34 {
		t.Error("Fetch did not retrieve correct data")
	}
}

func TestRead(t *testing.T) {
	mem := Memory{}

	// Test writing and reading from memory
	mem[10] = 0x55
	if got := mem.Read(10); got != 0x55 {
		t.Errorf("Expected to read 0x55 from memory, got %x", got)
	}
}

func TestLoadProgram(t *testing.T) {
	mem := Memory{}
	program := []byte{0x01, 0x02, 0x03, 0x04}

	mem.LoadProgram(program)

	for i, b := range program {
		if mem[i] != b {
			t.Errorf("Expected memory[%d] to be 0x%x, got 0x%x", i, b, mem[i])
		}
	}
}

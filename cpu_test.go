package main

import (
	"testing"
)

func TestNewCPU(t *testing.T) {
	cpu := NewCPU()

	if cpu.pc != 0 {
		t.Errorf("Expected PC to be 0, got %d", cpu.pc)
	}
	if cpu.sr != 0x3f {
		t.Errorf("Expected SR to be 0x3f, got %x", cpu.sr)
	}
}

func TestSetClearStatusRegister(t *testing.T) {
	cpu := NewCPU()

	// Set the zero flag
	cpu.set_z()
	if cpu.get_z() != 1 {
		t.Error("Expected Z flag to be set")
	}

	// Clear the zero flag
	cpu.clear_z()
	if cpu.get_z() != 0 {
		t.Error("Expected Z flag to be cleared")
	}
}

func TestIncrementXRegister(t *testing.T) {
	cpu := NewCPU()

	// Initial value of X register should be 0
	if cpu.xAddr() != 0 {
		t.Error("Expected initial X address to be 0")
	}

	cpu.incX()
	if cpu.xAddr() != 1 {
		t.Error("Expected X address to be 1 after increment")
	}
}

func TestDecrementYRegister(t *testing.T) {
	cpu := NewCPU()

	// Set Y register to 1
	cpu.dmem[y_low] = 1
	cpu.dmem[y_high] = 0

	cpu.decY()
	if cpu.yAddr() != 0 {
		t.Error("Expected Y address to be 0 after decrement")
	}
}

func TestBranchExecution(t *testing.T) {
	cpu := NewCPU()
	cpu.pc = 0x0020

	cpu.Branch(0x10)
	if cpu.pc != 0x0030 {
		t.Errorf("Expected PC to be 0x0030, got %x", cpu.pc)
	}
}

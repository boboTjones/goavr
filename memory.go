package main

import (
	"encoding/hex"
	"fmt"
)

var data []byte
var current []byte

// 4096 max size of program memory.
// Both data memory and program memory inherit
// from here.

type Memory [8192]byte

// Use Fetch() for grabbing 2 bytes from program memory.
// Increments the program counter.

func (mem *Memory) Fetch() {
	current = mem[cpu.pc:(cpu.pc + 2)]
	cpu.pc += 2
}

// Use Read() for reading a single byte from data memory
func (mem *Memory) Read(loc int) byte {
	return mem[loc]
}

// Loads the executable stuff into program memory.

func (mem *Memory) LoadProgram(data []byte) {
	for i, b := range data {
		//fmt.Println(i)
		mem[i] = b
	}
}

func (mem *Memory) Dump() string {
	return hex.Dump(mem[0:])
}

// Here but unused.
func (mem *Memory) Store(i uint16, b byte) {
	mem[i] = b
}

func (mem *Memory) printRegs() {
	var regs []string
	for i := 0; i < 32; i++ {
		regs = append(regs, fmt.Sprintf("r%d[%.2x]", i, mem[i]))
	}
	fmt.Println("Registers:")
	fmt.Printf("%v\n", regs[0:13])
	fmt.Printf("%v\n", regs[13:26])
	fmt.Printf("X:\t%v\t%v\n", regs[26], regs[27])
	fmt.Printf("Y:\t%v\t%v\n", regs[28], regs[29])
	fmt.Printf("Z:\t%v\t%v\n", regs[30], regs[31])
}

func (mem *Memory) printStack() {
	var stack []string
	se := b2u16little(stackEnd) + 1
	if cpu.sp.current() != 0 && cpu.sp.current() < se {
		for _, v := range mem[(cpu.sp.current() + 1):se] {
			stack = append(stack, fmt.Sprintf("%.2x ", v))
		}
	}
	fmt.Println("Stack:")
	fmt.Println(stack)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	x_low = iota + 26
	x_high
	y_low
	y_high
	z_low
	z_high
)

var programEnd uint16

var bitMasks = []byte{1, 2, 4, 8, 16, 32, 64, 128}

type CPU struct {
	pc      uint16
	sp      StackPointer
	sr      byte
	imem    Memory
	dmem    Memory
	objdump bool
}

// Set bits in status register
func (cpu *CPU) set_i() { cpu.dmem[cpu.sr] |= 128 }
func (cpu *CPU) set_t() { cpu.dmem[cpu.sr] |= 64 }
func (cpu *CPU) set_h() { cpu.dmem[cpu.sr] |= 32 }
func (cpu *CPU) set_s() { cpu.dmem[cpu.sr] |= 16 }
func (cpu *CPU) set_v() { cpu.dmem[cpu.sr] |= 8 }
func (cpu *CPU) set_n() { cpu.dmem[cpu.sr] |= 4 }
func (cpu *CPU) set_z() { cpu.dmem[cpu.sr] |= 2 }
func (cpu *CPU) set_c() { cpu.dmem[cpu.sr] |= 1 }

// Clear bits in status register
func (cpu *CPU) clear_i() { cpu.dmem[cpu.sr] &= 127 }
func (cpu *CPU) clear_t() { cpu.dmem[cpu.sr] &= 191 }
func (cpu *CPU) clear_h() { cpu.dmem[cpu.sr] &= 223 }
func (cpu *CPU) clear_s() { cpu.dmem[cpu.sr] &= 239 }
func (cpu *CPU) clear_v() { cpu.dmem[cpu.sr] &= 247 }
func (cpu *CPU) clear_n() { cpu.dmem[cpu.sr] &= 251 }
func (cpu *CPU) clear_z() { cpu.dmem[cpu.sr] &= 253 }
func (cpu *CPU) clear_c() { cpu.dmem[cpu.sr] &= 254 }

// Get bits from status register
func (cpu *CPU) get_i() byte { return (cpu.dmem[cpu.sr] & bitMasks[7]) >> 7 }
func (cpu *CPU) get_t() byte { return (cpu.dmem[cpu.sr] & bitMasks[6]) >> 6 }
func (cpu *CPU) get_h() byte { return (cpu.dmem[cpu.sr] & bitMasks[5]) >> 5 }
func (cpu *CPU) get_s() byte { return (cpu.dmem[cpu.sr] & bitMasks[4]) >> 4 }
func (cpu *CPU) get_v() byte { return (cpu.dmem[cpu.sr] & bitMasks[3]) >> 3 }
func (cpu *CPU) get_n() byte { return (cpu.dmem[cpu.sr] & bitMasks[2]) >> 2 }
func (cpu *CPU) get_z() byte { return (cpu.dmem[cpu.sr] & bitMasks[1]) >> 1 }
func (cpu *CPU) get_c() byte { return (cpu.dmem[cpu.sr] & bitMasks[0]) }

// Get the address in the X/Y/Z register
func (cpu *CPU) zAddr() uint16 { return b2u16little([]byte{cpu.dmem[z_low], cpu.dmem[z_high]}) }
func (cpu *CPU) xAddr() uint16 { return b2u16little([]byte{cpu.dmem[x_low], cpu.dmem[x_high]}) }
func (cpu *CPU) yAddr() uint16 { return b2u16little([]byte{cpu.dmem[y_low], cpu.dmem[y_high]}) }

// Increment the X/Y/Z registers

func (cpu *CPU) incX() {
	x := cpu.xAddr() + 1
	b := u16lil2byte(x)
	cpu.dmem[x_low] = b[0]
	cpu.dmem[x_high] = b[1]
}
func (cpu *CPU) incY() {
	y := cpu.yAddr() + 1
	b := u16lil2byte(y)
	cpu.dmem[y_low] = b[0]
	cpu.dmem[y_high] = b[1]
}
func (cpu *CPU) incZ() {
	z := cpu.zAddr() + 1
	b := u16lil2byte(z)
	cpu.dmem[z_low] = b[0]
	cpu.dmem[z_high] = b[1]
}

// Decrement the X/Y/Z registers

func (cpu *CPU) decX() {
	x := cpu.xAddr() - 1
	b := u16lil2byte(x)
	cpu.dmem[x_low] = b[0]
	cpu.dmem[x_high] = b[1]
}
func (cpu *CPU) decY() {
	y := cpu.yAddr() - 1
	b := u16lil2byte(y)
	cpu.dmem[y_low] = b[0]
	cpu.dmem[y_high] = b[1]
}
func (cpu *CPU) decZ() {
	z := cpu.zAddr() - 1
	b := u16lil2byte(z)
	cpu.dmem[z_low] = b[0]
	cpu.dmem[z_high] = b[1]
}

// Massaging the signed offset into the unsigned cpu to branch

func (cpu *CPU) Branch(offset int16) { cpu.pc = (cpu.pc + uint16(offset)) % 8192 }

// Get the return value of a function for testing
func (cpu *CPU) getReturnValue() uint16 { return b2u16little([]byte{cpu.dmem[24], cpu.dmem[25]}) }

/*
Golang Logical Operators: (because I'm tired of looking this shit up)
+    ADD
-    SUB
&    bitwise AND
|    bitwise OR
^    bitwise XOR
&^   bit clear (AND NOT)

<<   left shift
>>   right shift
*/

func (cpu *CPU) Step() {
	//defer handlePanic()
	cpu.imem.Fetch()
	cpu.Execute(dissAssemble(current))
}

func (cpu *CPU) Run() {
	for {
		cpu.Step()
		if cpu.pc == programEnd {
			fmt.Println(cpu.getReturnValue())
			os.Exit(0)
		}
	}
}

func (cpu *CPU) Noise() {
	cpu.objdump = true
	fmt.Printf("pc: %.4x\tsr: %.8b\tsp: %.4x\t\n", cpu.pc, cpu.dmem[cpu.sr], cpu.sp.current())
	cpu.dmem.printRegs()
	cpu.dmem.printStack()
	fmt.Println("---------------------------------")
}

func (cpu *CPU) Interactive() {
	fmt.Println("Type ? for help.")
	for {
		prompt := bufio.NewReader(os.Stdin)
		fmt.Print("$> ")

		response, err := prompt.ReadString('\n')

		check(err)
		// Ugh.
		r := strings.Split(response, "\n")

		switch r[0] {
		case "?":
			fmt.Println("g to run the whole program")
			fmt.Println("q to quit")
			fmt.Println("s to single step")
			fmt.Println("j prompts for a pc (in hex)")
			fmt.Println("d dumps the data memory")
			fmt.Println("p dumps the program mempory")
			fmt.Println("return jumps 5 instructions")
			fmt.Println("any number /n/ jumps /n/ instructions")
		case "g":
			for {
				cpu.Step()
				cpu.Noise()
				if cpu.pc == programEnd {
					break
				}
			}
		case "q":
			os.Exit(0)
		case "s":
			var b string
			for {
				cpu.Step()
				cpu.Noise()
				fmt.Scanf("%s", &b)
				if b == "x" {
					break
				}
				if cpu.pc == programEnd {
					break
				}
			}
		case "b":
			cpu.pc -= 2
		case "r":
			cpu.pc = 0x0026
		case "d":
			fmt.Println(cpu.dmem.Dump())
		case "p":
			fmt.Println(cpu.imem.Dump())
		case "j":
			var o uint16
			fmt.Println("Enter pc:")
			fmt.Scanf("%x", &o)
			for cpu.pc < (o + 2) {
				cpu.Step()
				if cpu.pc == programEnd {
					break
				}
			}
		default:
			var n int
			// default case: step n times
			if r[0] == "" {
				n = 5
			} else {
				n, err = strconv.Atoi(r[0])
				if err != nil {
					fmt.Println("Command not recognized.")
					break
				}
			}
			for i := 0; i < n; i++ {
				cpu.Step()
				cpu.Noise()
				if cpu.pc == programEnd {
					break
				}
			}
		}
	}
}

func (cpu *CPU) Execute(i Instr) {
	if cpu.objdump == true {
		fmt.Printf("%.4x\t%s\n", cpu.pc, i.objdump)
	}

	switch i.label {
	case INSN_UNK:
		os.Exit(1)
		return
	case INSN_NOP:
		// Duh.
		return
	case INSN_CLI:
		// Clear global interrupt
		cpu.clear_i()
		return
	case INSN_CLT:
		cpu.clear_t()
		return
	case INSN_SET:
		cpu.set_t()
		return
	case INSN_SEC:
		cpu.set_c()
		return
	case INSN_JMP:
		// we all know this doesn't work because
		// this version doesn't have a 22bit pc
		cpu.pc = uint16(i.k32)
		return
	case INSN_IJMP:
		// PC <- Z(15:0)
		// shift one bit because word addressing something something
		cpu.pc = cpu.zAddr() << 1
		return
	case INSN_RJMP:
		// PC <- PC + k + 1
		cpu.Branch(i.k16)
		return
	case INSN_ADD:
		// Rd <- Rd + Rr
		// Casting to uint16 to catch overflow for status registers
		r := uint16(cpu.dmem[i.dest]) + uint16(cpu.dmem[i.source])
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		// this is some lazy shit right here.
		if r > 0xff {
			cpu.set_v()
			cpu.set_c()
		} else {
			cpu.clear_v()
			cpu.clear_c()
		}
		s := cpu.get_n() ^ cpu.get_v()
		if s == 0 {
			cpu.clear_s()
		} else {
			cpu.set_s()
		}
		cpu.dmem[i.dest] = byte(r)
		return
	case INSN_ADC:
		// Rd <- Rd + Rr + C
		// Casting to uint16 to catch overflow for status registers
		r := uint16(cpu.dmem[i.dest]) + uint16(cpu.dmem[i.source]) + uint16(cpu.get_c())
		cpu.dmem[i.dest] = byte(r)
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if r > 0x00ff {
			cpu.set_v()
			cpu.set_c()
		} else {
			cpu.clear_v()
			cpu.clear_c()
		}
		if ((r & 0x0080) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		s := cpu.get_n() ^ cpu.get_v()
		if s == 0 {
			cpu.clear_s()
		} else {
			cpu.set_s()
		}
		return
	case INSN_AND:
		// Rd <- Rd & Rr
		cpu.clear_v()
		r := cpu.dmem[i.dest] & cpu.dmem[i.source]
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if (r & 0x80 >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		s := cpu.get_n() ^ cpu.get_v()
		if s == 0 {
			cpu.clear_s()
		} else {
			cpu.set_s()
		}
		cpu.dmem[i.dest] = r
		return
	case INSN_ANDI:
		// Rd <- Rd & K
		r := cpu.dmem[i.dest] & i.kdata
		cpu.clear_v()
		if (r&0x80)>>7 == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		s := cpu.get_n() ^ cpu.get_v()
		if s == 0 {
			cpu.clear_s()
		} else {
			cpu.set_s()
		}
		cpu.dmem[i.dest] = r
		return
	case INSN_EOR:
		// Rd <- Rd^Rr
		cpu.clear_v()
		r := cpu.dmem[i.dest] ^ cpu.dmem[i.source]
		if (r&0x80)>>7 == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		s := cpu.get_n() ^ cpu.get_v()
		if s == 0 {
			cpu.clear_s()
		} else {
			cpu.set_s()
		}
		cpu.dmem[i.dest] = r
		return
	case INSN_IN:
		// Rd <- I/O(A)
		cpu.dmem[i.dest] = cpu.dmem[i.ioaddr]
		return
	case INSN_OUT:
		// I/O(A) <- Rr
		switch i.ioaddr {
		case 0x3e: // high byte
			if cpu.dmem[i.ioaddr] == 0 {
				cpu.sp.setStackEnd(cpu.dmem[i.source], 1)
			}
		case 0x3d: // low byte
			if cpu.dmem[i.ioaddr] == 0 {
				cpu.sp.setStackEnd(cpu.dmem[i.source], 0)
			}
		}
		cpu.dmem[i.ioaddr] = cpu.dmem[i.source]
		return
	case INSN_LDI:
		// Rd <- K
		cpu.dmem[i.dest] = i.kdata
		return
	case INSN_SBI:
		// I/O(A,b) <- 1
		// bitMasks is a map that looks up the mask to isolate
		// the necessary bit
		cpu.dmem[i.ioaddr] |= bitMasks[i.registerBit]
		return
	case INSN_CBI:
		// I/O(A,b) <- 0
		cpu.dmem[i.ioaddr] ^= bitMasks[i.registerBit]
		//fmt.Printf("%.8b\n", cpu.dmem[i.ioaddr])
		return
	case INSN_SBIC:
		// If I/O(A,b) = 0 then PC <- PC + 2 (or 3) else PC <- PC + 1
		r := cpu.dmem[i.ioaddr] & bitMasks[i.registerBit] >> i.registerBit
		if r == 0 {
			// instructions are 1 word
			cpu.pc += 2
		}
		return
	case INSN_SBIS:
		// If I/O(A,b) = 1 then PC <- PC + 2
		r := cpu.dmem[i.ioaddr] & bitMasks[i.registerBit] >> i.registerBit
		if r == 1 {
			// instructions are 1 word
			cpu.pc += 2
		}
		return
	case INSN_BLD:
		// Rd(b) <- T
		// Copies the T Flag in the SREG (Status Register) to bit b in register Rd.
		t := (cpu.dmem[cpu.sr] & bitMasks[7]) >> i.registerBit
		cpu.dmem[i.dest] |= byte(t)
		return
	case INSN_BST:
		// T <- Rd(b)
		// Stores bit b from Rd to the T Flag in SREG (Status Register).
		t := cpu.dmem[i.dest] & bitMasks[i.registerBit] >> i.registerBit
		if t == 1 {
			cpu.set_t()
		} else {
			cpu.clear_t()
		}
		return
	case INSN_SBRC:
		// if Rr(b) = 0 then PC += 2
		r := (cpu.dmem[i.source] & bitMasks[i.registerBit]) >> i.registerBit
		if r == 0 {
			cpu.pc += 2
		}
		return
	case INSN_SBRS:
		// if Rr(b) = 1 then PC += 2
		r := (cpu.dmem[i.source] & bitMasks[i.registerBit]) >> i.registerBit
		if r == 1 {
			cpu.pc += 2
		}
		return
	case INSN_STS:
		// (k) <- Rr
		cpu.dmem[i.k16] = cpu.dmem[i.source]
		return
	case INSN_LDS:
		// Rd <- (k)
		cpu.dmem[i.dest] = cpu.dmem[i.k16]
		return
	case INSN_ADIW:
		// Rd+1:Rd <- Rd+1:Rd + K
		// low byte
		low := uint16(cpu.dmem[i.dest])
		// high byte
		high := uint16(cpu.dmem[i.dest+1])
		r := ((high << 8) | low) + uint16(i.kdata)
		cpu.dmem[i.dest] = byte(r & 0x00ff)
		cpu.dmem[i.dest+1] = byte(r >> 8)
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if (r&0x0080)>>7 == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r > 0xff {
			cpu.set_v()
			cpu.set_c()
		} else {
			cpu.clear_v()
			cpu.clear_c()
		}
		s := cpu.get_n() ^ cpu.get_v()
		if s == 0 {
			cpu.clear_s()
		} else {
			cpu.set_s()
		}
		return
	case INSN_SBIW:
		// Rd+1:Rd <- Rd+1:Rd - K
		if i.kdata > cpu.dmem[i.dest] {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		x := uint16(cpu.dmem[i.dest+1])<<8 | uint16(cpu.dmem[i.dest])
		r := x - uint16(i.kdata)
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if (r&0x0080)>>7 == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r > 0xff {
			cpu.set_v()
		} else {
			cpu.clear_v()
		}
		if uint16(i.kdata) > x {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		s := cpu.get_n() ^ cpu.get_v()
		if s == 0 {
			cpu.clear_s()
		} else {
			cpu.set_s()
		}
		cpu.dmem[i.dest] = byte(r & 0x00ff)
		cpu.dmem[i.dest+1] = byte(r >> 8)
		return
	case INSN_BRCC:
		// Branch if carry cleared
		if cpu.get_c() == 0 {
			cpu.Branch(i.k16)
		}
		return
	case INSN_BRCS:
		// Branch if carry set
		if cpu.get_c() == 1 {
			cpu.Branch(i.k16)
		}
		return
	case INSN_BREQ:
		//if Rd = Rr(Z=1) then PC <- PC + k + 1
		if cpu.get_z() == 1 {
			cpu.Branch(i.k16)
		}
		return
	case INSN_BRNE:
		// if (Z = 0) then PC <-  PC + k + 1
		if cpu.get_z() == 0 {
			cpu.Branch(i.k16)
		}
		return
	case INSN_BRLT:
		// if (S = 1) then PC <-  PC + k + 1
		if cpu.get_s() == 1 {
			cpu.Branch(i.k16)
		}
		return
	case INSN_BRGE:
		// if Rd >= Rr (N ^ V = 0 or S = 0) then PC += k
		if cpu.get_s() == 0 {
			cpu.Branch(i.k16)
		}
		return
	case INSN_BRPL:
		// if (N = 0) then PC <-  PC + k + 1
		if cpu.get_n() == 0 {
			cpu.Branch(i.k16)
		}
		return
	case INSN_BRTC:
		// if T = 0 then PC <- PC + k + 1
		if cpu.get_t() == 0 {
			cpu.Branch(i.k16)
		}
		return
	case INSN_BRTS:
		// if T = 1 then PC <- PC + k + 1
		if cpu.get_t() == 1 {
			cpu.Branch(i.k16)
		}
		return
	case INSN_BRMI:
		// if T = 1 then PC <- PC + k + 1
		if cpu.get_n() == 1 {
			cpu.Branch(i.k16)
		}
		return
	case INSN_COM:
		// Rd <- ^Rd
		r := 0xff - cpu.dmem[i.dest]
		cpu.dmem[i.dest] = r
		cpu.clear_v()
		cpu.set_c()
		if ((r & 0x80) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		return
	case INSN_CP:
		// Rd - Rr
		d := int16(cpu.dmem[i.dest])
		s := int16(cpu.dmem[i.source])
		if d >= s {
			cpu.set_s()
		} else {
			cpu.set_s()
		}
		if s > d {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := d - s
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if ((r & 0x0080) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r > 0x00ff {
			cpu.set_v()
		} else {
			cpu.clear_v()
		}
		return
	case INSN_CPC:
		// Rd - Rr - C
		d := int16(cpu.dmem[i.dest])
		s := int16(cpu.dmem[i.source])
		c := int16(cpu.get_c())
		if (s + c) > d {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := d - s - c
		if r != 0 {
			cpu.clear_z()
		}
		if ((r & 0x80) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r > 0xff {
			cpu.set_v()
		} else {
			cpu.clear_v()
		}
		if r < 0 {
			cpu.set_s()
		} else {
			cpu.clear_s()
		}
		return
	case INSN_CPI:
		// Rd - K
		// I can't tell from the doc, but I think this check
		// has to happen before K is subtracted. We'll see.
		if i.kdata > cpu.dmem[i.dest] {
			cpu.set_s()
			cpu.set_c()
		} else {
			cpu.clear_s()
			cpu.clear_c()
		}
		r := cpu.dmem[i.dest] - i.kdata
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if ((r & 0x80) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if (cpu.get_n() ^ cpu.get_v()) == 0 {
			cpu.clear_s()
		} else {
			cpu.set_s()
		}
		return
	case INSN_CPSE:
		// if Rd = Rr then PC <- PC + 2
		if cpu.dmem[i.dest] == cpu.dmem[i.source] {
			cpu.pc += 2
		}
		return
	case INSN_DEC:
		// Rd <- Rd - 1
		if cpu.dmem[i.dest] == 0x80 {
			cpu.set_v()
		} else {
			cpu.clear_v()
		}
		r := cpu.dmem[i.dest] - 1
		cpu.dmem[i.dest] = r
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if ((r & 0x80) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if (cpu.get_n() ^ cpu.get_v()) == 0 {
			cpu.clear_s()
		} else {
			cpu.set_s()
		}
		return
	case INSN_INC:
		// Rd <- Rd + 1
		if cpu.dmem[i.dest] == 0x7f {
			cpu.set_v()
		} else {
			cpu.clear_v()
		}
		r := cpu.dmem[i.dest] + 1
		cpu.dmem[i.dest] = r
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if ((r & 0x80) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if (cpu.get_n() ^ cpu.get_v()) == 0 {
			cpu.clear_s()
		} else {
			cpu.set_s()
		}
		return
	case INSN_LDDY:
		// Rd <- (Y + q)
		y := cpu.yAddr() + i.offset
		//fmt.Printf("%.4x\t%.4x\n", y, cpu.dmem[y])
		cpu.dmem[i.dest] = cpu.dmem[y]
		return
	case INSN_LDDZ:
		// Rd <- (Z + q)
		z := cpu.zAddr() + i.offset
		//fmt.Printf("%.4x\n", z)
		cpu.dmem[i.dest] = cpu.dmem[z]
		return
	case INSN_LDX:
		// Rd <- (X)
		cpu.dmem[i.dest] = cpu.dmem[cpu.xAddr()]
		return
	case INSN_LDXP:
		// Rd <- (X), X <- X + 1
		cpu.dmem[i.dest] = cpu.dmem[cpu.xAddr()]
		cpu.incX()
		return
	case INSN_LDXM:
		// X <- X-1, Rd <- (X)
		cpu.decX()
		cpu.dmem[i.dest] = cpu.dmem[cpu.xAddr()]
		return
	case INSN_LDY:
		// Rd <- (Y)
		cpu.dmem[i.dest] = cpu.dmem[cpu.yAddr()]
		return
	case INSN_LDYP:
		// Rd <- (Y), Y <- Y + 1
		cpu.dmem[i.dest] = cpu.dmem[cpu.yAddr()]
		// XXX TODO(ERIN) this could overflow into the high
		// byte someday.
		cpu.incY()
		return
	case INSN_LDYM:
		// Rd <- (Y), Y <- Y - 1
		// pre-decrement
		cpu.decY()
		cpu.dmem[i.dest] = cpu.dmem[cpu.yAddr()]
		return
	case INSN_LDZ:
		// Rd <- (Z) (dmem)
		cpu.dmem[i.dest] = cpu.dmem[cpu.zAddr()]
		return
	case INSN_LDZP:
		// Rd <- (Z) (dmem), Z <- Z - 1
		cpu.dmem[i.dest] = cpu.dmem[cpu.zAddr()]
		// post-decrement
		cpu.incZ()
		return
	case INSN_LDZM:
		// Rd <- (Z) (dmem), Z <- Z - 1
		// pre-decrement
		cpu.decZ()
		cpu.dmem[i.dest] = cpu.dmem[cpu.zAddr()]
		return
	case INSN_LPMZ:
		// Rd <- (Z) (imem)
		cpu.dmem[i.dest] = cpu.imem[cpu.zAddr()]
		return
	case INSN_LPMZP:
		// Rd <- (Z), Z <- Z + 1 (imem)
		cpu.dmem[i.dest] = cpu.imem[cpu.zAddr()]
		// post-increment
		cpu.incZ()
		return
	case INSN_LPM:
		// R0 <- (Z) (imem)
		cpu.dmem[0] = cpu.imem[cpu.zAddr()]
		return
	case INSN_MOV:
		// Rd <- Rr
		cpu.dmem[i.dest] = cpu.dmem[i.source]
		return
	case INSN_MOVW:
		// Rd+1:Rd <- Rr+1:Rr
		cpu.dmem[i.dest+1] = cpu.dmem[i.source+1]
		cpu.dmem[i.dest] = cpu.dmem[i.source]
		return
	case INSN_NEG:
		// Rd <- $00 - Rd
		// Replaces the contents of register Rd with its two's complement
		d := uint16(cpu.dmem[i.dest])
		r := 0x00 - d
		if r == 0 {
			cpu.set_z()
			cpu.clear_c()
		} else {
			cpu.clear_z()
			cpu.set_c()
		}
		if ((r & 0x0080) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r > 0xff {
			cpu.set_v()
		} else {
			cpu.clear_v()
		}
		cpu.dmem[i.dest] = byte(r)
		return
	case INSN_OR:
		// Rd <- Rd | Rr
		r := cpu.dmem[i.dest] | cpu.dmem[i.source]
		cpu.dmem[i.dest] = r
		cpu.clear_v()
		if ((r & 12) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		return
	case INSN_ORI:
		// Rd <- Rd | K
		cpu.clear_v()
		r := cpu.dmem[i.dest] | byte(i.kdata)
		cpu.dmem[i.dest] = r
		if ((r & 0x80) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		return
	case INSN_POP:
		// Rd <- Stack
		cpu.dmem[i.dest] = cpu.dmem[cpu.sp.current()+1]
		cpu.sp.inc(1)
		return
	case INSN_PUSH:
		// STACK <- Rr
		cpu.dmem[cpu.sp.current()] = cpu.dmem[i.source]
		cpu.sp.dec(1)
		return
	case INSN_RCALL:
		// PC <- PC + k + 1, STACK <- PC + 1, SP - 2
		// push the current PC onto the stack because
		// it is automaticaly incremented elsewhere.
		// low byte
		cpu.dmem[cpu.sp.current()] = byte(cpu.pc & 0x00ff)
		cpu.sp.dec(1)
		// high byte
		cpu.dmem[cpu.sp.current()] = byte(cpu.pc >> 8)
		cpu.sp.dec(1)
		// says +1, but that generates the wrong value
		// because the PC is incremented automaticaly anyway
		cpu.Branch(i.k16)
		return
	case INSN_ICALL:
		cpu.dmem[cpu.sp.current()] = byte(cpu.pc & 0x00ff)
		cpu.sp.dec(1)
		cpu.dmem[cpu.sp.current()] = byte(cpu.pc >> 8)
		cpu.sp.dec(1)
		cpu.pc = cpu.zAddr() << 1
		return
	case INSN_CALL:
		cpu.dmem[cpu.sp.current()] = byte(cpu.pc & 0x00ff)
		cpu.sp.dec(1)
		cpu.dmem[cpu.sp.current()] = byte(cpu.pc >> 8)
		cpu.sp.dec(1)
		cpu.pc = uint16(i.k32)
		return
	case INSN_RET:
		// PC <- Stack
		// r29
		h := uint16(cpu.dmem[cpu.sp.current()+1])
		cpu.sp.inc(1)
		// r 28
		l := uint16(cpu.dmem[cpu.sp.current()+1])
		cpu.sp.inc(1)
		cpu.pc = (h << 8) | l
		return
	case INSN_RETI:
		// PC <- Stack, enable interrupts
		low := cpu.dmem[cpu.sp.current()-1]
		high := cpu.dmem[cpu.sp.current()]
		cpu.pc = b2u16little([]byte{high, low})
		cpu.sp.inc(2)
		cpu.set_i()
		return
	case INSN_ROR:
		// Shifts all bits in Rd one place to the right.
		// The C Flag is shifted into bit 7 of Rd.
		// Bit 0 is shifted into the C Flag.
		// current C
		// new C flag
		x := cpu.dmem[i.dest] & 0x01
		r := (cpu.get_c() << 7) | (cpu.dmem[i.dest] >> 1)
		cpu.dmem[i.dest] = r
		if x == 1 {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if ((r & 0x80) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		return
	case INSN_LSR:
		// logical shift right Rd
		x := cpu.dmem[i.dest] & 0x01
		r := cpu.dmem[i.dest] >> 1
		cpu.dmem[i.dest] = r
		if x == 1 {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		cpu.clear_n()
		return
	case INSN_ASR:
		// Rd >> 1 but bit 7 remains constant. C becomes LSB of Rd.
		d := cpu.dmem[i.dest]
		c := d & 0x01
		b := d & 0x80
		r := (d >> 1) | b
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if c == 1 {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		if (r & 0x80) == 128 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if (cpu.get_n() ^ cpu.get_c()) == 1 {
			cpu.set_v()
		} else {
			cpu.clear_v()
		}
		if (cpu.get_n() ^ cpu.get_v()) == 1 {
			cpu.set_s()
		} else {
			cpu.clear_s()
		}
		cpu.dmem[i.dest] = r
		return
	case INSN_LSL:
		// logical shift left Rd
		x := cpu.dmem[i.dest] & 0x80 >> 7
		r := cpu.dmem[i.dest] << 1
		cpu.dmem[i.dest] = r
		if x == 1 {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if (r&0x80)>>7 == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		return
	case INSN_SUB:
		// Rd <- Rd - Rr
		d := cpu.dmem[i.dest]
		s := cpu.dmem[i.source]
		r := d - s
		if s > d {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if (r & 0x80 >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		cpu.dmem[i.dest] = r
		return
	case INSN_SBC:
		// Rd = Rd - Rr - C
		d := cpu.dmem[i.dest]
		s := cpu.dmem[i.source]
		c := cpu.get_c()
		r := d - s - c
		if (s + c) > d {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		if r != 0 {
			cpu.clear_z()
		}
		if ((r & 0x80) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		cpu.dmem[i.dest] = r
		return
	case INSN_SUBI:
		// Rd <- Rd - K
		k := uint16(i.kdata)
		d := uint16(cpu.dmem[i.dest])
		if k > d {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := d - k
		cpu.dmem[i.dest] = byte(r)
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if (r & 0x0080 >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r > 0xff {
			cpu.set_v()
		} else {
			cpu.clear_v()
		}
		if (cpu.get_n() ^ cpu.get_v()) == 1 {
			cpu.set_s()
		} else {
			cpu.clear_s()
		}
		return
	case INSN_SBCI:
		// Rd <- Rd - K - C
		c := uint16(cpu.get_c())
		k := uint16(i.kdata)
		d := uint16(cpu.dmem[i.dest])
		if (k + c) > d {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := d - k - c
		cpu.dmem[i.dest] = byte(r)
		if r == 0 {
			cpu.clear_z()
		}
		if (r & 0x0080 >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r > 0xff {
			cpu.set_v()
		} else {
			cpu.clear_v()
		}
		if (cpu.get_n() ^ cpu.get_v()) == 1 {
			cpu.set_s()
		} else {
			cpu.clear_s()
		}
		return
	case INSN_SEI:
		// set global interrupt flag
		cpu.set_i()
		return
	case INSN_MUL:
		// R1h:R0l <- Rx x Rr
		r := uint16(cpu.dmem[i.dest]) * uint16(cpu.dmem[i.source])
		cpu.dmem[1] = byte(r & 0xff00 >> 8)
		cpu.dmem[0] = byte(r & 0x00ff)
		if (r & 0x8000 >> 15) == 1 {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		return
	case INSN_STX:
		cpu.dmem[cpu.xAddr()] = cpu.dmem[i.source]
		return
	case INSN_STXP:
		// (X) <- Rr, X <- X + 1
		// 26 = low byte, 27 = high byte
		cpu.dmem[cpu.xAddr()] = cpu.dmem[i.source]
		// post-increment
		cpu.incX()
		return
	case INSN_STXM:
		// pre-decrement
		cpu.decX()
		cpu.dmem[cpu.xAddr()] = cpu.dmem[i.source]
		return
	case INSN_STY:
		cpu.dmem[cpu.yAddr()] = cpu.dmem[i.source]
		return
	case INSN_STYP:
		cpu.dmem[cpu.yAddr()] = cpu.dmem[i.source]
		// post-increment
		cpu.incY()
		return
	case INSN_STYM:
		// pre-decrement
		cpu.decY()
		cpu.dmem[cpu.yAddr()] = cpu.dmem[i.source]
		return
	case INSN_STDY:
		// (Y) <- Rr
		y := cpu.yAddr() + i.offset
		cpu.dmem[y] = cpu.dmem[i.source]
		return
	case INSN_STZ:
		cpu.dmem[cpu.zAddr()] = cpu.dmem[i.source]
		return
	case INSN_STZP:
		cpu.dmem[cpu.zAddr()] = cpu.dmem[i.source]
		// post-increment
		cpu.incZ()
		return
	case INSN_STZM:
		// pre-decrement
		cpu.decZ()
		cpu.dmem[cpu.zAddr()] = cpu.dmem[i.source]
		return
	case INSN_STDZ:
		z := cpu.zAddr() + i.offset
		cpu.dmem[z] = cpu.dmem[i.source]
		return
	default:
		fmt.Printf("CPU could not find instruction %s\n", getMnemonic(i.label))
		return
	}
}

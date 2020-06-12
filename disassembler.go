package main

import (
	"fmt"
)


func dissAssemble(b []byte) Instr {
	dm := map[byte]byte{
		0: 24,
		1: 26,
		2: 28,
		3: 30,
	}
	m := lookUp(b)
	inst := Instr{family: m.family, label: m.label}
	switch m.label {
	case INSN_UNK:
		inst.objdump = fmt.Sprintf("\nBZZZT! THANKS FOR PLAYING! %.4x\n", b2u16big(b))
		return inst
	case INSN_NOP:
		inst.objdump = fmt.Sprintf("%.4x\tnop\n", b2u16big(b))
		return inst
	case INSN_CLI:
		inst.objdump = fmt.Sprintf("%.4x\tcli\n", b2u16big(b))
		return inst
	case INSN_CLT:
		inst.objdump = fmt.Sprintf("%.4x\tclt\n", b2u16big(b))
		return inst
	case INSN_SET:
		inst.objdump = fmt.Sprintf("%.4x\tset\n", b2u16big(b))
		return inst
	case INSN_SEC:
		inst.objdump = fmt.Sprintf("%.4x\tsec\n", b2u16big(b))
		return inst
	case INSN_ADC:
		// 0001 11rd dddd rrrr
		inst.source = (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		inst.objdump = fmt.Sprintf("%.4x\tadc\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_EOR:
		// 0010 01rd dddd rrrr
		inst.source = (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		inst.objdump = fmt.Sprintf("%.4x\teor\tr%d, r%d\n", b2u16big(b), inst.source, inst.dest)
		return inst
	case INSN_OUT:
		// 1011 1AAr rrrr AAAA
		//out := (b[1] >> 3) & 0xff
		inst.ioaddr = ((b[1] & 0x06) << 3) | (b[0] & 0x0f)
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tout\t0x%.2x, r%d\t\t;%d\n", b2u16big(b), inst.ioaddr, inst.source, inst.ioaddr)
		return inst
	case INSN_LDI:
		// 1110 KKKK dddd KKKK
		inst.kdata = ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[0] & 0xf0) >> 4) + 0x10
		inst.objdump = fmt.Sprintf("%.4x\tldi\tr%d, 0x%.2x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	case INSN_RCALL:
		// 1101 kkkk kkkk kkkk
		k := (uint32(b[1]&0x0f)<<8 | uint32(b[0]))
		if ((k & 0x0800) >> 11) == 1 {
			inst.k16 = int16((k + 0xf000) << 1)
			inst.objdump = fmt.Sprintf("%.4x\trcall\t.%d\t;%.4x\n", b2u16big(b), inst.k16, (cpu.pc+uint16(inst.k16))%8192)
		} else {
			inst.k16 = int16(k << 1)
			inst.objdump = fmt.Sprintf("%.4x\trcall\t.+%d\t;%.4x\n", b2u16big(b), inst.k16, (cpu.pc+uint16(inst.k16))%8192)
		}
		return inst
	case INSN_SBI:
		// 1001 1010 AAAA Abbb
		inst.ioaddr = b[0] >> 3
		inst.registerBit = b[0] & 0x7
		inst.objdump = fmt.Sprintf("%.4x\tsbi\t0x%x, %d\n", b2u16big(b), inst.ioaddr, inst.registerBit)
		return inst
	case INSN_CBI:
		//1001 1000 AAAA Abbb
		inst.ioaddr = (b[0] & 0xf8) >> 3
		inst.registerBit = b[0] & 0x07
		inst.objdump = fmt.Sprintf("%.4x\tcbi\t0x%.2x, %d\n", b2u16big(b), inst.ioaddr, inst.registerBit)
		return inst
	case INSN_SBIC:
		// 1001 1001 AAAA Abbb
		inst.ioaddr = (b[0] & 0xf8) >> 3
		inst.registerBit = b[0] & 0x07
		inst.objdump = fmt.Sprintf("%.4x\tsbic\t0x%.2x, %d\n", b2u16big(b), inst.ioaddr, inst.registerBit)
		return inst
	case INSN_SBIS:
		// 1001 1011 AAAA Abbb
		inst.ioaddr = (b[0] & 0xf8) >> 3
		inst.registerBit = b[0] & 0x07
		inst.objdump = fmt.Sprintf("%.4x\tsbis\t0x%.2x, %d\n", b2u16big(b), inst.ioaddr, inst.registerBit)
		return inst
	case INSN_BLD:
		// 1111 100d dddd 0bbb
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.registerBit = b[0] & 0x07
		inst.objdump = fmt.Sprintf("%.4x\tbld\tr%d, %d\n", b2u16big(b), inst.dest, inst.registerBit)
		return inst
	case INSN_BST:
		// 1111 101d dddd 0bbb
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.registerBit = b[0] & 0x07
		inst.objdump = fmt.Sprintf("%.4x\tbst\tr%d, %d\n", b2u16big(b), inst.dest, inst.registerBit)
		return inst
	case INSN_SBRC:
		// 1111 110r rrrr 0bbb
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.registerBit = b[0] & 0x07
		inst.objdump = fmt.Sprintf("%.4x\tsbrc\tr%d, %d\n", b2u16big(b), inst.source, inst.registerBit)
		return inst
	case INSN_SBRS:
		// 1111 111r rrrr 0bbb
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.registerBit = b[0] & 0x07
		inst.objdump = fmt.Sprintf("%.4x\tsbrs\tr%d, %d\n", b2u16big(b), inst.source, inst.registerBit)
		return inst
	case INSN_STS:
		// 1001 001d dddd 0000 kkkk kkkk kkkk kkkk
		inst.source = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		cpu.imem.Fetch()
		inst.k16 = b2i16little(current)
		inst.objdump = fmt.Sprintf("%.4x\tsts\t0x%.4x, r%d\n", b2u16big(b), inst.k16, inst.source)
		return inst
	case INSN_LDS:
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		cpu.imem.Fetch()
		inst.k16 = b2i16little(current)
		inst.objdump = fmt.Sprintf("%.4x\tlds\tr%d, 0x%.4x\n", b2u16big(b), inst.dest, inst.k16)
		return inst
	case INSN_RJMP:
		// 1100 kkkk kkkk kkkk
		k := (uint32(b[1]&0x0f)<<8 | uint32(b[0]))
		if ((k & 0x800) >> 11) == 1 {
			inst.k16 = int16((k + 0xf000) << 1)
			inst.objdump = fmt.Sprintf("%.4x\trjmp\t.%d\n", b2u16big(b), inst.k16)
		} else {
			inst.k16 = int16(k << 1)
			inst.objdump = fmt.Sprintf("%.4x\trjmp\t.+%d\n", b2u16big(b), inst.k16)
		}
		return inst
	case INSN_IJMP:
		// 1001 0100 0000 1001
		inst.objdump = fmt.Sprintf("%.4x\tijmp\n", b2u16big(b))
		return inst
	case INSN_JMP:
		// 1001 010k kkkk 110k kkkk kkkk kkkk kkkk
		var k1, k2, k3 uint32
		k1 = uint32(b[1]&0x01) << 20
		k2 = uint32(b[0]&0xf0) << 12
		cpu.imem.Fetch()
		k3 = uint32(current[1])<<8 | uint32(current[0])
		inst.k32 = (k1 | k2 | k3) << 1
		inst.objdump = fmt.Sprintf("%.4x\tjmp\t0x%.8x\t;%d\n", b2u16big(b), inst.k32, inst.k32)
		return inst
	case INSN_CALL:
		// 1001 010k kkkk 111k kkkk kkkk kkkk kkkk
		var k1, k2, k3 uint32
		k1 = uint32(b[1]&0x01) << 20
		k2 = uint32(b[0]&0xf0) << 12
		cpu.imem.Fetch()
		k3 = uint32(current[1])<<8 | uint32(current[0])
		inst.k32 = (k1 | k2 | k3) << 1
		inst.objdump = fmt.Sprintf("%.4x\tcall\t0x%.8x\n", b2u16big(b), inst.k32)
		return inst		
	case INSN_ADD:
		// 0000 11rd dddd rrrr
		inst.source = (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		inst.objdump = fmt.Sprintf("%.4x\tadd\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_ADIW:
		// 1001 0110 KKdd KKKK
		// 24,26,28,30
		inst.kdata = ((b[0] & 0xc0) >> 2) | (b[0] & 0x0f)
		Rd := (b[0] & 0x30) >> 4
		inst.dest = dm[Rd]
		inst.objdump = fmt.Sprintf("%.4x\tadiw\tr%d, 0x%.2x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	case INSN_SBIW:
		// 1001 0111 KKdd KKKK
		inst.kdata = ((b[0] & 0xc0) >> 2) | (b[0] & 0x0f)
		Rd := (b[0] & 0x30) >> 4
		inst.dest = dm[Rd]
		inst.objdump = fmt.Sprintf("%.4x\tsbiw\tr%d, 0x%.2x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	case INSN_AND:
		// 0010 00rd dddd rrrr
		inst.source = (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		inst.objdump = fmt.Sprintf("%.4x\tand\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_ANDI:
		//0111 KKKK dddd KKKK
		inst.kdata = ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[0] & 0xf0) >> 4) + 0x10
		inst.objdump = fmt.Sprintf("%.4x\tandi\tr%d, 0x%.2x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	case INSN_BRCC:
		//1111 01kk kkkk k000
		// 64 ≤ k ≤ +63
		k := (b2u16little(b) & 0x03f8) >> 3
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
		} else {
			inst.k16 = int16(k << 1)
		}
		inst.objdump = fmt.Sprintf("%.4x\tbrcc\t.+%d\n", b2u16big(b), inst.k16)
		return inst
	case INSN_BRCS:
		// Supposed to be -64<k<+63, but avr-objdump doesn't display
		// these values this way.
		// 1111 00kk kkkk k000
		k := (b2u16little(b) & 0x03f8) >> 3
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrcs\t.%d\n", b2u16big(b), inst.k16)
		} else {
			inst.k16 = int16(k << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrcs\t.+%d\n", b2u16big(b), inst.k16)
		}
		return inst
	case INSN_BRMI:
		// 1111 00kk kkkk k010
		k := (b2u16little(b) & 0x03f8) >> 3
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrmi\t.%d\n", b2u16big(b), inst.k16)
		} else {
			inst.k16 = int16(k << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrmi\t.+%d\n", b2u16big(b), inst.k16)
		}
		return inst
	case INSN_BRGE:
		// 1111 01kk kkkk k100
		k := (b2u16little(b) & 0x03f8) >> 3
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrge\t.%d\n", b2u16big(b), inst.k16)
		} else {
			inst.k16 = int16(k << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrge\t.+%d\n", b2u16big(b), inst.k16)
		}
		return inst
	case INSN_BRNE:
		// 1111 01kk kkkk k001
		k := (b2u16little(b) & 0x03f8) >> 3
		// check to see if msb of k is 1
		// if it is, the result is negative.
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrne\t.%d\n", b2u16big(b), inst.k16)
		} else {
			inst.k16 = int16(k << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrne\t.+%d\n", b2u16big(b), inst.k16)
		}
		return inst
	case INSN_BRLT:
		// 1111 00kk kkkk k100
		k := (b2u16little(b) & 0x03f8) >> 3
		// check to see if msb of k is 1
		// if it is, the result is negative.
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrlt\t.%d\n", b2u16big(b), inst.k16)
		} else {
			inst.k16 = int16(k << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrlt\t.+%d\n", b2u16big(b), inst.k16)
		}
		return inst
	case INSN_BREQ:
		// Supposed to be -64<k<+63, but avr-objdump doesn't display
		// these values this way.
		// 1111 00kk kkkk k001
		k := (b2u16little(b) & 0x03f8) >> 3
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbreq\t.%d\n", b2u16big(b), inst.k16)
		} else {
			inst.k16 = int16(k << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbreq\t.+%d\n", b2u16big(b), inst.k16)
		}
		return inst
	case INSN_BRPL:
		// 1111 01kk kkkk k010
		k := (b2u16little(b) & 0x03f8) >> 3
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrpl\t.%d\n", b2u16big(b), inst.k16)
		} else {
			inst.k16 = int16(k << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrpl\t.+%d\n", b2u16big(b), inst.k16)
		}
		return inst
	case INSN_BRTC:
		// 1111 01kk kkkk k110
		k := (b2u16little(b) & 0x03f8) >> 3
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrtc\t.%d\n", b2u16big(b), inst.k16)
		} else {
			inst.k16 = int16(k << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrtc\t.+%d\n", b2u16big(b), inst.k16)
		}
		return inst
	case INSN_BRTS:
		// 1111 00kk kkkk k110
		k := (b2u16little(b) & 0x03f8) >> 3
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrts\t.%d\n", b2u16big(b), inst.k16)
		} else {
			inst.k16 = int16(k << 1)
			inst.objdump = fmt.Sprintf("%.4x\tbrts\t.+%d\n", b2u16big(b), inst.k16)
		}
		return inst
	case INSN_COM:
		// 1001 010d dddd 0000
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tcom\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_CP:
		// 0001 01rd dddd rrrr
		inst.source = (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		inst.objdump = fmt.Sprintf("%.4x\tcp\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_CPC:
		// 0000 01rd dddd rrrr
		inst.source = ((b[1] & 0x02) << 3) | b[0]&0x0f
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tcpc\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_CPI:
		// 0011 KKKK dddd KKKK
		inst.kdata = ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[0] & 0xf0) >> 4) + 0x10
		inst.objdump = fmt.Sprintf("%.4x\tcpi\tr%d, 0x%.2x\t;%d\n", b2u16big(b), inst.dest, inst.kdata, inst.kdata)
		return inst
	case INSN_CPSE:
		// 0001 00rd dddd rrrr
		inst.source = (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		inst.objdump = fmt.Sprintf("%.4x\tcpse\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_DEC:
		// 1001 010d dddd 1010
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tdec\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_INC:
		// 1001 010d dddd 1011
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tdec\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_IN:
		// 1011 0AAd dddd AAAA
		inst.ioaddr = ((b[1] & 0x06) << 3) | (b[0] & 0x0f)
		inst.dest = ((b[1] & 0xf1) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tin\tr%d, 0x%.2x\n", b2u16big(b), inst.dest, inst.ioaddr)
		return inst
	case INSN_LDDY:
		// 1001 000d dddd 1001
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.offset = m.offset
		inst.objdump = fmt.Sprintf("%.4x\tldd\tY+%d, r%d\n", b2u16big(b), inst.offset, inst.dest)
		return inst
	case INSN_LDDZ:
		// 1001 000d dddd 0001
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.offset = m.offset
		inst.objdump = fmt.Sprintf("%.4x\tldd\tr%d, Z+%d\n", b2u16big(b), inst.dest, inst.offset)
		return inst
	case INSN_LDX:
		// 1001 000d dddd 1100
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tld\tr%d, X\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LDXP:
		//1001 000d dddd 1101
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tld\tr%d, X+\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LDY:
		// 1000 000d dddd 1000
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tld\tr%d, Y\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LDYP:
		// 1001 000d dddd 1001
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tld\tr%d, Y+\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LDYM:
		// 1001 000d dddd 1010
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tld\tr%d, -Y\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LDZ:
		// 1000 000d dddd 0000
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tld\tr%d, Z\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LDZP:
		// 1001 000d dddd 0001
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tld\tr%d, Z+\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LDZM:
		// 1001 000d dddd 0010
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tld\tr%d, -Z\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LPMZ:
		//z  1001 000d dddd 0100
		// XXX ToDo not tested
		// XXX ToDo(Erin) Not sure this works.
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tlpm\tr%d, Z\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LPMZP:
		//z+ 1001 000d dddd 0101
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tlpm\tr%d, Z+\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LPM:
		// 1001 0101 1100 1000
		// XXX ToDo: not tested
		inst.dest = 0
		inst.objdump = fmt.Sprintf("%.4x\tlpm\n", b2u16big(b))
		return inst
	case INSN_LSR:
		//1001 010d dddd 0110
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tlsr\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_ASR:
		//1001 010d dddd 0101
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tasr\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LSL:
		//1001 11dd dddd 0110
		inst.dest = ((b[1] & 0x03) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tlsl\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_MOV:
		// 0010 11rd dddd rrrr
		inst.source = (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		inst.objdump = fmt.Sprintf("%.4x\tmov\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_MOVW:
		// 0000 0001 dddd rrrr
		inst.dest = (b[0] & 0xf0) >> 3
		inst.source = (b[0] & 0x0f) << 1
		inst.objdump = fmt.Sprintf("%.4x\tmovw\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_MUL:
		// 1001 11rd dddd rrrr
		inst.source = (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		inst.objdump = fmt.Sprintf("%.4x\tmul\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_NEG:
		// 1001 010d dddd s0001
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tneg\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_OR:
		// 0010 10rd dddd rrrr
		inst.source = (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		inst.objdump = fmt.Sprintf("%.4x\tor\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_ORI:
		// 0110 KKKK dddd KKKK
		inst.kdata = ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[0] & 0xf0) >> 4) + 0x10
		inst.objdump = fmt.Sprintf("%.4x\tori\tr%d, %.2x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	case INSN_POP:
		// 1001 000d dddd 1111
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tpop\t r%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_PUSH:
		//1001 001d dddd 1111
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tpush\tr%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_ICALL:
		inst.objdump = fmt.Sprintf("%.4x\ticall\n", b2u16big(b))
		return inst
	case INSN_RET:
		inst.objdump = fmt.Sprintf("%.4x\tret\n", b2u16big(b))
		return inst
	case INSN_RETI:
		inst.objdump = fmt.Sprintf("%.4x\treti\n", b2u16big(b))
		return inst
	case INSN_ROR:
		// 1001 010d dddd 0111
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tror\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_SBC:
		// 0000 10rd rrrr dddd
		inst.source = (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		inst.objdump = fmt.Sprintf("%.4x\tsbc\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_SBCI:
		// 0100 KKKK dddd KKKK
		inst.dest = ((b[0] & 0xf0) >> 4) + 0x10
		inst.kdata = (b[1]&0x0f)<<4 | (b[0] & 0x0f)
		inst.objdump = fmt.Sprintf("%.4x\tsbci\tr%d, 0x%x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	case INSN_SEI:
		// 1001 0100 0111 1000
		inst.objdump = fmt.Sprintf("%.4x\tsei\n", b2u16big(b))
		return inst
	case INSN_STDY:
		// 1001 001r rrrr 1001
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.offset = m.offset
		inst.objdump = fmt.Sprintf("%.4x\tstd\tY+%d, r%d\n", b2u16big(b), inst.offset, inst.source)
		return inst
	case INSN_STDZ:
		// 10q0 qq1r rrrr 0qqq
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.offset = m.offset
		inst.objdump = fmt.Sprintf("%.4x\tstd\tZ+%d, r%d\n", b2u16big(b), inst.offset, inst.source)
		return inst
	case INSN_STX:
		// 1001 001r rrrr 1100
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tst\tX, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STXP:
		// 1001 001r rrrr 1101
		//inst.source = (b2u16little(b) & 0x01f0) >> 4
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tst\tX+, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STXM:
		// 1001 001r rrrr 1110
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tst\t-X, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STY:
		//1001 001r rrrr 1000
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tst\tY, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STYP:
		//1001 001r rrrr 1001
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tst\tY+, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STYM:
		//1001 001r rrrr 1010
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tst\t-Y, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STZ:
		// 1000 001r rrrr 0000
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tst\tZ, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STZP:
		// 1001 001r rrrr 0001
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tst\tZ+, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STZM:
		// 1001 001r rrrr 0010
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.objdump = fmt.Sprintf("%.4x\tst\t-Z, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_SUB:
		// 0001 10rd dddd rrrr
		inst.source = (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		inst.objdump = fmt.Sprintf("%.4x\tsub\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_SUBI:
		// 0101 KKKK dddd KKKK
		inst.dest = ((b[0] & 0xf0) >> 4) + 0x10
		inst.kdata = (b[1]&0x0f)<<4 | (b[0] & 0x0f)
		inst.objdump = fmt.Sprintf("%.4x\tsubi\tr%d, 0x%x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	default:
		inst.objdump = fmt.Sprintf("None of the above. Got %s (0x%.4x)\n", m.mnemonic, b2u16big(b))
		return inst
	}

}

func lookUp(raw []byte) OpCode {
	var op OpCode
	b := b2u16little(raw)
	for _, entry := range OpCodeLookUpTable {
		v := b & entry.mask
		if v == entry.value {
			op = entry
			switch entry.mnemonic {
			case "std":
				return deConvoluter(raw, op)
			case "ldd":
				return deConvoluter(raw, op)
			}
			return op
		} else {
			op = OpCode{mnemonic: "unknown", value: b, label: INSN_UNK}
		}
	}
	return op
}

func getOffset(b []byte) uint16 {
	o0 := b[1] & 0x20
	o1 := (b[1] & 0x0c) << 1
	o2 := b[0] & 0x07
	o := o0 | o1 | o2
	return uint16(o)
}

func deConvoluter(b []byte, op OpCode) OpCode {
	x := b2u16little(b) & uint16(0xd208)
	offset := getOffset(b)
	switch x {
	case 0x8000:
		if offset == 0 {
			op.mnemonic = "ldz"
			op.label = INSN_LDZ
		} else {
			op.mnemonic = "lddz"
			op.offset = offset
			op.label = INSN_LDDZ
		}
	case 0x8008:
		if offset == 0 {
			op.mnemonic = "ldy"
			op.label = INSN_LDY
		} else {
			op.mnemonic = "lddy"
			op.offset = offset
			op.label = INSN_LDDY
		}
	case 0x8200:
		if offset == 0 {
			op.mnemonic = "stz"
			op.label = INSN_STZ
		} else {
			op.mnemonic = "stdz"
			op.offset = offset
			op.label = INSN_STDZ
		}
	case 0x8208:
		if offset == 0 {
			op.mnemonic = "sty"
			op.label = INSN_STY
		} else {
			op.mnemonic = "stdy"
			op.offset = offset
			op.label = INSN_STDY
		}
	default:
		op.mnemonic = "Unknown"
		op.value = b2u16little(b)
	}
	return op
}

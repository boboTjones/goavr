package main

//import "fmt"

var stackEnd = []byte{0x60, 0x04}

type StackPointer struct {
	high uint8
	low  uint8
}

func (sp *StackPointer) inc(x uint16) {
	d := b2u16little([]byte{cpu.dmem[sp.low], cpu.dmem[sp.high]}) + x
	r := u16big2byte(d)
	cpu.dmem[sp.high] = r[0]
	cpu.dmem[sp.low] = r[1]
}

func (sp *StackPointer) dec(x uint16) {
	d := b2u16little([]byte{cpu.dmem[sp.low], cpu.dmem[sp.high]}) - x
	r := u16big2byte(d)
	cpu.dmem[sp.high] = r[0]
	cpu.dmem[sp.low] = r[1]
}

func (sp *StackPointer) current() uint16 {
	return b2u16little([]byte{cpu.dmem[sp.low], cpu.dmem[sp.high]})
}

func (sp *StackPointer) set(b uint16) {
	r := u16big2byte(b)
	cpu.dmem[sp.high] = r[0]
	cpu.dmem[sp.low] = r[1]
}

func (sp *StackPointer) setStackEnd(s byte, x int) {
	stackEnd[x] = s
}

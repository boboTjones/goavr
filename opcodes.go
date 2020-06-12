package main

type family int

const (
	Arithmetic family = 0
	Branches          = 1
	Transfers         = 2
	BitWise           = 3
)

const (
	INSN_UNK = iota
	INSN_ADC
	INSN_ADD
	INSN_ADIW
	INSN_AND
	INSN_ANDI
	INSN_ASR
	INSN_BCLR
	INSN_BLD
	INSN_BRBC
	INSN_BRBS
	INSN_BRCC
	INSN_BRCS
	INSN_BREAK
	INSN_BREQ
	INSN_BRGE
	INSN_BRHC
	INSN_BRHS
	INSN_BRID
	INSN_BRIE
	INSN_BRLO
	INSN_BRLT
	INSN_BRMI
	INSN_BRNE
	INSN_BRPL
	INSN_BRSH
	INSN_BRTC
	INSN_BRTS
	INSN_BRVC
	INSN_BRVS
	INSN_BSET
	INSN_BST
	INSN_CALL
	INSN_CBI
	INSN_CBR
	INSN_CLC
	INSN_CLH
	INSN_CLI
	INSN_CLN
	INSN_CLR
	INSN_CLS
	INSN_CLT
	INSN_CLV
	INSN_CLZ
	INSN_COM
	INSN_CP
	INSN_CPC
	INSN_CPI
	INSN_CPSE
	INSN_DEC
	INSN_EICALL
	INSN_EIJMP
	INSN_ELPM
	INSN_EOR
	INSN_FMUL
	INSN_FMULS
	INSN_FMULSU
	INSN_ICALL
	INSN_IJMP
	INSN_IN
	INSN_INC
	INSN_JMP
	INSN_LAC
	INSN_LAS
	INSN_LAT
	INSN_LD
	INSN_LDD
	INSN_LDX
	INSN_LDXP
	INSN_LDXM
	INSN_LDDY
	INSN_LDY
	INSN_LDYP
	INSN_LDYM
	INSN_LDZ
	INSN_LDZP
	INSN_LDZM
	INSN_LDDZ
	INSN_LDI
	INSN_LDS
	INSN_LPM
	INSN_LPMZ
	INSN_LPMZP
	INSN_LSL
	INSN_LSR
	INSN_MOV
	INSN_MOVW
	INSN_MUL
	INSN_MULS
	INSN_MULSU
	INSN_NEG
	INSN_NOP
	INSN_OR
	INSN_ORI
	INSN_OUT
	INSN_POP
	INSN_PUSH
	INSN_RCALL
	INSN_RET
	INSN_RETI
	INSN_RJMP
	INSN_ROL
	INSN_ROR
	INSN_SBC
	INSN_SBCI
	INSN_SBI
	INSN_SBIC
	INSN_SBIS
	INSN_SBIW
	INSN_SBR
	INSN_SBRC
	INSN_SBRS
	INSN_SEC
	INSN_SEH
	INSN_SEI
	INSN_SEN
	INSN_SER
	INSN_SES
	INSN_SET
	INSN_SEV
	INSN_SEZ
	INSN_SLEEP
	INSN_SPM
	INSN_STD
	INSN_STDZ
	INSN_STDY
	INSN_STY
	INSN_STYP
	INSN_STYM
	INSN_STX
	INSN_STZ
	INSN_STZP
	INSN_STZM
	INSTN_STDY
	INSN_STXP
	INSN_STXM
	INSN_STS
	INSN_SUB
	INSN_SUBI
	INSN_SWAP
	INSN_TST
	INSN_WDR
	INSN_XCH
)

type OpCode struct {
	value    uint16
	mask     uint16
	mnemonic string
	offset   uint16
	family   family
	label    int
}

type Instr struct {
	family       family
	label        int
	offset       uint16
	dest         byte
	source       byte
	result       byte
	kdata        byte
	k16          int16
	k32          uint32
	ioaddr       byte
	iar          map[string]uint16
	displacement uint16
	registerBit  byte
	statusBit    uint16
	objdump      string
}

type Program []Instr

// Return a parsed opcode from a byte. Not used (yet)
func parseOpCode(b []byte) OpCode {
	o := OpCode{}
	return o
}

func getMnemonic(i int) string {
	m := "Not found."
	for _, e := range OpCodeLookUpTable {
		if e.label == i {
			m = e.mnemonic
		}
	}
	return m
}

// Put all your changes above this line cuz it's hairy
// down here.

var OpCodeLookUpTable = []OpCode{
	OpCode{
		mnemonic: "nop",
		mask:     0xffff,
		value:    0x0000,
		family:   BitWise,
		label:    INSN_NOP,
	},
	OpCode{
		mnemonic: "movw",
		mask:     0xff00,
		value:    0x0100,
		family:   Transfers,
		label:    INSN_MOVW,
	},
	OpCode{
		mnemonic: "eor",
		mask:     0xfc00,
		value:    0x2400,
		family:   Arithmetic,
		label:    INSN_EOR,
	},
	OpCode{
		mnemonic: "add",
		mask:     0xfc00,
		value:    0x0c00,
		family:   Arithmetic,
		label:    INSN_ADD,
	},
	OpCode{
		mnemonic: "adc",
		mask:     0xfc00,
		value:    0x1c00,
		family:   Arithmetic,
		label:    INSN_ADC,
	},
	OpCode{
		mnemonic: "and",
		mask:     0xfc00,
		value:    0x2000,
		family:   Arithmetic,
		label:    INSN_AND,
	},
	OpCode{
		mnemonic: "in",
		mask:     0xf800,
		value:    0xb000,
		family:   Transfers,
		label:    INSN_IN,
	},
	OpCode{
		mnemonic: "out",
		mask:     0xf800,
		value:    0xb800,
		family:   Transfers,
		label:    INSN_OUT,
	},
	OpCode{
		mnemonic: "rjmp",
		mask:     0xf000,
		value:    0xc000,
		family:   Branches,
		label:    INSN_RJMP,
	},
	OpCode{
		mnemonic: "ijmp",
		mask:     0xffff,
		value:    0x9409,
		family:   Branches,
		label:    INSN_IJMP,
	},
	OpCode{
		mnemonic: "rcall",
		mask:     0xf000,
		value:    0xd000,
		family:   Branches,
		label:    INSN_RCALL,
	},
	OpCode{
		mnemonic: "ldi",
		mask:     0xf000,
		value:    0xe000,
		family:   Transfers,
		label:    INSN_LDI,
	},
	OpCode{
		mnemonic: "cli",
		mask:     0xffff,
		value:    0x94f8,
		family:   BitWise,
		label:    INSN_CLI,
	},
	OpCode{
		mnemonic: "clt",
		mask:     0xffff,
		value:    0x94e8,
		family:   BitWise,
		label:    INSN_CLT,
	},
	OpCode{
		mnemonic: "set",
		mask:     0xffff,
		value:    0x9468,
		family:   BitWise,
		label:    INSN_SET,
	},
	OpCode{
		mnemonic: "sec",
		mask:     0xffff,
		value:    0x9408,
		family:   BitWise,
		label:    INSN_SEC,
	},
	OpCode{
		mnemonic: "sbi",
		mask:     0xff00,
		value:    0x9a00,
		family:   BitWise,
		label:    INSN_SBI,
	},
	OpCode{
		mnemonic: "cbi",
		mask:     0xff00,
		value:    0x9800,
		family:   BitWise,
		label:    INSN_CBI,
	},
	OpCode{
		mnemonic: "cpi",
		mask:     0xf000,
		value:    0x3000,
		family:   Branches,
		label:    INSN_CPI,
	},
	OpCode{
		mnemonic: "cpc",
		mask:     0xfc00,
		value:    0x0400,
		family:   Branches,
		label:    INSN_CPC,
	},
	OpCode{
		mnemonic: "subi",
		mask:     0xf000,
		value:    0x5000,
		family:   Arithmetic,
		label:    INSN_SUBI,
	},
	OpCode{
		mnemonic: "breq",
		mask:     0xfc07,
		value:    0xf001,
		family:   Branches,
		label:    INSN_BREQ,
	},
	OpCode{
		mnemonic: "brne",
		mask:     0xfc07,
		value:    0xf401,
		family:   Branches,
		label:    INSN_BRNE,
	},
	OpCode{
		mnemonic: "brlt",
		mask:     0xfc07,
		value:    0xf004,
		family:   Branches,
		label:    INSN_BRLT,
	},
	OpCode{
		mnemonic: "brmi",
		mask:     0xfc07,
		value:    0xf002,
		family:   Branches,
		label:    INSN_BRMI,
	},
	OpCode{
		mnemonic: "brpl",
		mask:     0xfc07,
		value:    0xf402,
		family:   Branches,
		label:    INSN_BRPL,
	},
	OpCode{
		mnemonic: "brcs",
		mask:     0xfc07,
		value:    0xf000,
		family:   Branches,
		label:    INSN_BRCS,
	},
	OpCode{
		mnemonic: "sbci",
		mask:     0xf000,
		value:    0x4000,
		family:   Arithmetic,
		label:    INSN_SBCI,
	},
	OpCode{
		mnemonic: "ori",
		mask:     0xf000,
		value:    0x6000,
		family:   Arithmetic,
		label:    INSN_ORI,
	},
	OpCode{
		mnemonic: "sbc",
		mask:     0xfc00,
		value:    0x0800,
		family:   Arithmetic,
		label:    INSN_SBC,
	},
	OpCode{
		mnemonic: "com",
		mask:     0xfe0f,
		value:    0x9400,
		family:   Arithmetic,
		label:    INSN_COM,
	},
	OpCode{
		mnemonic: "neg",
		mask:     0xfe0f,
		value:    0x9401,
		family:   Arithmetic,
		label:    INSN_NEG,
	},
	OpCode{
		mnemonic: "inc",
		mask:     0xfe0f,
		value:    0x9403,
		family:   Arithmetic,
		label:    INSN_INC,
	},
	OpCode{
		mnemonic: "asr",
		mask:     0xfe0f,
		value:    0x9405,
		family:   BitWise,
		label:    INSN_ASR,
	},
	OpCode{
		mnemonic: "lsr",
		mask:     0xfe0f,
		value:    0x9406,
		family:   BitWise,
		label:    INSN_LSR,
	},
	OpCode{
		mnemonic: "ror",
		mask:     0xfe0f,
		value:    0x9407,
		family:   BitWise,
		label:    INSN_ROR,
	},
	OpCode{
		mnemonic: "lsl",
		mask:     0xfe0f,
		value:    0x9c06,
		family:   BitWise,
		label:    INSN_LSL,
	},
	OpCode{
		mnemonic: "mov",
		mask:     0xfc00,
		value:    0x2c00,
		family:   Transfers,
		label:    INSN_MOV,
	},
	OpCode{
		mnemonic: "or",
		mask:     0xfc00,
		value:    0x2800,
		family:   Arithmetic,
		label:    INSN_OR,
	},
	OpCode{
		mnemonic: "push",
		mask:     0xfe0f,
		value:    0x920f,
		family:   Transfers,
		label:    INSN_PUSH,
	},
	OpCode{
		mnemonic: "pop",
		mask:     0xfe0f,
		value:    0x900f,
		family:   Transfers,
		label:    INSN_POP,
	},
	OpCode{
		mnemonic: "reti",
		mask:     0xffff,
		value:    0x9518,
		family:   Branches,
		label:    INSN_RETI,
	},
	OpCode{
		mnemonic: "icall",
		mask:     0xffff,
		value:    0x9509,
		family:   Branches,
		label:    INSN_ICALL,
	},
	OpCode{
		mnemonic: "ret",
		mask:     0xffff,
		value:    0x9508,
		family:   Branches,
		label:    INSN_RET,
	},
	OpCode{
		mnemonic: "sbis",
		mask:     0xff00,
		value:    0x9b00,
		family:   Branches,
		label:    INSN_SBIS,
	},
	OpCode{
		mnemonic: "cpse",
		mask:     0xfc00,
		value:    0x1000,
		family:   Branches,
		label:    INSN_CPSE,
	},
	OpCode{
		mnemonic: "sbiw",
		mask:     0xff00,
		value:    0x9700,
		family:   Arithmetic,
		label:    INSN_SBIW,
	},
	OpCode{
		mnemonic: "brcc",
		mask:     0xfc07,
		value:    0xf400,
		family:   Branches,
		label:    INSN_BRCC,
	},
	OpCode{
		mnemonic: "cp",
		mask:     0xfc00,
		value:    0x1400,
		family:   Branches,
		label:    INSN_CP,
	},
	OpCode{
		mnemonic: "adiw",
		mask:     0xff00,
		value:    0x9600,
		family:   Arithmetic,
		label:    INSN_ADIW,
	},
	OpCode{
		mnemonic: "andi",
		mask:     0xf000,
		value:    0x7000,
		family:   Arithmetic,
		label:    INSN_ANDI,
	},
	OpCode{
		mnemonic: "sbic",
		mask:     0xfd00,
		value:    0x9900,
		family:   Branches,
		label:    INSN_SBIC,
	},
	OpCode{
		mnemonic: "bst",
		mask:     0xfe08,
		value:    0xfa00,
		family:   BitWise,
		label:    INSN_BST,
	},
	OpCode{
		mnemonic: "bld",
		mask:     0xfe08,
		value:    0xf800,
		family:   BitWise,
		label:    INSN_BLD,
	},
	OpCode{
		mnemonic: "sei",
		mask:     0xffff,
		value:    0x9478,
		family:   BitWise,
		label:    INSN_SEI,
	},
	OpCode{
		mnemonic: "brge",
		mask:     0xfc07,
		value:    0xf404,
		family:   Branches,
		label:    INSN_BRGE,
	},
	OpCode{
		mnemonic: "brtc",
		mask:     0xfc07,
		value:    0xf406,
		family:   Branches,
		label:    INSN_BRTC,
	},
	OpCode{
		mnemonic: "brts",
		mask:     0xfc07,
		value:    0xf006,
		family:   Branches,
		label:    INSN_BRTS,
	},
	OpCode{
		mnemonic: "sbrc",
		mask:     0xfe08,
		value:    0xfc00,
		family:   Branches,
		label:    INSN_SBRC,
	},
	OpCode{
		mnemonic: "sbr",
		mask:     0xfe08,
		value:    0x6000,
		family:   Arithmetic,
		label:    INSN_SBR,
	},
	OpCode{
		mnemonic: "sub",
		mask:     0xfc00,
		value:    0x1800,
		family:   Arithmetic,
		label:    INSN_SUB,
	},
	OpCode{
		mnemonic: "dec",
		mask:     0xfe0f,
		value:    0x940a,
		family:   Arithmetic,
		label:    INSN_DEC,
	},
	OpCode{
		mnemonic: "mul",
		mask:     0xfc00,
		value:    0x9c00,
		family:   Arithmetic,
		label:    INSN_MUL,
	},
	OpCode{
		mnemonic: "sbrs",
		mask:     0xfe00,
		value:    0xfe00,
		family:   Branches,
		label:    INSN_SBRS,
	},
	// =======
	// Things that load stuff from data space or program memory into registers
	// These are tricky. the q values are interpolated into the other bits.
	// But applying the same mask as the other LD w/Z ops gives 0x8001.
	// I'm going to leave it this way until another opcode with that value comes
	// along (and hope that it won't).
	// XXX TODO: This screws up the actual mask value for this opcode. Which
	// Likely means that the way I'm doing all of these opcodes is wrong. Yay.
	// ========
	// LOAD FROM DATA SPACE
	// ========
	// LD Rd, X
	OpCode{
		mnemonic: "ldx",
		mask:     0xfe0f,
		value:    0x900c,
		family:   Transfers,
		label:    INSN_LDX,
	},
	// LD Rd, X+
	OpCode{
		mnemonic: "ldx+",
		mask:     0xfe0f,
		value:    0x900d,
		family:   Transfers,
		label:    INSN_LDXP,
	},
	// LD Rd, -X
	OpCode{
		mnemonic: "ldx-",
		mask:     0xfe0f,
		value:    0x900e,
		family:   Transfers,
		label:    INSN_LDXM,
	},
	// LD Rd, Z+
	OpCode{
		mnemonic: "ldz+",
		mask:     0xfe0f,
		value:    0x9001,
		family:   Transfers,
		label:    INSN_LDZP,
	},
	// LD Rd, -Z
	OpCode{
		mnemonic: "stz-",
		mask:     0xfe0f,
		value:    0x9002,
		family:   Transfers,
		label:    INSN_LDZM,
	},
	// LD Rd, Y+
	OpCode{
		mnemonic: "ldy+",
		mask:     0xfe0f,
		value:    0x9009,
		family:   Transfers,
		label:    INSN_LDYP,
	},
	// LD Rd, -Y
	OpCode{
		mnemonic: "ldy-",
		mask:     0xfe0f,
		value:    0x900a,
		family:   Transfers,
		label:    INSN_LDYM,
	},
	// LD Rd, Y
	// LD Rd, Z
	// Rd, Y+q
	// LDD Rd, Z+q
	OpCode{
		mnemonic: "ldd",
		mask:     0xd000,
		value:    0x8000,
		family:   Transfers,
		// Label set elsewhere
	},
	// ========
	// STORE TO DATA SPACE
	// ========
	// ST X, Rr
	OpCode{
		mnemonic: "stx",
		mask:     0xfe0f,
		value:    0x920c,
		family:   Transfers,
		label:    INSN_STX,
	},
	// ST X+, Rr
	OpCode{
		mnemonic: "stx+",
		mask:     0xfe0f,
		value:    0x920d,
		family:   Transfers,
		label:    INSN_STXP,
	},
	// ST -X, Rr
	OpCode{
		mnemonic: "stx-",
		mask:     0xfe0f,
		value:    0x920e,
		family:   Transfers,
		label:    INSN_STXM,
	},
	// ST Y+, Rr
	OpCode{
		mnemonic: "sty+",
		mask:     0xfe0f,
		value:    0x9209,
		family:   Transfers,
		label:    INSN_STYP,
	},
	// ST -Y, Rr
	OpCode{
		mnemonic: "sty-",
		mask:     0xfe0f,
		value:    0x920a,
		family:   Transfers,
		label:    INSN_STYM,
	},
	// ST Z+, Rr
	OpCode{
		mnemonic: "stz+",
		mask:     0xfe0f,
		value:    0x9201,
		family:   Transfers,
		label:    INSN_STZP,
	},
	// ST -Z, Rr
	OpCode{
		mnemonic: "stz-",
		mask:     0xfe0f,
		value:    0x9202,
		family:   Transfers,
		label:    INSN_STZM,
	},
	// ST Z, Rr
	// ST Y, Rr
	// STD Y+q, Rr
	// STD Z+q, Rr
	OpCode{
		mnemonic: "std",
		mask:     0xde00,
		value:    0x8200,
		family:   Transfers,
		// label set later
	},
	// ========
	// LOAD FROM PROGRAM MEMORY
	// ========
	// LPM Rd, Z
	OpCode{
		mnemonic: "lpmz",
		mask:     0xfe0f,
		value:    0x9004,
		family:   Transfers,
		label:    INSN_LPMZ,
	},
	// LPM Rd, Z+
	OpCode{
		mnemonic: "lpmz+",
		mask:     0xfe0f,
		value:    0x9005,
		family:   Transfers,
		label:    INSN_LPMZP,
	},
	OpCode{
		mnemonic: "lpm",
		mask:     0xffff,
		value:    0x95c8,
		family:   Transfers,
		label:    INSN_LPM,
	},
	// =======
	// END things that work with memory <-> registers
	// ======
	// 32 bit opcodes:
	OpCode{
		mnemonic: "lds",
		mask:     0xfe0f,
		value:    0x9000,
		family:   Transfers,
		label:    INSN_LDS,
	},
	OpCode{
		mnemonic: "sts",
		mask:     0xfe0f,
		value:    0x9200,
		family:   Transfers,
		label:    INSN_STS,
	},
	OpCode{
		mnemonic: "jmp",
		mask:     0xfe0e,
		value:    0x940c,
		family:   Branches,
		label:    INSN_JMP,
	},
	OpCode{
		mnemonic: "call",
		mask:     0xfe0e,
		value:    0x940e,
		family:   Branches,
		label:    INSN_CALL,
	},
}

package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/interfaces"
)

type LDR struct {
	instruction uint16
}

func (l LDR) Execute(cpu interfaces.CPU) {
	// Bits 10-8 are the destination register
	rd := uint8(l.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)

	// Bits 7-0 are the immediate value
	imm := (l.instruction & 0xFF) << 2

	fmt.Printf("ldr r%d, [pc, #0x%X]\n", rd, imm)
	memory := cpu.GetMMIO()

	address := cpu.ReadPC() + uint32(imm) - 2
	fmt.Printf("ldr r%d, [0x%X]\n", rd, address)

	read, err := memory.Read32(address)
	if err != nil {
		panic(err)
	}
	cpu.WriteRegister(rd, read)
}

type LDRR struct {
	instruction uint16
}

func (l LDRR) Execute(cpu interfaces.CPU) {
	// Bit 10 is the B bit, which determines whether this is a byte or word
	byte := l.instruction&(1<<10)>>10 == 1

	// Bits 8-6 are the offset register
	offsetRegister := l.instruction & (1<<8 | 1<<7 | 1<<6) >> 6

	// Bits 5-3 are the base register
	baseRegister := l.instruction & (1<<5 | 1<<4 | 1<<3) >> 3

	// Bits 2-0 are the destination/source register
	destinationSourceRegister := l.instruction & (1<<2 | 1<<1 | 1<<0)

	b := ""
	if byte {
		b = "b"
	}

	fmt.Printf("ldr%s r%d, [r%d, r%d]\n", b, destinationSourceRegister, baseRegister, offsetRegister)
	panic("Not implemented")
}

type STRR struct {
	instruction uint16
}

func (s STRR) Execute(cpu interfaces.CPU) {
	// Bit 10 is the B bit, which determines whether this is a byte or word
	byte := s.instruction&(1<<10)>>10 == 1

	// Bits 8-6 are the offset register
	offsetRegister := s.instruction & (1<<8 | 1<<7 | 1<<6) >> 6

	// Bits 5-3 are the base register
	baseRegister := s.instruction & (1<<5 | 1<<4 | 1<<3) >> 3

	// Bits 2-0 are the destination/source register
	destinationSourceRegister := s.instruction & (1<<2 | 1<<1 | 1<<0)

	b := ""
	if byte {
		b = "b"
	}

	fmt.Printf("str%s r%d, [r%d, r%d]\n", b, destinationSourceRegister, baseRegister, offsetRegister)

	memory := cpu.GetMMIO()

	offset := cpu.ReadRegister(uint8(offsetRegister))
	address := cpu.ReadRegister(uint8(baseRegister)) + offset
	fmt.Printf("offset=%d\n", offset)
	fmt.Printf("base address=0x%08X\n", int64(cpu.ReadRegister(uint8(baseRegister))))
	fmt.Printf("address=0x%08X\n", address)
	write := cpu.ReadRegister(uint8(destinationSourceRegister))

	if byte {
		write &= 0xFF
	}

	err := memory.Write32(uint32(address), write)
	if err != nil {
		panic(err)
	}
}

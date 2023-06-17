package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
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

	address := cpu.ReadPC() + uint32(imm)
	// Clear bit 1 of the address to ensure it's word aligned
	address &= 0xFFFFFFFC
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

type STRSP struct {
	instruction uint16
}

func (s STRSP) Execute(cpu interfaces.CPU) {
	fmt.Println("STRSP")

	// Bits 10-8 are the destination register
	rd := uint8(s.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)

	// Bits 7-0 are the immediate value
	imm := (s.instruction & 0xFF) << 2

	fmt.Printf("str r%d, [sp, #0x%X]\n", rd, imm)

	err := cpu.GetMMIO().Write32(cpu.ReadSP()+uint32(imm), cpu.ReadRegister(rd))
	if err != nil {
		panic(err)
	}

	cpu.SetN(cpu.ReadRegister(rd)&0x80000000 != 0)
}

type LDRSP struct {
	instruction uint16
}

func (l LDRSP) Execute(cpu interfaces.CPU) {
	fmt.Println("LDRSP")

	// Bits 10-8 are the destination register
	rd := uint8(l.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)

	// Bits 7-0 are the immediate value
	imm := (l.instruction & 0xFF) << 2

	fmt.Printf("ldr r%d, [sp, #0x%X]\n", rd, imm)

	mem, err := cpu.GetMMIO().Read32(cpu.ReadSP() + uint32(imm))
	if err != nil {
		panic(err)
	}
	cpu.WriteRegister(rd, mem)
}

type LDRH struct {
	instruction uint16
}

func (l LDRH) Execute(cpu interfaces.CPU) {
	fmt.Println("LDRH")

	// Bits 10-6 are the offset
	offset := l.instruction & (1<<10 | 1<<9 | 1<<8 | 1<<7 | 1<<6) >> 6

	// Bits 5-3 are the base register
	rb := uint8(l.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination/source register
	rd := uint8(l.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("ldrh r%d, [r%d, #0x%X]\n", rd, rb, offset)

	// Load the lower 16 bits of the address at rb + offset into rd
	readHW, err := cpu.GetMMIO().Read16(cpu.ReadRegister(rb) + uint32(offset))
	if err != nil {
		panic(err)
	}

	cpu.WriteRegister(rd, uint32(readHW))
}

type STRH struct {
	instruction uint16
}

func (s STRH) Execute(cpu interfaces.CPU) {
	fmt.Println("STRH")

	// Bits 10-6 are the offset
	offset := uint32(s.instruction & (1<<10 | 1<<9 | 1<<8 | 1<<7 | 1<<6) >> 6)

	// Bits 5-3 are the base register
	rb := uint8(s.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination/source register
	rd := uint8(s.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("strh r%d, [r%d, #0x%X]\n", rd, rb, offset)

	// Store the lower 16 bits of the rd into the address at rb + offset
	cpu.GetMMIO().Write16(cpu.ReadRegister(rb)+offset, uint16(cpu.ReadRegister(rd)&0xFFFF))
}

type LDRBImm struct {
	instruction uint16
}

func (l LDRBImm) Execute(cpu interfaces.CPU) {
	fmt.Println("LDRBImm")

	// Bits 10-6 are the offset
	offset := uint32(l.instruction & (1<<10 | 1<<9 | 1<<8 | 1<<7 | 1<<6) >> 6)

	// Bits 5-3 are the base register
	rb := uint8(l.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination/source register
	rd := uint8(l.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("ldrb r%d, [r%d, #0x%X]\n", rd, rb, offset)

	// Load the byte at rb + offset into rd
	readByte, err := cpu.GetMMIO().Read8(cpu.ReadRegister(rb) + offset)
	if err != nil {
		panic(err)
	}
	cpu.WriteRegister(rd, uint32(readByte))
}

type STRBImm struct {
	instruction uint16
}

func (s STRBImm) Execute(cpu interfaces.CPU) {
	fmt.Println("STRBImm")

	// Bits 10-6 are the offset
	offset := uint32(s.instruction & (1<<10 | 1<<9 | 1<<8 | 1<<7 | 1<<6) >> 6)

	// Bits 5-3 are the base register
	rb := uint8(s.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination/source register
	rd := uint8(s.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("strb r%d, [r%d, #0x%X]\n", rd, rb, offset)

	// Store the byte in rd into the address at rb + offset
	cpu.GetMMIO().Write8(cpu.ReadRegister(rb)+offset, uint8(cpu.ReadRegister(rd)&0xFF))
}

type LDRWImm struct {
	instruction uint16
}

func (l LDRWImm) Execute(cpu interfaces.CPU) {
	fmt.Println("LDRWImm")

	// Bits 10-6 are the offset
	offset := uint32(l.instruction & (1<<10 | 1<<9 | 1<<8 | 1<<7 | 1<<6) >> 6)

	// Bits 5-3 are the base register
	rb := uint8(l.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination/source register
	rd := uint8(l.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("ldr r%d, [r%d, #0x%X]\n", rd, rb, offset)

	// Load the word at rb + offset into rd
	mem, err := cpu.GetMMIO().Read32(cpu.ReadRegister(rb) + offset)
	if err != nil {
		panic(err)
	}

	cpu.WriteRegister(rd, mem)
}

type STRWImm struct {
	instruction uint16
}

func (s STRWImm) Execute(cpu interfaces.CPU) {
	fmt.Println("STRWImm")

	// Bits 10-6 are the offset
	offset := uint32(s.instruction & (1<<10 | 1<<9 | 1<<8 | 1<<7 | 1<<6) >> 6)

	// Bits 5-3 are the base register
	rb := uint8(s.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination/source register
	rd := uint8(s.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("str r%d, [r%d, #0x%X]\n", rd, rb, offset)

	// Store the word in rd into the address at rb + offset
	cpu.GetMMIO().Write32(cpu.ReadRegister(rb)+offset, cpu.ReadRegister(rd))
}

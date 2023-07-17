package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type LDR struct {
	instruction uint16
}

func (l LDR) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
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
	return
}

type LDRR struct {
	instruction uint16
}

func (l LDRR) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
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
	return
}

type STRR struct {
	instruction uint16
}

func (s STRR) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
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
	return
}

type STRSP struct {
	instruction uint16
}

func (s STRSP) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
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

	cpu.SetN(cpu.ReadRegister(rd)&(1<<31)>>31 != 0)
	return
}

type LDRSP struct {
	instruction uint16
}

func (l LDRSP) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
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
	return
}

type LDRH struct {
	instruction uint16
}

func (l LDRH) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("LDRH")

	// Bits 10-6 are the offset
	offset := (l.instruction & (1<<10 | 1<<9 | 1<<8 | 1<<7 | 1<<6) >> 6) << 1

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
	return
}

type STRH struct {
	instruction uint16
}

func (s STRH) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("STRH")

	// Bits 10-6 are the offset
	offset := uint32(s.instruction&(1<<10|1<<9|1<<8|1<<7|1<<6)>>6) << 1

	// Bits 5-3 are the base register
	rb := uint8(s.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination/source register
	rd := uint8(s.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("strh r%d, [r%d, #0x%X]\n", rd, rb, offset)

	// Store the lower 16 bits of the rd into the address at rb + offset
	err := cpu.GetMMIO().Write16(cpu.ReadRegister(rb)+offset, uint16(cpu.ReadRegister(rd)&0xFFFF))
	if err != nil {
		panic(err)
	}
	return
}

type LDRBImm struct {
	instruction uint16
}

func (l LDRBImm) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("LDRBImm")

	// Bits 10-6 are the offset
	offset := uint32(l.instruction&(1<<10|1<<9|1<<8|1<<7|1<<6)>>6) << 1

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
	return
}

type STRBImm struct {
	instruction uint16
}

func (s STRBImm) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("STRBImm")

	// Bits 10-6 are the offset
	offset := uint32(s.instruction & (1<<10 | 1<<9 | 1<<8 | 1<<7 | 1<<6) >> 6)

	// Bits 5-3 are the base register
	rb := uint8(s.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination/source register
	rd := uint8(s.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("strb r%d, [r%d, #0x%X]\n", rd, rb, offset)

	// Store the byte in rd into the address at rb + offset
	err := cpu.GetMMIO().Write8(cpu.ReadRegister(rb)+offset, uint8(cpu.ReadRegister(rd)&0xFF))
	if err != nil {
		panic(err)
	}
	return
}

type LDRWImm struct {
	instruction uint16
}

func (l LDRWImm) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("LDRWImm")

	// Bits 10-6 are the offset
	offset := uint32(l.instruction&(1<<10|1<<9|1<<8|1<<7|1<<6)>>6) << 1

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
	return
}

type STRWImm struct {
	instruction uint16
}

func (s STRWImm) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("STRWImm")

	// Bits 10-6 are the offset
	offset := uint32(s.instruction&(1<<10|1<<9|1<<8|1<<7|1<<6)>>6) << 2

	// Bits 5-3 are the base register
	rb := uint8(s.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination/source register
	rd := uint8(s.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("str r%d, [r%d, #0x%X]\n", rd, rb, offset)

	// Store the word in rd into the address at rb + offset
	err := cpu.GetMMIO().Write32(cpu.ReadRegister(rb)+offset, cpu.ReadRegister(rd))
	if err != nil {
		panic(err)
	}
	return
}

type LDMIA struct {
	instruction uint16
}

func (l LDMIA) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bits 10-8 are the base register
	rb := uint8(l.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)

	// Bits 7-0 are the register list
	registerList := uint16(l.instruction & (1<<7 | 1<<6 | 1<<5 | 1<<4 | 1<<3 | 1<<2 | 1<<1 | 1<<0))

	var popRegisters []uint8

	// Collect the registers to push in backwards order so that they are pushed in the correct order
	for i := 0; i < 8; i++ {
		if registerList&(1<<i)>>i == 1 {
			popRegisters = append(popRegisters, uint8(i))
		}
	}

	fmt.Printf("ldmia r%d!, {%v}\n", rb, popRegisters)

	address := cpu.ReadRegister(rb)
	for _, register := range popRegisters {
		// Load the word at address into register
		fmt.Printf("Loading word at 0x%X into r%d\n", address, register)
		mem, err := cpu.GetMMIO().Read32(address)
		if err != nil {
			panic(err)
		}
		address += 4
		cpu.WriteRegister(register, mem)
	}

	// Write the new address back to rb
	cpu.WriteRegister(rb, address)

	return
}

type STMIA struct {
	instruction uint16
}

func (s STMIA) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bits 10-8 are the base register
	rb := uint8(s.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)

	// Bits 7-0 are the register list
	registerList := uint16(s.instruction & (1<<7 | 1<<6 | 1<<5 | 1<<4 | 1<<3 | 1<<2 | 1<<1 | 1<<0))

	var pushRegisters []uint8

	// Collect the registers to push in backwards order so that they are pushed in the correct order
	for i := 0; i < 8; i++ {
		if registerList&(1<<i)>>i == 1 {
			pushRegisters = append(pushRegisters, uint8(i))
		}
	}

	fmt.Printf("stmia r%d!, {%v}\n", rb, pushRegisters)

	address := cpu.ReadRegister(rb)

	for _, register := range pushRegisters {
		// Store the word in register into address
		fmt.Printf("Storing word in r%d into 0x%X\n", register, address)
		err := cpu.GetMMIO().Write32(address, cpu.ReadRegister(register))
		if err != nil {
			panic(err)
		}
		address += 4
	}

	// Write the new address back to rb
	cpu.WriteRegister(rb, address)

	return
}

type STRNSH struct {
	instruction uint16
}

// strh unsigned offset
func (s STRNSH) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Store halfword:
	// Add Ro to base address in Rb. Store bits 0-
	// 15 of Rd at the resulting address

	// Bits 8-6 are the offset register
	ro := uint8(s.instruction & (1<<8 | 1<<7 | 1<<6) >> 6)

	// Bits 5-3 are the base register
	rb := uint8(s.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination/source register
	rd := uint8(s.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("strh r%d, [r%d, r%d]\n", rd, rb, ro)

	// Store the halfword in rd into the address at rb + ro
	err := cpu.GetMMIO().Write16(cpu.ReadRegister(rb)+cpu.ReadRegister(ro), uint16(cpu.ReadRegister(rd)))
	if err != nil {
		panic(err)
	}

	return
}

type LDRNSH struct {
	instruction uint16
}

// ldrh unsigned offset
func (l LDRNSH) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Load halfword:
	// Add Ro to base address in Rb. Load bits 0-
	// 15 of Rd from the resulting address, and set
	// bits 16-31 of Rd to 0.

	// Bits 8-6 are the offset register
	ro := uint8(l.instruction & (1<<8 | 1<<7 | 1<<6) >> 6)

	// Bits 5-3 are the base register
	rb := uint8(l.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination/source register
	rd := uint8(l.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("ldrh r%d, [r%d, r%d]\n", rd, rb, ro)

	// Load the halfword at rb + ro into rd
	mem, err := cpu.GetMMIO().Read16(cpu.ReadRegister(rb) + cpu.ReadRegister(ro))
	if err != nil {
		panic(err)
	}

	// Set the upper 16 bits of rd to 0
	cpu.WriteRegister(rd, uint32(mem))

	return
}

type LDSB struct {
	instrument uint16
}

// ldrsb signed offset
func (l LDSB) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	panic("LDSB not implemented")
}

type LDSH struct {
	instruction uint16
}

// ldrsh signed offset
func (l LDSH) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	panic("LDSH not implemented")
}

package arm

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type LDR struct {
	instruction uint32
}

func (ldr LDR) Execute(cpu interfaces.CPU) (repipeline bool) {
	immediate := ldr.instruction&(1<<25)>>25 == 0
	pre := ldr.instruction&(1<<24)>>24 == 1
	up := ldr.instruction&(1<<23)>>23 == 1
	word := ldr.instruction&(1<<22)>>22 == 0
	writeback := ldr.instruction&(1<<21)>>21 == 1

	rn := uint8((ldr.instruction >> 16) & 0xF)
	rd := uint8((ldr.instruction >> 12) & 0xF)

	memory := cpu.GetMMIO()

	fmt.Printf("Immediate: %t, Pre: %t, Up: %t, Word: %t, Writeback: %t\n", immediate, pre, up, word, writeback)

	offset := uint32(0)
	if immediate {
		offset = ldr.instruction & 0xFFF
	} else {
		offset, _ = unshiftRegister(ldr.instruction&0xFFF, cpu)
	}
	fmt.Printf("LDR r%d, [r%d, 0x%X]\n", rd, rn, offset)

	address := cpu.ReadRegister(rn)
	if pre {
		if up {
			address += offset
		} else {
			address -= offset
		}
	}

	// If the datatype is a byte, check if it's on a word boundary
	if !word && address%4 == 0 {
		// On a word boundary, bits 7-0 of the value in memory at address
		// are moved to the bottom bits of the result register
		read, err := memory.Read32(address)
		if err != nil {
			panic(err)
		}
		cpu.WriteRegister(rd, read&0xFF)
	} else if !word && address%4 == 1 {
		// On a word boundary plus one, bits 15-8 of the value in memory at address
		// are moved to the bottom bits of the result register
		read, err := memory.Read32(address)
		if err != nil {
			panic(err)
		}
		cpu.WriteRegister(rd, (read&0xFF00)>>8)
	} else if !word && address%4 == 2 {
		// On a word boundary plus two, bits 23-16 of the value in memory at address
		// are moved to the bottom bits of the result register
		read, err := memory.Read32(address)
		if err != nil {
			panic(err)
		}
		cpu.WriteRegister(rd, (read&0xFF0000)>>16)
	} else if !word && address%4 == 3 {
		// On a word boundary plus three, bits 31-24 of the value in memory at address
		// are moved to the bottom bits of the result register
		read, err := memory.Read32(address)
		if err != nil {
			panic(err)
		}
		cpu.WriteRegister(rd, (read&0xFF000000)>>24)
	} else if word && address%4 != 0 {
		panic("Unaligned word access")
	} else {
		// Otherwise, the value in memory at address is moved to the result register
		read, err := memory.Read32(address)
		if err != nil {
			panic(err)
		}
		cpu.WriteRegister(rd, read)
	}
	if !pre || writeback {
		cpu.WriteRegister(rn, address)
	}
	fmt.Printf("Address: 0x%X\n", address)
	return
}

type STR struct {
	instruction uint32
}

func (str STR) Execute(cpu interfaces.CPU) (repipeline bool) {
	immediate := str.instruction&(1<<25)>>25 == 0
	pre := str.instruction&(1<<24)>>24 == 1
	up := str.instruction&(1<<23)>>23 == 1
	word := str.instruction&(1<<22)>>22 == 0
	writeback := str.instruction&(1<<21)>>21 == 1

	rn := uint8((str.instruction >> 16) & 0xF)
	rd := uint8((str.instruction >> 12) & 0xF)

	memory := cpu.GetMMIO()

	fmt.Printf("Immediate: %t, Pre: %t, Up: %t, Word: %t, Writeback: %t\n", immediate, pre, up, word, writeback)

	offset := uint32(0)
	if immediate {
		offset = str.instruction & 0xFFF
	} else {
		offset, _ = unshiftRegister(str.instruction&0xFFF, cpu)
	}
	fmt.Printf("STR r%d, [r%d, 0x%X]\n", rd, rn, offset)

	address := cpu.ReadRegister(rn)
	if pre {
		if up {
			address += offset
		} else {
			address -= offset
		}
	}

	existingVal, err := memory.Read32(address)
	if err != nil {
		panic(err)
	}

	// If the datatype is a byte, check if it's on a word boundary
	if !word && address%4 == 0 {
		// On a word boundary, bits 7-0 of the value in the register
		// are moved to the bottom bits of the memory at address
		read := cpu.ReadRegister(rd)
		memory.Write32(address, (existingVal&0xFFFFFF00)|read)
	} else if !word && address%4 == 1 {
		// On a word boundary plus one, bits 15-8 of the value in the register
		// are moved to the bottom bits of the memory at address
		read := cpu.ReadRegister(rd)
		memory.Write32(address, (existingVal&0xFFFF00FF)|(read<<8))
	} else if !word && address%4 == 2 {
		// On a word boundary plus two, bits 23-16 of the value in the register
		// are moved to the bottom bits of the memory at address
		read := cpu.ReadRegister(rd)
		memory.Write32(address, (existingVal&0xFF00FFFF)|(read<<16))
	} else if !word && address%4 == 3 {
		// On a word boundary plus three, bits 31-24 of the value in the register
		// are moved to the bottom bits of the memory at address
		read := cpu.ReadRegister(rd)
		memory.Write32(address, (existingVal&0x00FFFFFF)|(read<<24))
	} else if word && address%4 != 0 {
		panic("Unaligned word access")
	} else {
		// Otherwise, the value in the register is moved to the memory at address
		read := cpu.ReadRegister(rd)
		memory.Write32(address, read)
	}
	if !pre || writeback {
		cpu.WriteRegister(rn, address)
	}
	fmt.Printf("Address: 0x%X\n", address)
	return
}

type LDM struct {
	instruction uint32
}

func (ldm LDM) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("LDM")

	// Bit 24 == 1 means pre-indexed addressing
	pre := ldm.instruction&(1<<24)>>24 == 1
	// Bit 23 == 1 means the offset is added to the base register (up)
	up := ldm.instruction&(1<<23)>>23 == 1
	// Bit 22 == 1 means to load the PSR or force user mode
	psr := ldm.instruction&(1<<22)>>22 == 1
	// Bit 21 == 1 means the base register is written back to
	writeback := ldm.instruction&(1<<21)>>21 == 1

	fmt.Printf("Pre: %t, Up: %t, PSR: %t, Writeback: %t\n", pre, up, psr, writeback)

	// Bits 19-16 are the base register
	rn := uint8((ldm.instruction >> 16) & 0xF)

	// Bits 15-0 are the register list
	registerList := ldm.instruction & 0xFFFF

	address := cpu.ReadRegister(rn)

	// If the PSR bit is set, we need to load the PSR
	if psr {
		panic("Not implemented")
	}

	var registers []uint8
	for i := uint8(0); i < 16; i++ {
		if registerList&(1<<i)>>i == 1 {
			registers = append(registers, i)
		}
	}

	for _, register := range registers {
		if pre {
			if up {
				address += 4
			} else {
				address -= 4
			}
		}
		value, err := cpu.GetMMIO().Read32(address)
		if err != nil {
			panic(err)
		}
		cpu.WriteRegister(register, value)
		if !pre {
			if up {
				address += 4
			} else {
				address -= 4
			}
		}
	}

	if !pre || writeback {
		cpu.WriteRegister(rn, address)
	}
	return
}

type STM struct {
	instruction uint32
}

func (stm STM) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("STM")

	// Bit 24 == 1 means pre-indexed addressing
	pre := stm.instruction&(1<<24)>>24 == 1
	// Bit 23 == 1 means the offset is added to the base register (up)
	up := stm.instruction&(1<<23)>>23 == 1
	// Bit 22 == 1 means to load the PSR or force user mode
	psr := stm.instruction&(1<<22)>>22 == 1
	// Bit 21 == 1 means the base register is written back to
	writeback := stm.instruction&(1<<21)>>21 == 1

	fmt.Printf("Pre: %t, Up: %t, PSR: %t, Writeback: %t\n", pre, up, psr, writeback)

	// Bits 19-16 are the base register
	rn := uint8((stm.instruction >> 16) & 0xF)

	// Bits 15-0 are the register list
	registerList := stm.instruction & 0xFFFF

	fmt.Printf("stm r%d, 0x%X\n", rn, registerList)

	address := cpu.ReadRegister(rn)

	// If the PSR bit is set, we need to store the PSR
	if psr {
		panic("Not implemented")
	}

	var pushRegisters []uint8

	// Collect the registers to push in backwards order so that they are pushed in the correct order
	for i := 15; i >= 0; i-- {
		if registerList&(1<<i)>>i == 1 {
			pushRegisters = append(pushRegisters, uint8(i))
		}
	}

	for _, reg := range pushRegisters {
		// If pre-indexed addressing is used, the base register is updated
		if pre {
			if up {
				address += 4
			} else {
				address -= 4
			}
		}
		err := cpu.GetMMIO().Write32(address, cpu.ReadRegister(reg))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Pushing register r%d @ %08X\n", reg, address)
		if !pre {
			if up {
				address += 4
			} else {
				address -= 4
			}
		}
	}

	if !pre || writeback {
		cpu.WriteRegister(rn, address)
	}

	return
}

type LDRSH struct {
	instruction uint32
}

func (ldrsh LDRSH) Execute(cpu interfaces.CPU) (repipeline bool) {
	// // Bit 24 == 1 means pre-indexed addressing
	// pre := ldrsh.instruction&(1<<24)>>24 == 1
	// // Bit 23 == 1 means the offset is added to the base register (up)
	// up := ldrsh.instruction&(1<<23)>>23 == 1
	// // Bit 21 == 1 means the base register is written back to
	// writeback := ldrsh.instruction&(1<<21)>>21 == 1

	// // Bits 19-16 are the base register
	// rn := uint8((ldrsh.instruction >> 16) & 0xF)

	// // Bits 15-12 are the destination register
	// rd := uint8((ldrsh.instruction >> 12) & 0xF)

	// // Bits 11-8 are the offset's high nibble
	// offsetHigh := uint8((ldrsh.instruction >> 8) & 0xF)

	// // Bits 3-0 are the offset's low nibble
	// offsetLow := uint8(ldrsh.instruction & 0xF)

	panic("LDRSH Not implemented")
}

type LDRSB struct {
	instruction uint32
}

func (ldrsb LDRSB) Execute(cpu interfaces.CPU) (repipeline bool) {
	// // Bit 24 == 1 means pre-indexed addressing
	// pre := ldrsb.instruction&(1<<24)>>24 == 1
	// // Bit 23 == 1 means the offset is added to the base register (up)
	// up := ldrsb.instruction&(1<<23)>>23 == 1
	// // Bit 21 == 1 means the base register is written back to
	// writeback := ldrsb.instruction&(1<<21)>>21 == 1

	// // Bits 19-16 are the base register
	// rn := uint8((ldrsb.instruction >> 16) & 0xF)

	// // Bits 15-12 are the destination register
	// rd := uint8((ldrsb.instruction >> 12) & 0xF)

	// // Bits 11-8 are the offset's high nibble
	// offsetHigh := uint8((ldrsb.instruction >> 8) & 0xF)

	// // Bits 3-0 are the offset's low nibble
	// offsetLow := uint8(ldrsb.instruction & 0xF)
	panic("LDRSB Not implemented")
}

type LDRH struct {
	instruction uint32
}

func (ldrh LDRH) Execute(cpu interfaces.CPU) (repipeline bool) {
	// Bit 24 == 1 means pre-indexed addressing
	pre := ldrh.instruction&(1<<24)>>24 == 1
	// Bit 23 == 1 means the offset is added to the base register (up)
	up := ldrh.instruction&(1<<23)>>23 == 1
	// Bit 21 == 1 means the base register is written back to
	writeback := ldrh.instruction&(1<<21)>>21 == 1

	// Bits 19-16 are the base register
	rn := uint8((ldrh.instruction >> 16) & 0xF)

	// Bits 15-12 are the destination register
	rd := uint8((ldrh.instruction >> 12) & 0xF)

	// Bits 11-8 are the offset's high nibble
	offsetHigh := uint8((ldrh.instruction >> 8) & 0xF)

	// Bits 3-0 are the offset's low nibble
	offsetLow := uint8(ldrh.instruction & 0xF)

	address := cpu.ReadRegister(rn)
	if pre {
		if up {
			address += uint32(offsetHigh<<4 | offsetLow)
		} else {
			address -= uint32(offsetHigh<<4 | offsetLow)
		}
	}

	fmt.Printf("ldrh r%d, [r%d, #%d]\n", rd, rn, offsetHigh<<4|offsetLow)

	// Load halfword from memory
	halfword, err := cpu.GetMMIO().Read16(address)
	if err != nil {
		panic(err)
	}

	cpu.WriteRegister(rd, uint32(halfword))

	if !pre || writeback {
		if up {
			address += uint32(offsetHigh<<4 | offsetLow)
		} else {
			address -= uint32(offsetHigh<<4 | offsetLow)
		}
		cpu.WriteRegister(rn, address)
	}

	return
}

type STRSH struct {
	instruction uint32
}

func (strsh STRSH) Execute(cpu interfaces.CPU) (repipeline bool) {
	// // Bit 24 == 1 means pre-indexed addressing
	// pre := strsh.instruction&(1<<24)>>24 == 1
	// // Bit 23 == 1 means the offset is added to the base register (up)
	// up := strsh.instruction&(1<<23)>>23 == 1
	// // Bit 21 == 1 means the base register is written back to
	// writeback := strsh.instruction&(1<<21)>>21 == 1

	// // Bits 19-16 are the base register
	// rn := uint8((strsh.instruction >> 16) & 0xF)

	// // Bits 15-12 are the destination register
	// rd := uint8((strsh.instruction >> 12) & 0xF)

	// // Bits 11-8 are the offset's high nibble
	// offsetHigh := uint8((strsh.instruction >> 8) & 0xF)

	// // Bits 3-0 are the offset's low nibble
	// offsetLow := uint8(strsh.instruction & 0xF)
	panic("STRSH Not implemented")
}

type STRSB struct {
	instruction uint32
}

func (strsb STRSB) Execute(cpu interfaces.CPU) (repipeline bool) {
	// // Bit 24 == 1 means pre-indexed addressing
	// pre := strsb.instruction&(1<<24)>>24 == 1
	// // Bit 23 == 1 means the offset is added to the base register (up)
	// up := strsb.instruction&(1<<23)>>23 == 1
	// // Bit 21 == 1 means the base register is written back to
	// writeback := strsb.instruction&(1<<21)>>21 == 1

	// // Bits 19-16 are the base register
	// rn := uint8((strsb.instruction >> 16) & 0xF)

	// // Bits 15-12 are the destination register
	// rd := uint8((strsb.instruction >> 12) & 0xF)

	// // Bits 11-8 are the offset's high nibble
	// offsetHigh := uint8((strsb.instruction >> 8) & 0xF)

	// // Bits 3-0 are the offset's low nibble
	// offsetLow := uint8(strsb.instruction & 0xF)
	panic("STRSB Not implemented")
}

type STRH struct {
	instruction uint32
}

func (strh STRH) Execute(cpu interfaces.CPU) (repipeline bool) {
	// Bit 24 == 1 means pre-indexed addressing
	pre := strh.instruction&(1<<24)>>24 == 1
	// Bit 23 == 1 means the offset is added to the base register (up)
	up := strh.instruction&(1<<23)>>23 == 1
	// Bit 21 == 1 means the base register is written back to
	writeback := strh.instruction&(1<<21)>>21 == 1

	// Bits 19-16 are the base register
	rn := uint8((strh.instruction >> 16) & 0xF)

	// Bits 15-12 are the destination register
	rd := uint8((strh.instruction >> 12) & 0xF)

	// Bits 11-8 are the offset's high nibble
	offsetHigh := uint8((strh.instruction >> 8) & 0xF)

	// Bits 3-0 are the offset's low nibble
	offsetLow := uint8(strh.instruction & 0xF)

	offset := uint32(offsetHigh)<<4 | uint32(offsetLow)

	address := cpu.ReadRegister(rn)
	if pre {
		if up {
			address += offset
		} else {
			address -= offset
		}
	}

	fmt.Printf("strh r%d, [r%d, #0x%X]\n", rd, rn, offset)

	// Store unsigned halfword
	cpu.GetMMIO().Write16(address, uint16(cpu.ReadRegister(rd)))

	if !pre {
		if up {
			address += offset
		} else {
			address -= offset
		}
	}

	if writeback || !pre {
		cpu.WriteRegister(rn, address)
	}

	return
}

type LDRSHRegisterOffset struct {
	instruction uint32
}

func (ldrsh LDRSHRegisterOffset) Execute(cpu interfaces.CPU) (repipeline bool) {
	panic("LDRSHRegisterOffset Not implemented")
}

type LDRSBRegisterOffset struct {
	instruction uint32
}

func (ldrsb LDRSBRegisterOffset) Execute(cpu interfaces.CPU) (repipeline bool) {
	panic("LDRSBRegisterOffset Not implemented")
}

type LDRHRegisterOffset struct {
	instruction uint32
}

func (ldrh LDRHRegisterOffset) Execute(cpu interfaces.CPU) (repipeline bool) {
	panic("LDRHRegisterOffset Not implemented")
}

type STRSHRegisterOffset struct {
	instruction uint32
}

func (strsh STRSHRegisterOffset) Execute(cpu interfaces.CPU) (repipeline bool) {
	panic("STRSHRegisterOffset Not implemented")
}

type STRSBRegisterOffset struct {
	instruction uint32
}

func (strsb STRSBRegisterOffset) Execute(cpu interfaces.CPU) (repipeline bool) {
	panic("STRSBRegisterOffset Not implemented")
}

type STRHRegisterOffset struct {
	instruction uint32
}

func (strh STRHRegisterOffset) Execute(cpu interfaces.CPU) (repipeline bool) {
	// Bit 24 == 1 means pre-indexed addressing
	pre := strh.instruction&(1<<24)>>24 == 1
	// Bit 23 == 1 means the offset is added to the base register (up)
	up := strh.instruction&(1<<23)>>23 == 1
	// Bit 21 == 1 means the base register is written back to
	writeback := strh.instruction&(1<<21)>>21 == 1

	// Bits 19-16 are the base register
	rn := uint8((strh.instruction >> 16) & 0xF)

	// Bits 15-12 are the destination register
	rd := uint8((strh.instruction >> 12) & 0xF)

	// Bits 3-0 are the offset register
	rm := uint8(strh.instruction & 0xF)

	offset := cpu.ReadRegister(rm)

	address := cpu.ReadRegister(rn)
	if pre {
		if up {
			address += offset
		} else {
			address -= offset
		}
	}

	fmt.Printf("strh r%d, [r%d, r%d]  # 0x%08x\n", rd, rn, rm, address)

	// Store unsigned halfword
	cpu.GetMMIO().Write16(address, uint16(cpu.ReadRegister(rd)))

	if !pre {
		if up {
			address += offset
		} else {
			address -= offset
		}
	}

	if writeback || !pre {
		cpu.WriteRegister(rn, address)
	}

	return
}

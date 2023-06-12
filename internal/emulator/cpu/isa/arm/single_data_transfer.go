package arm

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type LDR struct {
	instruction uint32
}

func (ldr LDR) Execute(cpu interfaces.CPU) {
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
}

type STR struct {
	instruction uint32
}

func (str STR) Execute(cpu interfaces.CPU) {
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
}

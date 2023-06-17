package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type B struct {
	instruction uint16
}

func (a B) Execute(cpu interfaces.CPU) {
	fmt.Println("B")
	fmt.Println("ConditionalBranch")
	// Bits 11-8 are the condition
	cond := a.instruction & (1<<11 | 1<<10 | 1<<9 | 1<<8) >> 8

	// Bits 7-0 are the 8-bit signed offset
	offset := int8(a.instruction & 0xFF)

	conditionPassed := false
	switch cond {
	case 0b0:
		conditionPassed = cpu.GetZ()
	case 0b1:
		conditionPassed = !cpu.GetZ()
	case 0b10:
		conditionPassed = cpu.GetC()
	case 0b11:
		conditionPassed = !cpu.GetC()
	case 0b100:
		conditionPassed = cpu.GetN()
	case 0b101:
		conditionPassed = !cpu.GetN()
	case 0b110:
		conditionPassed = cpu.GetV()
	case 0b111:
		conditionPassed = !cpu.GetV()
	case 0b1000:
		conditionPassed = cpu.GetC() && !cpu.GetZ()
	case 0b1001:
		conditionPassed = !cpu.GetC() || cpu.GetZ()
	case 0b1010:
		conditionPassed = cpu.GetN() == cpu.GetV()
	case 0b1011:
		conditionPassed = cpu.GetN() != cpu.GetV()
	case 0b1100:
		conditionPassed = !cpu.GetZ() && cpu.GetN() == cpu.GetV()
	case 0b1101:
		conditionPassed = cpu.GetZ() || cpu.GetN() != cpu.GetV()
	}

	if conditionPassed {
		cpu.WritePC(cpu.ReadPC() + uint32(offset))
	} else {
		fmt.Println("Branch condition not met")
	}
}

type BX struct {
	instruction uint16
}

func (a BX) Execute(cpu interfaces.CPU) {
	fmt.Println("BX")

	// Bits 7-6 are the hi operand flags
	hof := a.instruction & (1<<7 | 1<<6) >> 6
	rshs := uint8(a.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	if hof == 0b00 {
		fmt.Printf("Branching to r%d\n", rshs)
		addr := cpu.ReadRegister(rshs)
		fmt.Printf("Address: 0x%08X\n", addr)
		if addr&0b11 != 0 {
			fmt.Println("Setting THUMB mode")
			cpu.WritePC(addr - 1)
			cpu.SetThumbMode(true)
		} else {
			cpu.WritePC(addr)
			cpu.SetThumbMode(false)
		}
	} else if hof == 0b01 {
		fmt.Printf("Branching to r%d\n", rshs+8)
		addr := cpu.ReadRegister(rshs + 8)
		fmt.Printf("Address: 0x%08X\n", addr)
		if addr&0b11 != 0 {
			fmt.Println("Setting THUMB mode")
			cpu.WritePC(addr - 1)
			cpu.SetThumbMode(true)
		} else {
			cpu.WritePC(addr)
			cpu.SetThumbMode(false)
		}
	}
}

type LBL struct {
	instruction uint16
}

func (a LBL) Execute(cpu interfaces.CPU) {
	fmt.Println("LBL")

	// Bit 11 == 1 is low offset
	low := a.instruction&(1<<11)>>11 == 1

	// Bits 10-0 are the offset
	offset := uint16(a.instruction & 0x7FF)

	if low {
		offset = offset << 1
		fmt.Println("bl Low offset")
		// Take the LR
		lr := cpu.ReadLR()
		// Write the current PC to the LR
		cpu.WriteLR(cpu.ReadPC())
		// Add the offset to the LR and store it in the PC
		newPC := lr + uint32(offset)
		cpu.WritePC(newPC + 2)
		// Set bit 0 of the LR to 1
		cpu.WriteLR(cpu.ReadLR() | 1)
	} else {
		offset = offset << 12
		signedOffset := int16(offset)
		fmt.Println("bl High offset")
		fmt.Printf("Adding %d to PC (%08X)\n", signedOffset, cpu.ReadPC())
		newPC := cpu.ReadPC() + uint32(signedOffset)
		// Add the offset to the PC and store it in LR
		cpu.WriteLR(newPC)
	}
}

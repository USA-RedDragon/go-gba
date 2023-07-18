package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type UnconditionalBranch struct {
	instruction uint16
}

func (u UnconditionalBranch) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("UnconditionalBranch")

	// Bits 10-0 are the offset
	offset := int32(u.instruction & 0b11111111111)
	offset <<= 21
	offset >>= 20

	fmt.Printf("Offset: %d\n", offset)
	fmt.Printf("Address: 0x%08X\n", cpu.ReadPC()+uint32(offset))
	cpu.WritePC(cpu.ReadPC() + uint32(offset))
	return
}

type B struct {
	instruction uint16
}

func (a B) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("B")
	fmt.Println("ConditionalBranch")
	// Bits 11-8 are the condition
	cond := a.instruction & (1<<11 | 1<<10 | 1<<9 | 1<<8) >> 8

	fmt.Printf("Condition: 0b%04b\n", cond)

	// Bits 7-0 are the 8-bit signed offset
	offset := int8(a.instruction&0xFF) << 1

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
		if offset == 0 {
			repipeline = true
			return
		} else {
			cpu.WritePC(cpu.ReadPC() + uint32(offset))
		}
	} else {
		fmt.Println("Branch condition not met")
	}
	return
}

type BX struct {
	instruction uint16
}

func (a BX) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("BX")

	// Bits 7-6 are the hi operand flags
	rshs := uint8(a.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	fmt.Printf("BLUG bx r%d\n", rshs)

	addr := cpu.ReadRegister(rshs)

	if addr&1 != 1 {
		fmt.Println("ARM mode")
		cpu.SetThumbMode(false)
	}

	cpu.WritePC(addr)

	repipeline = true
	return
}

type LBL struct {
	instruction uint16
}

func (a LBL) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
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
		cpu.WritePC(newPC)
		// Set bit 0 of the LR to 1
		cpu.WriteLR((cpu.ReadLR() | 1) - 2)
	} else {
		offset = offset << 12
		signedOffset := int16(offset)
		fmt.Println("bl High offset")
		fmt.Printf("Adding %d to PC (%08X)\n", signedOffset, cpu.ReadPC())
		newPC := cpu.ReadPC() + uint32(signedOffset)
		// Add the offset to the PC and store it in LR
		cpu.WriteLR(newPC)
	}
	return
}

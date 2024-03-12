package arm

import (
	"fmt"
	"math/bits"

	"github.com/USA-RedDragon/go-gba/internal/emulator/cpu/isa"
	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

// https://iitd-plos.github.io/col718/ref/arm-instructionset.pdf, Figure 4-1
const (
	DataProcessingMask   = 0b0000_1100_0000_0000_0000_0000_0000_0000
	DataProcessingFormat = 0b0000_0000_0000_0000_0000_0000_0000_0000

	PSRTransferMRSMask   = 0b0000_1111_1011_1111_0000_0000_0000_0000
	PSRTransferMRSFormat = 0b0000_0001_0000_1111_0000_0000_0000_0000

	PSRTransferMSRMask   = 0b0000_1101_1011_0000_1111_0000_0000_0000
	PSRTransferMSRFormat = 0b0000_0001_0010_0000_1111_0000_0000_0000

	MultiplyMask   = 0b0000_1111_1100_0000_0000_0000_1111_0000
	MultiplyFormat = 0b0000_0000_0000_0000_0000_0000_1001_0000

	MultiplyLongMask   = 0b0000_1111_1000_0000_0000_0000_1111_0000
	MultiplyLongFormat = 0b0000_0000_1000_0000_0000_0000_1001_0000

	SingleDataSwapMask   = 0b0000_1111_1011_0000_0000_1111_1111_0000
	SingleDataSwapFormat = 0b0000_0001_0000_0000_0000_0000_1001_0000

	BranchExchangeMask   = 0b0000_1111_1111_1111_1111_1111_1111_0000
	BranchExchangeFormat = 0b0000_0001_0010_1111_1111_1111_0001_0000

	HalfwordDataTransferRegisterOffsetMask   = 0b0000_1110_0100_0000_0000_1111_1001_0000
	HalfwordDataTransferRegisterOffsetFormat = 0b0000_0000_0000_0000_0000_0000_1001_0000

	HalfwordDataTransferImmediateOffsetMask   = 0b0000_1110_0100_0000_0000_0000_1001_0000
	HalfwordDataTransferImmediateOffsetFormat = 0b0000_0000_0100_0000_0000_0000_1001_0000

	SingleDataTransferMask   = 0b0000_1100_0000_0000_0000_0000_0000_0000
	SingleDataTransferFormat = 0b0000_0100_0000_0000_0000_0000_0000_0000

	UndefinedMask   = 0b0000_1110_0000_0000_0000_0000_0001_0000
	UndefinedFormat = 0b0000_0110_0000_0000_0000_0000_0001_0000

	BlockDataTransferMask   = 0b0000_1110_0000_0000_0000_0000_0000_0000
	BlockDataTransferFormat = 0b0000_1000_0000_0000_0000_0000_0000_0000

	BranchMask           = 0b0000_1111_0000_0000_0000_0000_0000_0000
	BranchFormat         = 0b0000_1010_0000_0000_0000_0000_0000_0000
	BranchWithLinkFormat = 0b0000_1011_0000_0000_0000_0000_0000_0000

	// CoProcessorDataTransferMask   = 0b0000_1110_0000_0000_0000_0000_0000_0000
	// CoProcessorDataTransferFormat = 0b0000_1100_0000_0000_0000_0000_0000_0000

	// CoProcessorDataOperationMask   = 0b0000_1111_0000_0000_0000_0000_0001_0000
	// CoProcessorDataOperationFormat = 0b0000_1110_0000_0000_0000_0000_0000_0000

	// CoProcessorRegisterTransferMask   = 0b0000_1111_0000_0000_0000_0000_0001_0000
	// CoProcessorRegisterTransferFormat = 0b0000_1110_0000_0000_0000_0000_0001_0000

	SoftwareInterruptMask   = 0b0000_1111_0000_0000_0000_0000_0000_0000
	SoftwareInterruptFormat = 0b0000_1111_0000_0000_0000_0000_0000_0000
)

func DecodeInstruction(instruction uint32) isa.Instruction {
	// This function will check masks against the instruction to determine which
	// type of operation it is. Then, the opcode will be used to determine which
	// instruction to execute.
	switch {
	case instruction&BranchExchangeMask == BranchExchangeFormat:
		return BX{instruction}
	case instruction&BlockDataTransferMask == BlockDataTransferFormat:
		return matchBlockDataTransfer(instruction)
	case instruction&BranchMask == BranchFormat:
		return B{instruction}
	case instruction&BranchMask == BranchWithLinkFormat:
		return BL{instruction}
	case instruction&SoftwareInterruptMask == SoftwareInterruptFormat:
		return SWI{instruction}
	case instruction&UndefinedMask == UndefinedFormat:
		return matchUndefined(instruction)
	case instruction&SingleDataTransferMask == SingleDataTransferFormat:
		return matchSingleDataTransfer(instruction)
	case instruction&SingleDataSwapMask == SingleDataSwapFormat:
		return matchSingleDataSwap(instruction)
	case instruction&MultiplyMask == MultiplyFormat:
		// Bit 21 == 1 for multiple and accumulate, 0 for multiply
		accumulate := (instruction & (1 << 21)) != 0
		if accumulate {
			return MLA{instruction}
		}
		return MUL{instruction}
	case instruction&MultiplyLongMask == MultiplyLongFormat:
		return matchMultiplyLong(instruction)
	case instruction&HalfwordDataTransferRegisterOffsetMask == HalfwordDataTransferRegisterOffsetFormat:
		return matchHalfwordDataTransferRegisterOffset(instruction)
	case instruction&HalfwordDataTransferImmediateOffsetMask == HalfwordDataTransferImmediateOffsetFormat:
		return matchHalfwordDataTransferImmediateOffset(instruction)
	case instruction&PSRTransferMRSMask == PSRTransferMRSFormat:
		return MRS{instruction}
	case instruction&PSRTransferMSRMask == PSRTransferMSRFormat:
		return MSR{instruction}
	case instruction&DataProcessingMask == DataProcessingFormat:
		return matchDataProcessing(instruction)
	default:
		fmt.Println("Unknown instruction")
		panic(fmt.Sprintf("Unknown instruction: 0x%08X", instruction))
	}
}

func matchBlockDataTransfer(instruction uint32) isa.Instruction {
	// Get bit 20 == 1 for load, 0 for store
	load := (instruction & (1 << 20)) != 0

	if load {
		return LDM{instruction}
	}
	return STM{instruction}
}

func matchUndefined(_ uint32) isa.Instruction {
	fmt.Println("Undefined")
	return nil
}

func matchSingleDataTransfer(instruction uint32) isa.Instruction {
	// instruction is in little endian
	// Get bit 20 == 1 for load, 0 for store
	load := (instruction & (1 << 20)) != 0

	if load {
		return LDR{instruction}
	}
	return STR{instruction}
}

func matchSingleDataSwap(_ uint32) isa.Instruction {
	fmt.Println("Single Data Swap")
	return nil
}

func matchMultiplyLong(instruction uint32) isa.Instruction {
	fmt.Println("Multiply Long")
	// Bits 22 and 21 are the U and A flags, unsigned and accumulate
	ua := (instruction >> 21) & 0b11
	switch ua {
	case 0b00:
		// unsigned, multiply only
		return UMULL{instruction}
	case 0b01:
		// unsigned, multiply and accumulate
		return UMLAL{instruction}
	case 0b10:
		// signed, multiply only
		return SMULL{instruction}
	case 0b11:
		// signed, multiply and accumulate
		return SMLAL{instruction}
	}
	return nil
}

func matchHalfwordDataTransferRegisterOffset(instruction uint32) isa.Instruction {
	// Get bit 20 == 1 for load, 0 for store
	load := (instruction & (1 << 20)) != 0

	// Bit 6 is the s flag
	s := (instruction & (1 << 6)) != 0

	// Bit 5 is the h flag
	h := (instruction & (1 << 5)) != 0

	if load {
		if s {
			if h {
				return LDRSHRegisterOffset{instruction}
			}
			return LDRSBRegisterOffset{instruction}
		}
		if h {
			return LDRHRegisterOffset{instruction}
		}
		// SWP instruction
		panic("SWP instruction not implemented")
	} else {
		if s {
			if h {
				return STRSHRegisterOffset{instruction}
			}
			return STRSBRegisterOffset{instruction}
		}
		if h {
			return STRHRegisterOffset{instruction}
		}
		// SWP instruction
		panic("SWP instruction not implemented")
	}
}

func matchHalfwordDataTransferImmediateOffset(instruction uint32) isa.Instruction {
	// Get bit 20 == 1 for load, 0 for store
	load := (instruction & (1 << 20)) != 0

	// Bit 6 is the s flag
	s := (instruction & (1 << 6)) != 0

	// Bit 5 is the h flag
	h := (instruction & (1 << 5)) != 0

	if load {
		if s {
			if h {
				return LDRSH{instruction}
			}
			return LDRSB{instruction}
		}
		if h {
			return LDRH{instruction}
		}
		// SWP instruction
		panic("SWP instruction not implemented")
	} else {
		if s {
			if h {
				return STRSH{instruction}
			}
			return STRSB{instruction}
		}
		if h {
			return STRH{instruction}
		}
		// SWP instruction
		panic("SWP instruction not implemented")
	}
}

func matchDataProcessing(instruction uint32) isa.Instruction {
	// Opcode is bits 24-21
	opcode := (instruction & 0x01E00000) >> 21
	switch opcode {
	case 0b0000: // AND
		return AND{instruction}
	case 0b0001: // EOR
		return EOR{instruction}
	case 0b0010: // SUB
		return SUB{instruction}
	case 0b0011: // RSB
		return RSB{instruction}
	case 0b0100: // ADD
		return ADD{instruction}
	case 0b0101: // ADC
		return ADC{instruction}
	case 0b0110: // SBC
		return SBC{instruction}
	case 0b0111: // RSC
		return RSC{instruction}
	case 0b1000: // TST
		return TST{instruction}
	case 0b1001: // TEQ
		return TEQ{instruction}
	case 0b1010: // CMP
		return CMP{instruction}
	case 0b1011: // CMN
		return CMN{instruction}
	case 0b1100: // ORR
		return ORR{instruction}
	case 0b1101: // MOV
		return MOV{instruction}
	case 0b1110: // BIC
		return BIC{instruction}
	case 0b1111: // MVN
		return MVN{instruction}
	default:
		panic(fmt.Sprintf("Unknown opcode: 0x%04b", opcode))
	}
}

//nolint:golint,unparam
func unshiftRegister(instruction uint32, cpu interfaces.CPU) (uint32, bool) {
	carry := false
	// Shift is bits 11-4
	shift := (instruction & 0x00000FF0) >> 4
	// RM is bits 3-0
	rm := uint8(instruction & 0x0000000F)

	var shiftAmount uint32
	// fmt.Printf("Not calculating carry bit for now\n")

	if shift&0b0000_1001 == 1 {
		// Bits 11-8 refer to the register, of which the bottom byte is the shift amount.
		shiftRegister := uint8((shift & 0b1111_0000) >> 4)
		shiftAmount = cpu.ReadRegister(shiftRegister) & 0x000000FF
	} else {
		// Bits 11-7 are the shift amount.
		shiftAmount = (shift & 0b1111_1000) >> 3
	}
	switch (shift & 0b0000_0110) >> 1 {
	case 0b0000_0000: // Logical shift left
		if cpu.GetConfig().Debug {
			fmt.Printf("lsl r%d, #%d\n", rm, shiftAmount)
		}
		if shiftAmount == 0 {
			return cpu.ReadRegister(rm), cpu.GetC()
		}
		return cpu.ReadRegister(rm) << shiftAmount, carry
	case 0b0000_0001: // Logical shift right
		if shiftAmount == 0 {
			if cpu.GetConfig().Debug {
				fmt.Printf("lsr r%d, #32\n", rm)
			}
			// bit 31 of Rm is copied into the carry flag
			carry = (cpu.ReadRegister(rm) & 0x80000000) != 0
			return 0, carry
		}
		if cpu.GetConfig().Debug {
			fmt.Printf("lsr r%d, #%d\n", rm, shiftAmount)
		}
		return cpu.ReadRegister(rm) >> shiftAmount, carry
	case 0b0000_0010: // Arithmetic shift right
		// An arithmetic shift right (ASR) is similar to logical shift right, except that the high bits
		// are filled with bit 31 of Rm instead of zeros
		if cpu.GetConfig().Debug {
			fmt.Printf("asr r%d, #%d\n", rm, shiftAmount)
		}
		if shiftAmount == 0 {
			// bit 31 of Rm is copied into the carry flag
			// return all ones or all zeros depending on bit 31 of Rm
			carry = (cpu.ReadRegister(rm) & 0x80000000) != 0
			if carry {
				return 0xFFFFFFFF, carry
			}
			return 0, carry
		}
		carryBit := (cpu.ReadRegister(rm) & (1 << (shiftAmount - 1))) >> (shiftAmount - 1)

		// Shift right by shiftAmount
		shifted := cpu.ReadRegister(rm) >> shiftAmount
		// Get the sign bit
		highBit := shifted & (1 << 31)

		// Fill the high bits we shifted to zeros with the highBit
		for i := uint32(0); i < shiftAmount; i++ {
			shifted |= highBit << (31 - i)
		}

		return shifted, carryBit == 1
	case 0b0000_0011: // Rotate right
		if cpu.GetConfig().Debug {
			fmt.Printf("ror r%d, #%d\n", rm, shiftAmount)
		}
		if shiftAmount == 32 {
			// Result is the same as the original value, carry flag is bit 31 of Rm
			carry = cpu.ReadRegister(rm)&(1<<31) != 0
			return cpu.ReadRegister(rm), carry
		}
		if shiftAmount == 0 {
			// RRX
			rmValue := cpu.ReadRegister(rm)
			// existing carry flag becomes bit 31
			// bit 0 of Rm is copied into the carry flag
			carry = (rmValue & 0x00000001) != 0
			rmValue >>= 1
			if cpu.GetC() {
				rmValue |= 0x80000000
			} else {
				rmValue &= 0x7FFFFFFF
			}
			return rmValue, carry
		}
		return bits.RotateLeft32(cpu.ReadRegister(rm), -int(shiftAmount)), carry
	}

	return 0, carry
}

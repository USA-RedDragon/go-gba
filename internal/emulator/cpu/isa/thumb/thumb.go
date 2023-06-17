package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/cpu/isa"
)

// http://bear.ces.cwru.edu/eecs_382/ARM7-TDMI-manual-pt3.pdf, Figure 5-1
const (
	MoveShiftedRegisterMask                    uint16 = 0b1110_0000_0000_0000
	MoveShiftedRegisterFormat                  uint16 = 0b0000_0000_0000_0000
	AddSubtractMask                            uint16 = 0b1111_1000_0000_0000
	AddSubtractFormat                          uint16 = 0b0001_1000_0000_0000
	MoveCompareAddSubtractImmediateMask        uint16 = 0b1110_0000_0000_0000
	MoveCompareAddSubtractImmediateFormat      uint16 = 0b0010_0000_0000_0000
	AluOperationMask                           uint16 = 0b1111_1100_0000_0000
	AluOperationFormat                         uint16 = 0b0100_0000_0000_0000
	HiRegisterOperationsOrBranchExchangeMask   uint16 = 0b1111_1100_0000_0000
	HiRegisterOperationsOrBranchExchangeFormat uint16 = 0b0100_0100_0000_0000
	PCRelativeLoadMask                         uint16 = 0b1111_1000_0000_0000
	PCRelativeLoadFormat                       uint16 = 0b0100_1000_0000_0000
	LoadStoreWithRegisterOffsetMask            uint16 = 0b1111_0010_0000_0000
	LoadStoreWithRegisterOffsetFormat          uint16 = 0b0101_0000_0000_0000
	LoadStoreSignExtendedByteHalfwordMask      uint16 = 0b1111_0010_0000_0000
	LoadStoreSignExtendedByteHalfwordFormat    uint16 = 0b0101_0010_0000_0000
	LoadStoreWithImmediateOffsetMask           uint16 = 0b1110_0000_0000_0000
	LoadStoreWithImmediateOffsetFormat         uint16 = 0b0110_0000_0000_0000
	LoadStoreHalfwordMask                      uint16 = 0b1111_0000_0000_0000
	LoadStoreHalfwordFormat                    uint16 = 0b1000_0000_0000_0000
	SPRelativeLoadStoreMask                    uint16 = 0b1111_0000_0000_0000
	SPRelativeLoadStoreFormat                  uint16 = 0b1001_0000_0000_0000
	LoadAddressMask                            uint16 = 0b1111_0000_0000_0000
	LoadAddressFormat                          uint16 = 0b1010_0000_0000_0000
	AddOffsetToStackPointerMask                uint16 = 0b1111_1111_0000_0000
	AddOffsetToStackPointerFormat              uint16 = 0b1011_0000_0000_0000
	PushPopRegistersMask                       uint16 = 0b1111_0110_0000_0000
	PushPopRegistersFormat                     uint16 = 0b1011_0100_0000_0000
	MultipleLoadStoreMask                      uint16 = 0b1111_0000_0000_0000
	MultipleLoadStoreFormat                    uint16 = 0b1100_0000_0000_0000
	ConditionalBranchMask                      uint16 = 0b1111_0000_0000_0000
	ConditionalBranchFormat                    uint16 = 0b1101_0000_0000_0000
	SoftwareInterruptMask                      uint16 = 0b1111_1111_0000_0000
	SoftwareInterruptFormat                    uint16 = 0b1101_1111_0000_0000
	UnconditionalBranchMask                    uint16 = 0b1111_1000_0000_0000
	UnconditionalBranchFormat                  uint16 = 0b1110_0000_0000_0000
	LongBranchWithLinkMask                     uint16 = 0b1111_0000_0000_0000
	LongBranchWithLinkFormat                   uint16 = 0b1111_0000_0000_0000
)

func DecodeInstruction(instruction uint16) isa.Instruction {
	// This function will check masks against the instruction to determine which
	// type of operation it is. Then, the opcode will be used to determine which
	// instruction to execute.
	switch {
	case instruction&SoftwareInterruptMask == SoftwareInterruptFormat:
		return matchSoftwareInterrupt(instruction)
	case instruction&UnconditionalBranchMask == UnconditionalBranchFormat:
		return UnconditionalBranch{instruction}
	case instruction&ConditionalBranchMask == ConditionalBranchFormat:
		return B{instruction}
	case instruction&MultipleLoadStoreMask == MultipleLoadStoreFormat:
		return matchMultipleLoadStore(instruction)
	case instruction&LongBranchWithLinkMask == LongBranchWithLinkFormat:
		return LBL{instruction}
	case instruction&AddOffsetToStackPointerMask == AddOffsetToStackPointerFormat:
		return SUBSP{instruction}
	case instruction&PushPopRegistersMask == PushPopRegistersFormat:
		push := instruction&(1<<11)>>11 == 0
		if push {
			return PUSH{instruction}
		} else {
			return POP{instruction}
		}
	case instruction&LoadStoreHalfwordMask == LoadStoreHalfwordFormat:
		// Bit 11 = 1 for LDRH, 0 for STRH
		ldr := instruction&(1<<11)>>11 == 1
		if ldr {
			return LDRH{instruction}
		} else {
			return STRH{instruction}
		}
		return nil
	case instruction&SPRelativeLoadStoreMask == SPRelativeLoadStoreFormat:
		// Bit 11 == 1 for LDRPC, 0 for STRPC
		ldr := instruction&(1<<11)>>11 == 1
		if ldr {
			return LDRSP{instruction}
		} else {
			return STRSP{instruction}
		}
	case instruction&LoadAddressMask == LoadAddressFormat:
		return matchLoadAddress(instruction)
	case instruction&LoadStoreWithImmediateOffsetMask == LoadStoreWithImmediateOffsetFormat:
		// Bit 12 == 1 for LDRImmB, 0 for LDRImmW
		isByte := instruction&(1<<12)>>12 == 1
		// Bit 11 == 1 for LDRImm, 0 for STRImm
		ldr := instruction&(1<<11)>>11 == 1
		if isByte {
			if ldr {
				return LDRBImm{instruction}
			} else {
				return STRBImm{instruction}
			}
		} else {
			if ldr {
				return LDRWImm{instruction}
			} else {
				return STRWImm{instruction}
			}
		}
		return nil
	case instruction&LoadStoreWithRegisterOffsetMask == LoadStoreWithRegisterOffsetFormat:
		return matchLoadStoreWithRegisterOffset(instruction)
	case instruction&LoadStoreSignExtendedByteHalfwordMask == LoadStoreSignExtendedByteHalfwordFormat:
		return matchLoadStoreSignExtendedByteHalfword(instruction)
	case instruction&PCRelativeLoadMask == PCRelativeLoadFormat:
		return LDR{instruction}
	case instruction&HiRegisterOperationsOrBranchExchangeMask == HiRegisterOperationsOrBranchExchangeFormat:
		// bits 9-8 are the opcode
		opcode := instruction & (1<<9 | 1<<8) >> 8
		switch opcode {
		case 0:
			return ADDH{instruction}
		case 1:
			return CMPH{instruction}
		case 2:
			return MOVH{instruction}
		case 3:
			return BX{instruction}
		}
		return nil
	case instruction&AluOperationMask == AluOperationFormat:
		// Bits 9-6 are the opcode
		opcode := instruction & (1<<9 | 1<<8 | 1<<7 | 1<<6) >> 6
		switch opcode {
		case 0b0000:
			return AND{instruction}
		case 0b0001:
			return EOR{instruction}
		case 0b0010:
			return LSL{instruction}
		case 0b0011:
			return LSR{instruction}
		case 0b0100:
			return ASR{instruction}
		case 0b0101:
			return ADC{instruction}
		case 0b0110:
			return SBC{instruction}
		case 0b0111:
			return ROR{instruction}
		case 0b1000:
			return TST{instruction}
		case 0b1001:
			return NEG{instruction}
		case 0b1010:
			return CMPALU{instruction}
		case 0b1011:
			return CMN{instruction}
		case 0b1100:
			return ORR{instruction}
		case 0b1101:
			return MUL{instruction}
		case 0b1110:
			return BIC{instruction}
		case 0b1111:
			return MVN{instruction}
		}
		return nil
	case instruction&MoveCompareAddSubtractImmediateMask == MoveCompareAddSubtractImmediateFormat:
		op := instruction & (1<<12 | 1<<11) >> 11
		fmt.Printf("Instruction: 0x%04X\n", instruction)
		switch op {
		case 0:
			return MOV{instruction}
		case 1:
			return CMP{instruction}
		case 2:
			return ADD{instruction}
		case 3:
			return SUB{instruction}
		}
		return nil
	case instruction&AddSubtractMask == AddSubtractFormat:
		subtract := instruction&(1<<9)>>9 == 1
		if subtract {
			return SUB2{instruction}
		} else {
			return ADD2{instruction}
		}
	case instruction&MoveShiftedRegisterMask == MoveShiftedRegisterFormat:
		// Bits 12-11 are the opcode
		op := instruction & (1<<12 | 1<<11) >> 11

		switch op {
		case 0b00:
			return LSLMoveShifted{instruction}
		case 0b01:
			return LSRMoveShifted{instruction}
		case 0b10:
			return ASRMoveShifted{instruction}
		}
		return nil
	default:
		panic(fmt.Sprintf("Unknown THUMB instruction: %016b", instruction))
	}
}

func matchSoftwareInterrupt(instruction uint16) isa.Instruction {
	fmt.Println("SoftwareInterrupt")
	return nil
}

func matchMultipleLoadStore(instruction uint16) isa.Instruction {
	fmt.Println("MultipleLoadStore")
	return nil
}

func matchLoadAddress(instruction uint16) isa.Instruction {
	fmt.Println("LoadAddress")
	return nil
}

func matchLoadStoreWithRegisterOffset(instruction uint16) isa.Instruction {
	// Bit 11 is the L bit, which determines whether this is a load or store
	load := instruction&(1<<11)>>11 == 1
	if load {
		return LDRR{instruction}
	} else {
		return STRR{instruction}
	}
}

func matchLoadStoreSignExtendedByteHalfword(instruction uint16) isa.Instruction {
	fmt.Println("LoadStoreSignExtendedByteHalfword")
	return nil
}

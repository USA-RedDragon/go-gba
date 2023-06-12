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
	MoveCompareAddSubtractImmediateMask        uint16 = 0b1111_1000_0000_0000
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
		return matchUnconditionalBranch(instruction)
	case instruction&ConditionalBranchMask == ConditionalBranchFormat:
		return matchConditionalBranch(instruction)
	case instruction&MultipleLoadStoreMask == MultipleLoadStoreFormat:
		return matchMultipleLoadStore(instruction)
	case instruction&LongBranchWithLinkMask == LongBranchWithLinkFormat:
		return matchLongBranchWithLink(instruction)
	case instruction&AddOffsetToStackPointerMask == AddOffsetToStackPointerFormat:
		return matchAddOffsetToStackPointer(instruction)
	case instruction&PushPopRegistersMask == PushPopRegistersFormat:
		return matchPushPopRegisters(instruction)
	case instruction&LoadStoreHalfwordMask == LoadStoreHalfwordFormat:
		return matchLoadStoreHalfword(instruction)
	case instruction&SPRelativeLoadStoreMask == SPRelativeLoadStoreFormat:
		return matchSPRelativeLoadStore(instruction)
	case instruction&LoadAddressMask == LoadAddressFormat:
		return matchLoadAddress(instruction)
	case instruction&LoadStoreWithImmediateOffsetMask == LoadStoreWithImmediateOffsetFormat:
		return matchLoadStoreWithImmediateOffset(instruction)
	case instruction&LoadStoreWithRegisterOffsetMask == LoadStoreWithRegisterOffsetFormat:
		return matchLoadStoreWithRegisterOffset(instruction)
	case instruction&LoadStoreSignExtendedByteHalfwordMask == LoadStoreSignExtendedByteHalfwordFormat:
		return matchLoadStoreSignExtendedByteHalfword(instruction)
	case instruction&PCRelativeLoadMask == PCRelativeLoadFormat:
		return LDR{instruction}
	case instruction&HiRegisterOperationsOrBranchExchangeMask == HiRegisterOperationsOrBranchExchangeFormat:
		return matchHiRegisterOperationsOrBranchExchange(instruction)
	case instruction&AluOperationMask == AluOperationFormat:
		return matchAluOperation(instruction)
	case instruction&MoveCompareAddSubtractImmediateMask == MoveCompareAddSubtractImmediateFormat:
		op := instruction & (1<<12 | 1<<11) >> 11
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
		return matchAddSubtract(instruction)
	case instruction&MoveShiftedRegisterMask == MoveShiftedRegisterFormat:
		return matchMoveShiftedRegister(instruction)
	default:
		panic(fmt.Sprintf("Unknown THUMB instruction: 0x%04X", instruction))
	}
}

func matchSoftwareInterrupt(instruction uint16) isa.Instruction {
	fmt.Println("SoftwareInterrupt")
	return nil
}

func matchUnconditionalBranch(instruction uint16) isa.Instruction {
	fmt.Println("UnconditionalBranch")
	return nil
}

func matchConditionalBranch(instruction uint16) isa.Instruction {
	fmt.Println("ConditionalBranch")
	return nil
}

func matchMultipleLoadStore(instruction uint16) isa.Instruction {
	fmt.Println("MultipleLoadStore")
	return nil
}

func matchLongBranchWithLink(instruction uint16) isa.Instruction {
	fmt.Println("LongBranchWithLink")
	return nil
}

func matchAddOffsetToStackPointer(instruction uint16) isa.Instruction {
	fmt.Println("AddOffsetToStackPointer")
	return nil
}

func matchPushPopRegisters(instruction uint16) isa.Instruction {
	fmt.Println("PushPopRegisters")
	return nil
}

func matchLoadStoreHalfword(instruction uint16) isa.Instruction {
	fmt.Println("LoadStoreHalfword")
	return nil
}

func matchSPRelativeLoadStore(instruction uint16) isa.Instruction {
	fmt.Println("SPRelativeLoadStore")
	return nil
}

func matchLoadAddress(instruction uint16) isa.Instruction {
	fmt.Println("LoadAddress")
	return nil
}

func matchLoadStoreWithImmediateOffset(instruction uint16) isa.Instruction {
	fmt.Println("LoadStoreWithImmediateOffset")
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

func matchHiRegisterOperationsOrBranchExchange(instruction uint16) isa.Instruction {
	fmt.Println("HiRegisterOperationsOrBranchExchange")
	return nil
}

func matchAluOperation(instruction uint16) isa.Instruction {
	fmt.Println("AluOperation")
	return nil
}

func matchAddSubtract(instruction uint16) isa.Instruction {
	fmt.Println("AddSubtract")
	return nil
}

func matchMoveShiftedRegister(instruction uint16) isa.Instruction {
	fmt.Println("MoveShiftedRegister")
	return nil
}

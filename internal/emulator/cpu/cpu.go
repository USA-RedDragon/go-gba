package cpu

import (
	"fmt"
	"os"
	"time"

	"github.com/USA-RedDragon/go-gba/internal/config"
	"github.com/USA-RedDragon/go-gba/internal/emulator/cpu/isa/arm"
	"github.com/USA-RedDragon/go-gba/internal/emulator/cpu/isa/thumb"
	"github.com/USA-RedDragon/go-gba/internal/emulator/memory"
)

const (
	// BIOSROMSize is 16KB
	BIOSROMSize = 16 * 1024
	// OnChipRAMSize is 32KB
	OnChipRAMSize = 32 * 1024
	// OnBoardRAMSize is 256KB
	OnBoardRAMSize = 256 * 1024
	// IORAMSize is 1KB
	IORAMSize = 1 * 1024
	SP_REG    = 13
	LR_REG    = 14
	PC_REG    = 15
	CPSR_REG  = 16
)

type ARM7TDMI struct {
	// Registers R0-R16
	r [17]uint32

	sp_irq   uint32
	lr_irq   uint32
	spsr_irq uint32

	r8_fiq   uint32
	r9_fiq   uint32
	r10_fiq  uint32
	r11_fiq  uint32
	r12_fiq  uint32
	sp_fiq   uint32
	lr_fiq   uint32
	spsr_fiq uint32

	sp_svc   uint32
	lr_svc   uint32
	spsr_svc uint32

	sp_abt   uint32
	lr_abt   uint32
	spsr_abt uint32

	sp_und   uint32
	lr_und   uint32
	spsr_und uint32

	virtualMemory memory.MMIO

	biosROM    [BIOSROMSize]byte
	onChipRAM  [OnChipRAMSize]byte
	onBoardRAM [OnBoardRAMSize]byte
	ioRAM      [IORAMSize]byte

	halted bool
	exit   bool

	prefetchARMPipeline   [2]uint32
	prefetchThumbPipeline [2]uint16

	config *config.Config
}

// Enum for CPU mode
type cpuMode uint8

const (
	userMode       cpuMode = 0b10000
	fiqMode        cpuMode = 0b10001
	irqMode        cpuMode = 0b10010
	supervisorMode cpuMode = 0b10011
	abortMode      cpuMode = 0b10111
	undefinedMode  cpuMode = 0b11011
	systemMode     cpuMode = 0b11111
)

func NewARM7TDMI(config *config.Config) *ARM7TDMI {
	cpu := &ARM7TDMI{
		virtualMemory: memory.MMIO{},
		biosROM:       [BIOSROMSize]byte{},
		onChipRAM:     [OnChipRAMSize]byte{},
		config:        config,
	}
	cpu.virtualMemory.AddMMIO(cpu.biosROM[:], 0x00000000, BIOSROMSize)
	// 0x00004000-0x01FFFFFF is unused
	cpu.virtualMemory.AddMMIO(cpu.onBoardRAM[:], 0x02000000, OnBoardRAMSize)
	// 0x02040000-0x02FFFFFF is unused
	cpu.virtualMemory.AddMMIO(cpu.onChipRAM[:], 0x03000000, OnChipRAMSize)
	// 0x03008000-0x03FFFFFF is unused
	cpu.virtualMemory.AddMMIO(cpu.ioRAM[:], 0x04000000, IORAMSize)
	// 0x04000400-0x04FFFFFF is unused
	cpu.loadBIOSROM()
	cpu.Reset()
	return cpu
}

func (c *ARM7TDMI) RegisterMMIO(data []byte, address uint32, size uint32) {
	c.virtualMemory.AddMMIO(data, address, size)
}

func (c *ARM7TDMI) loadBIOSROM() {
	bios, err := os.ReadFile(c.config.BIOSPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load bios rom: %v", err))
	}
	if len(bios) != BIOSROMSize {
		panic(fmt.Sprintf("BIOS ROM size is %d, expected %d", len(bios), BIOSROMSize))
	}
	copy(c.biosROM[:], bios)
}

func (c *ARM7TDMI) Reset() {
	c.halted = true
	c.exit = false

	c.lr_svc = c.r[PC_REG]
	c.sp_svc = c.r[SP_REG]

	c.r[SP_REG] = 0x03007F00 // Stack pointer to the top of on-chip RAM

	// IRQs disabled, FIQs disabled, ARM mode, system mode
	c.r[CPSR_REG] = 0x1F

	c.r[PC_REG] = 0x00000000

	// Initialize the prefetch buffers
	var err error
	c.prefetchARMPipeline[0], err = c.virtualMemory.Read32(c.r[PC_REG])
	if err != nil {
		panic(fmt.Sprintf("Failed to read from memory: %v", err))
	}
	c.prefetchARMPipeline[1], err = c.virtualMemory.Read32(c.r[PC_REG] + 4)
	if err != nil {
		panic(fmt.Sprintf("Failed to read from memory: %v", err))
	}
	c.prefetchThumbPipeline[0], err = c.virtualMemory.Read16(c.r[PC_REG])
	if err != nil {
		panic(fmt.Sprintf("Failed to read from memory: %v", err))
	}
	c.prefetchThumbPipeline[1], err = c.virtualMemory.Read16(c.r[PC_REG] + 2)
	if err != nil {
		panic(fmt.Sprintf("Failed to read from memory: %v", err))
	}

	c.r[PC_REG] = 0x00000004 // Reset vector

	if c.config.Debug {
		fmt.Printf("Resetting CPU\n")
	}

	c.halted = false
	c.exit = false
}

func (c *ARM7TDMI) ReadSPSR() uint32 {
	switch cpuMode(c.r[CPSR_REG] & 0x1F) {
	case systemMode:
		panic("System mode does not have an SPSR")
	case userMode:
		panic("User mode does not have an SPSR")
	case fiqMode:
		return c.spsr_fiq
	case irqMode:
		return c.spsr_irq
	case supervisorMode:
		return c.spsr_svc
	case abortMode:
		return c.spsr_abt
	case undefinedMode:
		return c.spsr_und
	default:
		panic("Unknown CPU mode")
	}
}

func (c *ARM7TDMI) WriteSPSR(value uint32) {
	switch cpuMode(c.r[CPSR_REG] & 0x1F) {
	case systemMode:
		panic("System mode does not have an SPSR")
	case userMode:
		panic("User mode does not have an SPSR")
	case fiqMode:
		c.spsr_fiq = value
	case irqMode:
		c.spsr_irq = value
	case supervisorMode:
		c.spsr_svc = value
	case abortMode:
		c.spsr_abt = value
	case undefinedMode:
		c.spsr_und = value
	default:
		panic("Unknown CPU mode")
	}
}

func (c *ARM7TDMI) ReadHighRegister(reg uint8) uint32 {
	if !c.GetThumbMode() {
		panic("Cannot read high register in ARM mode, use ReadRegister")
	}
	if reg > 15 {
		panic(fmt.Sprintf("Invalid register number %d", reg))
	}
	return c.r[reg+8]
}

func (c *ARM7TDMI) WriteHighRegister(reg uint8, value uint32) {
	if !c.GetThumbMode() {
		panic("Cannot write high register in ARM mode, use WriteRegister")
	}
	if reg > 15 {
		panic(fmt.Sprintf("Invalid register number %d", reg))
	}
	c.r[reg+8] = value
}

func (c *ARM7TDMI) ReadRegister(reg uint8) uint32 {
	if c.GetThumbMode() {
		if reg > 7 && reg != PC_REG && reg != LR_REG && reg != SP_REG && reg != CPSR_REG {
			panic(fmt.Sprintf("Invalid register number %d", reg))
		}
		if reg == PC_REG {
			return c.ReadPC()
		} else if reg == LR_REG {
			return c.ReadLR()
		} else if reg == SP_REG {
			return c.ReadSP()
		} else if reg == CPSR_REG {
			return c.ReadCPSR()
		}
		return c.r[reg]
	} else {
		if reg > 16 {
			panic(fmt.Sprintf("Invalid register number %d", reg))
		}
		switch cpuMode(c.r[CPSR_REG] & 0x1F) {
		case systemMode:
			return c.r[reg]
		case userMode:
			return c.r[reg]
		case fiqMode:
			switch reg {
			case 8:
				return c.r8_fiq
			case 9:
				return c.r9_fiq
			case 10:
				return c.r10_fiq
			case 11:
				return c.r11_fiq
			case 12:
				return c.r12_fiq
			case 13:
				return c.sp_fiq
			case 14:
				return c.lr_fiq
			default:
				return c.r[reg]
			}
		case irqMode:
			switch reg {
			case 13:
				return c.sp_irq
			case 14:
				return c.lr_irq
			default:
				return c.r[reg]
			}
		case supervisorMode:
			switch reg {
			case 13:
				return c.sp_svc
			case 14:
				return c.lr_svc
			default:
				return c.r[reg]
			}
		case abortMode:
			switch reg {
			case 13:
				return c.sp_abt
			case 14:
				return c.lr_abt
			default:
				return c.r[reg]
			}
		case undefinedMode:
			switch reg {
			case 13:
				return c.sp_und
			case 14:
				return c.lr_und
			default:
				return c.r[reg]
			}
		default:
			panic("Unknown CPU mode")
		}
	}
}

func (c *ARM7TDMI) WriteRegister(reg uint8, value uint32) {
	if c.GetThumbMode() {
		if reg > 7 && reg != PC_REG && reg != LR_REG && reg != SP_REG && reg != CPSR_REG {
			panic(fmt.Sprintf("Invalid register number %d", reg))
		}
		if reg == PC_REG {
			c.WritePC(value)
		} else if reg == LR_REG {
			c.WriteLR(value)
		} else if reg == SP_REG {
			c.WriteSP(value)
		} else if reg == CPSR_REG {
			c.WriteCPSR(value)
		} else {
			c.r[reg] = value
		}
	} else {
		if reg > 16 {
			panic(fmt.Sprintf("Invalid register number %d", reg))
		}
		if reg == PC_REG {
			c.WritePC(value)
		} else if reg == LR_REG {
			c.WriteLR(value)
		} else {
			switch cpuMode(c.r[CPSR_REG] & 0x1F) {
			case systemMode:
				c.r[reg] = value
			case userMode:
				c.r[reg] = value
			case fiqMode:
				switch reg {
				case 8:
					c.r8_fiq = value
				case 9:
					c.r9_fiq = value
				case 10:
					c.r10_fiq = value
				case 11:
					c.r11_fiq = value
				case 12:
					c.r12_fiq = value
				case 13:
					c.sp_fiq = value
				case 14:
					c.lr_fiq = value
				default:
					c.r[reg] = value
				}
			case irqMode:
				switch reg {
				case 13:
					c.sp_irq = value
				case 14:
					c.lr_irq = value
				default:
					c.r[reg] = value
				}
			case supervisorMode:
				switch reg {
				case 13:
					c.sp_svc = value
				case 14:
					c.lr_svc = value
				default:
					c.r[reg] = value
				}
			case abortMode:
				switch reg {
				case 13:
					c.sp_abt = value
				case 14:
					c.lr_abt = value
				default:
					c.r[reg] = value
				}
			case undefinedMode:
				switch reg {
				case 13:
					c.sp_und = value
				case 14:
					c.lr_und = value
				default:
					c.r[reg] = value
				}
			}
		}
	}
}

func (c *ARM7TDMI) ReadSP() uint32 {
	switch cpuMode(c.r[CPSR_REG] & 0x1F) {
	case systemMode:
		return c.r[SP_REG]
	case userMode:
		return c.r[SP_REG]
	case fiqMode:
		return c.sp_fiq
	case irqMode:
		return c.sp_irq
	case supervisorMode:
		return c.sp_svc
	case abortMode:
		return c.sp_abt
	case undefinedMode:
		return c.sp_und
	default:
		panic("Unknown CPU mode")
	}
}

func (c *ARM7TDMI) WriteSP(value uint32) {
	switch cpuMode(c.r[CPSR_REG] & 0x1F) {
	case systemMode:
		c.r[SP_REG] = value
	case userMode:
		c.r[SP_REG] = value
	case fiqMode:
		c.sp_fiq = value
	case irqMode:
		c.sp_irq = value
	case supervisorMode:
		c.sp_svc = value
	case abortMode:
		c.sp_abt = value
	case undefinedMode:
		c.sp_und = value
	default:
		panic("Unknown CPU mode")
	}
}

func (c *ARM7TDMI) ReadLR() uint32 {
	switch cpuMode(c.r[CPSR_REG] & 0x1F) {
	case systemMode:
		return c.r[LR_REG]
	case userMode:
		return c.r[LR_REG]
	case fiqMode:
		return c.lr_fiq
	case irqMode:
		return c.lr_irq
	case supervisorMode:
		return c.lr_svc
	case abortMode:
		return c.lr_abt
	case undefinedMode:
		return c.lr_und
	default:
		panic("Unknown CPU mode")
	}
}

func (c *ARM7TDMI) WriteLR(value uint32) {
	switch cpuMode(c.r[CPSR_REG] & 0x1F) {
	case systemMode:
		c.r[LR_REG] = value
	case userMode:
		c.r[LR_REG] = value
	case fiqMode:
		c.lr_fiq = value
	case irqMode:
		c.lr_irq = value
	case supervisorMode:
		c.lr_svc = value
	case abortMode:
		c.lr_abt = value
	case undefinedMode:
		c.lr_und = value
	default:
		panic("Unknown CPU mode")
	}
}

func (c *ARM7TDMI) ReadPC() uint32 {
	return c.r[PC_REG]
}

func (c *ARM7TDMI) WritePC(value uint32) {
	// Mask out the bottom two bits in arm mode
	if c.GetThumbMode() {
		value &= 0xFFFFFFFE
	} else {
		value &= 0xFFFFFFFC
	}
	c.r[PC_REG] = value
}

func (c *ARM7TDMI) GetMMIO() *memory.MMIO {
	return &c.virtualMemory
}

func (c *ARM7TDMI) ReadCPSR() uint32 {
	return c.r[CPSR_REG]
}

func (c *ARM7TDMI) WriteCPSR(value uint32) {
	c.r[CPSR_REG] = value
}

func (c *ARM7TDMI) fetchARM() uint32 {
	c.r[PC_REG] += 4

	instruction := c.prefetchARMPipeline[0]
	c.prefetchARMPipeline[0] = c.prefetchARMPipeline[1]

	// Prefetch the next instruction
	var err error
	c.prefetchARMPipeline[1], err = c.virtualMemory.Read32(c.r[PC_REG])
	if err != nil {
		panic(fmt.Sprintf("Error reading instruction at 0x%08X: %v", c.r[PC_REG], err))
	}
	fmt.Printf("fetchARM: Prefetching arm instruction at 0x%08X\n", c.r[PC_REG])

	fmt.Printf("fetchARM: Prefetch: [0x%08x, 0x%08x]\n", c.prefetchARMPipeline[0], c.prefetchARMPipeline[1])

	return instruction
}

func (c *ARM7TDMI) FlushPipeline() {
	// Flush the pipeline
	var err error
	fmt.Printf("FlushPipeline: Prefetching arm instruction at 0x%08X\n", c.r[PC_REG])
	c.prefetchARMPipeline[0], err = c.virtualMemory.Read32(c.r[PC_REG])
	if err != nil {
		panic(fmt.Sprintf("Error reading instruction at 0x%08X: %v", c.r[PC_REG], err))
	}

	fmt.Printf("FlushPipeline: Prefetching thumb instruction at 0x%08X\n", c.r[PC_REG])
	c.prefetchThumbPipeline[0], err = c.virtualMemory.Read16(c.r[PC_REG])
	if err != nil {
		panic(fmt.Sprintf("Error reading instruction at 0x%08X: %v", c.r[PC_REG], err))
	}

	fmt.Printf("FlushPipeline: Prefetching arm instruction at 0x%08X\n", c.r[PC_REG]+4)
	c.prefetchARMPipeline[1], err = c.virtualMemory.Read32(c.r[PC_REG] + 4)
	if err != nil {
		panic(fmt.Sprintf("Error reading instruction at 0x%08X: %v", c.r[PC_REG]+4, err))
	}

	fmt.Printf("FlushPipeline: Prefetching thumb instruction at 0x%08X\n", c.r[PC_REG]+2)
	c.prefetchThumbPipeline[1], err = c.virtualMemory.Read16(c.r[PC_REG] + 2)
	if err != nil {
		panic(fmt.Sprintf("Error reading instruction at 0x%08X: %v", c.r[PC_REG]+2, err))
	}

	if c.GetThumbMode() {
		c.r[PC_REG] += 2
	} else {
		c.r[PC_REG] += 4
	}

	fmt.Printf("FlushPipeline: Prefetch: [0x%08x, 0x%08x]\n", c.prefetchARMPipeline[0], c.prefetchARMPipeline[1])
}

func (c *ARM7TDMI) fetchThumb() uint16 {
	instruction := c.prefetchThumbPipeline[0]
	c.prefetchThumbPipeline[0] = c.prefetchThumbPipeline[1]

	// Prefetch the next instruction
	var err error
	c.prefetchThumbPipeline[1], err = c.virtualMemory.Read16(c.r[PC_REG] + 2)
	if err != nil {
		panic(fmt.Sprintf("fetchThumb: Error reading instruction at 0x%08X: %v", c.r[PC_REG]+2, err))
	}
	fmt.Printf("fetchThumb: Prefetching thumb instruction at 0x%08X\n", c.r[PC_REG]+2)

	c.r[PC_REG] += 2

	fmt.Printf("fetchThumb: Prefetch: [0x%04x, 0x%04x]\n", c.prefetchThumbPipeline[0], c.prefetchThumbPipeline[1])

	return instruction
}

func (c *ARM7TDMI) stepARM() {
	// FETCH
	instruction := c.fetchARM()

	if c.config.TraceRegisters {
		fmt.Printf("\n\nInstruction 0x%08X\n", instruction)
		fmt.Printf("R0:  0x%08X\t\t R1: 0x%08X\t\tR3:   0x%08X\n", c.r[0], c.r[1], c.r[3])
		fmt.Printf("R4:  0x%08X\t\t R5: 0x%08X\t\tR6:   0x%08X\n", c.r[4], c.r[5], c.r[6])
		fmt.Printf("R7:  0x%08X\t\t R8: 0x%08X\t\tR9:   0x%08X\n", c.r[7], c.r[8], c.r[9])
		fmt.Printf("R10: 0x%08X\t\tR11: 0x%08X\t\tR12:  0x%08X\n", c.r[10], c.r[11], c.r[12])
		fmt.Printf("SP:  0x%08X\t\t LR: 0x%08X\t\tPC:   0x%08X\n", c.ReadSP(), c.ReadLR(), c.ReadPC())
		fmt.Println(c.prettyCPSR())
		if cpuMode(c.r[CPSR_REG]&0x1F) != systemMode && cpuMode(c.r[CPSR_REG]&0x1F) != userMode {
			fmt.Printf("SPSR: 0x%08X\n", c.ReadSPSR())
		}
	}

	// DECODE
	var condition uint32 = instruction >> 28
	conditionFailed := false
	// c.r[CSPR_REG] has condition flags N Z C V at bits 31-28
	// c.r[CSPR_REG] has control bits I F T at bits 7-5
	// c.r[CSPR_REG] has mode bits M4 M3 M2 M1 M0 at bits 4-0
	switch condition {
	case 0b0000 /* EQ, equal */ :
		// If the CSPR Z flag (bit 30) is set, then the instruction is executed
		if c.r[CPSR_REG]&(1<<30) == 0 {
			fmt.Printf("Skipping instruction 0x%08X becaue of EQ conditional\n", instruction)
			conditionFailed = true
		}
	case 0b0001 /* NE, not equal */ :
		// If the CSPR Z flag (bit 30) is clear, then the instruction is executed
		if c.r[CPSR_REG]&(1<<30) != 0 {
			fmt.Printf("Skipping instruction 0x%08X becaue of NE conditional\n", instruction)
			conditionFailed = true
		}
	case 0b0010 /* CS, unsigned higher or same */ :
		// If the CSPR C flag (bit 29) is set, then the instruction is executed
		if c.r[CPSR_REG]&(1<<29) == 0 {
			fmt.Printf("Skipping instruction 0x%08X becaue of CS conditional\n", instruction)
			conditionFailed = true
		}
	case 0b0011 /* CC, unsigned lower */ :
		// If the CSPR C flag (bit 29) is clear, then the instruction is executed
		if c.r[CPSR_REG]&(1<<29) != 0 {
			fmt.Printf("Skipping instruction 0x%08X becaue of CC conditional\n", instruction)
			conditionFailed = true
		}
	case 0b0100 /* MI, negative */ :
		// If the CSPR N flag (bit 31) is set, then the instruction is executed
		if c.r[CPSR_REG]&(1<<31) == 0 {
			fmt.Printf("Skipping instruction 0x%08X becaue of MI conditional\n", instruction)
			conditionFailed = true
		}
	case 0b0101 /* PL, positive or zero */ :
		// If the CSPR N flag (bit 31) is clear, then the instruction is executed
		if c.r[CPSR_REG]&(1<<31) != 0 {
			fmt.Printf("Skipping instruction 0x%08X becaue of PL conditional\n", instruction)
			conditionFailed = true
		}
	case 0b0110 /* VS, overflow */ :
		// If the CSPR V flag (bit 28) is set, then the instruction is executed
		if c.r[CPSR_REG]&(1<<28) == 0 {
			fmt.Printf("Skipping instruction 0x%08X becaue of VS conditional\n", instruction)
			conditionFailed = true
		}
	case 0b0111 /* VC, no overflow */ :
		// If the CSPR V flag (bit 28) is clear, then the instruction is executed
		if c.r[CPSR_REG]&(1<<28) != 0 {
			fmt.Printf("Skipping instruction 0x%08X becaue of VC conditional\n", instruction)
			conditionFailed = true
		}
	case 0b1000 /* HI, unsigned higher */ :
		// If the CSPR C flag (bit 29) is set and the CSPR Z flag (bit 30) is clear, then the instruction is executed
		if c.r[CPSR_REG]&(1<<29) == 0 || c.r[CPSR_REG]&(1<<30) != 0 {
			fmt.Printf("Skipping instruction 0x%08X becaue of HI conditional\n", instruction)
			conditionFailed = true
		}
	case 0b1001 /* LS, unsigned lower or same */ :
		// If the CSPR C flag (bit 29) is clear or the CSPR Z flag (bit 30) is set, then the instruction is executed
		if c.r[CPSR_REG]&(1<<29) != 0 || c.r[CPSR_REG]&(1<<30) == 0 {
			fmt.Printf("Skipping instruction 0x%08X becaue of LS conditional\n", instruction)
			conditionFailed = true
		}
	case 0b1010 /* GE, greater than or equal */ :
		// If the CSPR N flag (bit 31) is equal to the CSPR V flag (bit 28), then the instruction is executed
		if c.r[CPSR_REG]&(1<<31) != c.r[CPSR_REG]&(1<<28) {
			fmt.Printf("Skipping instruction 0x%08X becaue of GE conditional\n", instruction)
			conditionFailed = true
		}
	case 0b1011 /* LT, less than */ :
		// If the CSPR N flag (bit 31) is not equal to the CSPR V flag (bit 28), then the instruction is executed
		if c.r[CPSR_REG]&(1<<31) == c.r[CPSR_REG]&(1<<28) {
			fmt.Printf("Skipping instruction 0x%08X becaue of LT conditional\n", instruction)
			conditionFailed = true
		}
	case 0b1100 /* GT, greater than */ :
		// If the CSPR Z flag (bit 30) is clear, the CSPR N flag (bit 31) is equal to the CSPR V flag (bit 28), then the instruction is executed
		if c.r[CPSR_REG]&(1<<30) != 0 || c.r[CPSR_REG]&(1<<31) != c.r[CPSR_REG]&(1<<28) {
			fmt.Printf("Skipping instruction 0x%08X becaue of GT conditional\n", instruction)
			conditionFailed = true
		}
	case 0b1101 /* LE, less than or equal */ :
		// If the CSPR Z flag (bit 30) is set, or the CSPR N flag (bit 31) is not equal to the CSPR V flag (bit 28), then the instruction is executed
		if c.r[CPSR_REG]&(1<<30) == 0 || c.r[CPSR_REG]&(1<<31) == c.r[CPSR_REG]&(1<<28) {
			fmt.Printf("Skipping instruction 0x%08X becaue of LE conditional\n", instruction)
			conditionFailed = true
		}
	case 0b1110 /* AL, always */ :
		// The instruction is always executed
		break
	}

	// EXECUTE
	if !conditionFailed {
		oldPC := c.r[PC_REG]

		fmt.Printf("Executing instruction 0x%08X at 0x%08X\n", instruction, c.r[PC_REG])
		instr := arm.DecodeInstruction(instruction)
		if instr != nil {
			repipeline := instr.Execute(c)
			if repipeline {
				c.FlushPipeline()
			} else {
				if oldPC != c.r[PC_REG] {
					if c.config.Debug {
						fmt.Printf("Branching from 0x%08X to 0x%08X, flushing pipeline\n", oldPC, c.r[PC_REG])
					}
					c.FlushPipeline()
				}
			}
		} else {
			fmt.Printf("Unknown instruction 0x%08X\n", instruction)
			panic("")
		}
	}
}

func (c *ARM7TDMI) prettyCPSR() string {
	nMode := "-"
	if c.GetN() {
		nMode = "N"
	}
	zMode := "-"
	if c.GetZ() {
		zMode = "Z"
	}
	cMode := "-"
	if c.GetC() {
		cMode = "C"
	}
	vMode := "-"
	if c.GetV() {
		vMode = "V"
	}
	iMode := "-"
	if c.ReadCPSR()&(1<<7)>>7 == 1 {
		iMode = "I"
	}
	fMode := "-"
	if c.ReadCPSR()&(1<<6)>>6 == 1 {
		fMode = "F"
	}
	thumbMode := "-"
	if c.GetThumbMode() {
		thumbMode = "T"
	}

	opMode := "unknown"
	switch cpuMode(c.ReadCPSR() & 0x1F) {
	case userMode:
		opMode = "user"
	case fiqMode:
		opMode = "fiq"
	case irqMode:
		opMode = "irq"
	case supervisorMode:
		opMode = "supervisor"
	case abortMode:
		opMode = "abort"
	case undefinedMode:
		opMode = "undefined"
	case systemMode:
		opMode = "system"
	}

	return fmt.Sprintf(
		"CPSR: %08X [%s%s%s%s%s%s%s] [Mode: %s]",
		c.ReadCPSR(),
		nMode,
		zMode,
		cMode,
		vMode,
		iMode,
		fMode,
		thumbMode,
		opMode,
	)
}

func (c *ARM7TDMI) stepThumb() {
	// FETCH
	instruction := c.fetchThumb()

	if c.config.TraceRegisters {
		fmt.Printf("\n\nInstruction 0x%04X\n", instruction)
		fmt.Printf("R0:  0x%08X\t\t R1: 0x%08X\t\tR3:   0x%08X\n", c.r[0], c.r[1], c.r[3])
		fmt.Printf("R4:  0x%08X\t\t R5: 0x%08X\t\tR6:   0x%08X\n", c.r[4], c.r[5], c.r[6])
		fmt.Printf("R7:  0x%08X\t\t R8: 0x%08X\t\tR9:   0x%08X\n", c.r[7], c.r[8], c.r[9])
		fmt.Printf("R10: 0x%08X\t\tR11: 0x%08X\t\tR12:  0x%08X\n", c.r[10], c.r[11], c.r[12])
		fmt.Printf("SP:  0x%08X\t\t LR: 0x%08X\t\tPC:   0x%08X\n", c.ReadSP(), c.ReadLR(), c.ReadPC())
		fmt.Println(c.prettyCPSR())
		if cpuMode(c.r[CPSR_REG]&0x1F) != systemMode && cpuMode(c.r[CPSR_REG]&0x1F) != userMode {
			fmt.Printf("SPSR: 0x%08X\n", c.ReadSPSR())
		}
	}

	// DECODE
	instr := thumb.DecodeInstruction(instruction)

	// EXECUTE
	oldPC := c.r[PC_REG]
	if instr != nil {
		repipeline := instr.Execute(c)
		if repipeline {
			c.FlushPipeline()
		} else {
			if oldPC != c.r[PC_REG] {
				if c.config.Debug {
					fmt.Printf("Branching from 0x%08X to 0x%08X, flushing pipeline\n", oldPC, c.r[PC_REG])
				}
				c.FlushPipeline()
			}
		}
	} else {
		fmt.Printf("Unknown instruction 0x%04X\n", instruction)
		panic("")
	}
}

func (c *ARM7TDMI) SetThumbMode(value bool) {
	if value {
		c.r[CPSR_REG] |= 1 << 5
	} else {
		c.r[CPSR_REG] &^= 1 << 5
	}
}

func (c *ARM7TDMI) GetThumbMode() bool {
	if c.r[CPSR_REG]&(1<<5)>>5 == 0 {
		return false
	} else {
		return true
	}
}

// Run runs the CPU at a consistent 16.78MHz
func (c *ARM7TDMI) Run() {
	cycleTime := time.Second / 16777216
	prevTime := time.Now()
	for !c.exit {
		if !c.halted {
			// if c.r[CPSR_REG] bit 5 is set, the CPU is in thumb mode
			if c.r[CPSR_REG]&(1<<5)>>5 == 0 {
				c.stepARM()
			} else {
				c.stepThumb()
			}
		}

		time.Sleep(cycleTime - time.Since(prevTime))
		prevTime = time.Now()
	}
}

func (c *ARM7TDMI) Halt() {
	c.halted = true
}

func (c *ARM7TDMI) Unhalt() {
	c.halted = false
}

func (c *ARM7TDMI) Quit() {
	c.exit = true
}

func (c *ARM7TDMI) SetZ(value bool) {
	// Set bit 30 of CPSR to value
	if value {
		c.r[CPSR_REG] |= 1 << 30
	} else {
		c.r[CPSR_REG] &= ^uint32(1 << 30)
	}
}

func (c *ARM7TDMI) SetN(value bool) {
	// Set bit 31 of CPSR to value
	if value {
		c.r[CPSR_REG] |= 1 << 31
	} else {
		c.r[CPSR_REG] &= ^uint32(1 << 31)
	}
}

func (c *ARM7TDMI) SetC(value bool) {
	// Set bit 29 of CPSR to value
	if value {
		c.r[CPSR_REG] |= 1 << 29
	} else {
		c.r[CPSR_REG] &= ^uint32(1 << 29)
	}
}

func (c *ARM7TDMI) SetV(value bool) {
	// Set bit 28 of CPSR to value
	if value {
		c.r[CPSR_REG] |= 1 << 28
	} else {
		c.r[CPSR_REG] &= ^uint32(1 << 28)
	}
}

func (c *ARM7TDMI) GetZ() bool {
	// Return bit 30 of CPSR
	return c.r[CPSR_REG]&(1<<30)>>30 != 0
}

func (c *ARM7TDMI) GetN() bool {
	// Return bit 31 of CPSR
	return c.r[CPSR_REG]&(1<<31)>>31 != 0
}

func (c *ARM7TDMI) GetC() bool {
	// Return bit 29 of CPSR
	return c.r[CPSR_REG]&(1<<29)>>29 != 0
}

func (c *ARM7TDMI) GetV() bool {
	// Return bit 28 of CPSR
	return c.r[CPSR_REG]&(1<<28)>>28 != 0
}

func (c *ARM7TDMI) SetConditionCodes(res uint32, carry bool, overflow bool) {
	// Need to set CPSR Z, N, C, V
	// Z is set if the result is 0
	c.SetZ(res == 0)

	// N is set if the result is negative
	c.SetN(res&(1<<31)>>31 == 1)

	c.SetC(carry)
	c.SetV(overflow)
}

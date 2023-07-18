package memory

import (
	"fmt"
	"sort"

	"github.com/USA-RedDragon/go-gba/internal/config"
)

type mmioMapping struct {
	address uint32
	size    uint32

	data []byte
}

type MMIO struct {
	mmios  []mmioMapping
	Config *config.Config
}

func (h *MMIO) checkWritable(addr uint32) bool {
	// Addresses 0x00000000 - 0x00003FFF are not writable (BIOS)
	// Addresses 0x02000000 - 0x0203FFFF are writable (on-board WRAM)
	// Addresses 0x03000000 - 0x03007FFF are writable (on-chip WRAM)
	// Addresses 0x04000000 - 0x040003FE are writable (I/O registers)
	// Addresses 0x05000000 - 0x050003FF are writable (palette RAM)
	// Addresses 0x06000000 - 0x06017FFF are writable (VRAM)
	// Addresses 0x07000000 - 0x070003FF are writable (OAM)
	// Addresses 0x08000000 - 0x09FFFFFF are not writable (Game Pak ROM/FlashROM - waitstate 0)
	// Addresses 0x0A000000 - 0x0BFFFFFF are not writable (Game Pak ROM/FlashROM - waitstate 1)
	// Addresses 0x0C000000 - 0x0DFFFFFF are not writable (Game Pak ROM/FlashROM - waitstate 2)
	// Addresses 0x0E000000 - 0x0E00FFFF are writable (Game Pak SRAM)
	writable := true
	if addr >= 0x00000000 && addr < 0x00004000 {
		writable = false
	} else if addr >= 0x02000000 && addr < 0x02040000 {
		writable = true
	} else if addr >= 0x03000000 && addr < 0x03008000 {
		writable = true
	} else if addr >= 0x04000000 && addr < 0x04000400 {
		writable = true
	} else if addr >= 0x05000000 && addr < 0x05000400 {
		writable = true
	} else if addr >= 0x06000000 && addr < 0x06018000 {
		writable = true
	} else if addr >= 0x07000000 && addr < 0x07000400 {
		writable = true
	} else if addr >= 0x08000000 && addr < 0x0A000000 {
		writable = false
	} else if addr >= 0x0A000000 && addr < 0x0C000000 {
		writable = false
	} else if addr >= 0x0C000000 && addr < 0x0E000000 {
		writable = false
	} else if addr >= 0x0E000000 && addr < 0x0E010000 {
		writable = true
	}
	return writable
}

func (h *MMIO) mapMemory(addr uint32) uint32 {
	// account for memory mirroring
	// 0x02040000 - 0x02FFFFFF should map repeatedly to 0x02000000 - 0x0203FFFF
	if addr >= 0x02040000 && addr < 0x03000000 {
		mod := addr % 0x40000
		if h.Config.Debug {
			fmt.Printf("MMIO address 0x%08x mapped to 0x%08x\n", addr, 0x02000000+mod)
		}
		addr = 0x02000000 + mod
	}
	// 0x03008000 - 0x03FFFFFF should map repeatedly to 0x03000000 - 0x03007FFF
	if addr >= 0x03008000 && addr < 0x04000000 {
		mod := addr % 0x8000
		if h.Config.Debug {
			fmt.Printf("MMIO address 0x%08x mapped to 0x%08x\n", addr, 0x03000000+mod)
		}
		addr = 0x03000000 + mod
	}
	// 0x05000400 - 0x05FFFFFF should map repeatedly to 0x05000000 - 0x050003FF
	if addr >= 0x05000400 && addr < 0x06000000 {
		mod := addr % 0x400
		if h.Config.Debug {
			fmt.Printf("MMIO address 0x%08x mapped to 0x%08x\n", addr, 0x05000000+mod)
		}
		addr = 0x05000000 + mod
	}
	// 0x06018000 - 0x06FFFFFF should map repeatedly to 0x06000000 - 0x06017FFF
	if addr >= 0x06018000 && addr < 0x07000000 {
		mod := addr % 0x18000
		if h.Config.Debug {
			fmt.Printf("MMIO address 0x%08x mapped to 0x%08x\n", addr, 0x06000000+mod)
		}
		addr = 0x06000000 + mod
	}
	// 0x07000400 - 0x07FFFFFF should map repeatedly to 0x07000000 - 0x070003FF
	if addr >= 0x07000400 && addr < 0x08000000 {
		mod := addr % 0x400
		if h.Config.Debug {
			fmt.Printf("MMIO address 0x%08x mapped to 0x%08x\n", addr, 0x07000000+mod)
		}
		addr = 0x07000000 + mod
	}
	// 0x0A000000 - 0x0BFFFFFF should map to 0x08000000 - 0x09FFFFFF
	if addr >= 0x0A000000 && addr < 0x0C000000 {
		mod := addr % 0x2000000
		if h.Config.Debug {
			fmt.Printf("MMIO address 0x%08x mapped to 0x%08x\n", addr, 0x08000000+mod)
		}
		addr = 0x08000000 + mod
	}
	// 0x0C000000 - 0x0DFFFFFF should map to 0x08000000 - 0x09FFFFFF
	if addr >= 0x0C000000 && addr < 0x0E000000 {
		mod := addr % 0x2000000
		if h.Config.Debug {
			fmt.Printf("MMIO address 0x%08x mapped to 0x%08x\n", addr, 0x08000000+mod)
		}
		addr = 0x08000000 + mod
	}
	// 0xE000000 - 0xFFFFFFFF should map to 0x0E000000 - 0x0E00FFFF
	if addr >= 0x0E000000 && addr < 0x10000000 {
		mod := addr % 0x10000
		if h.Config.Debug {
			fmt.Printf("MMIO address 0x%08x mapped to 0x%08x\n", addr, 0x0E000000+mod)
		}
		addr = 0x0E000000 + mod
	}

	return addr
}

// findMMIO finds the MMIO device index that contains the given address
func (h *MMIO) findMMIOIndex(addr uint32) (int, error) {
	addr = h.mapMemory(addr)

	if addr >= 0x0E000000 && addr <= 0x0E00FFFF {
		return 9, nil
	} else if addr >= 0x08000000 && addr <= 0x09FFFFFF {
		return 8, nil
	} else if addr >= 0x07000000 && addr <= 0x070003FF {
		return 7, nil
	} else if addr >= 0x06000000 && addr <= 0x06017FFF {
		return 6, nil
	} else if addr >= 0x05000000 && addr <= 0x050003FF {
		return 5, nil
	} else if addr >= 0x04000410 && addr <= 0x04000411 {
		return 4, nil
	} else if addr >= 0x04000000 && addr <= 0x040003FE {
		return 3, nil
	} else if addr >= 0x03000000 && addr <= 0x03007FFF {
		return 2, nil
	} else if addr >= 0x02000000 && addr <= 0x0203FFFF {
		return 1, nil
	} else if addr < 0x00003FFF {
		return 0, nil
	}

	return 0, fmt.Errorf("MMIO address %08x not found", addr)
}

// Read8 reads a 8-bit value from the MMIO address space and returns it.
func (h *MMIO) Read8(addr uint32) (uint8, error) {
	if (addr >= 0x00004000 && addr <= 0x01FFFFFF) || (addr >= 0x10000000 && addr <= 0xFFFFFFFF) {
		return 0, nil
	}
	index, err := h.findMMIOIndex(addr)
	if err != nil {
		return 0, err
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return 0, fmt.Errorf("MMIO address %08x not found", addr)
	}
	return h.mmios[index].data[nonMapped], nil
}

// Write8 writes a 8-bit value to the MMIO address space.
func (h *MMIO) Write8(addr uint32, data uint8) error {
	if (addr >= 0x00004000 && addr <= 0x01FFFFFF) || (addr >= 0x10000000 && addr <= 0xFFFFFFFF) {
		return nil
	}
	if addr >= 0x06000000 && addr < 0x06FFFFFF {
		// VRAM is not byte-addressable
		return fmt.Errorf("VRAM address %08x not byte-addressable", addr)
	}
	index, err := h.findMMIOIndex(addr)
	if err != nil {
		return err
	}
	if h.Config.Debug {
		fmt.Printf("MMIO write: 0x%08x 0x%02x\n", addr, data)
	}
	if !h.checkWritable(addr) {
		return fmt.Errorf("MMIO address %08x not writable", addr)
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return fmt.Errorf("MMIO address %08x not found", addr)
	}
	h.mmios[index].data[nonMapped] = data
	return nil
}

// Read16 reads a 16-bit value from the MMIO address space and returns it.
func (h *MMIO) Read16(addr uint32) (uint16, error) {
	if (addr >= 0x00004000 && addr <= 0x01FFFFFF) || (addr >= 0x10000000 && addr <= 0xFFFFFFFF) {
		return 0, nil
	}
	index, err := h.findMMIOIndex(addr)
	if err != nil {
		return 0, err
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return 0, fmt.Errorf("MMIO address %08x not found", addr)
	}
	dataBytes := h.mmios[index].data[nonMapped : nonMapped+2]
	return uint16(dataBytes[0]) | uint16(dataBytes[1])<<8, nil
}

// Write16 writes a 16-bit value to the MMIO address space.
func (h *MMIO) Write16(addr uint32, data uint16) error {
	if (addr >= 0x00004000 && addr <= 0x01FFFFFF) || (addr >= 0x10000000 && addr <= 0xFFFFFFFF) {
		return nil
	}
	index, err := h.findMMIOIndex(addr)
	if err != nil {
		return err
	}
	if !h.checkWritable(addr) {
		panic(fmt.Errorf("MMIO address %08x not writable", addr))
		return fmt.Errorf("MMIO address %08x not writable", addr)
	}
	if h.Config.Debug {
		fmt.Printf("MMIO write: 0x%08x 0x%04x\n", addr, data)
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return fmt.Errorf("MMIO address %08x not found", addr)
	}
	h.mmios[index].data[nonMapped] = byte(data)
	h.mmios[index].data[nonMapped+1] = byte(data >> 8)
	return nil
}

// Read32 reads a 32-bit value from the MMIO address space and returns it.
func (h *MMIO) Read32(addr uint32) (uint32, error) {
	if (addr >= 0x00004000 && addr <= 0x01FFFFFF) || (addr >= 0x10000000 && addr <= 0xFFFFFFFF) {
		return 0, nil
	}
	index, err := h.findMMIOIndex(addr & ^uint32(3))
	if err != nil {
		return 0, err
	}
	if addr >= 0x03007FFF && addr < 0x04000000 {
		mod := addr % 0x8000
		addr = 0x03000000 + mod
	}
	nonMapped := (addr & ^uint32(3)) - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return 0, fmt.Errorf("MMIO address %08x not found", addr)
	}
	dataBytes := h.mmios[index].data[nonMapped : nonMapped+4]
	val := uint32(dataBytes[0]) | uint32(dataBytes[1])<<8 | uint32(dataBytes[2])<<16 | uint32(dataBytes[3])<<24
	if addr&3 > 0 { // https://github.com/jsmolka/gba-tests/blob/a6447c5404c8fc2898ddc51f438271f832083b7e/thumb/memory.asm#L72
		is := 8 * (uint(addr) & 3)
		is %= 32
		tmp0 := (val) >> (is)
		tmp1 := (val) << (32 - (is))
		return tmp0 | tmp1, nil
	}
	return val, nil
}

// Write32 writes a 32-bit value to the MMIO address space.
func (h *MMIO) Write32(addr uint32, data uint32) error {
	if (addr >= 0x00004000 && addr <= 0x01FFFFFF) || (addr >= 0x10000000 && addr <= 0xFFFFFFFF) {
		return nil
	}
	index, err := h.findMMIOIndex(addr)
	if err != nil {
		return err
	}
	if addr >= 0x03007FFF && addr < 0x04000000 {
		mod := addr % 0x8000
		addr = 0x03000000 + mod
	}
	if h.Config.Debug {
		fmt.Printf("MMIO write: 0x%08x 0x%08x\n", addr, data)
	}
	if !h.checkWritable(addr) {
		return fmt.Errorf("MMIO address %08x not writable", addr)
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return fmt.Errorf("MMIO address %08x not found", addr)
	}
	h.mmios[index].data[nonMapped] = byte(data)
	h.mmios[index].data[nonMapped+1] = byte(data >> 8)
	h.mmios[index].data[nonMapped+2] = byte(data >> 16)
	h.mmios[index].data[nonMapped+3] = byte(data >> 24)
	return nil
}

// AddMMIO adds a new MMIO device to the MMIO handler.
func (h *MMIO) AddMMIO(data []byte, address uint32, size uint32) {
	// Add the MMIO, but ensure that the entries are sorted by address.
	// This is required for the MMIO handler to work properly.

	mapping := mmioMapping{address, size, data}
	h.mmios = append(h.mmios, mapping)

	sort.Slice(h.mmios, func(i, j int) bool {
		return h.mmios[i].address < h.mmios[j].address
	})

	// Print the ordered MMIO list
	fmt.Println("MMIO mappings:")
	for _, mmio := range h.mmios {
		fmt.Printf("  %08x - %08x\n", mmio.address, mmio.address+mmio.size)
	}
}

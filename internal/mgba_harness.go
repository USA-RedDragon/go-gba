package internal

import (
	"bufio"
	"fmt"
	"os/exec"
)

type MGBAHarness struct {
	// The path to the BIOS file
	Bios string
	// The path to the ROM file
	Rom    string
	cmd    *exec.Cmd
	stdout *bufio.Scanner
	stderr *bufio.Scanner
	stdin  *bufio.Writer

	R0          uint32
	R1          uint32
	R2          uint32
	R3          uint32
	R4          uint32
	R5          uint32
	R6          uint32
	R7          uint32
	R8          uint32
	R9          uint32
	R10         uint32
	R11         uint32
	R12         uint32
	R13         uint32
	R14         uint32
	R15         uint32
	CPSR        uint32
	Cycle       uint
	Instruction uint32
}

func NewMGBAHarness(bios, rom string) *MGBAHarness {
	cmd := exec.Command("stdbuf", "-o0", "-i0", "-e0", "mgba", "-d", "-b", bios, rom)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	return &MGBAHarness{
		Bios:   bios,
		Rom:    rom,
		stdout: bufio.NewScanner(stdoutPipe),
		stderr: bufio.NewScanner(stderrPipe),
		stdin:  bufio.NewWriter(stdinPipe),
		cmd:    cmd,
	}
}

func (m *MGBAHarness) Start() error {
	return m.cmd.Start()
}

func (m *MGBAHarness) Stop() error {
	_, err := m.write([]byte("quit\n"))
	if err != nil {
		return err
	}
	return m.cmd.Process.Kill()
}

func (m *MGBAHarness) Wait() error {
	return m.cmd.Wait()
}

//nolint:golint,unparam
func (m *MGBAHarness) write(b []byte) (int, error) {
	i, err := m.stdin.Write(b)
	if err != nil {
		return i, err
	}
	return i, m.stdin.Flush()
}

func (m *MGBAHarness) readLine() (string, error) {
	if !m.stdout.Scan() {
		return "", m.stdout.Err()
	}
	return m.stdout.Text(), nil
}

func (m *MGBAHarness) parseRegisters() error {
	// This function grabs 7 lines from the output and tries to parse them
	// into the registers.
	// Example output:
	//  r0: 03000640   r1: 03001640   r2: 0300056A   r3: D557D557
	//  r4: 00000004   r5: B6D37B40   r6: 00000000   r7: 03000569
	//  r8: 0000001B   r9: 00000000  r10: 0300056A  r11: 00000001
	// r12: 000002E4  r13: 03007E64  r14: 00000004  r15: 00001078
	// cpsr: 2000005F [--C--F-]
	// Cycle: 569256
	//nolint:golint,dupword
	// 00001074:  E2588001     subs r8, r8, #1
	for i := 0; i < 7; i++ {
		line, err := m.readLine()
		if err != nil {
			return err
		}
		switch i {
		case 0:
			_, err := fmt.Sscanf(line, " r0: %X   r1: %X   r2: %X   r3: %X", &m.R0, &m.R1, &m.R2, &m.R3)
			if err != nil {
				return err
			}
		case 1:
			_, err := fmt.Sscanf(line, " r4: %X   r5: %X   r6: %X   r7: %X", &m.R4, &m.R5, &m.R6, &m.R7)
			if err != nil {
				return err
			}
		case 2:
			_, err := fmt.Sscanf(line, " r8: %X   r9: %X  r10: %X  r11: %X", &m.R8, &m.R9, &m.R10, &m.R11)
			if err != nil {
				return err
			}
		case 3:
			_, err := fmt.Sscanf(line, "r12: %X  r13: %X  r14: %X  r15: %X", &m.R12, &m.R13, &m.R14, &m.R15)
			if err != nil {
				return err
			}
		case 4:
			_, err := fmt.Sscanf(line, "cpsr: %X", &m.CPSR)
			if err != nil {
				return err
			}
		case 5:
			_, err := fmt.Sscanf(line, "Cycle: %d", &m.Cycle)
			if err != nil {
				return err
			}
		case 6:
			addr := 0
			_, err := fmt.Sscanf(line, "%X:  %X", &addr, &m.Instruction)
			if err != nil {
				return err
			}
			fmt.Printf("%X\n", m.Instruction)
		}
	}

	return nil
}

func (m *MGBAHarness) GetRegister(reg uint) uint32 {
	switch reg {
	case 0:
		return m.R0
	case 1:
		return m.R1
	case 2:
		return m.R2
	case 3:
		return m.R3
	case 4:
		return m.R4
	case 5:
		return m.R5
	case 6:
		return m.R6
	case 7:
		return m.R7
	case 8:
		return m.R8
	case 9:
		return m.R9
	case 10:
		return m.R10
	case 11:
		return m.R11
	case 12:
		return m.R12
	case 13:
		return m.R13
	case 14:
		return m.R14
	case 15:
		return m.R15
	case 16:
		return m.CPSR
	}
	return 0
}

func (m *MGBAHarness) Step() error {
	fmt.Println("Stepping")
	err := m.parseRegisters()
	if err != nil {
		return err
	}
	_, err = m.write([]byte("next\n"))
	return err
}

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/USA-RedDragon/go-gba/internal"
	"github.com/USA-RedDragon/go-gba/internal/config"
	"github.com/USA-RedDragon/go-gba/internal/emulator"
	"github.com/USA-RedDragon/go-gba/internal/emulator/cpu"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "go-gba [-b gba_bios.bin] [-r path to ROM]",
		RunE:              run,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}

	cmd.Flags().StringP("bios", "b", "", "path to the GBA BIOS")
	cmd.Flags().StringP("rom", "r", "", "path to the GBA ROM")
	cmd.Flags().BoolP("fullscreen", "f", false, "enable fullscreen")
	cmd.Flags().BoolP("trace-registers", "t", false, "trace CPU registers")
	cmd.Flags().BoolP("debug", "d", false, "enable debug logging")
	cmd.Flags().BoolP("diff", "c", false, "enable diff mode")
	cmd.Flags().BoolP("interactive", "i", false, "enable interactive mode, implies --cpu-only and --debug")
	cmd.Flags().Bool("cpu-only", false, "only run the CPU (for debugging)")
	cmd.Flags().Bool("no-gui", false, "disable the GUI (for debugging)")

	return cmd
}

// This mode will run `mgba` in debug mode with a given BIOS and ROM file. It will then
// start up our GBA emulator and step through both, comparing the memory and registers at each step.
// If there is a difference, it will print out the difference and exit.
func runComparer(cmd *cobra.Command) error {
	c := cpu.NewARM7TDMI(config.GetConfig(cmd))
	c.Step()

	mgba := internal.NewMGBAHarness("gba_bios.bin", "arm.gba")
	if err := mgba.Start(); err != nil {
		return err
	}
	defer (func() {
		_ = mgba.Stop()
	})()
	err := mgba.Step()
	for err == nil {
		err = mgba.Step()
		if err != nil {
			break
		}
		registers := uint(16)
		if c.GetThumbMode() {
			registers = 8
		}
		for i := uint(0); i < registers; i++ {
			if mgba.GetRegister(i) != c.ReadRegister(uint8(i)) && !(i > 12) {
				if i == 15 && c.GetThumbMode() {
					// mgba PC will be 2 bytes ahead of ours
					if mgba.GetRegister(i)-2 == c.ReadRegister(uint8(i)) {
						continue
					} else {
						fmt.Printf("Register %d: %08x != %08x\n", i, mgba.GetRegister(i), c.ReadRegister(uint8(i)))
						return nil
					}
				}
				fmt.Printf("Register %d: %08x != %08x\n", i, mgba.GetRegister(i), c.ReadRegister(uint8(i)))
				return nil
			}
			c.Step()
		}
	}
	if err != nil {
		return err
	}

	mgba.Wait()

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	diff, err := cmd.Flags().GetBool("diff")
	if err != nil {
		return err
	}
	if diff {
		return runComparer(cmd)
	}
	cpuOnly, err := cmd.Flags().GetBool("cpu-only")
	if err != nil {
		return err
	}
	interactive, err := cmd.Flags().GetBool("interactive")
	if err != nil {
		return err
	}
	if interactive || cpuOnly {
		c := cpu.NewARM7TDMI(config.GetConfig(cmd))
		if !interactive {
			c.Run()
		} else {
			c.InteractiveRun()
		}
		return nil
	}
	noGUI, err := cmd.Flags().GetBool("no-gui")
	if err != nil {
		return err
	}
	if noGUI {
		c := cpu.NewARM7TDMI(config.GetConfig(cmd))
		c.Run()
		return nil
	} else {
		config := config.GetConfig(cmd)

		emu := emulator.New(config)
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			for range ch {
				fmt.Println("Exiting")
				emu.Stop()
			}
		}()

		ebiten.SetWindowSize(int(config.Scale*240), int(config.Scale*160))
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		ebiten.SetFullscreen(config.Fullscreen)
		ebiten.SetScreenClearedEveryFrame(true)

		if config.ROMPath != "" {
			name := strings.TrimSuffix(filepath.Base(config.ROMPath), filepath.Ext(config.ROMPath))
			ebiten.SetWindowTitle(name + " | go-gba")
		} else {
			ebiten.SetWindowTitle("go-gba")
		}

		if err := ebiten.RunGame(emu); err != nil {
			return err
		}
	}

	return nil
}

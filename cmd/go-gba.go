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

	cmd.Flags().StringP("bios", "b", "gba_bios.bin", "path to the GBA BIOS")
	cmd.Flags().StringP("rom", "r", "", "path to the GBA ROM")
	cmd.Flags().BoolP("fullscreen", "f", false, "enable fullscreen")
	cmd.Flags().BoolP("trace-registers", "t", false, "trace CPU registers")
	cmd.Flags().BoolP("debug", "d", false, "enable debug logging")
	cmd.Flags().BoolP("interactive", "i", false, "enable interactive mode, implies --cpu-only and --debug")
	cmd.Flags().Bool("cpu-only", false, "only run the CPU (for debugging)")
	cmd.Flags().Bool("no-gui", false, "disable the GUI (for debugging)")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
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
			}
		}()

		ebiten.SetWindowSize(int(config.Scale*240), int(config.Scale*160))
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		ebiten.SetFullscreen(config.Fullscreen)
		ebiten.SetScreenClearedEveryFrame(false)

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

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/USA-RedDragon/go-gba/internal/config"
	"github.com/USA-RedDragon/go-gba/internal/emulator"

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

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
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

	return nil
}

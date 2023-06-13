package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/USA-RedDragon/go-gba/internal/config"

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
	fmt.Printf("%v\n", config)

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		for range ch {
			fmt.Println("Exiting")
		}
	}()

	return nil
}

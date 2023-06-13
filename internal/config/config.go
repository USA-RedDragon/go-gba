package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// Config stores the application configuration.
type Config struct {
	BIOSPath       string
	ROMPath        string
	Scale          float64
	TraceRegisters bool
	Debug          bool
	Fullscreen     bool
}

func loadConfigFromEnv() Config {
	scaleStr := os.Getenv("SCALE")
	scale, err := strconv.ParseFloat(scaleStr, 64)
	if err != nil {
		scale = 2.0
	} else if scale < 1.0 {
		scale = 1.0
	}

	tmpConfig := Config{
		BIOSPath:       os.Getenv("BIOS_PATH"),
		ROMPath:        os.Getenv("ROM_PATH"),
		TraceRegisters: os.Getenv("TRACE_REGISTERS") != "",
		Debug:          os.Getenv("DEBUG") != "",
		Scale:          scale,
		Fullscreen:     os.Getenv("FULLSCREEN") != "",
	}

	return tmpConfig
}

// GetConfig obtains the current configuration
func GetConfig(cmd *cobra.Command) *Config {
	currentConfig := loadConfigFromEnv()

	// Override with command line flags
	if cmd != nil {
		biosPath, err := cmd.Flags().GetString("bios")
		if err == nil && biosPath != "" {
			currentConfig.BIOSPath = biosPath
		}

		romPath, err := cmd.Flags().GetString("rom")
		if err == nil && romPath != "" {
			currentConfig.ROMPath = romPath
		}

		traceRegisters, err := cmd.Flags().GetBool("trace-registers")
		if err == nil && traceRegisters {
			currentConfig.TraceRegisters = traceRegisters
		}

		debug, err := cmd.Flags().GetBool("debug")
		if err == nil && debug {
			currentConfig.Debug = debug
		}

		scale, err := cmd.Flags().GetFloat64("scale")
		if err == nil {
			currentConfig.Scale = scale
		}

		fullscreen, err := cmd.Flags().GetBool("fullscreen")
		if err == nil {
			currentConfig.Fullscreen = fullscreen
		}
	}

	fmt.Println(currentConfig.ToString())

	return &currentConfig
}

// ToString returns a string representation of the configuration
func (config *Config) ToString() string {
	return "BIOSPath: " + config.BIOSPath + "\n" +
		"ROMPath: " + config.ROMPath + "\n" +
		"Scale: " + strconv.FormatFloat(config.Scale, 'f', 2, 64) + "\n" +
		"TraceRegisters: " + strconv.FormatBool(config.TraceRegisters) + "\n" +
		"Debug: " + strconv.FormatBool(config.Debug) + "\n" +
		"Fullscreen: " + strconv.FormatBool(config.Fullscreen) + "\n"
}

package main

import (
	"log"

	"github.com/USA-RedDragon/go-gba/cmd"
)

func main() {
	rootCmd := cmd.New()
	rootCmd.Version = "next"
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

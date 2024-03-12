package main

import (
	"log"

	"github.com/USA-RedDragon/go-gba/cmd"
)

// https://goreleaser.com/cookbooks/using-main.version/
//
//nolint:golint,gochecknoglobals
var (
	version = "dev"
	commit  = "none"
)

func main() {
	rootCmd := cmd.New(version, commit)
	rootCmd.Version = "next"
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

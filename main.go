package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/USA-RedDragon/go-gba/cpu"
)

func main() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Starting GBA emulator...")

	cpu := cpu.NewARM7TDMI()
	go cpu.Run()
	<-signalChannel
	cpu.Halt()
	os.Exit(0)
}

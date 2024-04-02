package main

import (
	"os"

	"github.com/richardmarshall/vclfiddle/internal/cmd"
)

func main() {
	if err := cmd.NewDefaultCommand().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

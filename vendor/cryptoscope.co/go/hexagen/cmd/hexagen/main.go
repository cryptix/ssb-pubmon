package main

import (
	"fmt"
	"image/png"
	"os"
	"strings"

	"cryptoscope.co/go/hexagen"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <id>\n", os.Args[0])
		return
	}

	g, err := hexagen.Generate(os.Args[1], 512)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
		return
	}

	// replace slashes, they are not allowed in filesystem context
	f, err := os.Create(strings.Replace(os.Args[1], "/", "|", -1) + ".png")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
		return
	}
	defer f.Close()

	if err := png.Encode(f, g); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
		return
	}

}

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var modes = map[string]func() error{
	"tagstats":     tagStats,
	"fanart":       fanart,
	"thumbAspects": thumbAspects,
	"mpaa":         mpaa,
}

func main() {
	if len(os.Args) != 2 {
		printUsage()
		return
	}

	mode := os.Args[1]
	modeFunc, exists := modes[mode]
	if !exists {
		printUsage()
		return
	}

	if err := modeFunc(); err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	availableModes := make([]string, 0, len(modes))
	for mode := range modes {
		availableModes = append(availableModes, mode)
	}
	fmt.Printf("Error: Invalid or missing mode argument.\nAvailable modes: %s\n", strings.Join(availableModes, ", "))
}

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	total := 0
	filenames := os.Args[1:]
	var hasErrorOccurred bool

	for _, filename := range filenames {
		counts, err := CountFile(filename)
		if err != nil {
			hasErrorOccurred = true
			fmt.Fprintln(os.Stderr, "wc:", err)
			continue
		}

		total += counts.Words
		fmt.Println(counts, filename)
	}

	if len(filenames) == 0 {
		fmt.Println(CountAll(os.Stdin))
		os.Exit(0)
	}

	if len(filenames) > 1 {
		fmt.Println(total, "total")
	}

	if hasErrorOccurred {
		os.Exit(1)
	}
}

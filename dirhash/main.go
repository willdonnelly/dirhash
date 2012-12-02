package main

import (
    "github.com/willdonnelly/dirhash"
	"flag"
	"fmt"
	"os"
)

func main() {
	var hashroot = flag.String("dir", ".", "the directory to generate a cryptographic hash of")
	flag.Parse()

	hash, err := dirhash.HashDir(*hashroot)
	if err != nil {
        fmt.Fprintf(os.Stderr, "error: %s\n", err)
        os.Exit(1)
	}

    fmt.Printf("%X\n", hash)
}

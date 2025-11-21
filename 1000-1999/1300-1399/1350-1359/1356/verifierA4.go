package main

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA4.go /path/to/candidate.go")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	refSrc, err := os.ReadFile(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read reference: %v\n", err)
		os.Exit(1)
	}
	candSrc, err := os.ReadFile(candidate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read candidate: %v\n", err)
		os.Exit(1)
	}

	refFmt := formatSource(refSrc)
	candFmt := formatSource(candSrc)
	if string(refFmt) != string(candFmt) {
		fmt.Fprintln(os.Stderr, "candidate does not match reference stub")
		os.Exit(1)
	}
	fmt.Println("Solution matches reference stub.")
}

func formatSource(src []byte) []byte {
	if formatted, err := format.Source(src); err == nil {
		return formatted
	}
	return src
}

func locateReference() (string, error) {
	candidates := []string{
		"1356A4.go",
		filepath.Join("1000-1999", "1300-1399", "1350-1359", "1356", "1356A4.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1356A4.go reference")
}

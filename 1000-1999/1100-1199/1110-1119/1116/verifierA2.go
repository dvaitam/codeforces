package main

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA2.go /path/to/candidate")
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
		fmt.Fprintf(os.Stderr, "candidate differs from reference stub\n")
		os.Exit(1)
	}
	fmt.Println("Solution matches reference stub.")
}

func formatSource(src []byte) []byte {
	formatted, err := format.Source(src)
	if err != nil {
		return src
	}
	return formatted
}

func locateReference() (string, error) {
	candidates := []string{
		"1116A2.go",
		filepath.Join("1000-1999", "1100-1199", "1110-1119", "1116", "1116A2.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not locate 1116A2.go")
}

package main

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC3.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	ref, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	refSrc, err := os.ReadFile(ref)
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
	if formatted, err := format.Source(src); err == nil {
		return formatted
	}
	return src
}

func locateReference() (string, error) {
	candidates := []string{
		"1116C3.go",
		filepath.Join("1000-1999", "1100-1199", "1110-1119", "1116", "1116C3.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not locate 1116C3.go")
}

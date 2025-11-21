package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const ref1662I = "1000-1999/1600-1699/1660-1669/1662/1662I.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read input:", err)
		os.Exit(1)
	}

	refBin, cleanup, err := buildReference(ref1662I)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}

	if err := compareTokens(refOut, candOut); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "reference output:")
		fmt.Fprintln(os.Stderr, refOut)
		fmt.Fprintln(os.Stderr, "candidate output:")
		fmt.Fprintln(os.Stderr, candOut)
		os.Exit(1)
	}

	fmt.Println("Accepted")
}

func buildReference(src string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-1662I-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	return out.String(), cmd.Run()
}

func compareTokens(expected, got string) error {
	refFields := strings.Fields(expected)
	candFields := strings.Fields(got)
	if len(refFields) != len(candFields) {
		return fmt.Errorf("output length mismatch: expected %d tokens, got %d", len(refFields), len(candFields))
	}
	for i := range refFields {
		if refFields[i] != candFields[i] {
			return fmt.Errorf("mismatch at token %d: expected %s, got %s", i+1, refFields[i], candFields[i])
		}
	}
	return nil
}

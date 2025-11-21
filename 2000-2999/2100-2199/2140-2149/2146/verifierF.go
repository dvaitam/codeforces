package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const ref2146F = "2000-2999/2100-2199/2140-2149/2146/2146F.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read input:", err)
		os.Exit(1)
	}
	testCount, err := readTestCount(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	refBin, cleanup, err := buildReference(ref2146F)
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

	refTokens := strings.Fields(refOut)
	candTokens := strings.Fields(candOut)
	if len(refTokens) != testCount {
		fmt.Fprintf(os.Stderr, "reference output token count mismatch: expected %d, got %d\n", testCount, len(refTokens))
		os.Exit(1)
	}
	if len(candTokens) != testCount {
		fmt.Fprintf(os.Stderr, "candidate output token count mismatch: expected %d, got %d\n", testCount, len(candTokens))
		os.Exit(1)
	}
	for i := 0; i < testCount; i++ {
		if refTokens[i] != candTokens[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %s, got %s\n", i+1, refTokens[i], candTokens[i])
			fmt.Fprintln(os.Stderr, "reference output:")
			fmt.Fprintln(os.Stderr, refOut)
			fmt.Fprintln(os.Stderr, "candidate output:")
			fmt.Fprintln(os.Stderr, candOut)
			os.Exit(1)
		}
	}

	fmt.Println("Accepted")
}

func readTestCount(data []byte) (int, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, fmt.Errorf("failed to read t: %v", err)
	}
	return t, nil
}

func buildReference(src string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-2146F-")
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

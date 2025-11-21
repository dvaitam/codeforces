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

const refSource = "1000-1999/1900-1999/1970-1979/1970/1970G3.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG3.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	inputBytes, err := io.ReadAll(bufio.NewReader(os.Stdin))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read input:", err)
		os.Exit(1)
	}
	input := string(inputBytes)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\ninput:\n%s\n", err, input)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput:\n%soutput:\n%s\n", err, input, candOut)
		os.Exit(1)
	}

	refFields := strings.Fields(refOut)
	candFields := strings.Fields(candOut)
	if len(refFields) != len(candFields) {
		fmt.Fprintf(os.Stderr, "output length mismatch: expected %d tokens, got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
			len(refFields), len(candFields), input, refOut, candOut)
		os.Exit(1)
	}
	for i := range refFields {
		if refFields[i] != candFields[i] {
			fmt.Fprintf(os.Stderr, "mismatch at token %d: expected %s, got %s\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
				i+1, refFields[i], candFields[i], input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Println("Accepted")
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "ref1970G3-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref1970G3.bin")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

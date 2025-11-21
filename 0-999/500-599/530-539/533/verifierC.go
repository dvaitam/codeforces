package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSourceC = "0-999/500-599/530-539/533/533C.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := readAll()
	if err != nil {
		fail("failed to read input: %v", err)
	}

	refBin, err := buildReference(refSourceC)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fail("reference run failed: %v", err)
	}

	userOut, err := runProgram(candidate, input)
	if err != nil {
		fail("candidate run failed: %v", err)
	}

	if normalize(refOut) != normalize(userOut) {
		fail("outputs differ\nexpected:\n%s\ngot:\n%s", refOut, userOut)
	}

	fmt.Println("OK")
}

func readAll() ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(os.Stdin)
	return buf.Bytes(), err
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "533C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(src))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func normalize(out string) string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

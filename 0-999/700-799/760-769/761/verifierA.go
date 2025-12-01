package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSourceA = "./761A.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceA)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateAllCases()
	for idx, input := range tests {
		expect, err := runProgram(refBin, input)
		if err != nil {
			fail("reference failed on case %d: %v\ninput: %s", idx+1, err, strings.TrimSpace(input))
		}
		got, err := runCandidate(candidate, input)
		if err != nil {
			fail("candidate failed on case %d: %v\ninput: %s", idx+1, err, strings.TrimSpace(input))
		}
		if normalize(expect) != normalize(got) {
			fail("case %d mismatch\ninput: %s\nexpected: %s\ngot: %s", idx+1, strings.TrimSpace(input), expect, got)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "761A-ref-*")
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

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	return runCommand(cmd, input)
}

func runCandidate(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	return runCommand(cmd, input)
}

func runCommand(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func generateAllCases() []string {
	tests := make([]string, 0, 10201)
	for a := 0; a <= 100; a++ {
		for b := 0; b <= 100; b++ {
			tests = append(tests, fmt.Sprintf("%d %d\n", a, b))
		}
	}
	return tests
}

func normalize(out string) string {
	return strings.ToUpper(strings.TrimSpace(out))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

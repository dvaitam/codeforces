package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const refPath = "2000-2999/2100-2199/2120-2129/2122/2122A.go"

type testCase struct {
	n, m int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := buildTests()
	input := renderInput(tests)

	expect, err := runBinary(refPath, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\ninput:\n%s\n", err, input)
		os.Exit(1)
	}

	actual, err := runBinary(bin, input)
	if err != nil {
		fmt.Printf("runtime error: %v\ninput:\n%s\n", err, input)
		os.Exit(1)
	}

	expLines := normalizeLines(expect)
	actLines := normalizeLines(actual)
	if len(expLines) != len(tests) {
		fmt.Printf("reference produced %d lines, expected %d\n", len(expLines), len(tests))
		os.Exit(1)
	}
	if len(actLines) != len(tests) {
		fmt.Printf("output has %d lines, expected %d\ninput:\n%s\nactual:\n%s\n", len(actLines), len(tests), input, actual)
		os.Exit(1)
	}

	for i := range tests {
		if expLines[i] != actLines[i] {
			fmt.Printf("case %d mismatch: expected %s got %s\nn=%d m=%d\nfull input:\n%s\n", i+1, expLines[i], actLines[i], tests[i].n, tests[i].m, input)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func buildTests() []testCase {
	rand.Seed(time.Now().UnixNano())
	var tests []testCase

	// Small exhaustive cases
	for n := 1; n <= 4; n++ {
		for m := 1; m <= 4; m++ {
			tests = append(tests, testCase{n: n, m: m})
		}
	}

	// Sample from statement
	tests = append(tests, testCase{n: 3, m: 3}, testCase{n: 1, m: 2})

	// Randomized cases up to max constraints
	for i := 0; i < 50; i++ {
		n := 1 + rand.Intn(100)
		m := 1 + rand.Intn(100)
		tests = append(tests, testCase{n: n, m: m})
	}

	return tests
}

func renderInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	}
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func normalizeLines(out string) []string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i := range lines {
		lines[i] = strings.ToUpper(strings.TrimSpace(lines[i]))
	}
	return lines
}

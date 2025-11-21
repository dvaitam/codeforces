package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "2000-2999/2000-2099/2030-2039/2031/2031C.go"

type testInput struct {
	raw string
	ns  []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	maxN := 0
	for _, tc := range tests {
		for _, n := range tc.ns {
			if n > maxN {
				maxN = n
			}
		}
	}
	squares := precomputeSquares(maxN)

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test #%d: %v\n", idx+1, err)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, tc.raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test #%d: %v\n", idx+1, err)
			os.Exit(1)
		}
		refLines := parseLines(refOut)
		candLines := parseLines(candOut)
		if len(refLines) != len(tc.ns) {
			fmt.Fprintf(os.Stderr, "reference produced %d lines, expected %d\n", len(refLines), len(tc.ns))
			os.Exit(1)
		}
		if len(candLines) != len(tc.ns) {
			fmt.Fprintf(os.Stderr, "candidate produced %d lines, expected %d\n", len(candLines), len(tc.ns))
			os.Exit(1)
		}
		for i, n := range tc.ns {
			if err := verifyCase(n, refLines[i], candLines[i], squares); err != nil {
				fmt.Fprintf(os.Stderr, "test #%d case #%d failed: %v\nInput:\n%s", idx+1, i+1, err, tc.raw)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d test inputs passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2031C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseLines(out string) []string {
	raw := strings.Split(out, "\n")
	res := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		res = append(res, line)
	}
	return res
}

func verifyCase(n int, refLine, candLine string, squares []bool) error {
	if refLine == "-1" {
		fields := strings.Fields(candLine)
		if len(fields) != 1 || fields[0] != "-1" {
			return fmt.Errorf("expected -1, got %q", candLine)
		}
		return nil
	}

	fields := strings.Fields(candLine)
	if len(fields) != n {
		return fmt.Errorf("expected %d integers, got %d", n, len(fields))
	}
	assign := make([]int, n)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid integer %q", f)
		}
		if val < 1 || val > 1_000_000 {
			return fmt.Errorf("value %d out of range", val)
		}
		assign[i] = val
	}
	return validateArrangement(assign, squares)
}

func validateArrangement(arr []int, squares []bool) error {
	posMap := make(map[int][]int)
	for i, v := range arr {
		posMap[v] = append(posMap[v], i+1)
	}
	for val, positions := range posMap {
		if len(positions) == 1 {
			return fmt.Errorf("value %d appears exactly once", val)
		}
		for i := 0; i < len(positions); i++ {
			for j := i + 1; j < len(positions); j++ {
				diff := positions[j] - positions[i]
				if diff >= len(squares) || !squares[diff] {
					return fmt.Errorf("value %d has positions %d and %d with non-square distance %d", val, positions[i], positions[j], diff)
				}
			}
		}
	}
	return nil
}

func precomputeSquares(n int) []bool {
	squares := make([]bool, n+2)
	for i := 1; i*i < len(squares); i++ {
		squares[i*i] = true
	}
	return squares
}

func buildTests() []testInput {
	return []testInput{
		buildTest("2\n3\n12\n"),
		buildTest("5\n1\n2\n3\n4\n5\n"),
		buildTest("3\n6\n8\n10\n"),
		buildTest("2\n100\n200\n"),
	}
}

func buildTest(input string) testInput {
	ns := parseNs(input)
	return testInput{raw: input, ns: ns}
}

func parseNs(input string) []int {
	reader := strings.NewReader(input)
	var t int
	fmt.Fscan(reader, &t)
	ns := make([]int, t)
	for i := 0; i < t; i++ {
		fmt.Fscan(reader, &ns[i])
	}
	return ns
}

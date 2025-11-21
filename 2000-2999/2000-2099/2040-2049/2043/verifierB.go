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

var oddDigits = []int{1, 3, 5, 7, 9}

type seqInfo struct {
	seq   []int
	start int
	cycle int
}

type ndPair struct {
	n int
	d int
}

type testData struct {
	input string
	cases []ndPair
}

var seqTable [10][10]seqInfo
var factThreshold [10]int

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	precomputeSequences()
	precomputeThresholds()

	tests := buildTests()
	for idx, tc := range tests {
		output, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		lines, err := extractLines(output, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output format on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		for i, line := range lines {
			expected := expectedDigits(tc.cases[i].n, tc.cases[i].d)
			if err := validateLine(line, expected); err != nil {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: %v\nInput:\n%s", idx+1, i+1, err, tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	switch filepath.Ext(target) {
	case ".go":
		cmd = exec.Command("go", "run", target)
	case ".py":
		cmd = exec.Command("python3", target)
	default:
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

func extractLines(out string, expected int) ([]string, error) {
	raw := strings.Split(out, "\n")
	lines := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d non-empty lines, got %d", expected, len(lines))
	}
	return lines, nil
}

func validateLine(line string, expected map[int]bool) error {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return fmt.Errorf("empty line")
	}
	seen := make(map[int]bool)
	for _, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid integer %q", f)
		}
		if val < 1 || val > 9 || val%2 == 0 {
			return fmt.Errorf("digit %d is not a valid odd digit", val)
		}
		if seen[val] {
			return fmt.Errorf("digit %d appears multiple times", val)
		}
		seen[val] = true
	}
	if len(seen) != len(expected) {
		return fmt.Errorf("expected digits %s, got %s", formatDigits(expected), formatDigits(seen))
	}
	for d := range expected {
		if !seen[d] {
			return fmt.Errorf("missing digit %d", d)
		}
	}
	return nil
}

func formatDigits(set map[int]bool) string {
	var values []string
	for _, d := range oddDigits {
		if set[d] {
			values = append(values, strconv.Itoa(d))
		}
	}
	return strings.Join(values, " ")
}

func expectedDigits(n, d int) map[int]bool {
	res := map[int]bool{1: true}
	for _, g := range []int{3, 5, 7, 9} {
		if divides(n, d, g) {
			res[g] = true
		}
	}
	return res
}

func divides(n, d, g int) bool {
	if g == 1 {
		return true
	}
	info := seqTable[d][g]
	start := info.start
	seq := info.seq
	if val, ok := smallFactorialValue(n); ok && val <= start {
		return seq[val-1] == 0
	}
	if info.cycle == 1 {
		return seq[start] == 0
	}
	factMod := factorialMod(n, info.cycle)
	startMod := (start + 1) % info.cycle
	offset := factMod - startMod
	offset %= info.cycle
	if offset < 0 {
		offset += info.cycle
	}
	idx := start + offset
	if idx >= len(seq) {
		idx = start + offset%info.cycle
	}
	return seq[idx] == 0
}

func smallFactorialValue(n int) (int, bool) {
	switch n {
	case 2:
		return 2, true
	case 3:
		return 6, true
	default:
		return 0, false
	}
}

func factorialMod(n, mod int) int {
	if mod == 1 {
		return 0
	}
	if n >= factThreshold[mod] {
		return 0
	}
	res := 1 % mod
	for i := 2; i <= n; i++ {
		res = (res * (i % mod)) % mod
	}
	return res
}

func precomputeThresholds() {
	factThreshold[0] = 1
	for m := 1; m <= 9; m++ {
		if m == 1 {
			factThreshold[m] = 1
			continue
		}
		prod := 1
		for n := 1; ; n++ {
			prod *= n
			if prod%m == 0 {
				factThreshold[m] = n
				break
			}
		}
	}
}

func precomputeSequences() {
	for d := 0; d <= 9; d++ {
		for _, g := range oddDigits {
			seqTable[d][g] = buildSequence(d, g)
		}
	}
}

func buildSequence(d, g int) seqInfo {
	seen := make([]int, g)
	for i := range seen {
		seen[i] = -1
	}
	seq := make([]int, 0, g+1)
	state := 0
	seen[state] = 0
	for {
		state = (state*10 + d) % g
		seq = append(seq, state)
		if seen[state] != -1 {
			start := seen[state]
			cycle := len(seq) - start
			return seqInfo{seq: seq, start: start, cycle: cycle}
		}
		seen[state] = len(seq)
	}
}

func buildTests() []testData {
	return []testData{
		buildTest("3\n2 6\n7 1\n8 5\n"),
		buildTest("5\n2 1\n2 9\n3 3\n50 4\n1000000000 5\n"),
		buildTest("4\n2 7\n3 7\n4 7\n6 7\n"),
	}
}

func buildTest(input string) testData {
	reader := strings.NewReader(input)
	var t int
	fmt.Fscan(reader, &t)
	cases := make([]ndPair, t)
	for i := 0; i < t; i++ {
		fmt.Fscan(reader, &cases[i].n, &cases[i].d)
	}
	return testData{input: input, cases: cases}
}

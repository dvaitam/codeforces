package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type testCase struct {
	n    int
	bits string
}

type operation struct {
	typ string
	a   int
	b   int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadDeterministicTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(1002))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		if err := checkCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed (n=%d, b=%s): %v\n", idx+1, tc.n, tc.bits, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func checkCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d\n%s\n", tc.n, tc.bits)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	ops, err := parseOperations(out, tc.n+1)
	if err != nil {
		return err
	}
	return verifyOracle(tc.n, tc.bits, ops)
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOperations(out string, totalQubits int) ([]operation, error) {
	lines := strings.Split(out, "\n")
	nextLine := func(idx *int) (string, bool) {
		for *idx < len(lines) {
			line := strings.TrimSpace(lines[*idx])
			*idx++
			if line != "" {
				return line, true
			}
		}
		return "", false
	}
	idx := 0
	first, ok := nextLine(&idx)
	if !ok {
		return nil, fmt.Errorf("empty output")
	}
	cnt, err := strconv.Atoi(first)
	if err != nil {
		return nil, fmt.Errorf("invalid operation count: %v", err)
	}
	if cnt < 0 {
		return nil, fmt.Errorf("negative operation count")
	}
	ops := make([]operation, 0, cnt)
	for len(ops) < cnt {
		line, ok := nextLine(&idx)
		if !ok {
			return nil, fmt.Errorf("expected %d operations, got %d", cnt, len(ops))
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "X":
			if len(fields) != 2 {
				return nil, fmt.Errorf("invalid X gate format: %q", line)
			}
			target, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, fmt.Errorf("invalid qubit index in %q: %v", line, err)
			}
			if target < 1 || target > totalQubits {
				return nil, fmt.Errorf("qubit index out of range in %q", line)
			}
			ops = append(ops, operation{typ: "X", a: target})
		case "CNOT":
			if len(fields) != 3 {
				return nil, fmt.Errorf("invalid CNOT format: %q", line)
			}
			ctrl, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, fmt.Errorf("invalid control index in %q: %v", line, err)
			}
			tgt, err := strconv.Atoi(fields[2])
			if err != nil {
				return nil, fmt.Errorf("invalid target index in %q: %v", line, err)
			}
			if ctrl < 1 || ctrl > totalQubits || tgt < 1 || tgt > totalQubits {
				return nil, fmt.Errorf("qubit index out of range in %q", line)
			}
			if ctrl == tgt {
				return nil, fmt.Errorf("control and target must differ in %q", line)
			}
			ops = append(ops, operation{typ: "CNOT", a: ctrl, b: tgt})
		default:
			return nil, fmt.Errorf("unsupported operation %q", fields[0])
		}
	}
	for ; idx < len(lines); idx++ {
		if strings.TrimSpace(lines[idx]) != "" {
			return nil, fmt.Errorf("unexpected extra output: %q", strings.TrimSpace(lines[idx]))
		}
	}
	return ops, nil
}

func verifyOracle(n int, bits string, ops []operation) error {
	totalQ := n + 1
	bVec, err := parseBitString(bits)
	if err != nil {
		return err
	}
	totalStates := 1 << totalQ
	for mask := 0; mask < totalStates; mask++ {
		state := maskToBits(totalQ, mask)
		origX := append([]int(nil), state[:n]...)
		origY := state[n]
		for _, op := range ops {
			applyOperation(state, op)
		}
		for i := 0; i < n; i++ {
			if state[i] != origX[i] {
				return fmt.Errorf("basis |x=%s,y=%d> changed qubit %d from %d to %d", bitsToString(origX), origY, i+1, origX[i], state[i])
			}
		}
		expectedY := origY
		for i := 0; i < n; i++ {
			if bVec[i] == 1 && origX[i] == 1 {
				expectedY ^= 1
			}
		}
		if state[n] != expectedY {
			return fmt.Errorf("basis |x=%s,y=%d>: expected output %d got %d", bitsToString(origX), origY, expectedY, state[n])
		}
	}
	return nil
}

func applyOperation(state []int, op operation) {
	switch op.typ {
	case "X":
		state[op.a-1] ^= 1
	case "CNOT":
		if state[op.a-1] == 1 {
			state[op.b-1] ^= 1
		}
	}
}

func parseBitString(bits string) ([]int, error) {
	vec := make([]int, len(bits))
	for i, ch := range bits {
		switch ch {
		case '0':
			vec[i] = 0
		case '1':
			vec[i] = 1
		default:
			return nil, fmt.Errorf("invalid bit %q in bitstring", ch)
		}
	}
	return vec, nil
}

func bitsToString(bits []int) string {
	var sb strings.Builder
	for _, v := range bits {
		if v == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return sb.String()
}

func maskToBits(total int, mask int) []int {
	bits := make([]int, total)
	for i := 0; i < total; i++ {
		if mask&(1<<i) != 0 {
			bits[i] = 1
		}
	}
	return bits
}

func loadDeterministicTests() ([]testCase, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("cannot determine verifier path")
	}
	path := filepath.Join(filepath.Dir(file), "testcasesD.txt")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tests []testCase
	scanner := bufio.NewScanner(f)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid testcase format on line %d", lineNo)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("invalid N on line %d: %v", lineNo, err)
		}
		bits := fields[1]
		if n < 1 || n > 8 {
			return nil, fmt.Errorf("N out of range on line %d", lineNo)
		}
		if len(bits) != n {
			return nil, fmt.Errorf("bitstring length mismatch on line %d", lineNo)
		}
		tests = append(tests, testCase{n: n, bits: bits})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return testCase{n: n, bits: sb.String()}
}

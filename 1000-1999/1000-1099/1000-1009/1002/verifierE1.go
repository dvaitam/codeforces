package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	tolerance = 1e-6
)

var invSqrt2 = 1 / math.Sqrt(2)

type operation struct {
	kind  string
	qubit int // -1 for oracle
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests, err := loadTests("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, n := range tests {
		input := fmt.Sprintf("%d\n", n)
		out, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (n=%d): runtime error: %v\n", idx+1, n, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		ops, err := parseOperations(out, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (n=%d): output parse error: %v\n", idx+1, n, err)
			fmt.Println("Candidate output:")
			fmt.Print(out)
			os.Exit(1)
		}
		if err := validateOperations(ops, n); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (n=%d) failed: %v\n", idx+1, n, err)
			fmt.Println("Candidate output:")
			fmt.Print(out)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func loadTests(path string) ([]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	cases := make([]int, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid N %q: %v", line, err)
		}
		cases = append(cases, n)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func runProgram(bin, input string) (string, error) {
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

func parseOperations(output string, n int) ([]operation, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanLines)

	var header string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		header = line
		break
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if header == "" {
		return nil, fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(header)
	if err != nil {
		return nil, fmt.Errorf("failed to parse operation count: %v", err)
	}
	lines := make([]string, 0, k)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(lines) != k {
		return nil, fmt.Errorf("expected %d operations, got %d", k, len(lines))
	}
	ops := make([]operation, 0, k)
	for idx, line := range lines {
		parts := strings.Fields(line)
		switch len(parts) {
		case 1:
			if strings.ToUpper(parts[0]) != "ORACLE" {
				return nil, fmt.Errorf("line %d: unknown operation %q", idx+2, line)
			}
			ops = append(ops, operation{kind: "ORACLE", qubit: -1})
		case 2:
			gate := strings.ToUpper(parts[0])
			if gate != "X" && gate != "H" {
				return nil, fmt.Errorf("line %d: unsupported gate %q", idx+2, gate)
			}
			qidx, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid qubit index %q", idx+2, parts[1])
			}
			if qidx < 1 || qidx > n+1 {
				return nil, fmt.Errorf("line %d: qubit index %d out of range [1,%d]", idx+2, qidx, n+1)
			}
			ops = append(ops, operation{kind: gate, qubit: qidx - 1})
		default:
			return nil, fmt.Errorf("line %d: malformed operation %q", idx+2, line)
		}
	}
	return ops, nil
}

func validateOperations(ops []operation, n int) error {
	totalQ := n + 1
	dim := 1 << totalQ
	for secret := 0; secret < (1 << n); secret++ {
		state := make([]complex128, dim)
		state[0] = 1
		for _, op := range ops {
			switch op.kind {
			case "X":
				applyX(state, op.qubit)
			case "H":
				applyH(state, op.qubit)
			case "ORACLE":
				state = applyOracle(state, n, secret)
			default:
				return fmt.Errorf("unsupported operation %s", op.kind)
			}
		}
		if err := checkMeasurement(state, n, secret); err != nil {
			return fmt.Errorf("secret %s: %w", formatBits(secret, n), err)
		}
	}
	return nil
}

func applyX(state []complex128, qubit int) {
	step := 1 << qubit
	for block := 0; block < len(state); block += 2 * step {
		for i := 0; i < step; i++ {
			a := block + i
			b := a + step
			state[a], state[b] = state[b], state[a]
		}
	}
}

func applyH(state []complex128, qubit int) {
	step := 1 << qubit
	for block := 0; block < len(state); block += 2 * step {
		for i := 0; i < step; i++ {
			a := block + i
			b := a + step
			pa := state[a]
			pb := state[b]
			state[a] = (pa + pb) * complex(invSqrt2, 0)
			state[b] = (pa - pb) * complex(invSqrt2, 0)
		}
	}
}

func applyOracle(state []complex128, n int, secret int) []complex128 {
	size := len(state)
	result := make([]complex128, size)
	inputMask := (1 << n) - 1
	for idx, amp := range state {
		x := idx & inputMask
		y := (idx >> n) & 1
		parity := bits.OnesCount(uint(x&secret)) & 1
		newY := y ^ parity
		newIdx := x | (newY << n)
		result[newIdx] += amp
	}
	return result
}

func checkMeasurement(state []complex128, n int, secret int) error {
	inputMask := (1 << n) - 1
	probs := make([]float64, 1<<n)
	for idx, amp := range state {
		x := idx & inputMask
		probs[x] += absSquared(amp)
	}
	target := secret
	if math.Abs(probs[target]-1) > tolerance {
		return fmt.Errorf("probability for %s is %.4f (expected 1)", formatBits(target, n), probs[target])
	}
	for i, p := range probs {
		if i == target {
			continue
		}
		if p > tolerance {
			return fmt.Errorf("non-zero probability %.4e for state %s", p, formatBits(i, n))
		}
	}
	return nil
}

func absSquared(c complex128) float64 {
	return real(c)*real(c) + imag(c)*imag(c)
}

func formatBits(x, n int) string {
	return fmt.Sprintf("%0*b", n, x)
}

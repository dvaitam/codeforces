package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/bits"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	name  string
	input string
}

type constraint struct {
	mask        int
	allowed     int
	size        int
	allowedBits int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}
	input := string(inputBytes)

	parsedInput, err := parseInput(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse input: %v\n", err)
		os.Exit(1)
	}
	validMask, solutions := enumerateSolutions(parsedInput)

	candOut, err := runProgram(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "solution runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(strings.TrimSpace(candOut), len(solutions))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse solution output: %v\n", err)
		os.Exit(1)
	}

	for i := 0; i < len(candAns); i++ {
		if i > 0 && candAns[i] <= candAns[i-1] {
			fmt.Fprintf(os.Stderr, "solutions must be in strictly increasing order\n")
			os.Exit(1)
		}
		if candAns[i] < 0 || candAns[i] >= len(validMask) {
			fmt.Fprintf(os.Stderr, "solution %d out of range\n", candAns[i])
			os.Exit(1)
		}
		if !validMask[candAns[i]] {
			fmt.Fprintf(os.Stderr, "reported solution %d is invalid\n", candAns[i])
			os.Exit(1)
		}
	}

	fmt.Println("Accepted")
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

type parsedInput struct {
	n           int
	values      []int
	constraints []constraint
}

func parseInput(input string) (*parsedInput, error) {
	reader := strings.NewReader(input)
	in := bufio.NewReader(reader)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return nil, err
	}
	total := 1 << n
	values := make([]int, total)
	for mask := 1; mask < total; mask++ {
		var val int
		if _, err := fmt.Fscan(in, &val); err != nil {
			return nil, err
		}
		values[mask] = val
	}
	var constraints []constraint
	for mask := 1; mask < total; mask++ {
		val := values[mask]
		size := bits.OnesCount(uint(mask))
		full := (1 << (size + 1)) - 1
		if val == full {
			continue
		}
		allowedBits := bits.OnesCount(uint(val))
		constraints = append(constraints, constraint{
			mask:        mask,
			allowed:     val,
			size:        size,
			allowedBits: allowedBits,
		})
	}
	sort.Slice(constraints, func(i, j int) bool {
		if constraints[i].size != constraints[j].size {
			return constraints[i].size < constraints[j].size
		}
		if constraints[i].allowedBits != constraints[j].allowedBits {
			return constraints[i].allowedBits < constraints[j].allowedBits
		}
		return constraints[i].mask < constraints[j].mask
	})
	return &parsedInput{
		n:           n,
		values:      values,
		constraints: constraints,
	}, nil
}

func enumerateSolutions(data *parsedInput) ([]bool, []int) {
	total := 1 << data.n
	popcount := make([]uint8, total)
	for i := 1; i < total; i++ {
		popcount[i] = popcount[i>>1] + uint8(i&1)
	}
	valid := make([]bool, total)
	solutions := make([]int, 0)
constraintsLoop:
	for mask := 0; mask < total; mask++ {
		for _, c := range data.constraints {
			count := int(popcount[mask&c.mask])
			if ((c.allowed >> count) & 1) == 0 {
				continue constraintsLoop
			}
		}
		valid[mask] = true
		solutions = append(solutions, mask)
	}
	return valid, solutions
}

func parseOutput(out string, expected int) ([]int, error) {
	if strings.TrimSpace(out) == "" {
		return nil, fmt.Errorf("empty output")
	}
	reader := bufio.NewReader(strings.NewReader(out))
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return nil, fmt.Errorf("failed to read number of solutions: %v", err)
	}
	if expected >= 0 && k != expected {
		return nil, fmt.Errorf("number of solutions mismatch: got %d, expected %d", k, expected)
	}
	ans := make([]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &ans[i]); err != nil {
			return nil, fmt.Errorf("failed to read solution %d: %v", i+1, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("failed to parse trailing output: %v", err)
		}
	} else {
		return nil, fmt.Errorf("unexpected token %q after reading solutions", extra)
	}
	return ans, nil
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

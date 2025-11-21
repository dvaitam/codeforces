package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const ref2152E = "2000-2999/2100-2199/2150-2159/2152/2152E.go"

type seqTest struct {
	n    int
	perm []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read input:", err)
		os.Exit(1)
	}
	tests, err := parseSeqInput(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	refBin, cleanup, err := buildReference(ref2152E)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()
	if _, err := runProgram(refBin, input); err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := verifySequences(candOut, tests); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "candidate output:")
		fmt.Fprintln(os.Stderr, candOut)
		os.Exit(1)
	}

	fmt.Println("Accepted")
}

func parseSeqInput(data []byte) ([]seqTest, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	tests := make([]seqTest, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return nil, fmt.Errorf("test %d: failed to read n: %v", i+1, err)
		}
		m := n*n + 1
		perm := make([]int, m)
		for j := 0; j < m; j++ {
			if _, err := fmt.Fscan(reader, &perm[j]); err != nil {
				return nil, fmt.Errorf("test %d: failed to read permutation value %d: %v", i+1, j+1, err)
			}
		}
		tests[i] = seqTest{n: n, perm: perm}
	}
	return tests, nil
}

func verifySequences(out string, tests []seqTest) error {
	tokens := strings.Fields(out)
	idx := 0
	for ti, tc := range tests {
		need := tc.n + 1
		if idx+need > len(tokens) {
			return fmt.Errorf("test %d: expected %d indices, got %d", ti+1, need, len(tokens)-idx)
		}
		indices := make([]int, need)
		for j := 0; j < need; j++ {
			val, err := strconv.Atoi(tokens[idx+j])
			if err != nil {
				return fmt.Errorf("test %d: invalid integer at position %d: %v", ti+1, j+1, err)
			}
			if val < 1 || val > len(tc.perm) {
				return fmt.Errorf("test %d: index %d out of range", ti+1, val)
			}
			indices[j] = val - 1
			if j > 0 && indices[j] <= indices[j-1] {
				return fmt.Errorf("test %d: indices must be strictly increasing", ti+1)
			}
		}
		idx += need
		values := make([]int, need)
		for j, pos := range indices {
			values[j] = tc.perm[pos]
		}
		if !isStrictlyMonotone(values) {
			return fmt.Errorf("test %d: sequence %v is not strictly monotone", ti+1, values)
		}
	}
	if idx != len(tokens) {
		return fmt.Errorf("unexpected extra tokens starting at %q", tokens[idx])
	}
	return nil
}

func isStrictlyMonotone(vals []int) bool {
	if len(vals) <= 1 {
		return true
	}
	inc := true
	for i := 1; i < len(vals); i++ {
		if vals[i] <= vals[i-1] {
			inc = false
			break
		}
	}
	if inc {
		return true
	}
	dec := true
	for i := 1; i < len(vals); i++ {
		if vals[i] >= vals[i-1] {
			dec = false
			break
		}
	}
	return dec
}

func buildReference(src string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-2152E-")
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

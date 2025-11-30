package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded gzipped+base64 testcases from testcasesB.txt.
const encodedTestcases = `
H4sIAA8IK2kC/z1O2w0AMQj6d5pefND9JzultDFEqmD51rLqcksDWRKwILLhzWayNfFmg1GCztI+dANUzW501dj0gMqgd97nj5DPX4rTU56bBuT17tx8oYw/0ydTX8wAAAA=
`

// verifyOutput ensures each line is a valid anti-Fibonacci permutation.
func verifyOutput(n int, out string) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != n {
		return fmt.Errorf("expected %d permutations got %d", n, len(lines))
	}
	seen := make(map[string]bool)
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != n {
			return fmt.Errorf("line %d: expected %d numbers got %d", i+1, n, len(fields))
		}
		perm := make([]int, n)
		used := make([]bool, n+1)
		for j, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return fmt.Errorf("line %d: invalid number", i+1)
			}
			if v < 1 || v > n || used[v] {
				return fmt.Errorf("line %d: invalid permutation", i+1)
			}
			used[v] = true
			perm[j] = v
			if j >= 2 && perm[j-2]+perm[j-1] == v {
				return fmt.Errorf("line %d: not anti-Fibonacci", i+1)
			}
		}
		key := strings.Join(fields, " ")
		if seen[key] {
			return fmt.Errorf("duplicate permutation on line %d", i+1)
		}
		seen[key] = true
	}
	return nil
}

// decodeTestcases returns the list of n values.
func decodeTestcases() ([]int, error) {
	data, err := base64.StdEncoding.DecodeString(encodedTestcases)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var out bytes.Buffer
	if _, err := out.ReadFrom(r); err != nil {
		return nil, err
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	var res []int
	for i, f := range fields {
		if i == 0 {
			continue // skip t
		}
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}

func runCandidate(bin string, n int) (string, error) {
	input := fmt.Sprintf("1\n%d\n", n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := decodeTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, n := range tests {
		out, err := runCandidate(bin, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := verifyOutput(n, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

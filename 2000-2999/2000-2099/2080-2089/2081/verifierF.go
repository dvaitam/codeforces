package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const refSourceF = "./2081F.go"

type testCase struct {
	n    int
	name string
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refYN, err := parseReference(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := validateCandidate(candOut, tests, refYN); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	outPath := "./ref_2081F.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSourceF)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(name string, n int) {
		tests = append(tests, testCase{name: name, n: n})
	}

	add("n1", 1)
	add("n2", 2)
	add("n3", 3)
	add("n4", 4)
	add("n5", 5)
	add("n6", 6)
	add("n7", 7)
	add("n8", 8)
	add("n10", 10)
	add("n12", 12)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	total := 0
	const maxSumN = 2800
	for len(tests) < 40 && total < maxSumN {
		n := rng.Intn(100) + 1
		if len(tests)%10 == 0 {
			n = rng.Intn(800) + 200
		}
		if total+n > maxSumN {
			break
		}
		add(fmt.Sprintf("random_%d", len(tests)), n)
		total += n
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	}
	return sb.String()
}

func parseReference(out string, tests []testCase) ([]bool, error) {
	tokens := strings.Fields(out)
	pos := 0
	ans := make([]bool, len(tests))
	for i, tc := range tests {
		if pos >= len(tokens) {
			return nil, fmt.Errorf("reference output ended early at case %d", i+1)
		}
		tn := strings.ToUpper(tokens[pos])
		pos++
		if tn == "NO" {
			ans[i] = false
			continue
		}
		if tn != "YES" {
			return nil, fmt.Errorf("case %d: expected YES/NO, got %q", i+1, tokens[pos-1])
		}
		ans[i] = true
		need := tc.n * tc.n
		if pos+need > len(tokens) {
			return nil, fmt.Errorf("case %d: expected %d matrix entries, got %d", i+1, need, len(tokens)-pos)
		}
		pos += need
	}
	return ans, nil
}

func validateCandidate(out string, tests []testCase, refYN []bool) error {
	tokens := strings.Fields(out)
	pos := 0
	for i, tc := range tests {
		if pos >= len(tokens) {
			return fmt.Errorf("candidate output ended early at case %d", i+1)
		}
		tn := strings.ToUpper(tokens[pos])
		pos++
		if !refYN[i] {
			if tn != "NO" {
				return fmt.Errorf("case %d (%s): expected NO (hot matrix impossible), got %q", i+1, tc.name, tokens[pos-1])
			}
			continue
		}
		if tn != "YES" {
			return fmt.Errorf("case %d (%s): expected YES, got %q", i+1, tc.name, tokens[pos-1])
		}
		need := tc.n * tc.n
		if pos+need > len(tokens) {
			return fmt.Errorf("case %d (%s): expected %d matrix entries, got %d", i+1, tc.name, need, len(tokens)-pos)
		}
		values := make([]int, need)
		for k := 0; k < need; k++ {
			v, err := strconv.Atoi(tokens[pos+k])
			if err != nil {
				return fmt.Errorf("case %d (%s): invalid integer %q", i+1, tc.name, tokens[pos+k])
			}
			values[k] = v
		}
		if err := checkMatrix(tc.n, values); err != nil {
			return fmt.Errorf("case %d (%s): %v", i+1, tc.name, err)
		}
		pos += need
	}
	if pos != len(tokens) {
		return fmt.Errorf("extra output tokens detected (%d unused)", len(tokens)-pos)
	}
	return nil
}

func checkMatrix(n int, vals []int) error {
	if len(vals) != n*n {
		return fmt.Errorf("matrix size mismatch, need %d values got %d", n*n, len(vals))
	}

	mat := vals
	colFreq := make([][]int, n)
	for j := 0; j < n; j++ {
		colFreq[j] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		rowFreq := make([]int, n)
		for j := 0; j < n; j++ {
			v := mat[i*n+j]
			if v < 0 || v >= n {
				return fmt.Errorf("value out of range at (%d,%d): %d", i+1, j+1, v)
			}
			rowFreq[v]++
			if rowFreq[v] > 1 {
				return fmt.Errorf("row %d is not a permutation (value %d repeats)", i+1, v)
			}
			colFreq[j][v]++
			if colFreq[j][v] > 1 && i < n-1 {
				// Early detection
				return fmt.Errorf("column %d has duplicate value %d", j+1, v)
			}
		}
	}
	for j := 0; j < n; j++ {
		for v := 0; v < n; v++ {
			if colFreq[j][v] != 1 {
				return fmt.Errorf("column %d missing value %d", j+1, v)
			}
		}
	}

	seenH := make([]bool, n*n)
	seenV := make([]bool, n*n)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := mat[i*n+j]
			if v+mat[i*n+(n-1-j)] != n-1 {
				return fmt.Errorf("horizontal symmetry failed at (%d,%d)", i+1, j+1)
			}
			if v+mat[(n-1-i)*n+j] != n-1 {
				return fmt.Errorf("vertical symmetry failed at (%d,%d)", i+1, j+1)
			}
			if j+1 < n {
				v2 := mat[i*n+j+1]
				idx := v*n + v2
				if seenH[idx] {
					return fmt.Errorf("duplicate horizontal pair (%d,%d) found", v, v2)
				}
				seenH[idx] = true
			}
			if i+1 < n {
				v2 := mat[(i+1)*n+j]
				idx := v*n + v2
				if seenV[idx] {
					return fmt.Errorf("duplicate vertical pair (%d,%d) found", v, v2)
				}
				seenV[idx] = true
			}
		}
	}

	return nil
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource    = "./2064B.go"
	totalNLimit  = 40000
	totalTCLimit = 300
)

type testCase struct {
	n int
	a []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	// Run reference once to ensure test input is solvable.
	if out, err := runProgram(refBin, input); err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, out)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	if err := validateOutput(candOut, tests); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2064B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", tc.n)
		for i, v := range tc.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func validateOutput(out string, tests []testCase) error {
	tokens := strings.Fields(out)
	expect := 0
	for _, tc := range tests {
		if bestLen(tc) == 0 {
			expect++
		} else {
			expect += 2
		}
	}
	if len(tokens) < expect {
		return fmt.Errorf("not enough output tokens: expected at least %d, got %d", expect, len(tokens))
	}

	idx := 0
	for caseIdx, tc := range tests {
		bLen := bestLen(tc)
		if bLen == 0 {
			val := tokens[idx]
			idx++
			if val != "0" {
				return fmt.Errorf("test %d: expected 0, got %q", caseIdx+1, val)
			}
			continue
		}

		lVal, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return fmt.Errorf("test %d: failed to parse l: %v", caseIdx+1, err)
		}
		idx++
		rVal, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return fmt.Errorf("test %d: failed to parse r: %v", caseIdx+1, err)
		}
		idx++

		if lVal < 1 || rVal < lVal || rVal > tc.n {
			return fmt.Errorf("test %d: invalid segment [%d,%d]", caseIdx+1, lVal, rVal)
		}

		if rVal-lVal+1 != bLen {
			return fmt.Errorf("test %d: segment length %d not optimal (expected %d)", caseIdx+1, rVal-lVal+1, bLen)
		}

		freq := make([]int, tc.n+1)
		for _, v := range tc.a {
			freq[v]++
		}
		for i := lVal - 1; i < rVal; i++ {
			if freq[tc.a[i]] != 1 {
				return fmt.Errorf("test %d: segment contains value %d with frequency %d (must be 1)", caseIdx+1, tc.a[i], freq[tc.a[i]])
			}
		}
	}

	if idx != len(tokens) {
		return fmt.Errorf("extra output tokens detected: used %d of %d", idx, len(tokens))
	}
	return nil
}

func bestLen(tc testCase) int {
	freq := make([]int, tc.n+1)
	for _, v := range tc.a {
		freq[v]++
	}
	best := 0
	cur := 0
	for _, v := range tc.a {
		if freq[v] == 1 {
			cur++
		} else {
			if cur > best {
				best = cur
			}
			cur = 0
		}
	}
	if cur > best {
		best = cur
	}
	return best
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	sumN := 0

	add := func(tc testCase) {
		if sumN+tc.n > totalNLimit || len(tests) >= totalTCLimit {
			return
		}
		tests = append(tests, tc)
		sumN += tc.n
	}

	// Deterministic edge cases.
	add(testCase{n: 1, a: []int{1}})                      // all unique
	add(testCase{n: 4, a: []int{2, 2, 2, 2}})             // no unique
	add(testCase{n: 5, a: []int{1, 2, 1, 2, 3}})          // single unique at end
	add(testCase{n: 6, a: []int{1, 2, 3, 1, 2, 3}})       // none removable helps
	add(testCase{n: 6, a: []int{1, 4, 5, 6, 1, 2}})       // unique block in middle
	add(testCase{n: 7, a: []int{1, 2, 3, 4, 1, 2, 3}})    // multiple unique length 1
	add(testCase{n: 8, a: []int{1, 2, 3, 4, 5, 1, 2, 3}}) // best is long prefix

	for len(tests) < totalTCLimit && sumN < totalNLimit {
		n := rng.Intn(300) + 1
		if sumN+n > totalNLimit {
			n = totalNLimit - sumN
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			mode := rng.Intn(5)
			switch mode {
			case 0:
				a[i] = 1
			case 1:
				a[i] = n
			default:
				a[i] = rng.Intn(n) + 1
			}
		}

		// Occasionally enforce large unique blocks.
		if rng.Intn(4) == 0 {
			start := rng.Intn(n)
			end := start + rng.Intn(n-start)
			val := n + 1
			for i := start; i <= end; i++ {
				val++
				a[i] = val
			}
		}

		add(testCase{n: n, a: a})
	}

	if len(tests) == 0 {
		tests = append(tests, testCase{n: 1, a: []int{1}})
	}
	return tests
}

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
	refSource   = "./2094F.go"
	maxTests    = 300
	cellLimit   = 40000
	randomCases = 250
)

type testCase struct {
	n int
	m int
	k int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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

	// Ensure the reference can handle the generated tests.
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
	tmp, err := os.CreateTemp("", "2094F-ref-*")
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
		fmt.Fprintf(&b, "%d %d %d\n", tc.n, tc.m, tc.k)
	}
	return b.String()
}

func validateOutput(out string, tests []testCase) error {
	tokens := strings.Fields(out)
	idx := 0
	for caseIdx, tc := range tests {
		totalCells := tc.n * tc.m
		if idx+totalCells > len(tokens) {
			return fmt.Errorf("test %d: not enough values, need %d more", caseIdx+1, totalCells-(len(tokens)-idx))
		}
		grid := make([][]int, tc.n)
		for i := 0; i < tc.n; i++ {
			grid[i] = make([]int, tc.m)
			for j := 0; j < tc.m; j++ {
				val, err := strconv.Atoi(tokens[idx])
				if err != nil {
					return fmt.Errorf("test %d: failed to parse integer at position %d: %v", caseIdx+1, idx+1, err)
				}
				grid[i][j] = val
				idx++
			}
		}

		if err := checkGrid(tc, grid); err != nil {
			return fmt.Errorf("test %d: %v", caseIdx+1, err)
		}
	}
	if idx != len(tokens) {
		return fmt.Errorf("extra output tokens detected: used %d of %d", idx, len(tokens))
	}
	return nil
}

func checkGrid(tc testCase, grid [][]int) error {
	cntPerValue := make([]int, tc.k+1)
	for i := 0; i < tc.n; i++ {
		if len(grid[i]) != tc.m {
			return fmt.Errorf("row %d has wrong length: got %d expected %d", i+1, len(grid[i]), tc.m)
		}
		for j := 0; j < tc.m; j++ {
			v := grid[i][j]
			if v < 1 || v > tc.k {
				return fmt.Errorf("cell (%d,%d) out of range: %d", i+1, j+1, v)
			}
			cntPerValue[v]++
			if i+1 < tc.n && grid[i+1][j] == v {
				return fmt.Errorf("adjacent cells (%d,%d) and (%d,%d) share value %d", i+1, j+1, i+2, j+1, v)
			}
			if j+1 < tc.m && grid[i][j+1] == v {
				return fmt.Errorf("adjacent cells (%d,%d) and (%d,%d) share value %d", i+1, j+1, i+1, j+2, v)
			}
		}
	}

	expected := (tc.n * tc.m) / tc.k
	for val := 1; val <= tc.k; val++ {
		if cntPerValue[val] != expected {
			return fmt.Errorf("value %d appears %d times, expected %d", val, cntPerValue[val], expected)
		}
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	sumCells := 0

	add := func(tc testCase) {
		if len(tests) >= maxTests || sumCells+tc.n*tc.m > cellLimit {
			return
		}
		tests = append(tests, tc)
		sumCells += tc.n * tc.m
	}

	// Deterministic edge cases.
	add(testCase{n: 1, m: 2, k: 2})
	add(testCase{n: 2, m: 1, k: 2})
	add(testCase{n: 2, m: 2, k: 2})
	add(testCase{n: 2, m: 2, k: 4})
	add(testCase{n: 3, m: 4, k: 6})
	add(testCase{n: 4, m: 4, k: 4})

	for attempts := 0; attempts < randomCases && len(tests) < maxTests && sumCells < cellLimit; attempts++ {
		n := rng.Intn(25) + 1
		m := rng.Intn(25) + 1
		if n*m == 0 {
			continue
		}
		if sumCells+n*m > cellLimit {
			break
		}

		// Choose a divisor k of n*m that is at least 2.
		product := n * m
		divs := divisors(product, rng)
		k := divs[rng.Intn(len(divs))]

		add(testCase{n: n, m: m, k: k})
	}

	if len(tests) == 0 {
		tests = append(tests, testCase{n: 1, m: 2, k: 2})
	}
	return tests
}

func divisors(x int, rng *rand.Rand) []int {
	var cand []int
	for d := 1; d*d <= x; d++ {
		if x%d == 0 {
			cand = append(cand, d)
			if d != x/d {
				cand = append(cand, x/d)
			}
		}
	}
	// filter k>=2
	dst := cand[:0]
	for _, v := range cand {
		if v >= 2 {
			dst = append(dst, v)
		}
	}
	if len(dst) == 0 {
		return []int{2}
	}
	// shuffle for variety
	rng.Shuffle(len(dst), func(i, j int) { dst[i], dst[j] = dst[j], dst[i] })
	return dst
}

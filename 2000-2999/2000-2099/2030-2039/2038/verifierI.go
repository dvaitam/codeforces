package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	refSource       = "2000-2999/2000-2099/2030-2039/2038/2038I.go"
	perTestBitLimit = 2_000_000
	totalBitLimit   = 1_800_000
)

type testCase struct {
	n    int
	m    int
	rows []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
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

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	expected := tokenize(refOut)

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	got := tokenize(candOut)

	if len(expected) != len(got) {
		fmt.Fprintf(os.Stderr, "wrong number of tokens: expected %d got %d\n", len(expected), len(got))
		os.Exit(1)
	}
	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "mismatch at token %d: expected %q got %q\n", i+1, expected[i], got[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2038I-ref-*")
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

func tokenize(s string) []string {
	return strings.Fields(strings.TrimSpace(s))
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalBits := 0

	addCase := func(n, m int, rows []string) {
		if rows == nil || len(rows) != n {
			panic("invalid rows for test case")
		}
		tests = append(tests, testCase{n: n, m: m, rows: rows})
		totalBits += n * m
	}

	tryAdd := func(n, m int, rows []string) bool {
		if n < 2 || m < 2 {
			return false
		}
		if n*m > perTestBitLimit {
			return false
		}
		if totalBits+n*m > totalBitLimit {
			return false
		}
		if m <= 60 {
			if limit := 1 << uint(m); n > limit {
				return false
			}
		}
		if rows == nil {
			rows = uniqueBinaryStrings(n, m, rng)
		}
		addCase(n, m, rows)
		return true
	}

	// Sample test from the statement.
	tryAdd(3, 5, []string{"10010", "01100", "10101"})
	// Small structured cases.
	tryAdd(2, 2, []string{"10", "01"})
	tryAdd(4, 4, []string{"0001", "0010", "0100", "1000"})
	tryAdd(8, 3, enumeratedStrings(8, 3))

	for attempts := 0; attempts < 2000 && totalBits < totalBitLimit; attempts++ {
		remaining := totalBitLimit - totalBits
		if remaining < 4 {
			break
		}
		maxM := 60
		if remaining < 200 {
			maxM = 20
		}
		m := rng.Intn(maxM-1) + 2
		if m > 25 && remaining < m*10 {
			m = 2 + rng.Intn(24)
		}
		maxN := remaining / m
		if maxN > 1000 {
			maxN = 1000
		}
		if m <= 20 {
			limit := 1 << uint(m)
			if maxN > limit {
				maxN = limit
			}
		}
		if maxN < 2 {
			continue
		}
		n := 2 + rng.Intn(maxN-1)
		if !tryAdd(n, m, nil) {
			continue
		}
	}

	if len(tests) == 0 {
		tryAdd(2, 2, []string{"10", "01"})
	}

	return tests
}

func uniqueBinaryStrings(n, m int, rng *rand.Rand) []string {
	if m <= 20 && n > (1<<uint(m))/2 {
		return enumeratedStrings(n, m)
	}
	res := make([]string, 0, n)
	seen := make(map[string]struct{}, n*2)
	for len(res) < n {
		bytes := make([]byte, m)
		for i := range bytes {
			if rng.Intn(2) == 1 {
				bytes[i] = '1'
			} else {
				bytes[i] = '0'
			}
		}
		s := string(bytes)
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		res = append(res, s)
	}
	return res
}

func enumeratedStrings(n, m int) []string {
	limit := 1 << uint(m)
	if n > limit {
		n = limit
	}
	res := make([]string, 0, n)
	format := fmt.Sprintf("%%0%db", m)
	for i := 0; i < n; i++ {
		res = append(res, fmt.Sprintf(format, i))
	}
	return res
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.m)
		for _, row := range tc.rows {
			fmt.Fprintln(&b, row)
		}
	}
	return b.String()
}

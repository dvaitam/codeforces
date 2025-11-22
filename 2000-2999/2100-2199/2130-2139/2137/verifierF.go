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
)

const (
	candUsage     = "usage: go run verifierF.go /path/to/binary"
	maxTests      = 140
	maxTotalN     = 20000
	maxTotalNSq   = 5_000_000
	refSource2137 = "2137F.go"
	refBinary2137 = "ref2137F.bin"
)

type testCase struct {
	n int
	a []int64
	b []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(candUsage)
		return
	}
	candidate := os.Args[1]

	// Build reference if available; otherwise skip. This keeps the verifier usable
	// even when the reference source is absent or empty.
	var refBinary string
	if _, err := os.Stat(refSource2137); err == nil {
		refBinary, _ = buildReference()
		if refBinary != "" {
			defer os.Remove(refBinary)
		}
	}

	tests := generateTests()
	input := formatInput(tests)
	expected := computeAll(tests)

	// Optionally validate the reference output when the binary exists.
	if refBinary != "" {
		if out, err := runProgram(refBinary, input); err == nil {
			if refAns, err := parseOutput(out, len(tests)); err == nil {
				if err := compareAnswers(refAns, expected); err != nil {
					fmt.Printf("reference failed validation: %v\n", err)
					return
				}
			}
		}
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}
	if err := compareAnswers(got, expected); err != nil {
		fmt.Printf("candidate failed validation: %v\n", err)
		fmt.Println("Input used:")
		fmt.Println(string(input))
		return
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2137, refSource2137)
	if out, err := cmd.CombinedOutput(); err != nil {
		// return empty string to indicate reference unavailable
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2137), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func compareAnswers(got, expected []int64) error {
	if len(got) != len(expected) {
		return fmt.Errorf("answer length mismatch")
	}
	for i := range got {
		if got[i] != expected[i] {
			return fmt.Errorf("test %d: expected %d, got %d", i+1, expected[i], got[i])
		}
	}
	return nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func computeAll(tests []testCase) []int64 {
	res := make([]int64, len(tests))
	for i, tc := range tests {
		res[i] = solve(tc)
	}
	return res
}

func solve(tc testCase) int64 {
	var total int64
	n := tc.n
	a := tc.a
	b := tc.b
	for l := 0; l < n; l++ {
		p := a[l]
		var cnt int64
		if b[l] == p {
			cnt++
		}
		total += cnt
		for r := l + 1; r < n; r++ {
			if a[r] > p {
				p = a[r]
				if b[r] == p {
					cnt++
				}
			} else {
				if b[r] <= p {
					cnt++
				}
			}
			total += cnt
		}
	}
	return total
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2137))
	var tests []testCase
	totalN := 0
	totalNSq := 0

	add := func(a, b []int64) {
		if len(a) != len(b) || len(a) == 0 {
			return
		}
		n := len(a)
		if len(tests) >= maxTests || totalN+n > maxTotalN || totalNSq+n*n > maxTotalNSq {
			return
		}
		tc := testCase{
			n: n,
			a: append([]int64(nil), a...),
			b: append([]int64(nil), b...),
		}
		tests = append(tests, tc)
		totalN += n
		totalNSq += n * n
	}

	// Deterministic cases
	add([]int64{1}, []int64{1})
	add([]int64{5, 3, 1}, []int64{4, 2, 1})
	add([]int64{7, 1, 1, 2}, []int64{10, 5, 8, 9})
	add([]int64{1, 2, 3, 4, 5}, []int64{5, 4, 3, 2, 1})
	add([]int64{2, 2, 2}, []int64{1, 2, 3})

	for len(tests) < maxTests && totalN < maxTotalN && totalNSq < maxTotalNSq {
		remainN := maxTotalN - totalN
		n := rnd.Intn(300) + 1
		if n > remainN {
			n = remainN
		}
		if totalNSq+n*n > maxTotalNSq {
			if totalNSq >= maxTotalNSq {
				break
			}
			for totalNSq+n*n > maxTotalNSq && n > 1 {
				n--
			}
			if n <= 0 {
				break
			}
		}

		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			limit := int64(2*n + 5)
			a[i] = rnd.Int63n(limit) + 1
			b[i] = rnd.Int63n(limit) + 1
		}
		add(a, b)
	}
	return tests
}

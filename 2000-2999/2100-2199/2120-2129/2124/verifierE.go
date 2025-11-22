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
	refSource2124E = "2124E.go"
	refBinary2124E = "ref2124E.bin"
	maxTests       = 150
	maxTotalN      = 45000
	maxOpsLimit    = 17
)

type testCase struct {
	n int
	a []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	refCounts, err := parseCounts(refOut, tests)
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	if err := validateCandidate(candOut, tests, refCounts); err != nil {
		fmt.Printf("candidate failed validation: %v\n", err)
		fmt.Println("Input used:")
		fmt.Println(string(input))
		return
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2124E, refSource2124E)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2124E), nil
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

func parseCounts(out string, tests []testCase) ([]int, error) {
	fields := strings.Fields(out)
	idx := 0
	counts := make([]int, len(tests))
	for t, tc := range tests {
		if idx >= len(fields) {
			return nil, fmt.Errorf("missing s for test %d", t+1)
		}
		s, err := strconv.Atoi(fields[idx])
		if err != nil {
			return nil, fmt.Errorf("test %d: cannot read s: %v", t+1, err)
		}
		idx++
		if s != -1 && (s < 1 || s > maxOpsLimit) {
			return nil, fmt.Errorf("test %d: invalid s %d", t+1, s)
		}
		counts[t] = s
		if s == -1 {
			continue
		}
		need := s * tc.n
		if idx+need > len(fields) {
			return nil, fmt.Errorf("test %d: not enough tokens for operations", t+1)
		}
		idx += need
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("unexpected extra tokens (%d remaining)", len(fields)-idx)
	}
	return counts, nil
}

func validateCandidate(out string, tests []testCase, refCounts []int) error {
	fields := strings.Fields(out)
	idx := 0
	if len(refCounts) != len(tests) {
		return fmt.Errorf("internal error: count mismatch")
	}
	for t, tc := range tests {
		if idx >= len(fields) {
			return fmt.Errorf("test %d: missing s", t+1)
		}
		s, err := strconv.Atoi(fields[idx])
		if err != nil {
			return fmt.Errorf("test %d: cannot parse s: %v", t+1, err)
		}
		idx++

		refS := refCounts[t]
		if refS == -1 && s != -1 {
			return fmt.Errorf("test %d: expected -1, got %d", t+1, s)
		}
		if refS != -1 && s == -1 {
			return fmt.Errorf("test %d: candidate printed -1 but solution exists", t+1)
		}
		if s == -1 {
			continue
		}
		if s != refS {
			return fmt.Errorf("test %d: s mismatch, expected %d, got %d", t+1, refS, s)
		}
		if s < 1 || s > maxOpsLimit {
			return fmt.Errorf("test %d: invalid s %d", t+1, s)
		}

		cur := append([]int64(nil), tc.a...)
		for op := 0; op < s; op++ {
			if idx+tc.n > len(fields) {
				return fmt.Errorf("test %d: insufficient tokens for operation %d", t+1, op+1)
			}
			b := make([]int64, tc.n)
			var total int64
			for i := 0; i < tc.n; i++ {
				val, err := strconv.ParseInt(fields[idx+i], 10, 64)
				if err != nil {
					return fmt.Errorf("test %d op %d idx %d: parse error: %v", t+1, op+1, i+1, err)
				}
				if val < 0 || val > cur[i] {
					return fmt.Errorf("test %d op %d idx %d: value %d out of range [0,%d]", t+1, op+1, i+1, val, cur[i])
				}
				b[i] = val
				total += val
			}
			idx += tc.n

			if total%2 != 0 {
				return fmt.Errorf("test %d op %d: sum %d is odd", t+1, op+1, total)
			}
			half := total / 2
			var pref int64
			found := false
			for i := 0; i < tc.n-1; i++ {
				pref += b[i]
				if pref == half {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("test %d op %d: no valid split", t+1, op+1)
			}

			for i := 0; i < tc.n; i++ {
				cur[i] -= b[i]
			}
		}

		for i, v := range cur {
			if v != 0 {
				return fmt.Errorf("test %d: position %d remains %d after %d ops", t+1, i+1, v, s)
			}
		}
	}
	if idx != len(fields) {
		return fmt.Errorf("extra tokens after processing all tests (%d remaining)", len(fields)-idx)
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
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2124))
	var tests []testCase
	totalN := 0

	add := func(a []int64) {
		if len(a) < 2 {
			return
		}
		if totalN+len(a) > maxTotalN || len(tests) >= maxTests {
			return
		}
		tc := testCase{n: len(a), a: append([]int64(nil), a...)}
		tests = append(tests, tc)
		totalN += len(a)
	}

	// Deterministic cases from statement and variants.
	add([]int64{1, 2, 3})
	add([]int64{2, 2})
	add([]int64{5, 3, 1, 5})
	add([]int64{3, 1, 1})
	add([]int64{10, 1})
	add([]int64{4, 4, 4, 4})
	add([]int64{7, 7, 1, 1})

	for len(tests) < maxTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		n := rnd.Intn(800) + 2
		if n > remain {
			n = remain
		}
		if n < 2 {
			break
		}

		a := make([]int64, n)
		mode := rnd.Intn(4)
		for i := 0; i < n; i++ {
			switch mode {
			case 0:
				a[i] = int64(rnd.Intn(5) + 1)
			case 1:
				a[i] = int64(rnd.Intn(1000) + 1)
			case 2:
				a[i] = int64(rnd.Intn(1_000_000) + 1)
			default:
				a[i] = rnd.Int63n(1_000_000_000_000) + 1
			}
		}
		add(a)
	}

	return tests
}

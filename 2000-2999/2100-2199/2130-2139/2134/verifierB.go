package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const refSource = "2134B.go"

type testCase struct {
	name string
	n    int
	k    int64
	a    []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate_binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	input := buildInput(tests)

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\ninput:\n%soutput:\n%s", err, input, refOut)
		os.Exit(1)
	}
	if err := validateOutputs(tests, refOut); err != nil {
		fmt.Fprintf(os.Stderr, "reference produced invalid output: %v\ninput:\n%soutput:\n%s", err, input, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput:\n%soutput:\n%s", err, input, candOut)
		os.Exit(1)
	}
	if err := validateOutputs(tests, candOut); err != nil {
		fmt.Fprintf(os.Stderr, "candidate output invalid: %v\ninput:\n%soutput:\n%s", err, input, candOut)
		os.Exit(1)
	}

	fmt.Printf("All %d testcases passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot locate verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2134B-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 32)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func validateOutputs(tests []testCase, output string) error {
	tokens := strings.Fields(output)
	idx := 0
	for caseIdx, tc := range tests {
		if idx+tc.n > len(tokens) {
			return fmt.Errorf("case %d: expected %d numbers, got %d", caseIdx+1, tc.n, len(tokens)-idx)
		}
		maxInc := int64(0)
		g := int64(0)
		limit := int64(1_000_000_000) + tc.k*tc.k
		for i := 0; i < tc.n; i++ {
			val, err := strconv.ParseInt(tokens[idx+i], 10, 64)
			if err != nil {
				return fmt.Errorf("case %d: invalid integer %q", caseIdx+1, tokens[idx+i])
			}
			if val < 1 || val > limit {
				return fmt.Errorf("case %d: value %d out of allowed range [1,%d]", caseIdx+1, val, limit)
			}
			inc := val - tc.a[i]
			if inc < 0 || inc%tc.k != 0 {
				return fmt.Errorf("case %d: value %d not reachable from %d with step %d", caseIdx+1, val, tc.a[i], tc.k)
			}
			steps := inc / tc.k
			if steps > tc.k {
				return fmt.Errorf("case %d: requires %d operations (>k=%d)", caseIdx+1, steps, tc.k)
			}
			if steps > maxInc {
				maxInc = steps
			}
			if g == 0 {
				g = val
			} else {
				g = gcd64(g, val)
			}
		}
		idx += tc.n
		if g <= 1 {
			return fmt.Errorf("case %d: gcd is %d, must be >1", caseIdx+1, g)
		}
		if maxInc > tc.k {
			return fmt.Errorf("case %d: operations needed %d exceed k=%d", caseIdx+1, maxInc, tc.k)
		}
	}
	if idx != len(tokens) {
		return fmt.Errorf("expected %d numbers, got %d", idx, len(tokens))
	}
	return nil
}

func gcd64(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func buildTests() []testCase {
	tests := []testCase{
		{
			name: "single_element",
			n:    1,
			k:    1,
			a:    []int64{1},
		},
		{
			name: "small_mixed",
			n:    3,
			k:    3,
			a:    []int64{2, 7, 1},
		},
		{
			name: "already_good",
			n:    4,
			k:    5,
			a:    []int64{10, 15, 20, 25},
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	for len(tests) < 80 && totalN < 100000 {
		n := rng.Intn(5000) + 1
		if totalN+n > 100000 {
			n = 100000 - totalN
		}
		k := int64(rng.Intn(1_000_000_000) + 1)
		a := make([]int64, n)
		for i := range a {
			a[i] = int64(rng.Intn(1_000_000_000) + 1)
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("rnd_%d", len(tests)+1),
			n:    n,
			k:    k,
			a:    a,
		})
		totalN += n
	}
	return tests
}

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

const referenceSolutionRel = "1000-1999/1900-1999/1910-1919/1912/1912H.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "1912H.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	name string
	n    int
	m    int
	pair [][2]int
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow((tc.n+tc.m)*8 + 32)
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for _, e := range tc.pair {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "1912H-ref-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_1912H")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parsePlan(out string, n int) (int, [][2]int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, nil, fmt.Errorf("empty output")
	}
	k64, err := strconv.ParseInt(fields[0], 10, 32)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse k: %v", err)
	}
	k := int(k64)
	if k == -1 {
		return -1, nil, nil
	}
	if k < 0 {
		return 0, nil, fmt.Errorf("k must be non-negative or -1, got %d", k)
	}
	expected := 1 + 2*k
	if len(fields) != expected {
		return 0, nil, fmt.Errorf("expected %d tokens, got %d", expected, len(fields))
	}
	edges := make([][2]int, k)
	for i := 0; i < k; i++ {
		a, err := strconv.Atoi(fields[1+2*i])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid city index %q", fields[1+2*i])
		}
		b, err := strconv.Atoi(fields[2+2*i])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid city index %q", fields[2+2*i])
		}
		if a < 1 || a > n || b < 1 || b > n {
			return 0, nil, fmt.Errorf("city out of range in edge %d: %d %d", i+1, a, b)
		}
		edges[i] = [2]int{a, b}
	}
	return k, edges, nil
}

func validatePlan(tc testCase, k int, edges [][2]int) error {
	if k != len(edges) {
		return fmt.Errorf("reported k=%d but got %d edges", k, len(edges))
	}
	n := tc.n
	next := make([]int, n)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		next[i] = -1
		order[i] = -1
	}
	for idx, e := range edges {
		a := e[0] - 1
		b := e[1] - 1
		if a == b {
			return fmt.Errorf("edge %d uses the same city as source and destination", idx+1)
		}
		if next[a] != -1 {
			return fmt.Errorf("city %d fires more than once", a+1)
		}
		next[a] = b
		order[a] = idx
	}

	for _, p := range tc.pair {
		start := p[0] - 1
		target := p[1] - 1
		cur := start
		lastOrder := -1
		for steps := 0; steps <= n; steps++ {
			if cur == target {
				break
			}
			if next[cur] == -1 {
				return fmt.Errorf("passenger %d->%d stuck at city %d", start+1, target+1, cur+1)
			}
			if order[cur] <= lastOrder {
				return fmt.Errorf("passenger %d->%d needs edge from %d after it has fired", start+1, target+1, cur+1)
			}
			lastOrder = order[cur]
			cur = next[cur]
		}
		if cur != target {
			return fmt.Errorf("passenger %d->%d does not reach destination", start+1, target+1)
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "sample_possible", n: 5, m: 6, pair: [][2]int{{1, 3}, {1, 2}, {2, 3}, {4, 2}, {1, 5}, {5, 1}}},
		{name: "sample_impossible", n: 3, m: 6, pair: [][2]int{{1, 2}, {1, 3}, {2, 1}, {2, 3}, {3, 1}, {3, 2}}},
		{name: "no_passengers", n: 4, m: 0, pair: nil},
		{name: "single_chain", n: 4, m: 3, pair: [][2]int{{1, 4}, {2, 4}, {3, 4}}},
		{name: "cycle_requirement", n: 4, m: 4, pair: [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 1}}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 120; i++ {
		n := rng.Intn(30) + 2
		m := rng.Intn(n*3 + 1)
		pairs := make([][2]int, m)
		for j := 0; j < m; j++ {
			a := rng.Intn(n) + 1
			b := rng.Intn(n-1) + 1
			if b >= a {
				b++
				if b > n {
					b = 1
				}
			}
			pairs[j] = [2]int{a, b}
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", i+1),
			n:    n,
			m:    m,
			pair: pairs,
		})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, input)
			os.Exit(1)
		}
		refK, _, err := parsePlan(refOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid reference output on test %d (%s): %v\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		out, err := runProgram(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, input)
			os.Exit(1)
		}
		k, edges, err := parsePlan(out, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d (%s): %v\n%s", idx+1, tc.name, err, out)
			os.Exit(1)
		}

		if refK == -1 {
			if k != -1 {
				fmt.Fprintf(os.Stderr, "test %d (%s): expected -1 (impossible), got k=%d\n", idx+1, tc.name, k)
				os.Exit(1)
			}
			continue
		}

		if k == -1 {
			fmt.Fprintf(os.Stderr, "test %d (%s): solution exists with k=%d but output claims impossible\n", idx+1, tc.name, refK)
			os.Exit(1)
		}
		if k != refK {
			fmt.Fprintf(os.Stderr, "test %d (%s): expected minimal k=%d, got k=%d\n", idx+1, tc.name, refK, k)
			os.Exit(1)
		}
		if err := validatePlan(tc, k, edges); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s): invalid plan: %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, out)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}

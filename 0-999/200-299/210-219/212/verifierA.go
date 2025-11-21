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

type testCase struct {
	n, m, k, t int
	edges      [][2]int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-212A-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleA")
	cmd := exec.Command("go", "build", "-o", outPath, "212A.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.n, tc.m, tc.k, tc.t))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
	}
	return sb.String()
}

func parseOutput(output string) (int, []int, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return 0, nil, fmt.Errorf("expected at least 2 lines, got: %s", output)
	}
	uneven, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return 0, nil, fmt.Errorf("invalid unevenness: %v", err)
	}
	assignStr := strings.TrimSpace(lines[1])
	parts := strings.Fields(assignStr)
	assign := make([]int, len(parts))
	for i, p := range parts {
		val, err := strconv.Atoi(p)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid assignment value: %s", p)
		}
		assign[i] = val
	}
	return uneven, assign, nil
}

func computeUnevenness(tc testCase, assign []int) (int, error) {
	if len(assign) != tc.k {
		return 0, fmt.Errorf("assignment length mismatch: expected %d, got %d", tc.k, len(assign))
	}
	degLeft := make([][]int, tc.n)
	degRight := make([][]int, tc.m)
	for i := 0; i < tc.n; i++ {
		degLeft[i] = make([]int, tc.t)
	}
	for i := 0; i < tc.m; i++ {
		degRight[i] = make([]int, tc.t)
	}
	for i := 0; i < tc.k; i++ {
		comp := assign[i] - 1
		if comp < 0 || comp >= tc.t {
			return 0, fmt.Errorf("invalid company index %d at edge %d", assign[i], i+1)
		}
		u := tc.edges[i][0] - 1
		v := tc.edges[i][1] - 1
		if u < 0 || u >= tc.n || v < 0 || v >= tc.m {
			return 0, fmt.Errorf("invalid edge endpoints (%d,%d)", tc.edges[i][0], tc.edges[i][1])
		}
		degLeft[u][comp]++
		degRight[v][comp]++
	}
	sum := 0
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.t; j++ {
			sum += degLeft[i][j] * degLeft[i][j]
		}
	}
	for i := 0; i < tc.m; i++ {
		for j := 0; j < tc.t; j++ {
			sum += degRight[i][j] * degRight[i][j]
		}
	}
	return sum, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 3, m: 5, k: 8, t: 2,
			edges: [][2]int{{1, 4}, {1, 3}, {3, 3}, {1, 2}, {1, 1}, {2, 1}, {1, 5}, {2, 2}},
		},
		{
			n: 1, m: 1, k: 1, t: 1,
			edges: [][2]int{{1, 1}},
		},
		{
			n: 2, m: 2, k: 2, t: 2,
			edges: [][2]int{{1, 1}, {2, 2}},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 50)
	sumK := 0
	for len(tests) < 50 && sumK < 5000 {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		k := rng.Intn(10) + 1
		t := rng.Intn(3) + 1
		if sumK+k > 5000 {
			k = 5000 - sumK
			if k <= 0 {
				break
			}
		}
		edges := make([][2]int, k)
		used := make(map[[2]int]bool)
		for i := 0; i < k; i++ {
			var u, v int
			for {
				u = rng.Intn(n) + 1
				v = rng.Intn(m) + 1
				if !used[[2]int{u, v}] {
					used[[2]int{u, v}] = true
					break
				}
			}
			edges[i] = [2]int{u, v}
		}
		tests = append(tests, testCase{n: n, m: m, k: k, t: t, edges: edges})
		sumK += k
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	input := buildInput(tests)

	expected, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
		os.Exit(1)
	}
	actual, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}

	expLines := strings.Split(strings.TrimSpace(expected), "\n")
	actLines := strings.Split(strings.TrimSpace(actual), "\n")
	if len(expLines) != len(actLines) {
		fmt.Fprintf(os.Stderr, "line count mismatch: expected %d, got %d\n", len(expLines), len(actLines))
		os.Exit(1)
	}

	lineIdx := 0
	for idx, tc := range tests {
		expUneven, _, err := parseOutput(strings.Join(expLines[lineIdx:lineIdx+2], "\n"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		actUneven, actAssign, err := parseOutput(strings.Join(actLines[lineIdx:lineIdx+2], "\n"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		lineIdx += 2
		computedUneven, err := computeUnevenness(tc, actAssign)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid assignment: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if computedUneven != actUneven {
			fmt.Fprintf(os.Stderr, "test %d: unevenness mismatch (declared %d, computed %d)\n", idx+1, actUneven, computedUneven)
			os.Exit(1)
		}
		if actUneven < expUneven {
			fmt.Fprintf(os.Stderr, "test %d: unevenness %d better than oracle %d (unexpected)\n", idx+1, actUneven, expUneven)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}

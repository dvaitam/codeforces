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

type edge struct {
	u int
	v int
	c int
}

type testCase struct {
	n     int
	edges []edge
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-720B-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", bin, "720B.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return bin, cleanup, nil
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
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.c))
	}
	return sb.String()
}

func parseAnswer(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	if val < 0 {
		return 0, fmt.Errorf("negative answer %d", val)
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 2,
			edges: []edge{
				{u: 1, v: 2, c: 1},
			},
		},
		{
			n: 3,
			edges: []edge{
				{1, 2, 1},
				{2, 3, 2},
				{3, 1, 3},
			},
		},
		{
			n: 4,
			edges: []edge{
				{1, 2, 1},
				{2, 3, 2},
				{3, 1, 3},
				{3, 4, 4},
			},
		},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(40) + 2
	edges := generateCactus(n, rng)
	assignColors(edges, rng)
	return testCase{n: n, edges: edges}
}

func generateCactus(n int, rng *rand.Rand) []edge {
	edges := make([]edge, 0, n+5)
	if n < 2 {
		return edges
	}
	nodeCount := 1
	for nodeCount < n {
		u := rng.Intn(nodeCount) + 1
		addCycle := false
		maxNew := n - nodeCount
		if maxNew >= 2 && rng.Intn(3) == 0 {
			addCycle = true
		}
		if addCycle {
			maxLen := maxNew + 1
			if maxLen > 6 {
				maxLen = 6
			}
			if maxLen >= 3 {
				L := rng.Intn(maxLen-2) + 3 // [3, maxLen]
				prev := u
				for i := 0; i < L-1; i++ {
					nodeCount++
					newV := nodeCount
					edges = append(edges, edge{u: prev, v: newV})
					prev = newV
				}
				edges = append(edges, edge{u: prev, v: u})
				continue
			}
		}
		nodeCount++
		edges = append(edges, edge{u: u, v: nodeCount})
	}
	return edges
}

func assignColors(edges []edge, rng *rand.Rand) {
	m := len(edges)
	if m == 0 {
		return
	}
	for i := range edges {
		edges[i].c = rng.Intn(m) + 1
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)

		expStr, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		exp, err := parseAnswer(expStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expStr)
			os.Exit(1)
		}

		gotStr, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		got, err := parseAnswer(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotStr, input)
			os.Exit(1)
		}

		if got != exp {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d got %d\ninput:\n%s\n", idx+1, exp, got, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

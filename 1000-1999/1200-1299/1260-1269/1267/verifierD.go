package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	refSource        = "1267D.go"
	tempOraclePrefix = "oracle-1267D-"
	randomTestCount  = 40
	maxEdges         = 264
	maxRandomN       = 60
)

type testInput struct {
	n    int
	want [][]int
	pass [][]int
}

type cdEdge struct {
	u int
	v int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, randomTestCount)...)

	for idx, tc := range tests {
		input := formatInput(tc)
		oracleOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		oracleStatus, err := parseStatus(oracleOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d: %v\noutput:\n%s", idx+1, err, oracleOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		candStatus, ct, edges, err := parseCandidate(candOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\noutput:\n%s", idx+1, err, candOut)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		if oracleStatus == "impossible" {
			if candStatus != "impossible" {
				fmt.Fprintf(os.Stderr, "test %d: oracle says Impossible but candidate says Possible\n", idx+1)
				fmt.Println("Input:")
				fmt.Print(input)
				fmt.Println("Candidate output:")
				fmt.Print(candOut)
				os.Exit(1)
			}
			continue
		}

		if candStatus != "possible" {
			fmt.Fprintf(os.Stderr, "test %d: oracle says Possible but candidate says Impossible\n", idx+1)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		if err := validateSolution(tc, ct, edges); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed validation: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Candidate output:")
			fmt.Print(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
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

func parseStatus(out string) (string, error) {
	reader := strings.NewReader(out)
	var status string
	if _, err := fmt.Fscan(reader, &status); err != nil {
		return "", fmt.Errorf("failed to read status: %v", err)
	}
	status = strings.ToLower(status)
	if status != "possible" && status != "impossible" {
		return "", fmt.Errorf("unexpected status %q", status)
	}
	return status, nil
}

func parseCandidate(out string, n int) (string, []int, []cdEdge, error) {
	reader := strings.NewReader(out)
	var status string
	if _, err := fmt.Fscan(reader, &status); err != nil {
		return "", nil, nil, fmt.Errorf("failed to read status: %v", err)
	}
	statusLower := strings.ToLower(status)
	if statusLower == "impossible" {
		return "impossible", nil, nil, nil
	}
	if statusLower != "possible" {
		return "", nil, nil, fmt.Errorf("unexpected status %q", status)
	}
	ct := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &ct[i]); err != nil {
			return "", nil, nil, fmt.Errorf("failed to read CT[%d]: %v", i+1, err)
		}
		if ct[i] != 0 && ct[i] != 1 {
			return "", nil, nil, fmt.Errorf("CT[%d] must be 0 or 1", i+1)
		}
	}
	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return "", nil, nil, fmt.Errorf("failed to read edges count: %v", err)
	}
	if m < 0 || m > maxEdges {
		return "", nil, nil, fmt.Errorf("edges count %d out of range [0,%d]", m, maxEdges)
	}
	edges := make([]cdEdge, m)
	for i := 0; i < m; i++ {
		if _, err := fmt.Fscan(reader, &edges[i].u, &edges[i].v); err != nil {
			return "", nil, nil, fmt.Errorf("failed to read edge %d: %v", i+1, err)
		}
	}
	return "possible", ct, edges, nil
}

func validateSolution(tc testInput, ct []int, edges []cdEdge) error {
	n := tc.n
	if len(ct) != n {
		return fmt.Errorf("expected %d CT values, got %d", n, len(ct))
	}
	if len(edges) > maxEdges {
		return fmt.Errorf("too many edges: %d > %d", len(edges), maxEdges)
	}
	for i, v := range ct {
		if v != 0 && v != 1 {
			return fmt.Errorf("CT[%d] must be 0 or 1", i+1)
		}
	}
	adj := make([][]int, n)
	for idx, e := range edges {
		if e.u < 1 || e.u > n || e.v < 1 || e.v > n {
			return fmt.Errorf("edge %d has invalid vertices (%d,%d)", idx+1, e.u, e.v)
		}
		if e.u == e.v {
			return fmt.Errorf("edge %d cannot be self-loop", idx+1)
		}
		adj[e.u-1] = append(adj[e.u-1], e.v-1)
	}
	for f := 0; f < 3; f++ {
		has := make([]bool, n)
		queue := []int{0}
		has[0] = true
		for len(queue) > 0 {
			s := queue[0]
			queue = queue[1:]
			for _, t := range adj[s] {
				if has[t] {
					continue
				}
				deliver := true
				if ct[s] == 1 && tc.pass[s][f] == 0 {
					deliver = false
				}
				if deliver {
					has[t] = true
					queue = append(queue, t)
				}
			}
		}
		for i := 0; i < n; i++ {
			want := tc.want[i][f]
			if want == 1 && !has[i] {
				return fmt.Errorf("feature %d missing on server %d", f+1, i+1)
			}
			if want == 0 && has[i] {
				return fmt.Errorf("feature %d wrongly deployed to server %d", f+1, i+1)
			}
		}
	}
	return nil
}

func deterministicTests() []testInput {
	tests := []testInput{
		{
			n: 2,
			want: [][]int{
				{1, 1, 1},
				{1, 1, 1},
			},
			pass: [][]int{
				{1, 1, 1},
				{1, 1, 1},
			},
		},
		{
			n: 3,
			want: [][]int{
				{1, 1, 1},
				{1, 0, 1},
				{0, 1, 0},
			},
			pass: [][]int{
				{1, 1, 1},
				{0, 1, 0},
				{1, 0, 1},
			},
		},
		{
			n: 4,
			want: [][]int{
				{1, 1, 1},
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
			pass: [][]int{
				{1, 1, 1},
				{1, 1, 0},
				{0, 1, 1},
				{1, 0, 1},
			},
		},
		{
			n: 5,
			want: [][]int{
				{1, 1, 1},
				{0, 1, 1},
				{1, 0, 1},
				{1, 1, 0},
				{0, 0, 1},
			},
			pass: [][]int{
				{1, 1, 1},
				{1, 0, 1},
				{0, 1, 1},
				{1, 1, 0},
				{0, 0, 1},
			},
		},
	}
	return tests
}

func randomTests(rng *rand.Rand, count int) []testInput {
	tests := make([]testInput, 0, count)
	for t := 0; t < count; t++ {
		n := rng.Intn(maxRandomN-2) + 2
		want := make([][]int, n)
		pass := make([][]int, n)
		for i := 0; i < n; i++ {
			want[i] = make([]int, 3)
			pass[i] = make([]int, 3)
		}
		for j := 0; j < 3; j++ {
			want[0][j] = 1
			pass[0][j] = 1
		}
		for i := 1; i < n; i++ {
			for j := 0; j < 3; j++ {
				want[i][j] = rng.Intn(2)
				pass[i][j] = rng.Intn(2)
			}
		}
		tests = append(tests, testInput{n: n, want: want, pass: pass})
	}
	return tests
}

func formatInput(tc testInput) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i := 0; i < tc.n; i++ {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.want[i][0], tc.want[i][1], tc.want[i][2])
	}
	for i := 0; i < tc.n; i++ {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.pass[i][0], tc.pass[i][1], tc.pass[i][2])
	}
	return sb.String()
}

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
	u, v int
	w    int64
}

type testCase struct {
	n     int
	edges []edge
}

const logN = 20

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1592D-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", path, "1592D.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return path, cleanup, nil
}

func runBinary(bin string, input string) (string, error) {
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

func parsePair(out string) (int, int, error) {
	fields := strings.Fields(out)
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("expected two integers, got %d tokens", len(fields))
	}
	a, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	b, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid integer %q", fields[1])
	}
	return a, b, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
	}
	return sb.String()
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

type lcaHelper struct {
	parent [][logN]int
	gcdUp  [][logN]int64
	depth  []int
}

func buildLCA(n int, edges []edge) *lcaHelper {
	adj := make([][]struct {
		to int
		w  int64
	}, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], struct {
			to int
			w  int64
		}{e.v, e.w})
		adj[e.v] = append(adj[e.v], struct {
			to int
			w  int64
		}{e.u, e.w})
	}
	parent := make([][logN]int, n+1)
	gcdUp := make([][logN]int64, n+1)
	depth := make([]int, n+1)
	var dfs func(int, int, int64)
	dfs = func(v, p int, w int64) {
		parent[v][0] = p
		gcdUp[v][0] = w
		for k := 1; k < logN; k++ {
			parent[v][k] = parent[parent[v][k-1]][k-1]
			gcdUp[v][k] = gcd(gcdUp[v][k-1], gcdUp[parent[v][k-1]][k-1])
		}
		for _, e := range adj[v] {
			if e.to == p {
				continue
			}
			depth[e.to] = depth[v] + 1
			dfs(e.to, v, e.w)
		}
	}
	dfs(1, 0, 0)
	return &lcaHelper{parent: parent, gcdUp: gcdUp, depth: depth}
}

func (h *lcaHelper) dist(u, v int) int64 {
	res := int64(0)
	if h.depth[u] < h.depth[v] {
		u, v = v, u
	}
	diff := h.depth[u] - h.depth[v]
	for k := logN - 1; k >= 0; k-- {
		if diff&(1<<k) != 0 {
			res = gcd(res, h.gcdUp[u][k])
			u = h.parent[u][k]
		}
	}
	if u == v {
		return res
	}
	for k := logN - 1; k >= 0; k-- {
		if h.parent[u][k] != h.parent[v][k] {
			res = gcd(res, h.gcdUp[u][k])
			res = gcd(res, h.gcdUp[v][k])
			u = h.parent[u][k]
			v = h.parent[v][k]
		}
	}
	res = gcd(res, h.gcdUp[u][0])
	res = gcd(res, h.gcdUp[v][0])
	return res
}

func bestDist(tc testCase, helper *lcaHelper) int64 {
	best := int64(0)
	for u := 1; u <= tc.n; u++ {
		for v := u + 1; v <= tc.n; v++ {
			val := helper.dist(u, v)
			if val > best {
				best = val
			}
		}
	}
	return best
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 2,
			edges: []edge{
				{u: 1, v: 2, w: 10},
			},
		},
		{
			n: 3,
			edges: []edge{
				{1, 2, 6},
				{2, 3, 9},
			},
		},
		{
			n: 4,
			edges: []edge{
				{1, 2, 7},
				{2, 3, 14},
				{2, 4, 21},
			},
		},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 2
	edges := make([]edge, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		w := int64(rng.Intn(1000) + 1)
		edges = append(edges, edge{u: u, v: v, w: w})
	}
	return testCase{n: n, edges: edges}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
		helper := buildLCA(tc.n, tc.edges)
		best := bestDist(tc, helper)
		if best == 0 {
			fmt.Fprintf(os.Stderr, "invalid test %d: best dist zero\n", idx+1)
			os.Exit(1)
		}

		// cross-check oracle
		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		expU, expV, err := parsePair(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}
		if expU < 1 || expU > tc.n || expV < 1 || expV > tc.n || expU == expV {
			fmt.Fprintf(os.Stderr, "oracle produced invalid nodes on test %d\noutput:\n%s\n", idx+1, expOut)
			os.Exit(1)
		}
		if helper.dist(expU, expV) != best {
			fmt.Fprintf(os.Stderr, "oracle output not optimal on test %d\n", idx+1)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotU, gotV, err := parsePair(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}
		if gotU < 1 || gotU > tc.n || gotV < 1 || gotV > tc.n || gotU == gotV {
			fmt.Fprintf(os.Stderr, "test %d: nodes out of range or equal (u=%d v=%d)\ninput:\n%s\n", idx+1, gotU, gotV, input)
			os.Exit(1)
		}
		val := helper.dist(gotU, gotV)
		if val != best {
			fmt.Fprintf(os.Stderr, "test %d: expected dist %d got %d\ninput:\n%s\n", idx+1, best, val, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

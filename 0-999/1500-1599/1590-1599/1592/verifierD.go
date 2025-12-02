package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
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

func runInteractive(bin string, tc testCase, helper *lcaHelper) (int, int, error) {
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return 0, 0, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, 0, err
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return 0, 0, err
	}

	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	// Write input (without weights)
	fmt.Fprintf(stdin, "%d\n", tc.n)
	for _, e := range tc.edges {
		fmt.Fprintf(stdin, "%d %d\n", e.u, e.v)
	}
	
	scanner := bufio.NewScanner(stdout)
	queries := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "!") {
			parts := strings.Fields(line)
			if len(parts) != 3 {
				return 0, 0, fmt.Errorf("invalid answer format: %s", line)
			}
			u, err := strconv.Atoi(parts[1])
			if err != nil { return 0, 0, fmt.Errorf("invalid node u: %s", parts[1]) }
			v, err := strconv.Atoi(parts[2])
			if err != nil { return 0, 0, fmt.Errorf("invalid node v: %s", parts[2]) }
			return u, v, nil
		} else if strings.HasPrefix(line, "?") {
			queries++
			if queries > 12 {
				return 0, 0, fmt.Errorf("too many queries")
			}
			parts := strings.Fields(line)
			if len(parts) < 2 {
				return 0, 0, fmt.Errorf("invalid query format: %s", line)
			}
			k, err := strconv.Atoi(parts[1])
			if err != nil { return 0, 0, fmt.Errorf("invalid count k: %s", parts[1]) }
			if len(parts) != k+2 {
				return 0, 0, fmt.Errorf("query token count mismatch: expected %d, got %d", k+2, len(parts))
			}
			nodes := make([]int, k)
			for i := 0; i < k; i++ {
				nodes[i], err = strconv.Atoi(parts[i+2])
				if err != nil { return 0, 0, fmt.Errorf("invalid node in query: %s", parts[i+2]) }
			}

			// Compute max dist
			ans := int64(0)
			for i := 0; i < k; i++ {
				for j := i + 1; j < k; j++ {
					d := helper.dist(nodes[i], nodes[j])
					if d > ans {
						ans = d
					}
				}
			}
			fmt.Fprintf(stdin, "%d\n", ans)
		} else {
			// Ignore empty lines?
			if line == "" { continue }
			return 0, 0, fmt.Errorf("unexpected output: %s", line)
		}
	}
	return 0, 0, fmt.Errorf("program exited without answer")
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

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		helper := buildLCA(tc.n, tc.edges)
		best := bestDist(tc, helper)
		if best == 0 {
			fmt.Fprintf(os.Stderr, "invalid test %d: best dist zero\n", idx+1)
			os.Exit(1)
		}

		gotU, gotV, err := runInteractive(target, tc, helper)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		
		if gotU < 1 || gotU > tc.n || gotV < 1 || gotV > tc.n || gotU == gotV {
			fmt.Fprintf(os.Stderr, "test %d: nodes out of range or equal (u=%d v=%d)\n", idx+1, gotU, gotV)
			os.Exit(1)
		}
		val := helper.dist(gotU, gotV)
		if val != best {
			fmt.Fprintf(os.Stderr, "test %d: expected dist %d got %d\n", idx+1, best, val)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

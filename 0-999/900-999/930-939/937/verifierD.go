package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ---------------- Helpers to compute expected output -----------------

type testCaseD struct {
	n     int
	edges [][]int
	start int
}

func generateD(rng *rand.Rand) testCaseD {
	n := rng.Intn(4) + 2 // 2..5
	edges := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		c := rng.Intn(n) // allow 0..n-1 edges
		edges[i] = make([]int, c)
		for j := 0; j < c; j++ {
			v := rng.Intn(n) + 1
			for v == i {
				v = rng.Intn(n) + 1
			}
			edges[i][j] = v
		}
	}
	start := rng.Intn(n) + 1
	return testCaseD{n: n, edges: edges, start: start}
}

func solveD(tc testCaseD) (string, []int) {
	n := tc.n
	adj := tc.edges
	vis := make([][2]bool, n+1)
	path := []int{tc.start}
	var ans []int
	var found bool
	var dfs func(u, mk int)
	dfs = func(u, mk int) {
		if found {
			return
		}
		for _, v := range adj[u] {
			if vis[v][mk^1] {
				continue
			}
			vis[v][mk^1] = true
			path = append(path, v)
			if mk == 0 && len(adj[v]) == 0 {
				ans = append([]int(nil), path...)
				found = true
				return
			}
			dfs(v, mk^1)
			path = path[:len(path)-1]
		}
	}
	dfs(tc.start, 0)
	if found {
		return "Win", ans
	}
	// tarjan to detect cycle
	dfn := make([]int, n+1)
	low := make([]int, n+1)
	inStack := make([]bool, n+1)
	stack := make([]int, 0)
	timer := 0
	maxComp := 0
	var tarjan func(u int)
	tarjan = func(u int) {
		timer++
		dfn[u] = timer
		low[u] = timer
		inStack[u] = true
		stack = append(stack, u)
		for _, v := range adj[u] {
			if dfn[v] == 0 {
				tarjan(v)
				if low[v] < low[u] {
					low[u] = low[v]
				}
			} else if inStack[v] {
				if dfn[v] < low[u] {
					low[u] = dfn[v]
				}
			}
		}
		if low[u] == dfn[u] {
			cnt := 0
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				inStack[w] = false
				cnt++
				if w == u {
					break
				}
			}
			if cnt > maxComp {
				maxComp = cnt
			}
		}
	}
	tarjan(tc.start)
	if maxComp > 1 {
		return "Draw", nil
	}
	return "Lose", nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func formatInput(tc testCaseD) string {
	var b strings.Builder
	m := 0
	for i := 1; i <= tc.n; i++ {
		m += len(tc.edges[i])
	}
	fmt.Fprintf(&b, "%d %d\n", tc.n, m)
	for i := 1; i <= tc.n; i++ {
		fmt.Fprintf(&b, "%d", len(tc.edges[i]))
		for _, v := range tc.edges[i] {
			fmt.Fprintf(&b, " %d", v)
		}
		b.WriteByte('\n')
	}
	fmt.Fprintf(&b, "%d\n", tc.start)
	return b.String()
}

func parsePath(line string) ([]int, error) {
	if strings.TrimSpace(line) == "" {
		return nil, fmt.Errorf("empty path line")
	}
	parts := strings.Fields(line)
	res := make([]int, len(parts))
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func runCase(bin string, tc testCaseD) error {
	input := formatInput(tc)
	expectedType, expectedPath := solveD(tc)
	gotStr, err := run(bin, input)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.TrimSpace(gotStr), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	outType := strings.TrimSpace(lines[0])
	if outType != expectedType {
		return fmt.Errorf("expected %s got %s", expectedType, outType)
	}
	if expectedType == "Win" {
		if len(lines) < 2 {
			return fmt.Errorf("missing path line")
		}
		path, err := parsePath(lines[1])
		if err != nil {
			return fmt.Errorf("bad path: %v", err)
		}
		if len(path) != len(expectedPath) {
			return fmt.Errorf("expected path %v got %v", expectedPath, path)
		}
		for i := range path {
			if path[i] != expectedPath[i] {
				return fmt.Errorf("expected path %v got %v", expectedPath, path)
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateD(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

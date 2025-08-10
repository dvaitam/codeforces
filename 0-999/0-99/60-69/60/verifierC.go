package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Edge struct {
	to   int
	g, l int64
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveC(r *bufio.Reader) string {
	var n, m int
	fmt.Fscan(r, &n, &m)
	adj := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var x, y int
		var g, l int64
		fmt.Fscan(r, &x, &y, &g, &l)
		adj[x] = append(adj[x], Edge{y, g, l})
		adj[y] = append(adj[y], Edge{x, g, l})
	}
	A := make([]int64, n+1)
	visited := make([]bool, n+1)
	var order []int
	var okFlag bool
	var dfs func(int)
	dfs = func(u int) {
		visited[u] = true
		order = append(order, u)
		for _, e := range adj[u] {
			v := e.to
			if visited[v] {
				if lcm(A[u], A[v]) != e.l || gcd(A[u], A[v]) != e.g {
					okFlag = false
				}
				continue
			}
			if e.l%A[u] != 0 {
				okFlag = false
				continue
			}
			A[v] = e.l / A[u] * e.g
			dfs(v)
			if !okFlag {
				return
			}
		}
	}
	for i := 1; i <= n; i++ {
		if visited[i] || len(adj[i]) == 0 {
			continue
		}
		var tmp, tmp2 int64
		for j, e := range adj[i] {
			if j == 0 {
				tmp = e.l
				tmp2 = e.g
			} else {
				tmp = gcd(tmp, e.l)
				tmp2 = lcm(tmp2, e.g)
			}
		}
		if tmp2 > tmp {
			return "NO"
		}
		var cands []int64
		for d := int64(1); d*d <= tmp; d++ {
			if tmp%d == 0 {
				if d%tmp2 == 0 {
					cands = append(cands, d)
				}
				oth := tmp / d
				if oth != d && oth%tmp2 == 0 {
					cands = append(cands, oth)
				}
			}
		}
		ok := false
		for _, D := range cands {
			A[i] = D
			order = order[:0]
			okFlag = true
			dfs(i)
			if okFlag {
				ok = true
				break
			}
			for _, u := range order {
				visited[u] = false
			}
		}
		if !ok {
			return "NO"
		}
	}
	for i := 1; i <= n; i++ {
		if A[i] == 0 {
			A[i] = 1
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&sb, "%d ", A[i])
	}
	return strings.TrimSpace(sb.String())
}

func generateCaseC(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	A := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		A[i] = int64(rng.Intn(20) + 1)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	edges := make([][2]int, 0, m)
	// create list of all pairs
	pairs := make([][2]int, 0, maxEdges)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			pairs = append(pairs, [2]int{i, j})
		}
	}
	rng.Shuffle(len(pairs), func(i, j int) { pairs[i], pairs[j] = pairs[j], pairs[i] })
	if m > len(pairs) {
		m = len(pairs)
	}
	pairs = pairs[:m]
	for _, p := range pairs {
		i, j := p[0], p[1]
		g := gcd(A[i], A[j])
		l := lcm(A[i], A[j])
		if rng.Float64() < 0.3 { // corrupt
			if rng.Intn(2) == 0 {
				g += int64(rng.Intn(3) + 1)
			} else {
				l += int64(rng.Intn(3) + 1)
			}
		}
		fmt.Fprintf(&sb, "%d %d %d %d\n", i, j, g, l)
		edges = append(edges, [2]int{i, j})
	}
	return sb.String()
}

// parse the test case input and validate solver output
type checkEdge struct {
	u, v int
	g, l *big.Int
}

func parseInput(tc string) (int, []checkEdge, error) {
	r := bufio.NewReader(strings.NewReader(tc))
	var n, m int
	if _, err := fmt.Fscan(r, &n, &m); err != nil {
		return 0, nil, err
	}
	edges := make([]checkEdge, m)
	for i := 0; i < m; i++ {
		var x, y int
		var g, l int64
		if _, err := fmt.Fscan(r, &x, &y, &g, &l); err != nil {
			return 0, nil, err
		}
		edges[i] = checkEdge{u: x, v: y, g: big.NewInt(g), l: big.NewInt(l)}
	}
	return n, edges, nil
}

func checkSolution(tc, out string) error {
	n, edges, err := parseInput(tc)
	if err != nil {
		return fmt.Errorf("invalid test case: %v", err)
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	if fields[0] != "YES" {
		return fmt.Errorf("expected YES, got %s", fields[0])
	}
	if len(fields) != n+1 {
		return fmt.Errorf("expected %d numbers, got %d", n, len(fields)-1)
	}
	vals := make([]*big.Int, n+1)
	for i := 1; i <= n; i++ {
		v, ok := new(big.Int).SetString(fields[i], 10)
		if !ok {
			return fmt.Errorf("invalid integer value")
		}
		if v.Sign() <= 0 {
			return fmt.Errorf("non-positive value")
		}
		vals[i] = v
	}
	for _, e := range edges {
		g := new(big.Int).GCD(nil, nil, vals[e.u], vals[e.v])
		if g.Cmp(e.g) != 0 {
			return fmt.Errorf("gcd mismatch on edge %d-%d", e.u, e.v)
		}
		l := new(big.Int).Mul(vals[e.u], vals[e.v])
		l.Div(l, g)
		if l.Cmp(e.l) != 0 {
			return fmt.Errorf("lcm mismatch on edge %d-%d", e.u, e.v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseC(rng)
		expect := solveC(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		expFields := strings.Fields(expect)
		if len(expFields) == 0 {
			fmt.Fprintf(os.Stderr, "internal error on case %d\n", i+1)
			os.Exit(1)
		}
		if expFields[0] == "NO" {
			if strings.TrimSpace(got) != "NO" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected NO got %s\ninput:\n%s", i+1, got, tc)
				os.Exit(1)
			}
			continue
		}
		if err := checkSolution(tc, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, tc, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

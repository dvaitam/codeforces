package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const solution1336FSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxN = 150005
const logN = 18

var (
	adj        [maxN][]edge
	up         [logN][maxN]int
	parentEdge [maxN]int
	depth      [maxN]int
)

type edge struct{ to, id int }

func dfs(v, p int) {
	for _, e := range adj[v] {
		if e.to == p {
			continue
		}
		up[0][e.to] = v
		parentEdge[e.to] = e.id
		depth[e.to] = depth[v] + 1
		for i := 1; i < logN; i++ {
			up[i][e.to] = up[i-1][up[i-1][e.to]]
		}
		dfs(e.to, v)
	}
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for i := 0; i < logN; i++ {
		if diff&(1<<i) > 0 {
			a = up[i][a]
		}
	}
	if a == b {
		return a
	}
	for i := logN - 1; i >= 0; i-- {
		if up[i][a] != up[i][b] {
			a = up[i][a]
			b = up[i][b]
		}
	}
	return up[0][a]
}

func getPathEdges(u, v int) []int {
	p := lca(u, v)
	res := make([]int, 0)
	x := u
	for x != p {
		res = append(res, parentEdge[x])
		x = up[0][x]
	}
	tmp := make([]int, 0)
	x = v
	for x != p {
		tmp = append(tmp, parentEdge[x])
		x = up[0][x]
	}
	for i := len(tmp) - 1; i >= 0; i-- {
		res = append(res, tmp[i])
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	for i := 1; i <= n; i++ {
		adj[i] = nil
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], edge{v, i})
		adj[v] = append(adj[v], edge{u, i})
	}
	depth[1] = 0
	dfs(1, 0)
	edgeTrav := make([][]int, n-1)
	for i := 0; i < m; i++ {
		var s, t int
		fmt.Fscan(reader, &s, &t)
		path := getPathEdges(s, t)
		for _, e := range path {
			edgeTrav[e] = append(edgeTrav[e], i)
		}
	}
	pairCount := make(map[uint64]int)
	ans := 0
	for _, list := range edgeTrav {
		L := len(list)
		for i := 0; i < L; i++ {
			for j := i + 1; j < L; j++ {
				a := list[i]
				b := list[j]
				if a > b {
					a, b = b, a
				}
				key := uint64(a)<<32 | uint64(b)
				pairCount[key]++
				if pairCount[key] == k {
					ans++
				}
			}
		}
	}
	fmt.Println(ans)
}
`

// Keep the embedded reference solution reachable.
var _ = solution1336FSource

type testCase struct {
	n     int
	m     int
	k     int
	edges [][2]int
	pairs [][2]int
}

const testcasesRaw = `5 3 5 1 2 2 3 2 4 3 5 5 3 4 5 3 4
4 4 3 1 2 1 3 3 4 4 3 4 2 2 4 2 1
4 5 2 1 2 2 3 1 4 4 2 4 2 1 2 2 3 3 1
2 4 2 1 2 2 1 1 2 2 1 1 2
2 6 1 1 2 1 2 2 1 2 1 1 2 1 2 2 1
4 4 1 1 2 1 3 2 4 4 3 2 4 1 3 3 2
2 5 2 1 2 1 2 2 1 2 1 2 1 2 1
3 3 2 1 2 1 3 2 1 3 2 1 2
2 3 1 1 2 1 2 2 1 2 1
6 3 5 1 2 1 3 1 4 2 5 1 6 2 4 6 2 1 4
3 4 2 1 2 2 3 1 3 2 1 2
3 2 1 1 2 2 1 1 2
3 6 3 1 2 2 3 2 4 3 5 1 6 1 2 6 3 5 3 6
5 3 3 1 2 2 3 2 4 3 5 4 1 4 3 4
2 4 3 1 2 2 1 1 2 2 1 2 1
2 3 2 1 2 2 1 1 2 1 1
2 3 1 1 2 2 1 1 2
2 4 1 1 2 2 1 1 2 2 2
2 3 3 1 2 2 1 1 2 1 1 2
3 6 1 1 2 1 3 1 4 1 5 2 6 1 4 5 6 2 2
3 3 3 1 2 2 1 1 2 1 2 3 2 1
2 5 1 1 2 1 1 2 1 1 2 1 2
3 4 3 1 2 2 1 1 3 1 3 1 2 1 3
2 5 2 1 2 2 1 1 2 1 1 2 1 2
2 4 3 1 2 2 1 1 2 2 1 2 2 1
3 3 1 1 2 1 3 2 1 3 1 2
4 4 5 1 2 1 3 1 4 2 3 2 3 4 1 2 2 3 1 2 1 4 2
4 4 3 1 2 1 3 1 4 1 3 1 3 2 3 1 4
4 3 4 1 2 2 3 3 4 1 4 2 4 3 1 3 3
3 5 5 1 2 2 3 2 1 3 1 3 1 2 2 3 2 1 3
5 5 1 1 2 1 3 1 4 1 5 2 2 3 4 5 4 5 3 4 2 5 1
4 3 3 1 2 1 3 2 4 2 2 2 4 2 3 1 3
4 5 3 1 2 2 3 3 4 2 4 3 4 1 4 3 2 2 1
4 3 3 1 2 1 3 2 4 1 2 4 3 4 3 2 3
2 4 1 1 2 1 2 1 2 2 1 1
3 5 3 1 2 2 3 1 4 2 2 3 3 4 4 1 2
2 5 1 1 2 1 2 1 2 2 1 1 2
5 3 1 1 2 1 3 1 4 3 5 4 2 4 1 2 1 3
5 4 4 1 2 1 3 1 4 1 5 1 2 3 4 1 2 5 3 4 2 3
3 4 4 1 2 1 3 2 3 3 1 2 2 1 2 2 1
5 4 4 1 2 1 3 1 4 2 5 1 2 3 3 4 3 4 4 4 4
4 3 3 1 2 1 3 2 3 3 3 1 2 2 1 2 3
2 3 2 1 2 2 1 1 2 1 1
2 6 2 1 2 1 2 1 2 2 3 1 1 1 1 2 2
4 5 3 1 2 1 3 2 3 3 2 4 1 3 1 3 3 3 4
2 3 1 1 2 1 2 2 1 2
4 3 1 1 2 2 3 1 3 3 3 3 4 2 4
2 4 3 1 2 1 2 2 1 1 2 2 1
5 4 2 1 2 2 3 2 4 3 5 2 1 2 2 3 4 4 5
5 6 5 1 2 2 3 1 4 2 5 3 2 4 4 5 1 5 3 1 5 3 1 2 3 3 4 2 4 5 3
3 3 1 1 2 2 1 2 3 3 1
4 4 3 1 2 1 3 3 4 1 3 2 3 2 4 3 2
4 3 2 1 2 2 3 3 2 3 2 1 3 1
5 6 2 1 2 2 3 3 4 3 5 4 1 3 3 5 1 2 3 4 4 5 4 3 5
2 4 2 1 2 2 1 1 2 1 1
2 3 2 1 2 2 1 1 2 1 1
6 5 5 1 2 1 3 2 4 3 5 4 6 5 3 2 2 5 2 4 3 4 4 1 1 2 3
5 4 3 1 2 1 3 1 4 2 5 1 3 2 1 4 3 3 4 4 2
5 6 1 1 2 1 3 1 4 1 5 1 2 3 2 4 1 5 1 1 2 1 3 2 4 3 5
5 6 1 1 2 2 3 2 4 2 5 2 2 3 2 4 4 5 4 5 5 5 2 5 3 4 4
5 6 3 1 2 2 3 1 4 3 5 2 4 1 3 5 3 3 4 5 2 4 3 5 3 4 4
5 5 2 1 2 1 3 1 4 2 5 3 3 4 2 3 2 4 3 5 1 2 2
4 4 5 1 2 1 3 2 4 1 4 2 4 3 1 3 4 1 1 2 3 4 2
4 6 4 1 2 1 3 2 4 2 3 3 2 5 4 3 5 1 2 4 2 3 4 6 4
4 6 2 1 2 1 3 1 4 3 2 3 3 5 1 6 4 3 4 5 3 3 4
5 3 3 1 2 2 3 1 4 1 3 2 1 3 2
3 3 3 1 2 1 2 1 3 2 1 3
5 4 3 1 2 1 3 2 4 1 5 1 2 1 3 4 3 5 4 2
2 3 2 1 2 2 1 1 2 2 2
4 3 3 1 2 2 3 1 3 1 3 3 2 3 1
3 4 4 1 2 1 3 1 4 3 4 2 2
5 6 5 1 2 1 3 1 4 2 5 3 3 4 3 5 3 5 3 5 1 2 3 4 3 5 1
3 6 4 1 2 1 3 2 3 3 2 4 1 5 4 3 5 2 2
5 6 2 1 2 2 3 1 4 3 5 1 2 2 3 2 4 2 5 2 5 3 4 3 5
5 6 5 1 2 2 3 1 4 1 5 2 3 3 2 3 3 4 3 4 4 4 2 3 5 3 5 3
4 4 3 1 2 1 3 1 4 1 2 2 3 3 4 4 3
5 3 5 1 2 1 3 2 4 3 5 5 3 4 5 3 4
3 3 4 1 2 2 1 1 3 1 3 2
4 6 3 1 2 2 3 3 4 1 4 2 4 3 5 4 6 2 2 3 4
4 6 2 1 2 1 3 2 4 3 2 4 3 5 1 4 1 2 4 6 1
`

type edgeLocal struct{ to, id int }

func parseTestcases() []testCase {
	reader := bufio.NewReader(strings.NewReader(testcasesRaw))
	var res []testCase
	for {
		var n, m, k int
		if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
			break
		}
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			if _, err := fmt.Fscan(reader, &edges[i][0], &edges[i][1]); err != nil {
				return res
			}
		}
		pairs := make([][2]int, m)
		for i := 0; i < m; i++ {
			if _, err := fmt.Fscan(reader, &pairs[i][0], &pairs[i][1]); err != nil {
				return res
			}
		}
		res = append(res, testCase{n: n, m: m, k: k, edges: edges, pairs: pairs})
	}
	return res
}

func expected(tc testCase) int {
	n := tc.n
	k := tc.k
	adj := make([][]edgeLocal, n+1)
	for i, e := range tc.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], edgeLocal{v, i})
		adj[v] = append(adj[v], edgeLocal{u, i})
	}
	const logN = 18
	up := make([][]int, logN)
	for i := range up {
		up[i] = make([]int, n+1)
	}
	parentEdge := make([]int, n+1)
	depth := make([]int, n+1)
	var dfs func(int, int)
	dfs = func(v, p int) {
		for _, e := range adj[v] {
			if e.to == p {
				continue
			}
			up[0][e.to] = v
			parentEdge[e.to] = e.id
			depth[e.to] = depth[v] + 1
			for i := 1; i < logN; i++ {
				up[i][e.to] = up[i-1][up[i-1][e.to]]
			}
			dfs(e.to, v)
		}
	}
	if n > 0 {
		dfs(1, 0)
	}
	lca := func(a, b int) int {
		if a <= 0 || b <= 0 || a > n || b > n {
			return 0
		}
		if depth[a] < depth[b] {
			a, b = b, a
		}
		diff := depth[a] - depth[b]
		for i := 0; i < logN; i++ {
			if diff&(1<<i) > 0 {
				a = up[i][a]
			}
		}
		if a == b {
			return a
		}
		for i := logN - 1; i >= 0; i-- {
			if up[i][a] != up[i][b] {
				a = up[i][a]
				b = up[i][b]
			}
		}
		return up[0][a]
	}
	getPath := func(u, v int) []int {
		if u <= 0 || v <= 0 || u > n || v > n {
			return nil
		}
		p := lca(u, v)
		if p == 0 {
			return nil
		}
		res := make([]int, 0)
		x := u
		for x != p {
			res = append(res, parentEdge[x])
			x = up[0][x]
		}
		tmp := make([]int, 0)
		x = v
		for x != p {
			tmp = append(tmp, parentEdge[x])
			x = up[0][x]
		}
		for i := len(tmp) - 1; i >= 0; i-- {
			res = append(res, tmp[i])
		}
		return res
	}
	edgeTrav := make([][]int, n-1)
	for idx, pr := range tc.pairs {
		path := getPath(pr[0], pr[1])
		for _, e := range path {
			if e >= 0 && e < len(edgeTrav) {
				edgeTrav[e] = append(edgeTrav[e], idx)
			}
		}
	}
	pairCount := make(map[uint64]int)
	ans := 0
	for _, list := range edgeTrav {
		L := len(list)
		for i := 0; i < L; i++ {
			for j := i + 1; j < L; j++ {
				a := list[i]
				b := list[j]
				if a > b {
					a, b = b, a
				}
				key := uint64(a)<<32 | uint64(b)
				pairCount[key]++
				if pairCount[key] == k {
					ans++
				}
			}
		}
	}
	return ans
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.k)
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	for _, p := range tc.pairs {
		fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
	}
	return sb.String()
}

func runCase(bin string, idx int, tc testCase) error {
	expect := expected(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(buildInput(tc))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("case %d failed: %v\nstderr: %s", idx, err, string(out))
	}
	gotStr := strings.TrimSpace(string(out))
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("case %d failed: invalid output %q", idx, gotStr)
	}
	if got != expect {
		return fmt.Errorf("case %d failed: expected %d got %d", idx, expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := parseTestcases()
	for i, tc := range testcases {
		if err := runCase(bin, i+1, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

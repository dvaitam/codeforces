package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Edge struct{ u, v int }

type Test struct {
	n     int
	edges []Edge
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.n))
	for _, e := range t.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	return sb.String()
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// solve is the CF-accepted algorithm for 955F embedded directly.
func solve(input string) string {
	buf := []byte(input)
	pos := 0
	nextInt := func() int {
		for pos < len(buf) && buf[pos] <= ' ' {
			pos++
		}
		if pos >= len(buf) {
			return 0
		}
		res := 0
		for pos < len(buf) && buf[pos] > ' ' {
			res = res*10 + int(buf[pos]-'0')
			pos++
		}
		return res
	}

	n := nextInt()
	if n == 0 {
		return "0"
	}

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		u := nextInt()
		v := nextInt()
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	depth := make([]int, n+1)
	dfs_in := make([]int, n+1)
	timer_in := 0

	euler := make([]int, 0, 2*n)
	first := make([]int, n+1)

	children := make([][]int, n+1)
	C := make([]int, n+1)

	var dfs func(u, p, d int)
	dfs = func(u, p, d int) {
		depth[u] = d
		timer_in++
		dfs_in[u] = timer_in

		first[u] = len(euler)
		euler = append(euler, u)

		for _, v := range adj[u] {
			if v != p {
				children[u] = append(children[u], v)
				C[u]++
				dfs(v, u, d+1)
				euler = append(euler, u)
			}
		}
	}
	dfs(1, 0, 1)

	m := len(euler)
	lg := make([]int, m+1)
	for i := 2; i <= m; i++ {
		lg[i] = lg[i/2] + 1
	}
	st := make([][]int, lg[m]+1)
	for i := range st {
		st[i] = make([]int, m)
	}
	for i := 0; i < m; i++ {
		st[0][i] = euler[i]
	}
	for i := 1; i <= lg[m]; i++ {
		for j := 0; j+(1<<i) <= m; j++ {
			u := st[i-1][j]
			v := st[i-1][j+(1<<(i-1))]
			if depth[u] < depth[v] {
				st[i][j] = u
			} else {
				st[i][j] = v
			}
		}
	}

	getLCA := func(u, v int) int {
		l := first[u]
		r := first[v]
		if l > r {
			l, r = r, l
		}
		k := lg[r-l+1]
		n1 := st[k][l]
		n2 := st[k][r-(1<<k)+1]
		if depth[n1] < depth[n2] {
			return n1
		}
		return n2
	}

	for u := 1; u <= n; u++ {
		sort.Slice(children[u], func(i, j int) bool {
			return C[children[u][i]] > C[children[u][j]]
		})
	}

	nodes := make([]int, n)
	for i := 0; i < n; i++ {
		nodes[i] = i + 1
	}
	sort.Slice(nodes, func(i, j int) bool {
		return depth[nodes[i]] > depth[nodes[j]]
	})

	V := make([][]int, n+1)
	for _, u := range nodes {
		for k := 2; k <= C[u]; k++ {
			V[k] = append(V[k], u)
		}
	}

	dp1 := make([]int, n+1)
	var dfs1 func(u int)
	dfs1 = func(u int) {
		dp1[u] = 1
		for _, v := range children[u] {
			dfs1(v)
			if dp1[v]+1 > dp1[u] {
				dp1[u] = dp1[v] + 1
			}
		}
	}
	dfs1(1)

	ans := int64(0)
	for u := 1; u <= n; u++ {
		ans += int64(dp1[u])
	}

	f := make([]int, n+1)
	var S [][]int
	vals := make([]int, 0, 32)

	for k := 2; k <= n; k++ {
		if len(V[k]) == 0 {
			ans += int64(n)
			continue
		}
		ans += int64(n)

		max_f := 0
		for _, u := range V[k] {
			vals = vals[:0]
			for _, v := range children[u] {
				if C[v] < k {
					break
				}
				vals = append(vals, f[v])
			}
			L := len(vals)
			if k <= L {
				sort.Slice(vals, func(i, j int) bool {
					return vals[i] > vals[j]
				})
				f[u] = vals[k-1] + 1
			} else {
				f[u] = 2
			}
			if f[u] > max_f {
				max_f = f[u]
			}
		}

		if cap(S) <= max_f {
			newS := make([][]int, max_f+1)
			copy(newS, S)
			S = newS
		} else {
			S = S[:max_f+1]
		}
		for mm := 2; mm <= max_f; mm++ {
			S[mm] = S[mm][:0]
		}

		for _, u := range V[k] {
			for mm := 2; mm <= f[u]; mm++ {
				S[mm] = append(S[mm], u)
			}
		}

		for mm := 2; mm <= max_f; mm++ {
			if len(S[mm]) == 0 {
				continue
			}
			sort.Slice(S[mm], func(i, j int) bool {
				return dfs_in[S[mm][i]] < dfs_in[S[mm][j]]
			})

			union_size := 0
			for i := 0; i < len(S[mm]); i++ {
				u := S[mm][i]
				union_size += depth[u]
				if i > 0 {
					lca := getLCA(S[mm][i-1], u)
					union_size -= depth[lca]
				}
			}
			ans += int64(union_size)
		}
	}

	return fmt.Sprintf("%d", ans)
}

func genTree(rng *rand.Rand, n int) []Edge {
	edges := make([]Edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, Edge{p, i})
	}
	return edges
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(5))
	tests := make([]Test, 0, 105)
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		edges := genTree(rng, n)
		tests = append(tests, Test{n: n, edges: edges})
	}
	tests = append(tests,
		Test{n: 1, edges: nil},
		Test{n: 2, edges: []Edge{{1, 2}}},
		Test{n: 3, edges: []Edge{{1, 2}, {1, 3}}},
		Test{n: 4, edges: []Edge{{1, 2}, {2, 3}, {3, 4}}},
		Test{n: 5, edges: []Edge{{1, 2}, {1, 3}, {3, 4}, {4, 5}}},
	)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		want := strings.TrimSpace(solve(input))
		got, err := runExe(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n" + input)
			os.Exit(1)
		}
		if want != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:%sexpected:%s\ngot:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const maxBit = 30

type Basis struct {
	b [maxBit + 1]int
}

func (bs *Basis) Add(x int) {
	for i := maxBit; i >= 0; i-- {
		if (x>>i)&1 == 0 {
			continue
		}
		if bs.b[i] == 0 {
			bs.b[i] = x
			return
		}
		if x^bs.b[i] < x {
			x ^= bs.b[i]
		} else {
			x ^= bs.b[i]
		}
	}
}

func (bs *Basis) Merge(o *Basis) {
	for i := maxBit; i >= 0; i-- {
		if o.b[i] != 0 {
			bs.Add(o.b[i])
		}
	}
}

func (bs *Basis) MaxXor() int {
	res := 0
	for i := maxBit; i >= 0; i-- {
		if bs.b[i] != 0 && (res^bs.b[i]) > res {
			res ^= bs.b[i]
		}
	}
	return res
}

func copyBasis(src *Basis) Basis {
	var d Basis
	d = *src
	return d
}

// ----- parsing -----

type testCaseE struct {
	n       int
	vals    []int
	edges   [][2]int
	queries [][2]int
}

func parseCasesE(path string) ([]testCaseE, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, err
	}
	cases := make([]testCaseE, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return nil, err
		}
		vals := make([]int, n+1)
		for j := 1; j <= n; j++ {
			if _, err := fmt.Fscan(in, &vals[j]); err != nil {
				return nil, err
			}
		}
		edges := make([][2]int, n-1)
		for j := 0; j < n-1; j++ {
			var u, v int
			if _, err := fmt.Fscan(in, &u, &v); err != nil {
				return nil, err
			}
			edges[j] = [2]int{u, v}
		}
		var q int
		if _, err := fmt.Fscan(in, &q); err != nil {
			return nil, err
		}
		queries := make([][2]int, q)
		for j := 0; j < q; j++ {
			var r, v int
			if _, err := fmt.Fscan(in, &r, &v); err != nil {
				return nil, err
			}
			queries[j] = [2]int{r, v}
		}
		cases[i] = testCaseE{n: n, vals: vals, edges: edges, queries: queries}
	}
	return cases, nil
}

// solver derived from 1778E.go
func solveE(tc testCaseE) []int {
	n := tc.n
	a := tc.vals
	g := make([][]int, n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	tin := make([]int, n+1)
	tout := make([]int, n+1)
	order := make([]int, 0, n)
	children := make([][]int, n+1)

	type item struct{ v, p, idx int }
	stack := []item{{1, 0, 0}}
	time := 0
	for len(stack) > 0 {
		cur := &stack[len(stack)-1]
		v := cur.v
		if cur.idx == 0 {
			tin[v] = time
			time++
			parent[v] = cur.p
			if cur.p != 0 {
				depth[v] = depth[cur.p] + 1
				children[cur.p] = append(children[cur.p], v)
			}
		}
		if cur.idx < len(g[v]) {
			to := g[v][cur.idx]
			cur.idx++
			if to == cur.p {
				continue
			}
			stack = append(stack, item{to, v, 0})
		} else {
			tout[v] = time - 1
			order = append(order, v)
			stack = stack[:len(stack)-1]
		}
	}

	down := make([]Basis, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		down[v].Add(a[v])
		for _, to := range children[v] {
			down[v].Merge(&down[to])
		}
	}

	up := make([]Basis, n+1)
	type qitem struct{ v int }
	q := []qitem{{1}}
	for len(q) > 0 {
		cur := q[len(q)-1]
		q = q[:len(q)-1]
		v := cur.v
		m := len(children[v])
		pref := make([]Basis, m+1)
		suff := make([]Basis, m+1)
		for i := 0; i < m; i++ {
			pref[i+1] = copyBasis(&pref[i])
			pref[i+1].Merge(&down[children[v][i]])
		}
		for i := m - 1; i >= 0; i-- {
			suff[i] = copyBasis(&suff[i+1])
			suff[i].Merge(&down[children[v][i]])
		}
		for i, to := range children[v] {
			tmp := copyBasis(&up[v])
			tmp.Add(a[v])
			bs := copyBasis(&pref[i])
			bs.Merge(&suff[i+1])
			tmp.Merge(&bs)
			up[to] = tmp
			q = append(q, qitem{to})
		}
	}

	LOG := 20
	upTable := make([][]int, LOG)
	for i := range upTable {
		upTable[i] = make([]int, n+1)
	}
	for i := 1; i <= n; i++ {
		upTable[0][i] = parent[i]
	}
	for k := 1; k < LOG; k++ {
		for i := 1; i <= n; i++ {
			upTable[k][i] = upTable[k-1][upTable[k-1][i]]
		}
	}
	isAncestor := func(u, v int) bool { return tin[u] <= tin[v] && tout[v] <= tout[u] }
	getKth := func(v, k int) int {
		for i := 0; i < LOG; i++ {
			if (k>>i)&1 == 1 {
				v = upTable[i][v]
			}
		}
		return v
	}

	totalBasis := down[1]

	res := make([]int, len(tc.queries))
	for idx, qv := range tc.queries {
		r := qv[0]
		v := qv[1]
		if r == v {
			res[idx] = totalBasis.MaxXor()
		} else if !isAncestor(v, r) {
			res[idx] = down[v].MaxXor()
		} else {
			child := getKth(r, depth[r]-depth[v]-1)
			res[idx] = up[child].MaxXor()
		}
	}
	return res
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCasesE("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	idxCase := 0
	for _, tc := range cases {
		expected := solveE(tc)
		for qi, val := range tc.queries {
			sb := strings.Builder{}
			fmt.Fprintf(&sb, "1\n%d\n", tc.n)
			for i := 1; i <= tc.n; i++ {
				if i > 1 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(tc.vals[i]))
			}
			sb.WriteByte('\n')
			for _, e := range tc.edges {
				fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
			}
			fmt.Fprintf(&sb, "1\n%d %d\n", val[0], val[1])
			// Wait - verifying expects q queries at once? Actually original solution expects q after edges.
			// We'll build input for entire case with q=1 each time to test separately.
			input := sb.String()
			gotStr, err := run(bin, input)
			idxCase++
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idxCase, err)
				os.Exit(1)
			}
			got, err := strconv.Atoi(strings.TrimSpace(gotStr))
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", idxCase, err)
				os.Exit(1)
			}
			if got != expected[qi] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idxCase, expected[qi], got)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", idxCase)
}

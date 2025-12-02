package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const maxBit = 30
const MOD int64 = 998244353

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
		x ^= bs.b[i]
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

type testCase struct {
	n       int
	vals    []int
	edges   [][2]int
	queries [][2]int
}

const testcaseData = `100
5
13 1 8 15 12
1 2
2 3
2 4
2 5
3
2 3
2 1
5 3
6
4 9 3 2 10 15
1 2
2 3
2 4
3 5
5 6
3
2 5
4 4
5 3
2
0 2
1 2
3
1 2
2 1
2 1
3
7 7 4
1 2
1 3
1
2 3
5
3 9 9 3 10
1 2
2 3
2 4
1 5
3
4 3
5 2
3 2
3
5 1 8
1 2
1 3
1
3 1
3
1 2 12
1 2
1 3
1
3 3
5
8 14 15 11 2
1 2
1 3
2 4
3 5
1
2 1
4
3 7 11 5
1 2
2 3
1 4
1
2 2
2
2 0
1 2
3
1 1
2 1
2 1
2
0 6
1 2
3
1 2
1 1
1 2
6
3 8 2 7 2 9
1 2
2 3
1 4
1 5
5 6
2
1 5
1 6
5
6 8 11 15 5
1 2
1 3
3 4
2 5
1
3 5
4
3 14 5 0
1 2
2 3
3 4
3
3 3
4 3
2 1
5
2 10 1 8 4
1 2
2 3
2 4
3 5
3
3 5
5 2
3 4
5
2 0 6 10 5
1 2
1 3
3 4
4 5
2
5 4
1 4
6
13 1 5 14 2 8
1 2
2 3
3 4
4 5
5 6
3
1 1
4 3
3 4
2
13 6
1 2
3
1 1
2 2
2 1
3
0 0 3
1 2
1 3
3
3 1
2 2
3 1
2
15 12
1 2
1
2 2
2
8 4
1 2
1
1 2
2
1 1
1 2
3
2 2
2 1
2 2
5
11 5 6 12 9
1 2
1 3
1 4
3 5
2
3 3
1 3
6
1 1 8 5 4 9
1 2
2 3
3 4
2 5
3 6
1
4 6
3
1 9 5
1 2
2 3
2
2 2
2 1
2
15 15
1 2
2
1 2
1 2
5
1 9 10 4 5
1 2
1 3
1 4
1 5
1
2 1
5
0 3 12 9 14
1 2
1 3
2 4
1 5
2
2 3
5 2
5
6 11 3 2 0
1 2
1 3
1 4
4 5
2
3 2
1 2
6
4 3 6 14 12 11
1 2
1 3
3 4
4 5
2 6
3
4 6
6 4
5 4
4
15 15 6 7
1 2
2 3
3 4
3
3 3
1 2
3 2
5
9 15 2 2 1
1 2
1 3
1 4
1 5
2
1 4
3 2
3
14 11 12
1 2
1 3
3
3 3
1 3
2 1
4
13 15 12 7
1 2
1 3
3 4
1
3 3
4
2 15 8 9
1 2
2 3
2 4
1
2 2
3
9 10 1
1 2
2 3
2
1 2
3 3
2
4 11
1 2
1
2 2
5
1 3 15 4 0
1 2
1 3
3 4
3 5
1
5 3
3
12 15 3
1 2
2 3
3
3 2
3 1
3 3
6
9 4 12 9 3 6
1 2
2 3
2 4
3 5
2 6
2
3 6
1 1
2
15 8
1 2
3
1 1
1 2
2 1
3
13 13 2
1 2
2 3
1
1 2
3
0 14 13
1 2
1 3
2
2 3
2 1
4
2 3 11 0
1 2
2 3
1 4
1
2 3
2
4 6
1 2
1
1 1
4
11 0 7 4
1 2
2 3
1 4
2
3 3
2 1
3
11 10 15
1 2
2 3
3
3 2
1 3
1 1
6
9 5 12 4 4 7
1 2
1 3
1 4
2 5
3 6
2
4 6
1 2
6
0 12 2 2 4 13
1 2
2 3
3 4
2 5
5 6
2
3 6
3 1
3
14 11 1
1 2
2 3
1
2 3
4
14 6 11 9
1 2
1 3
1 4
1
3 1
6
4 14 12 5 13 13
1 2
1 3
2 4
3 5
5 6
1
3 4
2
15 6
1 2
1
2 2
2
6 9
1 2
3
2 1
2 2
1 2
3
12 8 0
1 2
2 3
3
1 1
2 2
3 3
5
14 3 8 11 9
1 2
1 3
1 4
1 5
2
3 5
3 1
6
7 5 2 13 9 9
1 2
1 3
3 4
1 5
4 6
3
5 4
6 3
3 4
4
4 5 3 3
1 2
2 3
3 4
2
2 3
3 4
5
6 15 15 10 15
1 2
2 3
2 4
2 5
3
4 1
5 2
1 3
5
12 0 2 2 12
1 2
2 3
1 4
1 5
3
1 3
3 2
2 5
4
6 3 13 14
1 2
2 3
1 4
2
4 4
2 4
3
10 4 6
1 2
2 3
2
2 2
2 2
3
6 14 6
1 2
2 3
1
1 3
2
5 11
1 2
3
1 1
2 1
2 2
5
14 1 13 14 15
1 2
2 3
1 4
3 5
2
1 1
1 2
4
0 9 0 4
1 2
2 3
3 4
1
4 2
5
6 10 3 2 10
1 2
2 3
2 4
3 5
1
5 1
3
11 2 6
1 2
1 3
1
2 3
4
9 12 8 15
1 2
1 3
1 4
2
1 1
4 4
5
1 13 15 14 14
1 2
1 3
1 4
2 5
1
2 4
3
14 2 13
1 2
1 3
1
1 2
3
4 8 11
1 2
2 3
1
3 2
6
6 9 14 14 8 8
1 2
1 3
1 4
1 5
2 6
3
4 2
2 3
6 6
2
13 1
1 2
2
2 1
2 1
6
9 7 7 2 9 10
1 2
2 3
3 4
4 5
3 6
3
2 2
1 5
5 3
4
0 4 12 4
1 2
1 3
1 4
1
4 2
3
4 7 12
1 2
1 3
3
2 1
3 1
3 3
4
15 14 9 0
1 2
1 3
3 4
2
4 3
1 3
3
12 6 10
1 2
2 3
1
1 1
4
7 10 14 7
1 2
2 3
3 4
1
3 1
4
1 4 11 0
1 2
1 3
1 4
1
1 1
3
10 2 1
1 2
2 3
1
1 2
5
4 11 9 5 10
1 2
2 3
1 4
4 5
2
5 5
4 1
6
3 13 12 5 0 4
1 2
1 3
2 4
2 5
2 6
1
1 2
6
5 2 13 3 14 4
1 2
2 3
2 4
4 5
1 6
3
1 4
1 3
3 6
3
14 7 11
1 2
2 3
2
2 2
2 1
4
9 15 1 8
1 2
2 3
2 4
3
1 4
3 4
1 1
4
1 8 9 6
1 2
2 3
2 4
1
1 3
3
11 5 4
1 2
1 3
3
1 3
1 2
2 2
4
10 15 12 13
1 2
1 3
1 4
3
1 4
2 3
1 4
4
6 2 13 8
1 2
1 3
1 4
3
2 1
4 4
3 3
4
0 13 8 8
1 2
2 3
1 4
3
4 2
1 2
4 3
5
1 13 7 0 11
1 2
1 3
3 4
3 5
3
4 1
2 5
2 3
3
13 2 14
1 2
2 3
3
3 1
1 1
2 1
5
12 8 8 13 11
1 2
1 3
2 4
1 5
2
1 3
2 4
5
13 12 1 14 11
1 2
2 3
2 4
1 5
2
4 5
2 1
2
11 11
1 2
1
1 1
6
15 1 10 0 10 12
1 2
2 3
2 4
2 5
5 6
1
4 3
6
1 5 4 4 15 1
1 2
2 3
1 4
3 5
5 6
1
1 5
3
8 6 8
1 2
2 3
2
3 2
1 2`

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func expectedMoves(n, d int) int64 {
	if d == 0 {
		return 0
	}
	inv := make([]int64, n+1)
	inv[1] = 1
	for i := 2; i <= n; i++ {
		inv[i] = MOD - (MOD/int64(i))*inv[MOD%int64(i)]%MOD
	}
	A := make([]int64, n+1)
	B := make([]int64, n+1)
	if n >= 1 {
		A[1] = 1
	}
	for i := 1; i < n; i++ {
		denom := n - i
		x1 := (int64(n)*A[i] - int64(i)*A[i-1]) % MOD
		if x1 < 0 {
			x1 += MOD
		}
		A[i+1] = x1 * inv[denom] % MOD

		x2 := (int64(n)*B[i] - int64(i)*B[i-1] - int64(n)) % MOD
		if x2 < 0 {
			x2 += MOD
		}
		B[i+1] = x2 * inv[denom] % MOD
	}
	diff := (A[n] - A[n-1]) % MOD
	if diff < 0 {
		diff += MOD
	}
	x := (B[n-1] + 1 - B[n]) % MOD
	if x < 0 {
		x += MOD
	}
	x = x * modPow(diff, MOD-2) % MOD
	res := (A[d]*x + B[d]) % MOD
	return res
}

func solveCase(tc testCase) []int {
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
	for _, v := range order {
		down[v].Add(a[v])
		for _, to := range children[v] {
			down[v].Merge(&down[to])
		}
	}

	up := make([]Basis, n+1)
	type qitem struct{ v int }
	queue := []qitem{{1}}
	for len(queue) > 0 {
		cur := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
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
			queue = append(queue, qitem{to})
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

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseNum+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
		}
		pos++
		vals := make([]int, n+1)
		for i := 1; i <= n; i++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("case %d: missing val", caseNum+1)
			}
			v, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad val: %w", caseNum+1, err)
			}
			vals[i] = v
			pos++
		}
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			if pos+1 >= len(fields) {
				return nil, fmt.Errorf("case %d: missing edge", caseNum+1)
			}
			u, _ := strconv.Atoi(fields[pos])
			v, _ := strconv.Atoi(fields[pos+1])
			pos += 2
			edges[i] = [2]int{u, v}
		}
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing q", caseNum+1)
		}
		q, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad q: %w", caseNum+1, err)
		}
		pos++
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			if pos+1 >= len(fields) {
				return nil, fmt.Errorf("case %d: missing query", caseNum+1)
			}
			r, _ := strconv.Atoi(fields[pos])
			v, _ := strconv.Atoi(fields[pos+1])
			pos += 2
			queries[i] = [2]int{r, v}
		}
		cases = append(cases, testCase{n: n, vals: vals, edges: edges, queries: queries})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	total := 0
	for _, tc := range cases {
		expected := solveCase(tc)
		for qi, val := range tc.queries {
			var sb strings.Builder
			sb.WriteString("1\n")
			sb.WriteString(strconv.Itoa(tc.n))
			sb.WriteByte('\n')
			for i := 1; i <= tc.n; i++ {
				if i > 1 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(tc.vals[i]))
			}
			sb.WriteByte('\n')
			for _, e := range tc.edges {
				sb.WriteString(strconv.Itoa(e[0]))
				sb.WriteByte(' ')
				sb.WriteString(strconv.Itoa(e[1]))
				sb.WriteByte('\n')
			}
			sb.WriteString("1\n")
			sb.WriteString(strconv.Itoa(val[0]))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(val[1]))
			sb.WriteByte('\n')

			got, err := run(bin, sb.String())
			total++
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total, err)
				os.Exit(1)
			}
			if got != strconv.Itoa(expected[qi]) {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", total, expected[qi], got)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", total)
}

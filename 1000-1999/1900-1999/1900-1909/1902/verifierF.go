package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `
3
5 6 0
1 3
1 2
2
3 1 8
2 3 1

2
8 6
1 2
1
1 1 4

4
2 7 4 3
1 3
2 2
1 4
2
4 3 5
1 2 7

4
8 2 1 7
1 2
2 4
2 3
1
3 1 6

2
7 9
1 2
1
1 2 3

4
10 1 6 10
1 2
1 3
2 4
3
4 3 3
2 2 10
1 1 8

6
2 4 5 7 8 3
1 6
1 2
3 4
1 3
3 5
1
1 5 5

6
3 1 8 1 10 7
1 2
1 4
1 3
4 5
3 6
3
6 3 3
5 3 4
6 6 1

5
2 0 10 1 7
1 2
1 5
2 3
3 4
3
2 3 5
3 4 3
1 4 10

6
8 4 6 9 0 8
1 5
1 3
1 6
1 4
5 2
2
2 6 5
2 5 1

3
10 6 5
1 2
1 3
2
2 3 6
1 1 7

2
2 0
1 2
3
1 1 5
1 1 3
2 2 4

6
5 10 8 5 10 7
1 3
1 4
1 5
2 2
1 6
1
5 5 1

6
9 6 7 10 7 6
1 3
2 6
3 5
3 4
4 2
3
5 1 7
5 5 3
4 3 8

4
1 5 9 4
1 2
2 3
1 4
3
3 3 1
4 3 0
1 3 9

2
8 3
1 2
1
1 2 7

2
8 10
1 2
2
1 2 8
2 2 3

5
9 7 7 4 5
1 5
2 3
3 4
1 2
1
4 4 4

4
5 5 0 6
1 2
2 3
2 4
1
3 3 4

6
10 5 1 5 0 7
1 2
2 4
1 5
4 3
2 6
1
3 2 1

2
6 7
1 2
3
1 2 0
1 2 4
1 1 7

2
6 8
1 2
2
1 1 0
1 2 6

6
10 9 4 1 5 7
1 2
1 5
3 3
2 4
4 6
3
6 1 10
4 3 0
5 2 1

6
3 0 8 5 2 3
1 2
1 3
2 4
3 5
2 6
1
3 6 4

4
10 4 3 6
1 3
1 4
3 2
2
2 1 0
2 2 1

3
2 8 2
1 2
1 3
2
2 2 9
1 1 9

2
5 5
1 2
1
2 2 1

5
3 5 9 4 2
1 5
1 2
2 4
1 3
3
2 3 4
4 3 7
2 4 6

4
1 5 2 4
1 4
1 3
2 2
2
3 3 2
4 2 3

3
9 7 9
1 2
1 3
2
1 1 7
2 3 1

2
6 4
1 2
2
1 1 6
2 2 9

5
6 6 0 4 2
1 4
1 2
3 5
4 3
1
2 3 0

5
5 7 5 8 3
1 3
1 2
2 5
2 4
2
1 3 9
4 2 0

3
9 6 2
1 2
1 3
2
3 3 7
2 3 7

6
2 0 9 6 10 6
1 5
2 3
2 2
2 6
4 4
1
6 2 7

5
0 6 3 8 4
1 4
1 3
1 2
2 5
2
2 2 5
2 5 9

3
0 1 2
1 2
2 3
1
1 2 9

4
1 6 4 5
1 2
1 3
1 4
2
4 4 5
2 1 6

5
10 1 9 4 8
1 3
1 2
1 5
4 4
2
1 3 5
5 5 9

5
8 1 1 9 8
1 4
1 5
1 2
1 3
1
3 2 6

6
10 1 6 3 5 1
1 4
1 6
2 5
4 3
2 2
1
1 5 4

5
3 8 1 9 5
1 5
1 2
3 3
1 4
3
1 3 6
1 1 1
3 3 1

6
4 1 3 4 9 1
1 4
1 6
1 3
2 5
2 2
1
6 5 0

5
7 8 8 3 6
1 5
1 4
3 2
3 3
1
4 4 3

5
4 7 3 7 10
1 5
2 2
2 4
3 3
2
2 2 9
5 4 1

3
9 5 10
1 2
2 3
2
1 2 3
2 3 5

2
5 0
1 2
3
1 1 0
1 2 9
2 2 5

2
8 6
1 2
2
2 1 5
1 1 5

2
9 0
1 2
3
1 1 9
2 2 1
1 1 1

3
8 3 3
1 2
2 3
1
1 3 8

5
1 3 2 6 10
1 3
1 4
1 2
4 5
3
4 4 4
5 2 8
1 5 10

5
2 7 9 2 0
1 2
2 5
2 3
2 4
1
3 5 8

3
5 4 3
1 2
2 3
2
2 3 5
1 1 4

4
3 1 3 6
1 3
1 2
2 4
2
2 3 8
3 4 1

6
7 7 10 1 2 10
1 5
2 6
1 3
1 4
5 2
1
1 4 6

4
6 6 0 4
1 4
1 2
1 3
3
2 2 5
4 2 8
3 4 3

2
4 2
1 2
1
1 1 0

5
4 4 8 0 6
1 3
1 2
3 4
1 5
3
4 3 1
3 1 10
4 5 3

4
6 0 6 7
1 4
1 2
3 3
3
1 2 10
4 3 4
2 1 5

4
1 9 9 2
1 4
2 3
2 2
2
4 2 0
2 2 6

3
4 9 8
1 3
1 2
1
3 3 5

4
7 9 4 4
1 2
1 3
1 4
3
4 1 1
2 2 5
4 2 8

4
5 10 5 4
1 3
2 2
3 4
1
4 4 1

2
7 10
1 2
2
1 2 8
1 2 6

4
4 2 1 10
1 2
1 3
3 4
3
1 2 7
4 2 8
1 2 4

6
3 4 0 8 5 4
1 6
1 4
2 2
3 3
1 5
1
4 4 8

2
2 10
1 2
2
2 1 4
2 2 9

6
9 0 7 3 3 4
1 4
1 2
1 5
4 6
4 3
2
6 1 2
5 3 4

2
9 3
1 2
2
2 1 9
2 1 10

5
5 3 9 9 0
1 5
2 4
2 3
1 2
1
1 5 6

3
7 3 7
1 2
1 3
1
1 3 7

2
6 10
1 2
3
1 1 9
2 2 9
1 2 10

5
0 2 4 4 8
1 2
2 3
1 5
4 4
2
4 5 5
4 5 1

4
3 6 7 5
1 4
2 2
1 3
2
2 2 10
3 1 2

2
8 9
1 2
3
2 2 2
1 2 5
2 2 0

2
0 3
1 2
3
1 1 6
2 1 1
2 1 2

5
9 0 1 0 4
1 4
1 2
3 5
2 3
3
1 1 7
4 3 2
1 1 10

6
0 0 10 3 6 4
1 4
1 6
2 3
3 5
5 2
3
2 5 10
5 5 3
5 5 4

6
3 2 2 5 1 1
1 2
1 6
2 3
4 5
1 4
1
1 4 8

2
9 4
1 2
1
1 1 0

2
3 5
1 2
3
1 2 4
2 1 1
2 1 4

4
9 9 10 1
1 3
2 4
2 2
2
1 1 6
3 3 4

2
2 7
1 2
1
1 1 6

4
8 10 1 2
1 3
2 4
3 2
3
4 1 10
3 2 7
2 4 8

2
7 3
1 2
3
2 2 9
1 2 10
1 2 6

5
5 3 2 5 6
1 2
1 3
1 4
4 5
1
2 2 0

4
6 10 1 1
1 3
1 2
1 4
2
4 4 5
3 3 5

5
2 7 10 2 10
1 4
2 3
1 2
1 5
3
2 3 5
2 1 4
3 5 0

2
7 10
1 2
1
2 2 1

2
5 9
1 2
2
1 2 1
2 1 10

5
7 8 10 10 5
1 4
2 5
2 3
1 2
1
5 1 10

2
6 1
1 2
1
2 2 9

5
0 8 0 9 3
1 5
1 2
2 4
2 3
1
5 1 8

5
7 6 4 10 7
1 3
1 2
2 5
4 4
2
3 4 3
1 1 5

5
7 3 9 9 0
1 2
1 5
2 4
2 3
2
5 1 4
4 2 1

5
6 3 9 6 3
1 4
1 2
3 3
3 5
1
2 3 1`

type testCase struct {
	input    string
	expected string
}

// Basis implements xor linear basis for numbers up to 20 bits.
type Basis struct {
	b [20]int
}

func (bs *Basis) Add(x int) {
	for i := 19; i >= 0; i-- {
		if (x>>i)&1 == 0 {
			continue
		}
		if bs.b[i] != 0 {
			x ^= bs.b[i]
		} else {
			bs.b[i] = x
			return
		}
	}
}

func (bs *Basis) Merge(o *Basis) {
	for i := 19; i >= 0; i-- {
		if o.b[i] != 0 {
			bs.Add(o.b[i])
		}
	}
}

func (bs *Basis) Contains(x int) bool {
	for i := 19; i >= 0; i-- {
		if (x>>i)&1 == 0 {
			continue
		}
		if bs.b[i] == 0 {
			return false
		}
		x ^= bs.b[i]
	}
	return true
}

type SegTree struct {
	n    int
	tree []Basis
}

func NewSegTree(arr []int) *SegTree {
	n := len(arr) - 1
	st := &SegTree{n: n, tree: make([]Basis, 4*(n+2))}
	st.build(1, 1, n, arr)
	return st
}

func (st *SegTree) build(p, l, r int, arr []int) {
	if l == r {
		if arr[l] != 0 {
			st.tree[p].Add(arr[l])
		}
		return
	}
	mid := (l + r) >> 1
	st.build(p<<1, l, mid, arr)
	st.build(p<<1|1, mid+1, r, arr)
	st.tree[p] = Basis{}
	st.tree[p].Merge(&st.tree[p<<1])
	st.tree[p].Merge(&st.tree[p<<1|1])
}

func (st *SegTree) query(p, l, r, ql, qr int, res *Basis) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		res.Merge(&st.tree[p])
		return
	}
	mid := (l + r) >> 1
	if ql <= mid {
		st.query(p<<1, l, mid, ql, qr, res)
	}
	if qr > mid {
		st.query(p<<1|1, mid+1, r, ql, qr, res)
	}
}

type HLD struct {
	n      int
	adj    [][]int
	parent []int
	depth  []int
	heavy  []int
	head   []int
	pos    []int
	cur    int
}

func NewHLD(n int, adj [][]int) *HLD {
	h := &HLD{n: n, adj: adj}
	h.parent = make([]int, n+1)
	h.depth = make([]int, n+1)
	h.heavy = make([]int, n+1)
	h.head = make([]int, n+1)
	h.pos = make([]int, n+1)
	h.dfs()
	h.decompose()
	return h
}

func (h *HLD) dfs() {
	order := make([]int, 0, h.n)
	stack := []int{1}
	h.parent[1] = 0
	h.depth[1] = 0
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range h.adj[u] {
			if v != h.parent[u] {
				h.parent[v] = u
				h.depth[v] = h.depth[u] + 1
				stack = append(stack, v)
			}
		}
	}
	size := make([]int, h.n+1)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		size[u] = 1
		maxSize := 0
		heavyChild := 0
		for _, v := range h.adj[u] {
			if v != h.parent[u] {
				size[u] += size[v]
				if size[v] > maxSize {
					maxSize = size[v]
					heavyChild = v
				}
			}
		}
		h.heavy[u] = heavyChild
	}
}

func (h *HLD) decompose() {
	h.cur = 1
	type pair struct{ u, head int }
	stack := []pair{{1, 1}}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		u, hd := p.u, p.head
		for {
			h.head[u] = hd
			h.pos[u] = h.cur
			h.cur++
			for i := len(h.adj[u]) - 1; i >= 0; i-- {
				v := h.adj[u][i]
				if v != h.parent[u] && v != h.heavy[u] {
					stack = append(stack, pair{v, v})
				}
			}
			if h.heavy[u] == 0 {
				break
			}
			u = h.heavy[u]
		}
	}
}

func (h *HLD) queryPath(u, v int, st *SegTree) Basis {
	var res Basis
	for h.head[u] != h.head[v] {
		if h.depth[h.head[u]] < h.depth[h.head[v]] {
			u, v = v, u
		}
		st.query(1, 1, st.n, h.pos[h.head[u]], h.pos[u], &res)
		u = h.parent[h.head[u]]
	}
	if h.depth[u] > h.depth[v] {
		u, v = v, u
	}
	st.query(1, 1, st.n, h.pos[u], h.pos[v], &res)
	return res
}

func solveCase(n int, values []int, edges [][]int, queries [][3]int) string {
	allEqual := true
	for i := 1; i < len(values); i++ {
		if values[i] != values[0] {
			allEqual = false
			break
		}
	}
	var sb strings.Builder
	if allEqual {
		val := values[0]
		for i, q := range queries {
			if i > 0 {
				sb.WriteByte('\n')
			}
			if q[2] == 0 || q[2] == val {
				sb.WriteString("YES")
			} else {
				sb.WriteString("NO")
			}
		}
		return sb.String()
	}
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	h := NewHLD(n, adj)
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		arr[h.pos[i]] = values[i-1]
	}
	st := NewSegTree(arr)
	for i, q := range queries {
		bs := h.queryPath(q[0], q[1], st)
		if bs.Contains(q[2]) {
			sb.WriteString("YES")
		} else {
			sb.WriteString("NO")
		}
		if i+1 < len(queries) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func loadCases() ([]testCase, error) {
	blocks := strings.Split(strings.TrimSpace(testcaseData), "\n\n")
	cases := make([]testCase, 0, len(blocks))
	for idx, blk := range blocks {
		tokens := strings.Fields(blk)
		pos := 0
		if len(tokens) < 1 {
			return nil, fmt.Errorf("block %d: empty", idx+1)
		}
		n, err := strconv.Atoi(tokens[pos])
		if err != nil {
			return nil, fmt.Errorf("block %d: bad n: %w", idx+1, err)
		}
		pos++
		if pos+n > len(tokens) {
			return nil, fmt.Errorf("block %d: not enough values", idx+1)
		}
		values := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(tokens[pos+i])
			if err != nil {
				return nil, fmt.Errorf("block %d: bad value: %w", idx+1, err)
			}
			values[i] = val
		}
		pos += n
		edges := make([][]int, 0, n-1)
		for i := 0; i < n-1; i++ {
			if pos+1 >= len(tokens) {
				return nil, fmt.Errorf("block %d: missing edge %d", idx+1, i+1)
			}
			u, _ := strconv.Atoi(tokens[pos])
			v, _ := strconv.Atoi(tokens[pos+1])
			edges = append(edges, []int{u, v})
			pos += 2
		}
		if pos >= len(tokens) {
			return nil, fmt.Errorf("block %d: missing q", idx+1)
		}
		q, _ := strconv.Atoi(tokens[pos])
		pos++
		queries := make([][3]int, q)
		for i := 0; i < q; i++ {
			if pos+2 >= len(tokens) {
				return nil, fmt.Errorf("block %d: missing query %d", idx+1, i+1)
			}
			x, _ := strconv.Atoi(tokens[pos])
			y, _ := strconv.Atoi(tokens[pos+1])
			k, _ := strconv.Atoi(tokens[pos+2])
			queries[i] = [3]int{x, y, k}
			pos += 3
		}
		if pos != len(tokens) {
			return nil, fmt.Errorf("block %d: extra data", idx+1)
		}

		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range values {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		sb.WriteString(strconv.Itoa(q))
		sb.WriteByte('\n')
		for i, qu := range queries {
			sb.WriteString(fmt.Sprintf("%d %d %d", qu[0], qu[1], qu[2]))
			if i+1 < len(queries) {
				sb.WriteByte('\n')
			}
		}
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: solveCase(n, values, edges, queries),
		})
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
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load cases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

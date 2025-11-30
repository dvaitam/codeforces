package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const INF = 1000000007

type TestCase struct {
	n      int
	parent []int
	events []string
}

// Embedded copy of testcasesE.txt.
const testcasesEData = `
3
0 1 2
1
1 2

5
0 1 1 1 4
1
2 4 5 1 0

4
0 1 1 2
1
1 1

5
0 1 2 3 1
5
1 4
2 5 2 3 0
2 4 3 1 1
2 5 1 2 2
1 3

5
0 1 2 2 4
5
1 1
1 4
1 2
1 3
1 5

2
0 1
5
2 2 1 2 0
1 2
2 1 2 2 2
1 2
2 1 2 1 4

3
0 1 1
4
2 3 1 3 0
1 2
1 1
2 3 2 2 0

3
0 1 1
1
2 2 1 3 0

2
0 1
4
1 2
1 1
2 1 2 2 0
1 2

3
0 1 2
4
2 2 1 1 0
1 2
2 2 1 2 2
2 1 3 2 0

3
0 1 2
2
1 1
1 3

5
0 1 2 1 1
4
2 3 4 1 0
1 2
2 3 1 1 1
2 3 2 4 2

3
0 1 1
5
2 3 2 1 0
1 2
1 3
2 2 3 1 3
1 3

5
0 1 2 1 3
5
2 3 1 2 0
2 5 2 3 1
1 1
2 5 3 5 3
2 5 2 1 0

2
0 1
2
1 1
1 2

4
0 1 2 2
1
1 4

3
0 1 2
1
1 2

3
0 1 2
1
2 2 1 3 0

2
0 1
3
2 1 2 2 0
2 2 1 1 0
1 1

3
0 1 2
2
1 1
2 1 3 1 1

5
0 1 2 3 4
3
1 3
1 1
2 3 5 3 1

5
0 1 2 1 1
3
2 4 1 3 0
2 5 4 3 1
1 2

4
0 1 1 2
1
2 1 4 1 0

3
0 1 2
1
1 2

4
0 1 2 1
5
2 1 2 2 0
2 4 1 3 0
2 1 3 3 1
1 2
1 3

2
0 1
2
2 1 2 2 0
2 2 1 1 0

2
0 1
5
2 1 2 2 0
1 1
1 1
2 2 1 2 2
1 2

2
0 1
5
1 2
2 1 2 2 1
2 1 2 1 0
2 2 1 1 2
1 2

3
0 1 2
5
2 3 2 3 0
1 1
2 3 2 1 2
2 1 2 2 2
2 2 1 3 3

2
0 1
5
2 2 1 1 0
1 1
1 2
2 2 1 1 0
2 1 2 1 1

2
0 1
2
2 2 1 2 0
1 2

3
0 1 2
5
2 1 3 2 0
1 3
2 1 3 2 2
1 3
2 1 2 1 1

4
0 1 1 2
3
2 3 4 2 0
1 1
2 4 2 3 0

2
0 1
5
2 2 1 1 0
2 2 1 1 1
2 2 1 2 1
1 1
2 1 2 2 2

3
0 1 1
3
2 3 2 3 0
2 2 3 3 0
1 3

3
0 1 1
4
2 1 3 2 0
2 1 3 1 1
2 3 2 3 2
1 2

4
0 1 1 2
5
1 1
2 4 1 3 0
2 1 2 2 2
1 3
2 2 3 3 0

2
0 1
4
2 1 2 2 0
2 1 2 2 0
1 2
2 2 1 2 1

2
0 1
3
2 1 2 1 0
2 1 2 1 1
1 1

3
0 1 2
3
1 2
1 3
1 3

4
0 1 1 1
4
1 1
2 1 2 1 1
2 4 3 3 0
2 2 1 2 3

3
0 1 2
4
2 1 2 2 0
1 3
1 2
2 1 2 2 3

4
0 1 2 1
2
2 1 4 2 0
2 3 2 1 1

4
0 1 1 3
1
2 2 1 3 0

5
0 1 1 1 4
2
2 2 4 4 0
2 3 2 2 0

3
0 1 2
5
2 3 1 2 0
2 3 2 3 0
2 3 1 2 0
2 1 2 1 0
2 3 2 3 2

2
0 1
1
1 1

5
0 1 2 3 4
1
2 3 1 3 0

2
0 1
1
2 1 2 2 0

4
0 1 1 3
3
1 2
2 2 3 3 1
2 4 1 2 2

5
0 1 2 3 2
2
1 5
1 4

4
0 1 1 1
3
2 3 4 3 0
1 3
2 1 2 3 2

4
0 1 1 2
3
1 3
2 4 1 1 0
1 4

4
0 1 2 2
3
1 4
2 3 2 1 0
1 2

3
0 1 1
4
2 3 1 1 0
2 1 2 1 0
2 1 3 1 2
2 1 3 2 3

4
0 1 1 3
4
2 2 4 4 0
2 2 4 1 1
1 1
2 4 1 4 3

2
0 1
4
1 2
2 1 2 2 1
2 2 1 2 2
1 1

4
0 1 2 3
5
1 4
1 1
2 1 4 3 1
1 3
2 4 3 4 3

2
0 1
3
1 1
1 1
2 2 1 2 1

4
0 1 2 3
3
1 4
1 4
2 4 1 1 0

3
0 1 1
1
1 1

5
0 1 2 3 2
1
1 5

2
0 1
5
1 1
2 1 2 1 1
2 1 2 1 0
2 1 2 1 3
1 2

5
0 1 2 1 2
4
2 2 5 3 0
2 5 3 4 0
2 5 3 4 0
1 3

4
0 1 2 2
1
1 3

4
0 1 2 1
2
1 4
2 1 2 4 1

4
0 1 2 2
4
1 3
2 2 3 2 0
2 4 2 1 2
1 2

3
0 1 2
4
1 1
2 3 2 1 0
2 3 1 3 0
2 1 3 1 2

5
0 1 1 1 3
5
1 2
1 3
1 5
1 1
1 2

2
0 1
2
2 1 2 1 0
2 1 2 2 0

4
0 1 2 1
1
1 4

4
0 1 1 3
3
1 3
2 1 3 2 1
1 4

5
0 1 2 1 3
5
1 5
2 2 4 3 0
1 3
2 3 4 2 3
2 3 1 5 2

2
0 1
4
2 1 2 2 0
2 1 2 2 1
1 1
1 1

2
0 1
4
2 1 2 1 0
2 2 1 1 1
2 1 2 2 2
2 2 1 1 0

4
0 1 1 1
3
1 3
1 1
2 4 3 3 1

5
0 1 1 1 3
1
2 2 5 4 0

2
0 1
4
2 2 1 2 0
1 2
2 2 1 2 2
1 2

3
0 1 2
4
2 3 2 2 0
1 2
1 1
2 1 2 2 2

3
0 1 1
3
1 2
2 3 1 3 1
2 3 1 2 0

3
0 1 1
1
2 3 1 1 0

2
0 1
4
1 1
2 1 2 2 0
2 1 2 1 2
1 2

2
0 1
1
1 2

5
0 1 2 1 2
5
2 3 5 4 0
2 4 2 4 1
1 2
1 3
2 5 1 3 2

3
0 1 2
3
1 2
2 1 2 1 0
1 3

3
0 1 1
5
1 2
1 1
1 1
1 3
1 3

5
0 1 2 3 2
5
2 5 1 2 0
1 2
2 4 3 2 0
1 2
1 5

4
0 1 2 1
4
1 1
2 4 3 1 0
2 3 2 2 2
2 1 4 2 1

5
0 1 1 2 2
5
1 1
2 2 5 4 1
2 2 3 1 2
1 5
2 4 3 4 3

4
0 1 1 3
2
1 1
1 1

5
0 1 2 3 3
2
1 1
2 4 2 2 0

5
0 1 1 1 3
5
2 1 3 1 0
1 5
1 2
2 3 4 1 0
2 5 4 1 3

4
0 1 2 2
4
2 1 4 1 0
1 3
2 3 1 3 1
1 1

2
0 1
5
2 1 2 2 0
2 1 2 1 1
2 1 2 2 2
2 1 2 2 1
1 1

2
0 1
3
2 2 1 1 0
2 2 1 1 0
2 2 1 1 1

4
0 1 1 1
3
2 1 3 2 0
2 3 2 1 1
2 1 3 4 2

3
0 1 2
2
2 1 3 1 0
2 2 1 1 1

5
0 1 2 1 1
1
2 1 2 3 0

2
0 1
1
2 1 2 2 0

5
0 1 2 2 1
2
2 2 5 3 0
2 3 1 4 0
`

func parseTestcases() ([]TestCase, error) {
	sc := bufio.NewScanner(strings.NewReader(testcasesEData))
	nextLine := func() (string, bool) {
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line != "" {
				return line, true
			}
		}
		return "", false
	}

	var cases []TestCase
	for {
		line, ok := nextLine()
		if !ok {
			break
		}
		var n int
		if _, err := fmt.Sscan(line, &n); err != nil {
			return nil, fmt.Errorf("parse n: %w", err)
		}
		parentLine, ok := nextLine()
		if !ok {
			return nil, fmt.Errorf("missing parent line")
		}
		parentFields := strings.Fields(parentLine)
		if len(parentFields) != n {
			return nil, fmt.Errorf("expected %d parents got %d", n, len(parentFields))
		}
		parent := make([]int, n)
		for i, f := range parentFields {
			if _, err := fmt.Sscan(f, &parent[i]); err != nil {
				return nil, fmt.Errorf("parent[%d]: %w", i, err)
			}
		}
		mLine, ok := nextLine()
		if !ok {
			return nil, fmt.Errorf("missing m line")
		}
		var m int
		if _, err := fmt.Sscan(mLine, &m); err != nil {
			return nil, fmt.Errorf("parse m: %w", err)
		}
		events := make([]string, 0, m)
		for i := 0; i < m; i++ {
			ev, ok := nextLine()
			if !ok {
				return nil, fmt.Errorf("missing event %d", i)
			}
			events = append(events, ev)
		}
		cases = append(cases, TestCase{n: n, parent: parent, events: events})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func solveCase(tc TestCase) []int {
	n := tc.n
	adj := make([][]int, n)
	root := 0
	for i, p := range tc.parent {
		if p > 0 {
			p--
			adj[i] = append(adj[i], p)
			adj[p] = append(adj[p], i)
		} else {
			root = i
		}
	}

	parent := make([]int, n)
	depth := make([]int, n)
	heavy := make([]int, n)
	head := make([]int, n)
	pos := make([]int, n)
	sz := make([]int, n)
	posToNode := make([]int, n)
	for i := range heavy {
		heavy[i] = -1
	}
	var curPos int

	var dfs func(u, p int)
	dfs = func(u, p int) {
		parent[u] = p
		sz[u] = 1
		maxSz := 0
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			depth[v] = depth[u] + 1
			dfs(v, u)
			if sz[v] > maxSz {
				maxSz = sz[v]
				heavy[u] = v
			}
			sz[u] += sz[v]
		}
	}

	var decompose func(u, h int)
	decompose = func(u, h int) {
		head[u] = h
		pos[u] = curPos
		posToNode[curPos] = u
		curPos++
		if heavy[u] != -1 {
			decompose(heavy[u], h)
		}
		for _, v := range adj[u] {
			if v != parent[u] && v != heavy[u] {
				decompose(v, v)
			}
		}
	}

	dfs(root, -1)
	decompose(root, root)

	type Query struct{ a, b, k, y, t, idx int }
	attackTime := make([]int, n)
	for i := range attackTime {
		attackTime[i] = INF
	}
	var queries []Query
	answers := make([]int, 0, len(tc.events))
	for year, ev := range tc.events {
		fields := strings.Fields(ev)
		if len(fields) == 0 {
			continue
		}
		tp := 0
		fmt.Sscan(fields[0], &tp)
		if tp == 1 {
			var c int
			fmt.Sscan(ev, &tp, &c)
			attackTime[c-1] = year + 1
		} else {
			var a, b, k, y int
			fmt.Sscan(ev, &tp, &a, &b, &k, &y)
			queries = append(queries, Query{a - 1, b - 1, k, y, year + 1, len(answers)})
			answers = append(answers, -1)
		}
	}

	base := make([]int, n)
	for u := 0; u < n; u++ {
		base[pos[u]] = attackTime[u]
	}

	st := make([][]int, 4*n)
	var build func(node, l, r int)
	build = func(node, l, r int) {
		if l == r {
			st[node] = []int{base[l]}
			return
		}
		mid := (l + r) >> 1
		build(node<<1, l, mid)
		build(node<<1|1, mid+1, r)
		a, b := st[node<<1], st[node<<1|1]
		st[node] = make([]int, len(a)+len(b))
		i, j := 0, 0
		for k := 0; k < len(st[node]); k++ {
			if j >= len(b) || (i < len(a) && a[i] <= b[j]) {
				st[node][k] = a[i]
				i++
			} else {
				st[node][k] = b[j]
				j++
			}
		}
	}
	build(1, 0, n-1)

	var queryLE func(node, l, r, ql, qr, x int) int
	queryLE = func(node, l, r, ql, qr, x int) int {
		if qr < l || r < ql {
			return 0
		}
		if ql <= l && r <= qr {
			arr := st[node]
			lo, hi := 0, len(arr)
			for lo < hi {
				mid := (lo + hi) >> 1
				if arr[mid] <= x {
					lo = mid + 1
				} else {
					hi = mid
				}
			}
			return lo
		}
		mid := (l + r) >> 1
		return queryLE(node<<1, l, mid, ql, qr, x) + queryLE(node<<1|1, mid+1, r, ql, qr, x)
	}

	countBad := func(l, r, y, t int) int {
		if l > r {
			return 0
		}
		c1 := queryLE(1, 0, n-1, l, r, t)
		c2 := queryLE(1, 0, n-1, l, r, y)
		return c1 - c2
	}

	var lca func(u, v int) int
	lca = func(u, v int) int {
		for head[u] != head[v] {
			if depth[head[u]] > depth[head[v]] {
				u = parent[head[u]]
			} else {
				v = parent[head[v]]
			}
		}
		if depth[u] < depth[v] {
			return u
		}
		return v
	}

	for _, q := range queries {
		a, b, k, y, t := q.a, q.b, q.k, q.y, q.t
		type Seg struct{ l, r int; rev bool }
		var segs []Seg
		u := a
		l := lca(a, b)
		for head[u] != head[l] {
			segs = append(segs, Seg{pos[head[u]], pos[u], true})
			u = parent[head[u]]
		}
		if u != l {
			segs = append(segs, Seg{pos[l] + 1, pos[u], true})
		}
		var down []Seg
		v := b
		for head[v] != head[l] {
			down = append(down, Seg{pos[head[v]], pos[v], false})
			v = parent[head[v]]
		}
		if v != l {
			down = append(down, Seg{pos[l] + 1, pos[v], false})
		}
		for i := len(down) - 1; i >= 0; i-- {
			segs = append(segs, down[i])
		}

		totBad := 0
		totLen := 0
		for _, s := range segs {
			length := s.r - s.l + 1
			if length <= 0 {
				continue
			}
			totLen += length
			totBad += countBad(s.l, s.r, y, t)
		}
		totGood := totLen - totBad
		if k > totGood {
			answers[q.idx] = -1
			continue
		}
		rem := k
		found := -1
		for _, s := range segs {
			length := s.r - s.l + 1
			if length <= 0 {
				continue
			}
			bad := countBad(s.l, s.r, y, t)
			good := length - bad
			if rem > good {
				rem -= good
				continue
			}
			l0, r0 := s.l, s.r
			low, high := 0, r0-l0
			for low < high {
				mid := (low + high) >> 1
				var sl, sr int
				if s.rev {
					sl = r0 - mid
					sr = r0
				} else {
					sl = l0
					sr = l0 + mid
				}
				badm := countBad(sl, sr, y, t)
				goodm := (mid + 1) - badm
				if goodm >= rem {
					high = mid
				} else {
					low = mid + 1
					rem -= goodm
				}
			}
			var pidx int
			if s.rev {
				pidx = s.r - low
			} else {
				pidx = s.l + low
			}
			found = posToNode[pidx]
			break
		}
		if found < 0 {
			answers[q.idx] = -1
		} else {
			answers[q.idx] = found + 1
		}
	}
	return answers
}

func runCase(bin string, tc TestCase, idx int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.parent {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(v))
	}
	input.WriteByte('\n')
	input.WriteString(fmt.Sprintf("%d\n", len(tc.events)))
	for _, ev := range tc.events {
		input.WriteString(ev)
		input.WriteByte('\n')
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}

	gotFields := strings.Fields(strings.TrimSpace(out.String()))
	expect := solveCase(tc)
	if len(gotFields) != len(expect) {
		return fmt.Errorf("expected %d values got %d", len(expect), len(gotFields))
	}
	for i, f := range gotFields {
		var v int
		if _, err := fmt.Sscan(f, &v); err != nil {
			return fmt.Errorf("parse output at %d: %v", i, err)
		}
		if v != expect[i] {
			return fmt.Errorf("at %d expected %d got %d", i, expect[i], v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		if err := runCase(bin, tc, i+1); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type segment struct {
	l, r int
	a    int
	idx  int
	rank int
}

type event struct {
	pos int
	typ int
	idx int
}

type edge struct {
	u, v int
	w    int
}

type bit struct {
	n   int
	bit []int
}

func newBIT(n int) *bit {
	return &bit{n: n, bit: make([]int, n+2)}
}

func (b *bit) add(i, delta int) {
	for i <= b.n {
		b.bit[i] += delta
		i += i & -i
	}
}

func (b *bit) sum(i int) int {
	s := 0
	for i > 0 {
		s += b.bit[i]
		i &= i - 1
	}
	return s
}

func (b *bit) kth(k int) int {
	idx := 0
	bitMask := 1
	for bitMask<<1 <= b.n {
		bitMask <<= 1
	}
	for step := bitMask; step > 0; step >>= 1 {
		next := idx + step
		if next <= b.n && b.bit[next] < k {
			k -= b.bit[next]
			idx = next
		}
	}
	return idx + 1
}

type rankSet struct {
	cnt int
	m   map[int]struct{}
	any int
}

func (rs *rankSet) insert(id int) {
	if rs.m == nil {
		rs.m = make(map[int]struct{})
	}
	rs.m[id] = struct{}{}
	rs.cnt++
	rs.any = id
}

func (rs *rankSet) remove(id int) {
	if rs.m == nil {
		return
	}
	delete(rs.m, id)
	rs.cnt--
	if rs.cnt == 0 {
		rs.any = -1
	} else if rs.any == id {
		for k := range rs.m {
			rs.any = k
			break
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesE = `100
5
3 3 4
0 1 8
5 7 5
3 5 5
0 2 2
5
5 6 1
6 8 6
1 5 9
3 3 4
8 9 2
4
6 8 6
4 8 9
2 6 6
6 7 8
5
1 4 1
7 11 1
4 7 6
9 13 5
4 5 2
4
6 6 10
3 4 6
1 5 8
4 5 8
2
4 6 8
9 12 6
4
1 2 6
6 6 3
5 7 2
1 2 6
3
3 7 3
6 6 7
0 0 4
4
1 4 6
6 9 2
4 7 4
4 4 7
3
6 10 1
5 5 2
7 10 9
4
4 8 8
6 6 2
2 2 10
0 2 9
5
1 3 3
1 5 3
9 11 8
3 4 3
4 5 9
5
2 4 2
6 6 5
5 8 4
9 9 9
0 2 7
5
5 5 9
3 3 8
0 1 2
0 3 10
1 3 2
4
1 1 4
3 6 1
8 8 8
1 2 10
5
8 12 2
9 12 8
7 11 1
6 6 2
1 1 1
3
2 2 7
8 10 8
7 9 4
5
9 9 1
2 6 7
1 5 4
3 3 2
6 8 1
5
7 7 1
5 6 10
9 12 8
9 11 10
8 12 7
3
1 1 6
7 11 7
2 2 6
2
9 13 9
2 4 3
5
8 8 3
0 4 4
4 6 10
6 8 6
8 10 3
3
8 11 9
5 9 5
4 5 6
4
3 7 10
6 9 10
8 9 6
7 7 8
4
1 4 6
9 9 2
4 8 3
1 3 1
2
3 7 6
0 3 10
5
7 11 9
6 7 2
9 9 4
3 3 2
0 1 9
5
9 9 2
1 1 4
4 7 4
6 8 8
7 10 3
4
3 7 4
7 7 8
5 8 9
0 1 5
5
2 4 2
9 11 5
1 5 2
8 12 5
0 4 5
5
1 4 2
4 5 8
1 1 8
0 3 2
5 5 2
5
5 9 9
1 4 5
2 4 5
3 5 5
6 10 4
4
1 3 2
5 7 5
6 7 4
7 10 1
3
4 6 2
7 7 10
2 6 6
5
1 2 9
4 8 6
1 4 3
5 8 1
3 5 7
4
8 8 9
5 9 3
9 10 7
3 3 1
3
6 9 2
2 5 7
1 4 9
5
5 5 3
2 6 5
5 6 4
3 6 7
0 4 1
2
8 10 10
1 2 10
3
0 1 6
4 4 7
6 10 3
4
3 4 3
6 7 7
4 7 4
9 10 8
4
0 1 8
7 11 9
0 3 1
9 13 10
5
4 4 3
2 2 2
9 10 6
3 4 2
5 7 10
3
9 10 8
8 10 7
3 5 5
4
7 11 2
6 8 2
8 12 10
0 0 9
3
5 9 1
0 3 2
6 9 8
2
4 4 3
1 2 7
5
4 7 2
4 8 1
9 11 8
1 1 9
9 9 9
5
0 4 2
5 8 2
0 4 8
3 7 8
8 8 5
3
9 11 8
4 4 5
6 9 9
3
7 7 6
7 7 7
2 5 2
5
5 6 5
9 10 1
4 6 8
9 9 5
0 0 4
2
8 11 5
9 10 7
4
6 8 7
6 6 8
3 4 1
8 10 8
4
0 4 7
5 6 10
4 6 6
8 10 4
2
4 4 1
2 6 7
3
2 4 7
6 8 1
3 5 2
3
2 2 10
2 3 5
7 10 8
4
8 12 10
7 10 9
1 1 10
2 3 9
2
4 6 9
9 10 2
3
0 0 9
4 4 9
1 3 8
2
7 9 2
5 9 10
2
7 11 1
4 7 9
3
2 5 10
8 12 5
8 9 4
4
4 4 10
3 3 8
0 4 2
4 4 10
2
9 11 4
8 9 7
3
3 7 5
3 7 10
2 4 3
5
6 10 8
6 10 5
6 10 3
0 3 3
7 7 6
2
5 9 2
3 6 8
4
3 7 3
5 8 7
2 5 4
2 3 6
3
1 5 4
7 7 2
3 5 9
3
7 8 8
4 7 1
5 7 10
5
9 12 3
2 4 9
0 3 1
3 3 3
7 11 3
3
7 11 5
7 8 4
8 11 4
3
9 10 6
1 3 4
2 6 2
2
0 3 8
3 5 8
2
4 4 6
1 2 9
3
5 7 4
0 4 10
4 6 2
2
8 8 4
7 8 5
4
0 1 4
4 6 3
8 12 7
8 11 2
5
1 1 3
9 13 3
2 2 2
8 12 8
1 4 6
2
4 5 1
9 9 8
2
7 10 7
1 1 6
4
2 6 2
7 10 5
8 12 1
4 5 1
5
2 2 7
4 4 8
1 5 8
3 7 3
8 10 5
3
3 7 8
0 2 4
8 11 9
5
9 11 6
0 4 4
8 11 10
8 11 10
8 12 6
5
9 9 9
4 5 6
6 8 4
8 9 8
4 7 10
4
5 7 2
1 3 8
4 5 4
7 8 9
4
9 13 10
1 4 5
2 2 1
1 3 4
2
7 8 6
9 13 1
5
1 1 3
1 5 4
7 11 5
9 13 1
9 13 8
3
2 2 1
2 5 10
5 8 9
2
2 4 9
0 2 7
2
0 3 10
1 5 1
2
1 4 3
5 6 4
5
1 2 9
1 4 6
5 9 10
4 4 5
7 8 7
3
2 5 2
7 11 6
9 9 9
5
8 8 1
2 2 4
5 6 4
1 2 2
9 9 10
2
8 11 5
8 12 3`

func solveOne(segs []segment) int64 {
	n := len(segs)
	alls := make([]int, n)
	for i := range segs {
		alls[i] = segs[i].a
		segs[i].idx = i
	}
	tmp := make([]segment, n)
	copy(tmp, segs)
	sort.Slice(tmp, func(i, j int) bool { return tmp[i].l < tmp[j].l })
	curR := tmp[0].r
	for i := 1; i < n; i++ {
		if tmp[i].l > curR {
			return -1
		}
		if tmp[i].r > curR {
			curR = tmp[i].r
		}
	}

	sort.Ints(alls)
	uniq := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if i == 0 || alls[i] != alls[i-1] {
			uniq = append(uniq, alls[i])
		}
	}
	for i := 0; i < n; i++ {
		segs[i].rank = sort.SearchInts(uniq, segs[i].a) + 1 // 1-based
	}
	m := len(uniq)

	evs := make([]event, 0, 2*n)
	for i := 0; i < n; i++ {
		evs = append(evs, event{segs[i].l, 0, i})
		evs = append(evs, event{segs[i].r, 1, i})
	}
	sort.Slice(evs, func(i, j int) bool {
		if evs[i].pos == evs[j].pos {
			return evs[i].typ < evs[j].typ // start before end
		}
		return evs[i].pos < evs[j].pos
	})

	b := newBIT(m)
	rset := make([]rankSet, m+2)
	edges := make([]edge, 0, 2*n)

	getPred := func(rank int) int {
		cnt := b.sum(rank - 1)
		if cnt == 0 {
			return -1
		}
		pr := b.kth(cnt)
		return rset[pr].any
	}
	getSucc := func(rank int) int {
		total := b.sum(m)
		cnt := b.sum(rank)
		if cnt == total {
			return -1
		}
		su := b.kth(cnt + 1)
		return rset[su].any
	}

	for _, e := range evs {
		idx := e.idx
		r := segs[idx].rank
		if e.typ == 0 {
			if rset[r].cnt > 0 {
				edges = append(edges, edge{idx, rset[r].any, 0})
			}
			if p := getPred(r); p != -1 {
				w := abs(segs[idx].a - segs[p].a)
				edges = append(edges, edge{idx, p, w})
			}
			if s := getSucc(r); s != -1 {
				w := abs(segs[idx].a - segs[s].a)
				edges = append(edges, edge{idx, s, w})
			}
			rset[r].insert(idx)
			b.add(r, 1)
		} else {
			rset[r].remove(idx)
			b.add(r, -1)
			if p, s := getPred(r), getSucc(r); p != -1 && s != -1 {
				w := abs(segs[p].a - segs[s].a)
				edges = append(edges, edge{p, s, w})
			}
		}
	}

	sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })

	parent := make([]int, n)
	size := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	unite := func(a, b int) bool {
		ra, rb := find(a), find(b)
		if ra == rb {
			return false
		}
		if size[ra] < size[rb] {
			ra, rb = rb, ra
		}
		parent[rb] = ra
		size[ra] += size[rb]
		return true
	}

	cnt := 0
	var ans int64
	for _, e := range edges {
		if unite(e.u, e.v) {
			ans += int64(e.w)
			cnt++
			if cnt == n-1 {
				break
			}
		}
	}
	if cnt < n-1 {
		return -1
	}
	return ans
}

// Embedded solution from 1981E.go; used to derive expected outputs without
// calling an external oracle binary.
func embeddedSolution(input string) (string, error) {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return "", err
	}

	var sb strings.Builder

	for ; t > 0; t-- {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return "", err
		}
		segs := make([]segment, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(in, &segs[i].l, &segs[i].r, &segs[i].a); err != nil {
				return "", err
			}
			segs[i].idx = i
		}
		fmt.Fprintf(&sb, "%d\n", solveOne(segs))
	}

	return strings.TrimSpace(sb.String()), nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	input := strings.TrimSpace(testcasesE) + "\n"
	expected, err := embeddedSolution(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "embedded solver error:", err)
		os.Exit(1)
	}

	allOutput, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, allOutput)
		os.Exit(1)
	}

	outLines := strings.Split(strings.TrimSpace(allOutput), "\n")
	expectedLines := strings.Split(expected, "\n")
	if len(outLines) != len(expectedLines) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(expectedLines), len(outLines))
		os.Exit(1)
	}
	for i := range expectedLines {
		got := strings.TrimSpace(outLines[i])
		want := strings.TrimSpace(expectedLines[i])
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(expectedLines))
}

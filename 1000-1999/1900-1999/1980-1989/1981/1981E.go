package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type segment struct {
	l, r int
	a    int
	idx  int
	rank int
}

type event struct {
	pos int
	typ int //0=start,1=end
	idx int
}

type edge struct {
	u, v int
	w    int
}

// Fenwick tree for prefix sums and kth order statistics
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

// structure to keep active indices by rank
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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		segs := make([]segment, n)
		alls := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &segs[i].l, &segs[i].r, &segs[i].a)
			segs[i].idx = i
			alls[i] = segs[i].a
		}

		// connectivity check by union of intervals
		tmp := make([]segment, n)
		copy(tmp, segs)
		sort.Slice(tmp, func(i, j int) bool { return tmp[i].l < tmp[j].l })
		curR := tmp[0].r
		connected := true
		for i := 1; i < n; i++ {
			if tmp[i].l > curR {
				connected = false
				break
			}
			if tmp[i].r > curR {
				curR = tmp[i].r
			}
		}
		if !connected {
			fmt.Fprintln(out, -1)
			continue
		}

		// compress a values
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

		// build events
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

		// prepare structures
		b := newBIT(m)
		rset := make([]rankSet, m+2)
		edges := make([]edge, 0)

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
			if e.typ == 0 { // start
				if rset[r].cnt > 0 {
					edges = append(edges, edge{idx, rset[r].any, 0})
				}
				if p := getPred(r); p != -1 {
					w := segs[idx].a - segs[p].a
					if w < 0 {
						w = -w
					}
					edges = append(edges, edge{idx, p, w})
				}
				if s := getSucc(r); s != -1 {
					w := segs[idx].a - segs[s].a
					if w < 0 {
						w = -w
					}
					edges = append(edges, edge{idx, s, w})
				}
				rset[r].insert(idx)
				b.add(r, 1)
			} else { // end
				rset[r].remove(idx)
				b.add(r, -1)
				if p, s := getPred(r), getSucc(r); p != -1 && s != -1 {
					w := segs[p].a - segs[s].a
					if w < 0 {
						w = -w
					}
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
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, ans)
		}
	}
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"time"
)

type Fenwick struct {
	n    int
	tree []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int, n+1)}
}

func (f *Fenwick) Update(i, v int) {
	i++
	for ; i <= f.n; i += i & -i {
		f.tree[i] += v
	}
}

func (f *Fenwick) Query(i int) int {
	i++
	s := 0
	for ; i > 0; i -= i & -i {
		s += f.tree[i]
	}
	return s
}

type SegTree struct {
	n    int
	size int
	tree []int
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	return &SegTree{n: n, size: size, tree: make([]int, 2*size)}
}

func (s *SegTree) Update(i, v int) {
	i += s.size
	if s.tree[i] >= v {
		return
	}
	s.tree[i] = v
	for i >>= 1; i > 0; i >>= 1 {
		if s.tree[2*i] > s.tree[2*i+1] {
			s.tree[i] = s.tree[2*i]
		} else {
			s.tree[i] = s.tree[2*i+1]
		}
	}
}

func (s *SegTree) Query(l, r int) int {
	if l > r {
		return 0
	}
	l += s.size
	r += s.size
	res := 0
	for l <= r {
		if l&1 == 1 {
			if s.tree[l] > res {
				res = s.tree[l]
			}
			l++
		}
		if r&1 == 0 {
			if s.tree[r] > res {
				res = s.tree[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func uniqueInts(a []int) []int {
	n := 0
	for i := 0; i < len(a); i++ {
		if n == 0 || a[i] != a[n-1] {
			a[n] = a[i]
			n++
		}
	}
	return a[:n]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveE(n, ageK int, r []int, a []int, queries [][2]int) []int {
	ages := make([]int, n)
	copy(ages, a)
	sort.Ints(ages)
	ages = uniqueInts(ages)
	m := len(ages)
	idxAge := make([]int, n)
	for i := 0; i < n; i++ {
		idxAge[i] = sort.SearchInts(ages, a[i])
	}
	c := make([]int, n)
	fenw := NewFenwick(m)
	ord := make([]int, n)
	for i := range ord {
		ord[i] = i
	}
	sort.Slice(ord, func(i, j int) bool { return r[ord[i]] < r[ord[j]] })
	for i := 0; i < n; {
		j := i
		for j < n && r[ord[j]] == r[ord[i]] {
			j++
		}
		for k := i; k < j; k++ {
			fenw.Update(idxAge[ord[k]], 1)
		}
		for k := i; k < j; k++ {
			u := ord[k]
			lo := a[u] - ageK
			hi := a[u] + ageK
			lidx := sort.SearchInts(ages, lo)
			ridx := sort.Search(len(ages), func(x int) bool { return ages[x] > hi }) - 1
			if lidx < 0 {
				lidx = 0
			}
			if ridx >= 0 && lidx <= ridx {
				c[u] = fenw.Query(ridx)
				if lidx > 0 {
					c[u] -= fenw.Query(lidx - 1)
				}
			} else {
				c[u] = 0
			}
		}
		i = j
	}
	q := len(queries)
	type Query struct{ R0, lidx, ridx, idx int }
	qs := make([]Query, q)
	for i, qu := range queries {
		x, y := qu[0], qu[1]
		R0 := r[x]
		if r[y] > R0 {
			R0 = r[y]
		}
		lowAge := max(a[x]-ageK, a[y]-ageK)
		highAge := min(a[x]+ageK, a[y]+ageK)
		lidx := sort.SearchInts(ages, lowAge)
		ridx := sort.Search(len(ages), func(u int) bool { return ages[u] > highAge }) - 1
		qs[i] = Query{R0: R0, lidx: lidx, ridx: ridx, idx: i}
	}
	ordDesc := make([]int, n)
	copy(ordDesc, ord)
	sort.Slice(ordDesc, func(i, j int) bool { return r[ordDesc[i]] > r[ordDesc[j]] })
	sort.Slice(qs, func(i, j int) bool { return qs[i].R0 > qs[j].R0 })
	seg := NewSegTree(m)
	ans := make([]int, q)
	pi := 0
	for _, qu := range qs {
		for pi < n && r[ordDesc[pi]] >= qu.R0 {
			u := ordDesc[pi]
			seg.Update(idxAge[u], c[u])
			pi++
		}
		if qu.lidx > qu.ridx {
			ans[qu.idx] = -1
		} else {
			mv := seg.Query(qu.lidx, qu.ridx)
			if mv == 0 {
				ans[qu.idx] = -1
			} else {
				ans[qu.idx] = mv
			}
		}
	}
	return ans
}

func runCase(bin string, n, ageK int, r, a []int, queries [][2]int) error {
	var input bytes.Buffer
	fmt.Fprintf(&input, "%d %d\n", n, ageK)
	for i := 0; i < n; i++ {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprintf(&input, "%d", r[i])
	}
	input.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprintf(&input, "%d", a[i])
	}
	input.WriteByte('\n')
	q := len(queries)
	fmt.Fprintf(&input, "%d\n", q)
	for _, qu := range queries {
		fmt.Fprintf(&input, "%d %d\n", qu[0]+1, qu[1]+1)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(&out)
	got := make([]int, 0, q)
	for scanner.Scan() {
		var v int
		if _, err := fmt.Sscan(scanner.Text(), &v); err != nil {
			return fmt.Errorf("parse output: %v", err)
		}
		got = append(got, v)
	}
	if len(got) != q {
		return fmt.Errorf("expected %d lines, got %d", q, len(got))
	}
	want := solveE(n, ageK, r, a, queries)
	for i := 0; i < q; i++ {
		if got[i] != want[i] {
			return fmt.Errorf("expected %d got %d", want[i], got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	const tests = 100
	for t := 0; t < tests; t++ {
		n := rng.Intn(10) + 2
		ageK := rng.Intn(5) + 1
		r := make([]int, n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			r[i] = rng.Intn(20) + 1
			a[i] = rng.Intn(20) + 1
		}
		q := rng.Intn(5) + 1
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			x := rng.Intn(n)
			y := rng.Intn(n)
			for y == x {
				y = rng.Intn(n)
			}
			queries[i] = [2]int{x, y}
		}
		if err := runCase(bin, n, ageK, r, a, queries); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

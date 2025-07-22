package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Fenwick tree for sum queries and point updates
type Fenwick struct {
	n int
	t []int
}

func NewFenwick(n int) *Fenwick { return &Fenwick{n: n, t: make([]int, n+1)} }
func (f *Fenwick) Add(i, v int) {
	for ; i <= f.n; i += i & -i {
		f.t[i] += v
	}
}
func (f *Fenwick) Sum(i int) int {
	s := 0
	for ; i > 0; i -= i & -i {
		s += f.t[i]
	}
	return s
}
func (f *Fenwick) RangeSum(l, r int) int {
	if r < l {
		return 0
	}
	return f.Sum(r) - f.Sum(l-1)
}

type person struct {
	name   string
	parent int
}

type Query struct{ v, k int }

func expected(n int, people []person, queries []Query) []int {
	nameMap := make(map[string]int)
	nextID := 1
	nameID := make([]int, n+1)
	children := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		p := people[i-1]
		if p.parent != 0 {
			children[p.parent] = append(children[p.parent], i)
		}
		id, ok := nameMap[p.name]
		if !ok {
			id = nextID
			nameMap[p.name] = id
			nextID++
		}
		nameID[i] = id
	}
	depth := make([]int, n+1)
	tin := make([]int, n+1)
	tout := make([]int, n+1)
	tins := make([][]int, n+2)
	namesAtDepth := make([][]int, n+2)
	timev := 1
	type Frame struct{ u, idx int }
	stack := make([]Frame, 0, n)
	for i := 1; i <= n; i++ {
		if people[i-1].parent != 0 {
			continue
		}
		depth[i] = 0
		stack = append(stack, Frame{i, -1})
		for len(stack) > 0 {
			top := &stack[len(stack)-1]
			u := top.u
			if top.idx == -1 {
				tin[u] = timev
				timev++
				d := depth[u]
				tins[d] = append(tins[d], tin[u])
				namesAtDepth[d] = append(namesAtDepth[d], nameID[u])
				top.idx = 0
			} else if top.idx < len(children[u]) {
				v := children[u][top.idx]
				top.idx++
				depth[v] = depth[u] + 1
				stack = append(stack, Frame{v, -1})
			} else {
				tout[u] = timev - 1
				stack = stack[:len(stack)-1]
			}
		}
	}
	queriesByDepth := make([][]struct{ l, r, id int }, n+2)
	ans := make([]int, len(queries))
	for qi, q := range queries {
		td := depth[q.v] + q.k
		if td >= len(tins) || len(tins[td]) == 0 {
			ans[qi] = 0
			continue
		}
		arr := tins[td]
		l := lowerBound(arr, tin[q.v])
		r := upperBound(arr, tout[q.v]) - 1
		if l > r {
			ans[qi] = 0
		} else {
			queriesByDepth[td] = append(queriesByDepth[td], struct{ l, r, id int }{l + 1, r + 1, qi})
		}
	}
	maxNameID := nextID
	for d, qs := range queriesByDepth {
		if len(qs) == 0 {
			continue
		}
		sort.Slice(qs, func(i, j int) bool { return qs[i].r < qs[j].r })
		cnt := len(namesAtDepth[d])
		bit := NewFenwick(cnt)
		last := make([]int, maxNameID)
		qi := 0
		for i := 1; i <= cnt; i++ {
			id := namesAtDepth[d][i-1]
			if last[id] != 0 {
				bit.Add(last[id], -1)
			}
			bit.Add(i, 1)
			last[id] = i
			for qi < len(qs) && qs[qi].r == i {
				l := qs[qi].l
				r := qs[qi].r
				ans[qs[qi].id] = bit.RangeSum(l, r)
				qi++
			}
		}
	}
	return ans
}

func lowerBound(a []int, x int) int {
	l, r := 0, len(a)
	for l < r {
		m := l + (r-l)/2
		if a[m] < x {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}
func upperBound(a []int, x int) int {
	l, r := 0, len(a)
	for l < r {
		m := l + (r-l)/2
		if a[m] <= x {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

func runCase(bin string, n int, people []person, qs []Query) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, p := range people {
		sb.WriteString(p.name)
		sb.WriteByte(' ')
		sb.WriteString(fmt.Sprint(p.parent))
		if i+1 < len(people) {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(qs)))
	for i, q := range qs {
		sb.WriteString(fmt.Sprintf("%d %d", q.v, q.k))
		if i+1 < len(qs) {
			sb.WriteByte('\n')
		}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expect := expected(n, people, qs)
	gotFields := strings.Fields(strings.TrimSpace(out.String()))
	if len(gotFields) != len(expect) {
		return fmt.Errorf("expected %d numbers got %d", len(expect), len(gotFields))
	}
	for i, f := range gotFields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("non-int output %q", f)
		}
		if v != expect[i] {
			return fmt.Errorf("at %d expected %d got %d", i, expect[i], v)
		}
	}
	return nil
}

func randName(rng *rand.Rand) string {
	l := rng.Intn(5) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(15) + 1
		people := make([]person, n)
		for j := 0; j < n; j++ {
			name := randName(rng)
			parent := 0
			if j > 0 {
				parent = rng.Intn(j) + 1
			}
			people[j] = person{name, parent}
		}
		m := rng.Intn(15) + 1
		queries := make([]Query, m)
		for j := 0; j < m; j++ {
			v := rng.Intn(n) + 1
			k := rng.Intn(n)
			queries[j] = Query{v, k}
		}
		if err := runCase(bin, n, people, queries); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

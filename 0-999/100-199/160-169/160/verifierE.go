package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Bus struct {
	s, f int
	t    int
	id   int
}

type Query struct {
	l, r, b   int
	id        int
	low, high int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(46)
	var tests []string
	for t := 0; t < 100; t++ {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		buses := make([]Bus, n)
		usedT := make(map[int]bool)
		for i := 0; i < n; i++ {
			s := rand.Intn(20) + 1
			f := rand.Intn(20-s) + s + 1
			var ti int
			for {
				ti = rand.Intn(100) + 1
				if !usedT[ti] {
					usedT[ti] = true
					break
				}
			}
			buses[i] = Bus{s, f, ti, i + 1}
		}
		people := make([]Query, m)
		for i := 0; i < m; i++ {
			l := rand.Intn(20) + 1
			r := rand.Intn(20-l) + l + 1
			b := rand.Intn(100) + 1
			people[i] = Query{l, r, b, i, 0, n - 1}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, b := range buses {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", b.s, b.f, b.t))
		}
		for _, q := range people {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", q.l, q.r, q.b))
		}
		tests = append(tests, sb.String())
	}
	for i, input := range tests {
		expect := solveE(strings.NewReader(input))
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

// --- solution for problem E ---

type SegTree struct {
	n    int
	tree []int
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	return &SegTree{n: size, tree: make([]int, 2*size)}
}

func (st *SegTree) Update(pos, val int) {
	i := pos + st.n
	if st.tree[i] >= val {
		return
	}
	st.tree[i] = val
	for i >>= 1; i > 0; i >>= 1 {
		if st.tree[2*i] > st.tree[2*i+1] {
			st.tree[i] = st.tree[2*i]
		} else {
			st.tree[i] = st.tree[2*i+1]
		}
	}
}

func (st *SegTree) Query(r int) int {
	l := st.n
	rr := r + st.n
	maxv := 0
	for l <= rr {
		if l&1 == 1 {
			if st.tree[l] > maxv {
				maxv = st.tree[l]
			}
			l++
		}
		if rr&1 == 0 {
			if st.tree[rr] > maxv {
				maxv = st.tree[rr]
			}
			rr--
		}
		l >>= 1
		rr >>= 1
	}
	return maxv
}

func solveE(r io.Reader) string {
	in := bufio.NewReader(r)
	var n, m int
	fmt.Fscan(in, &n, &m)
	buses := make([]Bus, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &buses[i].s, &buses[i].f, &buses[i].t)
		buses[i].id = i + 1
	}
	queries := make([]Query, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &queries[i].l, &queries[i].r, &queries[i].b)
		queries[i].id = i
		queries[i].low = 0
		queries[i].high = n - 1
	}
	sort.Slice(buses, func(i, j int) bool { return buses[i].t < buses[j].t })
	tList := make([]int, n)
	for i := range buses {
		tList[i] = buses[i].t
	}
	for i := range queries {
		lb := sort.Search(n, func(j int) bool { return tList[j] >= queries[i].b })
		if lb >= n {
			queries[i].low = n
			queries[i].high = -1
		} else {
			queries[i].low = lb
			queries[i].high = n - 1
		}
	}
	coords := make([]int, 0, n+m)
	for _, b := range buses {
		coords = append(coords, b.s)
	}
	for _, q := range queries {
		coords = append(coords, q.l)
	}
	sort.Ints(coords)
	uniq := coords[:1]
	for i := 1; i < len(coords); i++ {
		if coords[i] != coords[i-1] {
			uniq = append(uniq, coords[i])
		}
	}
	getPos := func(x int) int { return sort.SearchInts(uniq, x) }
	sPos := make([]int, n)
	for i := range buses {
		sPos[i] = getPos(buses[i].s)
	}
	lPos := make([]int, m)
	for i := range queries {
		lPos[i] = getPos(queries[i].l)
	}
	type Task struct{ mid, qi int }
	for {
		tasks := make([]Task, 0, m)
		for i := range queries {
			if queries[i].low <= queries[i].high {
				mid := (queries[i].low + queries[i].high) >> 1
				tasks = append(tasks, Task{mid, i})
			}
		}
		if len(tasks) == 0 {
			break
		}
		sort.Slice(tasks, func(i, j int) bool { return tasks[i].mid < tasks[j].mid })
		st := NewSegTree(len(uniq))
		p := -1
		for _, t := range tasks {
			mid, qi := t.mid, t.qi
			for p < mid {
				p++
				st.Update(sPos[p], buses[p].f)
			}
			if st.Query(lPos[qi]) >= queries[qi].r {
				queries[qi].high = mid - 1
			} else {
				queries[qi].low = mid + 1
			}
		}
	}
	ans := make([]int, m)
	for i := range queries {
		if queries[i].low < n {
			ans[queries[i].id] = buses[queries[i].low].id
		} else {
			ans[queries[i].id] = -1
		}
	}
	var buf strings.Builder
	for i, v := range ans {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(strconv.Itoa(v))
	}
	buf.WriteByte('\n')
	return buf.String()
}

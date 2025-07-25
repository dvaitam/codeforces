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

type SegTree struct {
	n    int
	minv []int64
	maxv []int64
	lazy []int64
}

func NewSegTree(a []int64) *SegTree {
	n := len(a) - 1
	size := 4 * n
	st := &SegTree{n: n, minv: make([]int64, size), maxv: make([]int64, size), lazy: make([]int64, size)}
	var build func(node, l, r int)
	build = func(node, l, r int) {
		if l == r {
			st.minv[node] = a[l]
			st.maxv[node] = a[l]
			return
		}
		mid := (l + r) >> 1
		build(node<<1, l, mid)
		build(node<<1|1, mid+1, r)
		if st.minv[node<<1] < st.minv[node<<1|1] {
			st.minv[node] = st.minv[node<<1]
		} else {
			st.minv[node] = st.minv[node<<1|1]
		}
		if st.maxv[node<<1] > st.maxv[node<<1|1] {
			st.maxv[node] = st.maxv[node<<1]
		} else {
			st.maxv[node] = st.maxv[node<<1|1]
		}
	}
	build(1, 1, n)
	return st
}

func (st *SegTree) apply(node int, v int64) {
	st.minv[node] += v
	st.maxv[node] += v
	st.lazy[node] += v
}

func (st *SegTree) push(node int) {
	if st.lazy[node] != 0 {
		v := st.lazy[node]
		st.apply(node<<1, v)
		st.apply(node<<1|1, v)
		st.lazy[node] = 0
	}
}

func (st *SegTree) Update(node, l, r, ql, qr int, x int64) {
	if ql <= l && r <= qr {
		st.apply(node, x)
		return
	}
	st.push(node)
	mid := (l + r) >> 1
	if ql <= mid {
		st.Update(node<<1, l, mid, ql, qr, x)
	}
	if qr > mid {
		st.Update(node<<1|1, mid+1, r, ql, qr, x)
	}
	left, right := node<<1, node<<1|1
	if st.minv[left] < st.minv[right] {
		st.minv[node] = st.minv[left]
	} else {
		st.minv[node] = st.minv[right]
	}
	if st.maxv[left] > st.maxv[right] {
		st.maxv[node] = st.maxv[left]
	} else {
		st.maxv[node] = st.maxv[right]
	}
}

const INF = int(^uint(0) >> 1)

func (st *SegTree) QueryLeft(node, l, r int, y int64) int {
	if st.minv[node] > y || st.maxv[node] < y {
		return INF
	}
	if l == r {
		return l
	}
	st.push(node)
	mid := (l + r) >> 1
	res := st.QueryLeft(node<<1, l, mid, y)
	if res != INF {
		return res
	}
	return st.QueryLeft(node<<1|1, mid+1, r, y)
}

func (st *SegTree) QueryRight(node, l, r int, y int64) int {
	if st.minv[node] > y || st.maxv[node] < y {
		return -INF
	}
	if l == r {
		return l
	}
	st.push(node)
	mid := (l + r) >> 1
	res := st.QueryRight(node<<1|1, mid+1, r, y)
	if res != -INF {
		return res
	}
	return st.QueryRight(node<<1, l, mid, y)
}

func runCase(bin string, n, q int, a []int64, queries [][]int64) error {
	st := NewSegTree(append([]int64{0}, a...))
	var expect []string
	for _, qq := range queries {
		if int(qq[0]) == 1 {
			l := int(qq[1])
			r := int(qq[2])
			x := qq[3]
			st.Update(1, 1, st.n, l, r, x)
		} else {
			y := qq[1]
			if st.minv[1] > y || st.maxv[1] < y {
				expect = append(expect, "-1")
			} else {
				left := st.QueryLeft(1, 1, st.n, y)
				right := st.QueryRight(1, 1, st.n, y)
				if left > right {
					expect = append(expect, "-1")
				} else {
					expect = append(expect, strconv.Itoa(right-left))
				}
			}
		}
	}
	var input bytes.Buffer
	fmt.Fprintf(&input, "%d %d\n", n, q)
	for i := 0; i < n; i++ {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.FormatInt(a[i], 10))
	}
	input.WriteByte('\n')
	for _, qq := range queries {
		if qq[0] == 1 {
			fmt.Fprintf(&input, "1 %d %d %d\n", int(qq[1]), int(qq[2]), int(qq[3]))
		} else {
			fmt.Fprintf(&input, "2 %d\n", int(qq[1]))
		}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	outs := strings.Fields(strings.TrimSpace(string(out)))
	if len(outs) != len(expect) {
		return fmt.Errorf("expected %d answers got %d", len(expect), len(outs))
	}
	for i := 0; i < len(expect); i++ {
		if outs[i] != expect[i] {
			return fmt.Errorf("wrong answer on query %d: expected %s got %s", i+1, expect[i], outs[i])
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
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			fmt.Printf("case %d: invalid format\n", idx)
			os.Exit(1)
		}
		header := strings.Fields(strings.TrimSpace(parts[0]))
		if len(header) < 2 {
			fmt.Printf("case %d: bad header\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(header[0])
		q, _ := strconv.Atoi(header[1])
		if len(header) != 2+n {
			fmt.Printf("case %d: wrong array length\n", idx)
			os.Exit(1)
		}
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(header[2+i])
			a[i] = int64(v)
		}
		queryStrings := strings.Split(parts[1], ";")
		if len(queryStrings) != q {
			fmt.Printf("case %d: wrong query count\n", idx)
			os.Exit(1)
		}
		queries := make([][]int64, q)
		for i := 0; i < q; i++ {
			fs := strings.Fields(strings.TrimSpace(queryStrings[i]))
			if len(fs) == 0 {
				fmt.Printf("case %d: empty query\n", idx)
				os.Exit(1)
			}
			if fs[0] == "1" {
				if len(fs) != 4 {
					fmt.Printf("case %d: bad update query\n", idx)
					os.Exit(1)
				}
				l, _ := strconv.Atoi(fs[1])
				r, _ := strconv.Atoi(fs[2])
				x, _ := strconv.Atoi(fs[3])
				queries[i] = []int64{1, int64(l), int64(r), int64(x)}
			} else if fs[0] == "2" {
				if len(fs) != 2 {
					fmt.Printf("case %d: bad query2\n", idx)
					os.Exit(1)
				}
				y, _ := strconv.Atoi(fs[1])
				queries[i] = []int64{2, int64(y)}
			} else {
				fmt.Printf("case %d: unknown query type\n", idx)
				os.Exit(1)
			}
		}
		if err := runCase(bin, n, q, a, queries); err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

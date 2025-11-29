package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesRaw = `2 3 4 13|add 1 1 1|add 2 2 2|count 2 2
2 1 9 7|count 1 2
3 2 6 10 10|count 3 3|add 3 3 4
5 2 6 8 16 9 3|count 5 5|add 3 5 3
5 2 14 14 20 10 14|add 2 4 3|count 1 1
1 4 9|count 1 1|add 1 1 4|count 1 1|add 1 1 5
2 3 4 2|count 2 2|add 1 2 4|add 2 2 3
3 1 11 10 11|count 3 3
5 1 10 20 7 15 10|add 4 4 3
5 1 12 2 15 6 12|count 3 4
5 1 15 7 14 7 4|add 1 2 5
2 5 2 18|add 1 2 1|add 2 2 2|add 1 2 4|add 1 2 4|add 2 2 4
2 1 2 9|add 1 1 4
3 2 11 2 11|count 3 3|count 3 3
4 4 3 14 7 19|count 2 3|add 4 4 4|count 3 4|count 4 4
3 2 2 13 20|add 2 2 2|count 2 3
4 4 7 1 7 6|add 3 3 4|count 4 4|count 2 2|count 3 4
5 4 1 3 2 20 4|add 3 5 2|add 1 5 1|add 3 3 1|count 4 5
2 3 13 8|count 2 2|add 2 2 4|count 2 2
2 3 16 14|add 1 2 2|add 1 2 4|count 2 2
2 1 10 4|count 2 2
5 1 17 8 14 10 12|add 2 2 1
5 3 18 15 19 10 17|count 5 5|add 3 5 3|count 3 3
4 1 20 5 20 6|add 3 3 5
3 5 3 3 13|count 1 2|count 2 2|add 1 3 1|count 1 2|count 1 1
2 5 16 19|count 1 1|add 2 2 5|count 2 2|count 1 2|count 2 2
2 3 18 17|count 2 2|count 2 2|add 1 2 2
5 4 18 20 12 13 13|count 2 2|add 3 5 3|count 3 4|add 4 4 2
2 1 10 5|count 1 2
5 5 18 18 9 1 17|add 1 2 1|count 1 1|add 3 5 2|count 4 4|count 1 5
3 1 18 5 17|add 1 2 4
5 1 13 6 6 9 17|add 3 4 3
2 4 2 14|count 2 2|add 1 1 2|add 1 1 2|count 2 2
4 4 11 6 17 19|add 4 4 3|count 4 4|count 2 3|count 3 4
3 1 8 11 4|count 3 3
5 3 17 16 12 3 6|count 4 4|add 1 3 5|count 3 3
5 5 1 3 8 20 9|count 1 1|add 2 3 3|add 3 4 4|count 4 4|add 4 4 5
5 3 18 10 7 14 11|add 1 4 4|count 3 4|count 4 5
3 5 7 9 9|count 1 2|add 2 2 5|count 1 1|count 3 3|count 2 2
4 1 11 1 15 19|count 2 2
2 5 4 4|count 2 2|count 1 1|count 1 1|count 1 1|add 2 2 5
3 2 6 7 2|count 1 3|count 1 3
3 4 4 7 5|add 3 3 5|add 2 3 1|count 1 1|add 1 2 2
3 5 7 20 12|add 3 3 2|add 2 2 5|add 1 1 2|add 3 3 5|add 1 2 4
2 2 13 20|add 2 2 4|count 2 2
4 2 4 11 10 20|count 1 4|add 1 2 5
2 5 11 20|add 1 2 4|add 1 1 2|add 2 2 4|count 2 2|count 2 2
5 3 18 2 18 5 12|count 2 5|add 1 1 5|add 2 3 5
2 5 2 18|count 1 1|count 1 1|add 1 2 4|count 2 2|count 2 2
2 3 14 2|add 2 2 2|add 1 2 1|add 1 1 2
1 1 1|count 1 1
1 1 17|count 1 1
5 5 2 2 3 5 14|count 1 5|add 3 5 3|add 5 5 5|count 1 5|count 2 2
3 5 14 10 19|add 1 3 3|count 1 3|add 2 2 5|count 1 2|count 3 3
3 2 10 6 11|count 2 2|count 3 3
3 3 6 11 4|count 2 2|add 1 2 1|count 1 1
5 2 8 8 12 13 13|add 4 5 3|count 3 3
1 2 2|add 1 1 4|count 1 1
3 3 17 12 5|add 2 2 4|count 2 3|count 2 2
4 2 3 15 13 17|add 1 1 3|add 2 3 4
3 4 11 13 10|count 2 2|count 2 3|count 1 3|add 3 3 2
5 5 6 14 10 4 11|count 5 5|count 4 4|count 3 4|add 5 5 4|count 1 4
3 4 9 7 5|add 1 2 5|count 3 3|count 1 2|count 3 3
3 5 17 2 14|add 3 3 5|add 1 2 5|add 1 1 2|count 1 2|count 2 2
3 1 3 2 15|add 3 3 2
4 5 3 18 20 6|add 1 2 5|count 2 2|count 3 4|add 4 4 3|count 1 3
3 1 2 11 2|count 2 2
2 2 2 16|count 2 2|count 1 1
4 4 17 13 4 6|add 2 4 4|count 4 4|add 4 4 3|add 1 3 4
5 3 10 4 6 17 7|count 2 5|count 2 3|count 2 2
5 1 10 1 3 17 15|count 5 5
3 3 7 7 4|count 3 3|add 2 2 1|add 3 3 4
5 3 7 16 8 6 6|add 1 2 2|count 4 5|add 2 4 5
5 4 15 12 11 3 15|add 3 3 1|count 3 5|add 5 5 4|add 3 5 5
1 5 10|count 1 1|count 1 1|add 1 1 2|count 1 1|add 1 1 4
3 3 14 7 2|count 3 3|count 3 3|count 1 3
4 3 9 2 5 12|add 4 4 3|add 4 4 4|count 4 4
2 4 14 9|count 2 2|add 1 1 2|add 1 1 4|count 1 2
5 2 1 14 14 9 16|add 1 1 2|count 3 5
4 1 11 9 18 15|add 1 3 3
3 3 13 11 2|count 1 2|count 3 3|add 1 1 4
5 3 14 4 13 14 15|count 4 4|count 4 4|add 3 4 1
4 5 20 17 7 10|count 4 4|count 3 3|add 4 4 5|count 3 3|add 3 3 3
4 4 10 14 17 6|count 4 4|count 3 4|count 3 4|count 2 3
1 4 19|add 1 1 2|add 1 1 2|add 1 1 3|add 1 1 1
3 5 18 4 18|count 3 3|count 3 3|add 3 3 4|count 1 3|count 2 2
2 2 15 7|count 1 1|add 1 2 2
5 2 11 14 4 2 12|count 3 3|add 3 5 2
2 5 3 17|count 2 2|add 1 1 2|count 1 1|count 2 2|count 2 2
2 1 17 11|count 2 2
1 5 1|count 1 1|count 1 1|count 1 1|count 1 1|add 1 1 3
5 1 17 7 19 6 13|count 5 5
4 1 9 5 14 7|count 3 3
2 4 9 8|count 2 2|add 1 1 4|count 2 2|add 1 1 2
2 1 9 5|count 2 2
3 1 19 3 14|add 2 3 5
4 2 3 13 15 15|count 4 4|add 3 3 3
2 4 8 18|add 2 2 4|add 2 2 2|count 1 2|add 1 1 5
5 1 20 12 15 13 17|count 5 5
2 2 3 7|add 1 1 1|add 2 2 1`

// Embedded reference solution logic from 121E.go.
const INF = 1000000000

type SegTree struct {
	n     int
	arr   []int
	lucky []int
	mn    []int
	sum   []int
	lazy  []int
}

func NewSegTree(n int, arr []int, lucky []int) *SegTree {
	size := 4 * n
	st := &SegTree{
		n:     n,
		arr:   make([]int, n),
		lucky: lucky,
		mn:    make([]int, size),
		sum:   make([]int, size),
		lazy:  make([]int, size),
	}
	copy(st.arr, arr)
	st.build(1, 0, n-1)
	return st
}

func (st *SegTree) build(node, l, r int) {
	if l == r {
		d := st.nextDist(st.arr[l])
		st.mn[node] = d
		if d == 0 {
			st.sum[node] = 1
		}
		return
	}
	mid := (l + r) >> 1
	st.build(node<<1, l, mid)
	st.build(node<<1|1, mid+1, r)
	st.pull(node)
}

func (st *SegTree) pull(node int) {
	l, r := node<<1, node<<1|1
	if st.mn[l] < st.mn[r] {
		st.mn[node] = st.mn[l]
	} else {
		st.mn[node] = st.mn[r]
	}
	st.sum[node] = st.sum[l] + st.sum[r]
}

func (st *SegTree) push(node int) {
	if st.lazy[node] != 0 {
		d := st.lazy[node]
		for _, c := range []int{node << 1, node<<1 | 1} {
			st.lazy[c] += d
			st.mn[c] -= d
		}
		st.lazy[node] = 0
	}
}

func (st *SegTree) nextDist(x int) int {
	idx := sort.Search(len(st.lucky), func(i int) bool { return st.lucky[i] >= x })
	if idx < len(st.lucky) {
		return st.lucky[idx] - x
	}
	return INF
}

func (st *SegTree) Update(node, l, r, ql, qr, d int) {
	if r < ql || l > qr {
		return
	}
	if ql <= l && r <= qr && st.mn[node] > d {
		st.mn[node] -= d
		st.lazy[node] += d
		return
	}
	if l == r {
		st.arr[l] += st.lazy[node]
		st.lazy[node] = 0
		st.arr[l] += d
		dd := st.nextDist(st.arr[l])
		st.mn[node] = dd
		if dd == 0 {
			st.sum[node] = 1
		} else {
			st.sum[node] = 0
		}
		return
	}
	st.push(node)
	mid := (l + r) >> 1
	st.Update(node<<1, l, mid, ql, qr, d)
	st.Update(node<<1|1, mid+1, r, ql, qr, d)
	st.pull(node)
}

func (st *SegTree) Query(node, l, r, ql, qr int) int {
	if r < ql || l > qr {
		return 0
	}
	if ql <= l && r <= qr {
		return st.sum[node]
	}
	st.push(node)
	mid := (l + r) >> 1
	return st.Query(node<<1, l, mid, ql, qr) + st.Query(node<<1|1, mid+1, r, ql, qr)
}

func genLuck(cur, maxv int, out *[]int) {
	if cur > maxv {
		return
	}
	if cur > 0 {
		*out = append(*out, cur)
	}
	genLuck(cur*10+4, maxv, out)
	genLuck(cur*10+7, maxv, out)
}

func referenceSolve(n, m int, arr []int, ops []string) []int {
	var lucky []int
	genLuck(0, 10000, &lucky)
	sort.Ints(lucky)
	st := NewSegTree(n, arr, lucky)
	res := make([]int, 0)
	for _, op := range ops {
		parts := strings.Fields(op)
		if len(parts) == 0 {
			continue
		}
		if parts[0] == "add" {
			l, _ := strconv.Atoi(parts[1])
			r, _ := strconv.Atoi(parts[2])
			d, _ := strconv.Atoi(parts[3])
			st.Update(1, 0, n-1, l-1, r-1, d)
		} else { // count
			l, _ := strconv.Atoi(parts[1])
			r, _ := strconv.Atoi(parts[2])
			res = append(res, st.Query(1, 0, n-1, l-1, r-1))
		}
	}
	return res
}

func parseLine(line string) (int, int, []int, []string, error) {
	parts := strings.Split(line, "|")
	if len(parts) < 1 {
		return 0, 0, nil, nil, fmt.Errorf("invalid line")
	}
	header := strings.Fields(parts[0])
	if len(header) < 2 {
		return 0, 0, nil, nil, fmt.Errorf("invalid header")
	}
	n, err := strconv.Atoi(header[0])
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("invalid n")
	}
	m, err := strconv.Atoi(header[1])
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("invalid m")
	}
	if len(header) != 2+n {
		return 0, 0, nil, nil, fmt.Errorf("expected %d array values got %d", n, len(header)-2)
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i], _ = strconv.Atoi(header[2+i])
	}
	if len(parts)-1 != m {
		return 0, 0, nil, nil, fmt.Errorf("expected %d operations got %d", m, len(parts)-1)
	}
	ops := make([]string, m)
	for i := 0; i < m; i++ {
		ops[i] = strings.TrimSpace(parts[1+i])
	}
	return n, m, arr, ops, nil
}

func runCase(bin string, line string) error {
	n, m, arr, ops, err := parseLine(line)
	if err != nil {
		return err
	}
	expect := referenceSolve(n, m, arr, ops)

	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", n, m)
	for i, v := range arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')
	for _, op := range ops {
		input.WriteString(op)
		input.WriteByte('\n')
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotFields := strings.Fields(strings.TrimSpace(out.String()))
	if len(gotFields) != len(expect) {
		return fmt.Errorf("expected %d outputs got %d", len(expect), len(gotFields))
	}
	for i, g := range gotFields {
		val, err := strconv.Atoi(g)
		if err != nil {
			return fmt.Errorf("invalid output %q", g)
		}
		if val != expect[i] {
			return fmt.Errorf("query %d: expected %d got %d", i+1, expect[i], val)
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

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		if err := runCase(bin, line); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

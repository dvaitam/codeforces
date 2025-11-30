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

const testcasesE = `4 -9 6 -8 6 -2 4 10 4 7 2 7 5 -7 10 10 10
2 10 -3 10 -2 9 -6 10 -6
2 -7 -8 2 -8 9 -2 10 -2
5 8 3 8 4 4 -6 5 -6 0 -2 0 9 2 -10 2 0 -8 -4 -3 -4
2 0 4 1 4 3 8 3 9
3 -4 -1 -4 1 10 -2 10 2 -1 -8 -1 4
2 -10 2 -10 10 -1 10 8 10
5 4 9 4 10 -2 -10 0 -10 -2 -4 -2 4 6 2 6 5 -5 -10 -5 -3
3 -7 -10 -7 4 0 9 0 10 3 -6 3 -2
2 5 4 5 6 3 3 3 7
3 7 -8 9 -8 -10 -7 7 -7 10 5 10 7
4 -9 3 1 3 -2 -1 7 -1 0 7 1 7 -8 -10 3 -10
3 -4 -3 -4 3 -2 7 -2 8 2 -10 2 -6
1 -4 0 1 0
5 -1 6 1 6 9 7 10 7 -7 -3 -1 -3 -5 -5 -5 3 4 8 4 9
5 -10 -7 -10 -3 -8 2 -8 6 2 -3 2 5 -5 3 -5 9 3 8 8 8
2 1 6 1 9 10 4 10 10
3 -2 -5 -2 -2 0 4 1 4 5 9 5 10
2 4 -7 10 -7 2 -4 8 -4
5 5 -9 9 -9 4 1 4 2 -4 8 -2 8 -6 -8 -6 9 9 -2 9 2
1 -10 7 -10 9
3 9 -4 10 -4 4 6 4 9 -8 -9 4 -9
4 2 -3 2 2 7 -10 7 4 5 -6 5 7 8 9 8 10
4 -1 3 -1 10 -9 1 5 1 4 -1 4 1 -3 0 4 0
1 -6 -5 -6 9
4 -5 -7 5 -7 -5 -3 -4 -3 3 -5 3 0 4 -1 4 8
1 -5 5 10 5
1 -1 5 4 5
3 -9 8 -9 10 4 -8 4 -6 -3 4 -3 10
3 -8 2 6 2 -8 -9 4 -9 -1 -1 0 -1
1 -9 9 -4 9
5 -8 5 -3 5 -6 7 -2 7 9 -8 10 -8 5 -7 9 -7 6 4 8 4
5 1 -6 8 -6 6 -8 7 -8 2 -7 8 -7 7 8 7 10 -3 8 -3 9
2 1 -1 1 5 -10 -10 -10 -9
2 -4 -6 -4 0 5 5 7 5
3 -6 -9 -6 -3 0 -4 0 7 5 6 7 6
3 -4 -7 -4 1 2 -7 2 2 7 1 7 2
3 -9 2 5 2 -5 -2 -5 -1 -8 0 -8 7
1 -3 -5 -3 0
3 5 -5 7 -5 -1 3 10 3 -10 -2 -9 -2
5 -5 0 10 0 -3 8 -3 9 7 10 8 10 7 4 7 9 1 1 1 3
1 1 -5 1 3
2 9 -10 9 -3 3 2 3 5
2 2 4 2 10 -8 -2 -2 -2
5 -5 -1 -3 -1 -5 -2 2 -2 6 3 6 6 5 4 5 5 9 0 9 5
3 2 9 2 10 10 6 10 8 -2 3 -2 9
1 7 3 8 3
3 5 -10 9 -10 7 5 7 10 5 8 5 9
5 1 9 5 9 3 -8 3 0 9 -5 10 -5 0 5 8 5 5 -2 8 -2
5 1 -9 10 -9 2 7 2 10 -6 -1 -1 -1 0 -8 4 -8 5 -10 9 -10
1 4 -6 9 -6
2 -8 -1 -8 5 1 4 2 4
3 9 -7 10 -7 -1 8 -1 10 5 -10 9 -10
4 3 -8 10 -8 5 3 5 9 2 -6 2 8 -4 3 1 3
1 0 4 6 4
3 0 -3 0 0 -6 0 -6 4 -7 7 -7 10
5 0 -5 0 -3 -5 -2 -5 4 2 -7 8 -7 -1 0 -1 1 -1 10 10 10
2 -5 7 3 7 6 -8 9 -8
1 -8 2 5 2
4 -6 -6 -6 -4 -7 -6 6 -6 -8 -1 -4 -1 -5 2 4 2
4 0 -5 0 6 -2 -5 -2 10 -4 -8 7 -8 4 -8 4 3
5 -1 4 -1 5 -5 1 3 1 -8 -3 -3 -3 3 -7 3 5 0 -7 3 -7
2 7 9 9 9 0 8 0 10
5 -1 -10 -1 10 -4 -2 5 -2 4 -1 6 -1 4 -9 9 -9 3 0 3 10
3 5 0 7 0 -6 5 -5 5 -10 7 -5 7
1 -1 0 10 0
2 -3 -9 9 -9 2 1 8 1
4 -9 2 -1 2 -10 4 4 4 -8 -4 -3 -4 -2 -1 -2 7
3 5 0 5 5 4 7 4 9 0 2 0 7
1 -1 4 6 4
4 5 -5 7 -5 -4 9 -2 9 6 6 6 8 -7 5 9 5
2 -5 -7 -5 -5 -7 -8 -2 -8
2 1 10 7 10 2 -7 10 -7
2 8 6 8 8 8 9 10 9
5 -7 -6 -7 8 -6 -4 0 -4 -2 0 -2 3 0 7 0 8 -10 2 -8 2
2 3 -3 3 10 8 -7 8 4
1 9 5 9 8
2 -9 1 -9 7 7 -8 9 -8
4 -5 9 9 9 8 7 10 7 -8 2 -8 7 -1 -10 -1 3
1 -7 -6 -7 4
3 9 -6 9 10 -10 9 -10 10 -10 3 0 3
2 3 -8 10 -8 7 -6 7 8
3 1 7 4 7 -3 -8 -1 -8 -4 5 -3 5
4 -2 1 -2 6 6 -9 6 1 3 -3 3 10 -10 6 -10 10
3 1 8 6 8 -2 -5 2 -5 2 1 5 1
3 9 -5 9 9 -10 -8 -10 -2 3 -8 3 -5
5 -8 7 -1 7 4 7 4 10 -5 -5 8 -5 -2 4 2 4 -8 5 -7 5
1 5 -5 7 -5
5 8 8 8 9 7 1 8 1 10 -1 10 8 -1 8 0 8 1 -8 1 3
3 -7 -1 -7 9 -9 1 -8 1 -10 3 -10 5
4 -1 -10 -1 -4 7 -8 7 -7 -6 -8 7 -8 9 -5 9 7
2 -9 -10 4 -10 -1 -5 -1 9
4 3 0 8 0 1 6 5 6 5 3 10 3 0 -4 8 -4
2 -9 -4 -2 -4 0 -1 6 -1
3 8 -5 10 -5 -9 2 -9 4 -5 2 -5 6
4 -6 7 -6 10 9 1 9 2 6 -4 10 -4 -3 8 -3 9
5 -5 -6 0 -6 -7 0 -2 0 3 7 6 7 -6 1 -4 1 9 4 10 4
4 8 -3 10 -3 -7 0 2 0 1 -1 1 0 -7 1 -7 6
2 -3 -3 -3 1 -7 2 2 2
2 2 -10 2 8 3 7 9 7
1 -4 6 -4 9
1 3 6 9 6
5 -5 -6 -5 -2 -1 -5 -1 3 -3 7 -3 10 5 -8 7 -8 3 -2 3 4
1 1 4 10 4
2 2 4 10 4 9 -1 9 7
2 -9 3 -9 9 8 4 9 4
1 -3 2 2 2
1 10 8 10 10
5 -9 8 4 8 6 1 6 5 -7 -5 -3 -5 1 -1 9 -1 -6 3 0 3
3 0 -10 3 -10 0 6 10 6 7 -7 8 -7
4 -1 0 -1 10 -9 -4 4 -4 6 0 6 4 -2 -8 -2 7
2 4 -7 10 -7 0 -2 0 6
1 0 -3 0 3
3 2 -4 2 -3 10 8 10 10 9 -10 9 8
1 -10 3 -10 5
1 7 -9 9 -9
2 -3 8 -3 9 0 3 5 3
1 -2 4 9 4
4 4 -5 5 -5 -5 -10 5 -10 1 -8 2 -8 5 -3 8 -3
3 10 7 10 9 9 6 9 7 -4 -6 0 -6`

type dsu struct {
	p   []int
	sum []int64
}

func newDSU(n int, lengths []int64) *dsu {
	p := make([]int, n)
	sum := make([]int64, n)
	for i := 0; i < n; i++ {
		p[i] = i
		sum[i] = lengths[i]
	}
	return &dsu{p: p, sum: sum}
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	d.p[rb] = ra
	d.sum[ra] += d.sum[rb]
}

type segTree struct {
	n  int
	id []int
}

func newSegTree(n int) *segTree {
	size := 1
	for size < n {
		size <<= 1
	}
	id := make([]int, size*2)
	for i := range id {
		id[i] = -1
	}
	return &segTree{n: size, id: id}
}

func (st *segTree) update(i int, v int) {
	i += st.n
	st.id[i] = v
	for i >>= 1; i > 0; i >>= 1 {
		if st.id[i<<1] != -1 {
			st.id[i] = st.id[i<<1]
		} else {
			st.id[i] = st.id[i<<1|1]
		}
	}
}

func (st *segTree) collect(u, l, r, ql, qr int, res *[]int) {
	if ql > r || qr < l || st.id[u] == -1 {
		return
	}
	if l == r {
		*res = append(*res, st.id[u])
		return
	}
	mid := (l + r) >> 1
	st.collect(u<<1, l, mid, ql, qr, res)
	st.collect(u<<1|1, mid+1, r, ql, qr, res)
}

func solve(n int, coords []int) int64 {
	type hseg struct {
		x1, x2, y int
		id        int
	}
	type vseg struct {
		x, y1, y2 int
		id        int
	}
	hs := make([]hseg, 0, n)
	vs := make([]vseg, 0, n)
	lengths := make([]int64, n)
	for i := 0; i < n; i++ {
		x1, y1, x2, y2 := coords[4*i], coords[4*i+1], coords[4*i+2], coords[4*i+3]
		lengths[i] = int64(x2-x1) + int64(y2-y1)
		if y1 == y2 {
			hs = append(hs, hseg{x1, x2, y1, i})
		} else {
			vs = append(vs, vseg{x1, y1, y2, i})
		}
	}
	d := newDSU(n, lengths)
	ys := make([]int, len(hs))
	for i, h := range hs {
		ys[i] = h.y
	}
	sort.Ints(ys)
	ys = uniqueInts(ys)
	yIndex := make(map[int]int, len(ys))
	for i, y := range ys {
		yIndex[y] = i
	}
	const (
		evAdd   = 0
		evQuery = 1
		evRem   = 2
	)
	type event struct {
		x, typ, id int
		yi, yj     int
	}
	evs := make([]event, 0, len(hs)*2+len(vs))
	for _, h := range hs {
		yi := yIndex[h.y]
		evs = append(evs, event{h.x1, evAdd, h.id, yi, 0})
		evs = append(evs, event{h.x2, evRem, h.id, yi, 0})
	}
	for _, v := range vs {
		yl := sort.SearchInts(ys, v.y1)
		yr := sort.Search(len(ys), func(i int) bool { return ys[i] > v.y2 }) - 1
		if yl <= yr {
			evs = append(evs, event{v.x, evQuery, v.id, yl, yr})
		}
	}
	sort.Slice(evs, func(i, j int) bool {
		if evs[i].x != evs[j].x {
			return evs[i].x < evs[j].x
		}
		return evs[i].typ < evs[j].typ
	})
	st := newSegTree(len(ys))
	for _, e := range evs {
		switch e.typ {
		case evAdd:
			st.update(e.yi, e.id)
		case evRem:
			st.update(e.yi, -1)
		case evQuery:
			var ids []int
			st.collect(1, 0, st.n-1, e.yi, e.yj, &ids)
			for _, hid := range ids {
				d.union(e.id, hid)
			}
		}
	}
	var ans int64
	for i := 0; i < n; i++ {
		if d.p[i] == i && d.sum[i] > ans {
			ans = d.sum[i]
		}
	}
	return ans
}

func uniqueInts(a []int) []int {
	j := 0
	for i := 0; i < len(a); i++ {
		if i == 0 || a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseLine(line string) (int, []int, error) {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return 0, nil, fmt.Errorf("empty line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, err
	}
	if len(fields) != 1+4*n {
		return 0, nil, fmt.Errorf("expected %d numbers got %d", 1+4*n, len(fields))
	}
	vals := make([]int, 4*n)
	for i := 0; i < 4*n; i++ {
		v, convErr := strconv.Atoi(fields[1+i])
		if convErr != nil {
			return 0, nil, convErr
		}
		vals[i] = v
	}
	return n, vals, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesE))
	scanner.Buffer(make([]byte, 0, 1024), 1<<20)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, vals, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", idx, err)
			os.Exit(1)
		}
		expected := strconv.FormatInt(solve(n, vals), 10)

		var input strings.Builder
		input.WriteString(strconv.Itoa(n))
		input.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(vals[4*i]))
			input.WriteByte(' ')
			input.WriteString(strconv.Itoa(vals[4*i+1]))
			input.WriteByte(' ')
			input.WriteString(strconv.Itoa(vals[4*i+2]))
			input.WriteByte(' ')
			input.WriteString(strconv.Itoa(vals[4*i+3]))
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

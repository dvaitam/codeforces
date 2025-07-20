package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const inf = 2_000_000_005

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Mat struct{ a [2][2]int }

func newMat() Mat {
	var m Mat
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			m.a[i][j] = inf
		}
	}
	return m
}

var ident = Mat{a: [2][2]int{{0, inf}, {inf, 0}}}

func mul(u, v Mat) Mat {
	w := newMat()
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			w.a[i][j] = min(max(u.a[i][0], v.a[0][j]), max(u.a[i][1], v.a[1][j]))
		}
	}
	return w
}

type SegTree struct {
	N  int
	tr []Mat
}

func NewSegTree(n int, leaf []Mat) *SegTree {
	N := 1
	for N <= n {
		N <<= 1
	}
	st := &SegTree{N: N, tr: make([]Mat, 2*N+2)}
	for i := 1; i < 2*N; i++ {
		st.tr[i] = ident
	}
	for i := 1; i <= n; i++ {
		st.tr[i+N] = leaf[i]
	}
	for i := N - 1; i >= 1; i-- {
		st.tr[i] = mul(st.tr[i<<1], st.tr[i<<1|1])
	}
	return st
}

func (st *SegTree) Update(pos int, val Mat) {
	idx := pos + st.N
	st.tr[idx] = val
	for idx >>= 1; idx > 0; idx >>= 1 {
		st.tr[idx] = mul(st.tr[idx<<1], st.tr[idx<<1|1])
	}
}

type Op struct{ val, idx, x, y int }

func expectedE2(n int, a [][2]int) int {
	f := make([]Mat, n+2)
	for i := 1; i <= n; i++ {
		f[i] = newMat()
	}
	var ops []Op
	if n%2 == 0 {
		for i := 1; i <= n; i++ {
			if i%2 == 0 {
				f[i].a = [2][2]int{{0, 0}, {0, 0}}
				continue
			}
			for x := 0; x < 2; x++ {
				for y := 0; y < 2; y++ {
					l := a[i][x] + a[i+1][y]
					r := a[i][x^1] + a[i+1][y^1]
					f[i].a[x][y] = max(l, r)
					ops = append(ops, Op{min(l, r), i, x, y})
				}
			}
		}
	} else {
		a = append(a, [2]int{})
		a[n+1][0], a[n+1][1] = a[1][1], a[1][0]
		for i := 1; i <= n; i++ {
			for x := 0; x < 2; x++ {
				for y := 0; y < 2; y++ {
					if i&1 == 1 {
						f[i].a[x][y] = a[i][x] + a[i+1][y]
					} else {
						f[i].a[x][y] = a[i][x^1] + a[i+1][y^1]
					}
					ops = append(ops, Op{f[i].a[x][y], i, x, y})
				}
			}
		}
	}
	st := NewSegTree(n, f)
	ans := inf
	sort.Slice(ops, func(i, j int) bool { return ops[i].val < ops[j].val })
	for _, op := range ops {
		root := st.tr[1]
		v := min(root.a[0][0], root.a[1][1])
		if v == inf {
			break
		}
		if cur := v - op.val; cur < ans {
			ans = cur
		}
		f[op.idx].a[op.x][op.y] = inf
		st.Update(op.idx, f[op.idx])
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE2.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		var n int
		fmt.Sscan(fields[0], &n)
		if len(fields) != 1+2*n {
			fmt.Printf("test %d: invalid number of values\n", idx)
			os.Exit(1)
		}
		a := make([][2]int, n+2)
		for i := 1; i <= n; i++ {
			fmt.Sscan(fields[i], &a[i][0])
		}
		for i := 1; i <= n; i++ {
			fmt.Sscan(fields[n+i], &a[i][1])
		}
		expect := expectedE2(n, a)
		input := fmt.Sprintf("1\n%d\n", n)
		for i := 1; i <= n; i++ {
			input += fmt.Sprintf("%d ", a[i][0])
		}
		input += "\n"
		for i := 1; i <= n; i++ {
			input += fmt.Sprintf("%d ", a[i][1])
		}
		input += "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		res := strings.TrimSpace(out.String())
		var got int
		if _, err := fmt.Sscan(res, &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx, res)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

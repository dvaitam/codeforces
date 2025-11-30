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

const (
	inf = 2_000_000_005
	// Embedded copy of testcasesE2.txt to avoid filesystem dependency.
	testcasesE2Data = `
3 5 4 4 4 1 3
3 3 3 3 4 5 5
1 4 5
2 4 3 4 1
1 4 5
1 3 1
1 4 4
3 5 4 5 4 5 4
2 5 1 5 3
2 3 4 1 2
1 2 5
4 2 4 3 4 2 4 3 4
3 5 2 1 3 5 3
4 5 5 1 3 2 5 2 2
1 3 2
1 1 4
4 1 5 5 1 2 4 5 2
3 1 3 3 2 4 3
1 3 1
1 5 1
3 3 5 3 2 1 4
4 4 2 1 4 3 4 4 4
4 2 4 2 4 2 3 1 2
4 3 3 3 4 3 2 2 4
2 3 2 4 2
1 5 5
4 3 2 1 3 2 3 3 1
1 4 3
3 3 3 1 2 3 2
2 5 2 4 2
1 5 3
1 1 4
1 4 4
4 2 2 1 1 2 1 4 5
1 3 1
1 2 3
4 2 5 4 3 3 2 4 3
2 2 1 2 5
3 4 1 3 5 4 3
4 1 2 3 2 4 2 2 1
2 4 4 1 3
2 1 3 2 2
2 5 5 3 4
3 5 4 5 3 3 1
2 3 2 5 5
1 3 4
2 1 1 4 5
4 4 3 5 3 2 3 4 2
1 4 2
1 1 5
1 2 1
1 3 3
2 1 1 2 2
3 4 5 5 3 4 2
4 1 5 2 5 1 4 1 4
2 2 1 5 2
3 3 3 4 3 5 4
4 3 3 5 3 4 5 3 4
3 4 3 1 2 4 3
4 4 2 2 4 2 2 5 2
1 3 5
2 5 2 4 1
4 2 4 2 1 1 5 1 1
1 5 2
2 5 2 3 3
1 1 4
4 2 3 5 5 2 5 2 5
4 5 3 1 2 4 1 2 2
2 3 3 3 2
4 3 5 1 1 4 1 1 2
3 4 2 3 4 5 1
3 1 4 5 4 3 2
4 1 3 3 2 1 5 3 4
1 5 2
2 5 5 3 4
1 4 3
4 3 1 4 1 4 1 2 3
3 3 5 2 4 4 1
3 2 4 1 3 2 5
2 1 4 5 4
2 4 3 2 5
3 5 2 1 1 2 1
2 4 2 1 4
3 3 3 5 3 4 5
1 3 4
4 2 4 1 5 5 4 3 2
3 3 2 2 3 2 1
3 4 2 3 3 5 2
4 1 1 2 4 4 1 1 5
2 5 1 2 3
4 5 5 1 2 2 5 5 5
1 4 3
2 5 1 4 1
1 1 5
2 2 4 2 5
3 5 5 1 2 3 5
1 5 3
3 5 2 5 3 1 5
4 3 1 3 3 3 1 3 5
2 3 1 1 5
`
)

type testCase struct {
	n int
	a [][2]int
}

func parseTestcasesE2() ([]testCase, error) {
	lines := strings.Split(testcasesE2Data, "\n")
	cases := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			return nil, fmt.Errorf("line %d: empty testcase", i+1)
		}
		var n int
		if _, err := fmt.Sscan(fields[0], &n); err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", i+1, err)
		}
		expectedFields := 1 + 2*n
		if len(fields) != expectedFields {
			return nil, fmt.Errorf("line %d: expected %d values got %d", i+1, expectedFields, len(fields))
		}
		a := make([][2]int, n+2)
		idx := 1
		for j := 1; j <= n; j++ {
			if _, err := fmt.Sscan(fields[idx], &a[j][0]); err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d][0]: %w", i+1, j, err)
			}
			idx++
		}
		for j := 1; j <= n; j++ {
			if _, err := fmt.Sscan(fields[idx], &a[j][1]); err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d][1]: %w", i+1, j, err)
			}
			idx++
		}
		cases = append(cases, testCase{n: n, a: a})
	}
	return cases, nil
}

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

// Embedded solution logic from 1995E2.go, refactored to return the answer for a single test.
func solveCase1995E2(n int, a [][2]int) int {
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
					ops = append(ops, Op{val: min(l, r), idx: i, x: x, y: y})
				}
			}
		}
	} else {
		a[n+1][0], a[n+1][1] = a[1][1], a[1][0]
		for i := 1; i <= n; i++ {
			for x := 0; x < 2; x++ {
				for y := 0; y < 2; y++ {
					if i&1 == 1 {
						f[i].a[x][y] = a[i][x] + a[i+1][y]
					} else {
						f[i].a[x][y] = a[i][x^1] + a[i+1][y^1]
					}
					ops = append(ops, Op{val: f[i].a[x][y], idx: i, x: x, y: y})
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

// solve mirrors the solve() function from the embedded reference solution.
// It reads a single test case from the provided reader and writes the answer to the writer.
func solve(in *bufio.Reader, out *bufio.Writer) {
	var n int
	fmt.Fscan(in, &n)

	a := make([][2]int, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i][0])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i][1])
	}

	fmt.Fprintln(out, solveCase1995E2(n, a))
}

func expectedE2(n int, a [][2]int) int {
	cloned := make([][2]int, len(a))
	copy(cloned, a)
	return solveCase1995E2(n, cloned)
}

func buildInput(n int, a [][2]int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", n)
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&sb, "%d ", a[i][0])
	}
	sb.WriteString("\n")
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&sb, "%d ", a[i][1])
	}
	sb.WriteString("\n")
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcasesE2()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		expect := expectedE2(tc.n, tc.a)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(buildInput(tc.n, tc.a))
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		res := strings.TrimSpace(out.String())
		var got int
		if _, err := fmt.Sscan(res, &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx+1, res)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	C int64
	w []int64
}

// Embedded testcases from testcasesF.txt.
const testcaseData = `
4 2 0 1
4 9 1 6
5 17 0 1 12
6 16 11 7 7 5
3 3 1
5 11 7 8 3
6 4 4 4 1 2
4 0 0 0
3 16 10
4 3 0 3
6 9 5 9 1 3
3 20 12
5 16 10 7 10
5 2 1 0 0
3 14 14
5 10 4 7 0
5 10 10 3 5
5 11 11 6 3
3 12 10
5 12 12 9 0
3 10 0
6 4 0 1 0 4
5 6 6 5 2
4 18 2 7
6 10 10 10 10 4
3 2 1
5 20 12 6 5
6 0 0 0 0 0
3 11 0
4 7 7 4
5 18 11 17 10
6 19 16 1 7 4
4 7 1 2
3 1 0
6 14 10 7 6 0
5 16 6 4 9
5 12 9 8 12
3 14 9
4 12 4 8
5 3 2 2 2
6 14 1 14 9 13
5 10 1 6 5
3 3 2
3 17 6
3 0 0
5 14 1 8 10
3 15 7
5 16 9 8 6
6 19 14 15 9 10
3 18 9
3 10 1
4 12 7 4
5 8 3 2 8
6 14 4 7 1 1
4 17 1 8
5 10 5 0 6
5 5 2 3 0
5 10 4 9 10
5 2 1 0 1
4 11 9 4
5 15 0 14 5
4 11 3 1
3 12 7
6 20 19 1 10 0
3 14 4
4 5 3 0
6 17 1 8 3 4
3 0 0
6 16 9 2 10 15
5 8 5 3 0
4 14 8 11
3 12 5
3 16 6
3 17 7
3 17 10
4 7 7 7
6 5 3 2 1 0
5 9 5 4 3
6 0 0 0 0 0
4 0 0 0
4 2 0 0
4 13 6 0
5 9 5 9 0
4 18 0 11
5 10 1 9 10
5 11 1 6 3
5 18 8 4 6
4 19 4 16
3 10 9
3 19 5
5 8 7 4 7
4 12 6 5
3 10 2
5 5 2 3 0
5 11 4 10 5
4 13 3 3
4 18 14 13
4 0 0 0
3 14 6
3 1 1
`

func parseTestcases() ([]testCase, error) {
	data := strings.TrimSpace(testcaseData)
	if data == "" {
		return nil, fmt.Errorf("no test data")
	}
	lines := strings.Split(data, "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("case %d missing n/C", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		C, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d bad C: %v", i+1, err)
		}
		expect := 2 + max(0, n-2)
		if len(fields) != expect {
			return nil, fmt.Errorf("case %d expected %d values, got %d", i+1, expect, len(fields))
		}
		w := make([]int64, max(0, n-2))
		for j := 0; j < len(w); j++ {
			v, err := strconv.ParseInt(fields[2+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d bad w[%d]: %v", i+1, j, err)
			}
			w[j] = v
		}
		res = append(res, testCase{n: n, C: C, w: w})
	}
	return res, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Embedded solution of 1500F.
type node struct {
	l, r int64
	i    int
}

type reverseDeque struct {
	org int64
	fw  bool
	d   []node
}

func newReverseDeque() *reverseDeque {
	return &reverseDeque{org: 0, fw: true, d: make([]node, 0)}
}

func (rd *reverseDeque) empty() bool { return len(rd.d) == 0 }

func (rd *reverseDeque) clear() { rd.d = rd.d[:0] }

func (rd *reverseDeque) r2i(x node) node {
	if rd.fw {
		x.l += rd.org
		x.r += rd.org
	} else {
		tmp := x.l
		x.l = -x.r + rd.org
		x.r = -tmp + rd.org
	}
	return x
}

func (rd *reverseDeque) back() node {
	var x node
	if rd.fw {
		x = rd.d[len(rd.d)-1]
	} else {
		x = rd.d[0]
	}
	return rd.r2i(x)
}

func (rd *reverseDeque) popBack() {
	if rd.fw {
		rd.d = rd.d[:len(rd.d)-1]
	} else {
		rd.d = rd.d[1:]
	}
}

func (rd *reverseDeque) pushBack(x node) { rd.d = append(rd.d, x) }

func (rd *reverseDeque) rev(v int64) {
	if rd.fw {
		rd.org += v
	} else {
		rd.org -= v
	}
	rd.fw = !rd.fw
}

func solve(tc testCase) (bool, []int64) {
	n := tc.n
	w := make([]int64, n-1)
	for i := 0; i < n-2; i++ {
		w[i] = tc.w[i]
	}
	w[n-2] = tc.C

	d := newReverseDeque()
	g := make([]int64, n-1)
	pre := make([]int, n-1)
	d.pushBack(node{l: 0, r: tc.C, i: -1})
	for i := 0; i < n-1; i++ {
		v := w[i]
		for !d.empty() {
			x := d.back()
			if x.r <= v {
				break
			}
			d.popBack()
			if x.l <= v {
				x.r = v
				d.pushBack(x)
			}
		}
		if d.empty() {
			return false, nil
		}
		x := d.back()
		g[i] = x.r
		pre[i] = x.i
		d.rev(v)
		if g[i] == v {
			d.clear()
			d.pushBack(node{l: 0, r: v, i: i})
		} else {
			d.pushBack(node{l: v, r: v, i: i})
		}
	}
	dif := make([]int64, n-1)
	mx := make([]int, n-2)
	for pos := n - 2; pos >= 0; pos = pre[pos] {
		nx := pre[pos]
		if nx >= 0 {
			mx[nx] = 1
		}
		dif[pos] = g[pos]
		for i := pos - 1; i > nx; i-- {
			dif[i] = w[i] - dif[i+1]
		}
	}
	h := make([]int64, n)
	if n > 1 {
		h[1] = dif[0]
	}
	for i := 0; i < n-2; i++ {
		a := h[i] < h[i+1]
		if mx[i] == 1 {
			a = !a
		}
		if a {
			h[i+2] = h[i+1] + dif[i+1]
		} else {
			h[i+2] = h[i+1] - dif[i+1]
		}
	}
	mn := h[0]
	for _, v := range h {
		if v < mn {
			mn = v
		}
	}
	off := -mn
	for i := range h {
		h[i] += off
	}
	return true, h
}

func checkHeights(h []int64, w []int64) bool {
	n := len(h)
	for i := 0; i < n-2; i++ {
		mx := h[i]
		mn := h[i]
		if h[i+1] > mx {
			mx = h[i+1]
		}
		if h[i+2] > mx {
			mx = h[i+2]
		}
		if h[i+1] < mn {
			mn = h[i+1]
		}
		if h[i+2] < mn {
			mn = h[i+2]
		}
		if mx-mn != w[i] {
			return false
		}
	}
	return true
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.C)
	for i, v := range tc.w {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	output, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("runtime error: %v\n%s", err, string(ee.Stderr))
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expectYes, expectHeights := solve(tc)
		var expect strings.Builder
		if expectYes {
			expect.WriteString("YES\n")
			for i, v := range expectHeights {
				if i > 0 {
					expect.WriteByte(' ')
				}
				expect.WriteString(strconv.FormatInt(v, 10))
			}
		} else {
			expect.WriteString("NO")
		}
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect.String() {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expect.String(), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

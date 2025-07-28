package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseB struct {
	n int
	q int
	s string
	f string
	l []int
	r []int
}

func generateTests() []testCaseB {
	rand.Seed(42)
	tests := make([]testCaseB, 100)
	for i := range tests {
		n := rand.Intn(10) + 1 // 1..10
		q := rand.Intn(10) + 1 // 1..10
		b := testCaseB{n: n, q: q}
		sb := make([]byte, n)
		fb := make([]byte, n)
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				sb[j] = '0'
			} else {
				sb[j] = '1'
			}
			if rand.Intn(2) == 0 {
				fb[j] = '0'
			} else {
				fb[j] = '1'
			}
		}
		b.s = string(sb)
		b.f = string(fb)
		b.l = make([]int, q)
		b.r = make([]int, q)
		for j := 0; j < q; j++ {
			l := rand.Intn(n) + 1
			r := rand.Intn(n) + 1
			if l > r {
				l, r = r, l
			}
			b.l[j] = l
			b.r[j] = r
		}
		tests[i] = b
	}
	return tests
}

// solveB replicates the algorithm from 1477B.go
func solveB(t testCaseB) string {
	type SegTree struct {
		n    int
		sum  []int
		lazy []int
	}

	newSegTree := func(arr []int) *SegTree {
		n := len(arr) - 1
		st := &SegTree{n: n, sum: make([]int, 4*n+5), lazy: make([]int, 4*n+5)}
		for i := range st.lazy {
			st.lazy[i] = -1
		}
		var build func(node, l, r int, arr []int)
		build = func(node, l, r int, arr []int) {
			if l == r {
				st.sum[node] = arr[l]
				return
			}
			mid := (l + r) / 2
			build(node*2, l, mid, arr)
			build(node*2+1, mid+1, r, arr)
			st.sum[node] = st.sum[node*2] + st.sum[node*2+1]
		}
		build(1, 1, n, arr)
		return st
	}

	var push func(node, l, r int, st *SegTree)
	push = func(node, l, r int, st *SegTree) {
		if st.lazy[node] == -1 || l == r {
			return
		}
		val := st.lazy[node]
		mid := (l + r) / 2
		st.sum[node*2] = (mid - l + 1) * val
		st.sum[node*2+1] = (r - mid) * val
		st.lazy[node*2] = val
		st.lazy[node*2+1] = val
		st.lazy[node] = -1
	}

	var update func(node, l, r, ql, qr, val int, st *SegTree)
	update = func(node, l, r, ql, qr, val int, st *SegTree) {
		if ql > r || qr < l {
			return
		}
		if ql <= l && r <= qr {
			st.sum[node] = (r - l + 1) * val
			st.lazy[node] = val
			return
		}
		push(node, l, r, st)
		mid := (l + r) / 2
		update(node*2, l, mid, ql, qr, val, st)
		update(node*2+1, mid+1, r, ql, qr, val, st)
		st.sum[node] = st.sum[node*2] + st.sum[node*2+1]
	}

	var query func(node, l, r, ql, qr int, st *SegTree) int
	query = func(node, l, r, ql, qr int, st *SegTree) int {
		if ql > r || qr < l {
			return 0
		}
		if ql <= l && r <= qr {
			return st.sum[node]
		}
		push(node, l, r, st)
		mid := (l + r) / 2
		res := 0
		if ql <= mid {
			res += query(node*2, l, mid, ql, qr, st)
		}
		if qr > mid {
			res += query(node*2+1, mid+1, r, ql, qr, st)
		}
		return res
	}

	arr := make([]int, t.n+1)
	for i := 0; i < t.n; i++ {
		if t.f[i] == '1' {
			arr[i+1] = 1
		}
	}
	st := newSegTree(arr)
	ok := true
	for i := t.q - 1; i >= 0 && ok; i-- {
		l := t.l[i]
		r := t.r[i]
		ones := query(1, 1, t.n, l, r, st)
		length := r - l + 1
		if ones*2 == length {
			ok = false
			break
		}
		if ones*2 > length {
			update(1, 1, t.n, l, r, 1, st)
		} else {
			update(1, 1, t.n, l, r, 0, st)
		}
	}
	if ok {
		for i := 0; i < t.n; i++ {
			val := query(1, 1, t.n, i+1, i+1, st)
			if val != int(t.s[i]-'0') {
				ok = false
				break
			}
		}
	}
	if ok {
		return "YES"
	}
	return "NO"
}

func buildInput(tests []testCaseB) string {
	var b strings.Builder
	fmt.Fprintln(&b, len(tests))
	for _, t := range tests {
		fmt.Fprintf(&b, "%d %d\n%s\n%s\n", t.n, t.q, t.s, t.f)
		for i := 0; i < t.q; i++ {
			fmt.Fprintf(&b, "%d %d\n", t.l[i], t.r[i])
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	input := buildInput(tests)

	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "execution failed:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(&out)
	outputs := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			outputs = append(outputs, strings.ToUpper(line))
		}
	}
	if len(outputs) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d lines of output, got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}
	for i, t := range tests {
		exp := solveB(t)
		if outputs[i] != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", i+1, exp, outputs[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}

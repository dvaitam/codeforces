package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"

)

// ---------- embedded solver from accepted solution ----------

func solverMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type SegTree struct {
	maxVal []int
	lazy   []int
	n      int
}

func NewSegTree(n int) *SegTree {
	return &SegTree{
		maxVal: make([]int, 4*n+1),
		lazy:   make([]int, 4*n+1),
		n:      n,
	}
}

func (st *SegTree) push(node int) {
	if st.lazy[node] != 0 {
		st.maxVal[2*node] += st.lazy[node]
		st.lazy[2*node] += st.lazy[node]
		st.maxVal[2*node+1] += st.lazy[node]
		st.lazy[2*node+1] += st.lazy[node]
		st.lazy[node] = 0
	}
}

func (st *SegTree) Add(node, l, r, ql, qr, val int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.maxVal[node] += val
		st.lazy[node] += val
		return
	}
	st.push(node)
	mid := (l + r) / 2
	st.Add(2*node, l, mid, ql, qr, val)
	st.Add(2*node+1, mid+1, r, ql, qr, val)
	st.maxVal[node] = solverMax(st.maxVal[2*node], st.maxVal[2*node+1])
}

type Item struct {
	val, idx, depth int
}

func compute_depths(C []int) []int {
	L := len(C)
	D := make([]int, L+1)
	D[0] = 0
	if L == 0 {
		return D
	}

	st := NewSegTree(L)
	stack := make([]Item, 0, L)

	for i := 1; i <= L; i++ {
		val := C[i-1]
		for len(stack) > 0 && stack[len(stack)-1].val > val {
			stack = stack[:len(stack)-1]
		}

		left_idx := 1
		if len(stack) > 0 {
			left_idx = stack[len(stack)-1].idx + 1
		}

		right_idx := i - 1
		if left_idx <= right_idx {
			st.Add(1, 1, L, left_idx, right_idx, 1)
		}

		depth := 1
		if len(stack) > 0 {
			depth = stack[len(stack)-1].depth + 1
		}

		st.Add(1, 1, L, i, i, depth)
		stack = append(stack, Item{val: val, idx: i, depth: depth})

		D[i] = st.maxVal[1]
	}
	return D
}

func solveFCase(a []int) (int, int) {
	n := len(a)
	if n == 1 {
		return 1, 0
	}

	p0 := -1
	for i := 0; i < n; i++ {
		if a[i] == 1 {
			p0 = i
			break
		}
	}

	B := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		B[i] = a[(p0+1+i)%n]
	}

	D_pref := compute_depths(B)

	B_rev := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		B_rev[i] = B[n-2-i]
	}
	D_pref_rev := compute_depths(B_rev)

	minDepth := int(1e9)
	bestS := -1

	for k := 0; k < n; k++ {
		dY := D_pref[k]
		dX := D_pref_rev[n-1-k]
		depth := 1 + solverMax(dY, dX)
		if depth < minDepth {
			minDepth = depth
			bestS = (p0 + 1 + k) % n
		}
	}

	return minDepth, bestS
}

// ---------- verifier logic ----------

type testF struct {
	a []int
}

func genTestsF() []testF {
	rng := rand.New(rand.NewSource(122006))
	tests := make([]testF, 100)
	for i := range tests {
		n := rng.Intn(8) + 2
		a := make([]int, n)
		pos := rng.Intn(n)
		for j := range a {
			if j == pos {
				a[j] = 1
			} else {
				a[j] = rng.Intn(9) + 2
			}
		}
		tests[i] = testF{a: a}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsF()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, len(tc.a))
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}

	expectedRet := make([]int, len(tests))
	expectedK := make([]int, len(tests))
	for i, tc := range tests {
		r, k := solveFCase(tc.a)
		expectedRet[i] = r
		expectedK[i] = k
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, out.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i := range tests {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		r, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
			os.Exit(1)
		}
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		k, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
			os.Exit(1)
		}
		if r != expectedRet[i] || k != expectedK[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}

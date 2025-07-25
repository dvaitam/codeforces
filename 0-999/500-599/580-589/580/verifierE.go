package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	mod1 int64 = 1000000007
	mod2 int64 = 1000000009
	base int64 = 911382323
)

var (
	pow1  []int64
	pow2  []int64
	ones1 []int64
	ones2 []int64
)

type Hash struct {
	h1, h2 int64
	len    int
}

type SegTree struct {
	n    int
	h1   []int64
	h2   []int64
	lazy []int
}

func NewSegTree(arr []int) *SegTree {
	n := len(arr) - 1
	size := 4 * n
	st := &SegTree{
		n:    n,
		h1:   make([]int64, size),
		h2:   make([]int64, size),
		lazy: make([]int, size),
	}
	for i := range st.lazy {
		st.lazy[i] = -1
	}
	st.build(1, 1, n, arr)
	return st
}

func (st *SegTree) build(node, l, r int, arr []int) {
	if l == r {
		v := int64(arr[l] + 1)
		st.h1[node] = v
		st.h2[node] = v
		return
	}
	mid := (l + r) >> 1
	st.build(node<<1, l, mid, arr)
	st.build(node<<1|1, mid+1, r, arr)
	st.pull(node, l, r)
}

func (st *SegTree) apply(node, l, r, c int) {
	length := r - l + 1
	v := int64(c + 1)
	st.h1[node] = v * ones1[length] % mod1
	st.h2[node] = v * ones2[length] % mod2
	st.lazy[node] = c
}

func (st *SegTree) push(node, l, r int) {
	if st.lazy[node] != -1 && l != r {
		mid := (l + r) >> 1
		c := st.lazy[node]
		st.apply(node<<1, l, mid, c)
		st.apply(node<<1|1, mid+1, r, c)
		st.lazy[node] = -1
	}
}

func (st *SegTree) pull(node, l, r int) {
	left, right := node<<1, node<<1|1
	mid := (l + r) >> 1
	lenR := r - mid
	st.h1[node] = (st.h1[left]*pow1[lenR] + st.h1[right]) % mod1
	st.h2[node] = (st.h2[left]*pow2[lenR] + st.h2[right]) % mod2
}

func (st *SegTree) Update(node, l, r, ql, qr, c int) {
	if ql <= l && r <= qr {
		st.apply(node, l, r, c)
		return
	}
	st.push(node, l, r)
	mid := (l + r) >> 1
	if ql <= mid {
		st.Update(node<<1, l, mid, ql, qr, c)
	}
	if qr > mid {
		st.Update(node<<1|1, mid+1, r, ql, qr, c)
	}
	st.pull(node, l, r)
}

func merge(a, b Hash) Hash {
	if a.len == 0 {
		return b
	}
	if b.len == 0 {
		return a
	}
	return Hash{
		h1:  (a.h1*pow1[b.len] + b.h1) % mod1,
		h2:  (a.h2*pow2[b.len] + b.h2) % mod2,
		len: a.len + b.len,
	}
}

func (st *SegTree) Query(node, l, r, ql, qr int) Hash {
	if ql > qr {
		return Hash{0, 0, 0}
	}
	if ql <= l && r <= qr {
		return Hash{st.h1[node], st.h2[node], r - l + 1}
	}
	st.push(node, l, r)
	mid := (l + r) >> 1
	if qr <= mid {
		return st.Query(node<<1, l, mid, ql, qr)
	}
	if ql > mid {
		return st.Query(node<<1|1, mid+1, r, ql, qr)
	}
	leftHash := st.Query(node<<1, l, mid, ql, qr)
	rightHash := st.Query(node<<1|1, mid+1, r, ql, qr)
	return merge(leftHash, rightHash)
}

func solve(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	defer writer.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return ""
	}
	var str string
	fmt.Fscan(reader, &str)

	pow1 = make([]int64, n+2)
	pow2 = make([]int64, n+2)
	ones1 = make([]int64, n+2)
	ones2 = make([]int64, n+2)
	pow1[0], pow2[0] = 1, 1
	for i := 1; i <= n+1; i++ {
		pow1[i] = pow1[i-1] * base % mod1
		pow2[i] = pow2[i-1] * base % mod2
	}
	for i := 1; i <= n+1; i++ {
		ones1[i] = (ones1[i-1]*base + 1) % mod1
		ones2[i] = (ones2[i-1]*base + 1) % mod2
	}

	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		arr[i] = int(str[i-1] - '0')
	}
	st := NewSegTree(arr)

	total := m + k
	for ; total > 0; total-- {
		var tp int
		fmt.Fscan(reader, &tp)
		if tp == 1 {
			var l, r, c int
			fmt.Fscan(reader, &l, &r, &c)
			st.Update(1, 1, n, l, r, c)
		} else {
			var l, r, d int
			fmt.Fscan(reader, &l, &r, &d)
			if l > r-d {
				fmt.Fprintln(writer, "YES")
				continue
			}
			h1 := st.Query(1, 1, n, l, r-d)
			h2 := st.Query(1, 1, n, l+d, r)
			if h1.h1 == h2.h1 && h1.h2 == h2.h2 {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
	writer.Flush()
	return buf.String()
}

func buildCase(n, m, k int, str string, ops []string) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	sb.WriteString(str)
	sb.WriteByte('\n')
	for _, op := range ops {
		sb.WriteString(op)
		sb.WriteByte('\n')
	}
	input := sb.String()
	return testCase{input: input, expected: solve(input)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	m := rng.Intn(3) + 1
	k := rng.Intn(3) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	str := sb.String()
	var ops []string
	for i := 0; i < m; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		c := rng.Intn(2)
		ops = append(ops, fmt.Sprintf("1 %d %d %d", l, r, c))
	}
	for i := 0; i < k; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		d := 0
		if r-l > 0 {
			d = rng.Intn(r - l)
		}
		ops = append(ops, fmt.Sprintf("2 %d %d %d", l, r, d))
	}
	rng.Shuffle(len(ops), func(i, j int) { ops[i], ops[j] = ops[j], ops[i] })
	return buildCase(n, m, k, str, ops)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	ops := []string{"1 1 1 1", "2 1 1 0"}
	cases := []testCase{buildCase(1, 1, 1, "0", ops)}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

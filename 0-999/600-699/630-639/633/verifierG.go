package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	nG        int
	MG        int
	g         [][]int
	start     []int
	endt      []int
	flat      []int
	timer     int
	mask      *big.Int
	primeMask *big.Int
)

type SegmentTree struct {
	n    int
	tree []*big.Int
	lazy []int
}

func NewSegmentTree(vals []int) *SegmentTree {
	st := &SegmentTree{
		n:    len(vals),
		tree: make([]*big.Int, len(vals)*4),
		lazy: make([]int, len(vals)*4),
	}
	st.build(1, 0, st.n-1, vals)
	return st
}

func (st *SegmentTree) build(node, l, r int, vals []int) {
	if l == r {
		st.tree[node] = new(big.Int).SetBit(new(big.Int), vals[l], 1)
		return
	}
	mid := (l + r) / 2
	st.build(node*2, l, mid, vals)
	st.build(node*2+1, mid+1, r, vals)
	st.pull(node)
}

func (st *SegmentTree) pull(node int) {
	if st.tree[node] == nil {
		st.tree[node] = new(big.Int)
	}
	st.tree[node].Or(st.tree[node*2], st.tree[node*2+1])
}

func rotateInPlace(b *big.Int, k int) {
	k %= MG
	if k == 0 {
		return
	}
	tmp := new(big.Int).Set(b)
	b.Lsh(b, uint(k))
	tmp.Rsh(tmp, uint(MG-k))
	b.Or(b, tmp)
	b.And(b, mask)
}

func (st *SegmentTree) apply(node, k int) {
	rotateInPlace(st.tree[node], k)
	st.lazy[node] = (st.lazy[node] + k) % MG
}

func (st *SegmentTree) push(node int) {
	if st.lazy[node] != 0 {
		k := st.lazy[node]
		st.apply(node*2, k)
		st.apply(node*2+1, k)
		st.lazy[node] = 0
	}
}

func (st *SegmentTree) update(node, l, r, ql, qr, k int) {
	if ql <= l && r <= qr {
		st.apply(node, k)
		return
	}
	st.push(node)
	mid := (l + r) / 2
	if ql <= mid {
		st.update(node*2, l, mid, ql, qr, k)
	}
	if qr > mid {
		st.update(node*2+1, mid+1, r, ql, qr, k)
	}
	st.pull(node)
}

func (st *SegmentTree) Update(l, r, k int) {
	k %= MG
	if k < 0 {
		k += MG
	}
	st.update(1, 0, st.n-1, l, r, k)
}

func (st *SegmentTree) query(node, l, r, ql, qr int, res *big.Int) {
	if ql <= l && r <= qr {
		res.Or(res, st.tree[node])
		return
	}
	st.push(node)
	mid := (l + r) / 2
	if ql <= mid {
		st.query(node*2, l, mid, ql, qr, res)
	}
	if qr > mid {
		st.query(node*2+1, mid+1, r, ql, qr, res)
	}
}

func (st *SegmentTree) Query(l, r int) *big.Int {
	res := new(big.Int)
	st.query(1, 0, st.n-1, l, r, res)
	return res
}

func sievePrimes(m int) []int {
	isPrime := make([]bool, m)
	for i := 2; i < m; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i < m; i++ {
		if isPrime[i] {
			for j := i * i; j < m; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := []int{}
	for i := 2; i < m; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func bitCount(b *big.Int) int {
	cnt := 0
	for _, w := range b.Bits() {
		cnt += bits.OnesCount(uint(w))
	}
	return cnt
}

func dfs(u, p int) {
	start[u] = timer
	flat[timer] = u
	timer++
	for _, v := range g[u] {
		if v != p {
			dfs(v, u)
		}
	}
	endt[u] = timer - 1
}

func solveG(input string) string {
	reader := strings.NewReader(input)
	fmt.Fscan(reader, &nG, &MG)
	valsOrig := make([]int, nG+1)
	for i := 1; i <= nG; i++ {
		fmt.Fscan(reader, &valsOrig[i])
	}
	g = make([][]int, nG+1)
	for i := 0; i < nG-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	start = make([]int, nG+1)
	endt = make([]int, nG+1)
	flat = make([]int, nG)
	timer = 0
	dfs(1, 0)
	values := make([]int, nG)
	for i := 0; i < nG; i++ {
		node := flat[i]
		values[i] = valsOrig[node] % MG
	}
	mask = new(big.Int).Lsh(big.NewInt(1), uint(MG))
	mask.Sub(mask, big.NewInt(1))
	primeMask = new(big.Int)
	for _, p := range sievePrimes(MG) {
		primeMask.SetBit(primeMask, p, 1)
	}
	st := NewSegmentTree(values)
	var q int
	fmt.Fscan(reader, &q)
	var out strings.Builder
	for ; q > 0; q-- {
		var typ, v int
		fmt.Fscan(reader, &typ, &v)
		if typ == 1 {
			var x int
			fmt.Fscan(reader, &x)
			st.Update(start[v], endt[v], x%MG)
		} else if typ == 2 {
			res := st.Query(start[v], endt[v])
			res.And(res, primeMask)
			cnt := bitCount(res)
			out.WriteString(fmt.Sprintf("%d\n", cnt))
		}
	}
	return out.String()
}

func generateCaseG(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	M := rng.Intn(9) + 2
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		vals[i] = rng.Intn(M)
	}
	edges := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		edges[i-1][0] = i + 1
		edges[i-1][1] = rng.Intn(i) + 1
	}
	q := rng.Intn(4) + 1
	var qs []string
	hasType2 := false
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			v := rng.Intn(n) + 1
			x := rng.Intn(M)
			qs = append(qs, fmt.Sprintf("1 %d %d", v, x))
		} else {
			v := rng.Intn(n) + 1
			qs = append(qs, fmt.Sprintf("2 %d", v))
			hasType2 = true
		}
	}
	if !hasType2 {
		v := rng.Intn(n) + 1
		qs = append(qs, fmt.Sprintf("2 %d", v))
		q++
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, M))
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for _, qq := range qs {
		sb.WriteString(qq)
		sb.WriteByte('\n')
	}
	expect := solveG(sb.String())
	return sb.String(), expect
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCaseG(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

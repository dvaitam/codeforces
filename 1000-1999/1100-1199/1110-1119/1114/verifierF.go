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

type op struct {
	typ string
	l   int
	r   int
	x   int
}

type testcase struct {
	n   int
	q   int
	arr []int
	ops []op
}

const mod = 1000000007

// Precomputed primes up to 300 for factor masks.
var primes []int
var invs []int
var maskArr []uint64

const testcasesRaw = `2 3 2 7 MULTIPLY 1 1 1 MULTIPLY 2 2 4 TOTIENT 2 2
2 1 5 4 TOTIENT 1 2
3 2 3 5 5 TOTIENT 3 3 MULTIPLY 3 3 7
5 2 3 4 8 5 2 TOTIENT 5 5 MULTIPLY 3 5 5
5 2 7 7 10 5 7 MULTIPLY 2 4 5 TOTIENT 1 1
1 4 5 TOTIENT 1 1 MULTIPLY 1 1 7 TOTIENT 1 1 MULTIPLY 1 1 9
2 3 2 1 TOTIENT 2 2 MULTIPLY 1 2 8 MULTIPLY 2 2 5
3 1 6 5 6 TOTIENT 3 3
5 1 5 10 4 8 5 MULTIPLY 4 4 6
5 1 6 1 8 3 6 TOTIENT 3 4
5 1 8 4 7 4 2 MULTIPLY 1 2 10
2 5 1 9 MULTIPLY 1 2 1 MULTIPLY 2 2 4 MULTIPLY 1 2 7 MULTIPLY 1 2 8 MULTIPLY 2 2 8
2 1 1 5 MULTIPLY 1 1 7
3 2 6 1 6 TOTIENT 3 3 TOTIENT 3 3
4 4 2 7 4 10 TOTIENT 2 3 MULTIPLY 4 4 7 TOTIENT 3 4 TOTIENT 4 4
3 2 1 7 10 MULTIPLY 2 2 3 TOTIENT 2 3
4 4 4 1 4 3 MULTIPLY 3 3 7 TOTIENT 4 4 TOTIENT 2 2 TOTIENT 3 4
5 4 1 2 1 10 2 MULTIPLY 3 5 3 MULTIPLY 1 5 1 MULTIPLY 3 3 2 TOTIENT 4 5
2 3 7 4 TOTIENT 2 2 MULTIPLY 2 2 7 TOTIENT 2 2
2 3 8 7 MULTIPLY 1 2 3 MULTIPLY 1 2 8 TOTIENT 2 2
2 1 5 2 TOTIENT 2 2
5 1 9 4 7 5 6 MULTIPLY 2 2 1
5 3 9 8 10 5 9 TOTIENT 5 5 MULTIPLY 3 5 6 TOTIENT 3 3
4 1 10 3 10 3 MULTIPLY 3 3 10
3 5 2 2 7 TOTIENT 1 2 TOTIENT 2 2 MULTIPLY 1 3 1 TOTIENT 1 2 TOTIENT 1 1
2 5 8 10 TOTIENT 1 1 MULTIPLY 2 2 10 TOTIENT 2 2 TOTIENT 1 2 TOTIENT 2 2
2 3 9 9 TOTIENT 2 2 TOTIENT 2 2 MULTIPLY 1 2 3
5 4 9 10 6 7 7 TOTIENT 2 2 MULTIPLY 3 5 5 TOTIENT 3 4 MULTIPLY 4 4 3
2 1 5 3 TOTIENT 1 2
5 5 9 9 5 1 9 MULTIPLY 1 2 1 TOTIENT 1 1 MULTIPLY 3 5 3 TOTIENT 4 4 TOTIENT 1 5
3 1 9 3 9 MULTIPLY 1 2 8
5 1 7 3 3 5 9 MULTIPLY 3 4 6
2 4 1 7 TOTIENT 2 2 MULTIPLY 1 1 3 MULTIPLY 1 1 4 TOTIENT 2 2
4 4 6 3 9 10 MULTIPLY 4 4 5 TOTIENT 4 4 TOTIENT 2 3 TOTIENT 3 4
3 1 4 6 2 TOTIENT 3 3
5 3 9 8 6 2 3 TOTIENT 4 4 MULTIPLY 1 3 10 TOTIENT 3 3
5 5 1 2 4 10 5 TOTIENT 1 1 MULTIPLY 2 3 6 MULTIPLY 3 4 7 TOTIENT 4 4 MULTIPLY 4 4 9
5 3 9 5 4 7 6 MULTIPLY 1 4 7 TOTIENT 3 4 TOTIENT 4 5
3 5 4 5 5 TOTIENT 1 2 MULTIPLY 2 2 10 TOTIENT 1 1 TOTIENT 3 3 TOTIENT 2 2
4 1 6 1 8 10 TOTIENT 2 2
2 5 2 2 TOTIENT 2 2 TOTIENT 1 1 TOTIENT 1 1 TOTIENT 1 1 MULTIPLY 2 2 10
3 2 3 4 1 TOTIENT 1 3 TOTIENT 1 3
3 4 2 4 3 MULTIPLY 3 3 9 MULTIPLY 2 3 1 TOTIENT 1 1 MULTIPLY 1 2 4
3 5 4 10 6 MULTIPLY 3 3 4 MULTIPLY 2 2 9 MULTIPLY 1 1 4 MULTIPLY 3 3 9 MULTIPLY 1 2 7
2 2 7 10 MULTIPLY 2 2 8 TOTIENT 2 2
4 2 2 6 5 10 TOTIENT 1 4 MULTIPLY 1 2 10
2 5 6 10 MULTIPLY 1 2 7 MULTIPLY 1 1 3 MULTIPLY 2 2 7 TOTIENT 2 2 TOTIENT 2 2
5 3 9 1 9 3 6 TOTIENT 2 5 MULTIPLY 1 1 10 MULTIPLY 2 3 9
2 5 1 9 TOTIENT 1 1 TOTIENT 1 1 MULTIPLY 1 2 8 TOTIENT 2 2 TOTIENT 2 2
2 3 7 1 MULTIPLY 2 2 4 MULTIPLY 1 2 1 MULTIPLY 1 1 4
1 1 1 TOTIENT 1 1
1 1 9 TOTIENT 1 1
5 5 1 1 2 3 7 TOTIENT 1 5 MULTIPLY 3 5 5 MULTIPLY 5 5 10 TOTIENT 1 5 TOTIENT 2 2
3 5 7 5 10 MULTIPLY 1 3 6 TOTIENT 1 3 MULTIPLY 2 2 10 TOTIENT 1 2 TOTIENT 3 3
3 2 5 3 6 TOTIENT 2 2 TOTIENT 3 3
3 3 3 6 2 TOTIENT 2 2 MULTIPLY 1 2 1 TOTIENT 1 1
5 2 4 4 6 7 7 MULTIPLY 4 5 5 TOTIENT 3 3
1 2 1 MULTIPLY 1 1 7 TOTIENT 1 1
3 3 9 6 3 MULTIPLY 2 2 7 TOTIENT 2 3 TOTIENT 2 2
4 2 2 8 7 9 MULTIPLY 1 1 5 MULTIPLY 2 3 8
3 4 6 7 5 TOTIENT 2 2 TOTIENT 2 3 TOTIENT 1 3 MULTIPLY 3 3 4
5 5 3 7 5 2 6 TOTIENT 5 5 TOTIENT 4 4 TOTIENT 3 4 MULTIPLY 5 5 8 TOTIENT 1 4
3 4 5 4 3 MULTIPLY 1 2 9 TOTIENT 3 3 TOTIENT 1 2 TOTIENT 3 3
3 5 9 1 7 MULTIPLY 3 3 9 MULTIPLY 1 2 10 MULTIPLY 1 1 4 TOTIENT 1 2 TOTIENT 2 2
3 1 2 1 8 MULTIPLY 3 3 4
4 5 2 9 10 3 MULTIPLY 1 2 9 TOTIENT 2 2 TOTIENT 3 4 MULTIPLY 4 4 6 TOTIENT 1 3
3 1 1 6 1 TOTIENT 2 2
2 2 1 8 TOTIENT 2 2 TOTIENT 1 1
4 4 9 7 2 3 MULTIPLY 2 4 8 TOTIENT 4 4 MULTIPLY 4 4 5 MULTIPLY 1 3 7
5 3 5 2 3 9 4 TOTIENT 2 5 TOTIENT 2 3 TOTIENT 2 2
5 1 5 1 2 9 8 TOTIENT 5 5
3 3 4 4 2 TOTIENT 3 3 MULTIPLY 2 2 2 MULTIPLY 3 3 8
5 3 4 8 4 3 3 MULTIPLY 1 2 4 TOTIENT 4 5 MULTIPLY 2 4 10
5 4 8 6 6 2 8 MULTIPLY 3 3 2 TOTIENT 3 5 MULTIPLY 5 5 8 MULTIPLY 3 5 9
1 5 5 TOTIENT 1 1 TOTIENT 1 1 MULTIPLY 1 1 4 TOTIENT 1 1 MULTIPLY 1 1 7
3 3 7 4 1 TOTIENT 3 3 TOTIENT 3 3 TOTIENT 1 3
4 3 5 1 3 6 MULTIPLY 4 4 6 MULTIPLY 4 4 7 TOTIENT 4 4
2 4 7 5 TOTIENT 2 2 MULTIPLY 1 1 3 MULTIPLY 1 1 8 TOTIENT 1 2
5 2 1 7 7 5 8 MULTIPLY 1 1 3 TOTIENT 3 5
4 1 6 5 9 8 MULTIPLY 1 3 5
3 3 7 6 1 TOTIENT 1 2 TOTIENT 3 3 MULTIPLY 1 1 8
5 3 7 2 7 7 8 TOTIENT 4 4 TOTIENT 4 4 MULTIPLY 3 4 2
4 5 10 9 4 5 TOTIENT 4 4 TOTIENT 3 3 MULTIPLY 4 4 9 TOTIENT 3 3 MULTIPLY 3 3 5
4 4 5 7 9 3 TOTIENT 4 4 TOTIENT 3 4 TOTIENT 3 4 TOTIENT 2 3
1 4 10 MULTIPLY 1 1 3 MULTIPLY 1 1 3 MULTIPLY 1 1 6 MULTIPLY 1 1 2
3 5 9 2 9 TOTIENT 3 3 TOTIENT 3 3 MULTIPLY 3 3 7 TOTIENT 1 3 TOTIENT 2 2
2 2 8 4 TOTIENT 1 1 MULTIPLY 1 2 3
5 2 6 7 2 1 6 TOTIENT 3 3 MULTIPLY 3 5 4
2 5 2 9 TOTIENT 2 2 MULTIPLY 1 1 3 TOTIENT 1 1 TOTIENT 2 2 TOTIENT 2 2
2 1 9 6 TOTIENT 2 2
1 5 1 TOTIENT 1 1 TOTIENT 1 1 TOTIENT 1 1 TOTIENT 1 1 MULTIPLY 1 1 6
5 1 9 4 10 3 7 TOTIENT 5 5
4 1 5 3 7 4 TOTIENT 3 3
2 4 5 4 TOTIENT 2 2 MULTIPLY 1 1 8 TOTIENT 2 2 MULTIPLY 1 1 3
2 1 5 3 TOTIENT 2 2
3 1 10 2 7 MULTIPLY 2 3 9
4 2 2 7 8 8 TOTIENT 4 4 MULTIPLY 3 3 5
2 4 4 9 MULTIPLY 2 2 7 MULTIPLY 2 2 4 TOTIENT 1 2 MULTIPLY 1 1 9
5 1 10 6 8 7 9 TOTIENT 5 5
2 2 2 4 MULTIPLY 1 1 2 MULTIPLY 2 2 1`

var testcases = mustParseTestcases(testcasesRaw)

func init() {
	// sieve primes up to 300
	maxv := 300
	isComp := make([]bool, maxv+1)
	for i := 2; i <= maxv; i++ {
		if !isComp[i] {
			primes = append(primes, i)
		}
		for _, p := range primes {
			if i*p > maxv {
				break
			}
			isComp[i*p] = true
			if i%p == 0 {
				break
			}
		}
	}
	invs = make([]int, len(primes))
	for i, p := range primes {
		invs[i] = fpow(p, mod-2)
	}
	maskArr = make([]uint64, maxv+1)
	for v := 2; v <= maxv; v++ {
		var m uint64
		for i, p := range primes {
			if v%p == 0 {
				m |= 1 << uint(i)
			}
		}
		maskArr[v] = m
	}
}

func mustParseTestcases(raw string) []testcase {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	cases := make([]testcase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		scan := bufio.NewScanner(strings.NewReader(line))
		scan.Split(bufio.ScanWords)
		nextInt := func() int {
			if !scan.Scan() {
				panic(fmt.Sprintf("line %d: unexpected EOF", idx+1))
			}
			v, err := strconv.Atoi(scan.Text())
			if err != nil {
				panic(fmt.Sprintf("line %d: bad int %q: %v", idx+1, scan.Text(), err))
			}
			return v
		}

		n := nextInt()
		q := nextInt()
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = nextInt()
		}
		ops := make([]op, 0, q)
		for i := 0; i < q; i++ {
			if !scan.Scan() {
				panic(fmt.Sprintf("line %d: missing op at %d", idx+1, i+1))
			}
			typ := scan.Text()
			l := nextInt()
			r := nextInt()
			var x int
			if typ == "MULTIPLY" {
				x = nextInt()
			}
			ops = append(ops, op{typ: typ, l: l, r: r, x: x})
		}
		cases = append(cases, testcase{n: n, q: q, arr: arr, ops: ops})
	}
	return cases
}

// fpow computes a^b mod mod.
func fpow(a, b int) int {
	res := 1
	a %= mod
	for b > 0 {
		if b&1 != 0 {
			res = int((int64(res) * int64(a)) % mod)
		}
		a = int((int64(a) * int64(a)) % mod)
		b >>= 1
	}
	return res
}

type Node struct {
	prod     int
	mask     uint64
	lazyMul  int
	lazyMask uint64
}

type solver struct {
	n    int
	arr  []int
	tree []Node
}

func newSolver(arr []int) *solver {
	n := len(arr) - 1 // arr is 1-indexed
	s := &solver{n: n, arr: arr, tree: make([]Node, 4*(n+2))}
	s.build(1, 1, n)
	return s
}

func (s *solver) build(idx, l, r int) {
	s.tree[idx].lazyMul = 1
	s.tree[idx].lazyMask = 0
	if l == r {
		s.tree[idx].prod = s.arr[l]
		s.tree[idx].mask = s.getMask(s.arr[l])
		return
	}
	mid := (l + r) >> 1
	s.build(idx<<1, l, mid)
	s.build(idx<<1|1, mid+1, r)
	s.pull(idx)
}

func (s *solver) pull(idx int) {
	s.tree[idx].prod = int((int64(s.tree[idx<<1].prod) * int64(s.tree[idx<<1|1].prod)) % mod)
	s.tree[idx].mask = s.tree[idx<<1].mask | s.tree[idx<<1|1].mask
}

func (s *solver) push(idx, l, r int) {
	lm, lmask := s.tree[idx].lazyMul, s.tree[idx].lazyMask
	if lm == 1 && lmask == 0 {
		return
	}
	mid := (l + r) >> 1
	left := idx << 1
	s.applyMul(left, l, mid, lm, lmask)
	right := idx<<1 | 1
	s.applyMul(right, mid+1, r, lm, lmask)
	s.tree[idx].lazyMul = 1
	s.tree[idx].lazyMask = 0
}

func (s *solver) applyMul(idx, l, r, mul int, mask uint64) {
	s.tree[idx].prod = int((int64(s.tree[idx].prod) * int64(fpow(mul, r-l+1))) % mod)
	s.tree[idx].lazyMul = int((int64(s.tree[idx].lazyMul) * int64(mul)) % mod)
	s.tree[idx].mask |= mask
	s.tree[idx].lazyMask |= mask
}

func (s *solver) update(idx, l, r, ql, qr, mul int, mask uint64) {
	if ql <= l && r <= qr {
		s.applyMul(idx, l, r, mul, mask)
		return
	}
	s.push(idx, l, r)
	mid := (l + r) >> 1
	if ql <= mid {
		s.update(idx<<1, l, mid, ql, qr, mul, mask)
	}
	if qr > mid {
		s.update(idx<<1|1, mid+1, r, ql, qr, mul, mask)
	}
	s.pull(idx)
}

func (s *solver) queryProd(idx, l, r, ql, qr int) int {
	if ql <= l && r <= qr {
		return s.tree[idx].prod
	}
	s.push(idx, l, r)
	mid := (l + r) >> 1
	res := 1
	if ql <= mid {
		res = int((int64(res) * int64(s.queryProd(idx<<1, l, mid, ql, qr))) % mod)
	}
	if qr > mid {
		res = int((int64(res) * int64(s.queryProd(idx<<1|1, mid+1, r, ql, qr))) % mod)
	}
	return res
}

func (s *solver) queryMask(idx, l, r, ql, qr int) uint64 {
	if ql <= l && r <= qr {
		return s.tree[idx].mask
	}
	s.push(idx, l, r)
	mid := (l + r) >> 1
	var res uint64
	if ql <= mid {
		res |= s.queryMask(idx<<1, l, mid, ql, qr)
	}
	if qr > mid {
		res |= s.queryMask(idx<<1|1, mid+1, r, ql, qr)
	}
	return res
}

func (s *solver) getMask(x int) uint64 {
	if x < len(maskArr) {
		return maskArr[x]
	}
	var m uint64
	for i, p := range primes {
		if x%p == 0 {
			m |= 1 << uint(i)
		}
	}
	return m
}

func solve(tc testcase) []string {
	arr := make([]int, tc.n+1)
	copy(arr[1:], tc.arr)
	sol := newSolver(arr)
	outputs := []string{}
	for _, op := range tc.ops {
		if op.typ == "MULTIPLY" {
			sol.update(1, 1, sol.n, op.l, op.r, op.x, sol.getMask(op.x))
		} else { // TOTIENT
			pm := sol.queryMask(1, 1, sol.n, op.l, op.r)
			res := sol.queryProd(1, 1, sol.n, op.l, op.r)
			for i, p := range primes {
				if pm&(1<<uint(i)) != 0 {
					res = int((int64(res) * int64(invs[i]) % mod) * int64(p-1) % mod)
				}
			}
			outputs = append(outputs, strconv.Itoa(res))
		}
	}
	return outputs
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseCandidateOutput(out string) []string {
	return strings.Fields(out)
}

func checkCase(bin string, idx int, tc testcase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, op := range tc.ops {
		sb.WriteString(op.typ)
		sb.WriteByte(' ')
		sb.WriteString(fmt.Sprintf("%d %d", op.l, op.r))
		if op.typ == "MULTIPLY" {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(op.x))
		}
		sb.WriteByte('\n')
	}
	input := sb.String()

	expected := solve(tc)
	out, err := runCandidate(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\ninput:\n%s", err, input)
	}
	got := parseCandidateOutput(out)
	if len(got) != len(expected) {
		return fmt.Errorf("expected %d outputs, got %d\ninput:\n%s", len(expected), len(got), input)
	}
	for i := range expected {
		if got[i] != expected[i] {
			return fmt.Errorf("output %d: expected %s got %s\ninput:\n%s", i+1, expected[i], got[i], input)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range testcases {
		if err := checkCase(bin, i, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

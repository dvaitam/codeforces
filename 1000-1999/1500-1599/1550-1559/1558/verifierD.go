package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD int64 = 998244353

type pair struct {
	x int
	y int
}

type testCase struct {
	n   int
	ops []pair
}

// Embedded testcases from testcasesD.txt.
const testcaseData = `
100
3 2
2 1
3 1
6 0
5 2
2 1
3 2
5 1
3 1
6 3
2 1
5 3
6 1
4 3
2 1
3 2
4 1
4 0
2 0
5 1
4 3
5 2
4 2
5 3
6 4
2 1
3 1
5 3
6 5
6 4
2 1
3 1
4 1
5 4
5 0
4 0
5 1
2 1
5 3
2 1
4 3
5 3
6 2
3 1
6 3
2 0
2 0
3 1
3 2
3 2
2 1
3 2
3 1
3 2
6 3
2 1
4 2
6 2
4 3
2 1
3 2
4 3
4 0
5 4
2 1
3 2
4 3
5 3
5 0
6 0
2 1
2 1
4 2
2 1
3 2
4 2
3 1
4 1
6 5
2 1
3 2
4 1
5 4
6 1
2 1
2 1
5 1
2 1
3 2
2 1
3 1
2 0
4 1
3 2
2 1
2 1
4 2
3 2
4 2
5 4
2 1
3 1
4 2
5 4
6 5
2 1
3 1
4 1
5 3
6 1
3 0
2 0
5 4
2 1
3 2
4 1
5 1
6 4
2 1
4 1
5 3
6 5
5 0
4 1
2 1
6 0
3 0
4 1
2 1
6 3
2 1
4 3
6 5
6 3
2 1
3 1
5 2
2 0
2 0
5 0
2 1
2 1
4 0
4 3
2 1
3 2
4 1
4 3
2 1
3 1
4 3
3 0
4 3
2 1
3 1
4 2
5 2
4 2
5 4
2 0
3 2
2 1
3 1
5 1
2 1
4 2
2 1
3 2
2 1
2 1
5 1
3 2
2 0
6 0
5 1
2 1
2 0
2 1
2 1
2 0
6 4
2 1
3 2
4 3
6 2
2 0
3 2
2 1
3 2
4 3
2 1
3 1
4 2
3 1
3 1
5 1
3 1
5 3
2 1
4 2
5 2
3 2
2 1
3 2
6 4
2 1
3 2
4 2
6 4
3 0
6 2
3 2
4 1
5 3
2 1
3 1
5 4
2 1
2 1
2 0
4 1
2 1
3 0
6 5
2 1
3 2
4 1
5 1
6 1
2 1
2 1
4 3
2 1
3 2
4 2
5 3
2 1
4 2
5 2
6 2
2 1
4 3
5 0
5 4
2 1
3 2
4 3
5 3
2 1
2 1
6 2
2 1
4 3
3 0
`

// Combination utilities.
var fact, inv []int64

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func initComb(n int) {
	fact = make([]int64, n+1)
	inv = make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	inv[n] = modPow(fact[n], MOD-2)
	for i := n; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % MOD
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * inv[k] % MOD * inv[n-k] % MOD
}

type BIT struct {
	n    int
	tree []int
}

func (b *BIT) init(n int) { b.n = n; b.tree = make([]int, n+2) }
func (b *BIT) add(i, delta int) {
	for i <= b.n {
		b.tree[i] += delta
		i += i & -i
	}
}
func (b *BIT) kth(k int) int {
	idx := 0
	bit := 1
	for bit<<1 <= b.n {
		bit <<= 1
	}
	for bit > 0 {
		next := idx + bit
		if next <= b.n && b.tree[next] < k {
			k -= b.tree[next]
			idx = next
		}
		bit >>= 1
	}
	return idx + 1
}

// solve mirrors 1558D.go.
func solve(tc testCase) string {
	n := tc.n
	m := len(tc.ops)
	j := make([]int, n+1)
	for i := 1; i <= n; i++ {
		j[i] = i
	}
	for k := 0; k < m; k++ {
		x := tc.ops[k].x
		y := tc.ops[k].y
		j[x] = y
	}
	bit := BIT{}
	bit.init(n)
	for i := 1; i <= n; i++ {
		bit.add(i, 1)
	}
	p := make([]int, n+1)
	for i := n; i >= 1; i-- {
		pos := bit.kth(j[i])
		p[pos] = i
		bit.add(pos, -1)
	}
	r := 0
	for i := 1; i < n; i++ {
		if p[i] > p[i+1] {
			r++
		}
	}
	ans := comb(2*n-1-r, n)
	return strconv.FormatInt(ans, 10)
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos+1 >= len(fields) {
			return nil, fmt.Errorf("case %d missing n/m", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", caseIdx+1, err)
		}
		m, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, fmt.Errorf("case %d bad m: %v", caseIdx+1, err)
		}
		pos += 2
		ops := make([]pair, m)
		for i := 0; i < m; i++ {
			if pos+1 >= len(fields) {
				return nil, fmt.Errorf("case %d missing pair %d", caseIdx+1, i+1)
			}
			x, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d bad x: %v", caseIdx+1, err)
			}
			y, err := strconv.Atoi(fields[pos+1])
			if err != nil {
				return nil, fmt.Errorf("case %d bad y: %v", caseIdx+1, err)
			}
			ops[i] = pair{x: x, y: y}
			pos += 2
		}
		cases = append(cases, testCase{n: n, ops: ops})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra data at end")
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.ops)))
	for _, op := range tc.ops {
		sb.WriteString(fmt.Sprintf("%d %d\n", op.x, op.y))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	initComb(400000)

	for idx, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

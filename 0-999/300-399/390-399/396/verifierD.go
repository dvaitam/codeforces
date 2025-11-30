package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesD = `100
3
2 1 3
4
2 1 3 4
6
2 6 5 3 4 1
3
1 2 3
6
2 1 3 6 4 5
3
2 1 3
6
2 3 6 5 1 4
4
4 3 2 1
6
1 3 2 5 4 6
3
3 1 2
2
2 1
5
1 5 4 3 2
4
1 2 3 4
6
2 1 6 4 5 3
4
3 4 2 1
4
2 4 3 1
6
5 2 4 1 6 3
5
4 2 5 3 1
3
3 2 1
5
2 3 5 1 4
5
1 2 4 3 5
6
6 4 5 3 1 2
2
2 1
3
1 3 2
6
4 6 5 1 2 3
4
1 3 4 2
5
2 1 3 4 5
6
2 1 6 4 3 5
4
1 3 2 4
6
5 2 6 4 1 3
2
1 2
6
4 6 3 1 2 5
5
1 2 4 5 3
5
2 4 3 5 1
4
1 2 4 3
4
4 2 1 3
3
1 3 2
6
2 5 1 4 6 3
6
5 1 4 3 2 6
4
3 1 4 2
5
4 2 3 5 1
4
3 1 4 2
2
1 2
3
1 2 3
4
4 3 1 2
6
1 5 4 6 3 2
2
1 2
6
1 6 5 3 4 2
5
1 5 2 4 3
6
2 3 5 6 1 4
2
1 2
6
3 6 1 2 4 5
6
2 4 1 6 5 3
6
4 6 5 2 1 3
6
3 6 5 1 4 2
2
2 1
6
4 5 6 1 2 3
2
2 1
4
1 3 4 2
2
1 2
3
2 3 1
6
3 6 4 5 2 1
2
1 2
6
6 2 5 3 1 4
6
3 2 1 4 6 5
4
2 3 4 1
2
2 1
2
2 1
5
4 2 3 5 1
4
3 1 4 2
4
1 2 3 4
6
4 2 1 5 6 3
5
4 5 3 2 1
5
3 4 5 2 1
5
2 1 3 4 5
6
5 2 6 3 1 4
4
1 2 3 4
5
3 4 5 2 1
6
3 6 5 2 4 1
2
1 2
4
4 2 1 3
2
1 2
2
1 2
6
3 4 2 6 1 5
3
3 1 2
6
4 2 3 6 5 1
3
3 2 1
2
2 1
2
1 2
4
2 4 1 3
6
4 3 1 2 6 5
6
2 6 4 3 5 1
2
2 1
3
2 1 3
5
1 3 2 4 5
6
5 2 1 6 3 4
6
1 3 5 6 4 2
3
1 2 3
3
2 3 1
5
2 1 3 5 4`

const mod = 1000000007
const inv2 = 500000004
const inv4 = 250000002

type bit struct {
	n    int
	tree []int
}

func newBIT(n int) *bit {
	return &bit{n: n, tree: make([]int, n+1)}
}

func (b *bit) update(i, v int) {
	for ; i <= b.n; i += i & -i {
		b.tree[i] += v
	}
}

func (b *bit) sum(i int) int {
	s := 0
	for ; i > 0; i -= i & -i {
		s += b.tree[i]
	}
	return s
}

type testCase struct {
	n int
	p []int
}

// Embedded solver from 396D.go.
func solve(tc testCase) int64 {
	n := tc.n
	p := tc.p

	fact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}

	bit := newBIT(n)
	var sumYAll int64
	var ans, prefixInv int64

	for i := 1; i <= n; i++ {
		pi := p[i-1]
		usedLess := bit.sum(pi - 1)
		c := int64(pi - 1 - usedLess)
		U := int64(i - 1)
		rem := n - i
		factRem := fact[rem]

		termA := c * prefixInv % mod * factRem % mod
		termB := factRem * (c * (c - 1) % mod * inv2 % mod) % mod

		var suffixSum int64
		if rem >= 2 {
			suffixSum = factRem * int64(rem) % mod * int64(rem-1) % mod * inv4 % mod
		}
		suffixTotal := c * suffixSum % mod

		uC2 := U * (U - 1) % mod * inv2 % mod
		temp := (sumYAll - U - uC2) % mod
		if temp < 0 {
			temp += mod
		}
		missing := factRem * c % mod * temp % mod
		ans = (ans + termA + termB + suffixTotal + missing) % mod

		prefixInv = (prefixInv + int64((i-1)-usedLess)) % mod
		bit.update(pi, 1)
		sumYAll = (sumYAll + int64(pi)) % mod
	}
	ans = (ans + prefixInv) % mod
	return ans
}

func parseCases() ([]testCase, error) {
	scan := strings.Fields(testcasesD)
	if len(scan) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	pos := 0
	readInt := func() (int, error) {
		if pos >= len(scan) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(scan[pos])
		pos++
		return v, err
	}
	t, err := readInt()
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := readInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n", i+1)
		}
		p := make([]int, n)
		for j := 0; j < n; j++ {
			val, err := readInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: bad p[%d]", i+1, j+1)
			}
			p[j] = val
		}
		cases = append(cases, testCase{n: n, p: p})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Println("failed to load testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc)
		input := buildInput(tc)
		gotStr, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Printf("case %d: bad output\n", idx+1)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d cases passed\n", len(cases))
}

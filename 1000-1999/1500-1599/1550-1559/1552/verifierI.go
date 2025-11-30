package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD int64 = 998244353

type Set struct {
	bits     [2]uint64
	size     int
	children []*Set
}

type testCase struct {
	n    int
	sets [][]int
}

// Embedded testcases from testcasesI.txt.
const testcaseData = `
3 4 2 1 2 2 2 2 2 3 1 1 3
4 4 4 3 2 4 4 2 1 2 2 3 4 3 1 4 1
1 4 1 1 1 1 1 1 1 1
3 5 2 1 2 3 2 2 3 3 2 1 1 1 1 2 2 3
4 1 3 4 2 4
2 1 1 2
4 1 3 4 4 3
4 1 3 1 3 2
5 1 2 3 4
2 3 1 1 2 2 2 2 1 2
3 5 3 1 3 1 3 1 1 3 1 3 2 2 3 2 2 3
3 1 3 3 2 1
1 2 1 1 1 1
3 2 1 1 1 3
2 2 1 2 2 1 1
2 1 2 2 2
3 2 2 3 3 2 1 3
3 1 3 1 1 1
2 2 2 1 1 1 1
4 1 1 1
4 4 1 2 1 1 1 4 1 3
1 1 1 1
1 1 1 1
2 3 1 1 1 2 2 2 1
5 5 2 5 5 1 1 5 1 5 2 2 3 2 5 1 1 2
5 4 3 5 2 4 1 2 3 4 4 2 5 2 4 1 3 1
4 1 2 4 3
3 5 3 1 1 2 2 1 2 2 2 3 3 2 3 2 1 3
4 3 4 2 1 4 3 2 1 4 1 4
4 5 1 1 2 4 4 1 2 1 3 2 3 2
5 4 1 2 4 4 3 5 5 3 4 3 4 5 4 4 5 3 2
4 1 1 3
4 4 2 2 1 1 1 1 4 1 2
4 5 3 2 3 1 1 4 4 2 3 4 4 3 3 1 1 4 1 4 4 3
4 4 4 1 1 2 3 2 4 4 2 3 4 4 1 4 4 4
5 2 5 1 4 3 4 5 4 5 3 3 3
1 1 1 1
3 3 3 3 1 3 2 2 1 2 2 2
4 3 1 2 2 2 3 1 1
3 1 2 1 1
5 3 1 4 4 5 5 4 3 1 4
5 1 3 4 5 4
3 5 1 2 1 3 3 1 2 2 3 3 2 1 3 2 2 2
2 4 1 1 2 2 1 2 2 1 1 1
5 3 3 1 5 5 4 4 2 3 1 3 5 3 4
5 3 2 4 4 3 5 5 4 5 2 1 5 2 5
4 2 4 3 4 1 2 3 4 2 2
2 3 2 2 2 1 1 2 2 2
2 3 1 2 1 1 1 1
1 2 1 1 1 1
4 1 2 2 3
2 1 2 2 1
1 4 1 1 1 1 1 1 1 1
1 4 1 1 1 1 1 1 1 1
5 1 3 3 5 4
2 4 1 1 1 1 1 2 2 2 2
5 1 4 5 3 1 3
5 4 4 2 3 2 3 4 3 3 1 5 1 1 2 3 5
4 3 3 1 2 4 4 1 4 3 3 1 2
3 2 1 2 3 2 3 2
1 3 1 1 1 1 1 1
1 2 1 1 1 1
5 5 3 4 1 1 2 4 2 1 4 2 3 3 5 1 1 2 2 5
2 5 2 2 2 2 1 2 1 2 2 1 1 2 2 1
4 4 2 1 1 2 1 4 2 3 2 1 4
5 4 5 3 2 1 3 2 5 1 1 5 4 4 4 1 2 2 2 3 3 4 4
2 4 1 1 2 1 2 1 2 1 1
3 3 1 3 2 1 2 3 3 3 2
1 1 1 1
3 4 3 1 3 2 3 2 3 3 2 1 1 3 1 1 2
2 4 2 1 1 2 2 2 2 2 1 1 1
3 2 3 2 2 2 1 1
3 5 2 1 2 1 2 1 1 2 2 3 2 2 2
3 5 2 3 2 2 1 2 3 2 1 3 3 3 2 1 1 1
1 2 1 1 1 1
2 4 2 1 2 1 2 2 2 2 2 1 1
1 3 1 1 1 1 1 1
5 5 2 5 1 5 5 2 4 1 2 4 1 3 1 3 1 4 3 5 1 3
2 1 1 2
4 4 3 1 3 1 4 3 2 2 1 3 3 1 1 4 1 4 3 2
2 4 1 1 1 1 2 2 2 1 1
3 4 2 3 2 3 3 1 2 3 1 2 3 1 2
4 4 4 3 4 1 4 3 2 4 4 1 1 3 4 2 2
4 5 1 3 4 4 1 4 3 2 2 3 2 4 3 1 1
1 5 1 1 1 1 1 1 1 1 1 1
1 4 1 1 1 1 1 1 1 1
5 4 4 5 3 3 2 2 2 1 5 3 1 4 5 1 1 4
1 1 1 1
4 4 4 3 3 4 3 1 2 2 1 2 1 1
2 1 1 2
2 2 2 1 1 2 2 1
4 4 1 1 2 2 4 4 4 1 4 2 2 3 4
1 3 1 1 1 1 1 1
1 4 1 1 1 1 1 1 1 1
5 1 3 4 4 1
4 5 3 4 4 1 2 3 2 4 3 4 2 2 3 4 2 3 4 4 3 4 2
5 2 3 1 4 2 2 2 2
5 4 1 4 4 3 2 2 5 1 1 3 3 3 2
2 5 1 1 2 2 1 1 1 1 2 2 2 2
1 5 1 1 1 1 1 1 1 1 1 1
`

func newSet(nums []int) *Set {
	s := &Set{}
	for _, v := range nums {
		v--
		if v < 64 {
			s.bits[0] |= 1 << uint(v)
		} else {
			s.bits[1] |= 1 << uint(v-64)
		}
	}
	s.size = bits.OnesCount64(s.bits[0]) + bits.OnesCount64(s.bits[1])
	return s
}

func newUniversal(n int) *Set {
	nums := make([]int, n)
	for i := 1; i <= n; i++ {
		nums[i-1] = i
	}
	return newSet(nums)
}

func (s *Set) subsetOf(t *Set) bool {
	return s.bits[0]&^t.bits[0] == 0 && s.bits[1]&^t.bits[1] == 0
}

func (s *Set) intersects(t *Set) bool {
	return (s.bits[0]&t.bits[0]) != 0 || (s.bits[1]&t.bits[1]) != 0
}

var fac []int64

func calc(node *Set) int64 {
	sum := 0
	res := int64(1)
	for _, ch := range node.children {
		res = res * calc(ch) % MOD
		sum += ch.size
	}
	blocks := len(node.children) + node.size - sum
	if blocks < 0 {
		return 0
	}
	res = res * fac[blocks] % MOD
	return res
}

func solve(tc testCase) string {
	n := tc.n
	mp := make(map[[2]uint64]bool)
	sets := []*Set{}
	for _, nums := range tc.sets {
		st := newSet(nums)
		key := [2]uint64{st.bits[0], st.bits[1]}
		if !mp[key] {
			mp[key] = true
			sets = append(sets, st)
		}
	}

	for i := 0; i < len(sets); i++ {
		for j := i + 1; j < len(sets); j++ {
			a := sets[i]
			b := sets[j]
			if a.intersects(b) && !a.subsetOf(b) && !b.subsetOf(a) {
				return "0"
			}
		}
	}

	root := newUniversal(n)
	all := append([]*Set{root}, sets...)

	for _, child := range sets {
		var parent *Set
		parentSize := n + 1
		for _, p := range all {
			if p.size > child.size && child.subsetOf(p) && p.size < parentSize {
				parent = p
				parentSize = p.size
			}
		}
		if parent == nil {
			parent = root
		}
		parent.children = append(parent.children, child)
	}

	fac = make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}

	ans := calc(root)
	return strconv.FormatInt(ans, 10)
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("case %d too short", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", idx+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("case %d bad m: %v", idx+1, err)
		}
		pos := 2
		sets := make([][]int, 0, m)
		for i := 0; i < m; i++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("case %d missing q for set %d", idx+1, i+1)
			}
			q, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d bad q %d: %v", idx+1, i+1, err)
			}
			pos++
			if pos+q > len(fields) {
				return nil, fmt.Errorf("case %d not enough elements for set %d", idx+1, i+1)
			}
			nums := make([]int, q)
			for j := 0; j < q; j++ {
				val, err := strconv.Atoi(fields[pos+j])
				if err != nil {
					return nil, fmt.Errorf("case %d bad value %d in set %d: %v", idx+1, j+1, i+1, err)
				}
				nums[j] = val
			}
			sets = append(sets, nums)
			pos += q
		}
		if pos != len(fields) {
			return nil, fmt.Errorf("case %d has leftover data", idx+1)
		}
		cases = append(cases, testCase{n: n, sets: sets})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.sets)))
	for _, st := range tc.sets {
		sb.WriteString(strconv.Itoa(len(st)))
		for _, v := range st {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
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
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

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

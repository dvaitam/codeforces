package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// --- embedded correct solver (treap-based, matches CF-accepted solution) ---

type sItem struct {
	c int64
	q int64
}

type sNode struct {
	val int64
	sum int64
	pr  uint64
	l   *sNode
	r   *sNode
}

func sSum(t *sNode) int64 {
	if t == nil {
		return 0
	}
	return t.sum
}

func sPull(t *sNode) {
	if t != nil {
		t.sum = t.val + sSum(t.l) + sSum(t.r)
	}
}

var sSeed uint64

func sRnd() uint64 {
	sSeed ^= sSeed << 7
	sSeed ^= sSeed >> 9
	return sSeed
}

func sMerge(a, b *sNode) *sNode {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.pr > b.pr {
		a.r = sMerge(a.r, b)
		sPull(a)
		return a
	}
	b.l = sMerge(a, b.l)
	sPull(b)
	return b
}

func sSplitBySum(t *sNode, x int64) (*sNode, *sNode) {
	if t == nil {
		return nil, nil
	}
	ls := sSum(t.l)
	if x <= ls {
		a, b := sSplitBySum(t.l, x)
		t.l = b
		sPull(t)
		return a, t
	}
	if x <= ls+t.val {
		left := t.l
		t.l = nil
		sPull(t)
		return left, t
	}
	a, b := sSplitBySum(t.r, x-ls-t.val)
	t.r = a
	sPull(t)
	return t, b
}

func sCollect(t *sNode, res *[]int64) {
	if t == nil {
		return
	}
	sCollect(t.l, res)
	*res = append(*res, t.val)
	sCollect(t.r, res)
}

func sUpperBound(a []int64, x int64) int {
	l, r := 0, len(a)
	for l < r {
		m := (l + r) >> 1
		if a[m] <= x {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

func solve(shirts []sItem, budgets []int64) []int {
	sSeed = 88172645463393265

	sorted := make([]sItem, len(shirts))
	copy(sorted, shirts)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].q != sorted[j].q {
			return sorted[i].q < sorted[j].q
		}
		return sorted[i].c > sorted[j].c
	})

	var root *sNode
	for i := 0; i < len(sorted); i++ {
		left, right := sSplitBySum(root, sorted[i].c)
		node := &sNode{val: sorted[i].c, sum: sorted[i].c, pr: sRnd()}
		root = sMerge(sMerge(left, node), right)
	}

	seq := make([]int64, 0, len(sorted))
	sCollect(root, &seq)

	pref := make([]int64, len(sorted))
	var s int64
	for i := 0; i < len(sorted); i++ {
		s += seq[i]
		pref[i] = s
	}

	res := make([]int, len(budgets))
	for i, b := range budgets {
		res[i] = sUpperBound(pref, b)
	}
	return res
}

// --- verifier infrastructure ---

type testCase struct {
	input    string
	expected string
}

func run(bin, input string) (string, error) {
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

func genRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 1
	shirts := make([]sItem, n)
	for i := range shirts {
		shirts[i] = sItem{c: int64(rng.Intn(10) + 1), q: int64(rng.Intn(10) + 1)}
	}
	k := rng.Intn(6) + 1
	budgets := make([]int64, k)
	for i := range budgets {
		budgets[i] = int64(rng.Intn(50) + 1)
	}
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i, s := range shirts {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d %d", s.c, s.q))
		if i != n-1 {
			input.WriteByte(' ')
		}
	}
	input.WriteByte('\n')
	input.WriteString(fmt.Sprintf("%d\n", k))
	for i, b := range budgets {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.FormatInt(b, 10))
	}
	input.WriteByte('\n')
	answers := solve(shirts, budgets)
	var expected strings.Builder
	for i, v := range answers {
		if i > 0 {
			expected.WriteByte(' ')
		}
		expected.WriteString(strconv.Itoa(v))
	}
	expected.WriteByte('\n')
	return testCase{input: input.String(), expected: strings.TrimSpace(expected.String())}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, genRandomCase(rand.New(rand.NewSource(42))))
	for i := 0; i < 100; i++ {
		cases = append(cases, genRandomCase(rng))
	}

	for idx, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

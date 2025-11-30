package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesF = `100
4 4 42 35 38 14
2 1 38 34
3 2 3 2 17
4 4 31 48 34 34
2 1 1 26
5 6 35 47 48 20 26
4 3 37 44 33 37
5 7 19 20 12 30 24
3 3 7 10 17
2 1 46 0
3 2 14 49 44
5 6 9 41 22 8 16
4 2 37 25 34 36
4 4 40 30 26 40
4 5 7 23 43 39
5 10 15 42 50 31 11
4 1 37 25 19 24
5 4 22 37 44 15 42
3 1 6 45 27
2 1 28 33
3 1 18 46 26
3 1 17 13 48
2 1 38 5
4 5 30 12 19 7
2 1 32 27
3 3 23 9 21
4 6 26 26 35 28
4 4 35 43 3 1
5 6 43 23 41 34 32
4 2 35 13 37 45
4 5 45 44 44 16
2 1 23 36
2 1 37 9
3 3 32 48 33
5 8 5 35 20 34 19
2 1 39 36
2 1 13 6
2 1 30 44
2 1 33 23
3 2 46 44 43
2 1 0 1
2 1 37 2
5 5 16 32 47 16 3
2 1 41 42
2 1 29 46
4 6 50 28 39 31
4 4 12 42 22 39
4 6 0 37 35 1
4 2 36 15 16 36
2 1 10 34
3 3 19 29 47
4 3 12 14 48 11
2 1 7 18
5 5 17 21 17 13 46
2 1 4 46
4 1 32 34 1 0
4 6 26 16 13 35
3 3 35 16 13
2 1 34 40
2 1 13 4
5 2 10 32 7 36 0
2 1 23 39
3 3 4 49 21
5 7 23 39 1 16 46
3 3 32 5 25
2 1 40 33
5 8 12 45 35 0 15
2 1 30 26
5 9 14 38 42 8 15
2 1 45 34
4 5 29 29 13 22
2 1 5 3
2 1 21 1
2 1 15 48
4 2 16 48 2 49
4 4 18 12 39 49
2 1 4 41
4 3 6 45 0 32
2 1 2 38
2 1 21 8
2 1 10 34
4 2 24 1 28 36
2 1 11 17
3 1 7 5 32
3 3 34 44 4
5 1 40 30 20 35 33
2 1 38 37
2 1 29 13
2 1 30 5
5 5 10 14 6 1 4
2 1 48 48
3 2 3 44 32
4 4 46 43 16 49
4 4 43 41 32 29
5 8 22 8 12 49 1
5 1 23 30 32 11 44
3 3 18 13 46
4 1 47 48 6 46
5 5 22 29 50 6 47
4 1 13 22 6 46
`

type node struct {
	child [2]*node
}

func (n *node) insert(x int) {
	cur := n
	for i := 30; i >= 0; i-- {
		b := (x >> i) & 1
		if cur.child[b] == nil {
			cur.child[b] = &node{}
		}
		cur = cur.child[b]
	}
}

func (n *node) minXor(x int) int {
	cur := n
	res := 0
	for i := 30; i >= 0; i-- {
		b := (x >> i) & 1
		if cur.child[b] != nil {
			cur = cur.child[b]
		} else {
			res |= 1 << i
			cur = cur.child[1-b]
		}
	}
	return res
}

func expected(n, k int, a []int) string {
	var vals []int
	for l := 0; l < n; l++ {
		root := &node{}
		root.insert(a[l])
		minVal := int(^uint(0) >> 1)
		for r := l + 1; r < n; r++ {
			v := root.minXor(a[r])
			if v < minVal {
				minVal = v
			}
			root.insert(a[r])
			vals = append(vals, minVal)
		}
	}
	sort.Ints(vals)
	if k <= len(vals) {
		return fmt.Sprintf("%d", vals[k-1])
	}
	return "-1"
}

type testCase struct {
	n int
	k int
	a []int
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesF)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := nextInt()
	if err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, err
		}
		k, err := nextInt()
		if err != nil {
			return nil, err
		}
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j], err = nextInt()
			if err != nil {
				return nil, err
			}
		}
		tests[i] = testCase{n: n, k: k, a: a}
	}
	return tests, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outFields := strings.Fields(output)
	if len(outFields) != len(tests) {
		fmt.Printf("expected %d outputs, got %d\n", len(tests), len(outFields))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := expected(tc.n, tc.k, tc.a)
		if strings.TrimSpace(outFields[i]) != want {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, want, outFields[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

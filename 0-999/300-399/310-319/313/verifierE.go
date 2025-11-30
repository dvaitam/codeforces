package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Testcases embedded from testcasesE.txt (count + cases).
const rawTestcases = `100
1 3
0
0
4 4
3 3 0 2
0 2 3 0
6 2
0 0 1 1 0 0
0 1 1 1 0 1
7 8
3 3 1 1 6 2 7
7 2 0 1 0 7 2
4 4
0 0 2 1
2 2 3 0
7 7
0 5 2 4 5 2 3
4 2 3 5 2 0 6
1 5
1
0
1 3
1
0
2 7
4 1
5 1
5 5
2 3 3 4 4
3 1 1 0 1
6 7
1 2 6 1 0 4
3 3 4 2 1 1
5 2
0 0 1 0 0
0 0 0 0 0
3 5
4 2 0
3 0 1
1 8
3
2
1 6
3
2
5 6
4 1 4 4 3
2 0 1 0 4
3 5
1 0 4
0 2 1
2 4
2 3
3 1
5 7
4 4 3 3 2
3 1 4 4 3
2 5
3 1
0 2
4 7
6 5 0 4
5 0 0 5
7 4
1 1 2 3 0 0 3
3 2 2 3 2 3 1
3 2
0 0 0
1 1 1
7 8
3 6 0 5 6 7 6
4 1 2 4 2 6 4
3 7
4 0 2
3 2 6
7 6
0 3 1 4 2 1 0
4 5 4 2 0 3 5
6 7
3 2 5 4 5 6
3 6 1 2 0 0
2 6
4 5
5 5
7 2
1 1 0 1 1 1 1
0 1 1 0 1 1 1
5 3
1 2 1 0 2
0 0 1 1 2
3 8
4 2 4
0 0 3
6 4
0 3 0 2 1 3
1 1 2 0 0 2
8 5
2 3 2 4 4 1 0 0
1 3 2 0 4 4 2 4
3 7
3 3 1
2 2 2
6 3
1 0 0 1 1 1
1 1 2 0 2 0
4 3
1 0 2 1
1 0 0 0
6 5
2 4 2 2 1 0
4 2 1 0 3 3
2 5
3 0
0 3
7 3
1 0 2 1 1 0 2
1 2 1 2 1 1 1
8 4
2 0 0 0 0 1 3 3
0 3 0 1 1 3 2 1
6 3
0 1 1 0 0 1
0 1 0 2 0 1
1 8
1
4
5 8
4 2 3 6 7
0 4 6 1 4
6 6
4 3 2 2 0 3
5 1 2 3 0 1
1 2
0
1
7 3
1 0 0 0 1 1 1
2 0 1 1 1 1 1
7 6
4 2 4 3 3 1 5
5 1 5 3 2 1 0
3 3
1 1 1
0 1 0
1 3
0
2
6 6
0 2 4 5 3 3
0 1 2 5 5 1
4 6
2 1 0 0
0 1 0 5
1 5
0
2
4 6
2 5 5 3
1 1 5 4
3 5
3 1 0
4 4 4
7 6
5 3 1 3 4 4 4
5 0 1 4 1 2 4
3 6
3 2 5
4 3 5
7 4
1 3 3 3 1 3 0
0 3 0 2 2 2 3
5 8
1 0 1 4 2
4 3 1 1 3
1 4
3
2
2 6
4 0
0 5
3 7
1 0 3
1 5 6
6 2
0 0 1 0 0 1
0 1 1 0 0 0
5 5
4 0 3 1 3
3 4 1 3 0
1 6
1
4
6 2
0 0 0 0 0 1
1 0 0 1 0 0
4 5
4 3 2 2
1 2 2 3
6 3
0 2 0 1 2 2
0 1 1 0 0 0
6 7
2 5 4 6 5 0
2 1 6 3 1 5
1 7
1
4
4 4
3 3 1 1
1 3 1 0
1 5
1
0
4 7
4 3 1 3
0 6 3 4
2 4
3 0
2 3
1 5
1
1
6 5
2 0 0 1 4 0
0 3 2 1 0 2
4 4
3 3 1 1
3 0 3 3
1 5
1
0
4 3
1 2 2 0
0 1 2 1
7 3
0 0 0 0 1 2 0
0 0 1 1 1 0 2
8 6
4 5 5 2 5 0 5 2
2 0 2 0 5 0 3 2
1 5
4
2
3 3
2 0 2
2 2 0
5 7
3 6 3 1 1
2 6 1 5 0
4 3
0 1 1 1
1 2 0 0
4 4
1 0 0 2
1 0 3 2
6 6
0 5 5 3 2 2
3 2 3 2 2 1
5 7
5 3 6 5 6
1 0 0 2 6
6 8
2 7 4 6 3 0
2 4 2 6 4 4
1 7
2
0
6 7
2 6 4 5 6 2
5 4 1 5 5 6
4 4
0 1 2 1
2 2 1 3
6 8
6 5 5 6 6 4
3 1 4 1 5 7
3 6
3 4 0
2 1 1
2 4
2 0
2 0
3 4
3 0 1
3 1 2
1 2
1
1
6 2
0 1 0 1 1 1
0 1 0 0 0 1
2 8
1 4
1 5
3 5
0 0 1
0 1 1
4 6
2 0 4 1
5 5 2 5`

type testCase struct {
	n, m int
	A, B []int
}

func loadTestcases() ([]testCase, error) {
	reader := strings.NewReader(rawTestcases)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("read count: %w", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		var n, m int
		if _, err := fmt.Fscan(reader, &n, &m); err != nil {
			return nil, fmt.Errorf("case %d header: %w", i+1, err)
		}
		A := make([]int, n)
		B := make([]int, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(reader, &A[j]); err != nil {
				return nil, fmt.Errorf("case %d A[%d]: %w", i+1, j, err)
			}
		}
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(reader, &B[j]); err != nil {
				return nil, fmt.Errorf("case %d B[%d]: %w", i+1, j, err)
			}
		}
		cases = append(cases, testCase{n: n, m: m, A: A, B: B})
	}
	return cases, nil
}

type BIT struct {
	n    int
	tree []int
}

func newBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+1)}
}

func (b *BIT) update(i, v int) {
	for ; i <= b.n; i += i & -i {
		b.tree[i] += v
	}
}

func (b *BIT) sum(i int) int {
	res := 0
	for ; i > 0; i -= i & -i {
		res += b.tree[i]
	}
	return res
}

func (b *BIT) findByPrefix(k int) int {
	idx := 0
	bitMask := 1
	for bitMask<<1 <= b.n {
		bitMask <<= 1
	}
	for bitMask > 0 {
		next := idx + bitMask
		if next <= b.n && b.tree[next] < k {
			idx = next
			k -= b.tree[next]
		}
		bitMask >>= 1
	}
	return idx + 1
}

func solveCase(n, m int, A, B []int) string {
	a := append([]int(nil), A...)
	sort.Ints(a)
	bit := newBIT(m)
	for _, v := range B {
		bit.update(v+1, 1)
	}
	C := make([]int, n)
	for i, aVal := range a {
		bound := m - 1 - aVal
		var bIdx int
		if bound >= 0 {
			cnt := bit.sum(min(bound, m-1) + 1)
			if cnt > 0 {
				idx := bit.findByPrefix(cnt)
				bIdx = idx - 1
			} else {
				total := bit.sum(m)
				idx := bit.findByPrefix(total)
				bIdx = idx - 1
			}
		} else {
			total := bit.sum(m)
			idx := bit.findByPrefix(total)
			bIdx = idx - 1
		}
		bit.update(bIdx+1, -1)
		sum := aVal + bIdx
		if sum >= m {
			sum -= m
		}
		C[i] = sum
	}
	sort.Sort(sort.Reverse(sort.IntSlice(C)))
	var sb strings.Builder
	for i, v := range C {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := loadTestcases()
	if err != nil {
		fmt.Printf("failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range testcases {
		expected := solveCase(tc.n, tc.m, tc.A, tc.B)
		input := fmt.Sprintf("%d %d\n%s\n%s\n", tc.n, tc.m, strings.Join(intSliceToStrings(tc.A), " "), strings.Join(intSliceToStrings(tc.B), " "))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

func intSliceToStrings(arr []int) []string {
	res := make([]string, len(arr))
	for i, v := range arr {
		res[i] = strconv.Itoa(v)
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

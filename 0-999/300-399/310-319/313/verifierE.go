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

var rawTestcases = []string{
	"1 3\n0\n0",
	"4 4\n3 3 0 2\n0 2 3 0",
	"6 2\n0 0 1 1 0 0\n0 1 1 1 0 1",
	"7 8\n3 3 1 1 6 2 7\n7 2 0 1 0 7 2",
	"4 4\n0 0 2 1\n2 2 3 0",
	"7 7\n0 5 2 4 5 2 3\n4 2 3 5 2 0 6",
	"1 5\n1\n0",
	"1 3\n1\n0",
	"2 7\n4 1\n5 1",
	"5 5\n2 3 3 4 4\n3 1 1 0 1",
	"6 7\n1 2 6 1 0 4\n3 3 4 2 1 1",
	"5 2\n0 0 1 0 0\n0 0 0 0 0",
	"3 5\n4 2 0\n3 0 1",
	"1 8\n3\n2",
	"1 6\n3\n2",
	"5 6\n4 1 4 4 3\n2 0 1 0 4",
	"3 5\n1 0 4\n0 2 1",
	"2 4\n2 3\n3 1",
	"5 7\n4 4 3 3 2\n3 1 4 4 3",
	"2 5\n3 1\n0 2",
	"4 7\n6 5 0 4\n5 0 0 5",
	"7 4\n1 1 2 3 0 0 3\n3 2 2 3 2 3 1",
	"3 2\n0 0 0\n1 1 1",
	"7 8\n3 6 0 5 6 7 6\n4 1 2 4 2 6 4",
	"3 7\n4 0 2\n3 2 6",
	"7 6\n0 3 1 4 2 1 0\n4 5 4 2 0 3 5",
	"6 7\n3 2 5 4 5 6\n3 6 1 2 0 0",
	"2 6\n4 5\n5 5",
	"7 2\n1 1 0 1 1 1 1\n0 1 1 0 1 1 1",
	"5 3\n1 2 1 0 2\n0 0 1 1 2",
	"3 8\n4 2 4\n0 0 3",
	"6 4\n0 3 0 2 1 3\n1 1 2 0 0 2",
	"8 5\n2 3 2 4 4 1 0 0\n1 3 2 0 4 4 2 4",
	"3 7\n3 3 1\n2 2 2",
	"6 3\n1 0 0 1 1 1\n1 1 2 0 2 0",
	"4 3\n1 0 2 1\n1 0 0 0",
	"6 5\n2 4 2 2 1 0\n4 2 1 0 3 3",
	"2 5\n3 0\n0 3",
	"7 3\n1 0 2 1 1 0 2\n1 2 1 2 1 1 1",
	"8 4\n2 0 0 0 0 1 3 3\n0 3 0 1 1 3 2 1",
	"6 3\n0 1 1 0 0 1\n0 1 0 2 0 1",
	"1 8\n1\n4",
	"5 8\n4 2 3 6 7\n0 4 6 1 4",
	"6 6\n4 3 2 2 0 3\n5 1 2 3 0 1",
	"1 2\n0\n1",
	"7 3\n1 0 0 0 1 1 1\n2 0 1 1 1 1 1",
	"7 6\n4 2 4 3 3 1 5\n5 1 5 3 2 1 0",
	"3 3\n1 1 1\n0 1 0",
	"1 3\n0\n2",
	"6 6\n0 2 4 5 3 3\n0 1 2 5 5 1",
	"4 6\n2 1 0 0\n0 1 0 5",
	"1 5\n0\n2",
	"4 6\n2 5 5 3\n1 1 5 4",
	"3 5\n3 1 0\n4 4 4",
	"7 6\n5 3 1 3 4 4 4\n5 0 1 4 1 2 4",
	"3 6\n3 2 5\n4 3 5",
	"7 4\n1 3 3 3 1 3 0\n0 3 0 2 2 2 3",
	"5 8\n1 0 1 4 2\n4 3 1 1 3",
	"1 4\n3\n2",
	"2 6\n4 0\n0 5",
	"3 7\n1 0 3\n1 5 6",
	"6 2\n0 0 1 0 0 1\n0 1 1 0 0 0",
	"5 5\n4 0 3 1 3\n3 4 1 3 0",
	"1 6\n1\n4",
	"6 2\n0 0 0 0 0 1\n1 0 0 1 0 0",
	"4 5\n4 3 2 2\n1 2 2 3",
	"6 3\n0 2 0 1 2 2\n0 1 1 0 0 0",
	"6 7\n2 5 4 6 5 0\n2 1 6 3 1 5",
	"1 7\n1\n4",
	"4 4\n3 3 1 1\n1 3 1 0",
	"1 5\n1\n0",
	"4 7\n4 3 1 3\n0 6 3 4",
	"2 4\n3 0\n2 3",
	"1 5\n1\n1",
	"6 5\n2 0 0 1 4 0\n0 3 2 1 0 2",
	"4 4\n3 3 1 1\n3 0 3 3",
	"1 5\n1\n0",
	"4 3\n1 2 2 0\n0 1 2 1",
	"7 3\n0 0 0 0 1 2 0\n0 0 1 1 1 0 2",
	"8 6\n4 5 5 2 5 0 5 2\n2 0 2 0 5 0 3 2",
	"1 5\n4\n2",
	"3 3\n2 0 2\n2 2 0",
	"5 7\n3 6 3 1 1\n2 6 1 5 0",
	"4 3\n0 1 1 1\n1 2 0 0",
	"4 4\n1 0 0 2\n1 0 3 2",
	"6 6\n0 5 5 3 2 2\n3 2 3 2 2 1",
	"5 7\n5 3 6 5 6\n1 0 0 2 6",
	"6 8\n2 7 4 6 3 0\n2 4 2 6 4 4",
	"1 7\n2\n0",
	"6 7\n2 6 4 5 6 2\n5 4 1 5 5 6",
	"4 4\n0 1 2 1\n2 2 1 3",
	"6 8\n6 5 5 6 6 4\n3 1 4 1 5 7",
	"3 6\n3 4 0\n2 1 1",
	"2 4\n2 0\n2 0",
	"3 4\n3 0 1\n3 1 2",
	"1 2\n1\n1",
	"6 2\n0 1 0 1 1 1\n0 1 0 0 0 1",
	"2 8\n1 4\n1 5",
	"3 5\n0 0 1\n0 1 1",
	"4 6\n2 0 4 1\n5 5 2 5",
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

func parseCase(raw string) (int, int, []int, []int, error) {
	fields := strings.Fields(raw)
	if len(fields) < 2 {
		return 0, 0, nil, nil, fmt.Errorf("invalid case")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, nil, nil, err
	}
	m, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, nil, nil, err
	}
	if len(fields) != 2+2*n {
		return 0, 0, nil, nil, fmt.Errorf("expected %d numbers got %d", 2+2*n, len(fields))
	}
	A := make([]int, n)
	B := make([]int, n)
	for i := 0; i < n; i++ {
		val, _ := strconv.Atoi(fields[2+i])
		A[i] = val
	}
	for i := 0; i < n; i++ {
		val, _ := strconv.Atoi(fields[2+n+i])
		B[i] = val
	}
	return n, m, A, B, nil
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
	for idx, tc := range rawTestcases {
		n, m, A, B, err := parseCase(tc)
		if err != nil {
			fmt.Printf("case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected := solveCase(n, m, A, B)
		input := fmt.Sprintf("%d %d\n%s\n%s\n", n, m, strings.Join(intSliceToStrings(A), " "), strings.Join(intSliceToStrings(B), " "))
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
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
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

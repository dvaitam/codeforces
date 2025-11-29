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

type testCase struct {
	n int
	a []int
	b []int
}

const solution1038CSource = `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })

	var sumA, sumB int64
	i, j := 0, 0
	turn := 0 // 0: A's turn, 1: B's turn
	for i < n || j < n {
		if turn == 0 {
			if i < n && (j >= n || a[i] > b[j]) {
				sumA += int64(a[i])
				i++
			} else {
				j++
			}
			turn = 1
		} else {
			if j < n && (i >= n || b[j] > a[i]) {
				sumB += int64(b[j])
				j++
			} else {
				i++
			}
			turn = 0
		}
	}
	fmt.Fprint(writer, sumA-sumB)
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1038CSource

var testcases = []testCase{
	{n: 1, a: []int{3}, b: []int{3}},
	{n: 6, a: []int{6, 10, 9, 20, 7, 20}, b: []int{2, 19, 6, 14, 13, 17}},
	{n: 6, a: []int{18, 15, 17, 9, 2, 1}, b: []int{12, 15, 11, 13, 14, 17}},
	{n: 3, a: []int{18, 6, 8}, b: []int{8, 1, 6}},
	{n: 6, a: []int{6, 5, 17, 17, 12, 17}, b: []int{18, 6, 15, 14, 17, 12}},
	{n: 10, a: []int{12, 12, 15, 6, 13, 15, 17, 8, 16, 9}, b: []int{16, 17, 17, 12, 15, 15, 12, 19, 18, 15}},
	{n: 8, a: []int{8, 11, 6, 20, 9, 16, 10, 10}, b: []int{17, 18, 17, 17, 20, 19, 14, 10}},
	{n: 4, a: []int{16, 17, 12, 20}, b: []int{3, 11, 1, 7}},
	{n: 2, a: []int{2, 19}, b: []int{2, 9}},
	{n: 10, a: []int{8, 4, 17, 5, 9, 8, 7, 2, 14, 2}, b: []int{2, 12, 12, 6, 8, 1, 3, 4, 3, 1}},
	{n: 1, a: []int{1}, b: []int{12}},
	{n: 5, a: []int{5, 6, 6, 17, 1}, b: []int{13, 19, 2, 8, 5}},
	{n: 1, a: []int{1}, b: []int{12}},
	{n: 10, a: []int{4, 10, 11, 16, 1, 10, 15, 18, 20, 2}, b: []int{9, 13, 20, 5, 16, 8, 3, 11, 4, 1}},
	{n: 8, a: []int{5, 17, 19, 13, 16, 17, 11, 5}, b: []int{11, 9, 9, 20, 14, 1, 18, 5}},
	{n: 1, a: []int{9}, b: []int{2}},
	{n: 3, a: []int{6, 6, 4}, b: []int{15, 8, 17}},
	{n: 1, a: []int{8}, b: []int{8}},
	{n: 8, a: []int{3, 9, 3, 19, 8, 20, 20, 12}, b: []int{9, 14, 9, 17, 1, 5, 2, 13}},
	{n: 7, a: []int{6, 4, 17, 3, 8, 4, 4}, b: []int{1, 6, 8, 4, 7, 1, 17}},
	{n: 8, a: []int{15, 10, 18, 13, 7, 7, 14, 14}, b: []int{17, 1, 19, 19, 2, 14, 17, 19}},
	{n: 3, a: []int{4, 16, 12}, b: []int{1, 17, 4}},
	{n: 10, a: []int{12, 10, 12, 10, 1, 14, 4, 4, 10, 7}, b: []int{1, 15, 2, 14, 16, 15, 7, 19, 20, 3}},
	{n: 1, a: []int{10}, b: []int{1}},
	{n: 6, a: []int{10, 3, 8, 16, 7, 4}, b: []int{19, 12, 13, 15, 5, 12}},
	{n: 7, a: []int{4, 9, 4, 4, 3, 20, 11}, b: []int{13, 7, 4, 1, 20, 16, 2}},
	{n: 8, a: []int{10, 12, 15, 5, 12, 9, 16, 17}, b: []int{16, 14, 16, 10, 13, 8, 6, 16}},
	{n: 10, a: []int{9, 18, 14, 3, 19, 19, 4, 3, 12, 6}, b: []int{18, 5, 14, 3, 3, 2, 5, 10, 13, 8}},
	{n: 6, a: []int{15, 6, 17, 10, 4, 5}, b: []int{18, 14, 4, 11, 17, 8}},
	{n: 9, a: []int{9, 6, 6, 15, 8, 13, 12, 19, 5}, b: []int{15, 15, 1, 20, 13, 6, 13, 17, 2}},
	{n: 8, a: []int{9, 13, 9, 14, 16, 12, 18, 11}, b: []int{3, 8, 18, 20, 7, 13, 13, 1}},
	{n: 6, a: []int{15, 17, 15, 6, 4, 1}, b: []int{13, 7, 19, 20, 13, 7}},
	{n: 2, a: []int{13, 18}, b: []int{7, 9}},
	{n: 10, a: []int{19, 7, 16, 20, 5, 1, 20, 14, 16, 9}, b: []int{17, 19, 6, 15, 7, 3, 12, 1, 16, 18}},
	{n: 2, a: []int{19, 16}, b: []int{11, 15}},
	{n: 5, a: []int{17, 15, 1, 3, 20}, b: []int{12, 6, 13, 9, 5}},
	{n: 1, a: []int{6}, b: []int{16}},
	{n: 7, a: []int{15, 10, 5, 1, 10, 18, 15}, b: []int{1, 12, 2, 18, 13, 19, 15}},
	{n: 4, a: []int{10, 16, 5, 16}, b: []int{18, 10, 3, 9}},
	{n: 6, a: []int{10, 11, 10, 13, 17, 3}, b: []int{17, 7, 13, 20, 17, 5}},
	{n: 9, a: []int{3, 10, 2, 8, 15, 18, 8, 17, 9}, b: []int{2, 4, 4, 13, 12, 7, 11, 12, 3}},
	{n: 6, a: []int{15, 12, 6, 16, 15, 10}, b: []int{15, 5, 15, 7, 9, 11}},
	{n: 3, a: []int{4, 8, 16}, b: []int{7, 12, 6}},
	{n: 6, a: []int{5, 5, 8, 9, 18, 13}, b: []int{13, 11, 9, 20, 17, 19}},
	{n: 6, a: []int{13, 10, 18, 20, 3, 12}, b: []int{10, 13, 16, 6, 9, 12}},
	{n: 8, a: []int{16, 3, 6, 11, 13, 5, 1, 4}, b: []int{12, 6, 12, 3, 14, 1, 18, 11}},
	{n: 4, a: []int{20, 13, 18, 10}, b: []int{16, 5, 12, 11}},
	{n: 4, a: []int{16, 4, 5, 7}, b: []int{11, 9, 5, 14}},
	{n: 6, a: []int{9, 3, 11, 7, 8, 8}, b: []int{20, 2, 11, 12, 20, 2}},
	{n: 3, a: []int{6, 3, 14}, b: []int{15, 9, 5}},
	{n: 6, a: []int{17, 19, 4, 11, 20, 13}, b: []int{8, 2, 13, 16, 16, 20}},
	{n: 6, a: []int{18, 20, 20, 3, 19, 17}, b: []int{18, 16, 13, 15, 6, 14}},
	{n: 7, a: []int{17, 15, 2, 4, 15, 19, 5}, b: []int{4, 17, 6, 3, 13, 10, 15}},
	{n: 1, a: []int{9}, b: []int{4}},
	{n: 6, a: []int{8, 6, 1, 5, 14, 3}, b: []int{11, 15, 2, 16, 8, 3}},
	{n: 8, a: []int{5, 18, 1, 5, 17, 18, 2, 2}, b: []int{7, 18, 1, 17, 11, 17, 8, 5}},
	{n: 6, a: []int{16, 1, 5, 18, 4, 8}, b: []int{4, 15, 7, 2, 20, 7}},
	{n: 7, a: []int{11, 20, 13, 17, 17, 6, 17}, b: []int{4, 5, 7, 6, 13, 7, 10}},
	{n: 6, a: []int{14, 5, 14, 5, 13, 11}, b: []int{10, 4, 18, 4, 16, 9}},
	{n: 5, a: []int{17, 16, 9, 8, 14}, b: []int{5, 18, 4, 1, 20}},
	{n: 9, a: []int{7, 7, 7, 13, 19, 2, 5, 1, 9}, b: []int{16, 18, 2, 8, 5, 20, 11, 2, 7}},
	{n: 2, a: []int{5, 18}, b: []int{6, 3}},
	{n: 8, a: []int{10, 7, 6, 11, 9, 17, 19, 3}, b: []int{14, 14, 2, 15, 10, 4, 9, 1}},
	{n: 4, a: []int{14, 11, 9, 18}, b: []int{13, 19, 17, 7}},
	{n: 7, a: []int{5, 6, 15, 15, 12, 13, 16}, b: []int{20, 9, 20, 7, 19, 16, 15}},
	{n: 4, a: []int{16, 19, 11, 10}, b: []int{3, 6, 12, 20}},
	{n: 8, a: []int{8, 20, 19, 5, 10, 7, 18, 10}, b: []int{4, 1, 1, 7, 11, 2, 11, 18}},
	{n: 5, a: []int{11, 15, 3, 14, 16}, b: []int{1, 10, 19, 19, 5}},
	{n: 4, a: []int{5, 6, 20, 13}, b: []int{3, 19, 15, 9}},
	{n: 2, a: []int{16, 16}, b: []int{8, 5}},
	{n: 10, a: []int{10, 8, 7, 20, 11, 19, 20, 13, 17, 14}, b: []int{8, 7, 18, 2, 9, 8, 5, 20, 13, 14}},
	{n: 2, a: []int{15, 13}, b: []int{13, 16}},
	{n: 7, a: []int{10, 7, 8, 8, 2, 18, 17}, b: []int{3, 20, 18, 1, 2, 13, 14}},
	{n: 7, a: []int{8, 17, 9, 4, 12, 17, 12}, b: []int{17, 16, 19, 3, 15, 8, 9}},
	{n: 1, a: []int{1}, b: []int{16}},
	{n: 1, a: []int{5}, b: []int{5}},
	{n: 4, a: []int{11, 8, 18, 2}, b: []int{20, 5, 10, 4}},
	{n: 9, a: []int{18, 3, 5, 14, 5, 2, 10, 17, 9}, b: []int{16, 2, 18, 12, 11, 4, 20, 12, 4}},
	{n: 10, a: []int{12, 12, 9, 16, 10, 17, 20, 5, 1, 2}, b: []int{11, 14, 1, 12, 18, 2, 3, 18, 17, 20}},
	{n: 7, a: []int{14, 14, 8, 6, 6, 20, 2}, b: []int{1, 19, 12, 6, 10, 1, 2}},
	{n: 4, a: []int{19, 8, 13, 3}, b: []int{12, 4, 20, 3}},
	{n: 4, a: []int{8, 18, 7, 4}, b: []int{1, 13, 3, 17}},
	{n: 5, a: []int{19, 8, 2, 17, 17}, b: []int{17, 13, 14, 5, 5}},
	{n: 7, a: []int{5, 15, 12, 2, 19, 6, 17}, b: []int{15, 14, 20, 15, 6, 16, 20}},
	{n: 3, a: []int{12, 5, 1}, b: []int{9, 6, 5}},
	{n: 7, a: []int{19, 9, 15, 16, 15, 7, 14}, b: []int{14, 9, 8, 12, 2, 13, 20}},
	{n: 1, a: []int{14}, b: []int{10}},
	{n: 1, a: []int{18}, b: []int{16}},
	{n: 10, a: []int{9, 9, 8, 15, 15, 12, 17, 20, 15, 8}, b: []int{18, 18, 6, 15, 10, 12, 14, 4, 17, 8}},
	{n: 7, a: []int{4, 14, 20, 15, 20, 17, 15}, b: []int{3, 13, 15, 20, 12, 18, 12}},
	{n: 3, a: []int{5, 8, 6}, b: []int{14, 15, 16}},
	{n: 3, a: []int{14, 9, 11}, b: []int{19, 13, 10}},
	{n: 5, a: []int{11, 1, 13, 19, 2}, b: []int{7, 15, 4, 4, 1}},
	{n: 6, a: []int{11, 20, 11, 13, 6, 11}, b: []int{3, 17, 20, 16, 13, 20}},
	{n: 4, a: []int{15, 4, 20, 1}, b: []int{12, 1, 10, 16}},
	{n: 3, a: []int{2, 1, 11}, b: []int{13, 16, 20}},
	{n: 1, a: []int{16}, b: []int{19}},
	{n: 4, a: []int{8, 20, 11, 6}, b: []int{11, 10, 13, 19}},
	{n: 10, a: []int{16, 15, 9, 3, 17, 7, 19, 12, 8, 12}, b: []int{12, 6, 8, 18, 20, 8, 7, 19, 15, 8}},
	{n: 7, a: []int{9, 19, 7, 17, 6, 1, 13}, b: []int{16, 12, 6, 7, 20, 6, 16}},
}

func expected(n int, a, b []int) string {
	sa := append([]int(nil), a...)
	sb := append([]int(nil), b...)
	sort.Slice(sa, func(i, j int) bool { return sa[i] > sa[j] })
	sort.Slice(sb, func(i, j int) bool { return sb[i] > sb[j] })
	var sumA, sumB int64
	i, j := 0, 0
	turn := 0
	for i < n || j < n {
		if turn == 0 {
			if i < n && (j >= n || sa[i] > sb[j]) {
				sumA += int64(sa[i])
				i++
			} else {
				j++
			}
			turn = 1
		} else {
			if j < n && (i >= n || sb[j] > sa[i]) {
				sumB += int64(sb[j])
				j++
			} else {
				i++
			}
			turn = 0
		}
	}
	return fmt.Sprintf("%d", sumA-sumB)
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	for caseIdx, tc := range testcases {
		n := tc.n
		a := append([]int(nil), tc.a...)
		b := append([]int(nil), tc.b...)
		var sbInput strings.Builder
		sbInput.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sbInput.WriteByte(' ')
			}
			sbInput.WriteString(strconv.Itoa(a[i]))
		}
		sbInput.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				sbInput.WriteByte(' ')
			}
			sbInput.WriteString(strconv.Itoa(b[i]))
		}
		sbInput.WriteByte('\n')
		exp := expected(n, a, b)
		if err := runCase(exe, sbInput.String(), exp); err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

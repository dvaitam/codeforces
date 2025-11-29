package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var testcases = []struct{ w []int }{
	{w: []int{30, 3, 12, 32, 24, 10, 2, 13, 21, 34, 30, 25, 30, 25, 22, 9}},
	{w: []int{15, 22, 28, 22, 25, 6, 9, 24, 34, 9, 13, 14}},
	{w: []int{23, 32, 2, 20, 1, 6, 19, 16, 6, 8, 10, 8, 10, 19, 2, 4, 12, 20, 7, 5, 4, 11, 1, 15, 26, 18, 11, 1, 13}},
	{w: []int{17, 22, 20, 14, 13, 24, 13, 26, 12, 4, 19, 26, 18}},
	{w: []int{18, 2, 15, 5, 24, 19, 16, 29, 25, 1, 19, 9, 7, 13, 4, 25, 9, 27, 28, 18, 23, 4, 5, 24, 3, 11}},
	{w: []int{27, 20, 3, 1, 11, 16, 10, 20, 4, 20, 24, 24, 8, 18, 5, 19, 9, 4, 23, 19, 16, 10, 6, 1, 1, 11}},
	{w: []int{20, 9, 1, 17, 21, 19, 5, 7, 30, 20, 4, 15, 29, 14}},
	{w: []int{5, 12, 27, 6, 29, 7, 11, 14, 22, 24, 6, 10, 18, 27, 19, 1, 2, 5, 1, 7, 6, 5}},
	{w: []int{7, 19, 7, 13, 6, 18, 5, 10, 13}},
	{w: []int{17, 15}},
	{w: []int{15, 8, 14, 7, 16, 17, 1, 7, 38, 9, 27, 18}},
	{w: []int{36, 10, 33}},
	{w: []int{1, 26, 13, 36, 21, 34, 19, 20, 24, 10, 35, 12, 26, 4, 6, 4, 16}},
	{w: []int{34, 36, 2, 8, 20, 24, 15, 27, 23, 26, 30, 31, 50, 1, 19}},
	{w: []int{37, 31, 24, 40}},
	{w: []int{20, 20, 34, 50, 20, 14, 9, 11, 3, 46, 15, 25, 34, 22, 36, 9, 29, 11, 2, 36, 34, 26, 12, 35, 7, 8}},
	{w: []int{43, 28, 23, 45, 19, 12, 13, 3, 23, 26, 10, 29, 4, 29, 30, 36, 23, 34, 1, 19, 27, 35, 13, 25, 2, 23, 35, 7, 19, 8, 8, 6, 31, 6, 12, 24, 11, 7, 24, 12, 29, 30, 13, 3, 30}},
	{w: []int{1, 21, 27, 21, 3, 8, 35, 21, 6, 14, 10, 16, 4, 25, 11, 14, 13, 18, 33, 24, 27, 17, 14, 22, 6, 32, 20, 32, 33, 14, 19, 7, 10}},
	{w: []int{21, 1, 11, 20, 32, 5, 8, 2, 48, 29, 30, 12, 7, 37, 45, 34, 35, 23, 10, 12, 48, 11}},
	{w: []int{41, 38, 26, 42, 22, 15, 44, 23, 21, 47, 22, 24, 36, 13, 32, 7, 40, 36, 44, 4, 40, 26, 22, 38, 38, 8, 46, 16, 26, 36, 20, 34}},
	{w: []int{36, 2, 2, 1, 41, 13, 24, 24, 34, 3, 3, 14, 36, 20, 28, 33, 25, 23, 6, 32, 39, 2, 23, 2, 32, 37, 2, 40, 4, 28, 29, 38, 35, 3, 23, 9, 21, 8, 12, 32, 28, 10, 11, 11, 6, 12, 47, 40}},
	{w: []int{24, 25, 43, 3, 12, 1, 32, 11, 10, 35, 4, 44, 20, 47, 10, 22}},
	{w: []int{21, 21, 18, 19, 47, 33, 23, 44, 26, 31, 37, 24, 25}},
	{w: []int{38, 14, 22, 17, 35, 33, 1, 12, 19, 19, 36, 16, 37, 31, 11, 8, 11, 10, 2, 1, 10, 10}},
	{w: []int{10, 1, 4, 27, 10, 18, 16}},
	{w: []int{27, 7, 4, 14, 10, 15, 8, 9, 17, 17, 5}},
	{w: []int{18, 27, 15, 27, 3, 27, 11, 2, 2, 13, 10, 32, 11, 29, 28, 8, 11, 18, 16, 9, 1}},
	{w: []int{15, 29, 5, 11, 13, 25, 27, 26, 15, 19, 28, 9, 20, 9, 7, 17, 3, 13, 7, 4, 4, 10, 24, 16, 9, 19, 20, 12, 11, 10, 21}},
	{w: []int{18, 12, 18, 10, 17, 3, 10, 17, 1}},
	{w: []int{21, 15, 18, 16, 7, 27}},
	{w: []int{8, 17, 16, 17, 20, 20, 25, 5, 12, 2, 2, 15, 10, 18, 18, 9, 12, 27, 4, 4, 21, 10, 6}},
	{w: []int{8, 27, 34, 18}},
	{w: []int{22, 4, 26, 1, 28, 17, 4, 4, 18, 31, 5, 7, 8, 30, 28, 15, 24, 15, 24, 5, 12, 25, 16, 25, 7, 12, 16, 33, 7, 34, 29, 3}},
	{w: []int{12, 11, 8, 1, 21, 27, 17, 12, 16, 27, 7, 17, 23, 13, 8, 27, 10, 15, 28, 26, 15, 22, 18, 1, 25, 18, 26, 15}},
}

const testcasesCount = 34

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveC(w []int) int {
	n := len(w)
	freq := make([]int, n+1)
	for _, v := range w {
		if v >= 1 && v <= n {
			freq[v]++
		}
	}
	best := 0
	for s := 2; s <= 2*n; s++ {
		cnt := 0
		for i := 1; i < s-i && i <= n; i++ {
			j := s - i
			if j > n {
				continue
			}
			cnt += min(freq[i], freq[j])
		}
		if s%2 == 0 {
			i := s / 2
			if i >= 1 && i <= n {
				cnt += freq[i] / 2
			}
		}
		if cnt > best {
			best = cnt
		}
	}
	return best
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	if len(testcases) != testcasesCount {
		fmt.Fprintf(os.Stderr, "unexpected testcase count: got %d want %d\n", len(testcases), testcasesCount)
		os.Exit(1)
	}

	bin := os.Args[1]
	for idx, tc := range testcases {
		var sb strings.Builder
		n := len(tc.w)
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range tc.w {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		expected := strconv.Itoa(solveC(tc.w))
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

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

const testcasesD = `100
5 -1 -14 5 10 -11
3 -16 -19 5
10 -2 -17 -6 13 14 3 -3 -9 -14 -12
5 -19 -4 -3 -8 -10
6 -2 3 -15 1 4 12
5 -9 -5 10 -3 -15
10 -1 -20 -2 16 19 12 -8 6 7 15
6 7 8 -10 -6 -1 -4
2 -15 -18
9 -3 13 14 10 1 -11 -8 -16 6
5 8 -3 -9 2 7
11 0 15 -8 19 -14 -17 -6 -3 -5 14 -10
4 -2 9 -19 -18
7 -15 -2 0 -19 17 18 15
4 6 19 -16 -2
11 -8 8 -2 -12 -4 4 -10 1 -20 -9 -19
9 -10 3 18 -2 -14 8 -7 7 13
3 -17 -10 18
12 -11 18 -18 14 11 -5 0 17 -13 6 -4 19
8 -8 10 19 -5 8 6 11 -18
5 6 8 -5 7 -7
9 -8 -18 18 -4 16 -5 13 -7 -6
8 -4 -11 0 -17 17 -13 5 -18
9 4 -15 7 -7 -10 1 -2 10 0
8 13 -7 -3 1 5 11 -16 17
12 -8 -18 5 -12 -3 -17 -10 9 10 3 19 12
5 -20 -7 -10 19 -4
3 5 4 -6
10 -17 -8 -10 1 15 10 13 8 -19 -18
2 18 -13
9 15 -4 -12 -18 3 -15 13 -20 -1
7 -16 -15 14 9 4 -7 -1
8 -6 11 5 -14 -16 -13 3 12
8 6 8 -16 -8 -1 10 7 -13
10 -10 3 19 -9 -11 0 11 1 -4 -3
2 -10 -20
12 -1 -13 14 18 11 10 13 -16 -5 -7 7 -11
7 -6 -9 -20 -17 0 14 9
11 -1 12 8 17 5 -11 -4 3 1 -16 -7
3 18 -11 19
4 -2 3 -8 16
7 19 -15 -16 5 -9 1 3
7 -9 -1 -19 17 13 -15 2
3 -10 -9 17
9 16 -16 -13 -9 10 -6 -1 5 -5
9 -6 -1 3 19 0 14 13 8 5
10 5 0 -2 8 6 -20 -4 -9 9 2
10 19 3 5 4 -19 -11 12 -16 9 10
12 2 18 -1 -15 -4 10 -6 14 -16 11 9 -13
3 -1 -12 -17
4 5 17 19 15
10 -4 -19 13 -5 -10 -14 -7 -18 0 12
3 -3 -17 -2
11 -9 -11 6 -12 -15 3 -20 18 17 -16 -14
6 10 12 -16 4 -10 15
6 12 5 14 -1 18 1
4 4 -17 7 -19
6 -19 -1 -11 -15 -10 -13
11 -20 -6 18 15 19 10 -9 8 4 -10 7
4 13 17 -8 -14
9 19 2 -1 6 -18 -7 -4 17 10
12 1 -15 -5 0 -14 -18 16 10 2 9 14 18
2 11 13
10 18 -5 -18 -8 -16 1 -12 -1 -13 -4
10 -20 -16 -6 -3 -19 -18 19 10 -11 -14
7 -5 2 -1 4 5 -18 -9
8 11 -16 12 17 -1 -2 -7 7
7 -9 -7 -19 10 7 -3 1
8 6 0 -7 -4 -3 12 -15 -19
8 -4 18 -2 -12 0 -15 -10 14
5 -4 -7 11 -13 1
2 9 16
5 -10 -15 -5 13 -14
3 0 -2 16
11 -11 -16 -9 -19 17 -7 -4 -18 0 11 -10
11 3 -10 18 -7 -17 15 17 14 -4 2 9
12 11 -14 -8 -12 5 10 19 0 3 -7 13 17
2 -18 -9
4 2 -6 -2 -5
7 16 -8 2 4 3 14 -13
5 -8 2 -9 14 -7
2 -19 -13
5 -2 0 14 11 -18
6 5 -10 -12 4 9 2
7 -20 8 -12 18 -19 6 -7
3 1 -2 18
10 -17 -15 10 -10 13 19 -9 -5 1 -1
2 -19 4
8 -6 17 -12 -18 18 -9 -2 19
8 0 14 3 -10 2 18 15 -4
11 1 14 -18 18 -11 3 -8 11 10 6 -17
3 18 -8 -10
4 15 -12 18 -17
10 2 -15 -17 -11 14 -9 -4 17 -2 -5
6 -7 10 -2 13 -17 -12
6 6 -18 -20 2 13 3
3 -8 7 -15
8 -18 -2 -8 17 -9 16 -15 13
2 -20 -2
10 13 19 7 -16 -15 18 1 -10 -8 -4`

// expected computes parent list for the given sequence a using the same
// logic as 675D.go but with deterministic ordering.
func expected(a []int) []int {
	n := len(a)
	if n == 0 {
		return nil
	}
	sorted := []int{a[0]}
	idx := map[int]int{a[0]: 0}
	res := make([]int, 0, n-1)
	for i := 1; i < n; i++ {
		x := a[i]
		pos := lowerBound(sorted, x)
		var pred, succ *int
		if pos > 0 {
			p := sorted[pos-1]
			pred = &p
		}
		if pos < len(sorted) {
			s := sorted[pos]
			succ = &s
		}
		var parent int
		switch {
		case pred == nil:
			parent = *succ
		case succ == nil:
			parent = *pred
		default:
			if idx[*pred] > idx[*succ] {
				parent = *pred
			} else {
				parent = *succ
			}
		}
		res = append(res, parent)
		sorted = append(sorted, 0)
		copy(sorted[pos+1:], sorted[pos:])
		sorted[pos] = x
		idx[x] = i
	}
	return res
}

func lowerBound(a []int, x int) int {
	l, r := 0, len(a)
	for l < r {
		m := (l + r) >> 1
		if a[m] < x {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

func parseTests() ([][]int, error) {
	fields := strings.Fields(testcasesD)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty testcases")
	}
	ptr := 0
	t, err := strconv.Atoi(fields[ptr])
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	ptr++
	tests := make([][]int, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if ptr >= len(fields) {
			return nil, fmt.Errorf("unexpected end at case %d", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[ptr])
		if err != nil {
			return nil, fmt.Errorf("bad n at case %d: %w", caseIdx+1, err)
		}
		ptr++
		if ptr+n > len(fields) {
			return nil, fmt.Errorf("missing numbers at case %d", caseIdx+1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[ptr+i])
			if err != nil {
				return nil, fmt.Errorf("bad value at case %d idx %d: %w", caseIdx+1, i, err)
			}
			arr[i] = v
		}
		ptr += n
		tests = append(tests, arr)
	}
	return tests, nil
}

func runCandidate(bin string, input string) ([]int, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	var res []int
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	for i, arr := range tests {
		expect := expected(arr)
		var input strings.Builder
		input.WriteString(strconv.Itoa(len(arr)))
		input.WriteByte('\n')
		for j, v := range arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if len(got) != len(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\n", i+1, len(expect), len(got))
			os.Exit(1)
		}
		for j := range expect {
			if got[j] != expect[j] {
				fmt.Fprintf(os.Stderr, "case %d failed at pos %d: expected %d got %d\n", i+1, j+1, expect[j], got[j])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

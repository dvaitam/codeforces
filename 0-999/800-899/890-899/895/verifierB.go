package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcases = `
2 10 0 9 4
4 8 3 13 7 4 16
1 7 3 20
1 8 2 8
5 2 2 1 1 1 18 1
4 4 3 1 17 8 15
4 9 1 12 8 8 15
3 1 3 18 4 6
6 5 0 11 17 14 17 7 10
3 10 3 17 13 19
1 8 1 13
4 3 2 18 12 3 15
6 9 0 6 17 13 12 16 1
4 1 2 20 19 19 13
6 3 1 17 8 1 7 18 18
2 7 4 12 19
3 8 2 18 20 1
4 9 1 17 18 7 14
1 8 2 19
5 4 4 14 16 12 14 12
1 9 4 20
5 6 3 20 1 8 6 18
5 3 0 18 9 2 3 3
1 8 0 9
2 5 0 20 6
3 5 0 6 6 9
5 3 5 9 10 15 11 16
4 2 0 10 13 11 14
2 5 0 9 17
2 10 3 1 8
1 7 1 2
6 3 3 17 14 18 8 17 15
2 9 5 1 13
6 10 2 14 2 10 5 7 2
3 2 0 10 10 6
4 10 2 5 1 18 2
5 4 4 15 6 20 17 2
4 4 2 4 7 19 14
5 4 3 4 13 10 17 16
1 6 4 13
3 1 1 7 11 19
2 6 3 7 9
6 2 3 18 12 18 16 18 8
1 1 0 5
2 3 4 7 9
3 10 4 9 12 11
3 2 2 8 20 16
2 10 4 4 11
1 7 0 13
2 3 2 4 20
5 7 0 19 18 8 19 3
3 6 2 19 18 4
4 5 0 2 10 1 20
6 1 0 14 4 2 7 8 19
4 3 0 15 6 8 6
6 2 3 13 18 10 18 9 16
3 2 1 11 2 1
1 5 5 20
3 8 3 11 13 3
1 6 4 15
1 5 1 20
5 8 5 12 9 6 18 7
3 4 1 12 3 9
1 8 0 19
6 6 1 13 10 2 11 6 11
5 5 1 11 4 18 20 19
5 2 1 8 1 8 13 3
3 9 0 3 1 1
3 6 3 16 5 4
5 6 0 17 6 6 5 5
3 5 0 17 20 10
2 4 1 18 2
3 10 5 18 7 6
3 7 4 6 2 8
3 2 5 15 14 18
3 9 3 18 15 1
4 6 1 9 16 1 14
5 1 0 12 19 5 19 5
2 5 2 13 19
4 3 4 3 8 16 1
2 9 2 17 15
6 4 1 11 16 16 8 14 11
5 10 5 9 8 2 3 17
6 6 1 17 7 10 10 10 18
3 3 5 15 20 3
1 10 4 19
4 3 1 9 14 7 19
6 1 3 13 12 13 17 6 18
6 1 4 3 9 4 9 3 5
5 2 3 8 13 14 13 6
3 8 1 20 16 7
1 7 4 18
4 2 5 10 9 8 13
6 9 0 7 17 15 19 1 1
6 10 1 9 7 6 10 5 18
2 5 2 19 9
6 8 1 18 12 16 14 4 7
5 7 1 10 4 1 4 19
6 1 4 10 5 3 17 12 19
3 7 4 12 17 11

`

type testCase struct {
	n int
	x int64
	k int64
	a []int64
}

func parseCases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcases))
	scan.Split(bufio.ScanWords)
	var cases []testCase
	for {
		if !scan.Scan() {
			break
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("parse n: %w", err)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("missing x")
		}
		x, err := strconv.ParseInt(scan.Text(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse x: %w", err)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("missing k")
		}
		k, err := strconv.ParseInt(scan.Text(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse k: %w", err)
		}
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d: missing a[%d]", len(cases)+1, i)
			}
			a[i], err = strconv.ParseInt(scan.Text(), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d: parse a[%d]: %w", len(cases)+1, i, err)
			}
		}
		cases = append(cases, testCase{n: n, x: x, k: k, a: a})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func referenceSolve(tc testCase) string {
	arr := append([]int64(nil), tc.a...)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	ans := int64(0)
	for _, val := range arr {
		t := val / tc.x
		b := t - tc.k
		L := b*tc.x + 1
		R := (b + 1) * tc.x
		if R > val {
			R = val
		}
		if L < 1 {
			L = 1
		}
		if R >= L {
			l := sort.Search(len(arr), func(i int) bool { return arr[i] >= L })
			r := sort.Search(len(arr), func(i int) bool { return arr[i] > R })
			ans += int64(r - l)
		}
	}
	return strconv.FormatInt(ans, 10)
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.x, tc.k)
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatInt(v, 10))
		}
		input.WriteByte('\n')

		expected := referenceSolve(tc)
		out, stderr, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

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

const testcasesData = `
100
2 3 4
6 7 9
1 1 1
4 5 10
0 3 11
8 5 12
2 1 6
3 0 13
3 3 7
2 4 8
10 5 11
9 5 19
4 5 8
2 3 10
4 1 12
4 0 8
9 4 17
2 4 14
9 4 15
7 2 10
3 3 2
1 0 8
10 4 18
8 10 17
3 2 7
1 6 9
10 10 17
4 2 9
4 5 11
10 8 13
5 1 5
3 4 13
9 3 10
3 2 10
7 0 7
5 1 9
10 5 10
3 3 11
2 10 16
9 10 11
4 9 12
4 3 5
4 6 15
2 5 14
0 5 5
7 2 12
3 3 19
1 7 10
6 3 7
0 0 0
2 5 5
9 0 17
7 9 12
3 1 4
8 4 14
10 3 17
3 3 10
4 4 2
3 6 13
3 10 16
3 7 10
0 0 4
3 2 17
3 3 9
4 2 9
0 5 14
1 5 13
10 10 10
4 4 3
6 3 15
2 5 9
10 7 20
5 6 14
2 3 11
6 7 8
4 10 20
2 1 13
9 2 13
10 0 12
4 5 16
6 6 9
0 3 5
0 9 13
1 6 12
2 5 2
3 2 13
9 5 17
4 5 15
0 1 1
9 1 16
8 4 17
2 0 7
1 5 1
4 5 6
1 5 15
6 3 10
4 2 16
6 1 7
1 9 14
8 6 14
7 1 17
2 3 16
6 1 14
2 5 7
2 2 11
7 5 11
8 0 10
1 3 4
8 1 15
5 4 17
1 8 11
6 4 11
3 2 13
0 10 10
9 5 17
4 5 10
8 7 17
9 7 15
2 4 13
5 10 15
2 4 3
9 2 19
9 2 13
3 2 19
5 10 19
1 1 7
10 2 15
3 3 6
4 0 13
0 8 9
5 1 7
2 5 16
10 9 11
1 2 12
4 2 20
10 4 16
9 3 16
3 4 9
3 5 13
5 4 13
8 6 13
4 7 13
9 0 13
2 8 15
8 9 14
4 4 20
0 2 10
1 7 17
5 9 13
1 4 11
2 4 20
1 2 5
1 4 6
1 2 13
9 10 19
5 5 9
0 8 11
2 1 5
0 5 6
1 4 4
10 4 19
2 2 14
2 1 10
5 0 13
2 5 14
2 3 7
7 8 9
4 2 6
4 8 14
8 10 14
6 5 8
6 0 12
1 3 1
4 2 5
2 1 20
0 3 6
8 0 15
5 2 15
10 6 15
2 8 17
3 1 10
9 5 13
10 6 19
1 2 9
9 10 19
3 4 11
1 3 8
1 10 10
9 5 17
4 3 3
2 0 9
8 8 17
3 0 6
10 1 15
5 2 10
1 8 16
0 1 4
9 4 18
0 0 0
7 2 10
3 2 12
4 6 16
9 6 9
2 6 13
1 5 17
4 8 12
2 4 11
2 3 3
7 6 11
3 4 14
5 9 12
4 4 12
1 0 7
3 5 10
10 10 12
9 5 10
2 9 10
2 3 7
10 7 11
10 5 10
4 5 7
2 1 5
8 1 9
5 4 14
9 10 12
1 2 1
2 10 13
3 1 19
5 10 15
9 5 11
2 3 3
1 2 2
10 4 20
4 1 7
2 6 13
9 8 16
8 5 13
6 6 6
1 2 5
5 3 9
2 3 19
3 9 19
5 6 11
5 1 8
3 5 7
8 3 8
0 1 4
4 5 13
7 0 11
4 2 5
6 9 16
5 5 5
7 2 14
0 6 9
1 3 10
9 8 9
1 4 6
10 8 10
2 5 8
9 5 18
0 10 10
4 4 8
9 2 18
0 2 4
9 4 12
6 5 14
3 2 12
10 8 20
9 5 13
9 5 17
1 5 5
5 9 12
4 4 4
1 10 19
3 2 5
8 2 17
0 8 13
1 1 20
2 8 10
3 5 2
4 7 11
3 10 17
4 8 8
2 3 14
0 0 5
8 9 14
1 2 14
10 1 16
1 3 7
10 3 12
2 1 3
0 0 4
8 8 16
4 1 3
8 10 15
10 2 20
3 8 16
0 0 1
2 4 5
1 8 13
8 5 18
3 2 19
9 1 18
10 9 11
9 0 11
1 3 18
10 6 14
5 3 5
8 5 16
7 3 16
7 2 14
10 9 18
0 9 17
5 2 9
9 10 19
9 7 14
2 4 6
5 4 10
3 8 15
5 3 9
5 10 12
5 1 8
5 10 13
2 7 8
4 0 12
1 1 1
10 10 18
2 2 8
5 10 16
6 3 16
4 4 10
8 1 13
0 1 3
0 6 6
6 2 12
1 5 8
5 5 13
3 2 1
5 5 6
6 9 12
7 5 16
4 5 3
6 2 7
7 6 15
3 0 4
4 1 6
2 4 15
5 6 11
6 10 14
4 1 11
6 8 9
7 6 12
3 9 18
2 6 10
1 3 18
0 8 18
2 5 6
7 3 15
5 6 13
4 5 4
7 0 13
5 7 11
3 2 8
5 2 9
5 5 16
9 2 16
6 7 12
9 8 9
10 6 14
9 9 11
5 2 1
7 10 19
4 2 4
0 3 8
2 5 10
7 9 12
3 1 3
0 7 8
10 3 13
7 8 9
`

type station struct {
	a int64
	b int64
	c int64
}

type testCase struct {
	n        int
	t        int
	k        int64
	stations []station
}

func run(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseTestcases() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcasesData))
	scanner.Split(bufio.ScanWords)
	nextInt := func() (int64, error) {
		if !scanner.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return 0, err
		}
		return v, nil
	}

	totalCases64, err := nextInt()
	if err != nil {
		return nil, err
	}
	if totalCases64 < 0 {
		return nil, fmt.Errorf("negative test count")
	}
	totalCases := int(totalCases64)
	cases := make([]testCase, 0, totalCases)
	for i := 0; i < totalCases; i++ {
		n64, err := nextInt()
		if err != nil {
			return nil, err
		}
		tt64, err := nextInt()
		if err != nil {
			return nil, err
		}
		k, err := nextInt()
		if err != nil {
			return nil, err
		}
		if n64 < 0 || tt64 < 0 {
			return nil, fmt.Errorf("negative header at case %d", i+1)
		}
		n := int(n64)
		tc := testCase{
			n:        n,
			t:        int(tt64),
			k:        k,
			stations: make([]station, n),
		}
		for j := 0; j < n; j++ {
			a, err := nextInt()
			if err != nil {
				return nil, err
			}
			b, err := nextInt()
			if err != nil {
				return nil, err
			}
			c, err := nextInt()
			if err != nil {
				return nil, err
			}
			tc.stations[j] = station{a: a, b: b, c: c}
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

// solve mirrors the logic in 1007E.go so we do not need an external oracle.
func solve(tc testCase) int64 {
	n := tc.n
	t := tc.t
	k := tc.k
	a := make([]int64, n)
	b := make([]int64, n)
	c := make([]int64, n)
	for i, s := range tc.stations {
		a[i], b[i], c[i] = s.a, s.b, s.c
	}

	var sumd int64
	for i := 0; i < n; i++ {
		need := a[i] + int64(t)*b[i] - c[i]
		if need > 0 {
			sumd += need
		}
	}
	left := int64(0)
	right := (sumd + k - 1) / k
	right += int64(t)

	feasible := func(M int64) bool {
		used := int64(0)
		B := make([]int64, n)
		copy(B, a)
		for h := 1; h <= t; h++ {
			prefix := make([]int64, n+1)
			for i := 0; i < n; i++ {
				prefix[i+1] = prefix[i] + B[i]
			}
			allZero := true
			var maxCap int64
			for i := 0; i < n; i++ {
				req := B[i] + b[i] - c[i]
				if req < 0 {
					req = 0
				} else {
					allZero = false
				}
				needCap := prefix[i] + req
				if needCap > maxCap {
					maxCap = needCap
				}
			}
			if !allZero {
				R := (maxCap + k - 1) / k
				if used+R > M {
					return false
				}
				used += R
				rem := R * k
				for i := 0; i < n && rem > 0; i++ {
					take := B[i]
					if take > rem {
						take = rem
					}
					B[i] -= take
					rem -= take
				}
			}
			for i := 0; i < n; i++ {
				B[i] += b[i]
			}
		}
		return true
	}

	for left < right {
		mid := (left + right) / 2
		if feasible(mid) {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.t, tc.k))
	for _, s := range tc.stations {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", s.a, s.b, s.c))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		want := strconv.FormatInt(solve(tc), 10)
		input := formatInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if want != strings.TrimSpace(got) {
			fmt.Printf("case %d failed\nexpected:%s\ngot:%s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

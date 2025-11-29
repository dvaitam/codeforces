package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases previously stored in testcasesD.txt.
const testcasesDData = `8 6 24 6 6 7 4 1 7 10 9
10 7 26 10 1 9 2 6 1 7 3 2 3
2 9 10 4 6
8 10 3 4 7 6 5 1 2 8 5
13 4 6 8 1 3 1 3 10 1 2 6 1 4 4 8
9 5 19 6 9 6 8 10 9 1 4 3
7 4 4 3 2 1 10 9 2 10
20 7 15 8 7 1 10 9 6 9 4 3 10 5 9 5 7 10 9 1 8 10 10
8 7 9 9 5 2 6 10 5 9 5
15 8 4 10 5 6 2 7 6 1 2 6 10 2 9 4 10 1
1 1 10 8
4 9 22 1 8 4 5
17 5 26 5 4 7 10 8 8 5 6 1 10 10 10 5 2 6 2 3
13 8 9 6 5 10 4 3 9 8 8 5 8 2 2 3
18 1 25 5 6 6 6 1 7 6 3 5 8 1 9 9 5 6 5 10 6
3 6 27 1 8 10
6 6 29 10 5 5 8 1 8
6 4 12 4 2 1 7 8 7
20 1 11 1 1 8 5 4 3 8 2 7 9 1 5 2 3 1 1 5 7 9 10
10 2 11 8 10 6 5 6 4 1 9 4 8
17 1 22 7 6 1 9 4 2 9 4 1 9 6 3 4 10 8 9 8
14 3 15 5 3 1 6 5 6 5 4 7 1 1 2 3 10
11 10 6 1 1 9 3 10 3 2 3 3 9 3
14 7 1 6 5 6 10 1 4 10 1 6 6 6 2 10 6
12 2 14 4 6 10 5 3 4 9 3 10 3 9 2
11 10 19 9 10 2 10 3 6 5 8 5 8 3
13 7 12 2 6 3 5 3 6 8 2 6 6 5 6 3
14 4 7 3 10 8 10 7 10 4 1 6 3 1 8 7 10
4 1 13 5 10 4 7
1 10 8 7
20 8 27 8 10 7 7 7 1 5 3 3 6 2 9 10 8 3 10 7 4 3 5
1 8 26 6
16 6 17 4 5 10 1 8 3 5 4 6 1 5 10 7 4 1 4
20 7 5 5 10 9 4 3 7 10 8 6 8 10 3 1 10 7 7 6 5 3 5
16 5 13 10 4 8 4 3 8 8 1 7 4 9 5 7 2 5 5
6 3 25 4 5 3 3 6 4
19 10 26 3 1 2 3 5 10 5 5 3 1 5 10 6 2 7 5 6 10 9
17 5 28 2 9 9 1 1 6 7 10 8 6 10 4 2 7 7 2 10
9 9 12 3 1 1 3 7 6 1 6 6
17 10 30 10 7 9 2 2 10 9 4 4 6 9 7 3 1 2 1 6
7 7 19 2 7 4 6 2 2 3
9 4 30 3 4 7 7 4 4 1 9 5
14 4 17 2 10 6 6 5 9 4 8 10 4 10 1 10 2
9 7 21 2 1 2 6 6 2 10 5 2
8 5 26 10 2 2 7 8 3 6 4
5 4 5 1 10 1 8 8
18 9 8 7 4 4 10 2 6 3 5 2 8 8 4 8 5 8 4 8 3
20 4 27 10 7 6 7 5 7 3 3 10 9 8 2 4 10 6 10 5 1 5 8
7 7 8 5 10 6 1 6 1 7
1 10 26 9
1 1 2 7
20 5 15 6 1 9 7 2 7 5 1 6 1 1 6 2 10 1 9 2 10 9 1
7 10 19 7 5 2 2 3 2 3
18 4 8 8 10 5 7 4 1 9 9 8 2 4 3 7 8 2 4 4 2
5 7 28 8 7 5 10 4
18 1 18 7 3 8 10 3 1 5 7 10 8 8 10 7 6 9 4 2 5
18 9 8 6 5 4 9 7 5 5 7 6 6 4 2 5 5 4 2 4 7
18 3 11 4 7 8 4 5 9 6 7 2 10 8 8 2 3 2 9 5 5
1 7 13 4
1 10 6 1
14 7 11 7 7 1 5 8 10 10 1 2 1 3 10 9 4
6 6 14 3 9 5 8 4 2
9 3 20 4 4 3 1 1 8 5 5 9
1 7 15 10
6 1 6 1 9 7 6 2 6
1 8 24 9
8 5 13 1 7 8 8 3 4 2 9
17 2 17 4 9 7 6 5 4 10 1 6 5 2 10 10 3 2 8 1
16 6 9 7 3 1 9 3 3 7 4 5 10 9 10 1 1 4 1
19 6 26 5 8 10 5 5 9 6 6 3 2 8 9 2 2 3 3 6 7 4
18 5 24 6 6 5 9 1 9 2 7 2 10 7 8 2 2 8 3 7 3
14 9 2 6 7 5 5 3 2 10 8 9 6 10 4 9 1
8 8 24 9 7 3 9 9 1 4 5
18 4 9 1 9 6 5 7 3 1 4 10 5 4 8 5 1 1 7 7 9
1 3 22 2
12 7 8 5 7 9 6 10 3 2 6 4 1 7 7
9 6 7 9 4 10 10 10 1 8 4 4
9 3 28 3 8 4 1 2 4 7 7 7
1 3 20 5
10 1 19 4 5 8 7 2 10 6 2 8 6
8 10 19 1 2 6 8 4 5 5 2
1 2 28 10
14 4 16 4 8 9 7 9 8 3 3 1 9 8 2 7 6
6 2 19 10 6 8 10 4 8
14 1 6 5 5 9 5 7 10 1 8 2 8 5 8 9 6
16 10 3 5 4 7 8 6 8 3 8 6 2 5 10 4 10 3 9
10 2 6 1 9 10 9 8 9 8 5 3 1
14 6 10 7 1 1 1 4 7 2 1 7 8 10 4 2 5
1 3 16 10
2 2 1 5 3
15 3 29 4 4 4 9 1 1 8 8 6 3 2 1 10 1 1
8 7 25 5 1 8 5 2 9 3 2
12 6 21 9 4 9 10 10 4 9 9 5 9 4 3
6 6 3 2 5 7 7 9 5
17 3 4 6 7 2 3 1 7 9 1 10 5 10 6 4 1 1 1 1
8 6 4 2 5 5 8 1 2 7 4
9 5 29 10 10 2 3 9 2 6 6 5
12 1 3 7 5 6 5 1 3 8 6 4 4 2 7
12 9 22 2 3 3 10 4 9 2 5 9 7 1 6
7 8 19 10 4 8 9 2 8 4`

type testCase struct {
	n   int
	m   int
	k   int
	arr []int
}

// solve mirrors the logic from 1066D.go for a single test case.
func solve(n, m, k int, arr []int) int {
	check := func(s int) bool {
		sum := 0
		cnt := 0
		for i := s; i < n; i++ {
			sum += arr[i]
			if sum > k {
				cnt++
				sum = arr[i]
			}
		}
		cnt++
		return cnt <= m
	}

	l, r := 0, n-1
	for l < r {
		mid := (l + r) / 2
		if check(mid) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return n - l
}

func parseTestCases(data string) ([]testCase, error) {
	tokens := strings.Fields(data)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no embedded testcases found")
	}
	idx := 0
	cases := []testCase{}
	for idx < len(tokens) {
		if idx+3 >= len(tokens) {
			return nil, fmt.Errorf("missing header data at token %d", idx)
		}
		n, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("invalid n at token %d: %w", idx, err)
		}
		m, err := strconv.Atoi(tokens[idx+1])
		if err != nil {
			return nil, fmt.Errorf("invalid m at token %d: %w", idx+1, err)
		}
		k, err := strconv.Atoi(tokens[idx+2])
		if err != nil {
			return nil, fmt.Errorf("invalid k at token %d: %w", idx+2, err)
		}
		if idx+3+n > len(tokens) {
			return nil, fmt.Errorf("not enough numbers for array starting at token %d", idx)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(tokens[idx+3+i])
			if err != nil {
				return nil, fmt.Errorf("invalid arr value at token %d: %w", idx+3+i, err)
			}
			arr[i] = val
		}
		idx += 3 + n
		cases = append(cases, testCase{n: n, m: m, k: k, arr: arr})
	}
	return cases, nil
}

func runCase(bin string, tc testCase) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	for i, v := range tc.arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := strconv.Itoa(solve(tc.n, tc.m, tc.k, tc.arr))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestCases(testcasesDData)
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

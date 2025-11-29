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

const testcases = `4 10 9 3 6
10 8 10 2 10 1 8 5 9 4 4
8 9 9 8 7 3 4 3 9
7 1 2 3 10 1 5 1
5 8 10 7 7 7
10 8 3 6 2 1 3 8 4 5 7
5 7 9 7 10 6
9 10 7 10 4 6 1 5 10 3
6 9 10 10 2 4 10
5 5 2 2 8 8
2 6 2
7 3 1 5 7 7 2 1
10 10 1 7 10 6 9 5 9 4 1
5 1 2 2 10 9
1 4
7 5 10 5 3 1 6 6
6 3 7 7 8 9 7
10 9 2 10 9 5 7 4 5 7 5
9 5 9 6 1 7 10 6 1 7
10 10 3 1 6 8 6 6 10 5 8
1 10
1 1
6 5 8 5 10 10 6
3 6 3 6
6 10 5 5 7 2 1
10 3 5 9 4 5 4 6 3 7 2
2 10 6
6 4 8 3 2 6 4
10 8 5 4 2 1 9 4 6 10 3
5 6 2 10 6 10
3 7 5 9
5 8 6 7 5 7
10 7 1 7 3 4 1 8 10 9 7
9 4 1 8 9 5 9 6 4 2
10 5 2 4 1 1 9 4 7 10 1
1 8
2 3 9
5 4 1 9 9 7
1 10
2 6 3
5 9 8 1 6 4
4 2 9 2 3
4 5 3 1 8
10 7 1 5 4 5 10 9 9 7 1
8 6 1 1 3 1 2 1 2
8 1 2 9 9 8 6 3 6
2 6 7
7 10 5 6 5 4 6 7
2 3 9
1 7
2 10 3
1 6
8 10 9 7 1 10 7 1 6
8 6 7 7 8 1 4 4 9
5 10 2 7 4 7
3 1 6 6
9 5 2 8 2 9 7 2 6 10
9 2 10 1 8 3 4 7 1 9
2 10 2
7 3 1 6 2 1 2 8
5 10 5 2 1 10
9 9 4 2 9 2 9 1 9 6
10 3 2 4 3 4 8 10 7 5 6
10 7 6 9 7 2 7 9 4 7 3
7 10 10 9 8 3 7 3
3 2 8 8
9 8 10 3 3 5 4 3 10 9
6 4 9 5 7 10 10
10 5 4 5 1 5 8 7 4 3 10
6 4 6 8 3 7 8
10 4 8 10 9 1 8 2 7 1 8
4 4 2 4 5
4 4 5 3 3
10 1 5 3 1 6 3 7 2 2 2
2 5 5
1 6
8 10 6 1 1 6 6 7 7
8 2 4 10 8 7 3 9 6
2 5 2
7 2 8 9 5 2 9 6
6 8 5 5 2 6 10
9 9 2 8 9 6 1 5 10 3
3 3 6 8
2 2 9
3 6 10 7
9 5 3 8 8 5 3 2 2 3
9 9 10 7 6 2 5 5 7 1
3 1 8 9
5 4 9 6 6 7
8 9 2 6 8 2 3 5 10
2 2 10
2 3 4
10 7 7 3 10 10 3 7 4 9 9
3 10 3 4
5 6 5 1 8 7
7 6 9 10 5 8 9 5
8 1 10 4 1 2 4 8 3
9 8 4 4 9 4 1 9 8 2
10 5 3 3 8 2 10 1 1 6 10
4 9 2 8 9`

type testCase struct {
	n    int
	nums []int
}

func parseTestcases() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcases))
	var cases []testCase
	lineNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lineNum++
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", lineNum, err)
		}
		nums := make([]int, 0, n)
		for i, s := range parts[1:] {
			val, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse num %d: %w", lineNum, i+1, err)
			}
			nums = append(nums, val)
		}
		if len(nums) != n {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", lineNum, n, len(nums))
		}
		cases = append(cases, testCase{n: n, nums: nums})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	return cases, nil
}

func referenceSolve(nums []int) string {
	stack := make([]int, 0, len(nums))
	for _, v := range nums {
		m := v % 2
		l := len(stack)
		if l > 0 && stack[l-1] == m {
			stack = stack[:l-1]
		} else {
			stack = append(stack, m)
		}
	}
	if len(stack) > 1 {
		return "NO"
	}
	return "YES"
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
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.nums {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		want := referenceSolve(tc.nums)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, sb.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

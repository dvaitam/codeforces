package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
5 4 2 5 5 2 1 5 0 4
2 1 0 3 3 0 3 5 5 5
3 0 0 3 5 5 3 1 2 1
4 3 3 0 5 2 5 3 4 2
3 3 3 2 4 2 0 0 0 0
1 5 3 3 0 4 3 5 0 4
1 1 5 3 4 5 2 1 5 2
1 1 3 3 1 0 2 1 3 0
5 1 2 0 5 0 2 2 2 1
4 5 1 3 3 5 4 0 2 5
3 4 0 2 2 4 4 3 2 2
0 3 5 1 2 3 0 1 0 0
0 2 3 1 2 0 0 0 2 3
2 5 1 3 2 2 2 3 3 4
4 2 4 3 3 1 5 5 1 5
3 2 1 0 4 1 1 2 2 2
0 2 0 0 1 0 4 5 2 4
0 2 4 5 3 3 0 1 2 5
5 1 1 4 2 1 0 0 0 1
0 5 5 0 3 0 2 1 4 2
5 5 3 1 1 5 4 5 4 4
1 3 3 1 0 4 4 4 3 4
5 3 1 3 4 4 4 5 0 1
4 1 2 4 1 4 3 2 5 4
3 5 3 2 1 3 3 4 3 4
5 1 5 3 0 0 5 5 4 3
4 0 2 2 2 3 5 5 4 2
4 4 0 0 0 5 2 1 2 5
1 5 0 0 4 1 0 2 3 2
0 4 4 0 0 5 1 5 1 0
3 1 5 4 5 4 2 0 1 1
3 4 0 4 1 4 3 1 3 4
2 1 1 0 2 3 4 5 5 0
3 1 3 3 4 1 3 0 4 0
4 1 4 4 2 0 1 0 4 0
1 5 5 1 3 2 1 0 3 1
5 1 5 1 3 4 3 2 2 1
5 5 2 2 5 3 2 1 1 5
1 3 5 5 1 3 2 1 1 0
2 5 2 5 4 5 0 2 1 3
1 5 0 5 1 4 1 2 3 3
1 4 1 1 5 3 4 4 1 5
0 4 0 3 1 0 1 5 4 3
1 3 0 3 4 4 0 2 3 0
2 3 0 3 1 5 1 5 4 2
3 2 5 0 0 1 4 5 5 5
0 0 3 2 1 0 2 1 2 3
4 3 4 1 1 3 0 3 3 0
3 1 0 1 1 2 4 4 1 0
3 4 3 3 1 0 0 0 0 3
4 0 0 1 3 2 2 0 4 5
4 4 5 3 4 4 5 5 2 5
0 5 2 2 0 2 0 5 0 3
2 0 3 4 2 1 1 5 1 4
4 4 1 2 5 3 3 1 1 2
1 5 0 5 1 1 0 2 3 3
2 5 0 1 1 2 5 1 0 0
2 1 4 0 3 2 2 4 0 5
5 3 2 2 3 2 3 2 2 1
4 2 5 5 3 5 1 0 0 2
2 1 4 5 4 3 2 3 1 0
1 4 2 1 3 2 4 5 2 0
5 2 0 4 2 5 2 4 5 2
5 4 1 5 5 4 1 2 0 4
1 2 5 4 5 1 4 2 2 1
3 2 4 3 2 2 5 3 3 2
1 5 4 5 0 2 0 2 3 1
4 3 4 0 2 1 1 0 2 2
4 0 2 4 0 1 2 3 0 4
1 3 5 1 2 0 0 4 2 4
4 5 2 2 0 4 4 0 4 3
1 3 5 5 3 2 1 3 0 1
1 3 5 0 0 5 2 0 2 5
1 3 0 0 1 0 1 1 1 4
2 0 4 1 5 5 2 5 1 4
1 2 0 1 4 5 5 3 1 0
4 2 0 4 2 5 1 1 0 4
5 5 3 5 5 4 5 4 2 4
1 3 5 1 4 2 4 4 1 3
0 1 3 2 3 5 4 3 2 1
5 0 5 3 2 0 0 4 5 5
0 2 4 3 3 4 5 3 0 1
5 0 2 3 2 2 4 4 2 1
2 1 0 5 2 0 2 2 3 0
2 4 1 3 4 2 2 2 1 4
0 0 1 3 0 2 1 2 3 1
3 5 1 4 4 4 3 1 3 5
1 3 3 4 5 2 5 4 5 4
5 0 1 1 0 5 2 3 3 3
4 4 3 0 4 2 1 1 2 5
5 1 4 1 3 0 4 5 2 4
4 1 4 5 4 3 2 1 4 4
4 3 4 0 3 3 2 2 3 4
5 5 2 3 3 0 1 3 3 3
1 1 5 1 5 0 5 5 5 2
1 4 4 1 4 5 4 3 1 1
4 3 3 5 1 0 4 4 2 2
4 0 0 3 2 3 0 5 2 3
0 4 3 3 4 4 2 5 5 3
3 4 2 4 0 0 0 1 1 1`

func bestNumber(counts []int) string {
	best := "1" + strings.Repeat("0", counts[0]+1)
	for d := 1; d <= 9; d++ {
		digit := strconv.Itoa(d)
		cand := strings.Repeat(digit, counts[d]+1)
		if len(cand) < len(best) || (len(cand) == len(best) && cand < best) {
			best = cand
		}
	}
	return best
}

type testCase struct {
	counts []int
}

func parseTestcases(raw string) ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(raw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("no test count")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		counts := make([]int, 10)
		for j := 0; j < 10; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("missing digit count for case %d", i+1)
			}
			val, err := strconv.Atoi(scan.Text())
			if err != nil {
				return nil, fmt.Errorf("invalid count on case %d: %w", i+1, err)
			}
			counts[j] = val
		}
		tests = append(tests, testCase{counts: counts})
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test data:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		var input strings.Builder
		input.WriteString("1\n")
		for j, v := range tc.counts {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		expected := bestNumber(tc.counts)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		out, err := cmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

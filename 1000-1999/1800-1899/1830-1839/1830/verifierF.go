package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `5 5 3 5 5 5 4 4 1 2 1 3 15 7 12 17 3
5 4 1 2 4 4 2 3 2 2 2 4 19 14 4 4
1 1 1 1 5
2 5 3 3 5 5 5 6 12 9 0
3 7 2 3 3 3 3 5 19 18 0 19 10 2 9
3 5 4 5 2 5 4 4 1 8 0 11 12
1 7 3 6 18 0 14 1 5 19 6
1 4 4 4 16 11 16 8
4 2 2 2 1 2 1 1 2 2 4 10
3 2 2 2 2 2 1 1 9 15
2 1 1 1 1 1 19
3 5 4 5 2 2 1 4 10 6 4 18 4
4 2 1 2 2 2 1 2 2 2 14 19
2 8 8 8 8 8 9 15 12 4 3 12 17 5
4 6 2 2 4 5 5 6 1 3 18 1 9 11 17 8
4 5 3 4 2 2 4 5 3 5 8 14 9 16 20
3 6 3 5 6 6 3 4 14 11 10 16 4 16
2 4 3 4 3 3 13 5 19 18
5 7 3 7 5 7 3 3 2 3 5 6 19 20 5 7 5 20 1
4 4 2 2 2 2 3 3 4 4 17 1 13 14
3 7 6 6 5 5 2 7 11 0 11 12 8 13 3
5 6 1 5 5 6 1 3 5 6 5 6 11 4 13 13 18 20
5 6 4 4 2 6 4 6 4 4 2 6 2 11 0 12 3 10
5 3 2 3 2 3 1 2 2 3 3 3 12 5 19
5 5 3 4 3 4 1 3 3 4 3 3 15 0 3 19 14
2 5 1 2 4 4 15 17 17 8 7
4 1 1 1 1 1 1 1 1 1 16
5 2 1 1 2 2 2 2 1 2 2 2 18 4
1 1 1 1 3
5 4 1 4 4 4 4 4 4 4 3 4 2 9 0 13
5 5 4 5 2 3 4 5 3 5 2 5 18 17 1 2 7
3 2 1 1 2 2 1 2 15 1
1 2 1 1 4 9
4 3 1 3 1 2 1 1 3 3 1 17 3
4 2 2 2 1 2 1 1 2 2 10 19
4 8 6 8 6 6 3 8 5 6 18 20 16 9 17 17 19 7
3 2 1 2 2 2 1 2 14 12
2 3 3 3 3 3 2 18 1
1 2 2 2 12 14
4 6 1 5 1 5 6 6 1 2 13 7 3 2 15 6
2 7 3 4 3 5 19 11 12 12 4 11 20
3 7 7 7 5 5 5 7 6 5 12 2 3 1 1
2 4 2 2 4 4 11 0 13 15
3 7 3 6 4 4 2 3 20 5 2 11 12 15 4
5 5 1 3 2 4 2 2 4 4 3 4 10 1 0 14 15
2 2 2 2 2 2 5 1
2 1 1 1 1 1 9
2 8 6 7 2 7 14 5 7 5 6 1 20 19
4 5 1 4 1 4 4 4 1 1 17 16 18 11 3
1 4 4 4 15 1 7 20
1 8 7 7 1 8 13 14 9 20 1 1
2 3 3 3 1 3 7 2 10
1 2 1 1 6 19
1 8 6 7 13 18 12 16 20 4 0 9
4 4 1 2 3 4 1 3 4 4 4 4 17 2
4 4 1 4 2 3 1 4 2 2 6 9 11 3
1 7 7 7 12 17 0 9 20 4 18
1 8 3 8 20 0 1 13 8 3 12 17
2 6 1 3 2 3 8 15 10 9 1 2
5 1 1 1 1 1 1 1 1 1 1 1 4
3 4 3 4 1 2 2 4 3 8 16 19
2 3 2 3 2 2 10 10 14
2 8 7 7 8 8 5 17 15 15 19 0 3 4
4 1 1 1 1 1 1 1 1 1 6
2 5 2 4 5 5 7 11 11 7 1
5 8 7 8 3 4 4 7 7 7 2 3 12 3 9 6 10 15 15 17
5 2 1 2 1 1 1 2 1 2 1 2 16 12
5 8 7 8 1 3 7 8 1 4 5 7 2 0 2 11 19 18 13 9
1 5 2 2 1 2 9 11 19
5 4 4 4 1 2 3 4 4 4 4 4 5 20 16 6
5 6 4 5 1 6 5 6 2 2 2 4 9 8 19 2 5 2
4 3 1 3 1 3 3 3 2 2 13 15 16
5 4 3 4 3 3 3 3 4 4 2 2 9 18 17 19
5 3 2 2 2 2 3 3 3 3 1 2 4 16 2
3 4 4 4 2 4 4 4 0 6 5 3
1 6 2 4 19 3 18 14 10 5
5 2 2 2 1 1 1 2 2 2 2 2 7 12
5 8 7 7 2 6 6 6 1 8 7 8 2 15 9 0 3 15 11 6
3 7 2 4 4 4 5 5 9 13 10 2 16 13 1
2 7 2 4 6 6 4 8 3 17 15 13 5
5 7 7 7 1 4 3 4 7 7 1 7 4 11 0 0 8 17 14
4 2 1 1 1 2 1 2 1 2 19 15
5 1 1 1 1 1 1 1 1 1 1 1 12
1 3 2 2 9 6 13
1 6 2 5 13 6 8 4 15 11
3 1 1 1 1 1 1 1 1
4 3 1 3 1 3 2 3 3 3 4 14 7
1 6 5 5 0 5 13 2 8 4
3 6 2 3 3 3 3 3 3 20 12 11 9 11
4 8 4 8 3 5 1 1 7 7 7 15 1 0 0 6 13 7
5 4 3 3 2 2 4 4 1 3 4 4 3 20 3 11
5 5 5 5 4 5 2 2 4 4 2 5 3 16 9 17 11
4 3 2 3 1 2 1 3 3 3 18 14 1
1 2 1 2 4 11
3 5 3 4 2 2 3 4 0 20 8 12 5
2 5 5 5 2 3 18 1 11 17 3
1 8 1 8 20 1 0 18 13 3 13 16
5 8 4 8 2 7 2 7 2 2 2 7 14 9 20 16 17 2 20 8
1 7 7 7 8 14 9 9 0 10 16
4 8 7 7 4 6 6 8 4 7 2 4 11 2 20 13 14 19`

type testCase struct {
	input    string
	expected string
}

func maxCostNaive(n, m int, intervals [][2]int, p []int) int64 {
	if m > 20 {
		return 0
	}
	var best int64
	totalMasks := 1 << m
	for mask := 0; mask < totalMasks; mask++ {
		var sum int64
		for _, inter := range intervals {
			l, r := inter[0], inter[1]
			val := 0
			for x := r; x >= l; x-- {
				if mask&(1<<(x-1)) != 0 {
					val = p[x-1]
					break
				}
			}
			sum += int64(val)
		}
		if sum > best {
			best = sum
		}
	}
	return best
}

func solve(n, m int, intervals [][2]int, p []int) string {
	return strconv.FormatInt(maxCostNaive(n, m, intervals, p), 10)
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcaseData, "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tokens := strings.Fields(line)
		if len(tokens) < 2 {
			return nil, fmt.Errorf("case %d: not enough tokens", idx+1)
		}
		n, err := strconv.Atoi(tokens[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", idx+1, err)
		}
		m, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad m: %w", idx+1, err)
		}
		need := 2 + 2*n + m
		if len(tokens) != need {
			return nil, fmt.Errorf("case %d: expected %d tokens got %d", idx+1, need, len(tokens))
		}
		intervals := make([][2]int, n)
		for i := 0; i < n; i++ {
			a, err := strconv.Atoi(tokens[2+2*i])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad interval a: %w", idx+1, err)
			}
			b, err := strconv.Atoi(tokens[2+2*i+1])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad interval b: %w", idx+1, err)
			}
			intervals[i] = [2]int{a, b}
		}
		p := make([]int, m)
		base := 2 + 2*n
		for i := 0; i < m; i++ {
			val, err := strconv.Atoi(tokens[base+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad p value: %w", idx+1, err)
			}
			p[i] = val
		}

		var inBuilder strings.Builder
		fmt.Fprintf(&inBuilder, "1\n%d %d\n", n, m)
		for _, inter := range intervals {
			fmt.Fprintf(&inBuilder, "%d %d\n", inter[0], inter[1])
		}
		for i, v := range p {
			if i > 0 {
				inBuilder.WriteByte(' ')
			}
			inBuilder.WriteString(strconv.Itoa(v))
		}
		inBuilder.WriteByte('\n')

		cases = append(cases, testCase{
			input:    inBuilder.String(),
			expected: solve(n, m, intervals, p),
		})
	}
	return cases, nil
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
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

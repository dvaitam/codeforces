package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `3 3 3 2 1 2 3 1 2 1 3
6 3 4 6 5 3 2 1 1 4 2 3 3 5
3 2 1 3 2 2 3 2 3
6 2 3 1 5 4 6 2 5 6 1 4
6 3 1 3 6 4 2 5 4 6 2 4 3 5
8 1 7 8 5 6 3 4 1 2 2 8
3 1 3 1 2 2 3
3 3 1 3 2 1 3 1 2 1 3
8 2 6 5 1 4 7 8 3 2 3 4 3 6
2 3 1 2 1 2 1 2 1 2
5 3 4 1 5 3 2 3 4 3 4 4 5
4 3 2 1 4 3 2 3 2 3 1 2
2 1 2 1 1 2
6 4 6 4 1 3 2 5 5 6 4 5 4 5 2 6
5 4 5 3 4 2 1 2 5 4 5 4 5 1 2
4 3 4 1 3 2 1 3 2 3 2 3
4 1 2 1 3 4 2 4
2 4 2 1 1 2 1 2 1 2 1 2
4 4 3 2 1 4 3 4 1 3 3 4 2 3
8 2 2 3 1 7 6 4 5 8 2 4 1 6
4 1 3 1 2 4 3 4
3 2 1 2 3 2 3 1 2
2 1 1 2 1 2
2 3 2 1 1 2 1 2 1 2
5 2 5 2 1 4 3 4 5 1 2
6 3 1 2 3 6 4 5 1 3 3 5 4 5
8 2 1 3 4 5 8 2 7 6 7 8 3 8
2 2 2 1 1 2 1 2
8 4 7 1 6 3 4 8 5 2 7 8 7 8 6 7 5 7
6 4 1 2 6 4 3 5 4 5 3 6 3 6 3 4
5 1 4 1 3 2 5 3 5
3 3 2 1 3 1 3 1 3 2 3
3 3 2 3 1 1 3 1 2 1 3
7 1 5 6 3 4 2 1 7 5 7
8 4 7 1 5 3 2 6 8 4 6 8 5 8 5 7 7 8
6 4 5 1 2 4 6 3 3 4 5 6 5 6 4 6
6 1 4 3 6 1 5 2 5 6
2 3 1 2 1 2 1 2 1 2
4 2 3 2 4 1 3 4 1 4
3 2 2 3 1 1 3 1 2
4 1 2 1 4 3 3 4
7 4 7 5 4 3 1 6 2 6 7 5 7 2 4 3 7
6 1 1 3 6 5 2 4 5 6
8 4 3 7 5 8 1 4 2 6 6 7 3 5 1 3 1 8
6 1 3 4 5 1 6 2 4 6
4 2 3 4 1 2 3 4 2 4
6 1 1 5 4 3 2 6 4 6
8 1 5 2 8 7 6 1 3 4 3 4
3 1 3 1 2 1 2
7 1 6 4 1 7 2 5 3 1 3
6 3 3 2 4 6 1 5 2 4 3 4 3 5
8 4 3 6 8 5 4 2 1 7 5 8 3 5 6 8 3 5
3 1 1 3 2 2 3
8 4 4 1 6 8 2 5 3 7 1 5 3 8 3 5 5 7
2 2 2 1 1 2 1 2
8 1 8 3 7 2 5 4 1 6 2 7
8 1 4 1 5 8 6 3 7 2 6 7
7 1 4 5 1 3 6 2 7 6 7
7 3 5 3 4 1 6 7 2 6 7 6 7 6 7
2 2 2 1 1 2 1 2
4 4 2 3 1 4 1 2 2 3 2 3 2 3
6 3 2 5 6 1 3 4 3 4 5 6 1 2
2 2 1 2 1 2 1 2
8 4 1 8 7 6 5 4 2 3 1 5 6 7 4 5 6 8
3 1 1 3 2 1 2
5 2 2 3 4 1 5 3 4 1 5
5 2 4 1 3 2 5 2 3 3 4
5 3 4 2 1 3 5 3 5 3 4 2 4
6 2 2 3 5 1 6 4 2 4 5 6
6 1 2 4 6 1 3 5 5 6
6 2 4 6 2 1 5 3 3 4 4 6
6 1 6 5 1 4 3 2 3 6
6 3 5 6 3 4 2 1 4 5 3 4 2 4
3 1 2 3 1 1 3
6 4 2 5 4 3 6 1 2 3 1 2 2 6 2 3
8 3 5 4 7 1 2 3 8 6 5 6 2 3 7 8
6 4 4 2 1 6 5 3 5 6 2 6 2 6 5 6
6 2 6 1 2 4 5 3 3 4 3 5
8 3 1 6 2 3 7 8 5 4 7 8 3 4 7 8
8 3 7 6 5 3 1 8 2 4 1 8 1 7 6 7
3 2 1 3 2 2 3 2 3
4 1 2 4 1 3 1 2
5 1 3 1 5 2 4 2 4
4 3 3 4 1 2 2 3 2 3 2 4
6 4 2 3 6 4 1 5 4 6 5 6 1 2 3 4
3 2 1 3 2 2 3 2 3
7 3 6 5 2 3 1 4 7 6 7 4 7 3 5
6 2 5 6 2 1 3 4 5 6 5 6
8 4 6 1 2 7 8 3 5 4 4 5 7 8 7 8 4 7
3 2 1 3 2 1 3 2 3
5 4 2 1 5 3 4 4 5 2 5 2 3 4 5
3 1 2 3 1 2 3
4 3 2 1 3 4 2 3 3 4 1 3
2 2 2 1 1 2 1 2
3 1 2 3 1 1 2
6 4 1 2 4 5 6 3 5 6 3 6 1 6 3 5
2 1 1 2 1 2
2 2 2 1 1 2 1 2
2 1 2 1 1 2
6 4 4 3 5 2 1 6 2 6 5 6 4 6 5 6`

type testCase struct {
	input    string
	expected string
}

func bullyCount(p []int) int {
	n := len(p)
	arr := make([]int, n)
	copy(arr, p)
	steps := 0
	for {
		i := -1
		maxVal := -1
		for idx, val := range arr {
			if val != idx+1 && val > maxVal {
				maxVal = val
				i = idx
			}
		}
		if i == -1 {
			break
		}
		j := -1
		minVal := int(1e9)
		for k := i + 1; k < n; k++ {
			if arr[k] < minVal {
				minVal = arr[k]
				j = k
			}
		}
		arr[i], arr[j] = arr[j], arr[i]
		steps++
	}
	return steps
}

func solve(n, q int, perm []int, queries [][2]int) string {
	p := make([]int, n)
	copy(p, perm)
	var out strings.Builder
	for _, qr := range queries {
		x, y := qr[0]-1, qr[1]-1
		p[x], p[y] = p[y], p[x]
		fmt.Fprintf(&out, "%d\n", bullyCount(p))
	}
	return strings.TrimSpace(out.String())
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
		q, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad q: %w", idx+1, err)
		}
		need := 2 + n + 2*q
		if len(tokens) != need {
			return nil, fmt.Errorf("case %d: expected %d tokens got %d", idx+1, need, len(tokens))
		}
		perm := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(tokens[2+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad perm value: %w", idx+1, err)
			}
			perm[i] = val
		}
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			a, err := strconv.Atoi(tokens[2+n+2*i])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad query a: %w", idx+1, err)
			}
			b, err := strconv.Atoi(tokens[2+n+2*i+1])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad query b: %w", idx+1, err)
			}
			queries[i] = [2]int{a, b}
		}

		var inBuilder strings.Builder
		fmt.Fprintf(&inBuilder, "%d %d\n", n, q)
		for i, v := range perm {
			if i > 0 {
				inBuilder.WriteByte(' ')
			}
			inBuilder.WriteString(strconv.Itoa(v))
		}
		inBuilder.WriteByte('\n')
		for _, qr := range queries {
			fmt.Fprintf(&inBuilder, "%d %d\n", qr[0], qr[1])
		}
		cases = append(cases, testCase{
			input:    inBuilder.String(),
			expected: solve(n, q, perm, queries),
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
		fmt.Println("usage: verifierE /path/to/binary")
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

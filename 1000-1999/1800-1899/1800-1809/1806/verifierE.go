package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `3 2 3 4 1 1 1 2 3 3 2
5 1 4 1 5 4 4 1 1 1 4 3 4
5 1 2 4 3 3 4 1 1 3 2 2 3
5 2 2 4 4 2 2 1 1 2 1 5 2 2 3
3 3 5 4 5 1 1 3 2 3 2 2 3
3 3 1 5 4 1 1 3 2 2 3 3 2
5 3 3 4 1 2 2 1 1 1 2 3 2 4 2 3 2
5 3 2 1 2 5 4 1 1 2 4 2 3 2 3 3 2
4 1 4 3 1 5 1 1 2 2 3
5 3 5 3 2 5 3 1 1 1 1 4 5 2 4 5 3
5 3 1 1 4 5 2 1 1 1 1 5 2 4 3 2 3
4 2 5 2 5 5 1 1 1 4 2 3 2
4 2 2 3 1 5 1 1 3 2 3 2 3
5 1 5 2 3 1 5 1 1 3 2 2 3
3 3 4 3 2 1 1 3 2 3 2 3 2
4 1 1 5 1 3 1 1 1 4 3
3 3 5 3 4 1 1 3 2 2 3 2 3
3 3 4 3 1 1 1 2 3 3 2 2 3
3 2 5 1 4 1 1 3 2 3 2
5 3 1 4 3 2 5 1 1 1 2 3 2 2 3 2 4
5 1 3 5 4 1 5 1 1 1 4 3 4
5 2 1 2 5 3 2 1 1 1 2 3 2 4 3
3 1 5 2 5 1 1 2 3
5 1 4 5 5 4 1 1 1 1 1 2 5
3 3 4 4 3 1 1 3 2 2 3 2 3
3 2 2 5 5 1 1 2 3 3 2
4 3 3 2 1 4 1 1 3 3 2 3 2 3 2
4 2 5 4 4 5 1 1 2 2 3 2 3
3 3 2 5 5 1 1 3 2 3 2 2 3
5 3 1 1 2 5 1 1 1 2 1 2 5 5 3 5 2
3 2 4 3 3 1 1 3 2 2 3
4 2 2 4 2 3 1 1 1 3 2 4 2
3 3 1 3 2 1 1 3 2 3 2 3 2
5 3 4 4 2 5 4 1 1 1 1 3 4 4 5 5 3
4 1 3 5 5 1 1 1 1 4 2
3 1 4 4 4 1 1 3 2
3 1 2 4 1 1 1 3 2
4 2 4 5 1 4 1 1 3 2 3 2 3
4 2 3 2 4 1 1 1 3 2 3 2 3
3 1 3 1 3 1 1 3 2
5 2 4 4 5 5 1 1 1 1 3 3 4 4 3
4 1 1 2 2 1 1 1 2 2 3
4 2 4 4 5 3 1 1 3 3 2 3 2
4 3 1 2 5 1 1 1 2 3 2 2 3 3 2
4 1 3 3 4 4 1 1 2 2 3
4 3 2 1 4 5 1 1 3 3 2 2 3 2 3
4 3 3 4 5 2 1 1 3 3 2 3 2 2 3
3 3 5 3 4 1 1 3 2 2 3 2 3
5 2 2 2 5 3 5 1 1 1 1 3 5 4 5
3 2 4 3 1 1 1 2 3 3 2
4 3 3 4 4 3 1 1 3 2 3 3 2 3 2
5 2 3 2 1 1 1 1 1 3 1 5 2 2 5
5 3 2 5 5 5 2 1 1 2 3 2 3 2 3 5 4
3 3 5 3 2 1 1 2 3 2 3 3 2
4 1 5 3 5 2 1 1 1 2 4
5 1 2 4 5 1 5 1 1 2 2 3 2
5 1 3 2 1 3 4 1 1 2 3 4 5
5 1 2 4 1 5 1 1 1 3 2 4 5
4 3 5 2 2 1 1 1 2 3 2 3 2 2 3
5 1 4 2 4 1 1 1 1 3 4 3 2
3 1 4 1 2 1 1 2 3
5 3 4 3 4 5 1 1 1 2 3 2 3 4 5 2 3
5 2 1 5 3 3 4 1 1 2 3 2 3 4 5
3 2 2 1 4 1 1 2 3 2 3
4 3 1 3 2 2 1 1 2 3 2 2 3 3 2
3 3 4 3 2 1 1 3 2 2 3 3 2
3 1 3 4 5 1 1 3 2
5 2 5 3 5 4 2 1 1 1 1 5 4 2 5
3 3 5 2 2 1 1 2 3 2 3 3 2
5 1 3 3 3 4 4 1 1 1 3 2 4
3 1 5 2 4 1 1 2 3
5 3 5 5 2 3 5 1 1 3 4 3 2 3 2 3 2
4 1 4 1 4 4 1 1 1 4 3
5 3 5 5 2 2 1 1 1 2 1 5 3 3 5 5 3
3 3 4 3 2 1 1 2 3 2 3 3 2
3 2 4 2 1 1 1 3 2 2 3
5 2 4 3 5 1 1 1 1 3 4 2 3 3 2
5 3 1 4 1 3 5 1 1 1 3 4 2 2 4 2 3
4 2 2 2 3 5 1 1 2 2 3 2 3
4 3 1 5 3 1 1 1 1 4 2 2 4 2 4
5 2 2 4 5 1 5 1 1 2 1 2 3 5 3
3 2 4 2 4 1 1 2 3 3 2
4 2 5 3 4 4 1 1 3 3 2 2 3
5 2 1 3 4 5 3 1 1 3 2 2 3 2 3
4 3 2 1 4 4 1 1 2 3 2 2 3 2 3
5 3 3 2 1 4 1 1 1 3 4 2 3 3 2 2 3
5 3 2 3 4 4 5 1 1 1 4 4 3 3 4 4 3
5 3 3 3 4 5 2 1 1 2 1 3 5 2 5 2 3
3 2 1 1 2 1 1 2 3 3 2
3 2 4 4 3 1 1 2 3 2 3
5 3 3 3 3 4 1 1 1 3 3 5 4 3 2 5 4
5 1 4 1 1 1 3 1 1 1 3 2 4
3 2 1 5 2 1 1 2 3 2 3
3 1 4 4 1 1 1 3 2
3 3 1 2 3 1 1 2 3 2 3 3 2
5 2 2 5 1 2 5 1 1 2 3 5 4 2 3
3 2 2 2 2 1 1 3 2 3 2
3 1 2 3 3 1 1 2 3
3 3 3 3 2 1 1 2 3 2 3 3 2
4 1 4 1 4 5 1 1 2 3 2`

type testCase struct {
	input    string
	expected string
}

func solve(a []int64, parent []int, queries [][2]int) int64 {
	var result []int64
	for _, q := range queries {
		x, y := q[0], q[1]
		var ans int64
		for x > 0 && y > 0 {
			ans += a[x] * a[y]
			x = parent[x]
			y = parent[y]
		}
		result = append(result, ans)
	}
	// only last answer needed for single test? Actually each query outputs one line; concatenate outside.
	if len(result) == 0 {
		return 0
	}
	// This helper returns nothing; handled in caller.
	return 0
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcaseData, "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		p := 0
		n, err := strconv.Atoi(fields[p])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %w", idx+1, err)
		}
		p++
		q, err := strconv.Atoi(fields[p])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad q: %w", idx+1, err)
		}
		p++
		if len(fields) != 2+n+(n-1)+2*q {
			return nil, fmt.Errorf("line %d: expected %d values got %d", idx+1, 2+n+(n-1)+2*q, len(fields))
		}
		a := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			v, err := strconv.ParseInt(fields[p], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: bad a: %w", idx+1, err)
			}
			a[i] = v
			p++
		}
		parent := make([]int, n+1)
		parent[1] = 0
		for i := 2; i <= n; i++ {
			v, err := strconv.Atoi(fields[p])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad parent: %w", idx+1, err)
			}
			parent[i] = v
			p++
		}
		queries := make([][2]int, q)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(a[i], 10))
		}
		sb.WriteByte('\n')
		for i := 2; i <= n; i++ {
			if i > 2 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(parent[i]))
		}
		sb.WriteByte('\n')
		var out strings.Builder
		for i := 0; i < q; i++ {
			x, _ := strconv.Atoi(fields[p])
			y, _ := strconv.Atoi(fields[p+1])
			p += 2
			queries[i] = [2]int{x, y}
			sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
			var ans int64
			u, v := x, y
			for u > 0 && v > 0 {
				ans += a[u] * a[v]
				u = parent[u]
				v = parent[v]
			}
			out.WriteString(strconv.FormatInt(ans, 10))
			if i+1 < q {
				out.WriteByte('\n')
			}
		}
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: out.String(),
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

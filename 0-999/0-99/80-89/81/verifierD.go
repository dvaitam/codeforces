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

const testcases = `
5 4 2 2 4 1
3 4 5 4 3 6
6 2 3 6
6 2 2 5
8 5 5 2 6 4 3
4 3 2 1 3
4 4 3 4 6 4
7 1 7
4 2 4 1
4 2 6 6
5 1 5
4 5 5 6 6 3 1
6 2 2 5
7 1 7
5 2 6 4
8 2 5 6
5 1 5
5 5 5 2 4 2 5
5 3 2 4 2
8 4 1 4 4 6
5 3 2 5 5
3 2 5 6
5 1 5
5 3 4 3 1
7 5 6 4 3 4 6
6 2 2 6
8 4 1 3 2 3
6 5 5 3 2 1 3
3 2 5 2
3 4 4 5 1 5
7 5 4 4 1 1 5
4 3 6 1 2
5 2 3 2
3 5 5 5 5 3 3
5 2 2 5
8 4 1 4 2 6
3 3 2 4 4
8 1 8
5 5 4 2 1 1 5
6 2 5 6
7 2 5 6
7 1 7
8 2 6 3
6 4 4 5 3 4
8 5 4 1 1 5 1
5 4 6 5 2 2
4 1 4
4 1 4
8 5 5 3 1 1 5
4 5 1 2 2 3 4
8 2 2 6
4 3 1 1 5
6 1 6
8 5 2 3 4 5 5
3 2 6 4
4 3 3 3 1
3 3 5 5 5
4 2 4 6
7 4 6 2 6 4
3 5 2 2 6 5 3
5 3 4 3 5
6 3 5 5 3
6 2 6 1
5 3 6 4 6
6 3 3 4 1
6 1 6
3 4 1 4 1 2
5 4 3 6 4 1
8 1 8
7 3 2 2 4
4 4 1 3 5 1
3 1 3
4 2 2 2
5 3 1 4 1
4 5 2 6 4 3 5
6 5 6 1 2 3 4
4 3 4 4 6
8 2 7 1
5 5 1 6 4 5 2
5 2 1 5
7 2 6 2
8 3 2 5 1
7 1 7
3 5 1 1 6 3 6
4 4 6 3 5 2
3 3 4 5 2
7 1 7
3 3 6 4 1
7 1 7
8 5 4 4 1 5 6
7 3 4 5 2
8 5 4 1 4 1 3
3 4 5 5 3 5
3 2 2 2
6 5 5 5 1 4 6
3 1 3
6 5 6 6 1 1 5
4 3 4 1 2
6 4 6 6 5 2
4 2 5 6

`

// referenceSolve mirrors 81D.go greedy construction.
func referenceSolve(n, m int, cnt []int) ([]int, bool) {
	v := make([]int, m+1)
	for i := 1; i <= m; i++ {
		v[i] = cnt[i-1]
	}
	ans := make([]int, n)
	last := -1
	for i := 0; i < n; i++ {
		id := 0
		for j := 1; j <= m; j++ {
			if j == last || (i == n-1 && j == ans[0]) {
				continue
			}
			if v[j] > v[id] {
				id = j
			} else if v[j] == v[id] && j == ans[0] {
				id = j
			}
		}
		if v[id] == 0 {
			return nil, false
		}
		v[id]--
		ans[i] = id
		last = id
	}
	return ans, true
}

type testCase struct {
	n int
	m int
	a []int
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
			return nil, fmt.Errorf("missing m")
		}
		m, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("parse m: %w", err)
		}
		a := make([]int, m)
		for i := 0; i < m; i++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d: missing a[%d]", len(cases)+1, i)
			}
			a[i], err = strconv.Atoi(scan.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d: parse a[%d]: %w", len(cases)+1, i, err)
			}
		}
		cases = append(cases, testCase{n: n, m: m, a: a})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

// validateOutput checks if the contestant output satisfies the problem constraints.
func validateOutput(n, m int, a []int, out string, hasSolution bool) error {
	out = strings.TrimSpace(out)
	if out == "-1" {
		if hasSolution {
			return fmt.Errorf("solution exists but got -1")
		}
		return nil
	}

	parts := strings.Fields(out)
	if len(parts) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(parts))
	}
	used := make([]int, m+1)
	prev := -1
	first := -1
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return fmt.Errorf("output contains non-integer: %v", p)
		}
		if v < 1 || v > m {
			return fmt.Errorf("album %d out of range", v)
		}
		used[v]++
		if used[v] > a[v-1] {
			return fmt.Errorf("album %d used more than available", v)
		}
		if i == 0 {
			first = v
		} else if v == prev {
			return fmt.Errorf("adjacent photos from album %d", v)
		}
		prev = v
	}
	if n > 1 && prev == first {
		return fmt.Errorf("first and last photos from same album")
	}
	total := 0
	for i := 1; i <= m; i++ {
		total += used[i]
	}
	if total != n {
		return fmt.Errorf("expected %d photos, got %d", n, total)
	}
	return nil
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
		fmt.Println("usage: verifierD /path/to/binary")
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
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		_, hasSolution := referenceSolve(tc.n, tc.m, tc.a)
		out, stderr, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr)
			os.Exit(1)
		}
		if err := validateOutput(tc.n, tc.m, tc.a, out, hasSolution); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

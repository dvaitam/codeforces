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
4 -3 1 2 0 2 2 -2 -3 -2 -2 -2 -3
5 0 0 -3 3 1 -2 1 -2 -2 0 -3 -2 3 -3 1
2 1 -2 1 1 1 2
5 0 1 0 3 2 -1 2 3 0 2 2 1 3 -3 -1
3 3 3 -2 2 2 -3 -2 0 -2
5 2 -2 3 -3 -2 -2 0 0 -3 1 -1 3 2 -1 0
3 2 2 2 2 -1 2 -3 1 -3
4 -1 1 1 2 -3 -2 3 -2 0 -2 3 -3
4 -1 -2 1 3 2 -1 2 2 -2 -3 1 -1
3 2 2 -3 -1 3 3 -2 -2 -1
5 2 0 2 2 2 0 0 -3 -2 0 -3 -1 2 2 -2
5 -2 3 1 0 -1 -2 1 -2 1 0 -1 2 -2 -3 -2
2 0 3 1 0 -2 3
1 -1 3 0
5 -3 2 -1 2 -3 -3 0 -2 1 -3 -3 3 2 2 3
5 -3 0 3 -1 2 1 -1 2 1 3 0 -2 -2 2 -3
6 3 1 2 3 1 -2 0 2 2 -2 -3 -2 -2 2 1 2 3 3
2 3 2 -2 3 3 3
5 1 -2 -3 -3 -1 -3 -1 1 0 1 3 -1 1 0 -3
6 -3 -2 3 0 -3 -1 -1 -1 2 0 0 -3 -1 2 -2 -2 3 -3
5 2 -2 -2 3 1 -2 0 -1 3 -3 3 -1 0 1 1
5 -2 2 -1 2 1 2 1 0 -2 -1 -3 0 2 2 -1
2 3 2 -1 3 -1 -3
3 -2 0 -3 0 2 0 -3 2 3
3 1 1 3 0 0 1 3 1 2
5 -2 -2 3 -3 -1 -3 3 0 2 -1 -2 0 0 -1 1
4 -2 1 0 -1 3 -2 2 -3 -2 -1 -1 2
6 3 3 -3 2 -3 0 -1 -1 1 2 -2 -1 0 3 -2 2 -3 -2
3 3 2 0 0 -2 2 3 -1 0
4 0 -1 1 -3 -3 1 3 0 -2 -3 -3 3
5 3 2 0 2 0 -2 0 2 2 -3 0 -3 2 3 -1
5 -2 -1 -1 1 -3 1 2 2 -1 -3 -2 -3 -3 -1 3
3 0 3 -2 -2 -1 2 1 0 -1
5 -2 0 0 -3 -3 -3 -3 -1 2 -2 -1 1 -3 1 3
3 -3 -3 1 1 -1 2 -3 -2 3
1 2 -2 0
2 -2 0 -1 3 2 -2
4 2 3 1 -3 1 0 -3 -3 -3 -1 -3 2
4 -2 -1 2 0 -2 0 2 1 1 1 3 -1
2 0 2 3 -3 1 0
5 -2 -1 -1 1 -1 0 2 3 3 2 0 2 2 -1 -1
3 3 -2 -2 -2 0 -1 2 -3 -1
4 1 3 2 -1 2 -2 -3 -3 0 2 1 3
2 -2 3 0 2 1 -1
5 -2 -3 0 0 0 3 0 -1 -3 -2 1 -2 -3 0 -2
6 -3 2 1 2 2 1 0 1 -1 -2 -3 0 -3 1 0 2 2 -1
6 -3 -2 1 -1 3 3 3 3 3 0 0 2 1 2 -1 -3 -3 -3
3 0 2 0 0 3 2 2 3 2
2 -1 0 0 1 -3 2
4 2 -1 2 0 -2 0 2 3 0 0 2 3
5 2 1 -1 0 1 -1 -1 0 3 1 2 3 -2 -1 -3
4 0 0 -1 -3 2 -2 -1 -2 0 -3 1 -1
2 -2 3 0 -3 3 -3
2 1 -3 -2 2 -2 -3
6 1 2 0 2 0 -2 -1 0 0 -1 1 -3 -2 0 -2 -3 2 -3
4 3 2 -1 2 3 -1 -1 0 -3 1 3 -1
3 3 0 -3 -3 3 2 2 3 -1
2 0 2 -1 1 -2 0
1 0 -3 -3
1 -1 -3 -1
3 -3 -3 -3 1 -3 -2 3 1 3
2 0 2 -3 1 -2 -3
6 1 -2 3 -3 -3 2 -3 1 3 2 -3 -2 3 -3 0 0 -3 -3
4 2 3 0 1 2 -2 1 -3 -2 3 -1 2
4 1 -2 -3 2 1 -3 2 1 -2 -3 2 3
4 -1 1 -1 1 2 1 1 0 -1 -2 1 2
1 -2 1 0
3 2 -1 -1 1 -1 -3 -3 -3 -2
3 2 2 2 -2 2 -2 -2 -3 3
4 2 0 2 -2 -1 -1 3 -2 3 -3 -2 2
3 3 2 -1 -2 2 -3 1 3 0
2 2 -1 -1 -2 2 -3
5 2 2 -2 -1 3 0 3 -3 -1 3 0 1 -1 -3 2
3 -2 0 -3 3 0 1 -1 3 3
2 -1 -3 2 -1 2 -3
3 2 0 -3 3 2 -1 -2 2 3
6 -2 -1 -1 3 1 2 -3 -2 0 -2 3 1 1 1 -3 0 -1 3
6 2 -3 -1 -3 -3 3 1 -2 -1 -2 0 -2 1 1 3 -2 -3 2
5 3 3 -2 0 1 2 0 1 1 -1 3 -1 -1 -3 3
1 0 -1 -2
4 3 -2 2 2 2 3 2 3 2 2 2 0
5 -3 2 3 -2 1 2 0 0 3 0 -2 -2 2 0 0
4 3 2 2 2 2 -2 -1 -3 2 -2 1 -3
3 2 -2 -3 -2 -1 -2 0 -1 1
6 0 -1 1 2 0 0 0 -3 -1 1 -1 -1 -1 2 -1 0 1 1
1 0 2 1
5 -2 0 -1 2 2 3 0 2 -3 -3 0 1 -3 3 -1
6 3 -3 -3 0 2 3 0 1 -1 1 0 1 0 1 2 2 2 0
5 -2 3 1 3 -3 1 -1 -2 -1 0 -3 -2 -3 -2 3
5 0 2 -2 0 -1 3 0 3 2 -2 0 3 3 -1 0
3 2 1 2 3 -2 3 -3 -1 1
3 0 3 1 2 1 -3 -1 -2 1
2 -2 -3 -1 -2 -3 -3
4 0 -3 1 1 -3 -2 2 0 1 -2 1 -1
1 -1 -1 -2
4 -1 -3 -1 -3 -3 1 2 -3 -1 2 -1 1
3 -2 3 -3 2 3 0 0 0 2
2 -3 0 3 -1 0 0
6 0 0 -2 0 3 -2 3 0 2 3 -2 -1 -3 3 -2 0 -1 -1
3 3 3 0 -1 -2 -2 0 -2 2

`

const (
	mod  = 9999993
	base = 500000000
	inf  = 1000000000
)

type triple struct{ a, b, c int }
type entry struct {
	a, b, c int
	sta     int
}

type testCase struct {
	n int
	p []triple
}

func parseCases() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcases))
	scanner.Split(bufio.ScanLines)
	var cases []testCase
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("parse n: %w", err)
		}
		expected := 1 + 3*n
		if len(fields) != expected {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", len(cases)+1, expected, len(fields))
		}
		p := make([]triple, n)
		for i := 0; i < n; i++ {
			a, _ := strconv.Atoi(fields[1+3*i])
			b, _ := strconv.Atoi(fields[1+3*i+1])
			c, _ := strconv.Atoi(fields[1+3*i+2])
			p[i] = triple{a: a, b: b, c: c}
		}
		cases = append(cases, testCase{n: n, p: p})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func referenceSolve(p []triple) string {
	n := len(p)
	m := (n + 1) >> 1

	pow3 := make([]int, n+1)
	pow3[0] = 1
	for i := 1; i <= n; i++ {
		pow3[i] = pow3[i-1] * 3
	}

	buckets := make(map[int][]entry)
	ans := -inf * 2
	sta1, sta2 := 0, 0

	var dfs1 func(idx, end, a, b, c, code int)
	var dfs2 func(idx, end, a, b, c, code int)

	dfs1 = func(idx, end, a, b, c, code int) {
		if idx == end {
			x := (a - b + base) % mod
			buckets[x] = append(buckets[x], entry{a: a, b: b, c: c, sta: code})
			return
		}
		t := p[idx]
		dfs1(idx+1, end, a+t.a, b+t.b, c, code*3)
		dfs1(idx+1, end, a+t.a, b, c+t.c, code*3+1)
		dfs1(idx+1, end, a, b+t.b, c+t.c, code*3+2)
	}

	dfs2 = func(idx, end, a, b, c, code int) {
		if idx == end {
			x := (b - a + base) % mod
			best := -inf * 2
			bestSta := 0
			if list, ok := buckets[x]; ok {
				for _, e := range list {
					if e.b-e.c == c-b {
						if e.a > best {
							best = e.a
							bestSta = e.sta
						}
					}
				}
			}
			if best > -inf {
				if best+a > ans {
					ans = best + a
					sta1 = code
					sta2 = bestSta
				}
			}
			return
		}
		t := p[idx]
		dfs2(idx+1, end, a+t.a, b+t.b, c, code*3)
		dfs2(idx+1, end, a+t.a, b, c+t.c, code*3+1)
		dfs2(idx+1, end, a, b+t.b, c+t.c, code*3+2)
	}

	// second half in buckets
	dfs1(m, n, 0, 0, 0, 0)
	// first half queries
	dfs2(0, m, 0, 0, 0, 0)

	if ans < -inf/2 {
		return "Impossible"
	}
	var sb strings.Builder
	for i := 0; i < m; i++ {
		t := (sta1 / pow3[m-1-i]) % 3
		switch t {
		case 0:
			sb.WriteString("LM\n")
		case 1:
			sb.WriteString("LW\n")
		case 2:
			sb.WriteString("MW\n")
		}
	}
	for i := m; i < n; i++ {
		t := (sta2 / pow3[n-1-i]) % 3
		switch t {
		case 0:
			sb.WriteString("LM\n")
		case 1:
			sb.WriteString("LW\n")
		case 2:
			sb.WriteString("MW\n")
		}
	}
	return strings.TrimRight(sb.String(), "\n")
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
		fmt.Fprintf(&input, "%d\n", tc.n)
		for _, t := range tc.p {
			fmt.Fprintf(&input, "%d %d %d\n", t.a, t.b, t.c)
		}
		expected := referenceSolve(tc.p)
		out, stderr, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\ngot:\n%s\n", idx+1, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

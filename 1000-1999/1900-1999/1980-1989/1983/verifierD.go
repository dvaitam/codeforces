package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
100
4 1 9 5 2 1 8 5 3 3 5 7 2
3 1 1 2 1 2 1 2 1 1
1 5 5 5
1 1 1 1
4 1 1 1 1 1 1 1 1 1 1 1 1
5 1 2 9 1 1 1 1 10 1 1 10 1 1 1 1
2 1 1 1 1 1 1
5 1 1 2 8 1 1 3 3 2 4 1 1 6 1 4
1 1 1 1
2 1 2 1 2 2 1
2 1 3 2 2 1 3
1 5 5 5
4 2 1 1 3 1 2 1 3 1 4 1 1
4 1 1 6 1 2 5 1 1 1 4 3 1
3 1 1 7 4 4 1 2 2 5
4 10 1 2 2 6 1 1 7 1 1 12 1
4 7 7 2 1 2 5 3 7 2 1 5 9
3 7 2 1 3 6 1 1 3 6
1 5 5 5
5 1 2 1 1 3 1 1 4 1 1 1 1 3 2 1
1 2 2 2
3 3 6 2 8 2 1 3 2 6
5 3 1 3 4 3 1 1 6 5 1 4 1 1 2 6
1 3 3 3
2 2 8 7 3 5 5
1 1 1 1
5 1 1 1 16 1 14 2 2 1 1 13 1 3 1 2
1 1 1 1
4 1 10 1 8 3 12 3 2 6 1 1 12
3 2 1 1 3 2 1 3 2 1
4 1 1 1 1 2 2 2 2 3 3 3 3 4 4 4 4
5 6 3 4 8 7 1 4 8 5 2 4 2 2 2 2
2 4 1 1 2 3 1
1 3 3 3
1 1 1 1
5 2 5 8 8 2 2 5 11 1 1 9 1 3 1 6
3 5 2 3 2 2 6 2 5 7
1 2 2 2
3 1 2 3 4 5 6 7 8 9
4 1 2 1 1 3 3 2 1 2 4 4 6 2
3 10 1 2 6 1 1 5 8 1 1
3 11 1 8 1 1 2 2 1 5 1
2 1 4 2 2 3 3
5 2 1 3 1 2 1 4 1 2 1 2 1 6 1 2
5 2 1 1 1 1 8 1 1 9 1 1 6 1 1 3
5 1 2 7 1 1 2 3 3 4 5 1 2 3 1 3 6
4 1 1 2 2 3 4 1 1 2 2 1 1 1
4 1 1 1 2 2 2 1 1 1 1 1 1
4 3 1 3 1 1 1 2 3 1 4 3 1 4
3 1 1 3 5 4 1 2 6
2 2 4 2 1 3 4
4 7 5 2 1 1 9 1 3 2 2 2 2 2 1 8 1
1 3 3 3
5 1 1 4 4 1 1 2 2 4 3 1 5 1 1 3
5 1 3 6 1 1 6 4 3 6 2 1 4 1 3 1
2 1 3 1 2 1 1
1 10 10 10
1 2 2 2
3 4 1 1 2 3 4 5 1 2 3
1 2 2 2
4 1 2 1 2 3 3 3 4 2 1 1 1
5 2 1 1 1 1 4 6 3 1 7 6 2 1 4 7 3
1 2 2 2
2 1 1 1 1 1 1
5 1 1 1 6 2 4 4 1 1 6 4 1 1 6 5
1 1 1 1
5 1 1 2 3 4 1 1 3 2 5 1 1 7 1
5 8 2 1 1 4 2 3 6 5 4 3 6 1 1 4 4
5 5 1 5 2 5 2 2 7 2 1 8 1 1 3 9
3 8 1 2 2 2 1 2 2 8 8
5 3 3 2 3 1 4 2 4 1 1 2 5 1 2 1
5 3 2 5 1 1 2 5 2 2 2 2 4 1 3 4
4 1 1 1 1 1 1 4 1 3 2 1 2
3 10 7 6 4 6 1 1 6 4 1
4 9 3 1 2 2 2 2 1 1 1 2 2 1
2 3 8 3 1 8 1
4 2 1 1 1 1 2 1 2 1 2 3 2 3
5 7 4 1 2 2 1 3 1 2 2 1 3 1 3 3
1 2 2 2
1 4 4 4
2 1 1 1 3 4 3
1 7 7 7
5 1 1 1 1 2 2 2 2 1 3 3 3 3 4 4 4
5 1 1 4 1 3 8 7 1 7 6 1 1 5 6 1
4 3 1 5 6 2 8 3 8 2 1 1 3
5 9 1 2 9 1 1 7 2 1 2 1 2 2 1
5 3 1 1 1 4 3 1 1 2 7 1 7 2
4 3 2 2 2 3 3 3 1 1 1 3 4
3 1 1 1 1 1 2 2 2 3
2 1 1 2 5 3 3
4 3 1 4 2 1 5 2 6 2 8 1 1 6
3 2 1 1 2 6 1 1 4 4
2 5 1 1 2 4 3
5 2 3 4 1 1 6 6 6 1 1 8 5 2 1 4
1 2 2 2
5 5 1 3 1 4 2 2 4 4 2 1 5 3 1
5 4 1 1 2 5 2 2 6 3 3 2 2 2 2
3 6 8 3 1 4 8 7 1 1
5 3 1 1 1 2 6 1 1 2 2 6 1 3 1
2 3 1 1 2 3 1
5 4 2 1 1 3 3 3 2 1 1 4 3 2 1 4
3 1 1 7 2 2 1 1 8 4
2 1 1 1 2 3 4
2 1 1 4 6 3 1
5 1 1 1 6 4 4 3 3 1 3 4 6 2 1
1 1 1 1
3 1 1 1 2 2 2 3 3 3
4 1 1 2 2 5 5 1 4 5 5 1 3
1 2 2 2
1 3 3 3
3 1 1 2 1 1 4 2 1 1
1 1 1 1
5 1 1 5 3 5 1 3 1 3 1 3 1 3 1 3
3 1 1 7 7 1 1 3 4
3 4 1 2 2 4 2 1 3 2 1
1 2 2 2
4 4 2 4 1 1 2 4 2 4 3 4 1
4 4 2 1 2 3 3 3 4 4 1 1 1
1 5 5 5
5 1 1 5 3 1 1 4 2 4 6 5 2 6 3 4
2 1 1 1 1 1 1
3 2 3 5 6 3 1 1 3 6
2 2 4 6 2 2 1
1 1 1 1
1 1 1 1
3 2 1 2 2 3 2 2 2
1 1 1 1
`

type testCase struct {
	n int
	a []int
	b []int
}

func parseTests(raw string) ([]testCase, error) {
	tokens := strings.Fields(raw)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty testcases")
	}
	idx := 0
	t, err := strconv.Atoi(tokens[idx])
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	idx++
	cases := make([]testCase, 0, t)
	for c := 0; c < t; c++ {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("truncated before case %d", c+1)
		}
		n, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("bad n in case %d", c+1)
		}
		idx++
		a := make([]int, n)
		for i := 0; i < n; i++ {
			if idx >= len(tokens) {
				return nil, fmt.Errorf("missing a values in case %d", c+1)
			}
			v, err := strconv.Atoi(tokens[idx])
			if err != nil {
				return nil, fmt.Errorf("bad a value in case %d", c+1)
			}
			a[i] = v
			idx++
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			if idx >= len(tokens) {
				return nil, fmt.Errorf("missing b values in case %d", c+1)
			}
			v, err := strconv.Atoi(tokens[idx])
			if err != nil {
				return nil, fmt.Errorf("bad b value in case %d", c+1)
			}
			b[i] = v
			idx++
		}
		cases = append(cases, testCase{n: n, a: a, b: b})
	}
	return cases, nil
}

func solve(tc testCase) string {
	n := tc.n
	a := tc.a
	b := tc.b
	pos := make(map[int]int, n)
	for i := 0; i < n; i++ {
		pos[b[i]] = i
	}
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		p, ok := pos[a[i]]
		if !ok {
			return "NO"
		}
		perm[i] = p
	}
	visited := make([]bool, n)
	cycles := 0
	for i := 0; i < n; i++ {
		if !visited[i] {
			cycles++
			j := i
			for !visited[j] {
				visited[j] = true
				j = perm[j]
			}
		}
	}
	parity := (n - cycles) % 2
	if parity == 0 {
		return "YES"
	}
	return "NO"
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		want := solve(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %q got %q\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

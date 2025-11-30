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

const testcasesA = `4 7 1 5 9 8 7 5 8 6
5 4 9 3 5 3 2 10 5 9 10 3
3 2 2 6 8 9 2 6
4 6 10 4 9 8 8 9 5 1
5 1 2 7 1 10 8 6 4 6 2 4
5 4 4 3 9 8 2 2 6 9 8 2
3 9 5 2 9 6 9 4
5 9 10 5 8 2 10 7 6 10 4 5
2 4 3 1 10 5
4 2 2 3 3 1 2 9 7 9
3 9 4 4 10 7 10 5
4 8 6 2 6 10 2 8 10 6
2 4 1 5 2 4
3 3 6 7 1 2 3 4
1 10 9 10
1 1 2 4
5 10 2 7 2 6 2 1 10 1 4 3
1 8 4 1
1 9 7 10
1 5 2 4
1 5 6 7
2 1 9 8 1 10
1 7 4 5
3 8 10 3 4 1 3 3
3 9 5 2 10 8 3 1
4 7 10 9 5 6 7 5 3 9
1 8 2 6
1 9 5 3
2 8 6 10 5 6
5 10 3 5 7 7 2 1 10 4 6 3
2 4 8 7 10 7
1 7 10 7
1 3 8 2
3 3 8 9 8 9 10 1
1 8 6 5
4 1 7 4 9 2 3 1 7 7
3 1 4 1 1 9 10 2
2 2 10 4 5 5
2 2 8 7 2 1
3 8 2 5 3 9 6 2
2 5 1 1 1 4
3 9 6 6 10 1 10 8
4 7 6 9 3 4 7 10 5 1
2 3 5 6 6 6
1 6 10 1
1 5 3 3
5 5 6 7 9 3 5 2 8 4 1 5
2 9 2 5 7 6
3 7 2 2 9 8 8 6
3 2 8 2 8 7 1 5
3 3 3 10 7 2 2 2
2 4 1 7 1 2
4 9 9 5 8 8 10 4 7 2
3 4 5 10 3 7 4 6
1 2 1 9
4 4 2 8 7 5 4 1 4 10
2 2 4 8 7 6
5 3 2 10 8 3 10 7 7 9 8 6
4 8 4 9 10 4 1 6 6 6
1 9 3 5
5 3 7 10 5 8 2 2 9 1 2 4
2 1 5 1 8 6
2 3 8 6 9 7
5 9 1 10 2 9 10 2 7 4 5 9
5 7 8 7 10 10 4 1 1 3 5 9
5 5 6 2 8 5 5 7 7 7 1 3
2 4 5 6 1 1
4 7 3 8 10 2 3 6 7 1
5 8 7 8 1 2 8 3 1 1 10 10
2 6 2 9 6 4
4 8 2 1 10 8 10 6 2 10
3 3 7 5 2 9 4 1
4 8 6 4 8 6 2 1 1 8
3 1 9 10 10 4 4 2
5 9 7 9 5 2 3 7 10 7 2 2
4 2 2 7 3 1 8 7 7 1
4 6 5 2 6 2 2 6 1 6
3 3 1 4 6 2 10 3
2 1 4 2 1 5
3 1 10 4 3 3 8 2
4 6 5 3 1 4 6 6 8 5
3 9 6 3 10 2 2 9
5 5 3 7 3 3 4 6 9 4 4 3
3 6 7 1 3 10 1 7
1 2 3 7
3 9 7 3 10 7 5 6
1 4 8 6
5 1 7 7 1 7 6 8 4 6 5 8
1 3 2 5
1 9 10 3
4 7 3 7 7 3 4 8 6 9
2 6 8 2 8 4
3 1 8 10 8 1 4 5
1 5 9 10
2 7 8 2 8 4
5 7 5 1 2 5 1 1 5 7 9 10
4 8 2 5 6 5 4 10 2 1
1 5 5 9
3 2 9 4 3 2 7 5
3 9 3 10 9 4 9 2
`

type testCase struct {
	m int
	k int
	d []int
	s []int
}

func ceilDiv(x, y int) int {
	if x <= 0 {
		return 0
	}
	q := x / y
	if x%y != 0 {
		q++
	}
	return q
}

func solve(m, k int, d, s []int) int {
	if m == 0 {
		return 0
	}
	totalD := 0
	sumS := 0
	maxS := s[0]
	maxIdx := 0
	for i := 0; i < m; i++ {
		totalD += d[i]
		sumS += s[i]
		if s[i] > maxS {
			maxS = s[i]
			maxIdx = i
		}
	}

	fuel := s[0]
	waitTime := 0
	for i := 0; i < maxIdx; i++ {
		need := d[i]
		if fuel < need {
			rem := need - fuel
			waitBatches := ceilDiv(rem, s[i])
			waitTime += waitBatches * k
			fuel += waitBatches * s[i]
		}
		fuel = fuel - need + s[i+1]
	}

	extra := totalD - sumS
	if extra < 0 {
		extra = 0
	}
	waitTime += ceilDiv(extra, maxS) * k
	return waitTime + totalD
}

func parseTestcases() ([]testCase, error) {
	sc := bufio.NewScanner(strings.NewReader(testcasesA))
	sc.Buffer(make([]byte, 0, 1024), 1<<20)
	tests := []testCase{}
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		m, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		k, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		if len(parts) != 2+2*m {
			return nil, fmt.Errorf("line mismatch, expected %d numbers got %d", 2+2*m, len(parts))
		}
		d := make([]int, m)
		s := make([]int, m)
		for i := 0; i < m; i++ {
			d[i], _ = strconv.Atoi(parts[2+i])
		}
		for i := 0; i < m; i++ {
			s[i], _ = strconv.Atoi(parts[2+m+i])
		}
		tests = append(tests, testCase{m: m, k: k, d: d, s: s})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func computeExpected(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var m, k int
	if _, err := fmt.Fscan(reader, &m, &k); err != nil {
		return ""
	}
	d := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &d[i])
	}
	s := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &s[i])
	}
	return strconv.Itoa(solve(m, k, d, s))
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", tc.m, tc.k)
	for i, v := range tc.d {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	for i, v := range tc.s {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		input := buildInput(tc)
		expected := computeExpected(input)
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

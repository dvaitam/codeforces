package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input    string
	expected string
}

const testcaseData = `10 4 3 0 4 3 3 2 3 5 0
2 3 4
3 1 0 0
7 1 0 5 2 2 3 0
6 2 2 4 5 1 5
9 4 0 2 4 4 0 1 3 0
5 2 1 2 0 4
9 3 0 2 4 0 4 0 4 1
8 0 4 1 2 0 4 5 1
3 5 5 1
7 1 4 5 5 3 2 1
5 5 2 4 1 1
8 2 1 2 4 5 0 0 5
6 2 5 2 0 0 4
1 2
2 1 5
5 1 4 4 2 3
2 3 4
9 3 0 0 4 3 2 1 2 1
7 5 5 2 0 0 0 5
5 0 5 2 2 4
5 0 5 1 5 4
3 2 4 0
7 3 2 4 1 2 2 4
8 3 5 0 3 2 1 4 0
4 5 3 3 0
4 0 2 0 1
3 3 0 2
9 3 0 4 1 3 2 2 3 3
8 5 0 1 5 4 2 1 1
3 1 2 1
3 4 2 4
1 1
4 5 3 5 1
4 4 1 4 5
3 5 1 4
5 5 5 0 5 3
4 0 3 1 3
9 3 1 3 4 4 3 0 1 0
9 0 0 0 5 1 5 3 3 5
5 3 2 5 3 2
2 1 4
4 5 1 0 0
9 0 5 3 1 3 5 5 1 4
9 2 0 4 1 2 5 2 2 2
9 5 5 2 1 0 3 5 1 4
7 2 3 2 3 1 3 2
6 2 5 4 3 3 4
7 5 4 2 5 4 1 1
7 1 2 4 1 4 1 4
9 4 4 5 3 2 2 3 2 3
9 0 5 1 2 3 5 0 0 1
10 0 2 0 1 1 4 1 5 2 4
5 0 5 5 0 0
10 1 2 1 3 3 5 2 2 0 4
7 0 5 0 1 0 4 4
5 2 5 1 3 1
6 1 1 3 5 1 4
10 4 3 4 5 1 0 0 3 4 0
3 4 0 3
5 2 5 5 4 1
1 2
3 2 5 5
5 3 3 4 2 2
2 3 1
2 4 2
9 5 3 4 3 3 5 1 4 3
3 5 5 0
10 0 1 5 3 2 5 5 0 5 3
5 2 0 5 2 4
9 0 0 0 4 0 1 0 3 3
8 0 2 5 1 0 5 3 5
10 5 1 0 1 3 0 0 4 3 1
9 5 5 2 3 3 4 0 5 3
9 1 4 1 5 1 4 1 0 4
8 2 2 2 1 3 5 0 5
1 4
2 0 5
5 1 1 4 1 0
6 0 5 2 1 2 4
2 0 2
4 2 3 2 5
8 3 4 4 0 0 2 3 3
7 4 3 4 3 4 0 1
3 5 3 4
6 1 0 1 1 0 3
10 0 5 0 4 2 4 2 2 5 5
8 3 4 2 5 5 4 3 4
6 1 4 0 5 3 2
6 4 0 1 4 0 3
2 2 3
10 1 4 0 4 4 0 0 4 2 0
6 1 2 2 3 3 3
1 5
4 5 5 3 5
6 5 5 1 0 5 3
9 5 1 3 5 2 5 4 3 1
1 5
8 4 4 5 1 4 3 2 4
6 3 4 5 1 4 3`

func solveCase(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return "", fmt.Errorf("invalid line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	if len(fields) != n+1 {
		return "", fmt.Errorf("expected %d numbers, got %d", n, len(fields)-1)
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return "", err
		}
		a[i] = v
	}
	var starts, ends []int
	s := 0
	for i := 1; i <= n+1; i++ {
		x := 0
		if i <= n {
			x = a[i-1]
		}
		for s < x {
			starts = append(starts, i)
			s++
		}
		for s > x {
			ends = append(ends, i-1)
			s--
		}
	}
	res := []string{fmt.Sprintf("%d", len(starts))}
	for i := 0; i < len(starts); i++ {
		res = append(res, fmt.Sprintf("%d %d", starts[i], ends[i]))
	}
	return strings.Join(res, "\n"), nil
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		exp, err := solveCase(line)
		if err != nil {
			return nil, fmt.Errorf("case %d: %w", idx+1, err)
		}
		cases = append(cases, testCase{input: line + "\n", expected: strings.TrimSpace(exp)})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(tc.input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != tc.expected {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

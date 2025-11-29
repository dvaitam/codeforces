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

const testcases = `1 1
2 0
3 0
4 1
5 1
6 0
7 1
8 0
9 0
10 1
11 1
12 0
13 1
14 0
15 0
16 1
17 1
18 0
19 1
20 0
21 0
22 1
23 1
24 0
25 1
26 1
27 0
28 1
29 1
30 1
31 1
32 1
33 0
34 1
35 0
36 1
37 1
38 1
39 0
40 1
41 1
42 1
43 1
44 1
45 0
46 1
47 0
48 1
49 1
50 1
51 0
52 1
53 1
54 1
55 1
56 1
57 0
58 1
59 0
60 1
61 1
62 1
63 0
64 1
65 1
66 1
67 1
68 1
69 0
70 1
71 0
72 1
73 1
74 1
75 0
76 1
77 1
78 1
79 1
80 1
81 0
82 1
83 0
84 1
85 1
86 1
87 0
88 1
89 1
90 1
91 1
92 1
93 0
94 1
95 0
96 1
97 1
98 1
99 0
100 1`

func parseTestcases() (map[int]int, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcases))
	cases := make(map[int]int)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid test line %d", lineNum)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("parse n on line %d: %w", lineNum, err)
		}
		val, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("parse val on line %d: %w", lineNum, err)
		}
		cases[n] = val
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return cases, nil
}

func referenceOutputs() ([]int, error) {
	tests, err := parseTestcases()
	if err != nil {
		return nil, err
	}
	results := make([]int, 0, 30)
	for n := 21; n < 51; n++ {
		v, ok := tests[n]
		if !ok {
			return nil, fmt.Errorf("missing testcase for n=%d", n)
		}
		results = append(results, v)
	}
	return results, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	expectedSeq, err := referenceOutputs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	out, err := run(bin, "")
	if err != nil {
		fmt.Printf("runtime error: %v\n", err)
		fmt.Printf("Output: %s\n", out)
		os.Exit(1)
	}
	lines := strings.Fields(out)
	if len(lines) != len(expectedSeq) {
		fmt.Printf("output length mismatch: got %d want %d\n", len(lines), len(expectedSeq))
		os.Exit(1)
	}
	for i, exp := range expectedSeq {
		got, err := strconv.Atoi(lines[i])
		if err != nil {
			fmt.Printf("line %d: failed to parse %q\n", i+1, lines[i])
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("line %d failed: expected %d got %d\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d cases passed\n", len(expectedSeq))
}

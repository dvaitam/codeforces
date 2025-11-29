package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesB.txt to remove external dependency.
const testcasesRaw = `100
5
79
85
34
61
9
12
87
97
17
20
5
11
90
70
88
51
91
68
36
67
31
28
87
76
54
75
36
58
64
85
83
90
46
11
42
79
15
63
76
81
43
25
32
3
94
35
15
91
29
48
22
43
55
8
13
19
90
29
6
74
82
69
78
88
10
4
16
82
25
78
74
16
51
12
48
15
5
78
3
25
24
92
16
62
27
94
8
87
3
70
55
80
13
34
9
29
10
83
39`

func solve(n int) string {
	m := (n + 3) / 4
	return strings.Repeat("9", n-m) + strings.Repeat("8", m)
}

func parseTests() ([]int, error) {
	scan := strings.Fields(testcasesRaw)
	if len(scan) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	t, err := strconv.Atoi(scan[0])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	if len(scan) != t+1 {
		return nil, fmt.Errorf("expected %d tests got %d", t, len(scan)-1)
	}
	tests := make([]int, 0, t)
	for i := 0; i < t; i++ {
		n, err := strconv.Atoi(scan[i+1])
		if err != nil {
			return nil, fmt.Errorf("test %d: bad n: %v", i+1, err)
		}
		tests = append(tests, n)
	}
	return tests, nil
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Println("failed to parse embedded tests:", err)
		os.Exit(1)
	}

	for idx, n := range tests {
		input := []byte(fmt.Sprintf("1\n%d\n", n))
		expected := solve(n)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed: expected %q got %q\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

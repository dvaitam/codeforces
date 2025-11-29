package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
18
73
98
9
33
16
64
98
58
61
84
49
27
13
63
4
50
56
78
98
99
1
90
58
35
93
30
76
14
41
4
3
4
84
70
2
49
88
28
55
93
4
68
29
98
57
64
71
30
45
30
87
29
98
59
38
3
54
72
83
13
24
81
93
38
16
96
43
93
92
65
55
65
86
25
39
37
76
64
65
51
76
5
62
32
96
52
54
86
23
47
71
90
100
87
95
48
12
57
85`

// Embedded reference logic from 1038B.go.
func buildSolution(n int) string {
	var sb strings.Builder
	if n <= 2 {
		sb.WriteString("No")
		return sb.String()
	}
	sb.WriteString("Yes\n")
	if n%2 == 0 {
		fmt.Fprintf(&sb, "2 1 %d\n", n)
		fmt.Fprintf(&sb, "%d", n-2)
		for i := 2; i < n; i++ {
			fmt.Fprintf(&sb, " %d", i)
		}
	} else {
		k := (n + 1) / 2
		fmt.Fprintf(&sb, "1 %d\n", k)
		fmt.Fprintf(&sb, "%d", n-1)
		for i := 1; i <= n; i++ {
			if i == k {
				continue
			}
			fmt.Fprintf(&sb, " %d", i)
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(exe string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	expected := strings.TrimSpace(buildSolution(n))

	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}

	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("unexpected output.\nexpected:\n%s\ngot:\n%s", expected, got)
	}
	return nil
}

type testCase struct {
	n int
}

func parseTestcases(raw string) ([]testCase, error) {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcase data provided")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	if len(fields) != 1+t {
		return nil, fmt.Errorf("expected %d numbers, found %d", 1+t, len(fields))
	}
	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		v, err := strconv.Atoi(fields[1+i])
		if err != nil {
			return nil, fmt.Errorf("invalid n on case %d: %w", i+1, err)
		}
		tests = append(tests, testCase{n: v})
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Println("invalid test data:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		if err := runCase(exe, tc.n); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

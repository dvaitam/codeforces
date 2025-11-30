package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesA = `100
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
57
58
59
60
61
62
63
64
65
66
67
68
69
70
71
72
73
74
75
76
77
78
79
80
81
82
83
84
85
86
87
88
89
90
91
92
93
94
95
96
97
98
99
100`

func expected(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(i))
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func parseTests() ([]int, error) {
	fields := strings.Fields(testcasesA)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}
	ns := make([]int, t)
	if len(fields)-1 < t {
		return nil, fmt.Errorf("not enough cases")
	}
	for i := 0; i < t; i++ {
		n, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, err
		}
		ns[i] = n
	}
	return ns, nil
}

func buildInput(ns []int) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(ns)))
	sb.WriteByte('\n')
	for _, n := range ns {
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ns, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	input := buildInput(ns)
	output, err := runCandidate(bin, input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != len(ns) {
		fmt.Printf("expected %d outputs, got %d\n", len(ns), len(lines))
		os.Exit(1)
	}
	for i, n := range ns {
		want := expected(n)
		if strings.TrimSpace(lines[i]) != want {
			fmt.Printf("case %d failed\nn=%d\nexpected: %q\ngot: %q\n", i+1, n, want, lines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(ns))
}

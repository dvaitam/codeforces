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
50 98
54 6
34 66
63 52
39 62
46 75
28 65
18 37
18 97
13 80
33 69
91 78
19 40
13 94
10 88
43 61
72 13
46 56
41 79
82 27
60 7
89 70
52 49
75 42
54 14
35 7
88 25
9 76
98 59
67 68
1 3
41 56
25 52
77 15
30 69
39 63
39 68
67 87
89 56
48 75
20 69
92 100
93 51
38 52
29 81
85 77
79 67
12 8
33 20
4 58
1 70
95 42
6 35
28 32
19 36
88 58
68 20
5 81
22 72
46 20
6 82
10 82
92 59
4 3
53 19
44 32
33 49
56 34
68 6
3 10
43 100
80 86
41 43
11 82
61 23
100 59
73 42
57 47
84 40
5 54
46 74
63 16
92 82
44 75
36 9
65 45
63 37
70 39
57 45
75 60
72 88
89 40
80 91
38 19
2 30
30 84
44 44
14 44
83 12
64 4
`

func solve(n, m int) string {
	if n >= m && (n-m)%2 == 0 {
		return "Yes"
	}
	return "No"
}

func parseTests(raw string) ([][2]int, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	pairs := make([][2]int, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("bad testcase line: %q", line)
		}
		n, err1 := strconv.Atoi(fields[0])
		m, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("bad ints in line: %q", line)
		}
		pairs = append(pairs, [2]int{n, m})
	}
	return pairs, nil
}

func buildInput(n, m int) string {
	return fmt.Sprintf("1\n%d %d\n", n, m)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		n, m := tc[0], tc[1]
		input := buildInput(n, m)
		want := strings.ToLower(solve(n, m))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.ToLower(strings.TrimSpace(got)) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

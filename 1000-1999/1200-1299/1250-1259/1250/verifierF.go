package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
31
84
67
64
1
97
17
63
83
80
9
1
40
57
93
35
26
66
70
18
81
88
44
40
18
24
100
42
52
91
38
45
58
62
54
23
77
59
27
34
98
54
84
56
63
8
44
63
44
81
7
42
64
47
83
81
57
9
12
96
69
37
84
38
80
56
44
84
83
59
15
45
87
58
83
53
70
5
51
70
45
82
1
46
100
3
76
61
65
79
55
29
88
92
8
40
66
69
25
78
`

func solve(n int64) int64 {
	best := int64(1<<63 - 1)
	for i := int64(1); i*i <= n; i++ {
		if n%i == 0 {
			j := n / i
			p := 2 * (i + j)
			if p < best {
				best = p
			}
		}
	}
	return best
}

func parseTestcases() ([]int64, error) {
	lines := strings.Split(strings.TrimSpace(testcasesData), "\n")
	vals := make([]int64, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	return vals, nil
}

func run(bin string, input string) (string, error) {
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, n := range testcases {
		input := fmt.Sprintf("%d\n", n)
		want := strconv.FormatInt(solve(n), 10)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

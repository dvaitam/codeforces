package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solve matches the original 1180A solution logic.
func solve(n int) int {
	return 2*n*(n-1) + 1
}

const testcaseData = `
100
50
98
54
6
34
66
63
52
39
62
46
75
28
65
18
37
18
97
13
80
33
69
91
78
19
40
13
94
10
88
43
61
72
13
46
56
41
79
82
27
71
62
57
67
34
8
71
2
12
93
52
91
86
81
1
79
64
43
32
94
42
91
9
25
73
29
31
19
70
58
12
11
41
66
63
14
39
71
38
91
16
71
43
70
27
78
71
76
37
57
12
77
50
41
74
31
38
24
25
24
`

func loadTestcases() ([]int, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	T, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("bad testcase count: %v", err)
	}
	if len(fields)-1 != T {
		return nil, fmt.Errorf("testcase count mismatch: declared %d, have %d", T, len(fields)-1)
	}
	tests := make([]int, 0, T)
	for i := 1; i < len(fields); i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("bad testcase %d: %v", i, err)
		}
		tests = append(tests, v)
	}
	return tests, nil
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int) string {
	return fmt.Sprintf("%d", solve(n))
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for tc, n := range testcases {
		input := fmt.Sprintf("%d\n", n)
		want := expected(n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", tc+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput: %sexpected: %s\ngot: %s\n", tc+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

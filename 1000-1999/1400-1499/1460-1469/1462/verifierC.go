package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// solve mirrors 1462C.go logic.
func solve(x int) string {
	if x > 45 {
		return "-1"
	}
	var digits []int
	for d := 9; d >= 1 && x > 0; d-- {
		if x >= d {
			digits = append(digits, d)
			x -= d
		}
	}
	sort.Ints(digits)
	var sb strings.Builder
	for _, d := range digits {
		sb.WriteString(strconv.Itoa(d))
	}
	if sb.Len() == 0 {
		return "0"
	}
	return sb.String()
}

type testCase struct {
	input string
}

// Embedded testcases from testcasesC.txt.
const testcaseData = `
100
16
3
41
18
17
25
8
23
19
24
40
41
13
48
35
35
4
22
33
36
1
11
32
2
41
17
22
9
23
7
40
36
27
49
5
17
37
7
34
5
38
9
32
3
34
15
23
7
36
46
14
12
45
43
15
26
13
39
45
48
29
23
11
18
47
23
36
13
14
30
21
9
18
36
43
50
50
1
6
43
21
20
46
21
3
6
40
1
19
8
11
46
18
13
33
36
23
29
6
27
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	res := make([]testCase, 0, t)
	if len(fields) != 1+t {
		return nil, fmt.Errorf("expected %d test values, got %d", t, len(fields)-1)
	}
	for i := 0; i < t; i++ {
		x, err := strconv.Atoi(fields[1+i])
		if err != nil {
			return nil, fmt.Errorf("bad value at %d: %v", i+1, err)
		}
		input := fmt.Sprintf("1\n%d\n", x)
		res = append(res, testCase{input: input})
	}
	return res, nil
}

func runCandidate(bin, input string) (string, error) {
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

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expected, err := runReference(tc.input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, tc.input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

// runReference evaluates solve via 1462C semantics.
func runReference(input string) (string, error) {
	fields := strings.Fields(input)
	if len(fields) < 2 {
		return "", fmt.Errorf("invalid input")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	if t != 1 {
		return "", fmt.Errorf("unexpected t %d", t)
	}
	x, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", err
	}
	return solve(x), nil
}

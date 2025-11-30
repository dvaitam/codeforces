package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt.
const testcasesAData = `
29 X-- ++X X-- X++ X-- ++X X-- --X X-- ++X X-- --X X++ X++ X-- X++ X-- X++ --X X++ X-- X-- ++X --X ++X --X X-- X++ X++
26 X-- ++X --X X++ X-- X-- ++X --X --X X++ X++ X-- X-- --X --X ++X --X X++ --X --X ++X --X ++X X-- ++X X--
2 X-- ++X
49 X-- --X X++ --X ++X X-- ++X X-- --X ++X X-- X++ --X ++X X++ ++X X-- X-- --X X++ X++ X-- X++ ++X ++X ++X X++ ++X X++ --X X-- ++X X-- ++X --X ++X X-- ++X --X --X X++ --X ++X ++X X++ ++X --X X++ X--
23 X-- ++X ++X X++ X++ X++ ++X ++X ++X ++X ++X X++ --X ++X X++ X-- --X ++X X-- --X X-- --X ++X
13 --X X-- ++X X++ X++ --X X++ ++X --X --X ++X X-- --X
13 --X X++ ++X --X --X X++ ++X ++X X-- --X X-- ++X --X
37 X++ X++ X++ ++X --X --X X++ ++X ++X X++ X++ X++ ++X ++X X-- X++ X-- --X ++X ++X X++ ++X X++ --X X++ ++X ++X X++ X-- --X X-- --X X-- ++X X++ --X --X
8 X++ X++ X-- --X --X X++ X++ --X
32 X++ ++X X-- ++X X-- ++X ++X --X ++X X-- ++X --X ++X X-- ++X X++ X-- X++ ++X ++X --X X-- --X X-- X++ X++ X++ --X X-- --X ++X X++
18 --X X++ ++X ++X X++ X-- ++X X++ --X X++ --X --X X++ X++ --X --X X++ --X
8 X++ --X --X --X ++X X-- X-- X--
43 X++ X-- --X ++X X-- ++X --X X++ X-- --X ++X X-- X-- ++X --X ++X X++ ++X X++ X-- ++X X-- ++X ++X X-- ++X --X X-- --X --X ++X X++ X-- X++ X++ X-- ++X X-- --X ++X X++ X-- X--
27 X-- --X X++ --X ++X ++X --X --X ++X X-- --X X++ --X ++X ++X ++X X++ ++X X++ X++ X++ X++ --X --X --X ++X X--
11 X++ ++X X++ X++ --X ++X X-- ++X --X ++X ++X
4 --X --X X-- ++X
34 --X --X X-- ++X X-- --X X++ --X X++ ++X X++ --X X-- ++X X++ --X ++X X-- X-- X-- ++X X-- X++ X-- --X X-- ++X X-- X++ X-- --X X-- --X --X
24 X-- --X ++X ++X X++ --X --X ++X ++X X-- X++ ++X ++X X++ X-- --X ++X ++X X-- X++ X-- X++ X++ --X
48 --X X++ --X X-- X++ X++ X-- X-- --X --X X++ ++X X-- ++X X++ X++ --X --X X++ ++X ++X X-- X++ X++ X++ ++X X-- X++ --X X-- --X --X ++X X-- --X ++X ++X X-- X++ X-- X-- X-- ++X --X X++ --X --X --X
1 X--
19 --X --X X++ X-- ++X X++ --X ++X X-- X-- ++X X++ --X --X X++ ++X --X ++X X--
4 ++X X++ X++ --X
32 ++X --X --X X-- --X --X --X ++X --X ++X --X --X --X --X X++ X++ X++ X-- X++ X-- --X X-- X-- X++ X++ --X X++ X++ X-- X-- --X X++
43 --X --X --X X++ X++ --X X++ X-- X++ ++X ++X ++X --X X-- ++X ++X X++ ++X X++ --X --X --X X++ X-- ++X X-- --X ++X X-- X++ --X --X X-- ++X X++ ++X X++ ++X --X ++X ++X X++ X++
14 X-- --X ++X --X X-- X-- ++X X-- ++X ++X ++X X++ ++X X++
40 X++ X-- ++X X++ --X --X X-- ++X X-- X-- X-- --X ++X ++X X++ ++X X-- X++ --X ++X --X --X --X X-- ++X --X X-- X++ ++X X-- ++X X-- X-- X-- ++X --X ++X ++X --X X--
12 X++ X++ X-- --X X++ X-- ++X X-- ++X X-- X-- ++X
3 --X --X --X
26 ++X --X X++ X-- X++ X-- ++X --X X-- ++X ++X X-- ++X ++X --X ++X ++X X++ ++X X-- --X --X --X X-- X-- X--
1 X++
25 ++X X-- --X --X --X ++X X-- --X --X X-- ++X ++X X++ X-- ++X ++X ++X X-- --X --X ++X ++X --X --X X--
26 X++ X++ X++ ++X X++ --X ++X X-- --X --X ++X X-- X++ X++ ++X X++ --X ++X ++X X-- X-- ++X X++ --X --X ++X
6 ++X --X X++ X++ X-- X++
16 X-- ++X ++X ++X --X X++ --X --X --X X-- ++X --X X-- X-- --X X++
20 X++ ++X ++X --X X++ ++X --X --X --X --X ++X --X --X ++X --X X-- --X --X --X X--
43 X-- X-- X-- X++ X-- X-- ++X --X X++ X-- X++ ++X ++X X-- ++X X-- X-- X++ X-- X-- ++X --X ++X X++ --X --X --X --X --X --X ++X ++X --X --X ++X --X X-- X++ X++ X++ --X X++ X--
17 X++ X++ --X ++X --X --X X-- --X --X --X X++ X++ X-- X++ X-- X-- X--
44 X-- X++ --X --X ++X X-- --X X-- --X ++X ++X X++ ++X ++X X-- --X X++ X-- ++X X++ ++X --X ++X ++X X-- ++X X++ X++ X-- X++ --X ++X ++X X-- ++X --X X++ --X X++ --X ++X --X --X --X
14 --X ++X X++ X++ X-- --X --X X-- X-- --X --X ++X X-- X--
39 X-- --X ++X --X X-- X-- --X --X X++ --X X-- --X X-- ++X ++X X-- --X ++X ++X X-- ++X --X --X ++X X-- ++X X++ ++X X-- ++X X-- X-- --X ++X X++ X-- X++ --X X++
43 X++ X++ ++X X++ --X ++X X++ X++ ++X ++X X-- ++X ++X --X X++ ++X X++ ++X --X X-- X++ X-- ++X ++X --X --X --X X-- X-- X++ X-- X-- ++X ++X --X ++X ++X --X ++X X-- --X --X ++X
30 X-- --X ++X X-- X-- ++X X-- X-- --X X++ ++X ++X X++ ++X --X X-- --X ++X X++ ++X --X X++ --X X-- --X X++ --X X-- X-- X--
25 X-- ++X X++ --X X++ X++ ++X --X ++X X-- --X ++X X++ X++ X++ ++X X-- --X X++ X-- ++X ++X ++X --X X++
16 X-- ++X X-- --X --X X++ X-- --X --X X++ X-- ++X X++ X++ --X ++X
19 ++X --X --X ++X X-- X++ X++ X++ ++X --X ++X X-- --X ++X ++X --X --X ++X X--
49 ++X X-- --X --X X++ X++ ++X ++X X++ ++X X-- --X ++X --X ++X ++X ++X X-- ++X ++X ++X X-- X-- --X X++ X++ X-- ++X X++ X-- X++ X-- --X --X X-- --X X++ X++ ++X --X --X --X X-- X-- X++ X-- ++X ++X --X
35 X++ X++ --X X-- --X X++ X-- --X X++ --X ++X ++X X++ --X ++X ++X X++ X-- X-- X-- --X X++ X-- X++ X++ ++X ++X X-- X++ ++X X++ X-- ++X X++ --X
16 --X X-- X-- ++X X++ X++ ++X --X X++ X-- --X X-- X++ X++ X-- --X
7 --X X++ X++ ++X ++X ++X X++
4 ++X X++ X-- --X
43 X++ ++X X++ ++X X-- X-- --X --X --X X-- --X X-- X-- X++ --X ++X X-- X-- X++ X++ X++ X-- X++ X++ --X X-- --X --X --X ++X X++ ++X ++X X++ --X ++X --X X++ --X ++X X-- X++ X++
5 ++X X-- --X X-- X++
36 --X --X ++X ++X X++ --X X++ ++X X++ ++X X-- --X ++X --X X++ --X ++X ++X X++ X++ --X X-- --X X++ X++ --X --X X-- X-- X-- X++ --X X-- X-- --X X++
38 --X ++X ++X --X X-- ++X X++ --X ++X --X --X X++ ++X X-- X-- ++X ++X --X X++ --X --X ++X --X --X X-- X-- --X ++X X-- ++X ++X X++ --X X++ X++ X-- ++X --X
28 X++ X-- X-- X++ --X ++X X-- ++X X++ ++X --X --X ++X ++X --X --X ++X ++X ++X --X X-- --X X++ ++X ++X X-- ++X X++
8 ++X ++X --X --X X++ X++ ++X X--
30 X-- X-- X++ X++ X++ X++ ++X ++X --X X++ --X --X X++ --X X++ X-- --X X-- --X X++ X-- --X --X --X ++X ++X X++ --X X-- X++
19 X-- X++ --X ++X --X X-- --X ++X X-- --X --X ++X X-- X-- X-- --X X-- X++ X--
34 X++ ++X X-- X++ ++X X-- X++ X++ X++ --X --X X++ --X --X X-- --X --X X++ X++ --X X-- X++ X++ X++ --X X++ ++X ++X X++ ++X ++X X-- X++ X--
19 X-- X-- ++X ++X X-- X++ ++X X-- X-- ++X --X ++X X-- X-- --X X-- --X X++ --X
31 --X X-- ++X ++X X-- X-- ++X ++X ++X --X X++ --X X++ X++ ++X ++X X-- X++ X-- X++ ++X X-- X++ X-- ++X X-- --X ++X --X --X X++
3 ++X --X X--
33 --X X-- --X X-- X++ X-- X++ X++ X-- ++X --X --X --X X-- --X X++ --X --X ++X ++X X++ X++ X-- ++X --X X-- X++ --X ++X X++ X-- X-- X--
28 --X X-- X-- X++ --X X-- --X ++X ++X X-- ++X ++X ++X ++X --X --X --X X-- X++ X++ X-- X-- X-- X++ --X --X X++ ++X
47 X++ X-- --X ++X X-- --X --X X-- ++X --X ++X X++ X++ --X ++X ++X --X X-- X-- X++ ++X --X --X --X X++ --X X++ X-- X-- X-- X-- X++ --X X-- --X --X X-- ++X X++ --X ++X --X ++X X++ --X --X X--
26 ++X X++ X++ --X ++X X-- --X --X X-- X-- --X X++ ++X X++ X-- ++X --X X-- X-- X++ X++ --X ++X X++ --X X--
36 ++X X++ --X --X X++ X++ X++ --X X-- --X X++ X++ X-- --X ++X --X ++X X-- ++X X++ --X ++X --X ++X X++ X-- X++ ++X X-- --X X++ ++X X-- X++ ++X X--
14 X++ X-- X-- ++X X++ --X ++X X++ ++X --X X++ ++X X++ --X
7 X-- X++ ++X X-- X++ ++X X--
27 X-- X-- --X X-- X++ --X X++ X-- X-- ++X ++X X++ --X X++ --X --X --X ++X --X --X X-- X-- ++X X++ --X --X X--
26 X-- X-- X-- X++ --X X-- X++ X++ --X X++ X-- --X X++ ++X X++ X++ X-- X++ --X X-- X-- --X X-- ++X X++ X++
47 X++ X-- --X X-- X-- X++ ++X --X X-- X++ --X ++X X++ --X ++X ++X X-- --X --X ++X X++ X++ --X --X X-- X++ X++ X++ --X X++ --X --X --X X-- --X ++X X-- X-- X-- --X X++ --X --X X++ X-- X++ --X
49 ++X X++ X-- X-- --X X-- X-- X-- X++ --X X-- X-- X++ X-- X-- X-- X-- X++ ++X ++X --X --X ++X X++ X-- X-- ++X X-- X++ X-- ++X --X X-- ++X X-- X++ ++X --X X-- X-- X-- ++X ++X --X --X X++ --X ++X X--
10 X-- --X X++ X++ X-- --X X++ --X --X --X
42 ++X X-- ++X --X X-- ++X ++X X++ X-- X-- X++ X++ ++X X++ X++ ++X ++X ++X X++ ++X ++X X++ --X X++ --X X++ --X ++X X-- ++X X++ X++ --X ++X X-- X-- X++ X-- X++ X-- ++X X++
29 --X --X X-- X++ --X ++X --X --X X++ X++ X-- X-- --X --X ++X ++X X-- X++ --X --X --X ++X X++ X-- X-- X++ ++X X++ --X
43 X-- X-- X-- ++X X++ X-- ++X ++X --X ++X X++ X++ X-- --X X-- --X X-- X-- --X X-- X++ ++X X-- ++X --X --X X++ X-- --X --X ++X --X --X X-- --X X-- --X ++X ++X X++ X++ X-- --X
37 ++X X++ ++X X-- --X X++ ++X --X --X --X X-- ++X X-- X-- --X ++X --X X++ X++ X++ ++X X-- X++ ++X --X X-- --X X-- ++X X++ --X --X ++X --X ++X X-- X--
39 ++X X-- X-- ++X X-- X-- --X X-- X-- ++X X++ X-- ++X --X X-- ++X X-- ++X --X --X ++X X-- ++X X++ ++X X-- X-- X++ ++X ++X X++ X-- X++ --X X++ --X X-- X++ ++X
13 ++X --X --X --X ++X X-- ++X --X X-- --X ++X --X X++
23 X++ ++X X-- ++X --X ++X --X X-- X-- ++X --X ++X X++ ++X X-- X-- ++X X-- ++X --X ++X X++ X++
6 X-- X++ --X --X X-- ++X
10 --X X-- --X --X X++ ++X --X X-- --X X--
3 --X --X ++X
34 ++X X-- ++X --X --X --X --X X++ X-- --X X++ ++X --X X++ X++ X++ --X --X --X X-- --X X++ ++X ++X ++X ++X X-- X-- X-- X-- ++X X-- X-- X--
40 --X --X ++X ++X X-- ++X X++ X-- ++X X-- --X X++ --X ++X X++ X-- --X ++X X-- X-- X++ ++X X++ ++X X++ ++X X++ ++X --X --X --X ++X X-- X++ X++ --X ++X ++X --X ++X
32 X++ X++ --X X++ X++ ++X X-- X++ ++X X-- X++ X-- --X X++ X++ ++X X-- ++X ++X ++X X-- X-- ++X ++X --X --X --X X++ X++ ++X X-- X++
46 X-- --X X++ X++ X-- ++X X-- X-- X++ --X ++X X-- X++ X-- --X X++ ++X ++X ++X --X ++X ++X --X X-- X++ ++X ++X ++X --X X-- X-- --X --X X++ X++ --X --X ++X X-- ++X X-- X++ X-- X++ --X X--
20 X-- --X --X X++ X++ X++ ++X --X ++X X-- X++ --X X-- X-- ++X X-- --X --X ++X X++
40 ++X X++ --X ++X ++X X++ X-- ++X ++X X-- X-- X-- --X X++ X-- X++ --X --X X-- --X --X --X ++X X++ X-- ++X X-- X++ X++ X-- ++X ++X ++X X-- X-- --X X++ --X --X X++
48 X++ --X X-- --X ++X --X X-- --X --X X-- X-- X++ X++ X-- --X X-- X++ X-- X++ X-- X++ X-- X++ X++ --X --X X++ ++X ++X --X X-- X-- X++ ++X X-- ++X --X --X X-- ++X X++ --X X++ X++ ++X ++X X-- --X
1 --X
16 ++X X++ X++ --X X++ X-- X++ --X X++ --X X++ ++X X-- X-- X-- ++X
17 ++X X-- --X X-- ++X X-- X-- --X X-- ++X X++ --X --X X-- X++ X++ ++X
5 --X --X --X X-- --X
27 X-- X-- ++X X++ --X X++ ++X --X --X X++ ++X ++X X++ X-- X-- X-- ++X X++ X-- ++X --X ++X --X X++ X++ --X X--
39 --X --X --X ++X X-- X-- X-- ++X X++ X-- --X --X ++X X++ ++X --X X-- X-- X++ X-- X-- X-- --X ++X --X --X --X X-- X++ X-- X++ ++X X++ X-- X++ X-- --X X++ X++
43 X++ X++ X-- --X ++X --X --X --X ++X X-- --X X++ X-- --X X++ --X --X X-- X-- ++X ++X --X X-- --X X-- X++ X-- ++X --X --X --X X-- X-- --X X++ ++X X-- X-- X-- ++X X-- ++X X--
20 X-- ++X X++ --X ++X --X --X X-- X-- X-- X++ --X ++X X-- --X X-- --X --X ++X X--
13 ++X --X --X X++ ++X ++X X++ ++X --X X-- X++ X-- X++
`

// expected computes the result following 282A.go logic.
func expected(ops []string) int {
	x := 0
	for _, op := range ops {
		if len(op) >= 2 && op[1] == '+' {
			x++
		} else {
			x--
		}
	}
	return x
}

type testCase struct {
	ops []string
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesAData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("line %d: expected %d ops, got %d", idx+1, n, len(fields)-1)
		}
		cases = append(cases, testCase{ops: fields[1:]})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (int, error) {
	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(tc.ops))
	for _, op := range tc.ops {
		input.WriteString(op)
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	valStr := strings.TrimSpace(out.String())
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("parse output %q: %v", valStr, err)
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := expected(tc.ops)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

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
13 14
2 9
17 16
13 10
16 12
19 7
17 5
10 5
4 20
9 18
20 5
10 4
3 11
16 18
4 12
14 11
20 7
18 16
15 17
9 2
18 1
3 13
1 20
16 11
8 11
3 7
19 8
8 5
18 15
3 3
11 17
16 4
10 18
10 4
18 11
18 7
20 18
19 10
15 3
20 13
11 19
8 10
6 7
6 2
20 9
16 3
3 5
5 2
3 18
13 17
9 17
8 7
19 14
19 9
15 16
12 3
11 20
4 16
19 11
7 8
1 9
4 8
12 6
11 14
2 4
5 8
2 19
18 20
3 1
4 7
20 19
4 13
3 12
4 2
20 1
7 6
4 16
7 2
1 18
14 20
4 9
3 8
3 10
12 14
6 2
17 15
2 20
4 13
7 9
12 16
19 6
7 2
6 6
11 17
9 4
20 15
6 1
16 14
19 17
10 12
`

type testCase struct {
	a int64
	b int64
}

func expected(tc testCase) string {
	a, b := tc.a, tc.b
	switch {
	case a == b:
		return "0"
	case a < b:
		d := b - a
		if d%2 == 1 {
			return "1"
		}
		return "2"
	default:
		d := a - b
		if d%2 == 0 {
			return "1"
		}
		return "2"
	}
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesData), "\n")
	cases := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid testcase line: %q", line)
		}
		a, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}
		b, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}
		cases = append(cases, testCase{a: a, b: b})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := fmt.Sprintf("1\n%d %d\n", tc.a, tc.b)
		want := expected(tc)
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

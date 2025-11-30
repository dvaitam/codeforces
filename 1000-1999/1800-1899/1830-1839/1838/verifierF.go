package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `100
10
5
6
9
1
8
4
1
3
2
6
8
4
7
9
2
10
4
1
4
7
5
3
7
3
2
3
10
10
8
3
3
1
1
4
4
3
3
5
6
4
9
4
3
4
7
5
1
6
7
3
3
5
2
6
5
10
10
1
10
6
2
5
6
5
8
6
3
8
8
3
1
5
1
6
7
1
9
7
6
7
10
1
8
1
3
10
4
2
4
8
6
9
6
9
5
8
2
10
6`

const expectedOutput = "1 1 ^"

type testCase struct {
	input    string
	expected string
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing value", i+1)
		}
		val := fields[pos]
		pos++
		cases = append(cases, testCase{
			input:    val + "\n",
			expected: expectedOutput,
		})
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
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

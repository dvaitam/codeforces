package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `
100
bca
bca
abc
bac
cab
bca
bca
bac
bca
bac
cab
acb
cab
acb
bac
acb
abc
cab
bac
cab
cba
cab
acb
bac
abc
cba
abc
cba
bac
bca
cab
abc
bac
bca
bac
cab
cba
acb
cab
bca
bca
cab
bac
abc
cab
abc
abc
cba
bca
cba
cba
cba
abc
cab
bca
bac
acb
cba
bac
cba
abc
acb
cab
acb
acb
acb
cab
bca
abc
abc
bac
cab
bca
abc
bac
cab
bac
cba
abc
cab
bac
cab
acb
cab
cab
cab
bac
bca
abc
cab
bca
bac
cab
acb
bac
acb
acb
acb
abc
cab
`

func solve(s string) string {
	target := "abc"
	diff := make([]int, 0, 2)
	for i := 0; i < 3; i++ {
		if s[i] != target[i] {
			diff = append(diff, i)
		}
	}
	ok := false
	if len(diff) == 0 {
		ok = true
	} else if len(diff) == 2 {
		i, j := diff[0], diff[1]
		if s[i] == target[j] && s[j] == target[i] {
			ok = true
		}
	}
	if ok {
		return "YES"
	}
	return "NO"
}

type testCase struct {
	input    string
	expected string
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	if len(fields) != t+1 {
		return nil, fmt.Errorf("expected %d cases, got %d", t, len(fields)-1)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		s := fields[i+1]
		input := fmt.Sprintf("1\n%s\n", s)
		cases = append(cases, testCase{
			input:    input,
			expected: solve(s),
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
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load cases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

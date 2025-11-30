package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcaseData = `100
m
y
n
b
i
q
p
m
z
j
p
l
s
g
q
e
j
e
y
d
t
z
i
r
w
z
t
e
j
d
x
c
v
k
p
r
d
l
n
k
t
u
g
r
p
o
q
i
b
z
r
a
c
x
m
w
z
v
u
a
t
p
k
h
x
k
w
c
g
s
h
h
z
e
z
r
o
c
c
k
q
p
d
j
r
j
w
d
r
k
r
g
z
t
r
s
j
o
c
t`

type testCase struct {
	input    string
	expected string
}

func solve(s string) string {
	if strings.ContainsRune("codeforces", rune(s[0])) {
		return "YES"
	}
	return "NO"
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	cases := make([]testCase, 0, len(lines)-1)
	for i, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			return nil, fmt.Errorf("case %d: empty line", i+1)
		}
		input := fmt.Sprintf("1\n%s\n", line)
		cases = append(cases, testCase{
			input:    input,
			expected: solve(line),
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
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

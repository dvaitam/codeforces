package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `13 )())))))(()((
20 )()(())()))()))((()(
20 ))()((((()(())())()(
20 ))())()(((())(((((()
17 )(()))))()())((()
4 ()()
14 (((((((((()()(
2 ((
6 ()((()
20 ()((()))(()(()()))((
2 ((
11 )()(())))))
5 ()()(
18 )(())))()))((()(((
15 ))())(()()())((
16 )))()(((()))((((
17 (((())(())(())()(
17 )(()(((()))())))((
7 ))((())
11 )()(()(()))
18 ()()(()(()))))(())
11 )()())())((
19 )(((((()(())))()()(
9 ()()((()(
4 )))(
2 ((
4 ()))
18 (()())))))((()))((
9 ()))(((((
5 ()())
6 ()))((
17 ()())))(((()))())
10 )))(((())(
2 ))
5 )(())
2 ))
15 (()(((()()())((
20 ))()())((()))())((()
9 (((())(()
19 )(()(()(()))()))()(
4 )())
6 (()(((
1 (
4 ())(
20 ((()()))((())))))(((
18 )()((()((()))((()(
3 ())
18 )()))(())())()))()
10 )((()(())(
14 )(())())()()()
20 )(()()())()())(()(()
13 ))()))(((()))
4 ((()
10 )((())))))
19 (((()))())))())))()
10 ()((()))((
3 )()
2 ((
9 )(()(()))
13 ()))()()((())
13 )))(()(()((((
12 ((()())))())
16 ))())(((()()((()
8 )()()(()
11 )))((()(()(
7 )))))))
8 ()(()))(
14 )))(((((()()()
18 )((()(())))()())))
9 ((((()(()
1 )
2 ()
9 ()()((())
8 )))((())
19 (()((((()(((())()((
17 ))))((()))()()())
13 ((()())())()(
12 (()()(((((()
3 ())
5 ())()
10 ()))()))((
14 )((((()(((((((
14 ()(()))(()())(
15 ()()))))()))()(
15 )()))(()())()))
7 ()()(()
1 (
19 ())))))))(((()()())
20 (())((((()))))()))))
7 )((()))
2 )(
1 )
17 (())((()()()()(((
15 ()))))))()()()(
13 )))())())())(
2 ()
12 (((()()())()
14 (())(((()(()()
19 ((()()))))()(((()((
15 ()()()(()((()))`

// Embedded reference logic from 1322A.go.
func solve(n int, s string) int {
	if n%2 == 1 {
		return -1
	}
	count := 0
	for _, ch := range s {
		if ch == '(' {
			count++
		} else {
			count--
		}
	}
	if count != 0 {
		return -1
	}

	balance := 0
	start := -1
	ans := 0
	for i, ch := range s {
		if ch == '(' {
			balance++
		} else {
			balance--
		}
		if balance < 0 && start == -1 {
			start = i
		}
		if balance == 0 && start != -1 {
			ans += i - start + 1
			start = -1
		}
	}
	return ans
}

type testCase struct {
	n int
	s string
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n", idx+1)
		}
		res = append(res, testCase{n: n, s: fields[1]})
	}
	return res, nil
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test data:", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expected := fmt.Sprintf("%d", solve(tc.n, tc.s))
		input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

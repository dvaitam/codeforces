package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcases = `
1 1 1
6 2 6
5 3 5
4 1 2
7 6 4
9 6 9
8 5 1
1 1 1
6 4 4
9 3 9
3 1 1
1 1 1
3 1 3
9 6 9
9 3 8
7 6 5
6 5 3
6 4 2
7 6 6
8 4 8
5 4 5
9 6 8
8 6 8
8 4 6
3 3 2
8 5 5
9 9 9
9 7 5
4 4 3
10 2 6
1 1 1
1 1 1
10 4 2
9 3 5
4 2 1
7 6 7
1 1 1
6 2 2
1 1 1
2 1 1
1 1 1
3 1 3
3 3 3
1 1 1
4 2 1
1 1 1
5 3 4
1 1 1
9 1 5
7 7 5
3 2 1
2 2 1
1 1 1
9 7 8
9 6 3
6 3 3
10 7 1
9 3 1
5 1 2
3 1 1
8 4 1
4 2 4
2 2 1
10 4 10
10 6 5
7 3 5
1 1 1
7 4 2
2 1 1
2 1 1
3 1 1
4 1 4
8 5 7
4 2 4
7 5 1
10 10 1
7 5 5
3 1 3
8 6 1
9 2 6
5 3 3
1 1 1
2 2 1
1 1 1
7 6 4
8 4 2
1 1 1
6 3 6
2 1 2
4 1 3
7 6 4
3 2 2
2 2 1
2 1 2
7 2 6
2 1 2
1 1 1
6 4 2
6 3 4
9 8 7

`

type testCase struct {
	n int
	a int
	b int
}

func parseCases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcases))
	scan.Split(bufio.ScanWords)
	cases := []testCase{}
	for {
		if !scan.Scan() {
			break
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("parse n: %w", err)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d: missing a", len(cases)+1)
		}
		a, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("parse a: %w", err)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d: missing b", len(cases)+1)
		}
		b, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("parse b: %w", err)
		}
		cases = append(cases, testCase{n: n, a: a, b: b})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func referenceSolve(n, a, b int) string {
	x, y := -1, -1
	for i := 0; i <= n/a; i++ {
		rem := n - a*i
		if rem%b == 0 {
			x = i
			y = rem / b
			break
		}
	}
	if y == -1 {
		return "-1"
	}
	idx := 1
	var sb strings.Builder
	for cnt := 0; cnt < x; cnt++ {
		l := idx
		r := idx + a - 1
		for j := l + 1; j <= r; j++ {
			sb.WriteString(strconv.Itoa(j))
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(l))
		sb.WriteByte(' ')
		idx += a
	}
	for cnt := 0; cnt < y; cnt++ {
		l := idx
		r := idx + b - 1
		for j := l + 1; j <= r; j++ {
			sb.WriteString(strconv.Itoa(j))
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(l))
		sb.WriteByte(' ')
		idx += b
	}
	return strings.TrimSpace(sb.String())
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := fmt.Sprintf("%d %d %d\n", tc.n, tc.a, tc.b)
		expected := referenceSolve(tc.n, tc.a, tc.b)
		out, stderr, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

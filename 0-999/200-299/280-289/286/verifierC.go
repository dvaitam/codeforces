package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesC = `4 2 1 4 2 0 
1 5 1 1
5 1 2 1 3 3 3 2 1 3
8 1 5 1 4 2 3 3 4 2 4 7
1 2 0 
6 5 3 1 5 4 2 0 
8 4 5 5 3 3 4 3 2 8 1 4 6 8 3 7 5 2
3 2 4 3 2 3 2
10 5 2 3 4 4 1 1 5 2 3 2 4 10
8 4 5 4 1 4 5 4 1 2 8 1
5 2 4 5 4 5 4 1 5 2 3
5 4 1 4 2 5 5 1 2 5 4 3
6 1 2 1 1 5 5 0 
4 1 5 2 3 2 2 1
8 4 1 1 3 4 1 3 2 8 6 1 2 3 7 4 5 8
5 5 3 3 5 1 5 5 4 3 2 1
6 5 2 2 4 5 3 0 
3 2 3 3 2 2 1
6 5 1 1 3 2 2 4 3 6 4 5
3 3 1 4 1 1
5 2 5 1 3 4 2 3 4
2 1 5 1 2
6 3 1 4 1 4 4 0 
5 3 2 2 5 4 5 1 5 4 3 2
1 4 0 
2 4 5 2 2 1
8 5 2 4 1 3 2 3 5 2 7 2
6 1 1 1 5 4 2 0 
8 4 3 2 1 2 5 2 1 3 8 4 3
9 2 1 5 4 2 5 4 4 5 7 6 8 4 9 2 5 1
6 3 3 1 5 2 3 4 2 4 3 5
8 1 1 5 1 1 2 2 1 4 1 7 4 3
3 2 4 3 3 3 1 2
9 5 1 4 2 3 5 5 4 4 6 4 1 6 8 2 3
9 5 3 3 1 4 3 3 4 4 6 1 3 6 2 7 8
6 1 1 4 4 2 4 6 5 1 2 4 3 6
10 4 4 4 1 1 4 2 1 1 5 9 3 6 2 9 5 10 8 4 7
2 1 5 2 2 1
2 5 3 0 
7 3 1 5 2 1 4 4 5 7 2 4 3 5
2 1 1 1 2
1 5 0 
4 1 5 5 4 4 3 1 4 2
10 4 1 1 4 1 1 4 2 1 4 6 7 1 8 10 3 6
2 3 1 0 
6 1 3 3 2 1 2 6 3 1 2 5 6 4
2 1 3 1 1
10 2 2 2 4 1 4 3 3 2 1 3 6 10 8
5 3 5 3 2 5 0 
2 5 5 1 1
7 2 2 2 3 5 2 2 2 3 7
7 1 2 5 1 4 1 1 2 4 3
9 4 2 5 4 3 3 1 2 4 5 9 1 4 7 8
7 3 4 2 3 3 4 1 2 7 1
5 1 5 5 2 4 3 2 4 5
3 2 4 3 1 2
8 1 4 2 3 1 4 5 4 0 
4 3 1 3 5 4 2 4 3 1
8 2 5 4 3 1 1 3 1 0 
5 4 5 5 4 4 0 
5 3 3 2 5 1 0 
2 3 3 2 2 1
9 2 2 1 4 3 3 5 2 5 8 4 2 9 6 5 7 3 8
5 4 3 5 2 2 0 
2 4 4 2 2 1
9 3 3 4 4 2 4 4 5 3 7 1 8 3 2 4 9 7
4 1 3 4 4 0 
9 1 1 4 1 3 1 1 5 1 4 5 4 2 9
5 2 1 4 4 3 3 2 3 5
7 2 4 2 5 3 2 2 2 4 3
7 4 4 4 2 2 4 2 0 
7 1 2 1 2 3 1 2 3 5 3 7
2 5 3 1 2
8 1 5 5 4 5 4 4 3 7 4 3 7 1 5 8 6
6 1 3 1 2 1 4 5 2 5 4 3 1
8 2 3 5 1 5 1 3 3 8 8 3 7 1 5 4 2 6
4 5 3 2 2 2 3 2
9 4 3 4 3 2 1 3 5 1 0 
8 4 4 1 4 4 4 4 1 1 2
4 1 2 4 2 3 1 2 3
1 2 0 
8 2 2 3 3 3 4 1 5 4 4 6 3 8
9 5 4 5 3 3 2 1 1 5 1 3
7 2 2 3 1 5 5 4 0 
2 4 3 0 
10 3 2 5 3 2 2 1 5 3 3 3 6 8 5
10 2 2 1 5 5 3 3 5 1 2 6 3 10 2 8 7 4
10 2 2 2 2 4 3 5 5 2 4 1 10
1 5 1 1
8 3 1 2 5 2 4 4 5 5 2 3 8 5 4
4 3 3 4 1 1 1
6 2 3 4 2 3 3 5 2 3 1 6 4
3 3 1 4 0 
1 2 0 
1 2 1 1
1 3 1 1
4 4 4 2 3 2 2 3
6 4 4 1 4 3 5 4 6 4 1 3
2 4 4 0 
1 5 0 `

// Embedded solver from 286C.go.
func solveCase(n int, p []int, forced []int) (bool, []int) {
	isForcedClose := make([]bool, n)
	for _, q := range forced {
		if q >= 1 && q <= n {
			isForcedClose[q-1] = true
		}
	}
	res := make([]int, n)
	stack := make([]int, 0, n)
	for i := n - 1; i >= 0; i-- {
		if isForcedClose[i] {
			stack = append(stack, p[i])
			res[i] = -p[i]
		} else {
			if len(stack) > 0 && stack[len(stack)-1] == p[i] {
				stack = stack[:len(stack)-1]
				res[i] = p[i]
			} else {
				stack = append(stack, p[i])
				res[i] = -p[i]
			}
		}
	}
	if len(stack) != 0 {
		return false, nil
	}
	return true, res
}

type testCase struct {
	n   int
	p   []int
	t   int
	qes []int
}

func parseCases() ([]testCase, error) {
	data := strings.TrimSpace(testcasesC)
	if data == "" {
		return nil, fmt.Errorf("no testcases provided")
	}
	lines := strings.Split(data, "\n")
	cases := make([]testCase, 0, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d malformed", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n", i+1)
		}
		if len(fields) < 1+n+1 {
			return nil, fmt.Errorf("line %d: not enough fields", i+1)
		}
		p := make([]int, n)
		for j := 0; j < n; j++ {
			val, err := strconv.Atoi(fields[1+j])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad p[%d]", i+1, j+1)
			}
			p[j] = val
		}
		tIdx := 1 + n
		t, err := strconv.Atoi(fields[tIdx])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad t", i+1)
		}
		if len(fields) != 1+n+1+t {
			return nil, fmt.Errorf("line %d: expected %d fields got %d", i+1, 1+n+1+t, len(fields))
		}
		qes := make([]int, t)
		for j := 0; j < t; j++ {
			val, err := strconv.Atoi(fields[tIdx+1+j])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad q[%d]", i+1, j+1)
			}
			qes[j] = val
		}
		cases = append(cases, testCase{n: n, p: p, t: t, qes: qes})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(tc.t))
	sb.WriteByte('\n')
	for i, q := range tc.qes {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(q))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func expectedOutput(tc testCase) string {
	ok, res := solveCase(tc.n, tc.p, tc.qes)
	if !ok {
		return "NO"
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := buildInput(tc)
		expect := expectedOutput(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

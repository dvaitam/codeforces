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
100
2 5 2
2 1 7
2 4 10
4 2 2
2 1 6
2 1 7
3 2 10
1 2 0
1 1 10
1 2 10
1 2 0
1 2 7
1 2 3
1 2 4
1 2 8
1 1 10
2 1 5
5 4 9
1 3 4
2 4 9
1 4 3
2 4 10
1 3 8
2 1 7
1 2 8
2 3 7
1 4 0
3 5 10
2 2 2
1 1 3
1 2 8
2 5 5
2 3 10
1 2 8
1 3 8
1 2 0
2 3 9
1 3 6
4 3 7
2 1 8
2 2 9
1 2 10
1 2 1
2 1 10
1 1 0
2 1 4
2 3 2
1 2 4
1 1 2
3 5 3
2 3 7
2 4 7
1 1 4
4 3 7
1 3 1
2 3 8
1 4 0
1 1 6
1 1 2
2 3 8
2 3 3
5 4 4
1 4 10
2 4 0
2 2 3
1 3 1
1 3 5
1 1 9
2 1 0
1 1 9
2 1 9
1 1 3
3 1 4
2 1 7
1 3 6
2 1 0
2 1 4
1 2 4
2 1 5
2 1 4
1 1 8
2 2 8
2 1 1
1 1 2
2 5 4
2 3 9
2 3 5
2 1 4
1 2 2
5 1 6
1 4 1
2 1 2
2 1 9
2 1 9
1 5 1
2 1 4
5 1 8
2 1 0
2 1 9
1 1 6
1 1 3
1 5 6
1 1 7
1 2 2
1 4 6
5 3 9
2 3 7
2 1 3
2 1 0
1 3 9
2 2 6
2 2 1
1 3 9
2 1 4
1 5 8
4 3 5
1 2 4
1 2 5
1 3 1
2 1 10
2 1 6
3 1 6
1 2 9
2 1 5
1 3 9
1 1 3
1 1 6
1 2 8
1 1 1
1 1 5
4 4 3
1 3 1
1 2 2
1 3 4
1 5 10
2 2 3
1 1 5
1 1 4
2 5 2
1 1 4
1 1 6
2 5 7
2 1 6
2 2 4
2 1 10
4 5 1
1 3 9
2 5 3
1 2 4
2 5 6
1 1 3
4 1 3
2 1 10
1 2 5
2 1 3
4 3 9
2 3 3
1 1 8
2 1 8
1 3 4
2 3 5
1 4 9
1 1 9
2 1 2
2 2 3
5 1 8
2 1 6
1 5 0
1 3 10
1 3 1
1 5 10
1 4 3
2 1 6
1 3 7
2 5 8
1 1 6
2 1 10
2 3 3
2 5 0
1 2 9
1 1 10
1 2 3
1 2 2
5 2 5
2 2 10
2 1 8
2 2 6
1 2 9
2 1 4
1 1 2
1 1 10
1 1 8
3 5 5
2 5 10
2 5 5
1 1 7
2 3 4
2 3 10
5 4 2
2 4 3
1 3 10
1 5 8
2 3 2
2 5 10
1 1 8
1 1 9
2 4 5
1 1 3
2 1 6
1 1 4
2 1 10
1 2 4
2 1 2
2 1 7
1 2 8
1 2 8
1 2 1
2 1 10
2 1 2
1 1 6
2 1 3
5 2 4
2 2 1
1 5 10
2 2 8
1 2 4
5 3 6
1 4 8
2 1 7
2 3 5
1 3 9
1 1 9
2 2 6
2 3 4
1 1 9
2 3 2
2 2 8
1 1 2
4 3 5
2 1 1
1 3 1
1 2 6
2 2 1
1 1 0
5 1 4
1 4 8
2 1 10
2 1 9
1 1 3
4 2 8
2 2 2
1 2 4
2 2 3
2 2 5
2 1 3
1 1 0
1 4 5
2 2 3
4 2 3
1 1 6
1 1 9
2 2 2
1 4 5
1 1 8
1 1 0
2 1 6
1 1 0
2 2 4
2 4 7
2 3 4
1 1 0
1 2 6
1 2 8
2 2 8
2 1 4
1 2 2
1 2 1
1 1 0
1 1 9
2 1 1
2 1 2
1 1 10
1 1 7
1 1 1
2 1 1
2 1 6
1 1 5
1 1 6
1 3 2
2 1 8
1 1 5
5 4 10
2 1 2
2 1 4
1 2 5
2 3 1
2 3 2
1 1 4
2 4 4
2 3 4
2 1 8
1 2 5
3 3 10
1 2 4
2 2 5
2 1 9
1 1 0
2 3 4
1 3 9
2 2 10
2 2 4
2 3 5
1 1 2
3 2 10
1 1 2
2 1 1
2 1 3
2 1 10
1 1 3
1 3 6
1 3 5
2 2 3
1 3 7
1 2 7
3 5 4
2 1 4
2 2 0
2 5 7
1 2 9
5 1 1
2 1 0
2 3 1
1 2 8
3 5 10
2 5 6
2 5 10
2 4 4
1 3 7
1 3 2
2 1 6
1 2 6
2 3 10
1 1 1
1 2 4
4 3 6
2 2 6
2 1 7
2 1 6
1 1 2
2 2 2
2 2 4
5 3 7
2 2 5
2 1 7
2 3 6
1 1 2
1 2 3
1 1 4
1 4 1
4 2 1
1 4 9
1 5 4
2 3 0
1 1 10
1 1 10
2 2 7
1 2 2
2 1 10
2 2 10
5 4 7
1 5 7
1 2 6
1 2 8
2 4 8
2 4 10
1 5 5
2 1 0
3 3 1
2 2 1
2 5 5
2 2 6
1 1 4
1 2 8
1 1 6
2 3 7
3 3 8
1 2 5
2 1 5
1 3 2
2 3 2
1 3 5
1 3 3
2 3 7
2 2 1
2 2 5
1 1 10
1 1 10
1 1 9
1 2 6
2 1 0
3 5 4
1 3 3
2 3 4
2 1 1
2 3 2
1 3 3
1 1 1
1 1 4
2 1 4
5 1 6
1 1 2
2 1 10
1 1 10
2 1 0
2 1 5
1 5 4
4 1 10
2 1 8
2 1 3
2 1 0
1 2 3
1 4 4
1 3 8
2 1 7
1 4 1
2 1 10
2 1 9
3 4 3
1 1 9
1 3 3
2 3 5
3 2 3
2 2 6
1 3 2
2 2 10
5 1 9
1 2 6
1 4 0
2 1 4
2 1 6
2 1 1
2 1 10
1 1 1
1 5 8
2 1 9
3 4 4
1 1 8
2 2 1
1 3 5
2 3 4
1 5 7
2 4 5
2 3 7
1 1 0
2 1 8
1 1 5
1 1 7
1 1 10
2 4 4
1 1 5
2 2 10
2 4 5
2 2 6
4 4 9
1 4 4
1 2 0
2 4 1
1 1 2
2 4 10
2 2 2
1 3 0
2 4 10
1 4 0
5 2 7
1 2 5
1 1 8
1 2 6
1 5 3
2 1 0
1 5 9
1 2 6
4 1 10
1 4 1
1 4 0
1 1 0
2 1 4
2 1 9
1 3 8
2 1 8
1 4 6
1 4 4
2 1 4
5 3 3
1 3 5
1 3 4
2 2 6
3 5 8
1 1 2
2 2 3
1 3 8
1 3 6
1 3 2
2 4 3
1 3 1
1 3 0
1 4 7
2 2 9
1 1 3
2 2 4
1 1 5
1 1 4
1 1 7
2 1 8
3 2 8
1 2 9
1 3 5
2 2 9
1 1 5
1 1 10
1 1 6
1 3 3
2 1 3
4 5 3
2 1 1
1 4 10
2 5 9
2 3 1
1 1 6
3 4 7
2 4 1
1 1 8
1 2 0
2 2 6
2 2 1
1 3 0
2 4 3
2 5 9
1 2 8
2 2 3
2 5 1
2 1 7
1 1 2
2 4 0
1 2 1
2 5 9
2 4 4
3 4 1
2 1 6
3 5 6
1 3 9
2 1 9
2 4 6
1 3 9
1 1 9
1 1 5
3 3 9
1 3 5
1 2 10
1 3 7
2 3 8
1 1 5
2 1 2
1 3 1
2 2 8
2 2 5
3 5 6
1 3 1
1 2 6
2 5 9
1 1 2
2 4 1
1 2 8
3 2 4
1 2 7
1 3 8
2 1 8
1 3 5
3 3 9
1 1 7
2 3 0
1 2 6
1 3 0
2 1 3
1 1 9
1 3 8
1 3 4
2 1 0
3 4 6
2 2 10
1 1 7
1 1 4
1 1 3
1 3 2
1 3 3
2 4 5
1 2 8
2 3 6
1 1 9
1 1 10
2 4 9
4 3 1
2 1 1
5 1 7
2 1 1
2 1 8
2 1 9
2 1 2
2 1 9
2 1 4
2 1 9
1 4 8
2 2 8
2 2 7
2 2 6
2 1 10
1 1 0
2 1 4
2 1 5
2 3 9
1 2 2
2 1 10
1 1 4
4 5 6
1 1 10
1 3 4
2 4 6
2 1 3
2 2 9
1 2 10
2 2 7
2 1 9
1 2 5
1 2 7
2 1 10
1 1 5
2 2 8
2 2 9
1 5 5
2 4 0
2 1 10
2 1 8
1 1 6
2 2 0

`

type op struct {
	typ   int
	idx   int
	color int
}

type testCase struct {
	n, m int
	ops  []op
}

func parseCases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcases))
	scan.Split(bufio.ScanWords)
	nextInt := func() (int, error) {
		if !scan.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(scan.Text())
		if err != nil {
			return 0, err
		}
		return v, nil
	}

	t, err := nextInt()
	if err != nil {
		return nil, fmt.Errorf("read test count: %w", err)
	}
	cases := make([]testCase, 0, t)
	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		n, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read n: %w", caseIdx, err)
		}
		m, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read m: %w", caseIdx, err)
		}
		k, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read k: %w", caseIdx, err)
		}
		ops := make([]op, k)
		for i := 0; i < k; i++ {
			typ, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d op %d: read typ: %w", caseIdx, i+1, err)
			}
			idx, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d op %d: read idx: %w", caseIdx, i+1, err)
			}
			color, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d op %d: read color: %w", caseIdx, i+1, err)
			}
			ops[i] = op{typ: typ, idx: idx, color: color}
		}
		cases = append(cases, testCase{n: n, m: m, ops: ops})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func referenceSolve(tc testCase) string {
	n, m := tc.n, tc.m
	rTime := make([]int, n+1)
	rColor := make([]int, n+1)
	cTime := make([]int, m+1)
	cColor := make([]int, m+1)

	for t, op := range tc.ops {
		time := t + 1
		if op.typ == 1 {
			rTime[op.idx] = time
			rColor[op.idx] = op.color
		} else {
			cTime[op.idx] = time
			cColor[op.idx] = op.color
		}
	}

	var sb strings.Builder
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if j > 1 {
				sb.WriteByte(' ')
			}
			var val int
			if rTime[i] > cTime[j] {
				val = rColor[i]
			} else {
				val = cColor[j]
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return out.String(), errBuf.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.m, len(tc.ops))
		for _, op := range tc.ops {
			fmt.Fprintf(&input, "%d %d %d\n", op.typ, op.idx, op.color)
		}

		expected := referenceSolve(tc)
		out, stderr, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\ngot:\n%s\n", idx+1, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

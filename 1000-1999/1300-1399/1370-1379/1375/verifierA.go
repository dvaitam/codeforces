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
100
9 3 -9 -2 6 5 2 -1 5 1
5 6 -6 -1 -6 -7
7 7 9 -6 -1 -7 -8 0
9 7 -7 1 3 0 9 10 -4 7
9 4 6 -2 -9 7 -10 -8 2 10
3 9 5 0
5 0 -8 -4 8 -3
5 -6 7 4 -8 -8
7 6 5 -7 -1 7 -1 -7
7 7 -4 9 7 8 -1 4
3 9 2 0
5 -1 -5 -4 -5 -9
7 5 -8 -8 -6 -6 -9 -8
9 6 -2 6 -3 -4 8 3 8 -2
9 5 10 1 -8 0 9 -7 5 8
7 -4 -3 -10 -2 -7 -3 1
5 0 3 -9 -7 -6
5 -9 8 10 7 9
3 -10 -7 10
5 9 8 -7 2 -8
7 -7 -9 9 -10 -4 -5 -7
9 -4 -9 -10 7 3 9 -7 -2 -8
5 -8 10 -1 1 3
5 -9 6 4 -9 9
3 2 -4 -2
7 5 8 -5 -4 -9 -5 -5
7 6 -2 -7 9 4 -5 -10
9 3 8 6 -1 10 1 2 -2 -6
3 4 -8 0
3 7 -2 -6
5 5 1 9 -1 1
5 -1 2 3 10 -8
3 9 -4 0
5 -3 -3 10 4 2
9 -9 2 8 3 -9 -5 4 -8 -2
5 4 6 5 7 9
3 -9 5 0
7 4 -9 3 -4 7 10 -8
5 -10 2 3 0 -10
5 -10 -10 6 9 -7
5 -7 9 10 -4 -1
7 -5 -7 5 2 10 -8 -10
7 4 -7 -2 -6 10 6 10
7 -7 -6 -2 -10 -9 -9 -4
7 7 0 1 8 -9 9 10
9 10 4 10 3 1 7 -5 -4 2
7 -10 -6 -6 -2 0 0 1
3 0 9 -9
3 -2 -5 -6
7 1 2 7 -6 -1 -7 5
5 -9 -1 -5 6 -8
7 2 0 -1 3 -7 -7 7
9 5 0 0 -7 5 -7 5 3 -9
7 0 -6 -5 10 8 2 10
3 -8 -8 -4
5 -9 2 -10 -7 2
7 4 5 8 -4 3 -8 1
5 -2 8 -5 3 -4
7 -7 -8 -10 6 4 -4 -7
9 2 -2 -4 10 -9 -4 9 -6 -7
5 4 2 1 7 -6
3 9 5 -6
9 10 3 6 5 0 5 5 10 -4
5 -10 0 0 0 -9
5 -2 9 -6 2 8
7 5 -8 -8 6 -9 -8 -3
5 -9 -1 -10 4 0
5 -6 10 4 1 6
9 6 6 -9 8 -8 6 9 -8 3
5 -1 7 9 3 5
9 9 8 -3 -10 -10 -5 -1 6 8
7 0 -8 5 -2 -1 3 2
9 -9 -5 10 -6 -3 -1 0 -9 -9
9 3 -6 5 9 -8 -6 1 3 -9
9 2 4 -9 -7 5 -6 -10 -9 9
5 10 0 -7 7 10
7 -4 2 5 -7 -9 9 4
7 10 -7 9 -1 -6 2 -1
3 6 -4 -9
9 4 1 -4 4 1 10 -8 -9 -9
9 -2 -10 6 8 8 -4 -3 -8 10
9 6 -1 -7 -6 3 8 3 -8 -7
9 -8 -7 3 -6 -10 4 3 3 -10
9 0 -2 -8 1 -8 -7 1 -10 1
7 -5 -10 -3 1 -8 9 -6
5 -10 -4 -7 -10 -1
7 -10 9 -3 -6 -5 4 -7
9 1 -2 -6 -10 -4 1 0 5 -1
7 7 10 0 -5 8 -8 -7
7 -5 2 -6 -6 -3 0 6
5 -3 -5 -1 1 3
3 -6 9 -10
9 -8 -8 -6 3 -1 7 3 -6 8
9 -1 10 1 -8 -3 4 10 1 10
3 2 3 -10
9 0 4 -4 1 -1 5 -8 -5 -7
7 -7 7 9 -6 4 2 -5
9 3 -5 -3 4 0 6 -6 1 4
3 5 -4 -1
3 4 9 4
`

type testCase struct {
	n   int
	arr []int
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

// solve replicates 1375A logic: make even indices positive, odd indices negative.
func solve(tc testCase) []int {
	res := make([]int, tc.n)
	for i, x := range tc.arr {
		if i%2 == 0 {
			if x < 0 {
				x = -x
			}
		} else {
			if x > 0 {
				x = -x
			}
		}
		res[i] = x
	}
	return res
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, err
	}
	pos++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("unexpected EOF at case %d", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, err
		}
		pos++
		if pos+n > len(fields) {
			return nil, fmt.Errorf("case %d missing numbers", i+1)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[pos+j])
			if err != nil {
				return nil, err
			}
			arr[j] = v
		}
		pos += n
		cases = append(cases, testCase{n: n, arr: arr})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("unexpected extra data after %d cases", t)
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		want := solve(tc)
		gotStr, err := run(bin, buildInput(tc))
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		outFields := strings.Fields(gotStr)
		if len(outFields) != tc.n {
			fmt.Printf("case %d failed: expected %d numbers got %d\n", idx+1, tc.n, len(outFields))
			os.Exit(1)
		}
		for i := 0; i < tc.n; i++ {
			val, err := strconv.Atoi(outFields[i])
			if err != nil {
				fmt.Printf("case %d invalid output %q\n", idx+1, outFields[i])
				os.Exit(1)
			}
			if val != want[i] {
				fmt.Printf("case %d failed at index %d: expected %d got %d\n", idx+1, i, want[i], val)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `2 1 1 3 2 2 1 2 1 1
14 6 8 5 8 7 8 4 0 0 5 7 5 6 6 2 4 5 1 2
6 5 2 2 8 8 5 2 4 5 3 6
12 9 5 5 7 2 6 7 8 3 7 4 7 3 4 11 9 12 4 11
8 5 2 9 4 7 4 4 8 4 2 5 1 8 6 6 4 4
4 0 9 0 4 2 1 4 1 4
9 3 3 0 6 0 0 5 5 2 2 2 2 2 2
1 0 1 1 1
5 2 2 8 0 6 1 2 2
2 0 5 1 1 2
16 0 4 7 8 9 0 4 6 9 2 7 3 1 5 1 0 4 9 10 1 16 9 16 3 6
11 4 4 9 6 0 8 2 0 4 0 2 2 2 3 2 9
2 3 3 4 2 2 1 1 1 2 1 2
1 2 1 1 1
6 1 8 1 3 1 1 1 2 2
4 3 0 8 7 4 3 4 1 4 2 2 1 4
14 8 0 9 9 0 6 8 9 2 1 7 5 0 8 1 5 8
12 4 0 6 1 1 4 3 0 7 0 6 7 4 10 11 1 1 1 4 5 8
3 3 7 3 1 2 3
15 2 5 6 1 4 1 1 1 9 5 6 3 1 0 9 4 12 12 5 12 8 11 13 14
12 4 7 8 7 6 7 4 6 3 2 7 9 3 1 8 2 2 3 6
5 6 1 1 0 2 3 2 3 2 5 2 3
10 1 2 8 6 1 5 8 3 8 4 2 8 9 7 8
12 9 2 7 7 0 9 6 2 6 8 0 7 3 3 10 4 11 9 12
11 1 3 8 9 3 6 6 0 5 7 8 4 2 3 7 7 10 11 2 9
4 6 8 3 4 2 3 4 1 1
14 7 4 8 9 2 7 3 1 5 0 7 8 1 9 4 8 11 9 12 1 8 10 10
12 2 6 4 2 0 2 7 6 7 4 2 0 3 1 8 1 4 5 12
15 3 4 7 2 7 8 4 1 4 5 4 5 4 6 8 1 7 8
5 8 1 4 0 3 4 1 4 2 5 1 1 4 4
12 3 5 5 1 5 7 5 2 7 7 4 7 2 2 9 6 9
6 1 3 7 3 5 2 3 2 2 3 3 3 6
13 6 5 4 9 8 9 5 6 4 8 9 1 5 3 4 11 5 6 8 11
16 1 2 5 6 2 0 1 5 2 5 1 6 0 8 5 3 4 1 16 3 10 6 9 15 16
16 1 2 3 5 4 2 6 5 4 1 5 3 3 3 9 0 3 11 14 1 16 3 4
3 6 7 4 2 1 2 2 3
8 0 6 7 7 9 5 8 9 1 1 8
15 2 6 6 8 7 0 1 7 9 2 1 8 2 1 6 3 1 8 2 5 4 7
6 0 2 6 1 5 7 1 2 3
3 7 2 8 1 3 3
2 0 3 1 1 2
5 5 7 0 2 8 1 1 1
15 3 0 9 3 6 5 9 6 8 8 2 8 1 2 3 2 4 11 6 9
14 2 6 2 6 5 4 1 8 1 7 4 4 8 7 3 7 8 12 13 1 1
7 3 3 6 9 0 2 0 3 4 7 1 4 2 5
5 9 5 0 3 1 2 1 4 4 4
10 3 2 5 4 8 9 1 6 6 0 4 6 9 5 5 4 4 2 9
9 8 6 9 8 3 6 2 2 7 4 4 7 2 9 8 9 1 8
16 9 5 4 1 2 5 9 7 3 9 9 2 4 3 8 4 1 1 1
7 5 0 5 8 4 5 7 1 4 5
1 4 2 1 1 1 1
3 9 7 4 1 2 3
8 2 9 4 3 3 9 5 9 4 1 8 5 6 5 5 2 3
13 6 1 7 6 6 7 6 4 3 3 3 0 8 1 1 1
13 6 6 3 8 4 1 5 8 5 8 7 9 1 4 5 6 1 1 1 8 11 12
5 3 5 3 8 0 2 2 5 5 5
3 2 6 2 1 2 3
16 0 8 5 5 1 9 5 1 9 5 5 4 7 4 8 9 2 2 2 7 10
1 5 1 1 1
14 6 3 2 2 9 0 0 9 5 2 4 0 0 3 2 1 8 2 5
3 3 3 8 2 1 1 1 2
9 9 3 0 8 8 8 6 6 2 2 1 8 2 9
2 9 2 4 1 2 2 2 2 2 1 1
9 2 2 6 9 4 7 7 7 3 4 2 9 6 7 7 7 7 7
10 0 8 7 9 4 4 3 7 7 5 4 9 10 8 9 7 10 4 7
4 8 3 6 1 4 1 4 1 4 4 4 3 4
12 8 5 2 2 3 2 6 7 7 2 6 4 3 3 10 6 9 7 7
2 3 7 1 1 1
12 5 9 5 6 2 5 1 8 9 7 6 9 2 1 8 6 6
1 4 4 1 1 1 1 1 1 1 1
7 3 9 5 2 5 4 6 4 3 4 5 5 5 5 2 3
12 5 2 3 8 9 3 3 9 7 3 6 4 2 1 2 4 11
12 2 3 9 2 7 9 0 2 3 3 0 9 1 2 9
6 4 6 9 0 0 5 1 1 2
12 7 6 1 1 7 2 2 7 1 4 2 7 1 5 8
1 8 2 1 1 1 1
10 7 3 0 5 7 0 2 3 9 5 3 4 4 4 5 6 6
4 4 1 5 1 2 1 2 3 4
7 4 8 0 0 5 2 1 1 6 6
16 6 5 7 4 4 1 2 1 0 1 2 5 2 4 9 8 2 8 8 1 16
11 1 2 6 3 6 0 5 6 7 8 5 2 11 11 3 3
13 2 6 3 4 9 8 1 1 8 9 9 2 3 1 3 6
8 3 7 6 7 1 5 3 0 3 7 7 2 2 3 6
8 5 6 3 0 2 2 2 2 4 1 2 4 7 1 8 7 7
1 1 1 1 1
1 3 4 1 1 1 1 1 1 1 1
10 9 6 7 7 6 4 6 2 6 4 4 3 10 5 8 10 10 5 5
14 9 7 1 8 8 7 8 4 2 3 5 1 3 3 4 9 10 9 10 7 14 2 3
10 7 2 8 1 5 2 3 8 0 6 1 3 10
10 3 2 6 3 3 9 3 8 4 8 3 2 5 10 10 3 10
2 0 7 2 1 2 2 2
11 2 7 0 3 5 4 7 4 8 4 0 3 7 7 3 4 4 11
5 6 4 1 2 0 4 5 5 1 4 1 1 2 3
16 6 9 9 5 1 0 1 3 0 8 8 0 2 0 4 4 1 1 1
12 4 1 6 2 6 0 3 3 2 9 6 8 1 2 5
8 9 3 0 3 3 2 9 0 2 4 4 3 4
11 1 3 2 9 3 1 5 7 2 3 3 4 1 4 4 4 6 9 5 8
2 4 1 4 2 2 1 2 1 2 1 1
5 1 7 1 9 2 4 2 2 2 3 3 4 1 4
14 3 1 5 5 4 6 2 3 9 4 7 5 0 0 1 4 4`

// Embedded reference logic from 1189C.go.
func buildDP(s []uint8) ([][]uint8, [][]int) {
	n := len(s)
	maxk := bits.Len(uint(n)) - 1
	dpVal := make([][]uint8, maxk+1)
	dpCarry := make([][]int, maxk+1)
	dpVal[0] = make([]uint8, n)
	dpCarry[0] = make([]int, n)
	copy(dpVal[0], s)
	for k := 1; k <= maxk; k++ {
		length := 1 << k
		half := 1 << (k - 1)
		size := n - length + 1
		if size <= 0 {
			dpVal[k] = nil
			dpCarry[k] = nil
			continue
		}
		dpVal[k] = make([]uint8, size)
		dpCarry[k] = make([]int, size)
		prevVal := dpVal[k-1]
		prevCarry := dpCarry[k-1]
		for i := 0; i < size; i++ {
			l := prevVal[i]
			r := prevVal[i+half]
			sum := int(l) + int(r)
			c := 0
			if sum >= 10 {
				c = 1
				sum -= 10
			}
			dpVal[k][i] = uint8(sum)
			dpCarry[k][i] = prevCarry[i] + prevCarry[i+half] + c
		}
	}
	return dpVal, dpCarry
}

func answerQueries(dpCarry [][]int, queries [][2]int) []int {
	res := make([]int, len(queries))
	for i, q := range queries {
		l := q[0] - 1
		length := q[1] - q[0]
		k := bits.TrailingZeros(uint(length + 1))
		res[i] = dpCarry[k][l]
	}
	return res
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

type testCase struct {
	n       int
	digits  []int
	queries [][2]int
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		pos := 0
		n, err := strconv.Atoi(parts[pos])
		if err != nil || n <= 0 {
			return nil, fmt.Errorf("line %d: invalid n", idx+1)
		}
		pos++
		if len(parts) < pos+n+1 {
			return nil, fmt.Errorf("line %d: missing digits/queries", idx+1)
		}
		digits := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[pos+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid digit", idx+1)
			}
			digits[i] = v
		}
		pos += n
		q, err := strconv.Atoi(parts[pos])
		if err != nil || q < 0 {
			return nil, fmt.Errorf("line %d: invalid q", idx+1)
		}
		pos++
		if len(parts) != pos+2*q {
			return nil, fmt.Errorf("line %d: query count mismatch", idx+1)
		}
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			l, err1 := strconv.Atoi(parts[pos])
			r, err2 := strconv.Atoi(parts[pos+1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("line %d: invalid query values", idx+1)
			}
			queries[i] = [2]int{l, r}
			pos += 2
		}
		res = append(res, testCase{n: n, digits: digits, queries: queries})
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test data:", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for i, d := range tc.digits {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(d))
		}
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "%d\n", len(tc.queries))
		for _, q := range tc.queries {
			fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
		}

		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		outLines := strings.Fields(out)
		if len(outLines) != len(tc.queries) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d lines got %d\n", idx+1, len(tc.queries), len(outLines))
			os.Exit(1)
		}

		digitsBytes := make([]uint8, len(tc.digits))
		for i, v := range tc.digits {
			digitsBytes[i] = uint8(v)
		}
		_, dpCarry := buildDP(digitsBytes)
		expected := answerQueries(dpCarry, tc.queries)

		for i := 0; i < len(tc.queries); i++ {
			val, err := strconv.Atoi(outLines[i])
			if err != nil || val != expected[i] {
				fmt.Fprintf(os.Stderr, "test %d query %d failed: expected %d got %s\n", idx+1, i+1, expected[i], outLines[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

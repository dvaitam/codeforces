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
3
9 1 4
2
7 7
8
6 3 1 7 0 6 6 9
1
7
5
3 9 1 5 0
1
0
9
0 6 3 6 0 8 3 7 7
9
3 5 3 3 7 4 0 6 8
2
2 4
2
5 8
7
8 3 4 4 9 7 8
7
9 0 7 3 6 6 2
6
8 5 1 7 8 1
3
8 6 5
8
0 7 0 4 9 9 9 6
3
2 8 3
1
3
9
8 3 6 8 5 9 5 7 4
9
9 0 6 8 2 8 8 3 6
1
7
6
9 8 3 8 6 7
6
6 5 0 8 8 9
10
5 7 9 0 3 2 8 9 2 1
9
4 0 1 1 0 7 0 4 3
5
1 9 2 5 4
2
2 2
5
8 2 4 4 7
6
7 7 1 0 4 6
6
6 3 4 1 4 8
4
9 6 0 3
1
6
3
0 2 7
9
6 8 3 8 7 3 8 0 6
10
5 6 0 4 2 3 0 4 1 1
5
4 2 6 9 4
3
0 8 0
10
3 9 7 2 9 8 0 6 3 5
2
3 9
7
9 3 7 1 6 4 8
8
0 5 9 6 4 0 2 3
6
9 2 5 6 3 4
2
6 8
6
8 7 8 3 1 0
2
2 2
3
8 3 4
6
9 8 4 5 5 5
2
4 3
10
7 2 9 8 1 5 0 6 1 6
3
2 5 1
10
9 6 1 9 8 3 9 1 4 5
5
9 8 1 7 4
2
0 4
1
9
1
1
7
1 0 3 3 9 6 2
2
7 2
4
2 1 6 6
9
4 8 4 7 5 1 3 5 0
1
0
5
9 5 7 6 5
7
1 1 5 9 7 1 4
4
9 8 7 5
5
2 8 3 4 3
4
5 1 4 1
8
1 9 5 3 6 4 0 5
3
5 9 4
4
5 1 8 9
10
9 1 3 3 0 3 6 1 4 8
2
1 0
1
4
6
7 7 2 1 8 5
2
8 2
3
2 2 5
5
1 8 9 4 2
4
2 8 0 5
10
8 3 2 4 6 8 2 0 3 4
2
7 6
9
4 8 7 8 7 0 6 5 2
5
7 0 6 9 0
1
5
10
2 9 2 2 4 4 6 9 6 2
10
1 3 7 0 2 8 5 8 7 3
4
5 7 7 3
7
5 8 9 4 3 0 1
9
5 2 8 3 4 4 4 8 5
3
7 9 1
2
9 8
10
6 2 2 4 6 3 9 0 7 6
6
6 8 2 8 0 8
2
4 1
5
1 2 9 1 7
4
6 6 6 2
6
7 2 9 7 3 1
7
9 8 6 1 4 4 3
7
8 0 3 8 7 9 0
1
9
4
4 3 2 4
3
8 3 4
5
9 4 7 2 8
6
7 6 1 3 9 6
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

func solveCase(arr []int) int {
	mismatchEven := 0
	mismatchOdd := 0
	for i, v := range arr {
		if i%2 != v%2 {
			if i%2 == 0 {
				mismatchEven++
			} else {
				mismatchOdd++
			}
		}
	}
	if mismatchEven == mismatchOdd {
		return mismatchEven
	}
	return -1
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
		// Allow trailing whitespace/numbers? treat as error to catch malformed embedding.
		return nil, fmt.Errorf("unexpected extra data after %d cases", t)
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	var inputBuilder strings.Builder
	inputBuilder.WriteString(strconv.Itoa(len(testcases)))
	inputBuilder.WriteByte('\n')
	for _, tc := range testcases {
		inputBuilder.WriteString(strconv.Itoa(tc.n))
		inputBuilder.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				inputBuilder.WriteByte(' ')
			}
			inputBuilder.WriteString(strconv.Itoa(v))
		}
		inputBuilder.WriteByte('\n')
	}
	input := inputBuilder.String()

	buildCaseInput := func(tc testCase) string {
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

	for idx, tc := range testcases {
		want := solveCase(tc.arr)
		gotStr, err := run(bin, buildCaseInput(tc))
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
		if err != nil || got != want {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, want, strings.TrimSpace(gotStr))
			os.Exit(1)
		}
	}

	// Also run full batch to ensure handling multiple cases works.
	batchOut, err := run(bin, input)
	if err != nil {
		fmt.Printf("batch run failed: %v\n", err)
		os.Exit(1)
	}
	outFields := strings.Fields(batchOut)
	if len(outFields) != len(testcases) {
		fmt.Printf("batch output count mismatch: expected %d got %d\n", len(testcases), len(outFields))
		os.Exit(1)
	}
	for i, f := range outFields {
		got, err := strconv.Atoi(f)
		if err != nil || got != solveCase(testcases[i].arr) {
			fmt.Printf("batch case %d mismatch: expected %d got %s\n", i+1, solveCase(testcases[i].arr), f)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

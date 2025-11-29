package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n   int
	arr []int
}

const testcasesRaw = `100
9 9 4 5 6 8 1 7 3 2
9 8 5 1 2 6 4 9 3 7
7 6 1 7 4 3 2 5
9 2 3 7 8 9 6 1 5 4
6 6 5 4 2 3 1
3 1 2 3
5 3 1 2 4 5
3 3 1 2
2 1 2
8 8 6 2 5 3 1 4 7
3 2 3 1
6 2 4 1 6 5 3
4 4 3 2 1
6 4 1 3 6 2 5
8 5 1 2 4 8 3 6 7
4 2 4 3 1
8 1 6 8 3 2 4 5 7
8 6 1 2 3 5 7 8 4
9 8 7 9 2 1 4 6 3 5
9 4 2 3 6 8 5 9 1 7
2 1 2
2 2 1
2 1 2
6 1 4 3 5 2 6
5 5 3 2 4 1
8 7 1 5 2 6 4 8 3
10 1 5 7 10 9 8 2 4 3 6
8 5 1 3 6 8 2 7 4
8 3 5 4 7 8 6 2 1
4 1 4 3 2
3 1 2 3
7 5 2 3 6 1 7 4
9 7 9 4 1 6 3 2 8 5
4 4 2 1 3
2 2 1
3 3 1 2
8 6 7 5 1 3 2 8 4
7 6 1 2 7 4 5 3
10 8 9 7 3 5 2 10 6 4 1
6 1 2 4 5 3 6
6 5 6 1 2 3 4
10 5 3 10 9 7 6 4 8 1 2
9 3 5 4 6 9 7 1 2 8
5 5 3 2 1 4
2 2 1
5 3 5 1 2 4
7 2 5 6 7 1 4 3
6 1 2 3 6 5 4
6 6 4 3 5 1 2
4 2 3 1 4
6 2 3 5 4 1 6
8 4 2 7 3 8 6 1 5
5 4 5 3 1 2
7 7 1 5 6 4 3 2
4 4 2 3 1
7 3 7 4 2 1 6 5
4 2 3 4 1
9 3 1 5 8 7 2 6 9 4
9 8 6 9 4 3 7 5 1 2
2 2 1
10 4 2 7 6 1 10 5 9 8 3
5 4 1 2 5 3
5 2 5 4 3 1
7 6 1 5 3 2 4 7
6 2 4 5 6 3 1
7 2 3 6 5 7 4 1
2 2 1
7 2 5 6 4 3 7 1
9 1 5 4 6 2 9 8 3 7
8 2 6 8 7 5 3 4 1
8 4 6 3 8 5 1 2 7
3 3 1 2
4 4 3 1 2
10 9 6 10 4 1 5 8 7 2 3
2 1 2
7 2 5 3 1 4 6 7
9 7 5 1 4 9 8 3 6 2
4 1 3 2 4
9 6 2 8 4 9 5 3 1 7
2 1 2
8 6 8 5 1 7 3 4 2
2 1 2
6 1 4 5 3 2 6
5 4 2 5 3 1
4 3 1 4 2
2 2 1
7 1 4 2 5 6 7 3
8 8 4 6 7 5 2 1 3
7 2 7 5 3 4 6 1
3 1 2 3
6 4 3 1 6 5 2
3 1 2 3
8 8 2 6 4 1 7 3 5
6 4 1 2 3 6 5
8 6 4 8 2 5 1 7 3
7 2 6 3 7 5 1 4
2 1 2
10 2 5 10 1 9 7 8 6 4 3
5 5 1 4 2 3
9 3 9 6 2 1 5 7 8 4`

func parseTestcases(raw string) []testCase {
	fields := strings.Fields(raw)
	if len(fields) < 1 {
		panic("no testcase data")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		panic(fmt.Sprintf("bad t: %v", err))
	}
	fields = fields[1:]
	tests := make([]testCase, 0, t)
	idx := 0
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if idx >= len(fields) {
			panic("unexpected end of data")
		}
		n, err := strconv.Atoi(fields[idx])
		if err != nil {
			panic(fmt.Sprintf("bad n at case %d: %v", caseIdx+1, err))
		}
		idx++
		if idx+n > len(fields) {
			panic(fmt.Sprintf("case %d missing array elements", caseIdx+1))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[idx+i])
			if err != nil {
				panic(fmt.Sprintf("bad value at case %d pos %d: %v", caseIdx+1, i+1, err))
			}
			arr[i] = v
		}
		idx += n
		tests = append(tests, testCase{n: n, arr: arr})
	}
	if idx != len(fields) {
		panic("extra data after parsing testcases")
	}
	return tests
}

// Embedded solver logic from 1364B.go.
func solve(tc testCase) []int {
	p := tc.arr
	ans := []int{p[0]}
	for i := 1; i < tc.n-1; i++ {
		if (p[i]-p[i-1])*(p[i+1]-p[i]) < 0 {
			ans = append(ans, p[i])
		}
	}
	ans = append(ans, p[tc.n-1])
	return ans
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := parseTestcases(testcasesRaw)
	for i, tc := range tests {
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j, v := range tc.arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		expectedArr := solve(tc)
		expectedStr := fmt.Sprintf("%d\n", len(expectedArr))
		for j, v := range expectedArr {
			if j > 0 {
				expectedStr += " "
			}
			expectedStr += strconv.Itoa(v)
		}
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(gotLines) != 2 {
			fmt.Printf("test %d failed: expected 2 lines got %d\n", i+1, len(gotLines))
			os.Exit(1)
		}
		if strings.TrimSpace(gotLines[0]) != strings.Fields(expectedStr)[0] || strings.TrimSpace(gotLines[1]) != strings.Join(strings.Fields(expectedStr)[1:], " ") {
			fmt.Printf("test %d failed\ninput: %sexpected:\n%s\ngot:\n%s\n", i+1, input.String(), strings.TrimSpace(expectedStr), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

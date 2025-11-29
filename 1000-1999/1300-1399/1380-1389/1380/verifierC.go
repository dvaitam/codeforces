package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesData = `
100
1 2
2
6 3
5 5 10 4 10 1
10 3
7 7 9 6 9 8 9 5 1 1
6 8
6 7 7 9 3 9
3 4
4 1 3
6 3
3 9 9 6 9 9
3 8
7 9 6
10 6
6 8 3 7 8 9 4 8 5 8
9 9
6 8 8 6 10 9 8 8 4
6 3
10 5 8 5 5 9
9 9
9 10 10 7 5 4 8 9 6
10 2
6 1 4 2 1 10 1 5 10 4
2 9
3 5
4 4
1 7 1 1
6 6
3 4 1 2 2 2
1 1
1
6 5
3 3 3 9 1 7
10 1
4 3 1 1 6 10 2 5 6 8
1 5
8
9 10
1 5 7 10 3 8 4 2 6
2 1
8 3
9 10
7 8 9 6 3 6 5 5 10
7 1
9 3 1 5 1 3 3
3 2
8 4 9
1 4
4
8 2
5 2 10 4 10 10 6 5
7 5
9 1 3 1 7 7 3
2 9
2 4
2 2
1 3
4 2
4 1 9 8
8 5
9 7 4 4 7 7 9 1
10 10
1 7 9 10 3 2 8 6 1 9
2 10
6 5
6 5
1 7 2 2 5 4
1 8
1
7 8
8 4 10 10 2 1 5
1 6
5
2 4
8 4
2 10
6 7
8 3
6 7 2 5 2 2 2 10
6 7
4 2 1 10 8 1
8 5
6 8 3 6 5 8 9 8
7 8
5 7 4 3 8 10 5
9 7
2 10 10 2 2 6 3 9 3
7 2
2 1 3 5 7 4 6
8 3
9 5 2 3 9 7 2 6
9 4
9 5 3 3 8 4 7 6 10
3 8
8 1 10
7 3
7 9 1 8 5 7 5
7 8
6 9 6 2 4 9 10
4 7
7 1 6 8
9 8
3 2 1 7 4 10 10 7 4
2 7
9 4
5 10
10 4 8 10 3
1 10
7
8 5
9 10 3 8 4 2 6 1
8 9
2 10 8 6 8 5 9 8
1 2
10
6 3
7 5 3 1 3 8
7 8
5 3 1 5 9 8 1
6 1
9 7 10 8 4 5
8 3
8 9 5 2 5 6 5 6
5 7
9 2 9 4 7
10 9
3 9 2 5 1 4 8 9 4 9
5 1
2 2 7 6 4
6 6
2 6 8 6 3 8
8 5
8 3 8 4 5 6 3 2
4 8
4 6 3 6
3 3
4 5 9
7 7
6 5 10 9 10 6 7
5 9
10 2 6 5 7
8 3
5 6 8 8 2 3 6 7
3 1
2 6 3
6 2
7 1 9 6 4 10
7 9
5 8 3 6 6 4 8
2 3
4 6
5 3
7 6 5 2 6
4 4
4 10 1 6
6 10
1 3 3 2 7 8
5 3
6 9 10 2 6
10 7
4 1 7 8 8 10 6 9 10 10
2 10
9 9
8 7
8 3 7 7 9 8 1 2
8 10
3 2 9 3 2 7 5 8
1 5
2
6 4
3 1 3 7 2 6
8 1
8 4 2 8 3 9 1 3
9 9
1 1 4 9 1 9 6 9 4
3 6
8 1 3
9 2
4 2 8 4 1 10 4 7 6
10 7
9 9 3 9 2 3 4 3 7 4
5 6
7 3 7 3 7
6 5
2 9 2 8 5 5
9 8
5 4 7 3 9 2 1 10 9
4 4
4 7 10 1
3 1
5 8 9
1 4
3
10 6
1 4 2 3 9 3 2 8 5 4
3 6
5 9 10
2 7
7 1
`

type testCase struct {
	n   int
	x   int
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

func solve(tc testCase) int {
	a := append([]int(nil), tc.arr...)
	sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
	teams := 0
	size := 0
	for _, v := range a {
		size++
		if size*v >= tc.x {
			teams++
			size = 0
		}
	}
	return teams
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
		if pos+2 > len(fields) {
			return nil, fmt.Errorf("unexpected EOF at case %d", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, err
		}
		x, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, err
		}
		pos += 2
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
		cases = append(cases, testCase{n: n, x: x, arr: arr})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("unexpected extra data after %d cases", t)
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.x))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		want := solve(tc)
		val, err := strconv.Atoi(strings.TrimSpace(got))
		if err != nil || val != want {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, want, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

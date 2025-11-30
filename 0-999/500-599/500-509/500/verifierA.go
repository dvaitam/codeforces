package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesA = `100
14 14
7 1 5 9 8 7 7 3 4 3 3 1 1
11 4
2 5 3 3 1 1 3 2 1 1
15 7
10 11 4 9 8 8 5 1 5 1 1 3 2 1
17 12
8 12 6 12 2 4 10 4 4 7 2 5 4 1 1 1
18 17
4 10 9 5 12 2 9 6 9 4 7 5 5 3 2 1 1
12 11
4 5 3 4 7 2 1 3 2 1 1
6 3
1 1 3 2 1
18 9
7 14 10 5 8 8 11 6 2 6 5 1 4 3 1 1 1
10 3
4 6 7 2 3 4 1 1 1
9 2
2 1 1 2 1 2 1 1
5 2
1 1 1 1
17 8
2 15 11 1 9 7 10 2 5 1 2 1 3 2 2 1
3 3
1 1
14 5
5 6 8 10 3 4 7 1 2 2 2 2 1
16 12
3 1 8 11 7 10 9 5 6 3 4 3 1 1 1
4 3
3 1 1
6 3
4 3 3 2 1
20 6
10 13 14 3 1 10 4 12 6 3 4 4 6 4 4 4 1 2 1
3 2
2 1
10 4
8 8 5 5 1 1 2 2 1
16 2
13 14 13 7 4 9 2 3 1 4 4 3 1 1 1
2 2
1
5 3
3 2 1 1
17 14
3 1 5 8 2 5 3 9 6 1 2 3 1 1 1 1
10 10
6 6 5 1 5 4 3 2 1
13 10
3 4 7 5 1 2 2 3 3 2 2 1
12 11
1 1 5 3 2 5 3 3 2 1 1
5 5
2 1 2 1
18 4
10 13 14 6 5 7 2 2 9 8 4 3 3 1 2 1 1
15 2
5 6 12 11 3 3 7 7 6 1 1 1 1 1
3 3
1 1
14 10
9 5 8 8 4 7 1 3 2 3 3 1 1
8 4
1 1 1 4 3 1 1
17 14
9 4 11 1 4 10 3 2 4 4 4 3 2 1 2 1
20 14
14 17 16 11 14 8 8 11 11 4 9 4 1 3 3 3 1 1 1
6 5
5 3 3 2 1
4 4
1 1 1
6 2
3 1 2 2 1
6 5
3 4 3 1 1
18 4
14 7 5 9 10 7 8 7 4 1 6 1 2 3 3 2 1
4 3
2 2 1
14 14
7 1 3 3 4 5 6 3 1 1 2 2 1
17 4
5 13 6 7 1 10 8 7 8 1 1 4 2 1 1 1
12 3
9 6 4 7 7 4 1 1 3 2 1
5 4
2 2 2 1
18 8
2 13 8 6 13 4 8 6 2 1 1 4 3 1 3 1 1
4 4
3 2 1
5 3
4 3 2 1
5 5
1 1 2 1
2 2
1
15 2
8 6 12 5 2 6 2 1 3 1 3 2 1 1
9 7
2 5 2 2 1 1 1 1
11 7
1 4 3 2 4 1 4 2 2 1
2 2
1
12 9
5 5 9 6 2 5 1 1 3 2 1
14 4
3 4 6 9 4 4 7 2 3 3 2 1 1
2 2
1
4 2
2 2 1
6 6
4 3 3 2 1
9 9
6 6 5 1 4 2 1 1
12 9
4 6 5 8 1 2 1 3 1 1 1
14 4
13 7 7 3 4 8 3 5 2 3 2 1 1
8 4
1 6 4 4 1 1 1
5 4
2 2 2 1
17 9
13 5 11 1 2 5 1 1 5 4 5 5 4 2 1 1
13 6
11 4 10 2 1 1 3 3 3 1 1 1
4 3
2 2 1
20 18
7 18 4 14 11 9 7 12 5 5 8 6 5 6 2 2 1 1 1
14 11
8 3 9 5 6 8 6 4 2 4 2 2 1
3 3
2 1
17 3
7 1 6 8 7 1 9 2 2 6 6 4 1 2 1 1
2 2
1
9 4
5 2 1 4 4 3 2 1
7 4
4 4 2 2 1 1
6 3
2 4 2 2 1
17 14
8 13 4 8 4 10 1 7 1 2 6 1 2 2 1 1
9 6
2 6 5 3 3 2 2 1
18 15
15 16 5 12 8 4 6 5 1 1 1 2 3 1 2 1 1
4 3
3 1 1
19 9
15 7 11 10 2 10 2 6 6 9 8 3 3 1 1 1 2 1
8 8
5 3 2 2 2 2 1
18 14
9 16 6 14 12 4 1 5 9 2 1 4 4 4 1 2 1
16 9
2 2 2 4 2 3 7 4 4 5 1 4 3 2 1
7 3
4 2 2 2 2 1
15 3
9 5 10 9 4 5 8 5 5 4 3 2 1 1
5 2
2 3 2 1
8 4
6 6 1 4 1 1 1
10 3
6 4 6 6 5 3 1 1 1
18 11
11 8 6 11 8 5 10 3 3 1 5 5 3 3 3 1 1
14 4
3 9 2 3 4 8 5 6 2 2 3 1 1
14 7
10 10 3 8 2 1 5 5 3 4 2 2 1
9 4
8 7 6 4 3 3 1 1
6 6
4 2 2 2 1
3 2
1 1
9 7
8 6 6 2 3 2 1 1
2 2
1
6 4
1 4 3 1 1
9 2
1 2 6 3 1 1 2 1
6 3
4 4 1 2 1
7 7
3 4 4 1 2 1`

// Embedded solver from 500A.go.
func solveCase(n, t int, arr []int) string {
	pos := 1
	for pos < t {
		pos += arr[pos-1]
	}
	if pos == t {
		return "YES"
	}
	return "NO"
}

type testCase struct {
	n   int
	t   int
	arr []int
}

func parseCases() ([]testCase, error) {
	fields := strings.Fields(testcasesA)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	idx := 0
	readInt := func() (int, error) {
		if idx >= len(fields) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(fields[idx])
		idx++
		return v, err
	}
	tCount, err := readInt()
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	cases := make([]testCase, 0, tCount)
	for i := 0; i < tCount; i++ {
		n, err1 := readInt()
		t, err2 := readInt()
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("case %d: bad n or t", i+1)
		}
		arr := make([]int, n-1)
		for j := 0; j < n-1; j++ {
			val, err := readInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: bad arr[%d]", i+1, j+1)
			}
			arr[j] = val
		}
		cases = append(cases, testCase{n: n, t: t, arr: arr})
	}
	return cases, nil
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
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solveCase(tc.n, tc.t, tc.arr)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.t)
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

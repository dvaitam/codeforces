package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solve mirrors 1499C.go for one test case.
func solve(costs []int64) int64 {
	n := len(costs)
	minOdd := costs[0]
	minEven := costs[1]
	sumOdd := costs[0]
	sumEven := costs[1]
	cntOdd := 1
	cntEven := 1
	ans := (int64(n)-int64(cntOdd))*minOdd + (int64(n)-int64(cntEven))*minEven + sumOdd + sumEven
	for i := 2; i < n; i++ {
		v := costs[i]
		if i%2 == 0 {
			sumOdd += v
			cntOdd++
			if v < minOdd {
				minOdd = v
			}
		} else {
			sumEven += v
			cntEven++
			if v < minEven {
				minEven = v
			}
		}
		cur := (int64(n)-int64(cntOdd))*minOdd + (int64(n)-int64(cntEven))*minEven + sumOdd + sumEven
		if cur < ans {
			ans = cur
		}
	}
	return ans
}

type testCase struct {
	n   int
	arr []int64
}

// Embedded testcases from testcasesC.txt.
const testcaseData = `
100
2
6 6
7
11 48 43 20 17 39 14
2
38 44
4
28 41 26 47
10
24 35 29 33 18 3 2 24 30 21
8
28 34 11 36 12 16 15 2
4
21 12 9 33
10
24 33 44 36 12 29 27 48 34 49
7
38 23 24 29 11 49 26
9
42 34 16 32 18 32 33 33 23
9
30 23 37 47 36 47 30 32 43
5
21 45 11 40 18
9
20 20 46 33 36 34 33 42 40
8
20 47 14 32 33 24 44 40
3
22 47 1
5
48 7 4 37 42
2
18 38
5
44 7 49 34 9
6
16 14 4 28 46 49
2
4 24
7
12 16 44 2 6 8 5
2
3 47
2
24 17
4
11 48 12 34
2
25 38
2
16 10
2
1 23
3
19 22 32
2
20 29
10
50 39 48 3 17 49 26 40 46 10
9
15 6 43 44 21 7 2 29 9
10
38 50 26 32 33 21 10 22 17 17
8
42 2 45 36 9 43 4 17
2
9 11
4
7 30 41 15
10
46 3 16 15 46 29 5 17 6 38
5
40 40 46 24 17
8
18 34 49 1 10 3 25 27
4
8 33 47 6
5
48 28 5 25 8
4
28 2 36 28
9
10 5 46 20 10 16 49 4 39
3
17 29 13
6
32 40 46 5 36 26
9
1 32 1 37 24 35 14 28 7
6
20 50 39 20 35 8
9
37 38 12 44 26 43 22 12 16
10
12 47 5 35 4 32 13 5 32 32
7
44 44 34 37 30 33 42
7
44 8 8 5 16 49 6
6
18 10 42 14 46 37
3
48 31 7
2
39 37
2
44 7
2
25 18
4
15 12 5 12
7
19 19 15 9 15 49 18
7
26 29 36 4 5 16 47
6
28 24 9 15 18 50
4
28 18 4 17
3
47 47 7
5
4 18 12 48 42
9
17 2 9 23 42 6 3 37 25
5
25 19 3 4 32
9
39 25 32 23 37 34 17 2 44
3
26 15 27
10
31 8 11 25 7 43 10 2 11 28
9
22 19 3 27 26 2 32 9 22
4
41 46 7 13
5
26 18 38 22 33
4
46 44 26 40
4
11 40 20 42
6
16 4 18 34 8 23
9
6 10 14 1 19 32 6 12 9
10
1 18 17 34 2 23 6 19 33 13
7
15 35 12 38 32 11 33
5
32 5 15 1 50
10
3 1 6 21 11 1 15 39 16 19
2
24 41
5
21 2 46 2 23
8
10 18 50 26 16 40 27 3
5
11 13 9 11 14
10
28 4 43 40 6 31 11 2 17 6
9
1 43 30 29 11 27 20 7 1
2
8 26
5
38 30 45 41 42
10
13 21 26 3 31 1 35 11 25 7
7
29 4 9 50 1 15 20
5
33 9 50 36 45
8
22 39 6 28 2 37 32 7
3
8 16 1
2
43 10
10
21 12 49 24 1 29 36 46 7 18
9
12 39 13 12 47 44 49 5 33
3
9 2 40
7
31 22 34 40 27 17 28
9
8 3 32 37 5 28 20 36 34
5
10 17 42 28 9
7
32 41 35 10 42 44 30
2
37 21
5
8 15 43 9 21
3
21 2 37
7
27 14 49 33 16 31 34
10
38 14 11 30 18 20 43 11 2 25
6
2 10 35 28 12 16
2
23 18
6
48 2 17 45 48 50
8
19 17 21 33 2 30 49 6
3
50 39 15
7
4 46 50 43 32 35 50
2
9 13
9
20 7 44 49 2 49 46 30 37
7
29 10 5 21 30 41 38
10
1 42 11 34 28 3 39 1 39 44
8
8 41 18 2 15 39 50 1
8
37 3 37 29 31 39 35 11
5
44 26 19 38 50
6
28 22 18 31 49 18
2
7 21
3
9 32 19
6
42 15 23 13 48 30
9
15 4 38 24 30 46 23 7 50
9
26 40 6 32 25 35 12 2 5
7
46 10 43 31 26 27 35
4
23 19 41 6
4
20 27 10 17
5
21 40 14 25 20
3
23 8 1
5
10 8 49 44 4
2
20 9
7
20 31 17 14 11 27 6
5
40 43 2 11 18
8
31 31 31 34 1 13 19 10
6
28 21 35 19 23 6
5
35 35 50 9 17
7
41 11 26 28 22 44 12
7
48 39 46 50 35 16 32
5
12 15 41 31 45
5
34 20 31 8 30
6
17 6 22 31 39 11
2
30 23
5
30 36 14 5 28
9
33 22 6 20 22 37 32 32 27
6
26 33 42 35 17 46
10
45 7 19 22 46 35 43 29 34 16
4
26 50 5 17
2
42 20
3
49 42 42
7
34 3 18 5 39 50 24
2
12 33
5
18 50 35 2 49
8
10 40 13 4 46 9 37 15
4
29 15 28 49
6
29 12 1 35 18 23
4
30 6 17 42
6
50 8 43 26 50 50
7
45 30 39 18 32 13 34
2
27 2
6
26 12 42 37 4 14
3
31 30 36
8
18 2 37 46 14 30 34 14
10
25 20 22 6 2 40 45 5 41 42
2
12 48
9
30 42 42 45 30 25 49 50 29
2
35 25
2
22 7
10
19 21 23 49 10 24 33 33 25 14
6
45 43 38 40 12 50
8
18 41 10 14 10 8 23 34
6
17 10 17 31 25 20
3
33 42 21
9
31 37 9 41 5 30 25 18 32
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	pos++
	res := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d missing n", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		pos++
		if pos+n > len(fields) {
			return nil, fmt.Errorf("case %d missing values", i+1)
		}
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			val, err := strconv.ParseInt(fields[pos+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d bad value %d: %v", i+1, j+1, err)
			}
			arr[j] = val
		}
		pos += n
		res = append(res, testCase{n: n, arr: arr})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra tokens at end")
	}
	return res, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	input := sb.String()

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expected := strconv.FormatInt(solve(tc.arr), 10)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput: 1 case with n=%d\nexpected: %s\ngot: %s\n", idx+1, tc.n, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
